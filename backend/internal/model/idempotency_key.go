package model

import (
	"time"
)

// IdempotencyKey prevents duplicate request processing.
type IdempotencyKey struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Key            string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	Endpoint       string    `gorm:"size:255;not null" json:"endpoint"`
	ResponseStatus int       `gorm:"column:response_code;not null" json:"response_status"`
	ResponseBody   string    `gorm:"type:text;not null" json:"response_body"`
	ExpiresAt      time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName specifies the database table name for IdempotencyKey.
func (IdempotencyKey) TableName() string {
	return "idempotency_keys"
}
