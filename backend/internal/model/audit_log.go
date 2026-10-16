package model

import (
	"time"
)

// AuditLog records important actions for accountability and traceability.
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string    `gorm:"size:50;not null;index" json:"action"`
	EntityType string    `gorm:"size:50;not null;index" json:"entity_type"`
	EntityID   uint      `gorm:"index" json:"entity_id"`
	OldValue   *string   `gorm:"type:json" json:"old_value,omitempty"`
	NewValue   *string   `gorm:"type:json" json:"new_value,omitempty"`
	IPAddress  string    `gorm:"size:45" json:"ip_address,omitempty"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// TableName specifies the database table name for AuditLog.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Audit actions
const (
	AuditActionCreate        = "CREATE"
	AuditActionUpdate        = "UPDATE"
	AuditActionDelete        = "DELETE"
	AuditActionLogin         = "LOGIN"
	AuditActionLoginFailed   = "LOGIN_FAILED"
	AuditActionLogout        = "LOGOUT"
	AuditActionStatusChange  = "STATUS_CHANGE"
	AuditActionPayment       = "PAYMENT"
	AuditActionPaymentCancel = "PAYMENT_CANCEL"
	AuditActionOrderCancel   = "ORDER_CANCEL"
	AuditActionExpenseCancel = "EXPENSE_CANCEL"
	AuditActionDiscount      = "DISCOUNT"
	AuditActionAIChat        = "AI_CHAT"
)

// Audit entity types
const (
	AuditEntityUser           = "user"
	AuditEntityOrder          = "order"
	AuditEntityPayment        = "payment"
	AuditEntityExpense        = "expense"
	AuditEntityService        = "service"
	AuditEntityCustomer       = "customer"
	AuditEntitySetting        = "setting"
	AuditEntityAIConversation = "ai_conversation"
)
