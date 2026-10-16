package model

import (
	"time"
)

// Customer stores laundry customer information.
type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;index" json:"name"`
	Phone     string    `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Address   string    `gorm:"type:text" json:"address,omitempty"`
	Notes     string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the database table name for Customer.
func (Customer) TableName() string {
	return "customers"
}
