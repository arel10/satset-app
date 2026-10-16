package model

import (
	"time"
)

// OrderStatusHistory records each status change for an order.
type OrderStatusHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"index;not null" json:"order_id"`
	FromStatus string    `gorm:"size:20;not null" json:"from_status"`
	ToStatus   string    `gorm:"size:20;not null" json:"to_status"`
	ChangedBy  uint      `gorm:"not null" json:"changed_by"`
	Changer    User      `gorm:"foreignKey:ChangedBy" json:"changer,omitempty"`
	Notes      string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the database table name for OrderStatusHistory.
func (OrderStatusHistory) TableName() string {
	return "order_status_history"
}
