package model

import "time"

type Order struct {
	ID               string      `json:"id" gorm:"primaryKey;type:varchar(20)"`
	BuyerID          string      `json:"buyer_id" gorm:"type:varchar(36)"`
	TotalHarga       float64     `json:"total_harga" gorm:"type:decimal(14,2);not null"`
	Status           string      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	PaymentToken     string      `json:"payment_token" gorm:"type:varchar(255)"`
	PaymentMethod    string      `json:"payment_method" gorm:"type:varchar(50)"`
	PaymentReference string      `json:"payment_reference" gorm:"type:varchar(100)"`
	PaidAt           *time.Time  `json:"paid_at"`
	CreatedAt        time.Time   `json:"created_at" gorm:"autoCreateTime"`

	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Buyer *User       `json:"buyer,omitempty" gorm:"foreignKey:BuyerID"`
}

type OrderItem struct {
	ID         string   `json:"id" gorm:"primaryKey;type:varchar(36);default:(UUID())"`
	OrderID    string   `json:"order_id" gorm:"type:varchar(20)"`
	ProductID  string   `json:"product_id" gorm:"type:varchar(36)"`
	Jumlah     float64  `json:"jumlah" gorm:"type:decimal(10,2);not null"`
	HargaUnit  float64  `json:"harga_unit" gorm:"type:decimal(12,2);not null"`
	
	Product    *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
