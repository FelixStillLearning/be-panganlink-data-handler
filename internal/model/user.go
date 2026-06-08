package model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // Not exposed in JSON
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // petani, pembeli, admin
	Location  string    `gorm:"type:varchar(100)" json:"location"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
