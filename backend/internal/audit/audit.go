package audit

import (
	"encoding/json"
	"log"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// Logger provides methods to create audit log entries.
type Logger struct {
	db *gorm.DB
}

// NewLogger creates a new audit Logger.
func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

// Log creates a new audit log entry.
func (l *Logger) Log(userID uint, action, entityType string, entityID uint, oldValue, newValue interface{}, ipAddress string) {
	entry := model.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
	}

	if oldValue != nil {
		if data, err := json.Marshal(oldValue); err == nil {
			s := string(data)
			entry.OldValue = &s
		}
	}

	if newValue != nil {
		if data, err := json.Marshal(newValue); err == nil {
			s := string(data)
			entry.NewValue = &s
		}
	}

	if err := l.db.Create(&entry).Error; err != nil {
		log.Printf("[AUDIT] Failed to create audit log: %v", err)
	}
}

// LogWithTx creates a new audit log entry within an existing transaction.
func (l *Logger) LogWithTx(tx *gorm.DB, userID uint, action, entityType string, entityID uint, oldValue, newValue interface{}, ipAddress string) error {
	entry := model.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
	}

	if oldValue != nil {
		if data, err := json.Marshal(oldValue); err == nil {
			s := string(data)
			entry.OldValue = &s
		}
	}

	if newValue != nil {
		if data, err := json.Marshal(newValue); err == nil {
			s := string(data)
			entry.NewValue = &s
		}
	}

	return tx.Create(&entry).Error
}
