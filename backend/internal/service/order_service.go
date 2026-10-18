package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"gorm.io/gorm"

	"satset-app/internal/audit"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// OrderService handles order/POS logic.
type OrderService struct {
	orderRepo    *repository.OrderRepository
	serviceRepo  *repository.ServiceRepository
	customerRepo *repository.CustomerRepository
	paymentRepo  *repository.PaymentRepository
	auditLogger  *audit.Logger
}

// NewOrderService creates a new OrderService.
func NewOrderService(
	orderRepo *repository.OrderRepository,
	serviceRepo *repository.ServiceRepository,
	customerRepo *repository.CustomerRepository,
	paymentRepo *repository.PaymentRepository,
	auditLogger *audit.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		serviceRepo:  serviceRepo,
		customerRepo: customerRepo,
		paymentRepo:  paymentRepo,
		auditLogger:  auditLogger,
	}
}

// CreateOrderRequest represents the POS order creation request.
type CreateOrderRequest struct {
	CustomerID          uint               `json:"customer_id" binding:"required"`
	Items               []OrderItemRequest `json:"items" binding:"required,min=1"`
	Discount            int64              `json:"discount" binding:"min=0"`
	Notes               string             `json:"notes"`
	EstimatedCompletion *time.Time         `json:"estimated_completion"`
}

