package repository

import (
	"time"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id string) (*model.Order, error)
	FindByBuyerID(buyerID string) ([]model.Order, error)
	FindByPetaniID(petaniID string) ([]model.Order, error)
	UpdatePaymentStatus(orderID, status, method, reference string, paidAt *time.Time) error
	UpdateStatus(orderID, status string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Items").Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByBuyerID(buyerID string) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Preload("Items").Where("buyer_id = ?", buyerID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByPetaniID(petaniID string) ([]model.Order, error) {
	var orders []model.Order
	// Need to join order_items and products to filter by product.user_id = petaniID
	query := r.db.Distinct("orders.*").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.user_id = ?", petaniID).
		Preload("Items")

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdatePaymentStatus(orderID, status, method, reference string, paidAt *time.Time) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"status":            status,
		"payment_method":    method,
		"payment_reference": reference,
		"paid_at":           paidAt,
	}).Error
}

func (r *orderRepository) UpdateStatus(orderID, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
