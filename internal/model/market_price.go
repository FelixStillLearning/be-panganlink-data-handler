package model

import "time"

type MarketPrice struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Date      time.Time `gorm:"type:date;not null" json:"date"`
	Komoditas string    `gorm:"type:varchar(50);not null" json:"item"`
	Price     float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Region    string    `gorm:"type:varchar(100);default:'Nasional'" json:"region"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
