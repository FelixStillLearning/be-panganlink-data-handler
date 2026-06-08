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
}

type orderService struct {
	repo    repository.OrderRepository
	payment PaymentService
}

func NewOrderService(r repository.OrderRepository, p PaymentService) OrderService {
	return &orderService{repo: r, payment: p}
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

	// Buat di DB
	if err := s.repo.Create(order); err != nil {
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
	case "cancel", "expire", "deny":
		status = "rejected"
	case "pending":
		status = "pending"
	}

	return s.repo.UpdatePaymentStatus(orderID, status, "midtrans", "", paidAt)
}
