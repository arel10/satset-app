package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"satset-app/internal/audit"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// ServiceService handles laundry service management logic.
type ServiceService struct {
	serviceRepo *repository.ServiceRepository
	auditLogger *audit.Logger
}

// NewServiceService creates a new ServiceService.
func NewServiceService(serviceRepo *repository.ServiceRepository, auditLogger *audit.Logger) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepo,
		auditLogger: auditLogger,
	}
}

// CreateServiceRequest represents the data needed to create a laundry service.
type CreateServiceRequest struct {
	Name              string `json:"name" binding:"required,max=100"`
	Category          string `json:"category" binding:"max=50"`
	Unit              string `json:"unit" binding:"required"`
	PricePerUnit      int64  `json:"price_per_unit" binding:"required,min=0"`
	EstimatedDuration int    `json:"estimated_duration" binding:"min=0"` // hours
	Description       string `json:"description"`
}

// UpdateServiceRequest represents updatable service fields.
type UpdateServiceRequest struct {
	Name              *string `json:"name" binding:"omitempty,max=100"`
	Category          *string `json:"category" binding:"omitempty,max=50"`
	Unit              *string `json:"unit"`
	PricePerUnit      *int64  `json:"price_per_unit" binding:"omitempty,min=0"`
	EstimatedDuration *int    `json:"estimated_duration" binding:"omitempty,min=0"`
	Description       *string `json:"description"`
	IsActive          *bool   `json:"is_active"`
}

// CreateService creates a new laundry service.
func (s *ServiceService) CreateService(ctx context.Context, req CreateServiceRequest, userID uint, ipAddress string) (*model.Service, error) {
	if !model.IsValidUnit(req.Unit) {
		return nil, fmt.Errorf("satuan tidak valid: %s (gunakan: kg, pcs, pasang)", req.Unit)
	}

	if req.PricePerUnit < 0 {
		return nil, fmt.Errorf("harga tidak boleh negatif")
	}

	svc := &model.Service{
		Name:              req.Name,
		Category:          req.Category,
		Unit:              req.Unit,
		PricePerUnit:      req.PricePerUnit,
		EstimatedDuration: req.EstimatedDuration,
		IsActive:          true,
		Description:       req.Description,
	}

	if err := s.serviceRepo.Create(ctx, svc); err != nil {
		return nil, fmt.Errorf("gagal membuat layanan: %w", err)
	}

	s.auditLogger.Log(userID, model.AuditActionCreate, model.AuditEntityService, svc.ID, nil,
		map[string]interface{}{"name": svc.Name, "price": svc.PricePerUnit, "unit": svc.Unit}, ipAddress)

	return svc, nil
}

// GetService retrieves a service by ID.
func (s *ServiceService) GetService(ctx context.Context, id uint) (*model.Service, error) {
	svc, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("layanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data layanan: %w", err)
	}
	return svc, nil
}

// ListServices retrieves services with optional active-only filter and pagination.
func (s *ServiceService) ListServices(ctx context.Context, activeOnly bool, page, perPage int) ([]model.Service, int64, error) {
	offset := (page - 1) * perPage
	return s.serviceRepo.FindAll(ctx, activeOnly, offset, perPage)
}

// GetActiveServices retrieves only active services (for POS).
func (s *ServiceService) GetActiveServices(ctx context.Context) ([]model.Service, error) {
	return s.serviceRepo.FindActive(ctx)
}

// UpdateService updates a service's information.
func (s *ServiceService) UpdateService(ctx context.Context, id uint, req UpdateServiceRequest, userID uint, ipAddress string) (*model.Service, error) {
	svc, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("layanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data layanan: %w", err)
	}

	oldPrice := svc.PricePerUnit
	oldActive := svc.IsActive

	if req.Name != nil {
		svc.Name = *req.Name
	}
	if req.Category != nil {
		svc.Category = *req.Category
	}
	if req.Unit != nil {
		if !model.IsValidUnit(*req.Unit) {
			return nil, fmt.Errorf("satuan tidak valid: %s", *req.Unit)
		}
		svc.Unit = *req.Unit
	}
	if req.PricePerUnit != nil {
		if *req.PricePerUnit < 0 {
			return nil, fmt.Errorf("harga tidak boleh negatif")
		}
		svc.PricePerUnit = *req.PricePerUnit
	}
	if req.EstimatedDuration != nil {
		svc.EstimatedDuration = *req.EstimatedDuration
	}
	if req.Description != nil {
		svc.Description = *req.Description
	}
	if req.IsActive != nil {
		svc.IsActive = *req.IsActive
	}

	if err := s.serviceRepo.Update(ctx, svc); err != nil {
		return nil, fmt.Errorf("gagal memperbarui layanan: %w", err)
	}

	s.auditLogger.Log(userID, model.AuditActionUpdate, model.AuditEntityService, svc.ID,
		map[string]interface{}{"price": oldPrice, "is_active": oldActive},
		map[string]interface{}{"price": svc.PricePerUnit, "is_active": svc.IsActive}, ipAddress)

	return svc, nil
}
