package model

import (
	"time"
)

// OrderItem represents one line item in an order with snapshot pricing.
type OrderItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	ServiceID    uint      `gorm:"not null" json:"service_id"`
	ServiceName  string    `gorm:"size:100;not null" json:"service_name"`
	ServiceUnit  string    `gorm:"size:20;not null" json:"service_unit"`
	PricePerUnit int64     `gorm:"not null" json:"price_per_unit"`
	Quantity     float64   `gorm:"type:decimal(10,3);not null" json:"quantity"`
	Subtotal     int64     `gorm:"not null" json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName specifies the database table name for OrderItem.
func (OrderItem) TableName() string {
	return "order_items"
}
