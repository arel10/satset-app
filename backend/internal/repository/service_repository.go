package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// ServiceRepository handles database operations for laundry services.
type ServiceRepository struct {
	db *gorm.DB
}

// NewServiceRepository creates a new ServiceRepository.
func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create inserts a new service.
func (r *ServiceRepository) Create(ctx context.Context, service *model.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

// FindByID retrieves a service by ID.
func (r *ServiceRepository) FindByID(ctx context.Context, id uint) (*model.Service, error) {
	var service model.Service
	err := r.db.WithContext(ctx).First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// FindAll retrieves all services with optional active-only filter and pagination.
func (r *ServiceRepository) FindAll(ctx context.Context, activeOnly bool, offset, limit int) ([]model.Service, int64, error) {
	var services []model.Service
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Service{})
	if activeOnly {
		db = db.Where("is_active = ?", true)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

// FindActive retrieves only active services (for POS order creation).
func (r *ServiceRepository) FindActive(ctx context.Context) ([]model.Service, error) {
	var services []model.Service
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&services).Error
	return services, err
}

// Update updates a service record.
func (r *ServiceRepository) Update(ctx context.Context, service *model.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
}

// FindTopServices returns the most ordered services within a date range.
func (r *ServiceRepository) FindTopServices(ctx context.Context, startDate, endDate string, limit int) ([]TopServiceResult, error) {
	var results []TopServiceResult
	err := r.db.WithContext(ctx).
		Raw(`SELECT oi.service_id, oi.service_name, oi.service_unit,
		     COUNT(*) as order_count,
		     SUM(oi.quantity) as total_quantity,
		     SUM(oi.subtotal) as total_revenue
		     FROM order_items oi
		     JOIN orders o ON oi.order_id = o.id
		     WHERE o.order_date >= ? AND o.order_date < ?
		     AND o.work_status != 'DIBATALKAN'
		     GROUP BY oi.service_id, oi.service_name, oi.service_unit
		     ORDER BY order_count DESC
		     LIMIT ?`, startDate, endDate, limit).
		Scan(&results).Error
	return results, err
}

// TopServiceResult is a view model for service ordering statistics.
type TopServiceResult struct {
	ServiceID     uint    `json:"service_id"`
	ServiceName   string  `json:"service_name"`
	ServiceUnit   string  `json:"service_unit"`
	OrderCount    int     `json:"order_count"`
	TotalQuantity float64 `json:"total_quantity"`
	TotalRevenue  int64   `json:"total_revenue"`
}
