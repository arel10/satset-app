package service

import (
	"context"

	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// AuditService handles business logic for audit logs.
type AuditService struct {
	auditRepo *repository.AuditRepository
}

// NewAuditService creates a new AuditService instance.
func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

// ListAuditLogs retrieves audit logs with filter and pagination.
func (s *AuditService) ListAuditLogs(ctx context.Context, filter repository.AuditFilter, page, perPage int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage
	return s.auditRepo.FindAll(ctx, filter, offset, perPage)
}
