package model

import "time"

type Komoditas struct {
	ID       string `gorm:"type:varchar(20);primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(50);not null" json:"nama"`
	Satuan   string `gorm:"type:varchar(20);not null" json:"satuan"`
	Kategori string `gorm:"type:varchar(50)" json:"kategori"`
}

func (Komoditas) TableName() string {
	return "komoditas"
}

type Product struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID      string    `gorm:"type:varchar(36)" json:"user_id"`
	KomoditasID string    `gorm:"type:varchar(20)" json:"komoditas_id"`
	Harga       float64   `gorm:"type:decimal(12,2);not null" json:"harga"`
	Stok        float64   `gorm:"type:decimal(10,2);not null" json:"stok"`
	FotoUrl     string    `gorm:"type:text" json:"foto_url"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Komoditas Komoditas `gorm:"foreignKey:KomoditasID" json:"komoditas"`
	User      User      `gorm:"foreignKey:UserID" json:"petani"`
}
