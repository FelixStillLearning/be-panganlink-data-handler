package service

import (
	"fmt"
	"time"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
)

type OrderService interface {
	Checkout(buyerID string, items []model.OrderItem) (*model.Order, error)
	GetBuyerOrders(buyerID string) ([]model.Order, error)
	GetPetaniOrders(petaniID string) ([]model.Order, error)
	UpdateOrderStatus(orderID, status string) error
	HandleMidtransWebhook(payload map[string]interface{}) error
	CancelExpiredOrders() error
}

type orderService struct {
	repo      repository.OrderRepository
	payment   PaymentService
	notifRepo repository.NotificationRepository
}

func NewOrderService(r repository.OrderRepository, p PaymentService, nr repository.NotificationRepository) OrderService {
	return &orderService{repo: r, payment: p, notifRepo: nr}
}

func (s *orderService) Checkout(buyerID string, items []model.OrderItem) (*model.Order, error) {
	var totalHarga float64
	for _, item := range items {
		totalHarga += item.Jumlah * item.HargaUnit
	}

	orderID := fmt.Sprintf("ORD-%d", time.Now().Unix())

	order := &model.Order{
		ID:         orderID,
		BuyerID:    buyerID,
		TotalHarga: totalHarga,
		Status:     "pending",
		Items:      items,
	}

	// Buat di DB dengan Transaksi (Lock Stok & ACID)
	if err := s.repo.CreateCheckoutTransaction(order); err != nil {
		return nil, err
	}

	// Dapatkan Snap Token dari Midtrans
	token, err := s.payment.InitiatePayment(orderID, totalHarga, "PanganLink Buyer", "buyer@example.com")
	if err != nil {
		return nil, err
	}

	order.PaymentToken = token
	s.repo.UpdatePaymentStatus(orderID, "pending", "", "", nil)

	return order, nil
}

func (s *orderService) GetBuyerOrders(buyerID string) ([]model.Order, error) {
	return s.repo.FindByBuyerID(buyerID)
}

func (s *orderService) GetPetaniOrders(petaniID string) ([]model.Order, error) {
	return s.repo.FindByPetaniID(petaniID)
}

func (s *orderService) UpdateOrderStatus(orderID, status string) error {
	return s.repo.UpdateStatus(orderID, status)
}

func (s *orderService) HandleMidtransWebhook(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)

	if !s.payment.VerifyPaymentSignature(orderID, statusCode, grossAmount, signatureKey) {
		return fmt.Errorf("invalid signature")
	}

	var status string
	var paidAt *time.Time

	switch transactionStatus {
	case "capture", "settlement":
		status = "paid"
		now := time.Now()
		paidAt = &now

		// Create notification for Buyer
		order, err := s.repo.FindByID(orderID)
		if err == nil && order != nil {
			s.notifRepo.Create(&model.Notification{
				UserID:  order.BuyerID,
				Title:   "Pembayaran Berhasil",
				Message: fmt.Sprintf("Pembayaran untuk pesanan %s telah berhasil dikonfirmasi.", orderID),
			})

			// Notify Petani (from the first item as simplified approach)
			if len(order.Items) > 0 {
				// Ideally we fetch product to get PetaniID, but this needs ProductRepository
				// For now, we only notify Buyer to avoid circular dependencies
			}
		}
	case "cancel", "expire", "deny":
		status = "rejected"
	case "pending":
		status = "pending"
	}

	return s.repo.UpdatePaymentStatus(orderID, status, "midtrans", "", paidAt)
}

func (s *orderService) CancelExpiredOrders() error {
	// Let's say expiration time is 24 hours
	expiryTime := time.Now().Add(-24 * time.Hour)
	orders, err := s.repo.FindExpiredOrders(expiryTime)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := s.repo.CancelExpiredOrderTransaction(&order); err != nil {
			// Log error but continue to next order
			fmt.Printf("Failed to cancel expired order %s: %v\n", order.ID, err)
		} else {
			// Notify buyer that order is canceled
			s.notifRepo.Create(&model.Notification{
				UserID:  order.BuyerID,
				Title:   "Pesanan Dibatalkan Otomatis",
				Message: fmt.Sprintf("Pesanan %s telah dibatalkan karena melewati batas waktu pembayaran 24 jam.", order.ID),
			})
		}
	}
	return nil
}
