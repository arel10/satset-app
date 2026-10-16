package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// AuditRepository handles database operations for audit logs.
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new AuditRepository.
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// FindAll retrieves audit logs with filters and pagination.
func (r *AuditRepository) FindAll(ctx context.Context, filter AuditFilter, offset, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if filter.UserID > 0 {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		db = db.Where("action = ?", filter.Action)
	}
	if filter.EntityType != "" {
		db = db.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID > 0 {
		db = db.Where("entity_id = ?", filter.EntityID)
	}
	if filter.StartDate != "" {
		db = db.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		db = db.Where("created_at < ?", filter.EndDate+" 23:59:59")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("User").
		Where(db).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// AuditFilter defines filter criteria for listing audit logs.
type AuditFilter struct {
	UserID     uint
	Action     string
	EntityType string
	EntityID   uint
	StartDate  string
	EndDate    string
}
