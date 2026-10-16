package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"satset-app/internal/model"
)

// OrderRepository handles database operations for orders.
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// DB returns the underlying gorm.DB for transaction support.
func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

// GenerateOrderNumber generates the next order number for today within a transaction.
// Uses SELECT FOR UPDATE to prevent race conditions.
func (r *OrderRepository) GenerateOrderNumber(ctx context.Context, tx *gorm.DB, orderDate time.Time) (string, error) {
	dateStr := orderDate.Format("2006-01-02")
	dateShort := orderDate.Format("060102") // YYMMDD

	var seq model.OrderSequence
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("date = ?", dateStr).
		First(&seq).Error

	if err == gorm.ErrRecordNotFound {
		// Create new sequence for this date
		seq = model.OrderSequence{
			Date:         dateStr,
			LastSequence: 1,
		}
		if err := tx.WithContext(ctx).Create(&seq).Error; err != nil {
			// Handle race condition: another concurrent transaction created the sequence
			var retrySeq model.OrderSequence
			if errRetry := tx.WithContext(ctx).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("date = ?", dateStr).
				First(&retrySeq).Error; errRetry == nil {
				retrySeq.LastSequence++
				if errSave := tx.WithContext(ctx).Save(&retrySeq).Error; errSave != nil {
					return "", fmt.Errorf("failed to update order sequence on retry: %w", errSave)
				}
				seq = retrySeq
			} else {
				return "", fmt.Errorf("failed to create order sequence: %w", err)
			}
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to lock order sequence: %w", err)
	} else {
		seq.LastSequence++
		if err := tx.WithContext(ctx).Save(&seq).Error; err != nil {
			return "", fmt.Errorf("failed to update order sequence: %w", err)
		}
	}

	return fmt.Sprintf("LDR-%s-%03d", dateShort, seq.LastSequence), nil
}

// Create inserts a new order with its items within a transaction.
func (r *OrderRepository) Create(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}

// FindByID retrieves an order by ID with related data.
func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Preload("Payments").
		Preload("Payments.Receiver").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("StatusHistory.Changer").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByIDForUpdate retrieves an order by ID within a transaction using SELECT ... FOR UPDATE.
func (r *OrderRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.Order, error) {
	var order model.Order
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Preload("Payments").
		Preload("Payments.Receiver").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("StatusHistory.Changer").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByOrderNumber retrieves an order by its order number.
func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Preload("Payments").
		Preload("Payments.Receiver").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("StatusHistory.Changer").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// OrderFilter defines filter criteria for listing orders.
type OrderFilter struct {
	Query         string // search by order_number, customer name, phone
	WorkStatus    string
	PaymentStatus string
	StartDate     string // YYYY-MM-DD
	EndDate       string // YYYY-MM-DD
	Overdue       bool
	CustomerID    uint
}

// FindAll retrieves orders with filters and pagination.
func (r *OrderRepository) FindAll(ctx context.Context, filter OrderFilter, offset, limit int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Order{})

	if filter.Query != "" {
		search := "%" + filter.Query + "%"
		db = db.Joins("LEFT JOIN customers ON customers.id = orders.customer_id").
			Where("orders.order_number LIKE ? OR customers.name LIKE ? OR customers.phone LIKE ?",
				search, search, search)
	}

	if filter.WorkStatus != "" {
		db = db.Where("orders.work_status = ?", filter.WorkStatus)
	}

	if filter.PaymentStatus != "" {
		db = db.Where("orders.payment_status = ?", filter.PaymentStatus)
	}

	if filter.StartDate != "" {
		db = db.Where("orders.order_date >= ?", filter.StartDate)
	}

	if filter.EndDate != "" {
		db = db.Where("orders.order_date < ?", filter.EndDate+" 23:59:59")
	}

	if filter.Overdue {
		db = db.Where("orders.estimated_completion < ? AND orders.work_status NOT IN (?, ?, ?)",
			time.Now().UTC(), model.WorkStatusSiapDiambil, model.WorkStatusSelesai, model.WorkStatusDibatalkan)
	}

	if filter.CustomerID > 0 {
		db = db.Where("orders.customer_id = ?", filter.CustomerID)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Where(db).
		Order("orders.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Update updates an order record.
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateWithTx updates an order record within a transaction.
func (r *OrderRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	return tx.WithContext(ctx).Save(order).Error
}

// CreateStatusHistory inserts a status change record.
func (r *OrderRepository) CreateStatusHistory(ctx context.Context, tx *gorm.DB, history *model.OrderStatusHistory) error {
	return tx.WithContext(ctx).Create(history).Error
}

// CountByStatus returns order counts grouped by work status.
func (r *OrderRepository) CountByStatus(ctx context.Context, startDate, endDate string) ([]StatusCount, error) {
	var results []StatusCount
	db := r.db.WithContext(ctx).Model(&model.Order{})

	if startDate != "" {
		db = db.Where("order_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("order_date < ?", endDate+" 23:59:59")
	}

	err := db.Select("work_status, COUNT(*) as count").
		Group("work_status").
		Scan(&results).Error
	return results, err
}

// CountActive returns the count of active (non-terminal) orders.
func (r *OrderRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Where("work_status NOT IN (?, ?)", model.WorkStatusSelesai, model.WorkStatusDibatalkan).
		Count(&count).Error
	return count, err
}

// CountOverdue returns the count of overdue orders.
func (r *OrderRepository) CountOverdue(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Where("estimated_completion < ? AND work_status NOT IN (?, ?, ?)",
			time.Now().UTC(), model.WorkStatusSiapDiambil, model.WorkStatusSelesai, model.WorkStatusDibatalkan).
		Count(&count).Error
	return count, err
}

// CountReadyForPickup returns the count of orders ready for pickup.
func (r *OrderRepository) CountReadyForPickup(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Where("work_status = ?", model.WorkStatusSiapDiambil).
		Count(&count).Error
	return count, err
}

// CountUnpaid returns the count of unpaid orders.
func (r *OrderRepository) CountUnpaid(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Where("payment_status = ? AND work_status != ?", model.PaymentStatusBelumLunas, model.WorkStatusDibatalkan).
		Count(&count).Error
	return count, err
}

// GetRevenueSummary returns revenue summary for a date range.
func (r *OrderRepository) GetRevenueSummary(ctx context.Context, startDate, endDate string) (*RevenueSummary, error) {
	var result RevenueSummary
	err := r.db.WithContext(ctx).
		Raw(`SELECT 
		     COUNT(*) as total_orders,
		     COALESCE(SUM(CASE WHEN work_status != 'DIBATALKAN' THEN total ELSE 0 END), 0) as total_revenue,
		     COALESCE(SUM(CASE WHEN work_status != 'DIBATALKAN' THEN discount ELSE 0 END), 0) as total_discount
		     FROM orders
		     WHERE order_date >= ? AND order_date < ?`,
			startDate, endDate+" 23:59:59").
		Scan(&result).Error
	return &result, err
}

// GetDailyRevenue returns daily revenue breakdown.
func (r *OrderRepository) GetDailyRevenue(ctx context.Context, startDate, endDate string) ([]DailyRevenue, error) {
	var results []DailyRevenue
	err := r.db.WithContext(ctx).
		Raw(`SELECT DATE(order_date) as date, 
		     COUNT(*) as order_count,
		     COALESCE(SUM(total), 0) as revenue
		     FROM orders
		     WHERE order_date >= ? AND order_date < ?
		     AND work_status != 'DIBATALKAN'
		     GROUP BY DATE(order_date)
		     ORDER BY date ASC`, startDate, endDate+" 23:59:59").
		Scan(&results).Error
	return results, err
}

// FindOverdueOrders returns overdue orders with customer info.
func (r *OrderRepository) FindOverdueOrders(ctx context.Context, limit int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Items").
		Where("estimated_completion < ? AND work_status NOT IN (?, ?, ?)",
			time.Now().UTC(), model.WorkStatusSiapDiambil, model.WorkStatusSelesai, model.WorkStatusDibatalkan).
		Order("estimated_completion ASC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// FindUnpaidOrders returns unpaid orders with customer info.
func (r *OrderRepository) FindUnpaidOrders(ctx context.Context, limit int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Items").
		Where("payment_status = ? AND work_status != ?", model.PaymentStatusBelumLunas, model.WorkStatusDibatalkan).
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// FindByCustomerID retrieves orders for a specific customer.
func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID uint, limit int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// GetTodayOrderCount returns today's order count.
func (r *OrderRepository) GetTodayOrderCount(ctx context.Context, today string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Where("DATE(order_date) = ?", today).
		Count(&count).Error
	return count, err
}

// CheckIdempotencyKey checks if an idempotency key exists and returns it.
func (r *OrderRepository) CheckIdempotencyKey(ctx context.Context, key string) (*model.IdempotencyKey, error) {
	var ik model.IdempotencyKey
	err := r.db.WithContext(ctx).Where("`key` = ? AND expires_at > ?", key, time.Now().UTC()).First(&ik).Error
	if err != nil {
		return nil, err
	}
	return &ik, nil
}

// SaveIdempotencyKey stores an idempotency key.
func (r *OrderRepository) SaveIdempotencyKey(ctx context.Context, tx *gorm.DB, ik *model.IdempotencyKey) error {
	return tx.WithContext(ctx).Create(ik).Error
}

// DeleteExpiredIdempotencyKeys cleans up expired idempotency keys.
func (r *OrderRepository) DeleteExpiredIdempotencyKeys(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&model.IdempotencyKey{}).Error
}

// View models for query results

// StatusCount holds the count for each work status.
type StatusCount struct {
	WorkStatus string `json:"work_status"`
	Count      int64  `json:"count"`
}

// RevenueSummary holds aggregated revenue data.
type RevenueSummary struct {
	TotalOrders   int64 `json:"total_orders"`
	TotalRevenue  int64 `json:"total_revenue"`
	TotalDiscount int64 `json:"total_discount"`
}

// DailyRevenue holds per-day revenue data.
type DailyRevenue struct {
	Date       string `json:"date"`
	OrderCount int64  `json:"order_count"`
	Revenue    int64  `json:"revenue"`
}