// OrderItemRequest represents one line item in the order.
type OrderItemRequest struct {
	ServiceID uint    `json:"service_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
}

// CreateOrder creates a new order with items in a single transaction.
// Supports idempotency via idempotency key.
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest, userID uint, ipAddress string, idempotencyKey string) (*model.Order, error) {
	// Check idempotency key
	if idempotencyKey != "" {
		existing, err := s.orderRepo.CheckIdempotencyKey(ctx, idempotencyKey)
		if err == nil && existing != nil {
			// Return cached response
			var order model.Order
			if err := json.Unmarshal([]byte(existing.ResponseBody), &order); err == nil {
				return &order, nil
			}
		}
	}

	// Validate customer exists
	customer, err := s.customerRepo.FindByID(ctx, req.CustomerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pelanggan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa pelanggan: %w", err)
	}
	_ = customer

	// Validate and calculate items
	var items []model.OrderItem
	var subtotal int64
	var maxDuration int

	for i, itemReq := range req.Items {
		svc, err := s.serviceRepo.FindByID(ctx, itemReq.ServiceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("layanan item ke-%d tidak ditemukan", i+1)
			}
			return nil, fmt.Errorf("gagal memeriksa layanan: %w", err)
		}

		if !svc.IsActive {
			return nil, fmt.Errorf("layanan '%s' tidak aktif", svc.Name)
		}

		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("jumlah item ke-%d harus lebih dari 0", i+1)
		}

		// Calculate item subtotal: quantity * price_per_unit (rounded)
		itemSubtotal := int64(math.Round(itemReq.Quantity * float64(svc.PricePerUnit)))

		items = append(items, model.OrderItem{
			ServiceID:    svc.ID,
			ServiceName:  svc.Name,
			ServiceUnit:  svc.Unit,
			PricePerUnit: svc.PricePerUnit,
			Quantity:     itemReq.Quantity,
			Subtotal:     itemSubtotal,
		})

		subtotal += itemSubtotal
		if svc.EstimatedDuration > maxDuration {
			maxDuration = svc.EstimatedDuration
		}
	}

	// Validate discount
	if req.Discount < 0 {
		return nil, fmt.Errorf("diskon tidak boleh negatif")
	}
	if req.Discount > subtotal {
		return nil, fmt.Errorf("diskon tidak boleh melebihi subtotal")
	}

	total := subtotal - req.Discount

	// Calculate estimated completion
	var estimatedCompletion time.Time
	if req.EstimatedCompletion != nil {
		estimatedCompletion = req.EstimatedCompletion.UTC()
	} else {
		if maxDuration == 0 {
			maxDuration = 24 // default 24 hours
		}
		estimatedCompletion = time.Now().UTC().Add(time.Duration(maxDuration) * time.Hour)
	}

	// Create order in transaction
	var order *model.Order

	err = s.orderRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		// Generate order number (with lock)
		orderNumber, err := s.orderRepo.GenerateOrderNumber(ctx, tx, now)
		if err != nil {
			return fmt.Errorf("gagal membuat nomor pesanan: %w", err)
		}

		order = &model.Order{
			OrderNumber:         orderNumber,
			CustomerID:          req.CustomerID,
			OrderDate:           now,
			EstimatedCompletion: estimatedCompletion,
			Subtotal:            subtotal,
			Discount:            req.Discount,
			Total:               total,
			Notes:               req.Notes,
			WorkStatus:          model.WorkStatusDiterima,
			PaymentStatus:       model.PaymentStatusBelumLunas,
			CreatedBy:           userID,
		}

		if err := s.orderRepo.Create(ctx, tx, order); err != nil {
			return fmt.Errorf("gagal membuat pesanan: %w", err)
		}

		// Create order items
		for i := range items {
			items[i].OrderID = order.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return fmt.Errorf("gagal membuat item pesanan: %w", err)
		}

		// Create initial status history
		history := &model.OrderStatusHistory{
			OrderID:    order.ID,
			FromStatus: "",
			ToStatus:   model.WorkStatusDiterima,
			ChangedBy:  userID,
			Notes:      "Pesanan baru dibuat",
		}
		if err := s.orderRepo.CreateStatusHistory(ctx, tx, history); err != nil {
			return fmt.Errorf("gagal mencatat riwayat status: %w", err)
		}

		// Save idempotency key
		if idempotencyKey != "" {
			orderJSON, _ := json.Marshal(order)
			ik := &model.IdempotencyKey{
				Key:            idempotencyKey,
				UserID:         userID,
				Endpoint:       "/api/orders",
				ResponseStatus: http.StatusCreated,
				ResponseBody:   string(orderJSON),
				ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
			}
			if err := s.orderRepo.SaveIdempotencyKey(ctx, tx, ik); err != nil {
				return fmt.Errorf("gagal menyimpan idempotency key: %w", err)
			}
		}

		// Audit log
		return s.auditLogger.LogWithTx(tx, userID, model.AuditActionCreate, model.AuditEntityOrder, order.ID, nil,
			map[string]interface{}{
				"order_number": order.OrderNumber,
				"customer_id":  order.CustomerID,
				"total":        order.Total,
				"item_count":   len(items),
			}, ipAddress)
	})

	if err != nil {
		return nil, err
	}

	order.Items = items

	// Log discount if applied
	if req.Discount > 0 {
		s.auditLogger.Log(userID, model.AuditActionDiscount, model.AuditEntityOrder, order.ID, nil,
			map[string]interface{}{"discount": req.Discount, "order_number": order.OrderNumber}, ipAddress)
	}

	return order, nil
}

// GetOrder retrieves an order by ID with full details.
func (s *OrderService) GetOrder(ctx context.Context, id uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pesanan: %w", err)
	}
	return order, nil
}

// GetOrderByNumber retrieves an order by its order number.
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	order, err := s.orderRepo.FindByOrderNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pesanan: %w", err)
	}
	return order, nil
}

// ListOrders retrieves orders with filters and pagination.
func (s *OrderService) ListOrders(ctx context.Context, filter repository.OrderFilter, page, perPage int) ([]model.Order, int64, error) {
	offset := (page - 1) * perPage
	return s.orderRepo.FindAll(ctx, filter, offset, perPage)
}

// UpdateWorkStatusRequest represents a status change request.
type UpdateWorkStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// UpdateWorkStatus changes the work status of an order.
func (s *OrderService) UpdateWorkStatus(ctx context.Context, orderID uint, req UpdateWorkStatusRequest, userID uint, userRole string, ipAddress string) (*model.Order, error) {
	var order *model.Order
	isOwner := userRole == model.RoleOwner

	err := s.orderRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("pesanan tidak ditemukan")
			}
			return fmt.Errorf("gagal mengambil data pesanan: %w", err)
		}

		if !model.IsValidTransition(order.WorkStatus, req.Status, isOwner) {
			return fmt.Errorf("perubahan status dari %s ke %s tidak diizinkan", order.WorkStatus, req.Status)
		}

		oldStatus := order.WorkStatus
		order.WorkStatus = req.Status

		if err := s.orderRepo.UpdateWithTx(ctx, tx, order); err != nil {
			return fmt.Errorf("gagal memperbarui status: %w", err)
		}

		// Record status history
		history := &model.OrderStatusHistory{
			OrderID:    order.ID,
			FromStatus: oldStatus,
			ToStatus:   req.Status,
			ChangedBy:  userID,
			Notes:      req.Notes,
		}
		if err := s.orderRepo.CreateStatusHistory(ctx, tx, history); err != nil {
			return fmt.Errorf("gagal mencatat riwayat status: %w", err)
		}

		return s.auditLogger.LogWithTx(tx, userID, model.AuditActionStatusChange, model.AuditEntityOrder, order.ID,
			map[string]string{"status": oldStatus},
			map[string]string{"status": req.Status, "notes": req.Notes}, ipAddress)
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// CancelOrderRequest represents a cancellation request.
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// CancelOrder cancels an order.
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, req CancelOrderRequest, userID uint, userRole string, ipAddress string) (*model.Order, error) {
	var order *model.Order
	isOwner := userRole == model.RoleOwner
	now := time.Now().UTC()

	err := s.orderRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("pesanan tidak ditemukan")
			}
			return fmt.Errorf("gagal mengambil data pesanan: %w", err)
		}

		if !model.IsValidTransition(order.WorkStatus, model.WorkStatusDibatalkan, isOwner) {
			return fmt.Errorf("pesanan dengan status %s tidak dapat dibatalkan", order.WorkStatus)
		}

		oldStatus := order.WorkStatus
		order.WorkStatus = model.WorkStatusDibatalkan
		order.CancelledBy = &userID
		order.CancelledAt = &now
		order.CancelReason = req.Reason

		if err := s.orderRepo.UpdateWithTx(ctx, tx, order); err != nil {
			return fmt.Errorf("gagal membatalkan pesanan: %w", err)
		}

		history := &model.OrderStatusHistory{
			OrderID:    order.ID,
			FromStatus: oldStatus,
			ToStatus:   model.WorkStatusDibatalkan,
			ChangedBy:  userID,
			Notes:      "Dibatalkan: " + req.Reason,
		}
		if err := s.orderRepo.CreateStatusHistory(ctx, tx, history); err != nil {
			return fmt.Errorf("gagal mencatat riwayat: %w", err)
		}

		return s.auditLogger.LogWithTx(tx, userID, model.AuditActionOrderCancel, model.AuditEntityOrder, order.ID,
			map[string]string{"status": oldStatus},
			map[string]interface{}{
				"status":       model.WorkStatusDibatalkan,
				"reason":       req.Reason,
				"order_number": order.OrderNumber,
			}, ipAddress)
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrdersByCustomer retrieves orders for a specific customer.
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uint, limit int) ([]model.Order, error) {
	return s.orderRepo.FindByCustomerID(ctx, customerID, limit)
}
