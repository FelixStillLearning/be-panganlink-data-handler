package model

import "time"

type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36);default:(UUID())"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);index"`
	Title     string    `json:"title" gorm:"type:varchar(100)"`
	Message   string    `json:"message" gorm:"type:text"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
