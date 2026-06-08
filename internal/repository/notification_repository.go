package repository

import (
	"github.com/example/be-panganlink-data-handler/internal/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notif *model.Notification) error
	FindByUserID(userID string) ([]model.Notification, error)
	MarkAsRead(id string, userID string) error
	MarkAllAsRead(userID string) error
}

type notificationRepository struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) NotificationRepository { return &notificationRepository{db} }

func (r *notificationRepository) Create(n *model.Notification) error { return r.db.Create(n).Error }

func (r *notificationRepository) FindByUserID(userID string) ([]model.Notification, error) {
	var notifs []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) MarkAsRead(id, userID string) error {
	return r.db.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID string) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}
