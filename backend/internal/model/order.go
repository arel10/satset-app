package model

import (
	"time"
)

// Order represents a laundry order/transaction.
type Order struct {
	ID                  uint                 `gorm:"primaryKey" json:"id"`
	OrderNumber         string               `gorm:"uniqueIndex;size:20;not null" json:"order_number"`
	CustomerID          uint                 `gorm:"index;not null" json:"customer_id"`
	Customer            Customer             `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	OrderDate           time.Time            `gorm:"index;not null" json:"order_date"`
	EstimatedCompletion time.Time            `json:"estimated_completion"`
	Subtotal            int64                `gorm:"not null" json:"subtotal"`
	Discount            int64                `gorm:"default:0" json:"discount"`
	Total               int64                `gorm:"not null" json:"total"`
	Notes               string               `gorm:"type:text" json:"notes,omitempty"`
	WorkStatus          string               `gorm:"size:20;not null;index;default:DITERIMA" json:"work_status"`
	PaymentStatus       string               `gorm:"size:20;not null;index;default:BELUM_LUNAS" json:"payment_status"`
	CreatedBy           uint                 `gorm:"not null" json:"created_by"`
	Creator             User                 `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CancelledBy         *uint                `json:"cancelled_by,omitempty"`
	CancelledAt         *time.Time           `json:"cancelled_at,omitempty"`
	CancelReason        string               `gorm:"type:text" json:"cancel_reason,omitempty"`
	Items               []OrderItem          `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payments            []Payment            `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	StatusHistory       []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"status_history,omitempty"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

// TableName specifies the database table name for Order.
func (Order) TableName() string {
	return "orders"
}

// Work statuses
const (
	WorkStatusDiterima    = "DITERIMA"
	WorkStatusDicuci      = "DICUCI"
	WorkStatusDikeringkan = "DIKERINGKAN"
	WorkStatusDisetrika   = "DISETRIKA"
	WorkStatusDikemas     = "DIKEMAS"
	WorkStatusSiapDiambil = "SIAP_DIAMBIL"
	WorkStatusSelesai     = "SELESAI"
	WorkStatusDibatalkan  = "DIBATALKAN"
)

// Payment statuses
const (
	PaymentStatusBelumLunas = "BELUM_LUNAS"
	PaymentStatusLunas      = "LUNAS"
)

// ValidWorkStatuses is the ordered list of work statuses in the standard flow.
var ValidWorkStatuses = []string{
	WorkStatusDiterima,
	WorkStatusDicuci,
	WorkStatusDikeringkan,
	WorkStatusDisetrika,
	WorkStatusDikemas,
	WorkStatusSiapDiambil,
	WorkStatusSelesai,
}

// AllowedTransitions defines the valid state transitions for work status.
// Each key maps to the list of statuses it can transition to.
var AllowedTransitions = map[string][]string{
	WorkStatusDiterima:    {WorkStatusDicuci, WorkStatusDikeringkan, WorkStatusDisetrika, WorkStatusDikemas},
	WorkStatusDicuci:      {WorkStatusDikeringkan, WorkStatusDisetrika, WorkStatusDikemas},
	WorkStatusDikeringkan: {WorkStatusDisetrika, WorkStatusDikemas},
	WorkStatusDisetrika:   {WorkStatusDikemas},
	WorkStatusDikemas:     {WorkStatusSiapDiambil},
	WorkStatusSiapDiambil: {WorkStatusSelesai},
	// WorkStatusSelesai and WorkStatusDibatalkan are terminal states
}

// IsValidTransition checks if a work status transition is allowed.
// Owner can override certain transitions via the isOwner flag.
func IsValidTransition(from, to string, isOwner bool) bool {
	// Owner can cancel any non-terminal order
	if to == WorkStatusDibatalkan && isOwner {
		return from != WorkStatusDibatalkan && from != WorkStatusSelesai
	}

	// Kasir can also cancel orders that are still DITERIMA
	if to == WorkStatusDibatalkan && from == WorkStatusDiterima {
		return true
	}

	allowed, exists := AllowedTransitions[from]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	// Owner can force certain transitions
	if isOwner {
		return isValidWorkStatus(to) && to != WorkStatusDibatalkan
	}

	return false
}

// isValidWorkStatus checks if a status string is a valid work status.
func isValidWorkStatus(status string) bool {
	allStatuses := append(ValidWorkStatuses, WorkStatusDibatalkan)
	for _, s := range allStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsTerminalStatus returns true if the status is a terminal state.
func IsTerminalStatus(status string) bool {
	return status == WorkStatusSelesai || status == WorkStatusDibatalkan
}

// IsActiveOrder returns true if the order is not cancelled or completed.
func IsActiveOrder(status string) bool {
	return !IsTerminalStatus(status)
}
