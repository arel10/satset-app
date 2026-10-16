package model

import (
	"time"
)

// Payment records a payment made against an order.
type Payment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	OrderID       uint       `gorm:"index;not null" json:"order_id"`
	Order         *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Amount        int64      `gorm:"not null" json:"amount"`
	PaymentMethod string     `gorm:"size:20;not null" json:"payment_method"`
	PaidAt        time.Time  `gorm:"index;not null" json:"paid_at"`
	ReceivedBy    uint       `gorm:"not null" json:"received_by"`
	Receiver      User       `gorm:"foreignKey:ReceivedBy" json:"receiver,omitempty"`
	Notes         string     `gorm:"type:text" json:"notes,omitempty"`
	IsCancelled   bool       `gorm:"default:false;not null" json:"is_cancelled"`
	CancelledBy   *uint      `json:"cancelled_by,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CancelReason  string     `gorm:"type:text" json:"cancel_reason,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName specifies the database table name for Payment.
func (Payment) TableName() string {
	return "payments"
}

// Payment methods
const (
	PaymentMethodTunai    = "TUNAI"
	PaymentMethodTransfer = "TRANSFER"
	PaymentMethodQRIS     = "QRIS"
	PaymentMethodLainnya  = "LAINNYA"
)

// ValidPaymentMethods is a list of valid payment methods.
var ValidPaymentMethods = []string{
	PaymentMethodTunai,
	PaymentMethodTransfer,
	PaymentMethodQRIS,
	PaymentMethodLainnya,
}

// IsValidPaymentMethod checks whether a payment method string is valid.
func IsValidPaymentMethod(method string) bool {
	for _, m := range ValidPaymentMethods {
		if m == method {
			return true
		}
	}
	return false
}
