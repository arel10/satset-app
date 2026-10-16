package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// PaymentRepository handles database operations for payments.
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PaymentRepository.
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create inserts a new payment record within a transaction.
func (r *PaymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *model.Payment) error {
	return tx.WithContext(ctx).Create(payment).Error
}

// FindByID retrieves a payment by ID.
func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).
		Preload("Receiver").
		First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByOrderID retrieves all payments for an order.
func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uint) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).
		Preload("Receiver").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&payments).Error
	return payments, err
}

// FindRecent retrieves recent payments with order and customer info.
func (r *PaymentRepository) FindRecent(ctx context.Context, limit int) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).
		Preload("Receiver").
		Preload("Order").
		Preload("Order.Customer").
		Where("is_cancelled = ?", false).
		Order("paid_at DESC").
		Limit(limit).
		Find(&payments).Error
	return payments, err
}

// Update updates a payment record.
func (r *PaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// UpdateWithTx updates a payment record within a transaction.
func (r *PaymentRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, payment *model.Payment) error {
	return tx.WithContext(ctx).Save(payment).Error
}

// GetValidPaymentTotal returns the sum of valid (non-cancelled) payments for an order.
func (r *PaymentRepository) GetValidPaymentTotal(ctx context.Context, orderID uint) (int64, error) {
	var total *int64
	err := r.db.WithContext(ctx).Model(&model.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("order_id = ? AND is_cancelled = ?", orderID, false).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

// GetPaymentSummary returns payment summary within a date range.
func (r *PaymentRepository) GetPaymentSummary(ctx context.Context, startDate, endDate string) (*PaymentSummary, error) {
	var result PaymentSummary
	err := r.db.WithContext(ctx).
		Raw(`SELECT 
		     COUNT(*) as total_payments,
		     COALESCE(SUM(CASE WHEN is_cancelled = false THEN amount ELSE 0 END), 0) as total_received,
		     COALESCE(SUM(CASE WHEN is_cancelled = false AND payment_method = 'TUNAI' THEN amount ELSE 0 END), 0) as cash_received,
		     COALESCE(SUM(CASE WHEN is_cancelled = false AND payment_method = 'TRANSFER' THEN amount ELSE 0 END), 0) as transfer_received,
		     COALESCE(SUM(CASE WHEN is_cancelled = false AND payment_method = 'QRIS' THEN amount ELSE 0 END), 0) as qris_received,
		     COALESCE(SUM(CASE WHEN is_cancelled = false AND payment_method = 'LAINNYA' THEN amount ELSE 0 END), 0) as other_received
		     FROM payments
		     WHERE paid_at >= ? AND paid_at < ?`, startDate, endDate+" 23:59:59").
		Scan(&result).Error
	return &result, err
}

// GetDailyPayments returns daily payment breakdown.
func (r *PaymentRepository) GetDailyPayments(ctx context.Context, startDate, endDate string) ([]DailyPayment, error) {
	var results []DailyPayment
	err := r.db.WithContext(ctx).
		Raw(`SELECT DATE(paid_at) as date,
		     COUNT(*) as payment_count,
		     COALESCE(SUM(amount), 0) as total_amount
		     FROM payments
		     WHERE paid_at >= ? AND paid_at < ?
		     AND is_cancelled = false
		     GROUP BY DATE(paid_at)
		     ORDER BY date ASC`, startDate, endDate+" 23:59:59").
		Scan(&results).Error
	return results, err
}

// GetValidPaymentTotalWithTx is GetValidPaymentTotal but within a transaction.
func (r *PaymentRepository) GetValidPaymentTotalWithTx(ctx context.Context, tx *gorm.DB, orderID uint) (int64, error) {
	var total *int64
	err := tx.WithContext(ctx).Model(&model.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("order_id = ? AND is_cancelled = ?", orderID, false).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

// PaymentSummary holds aggregated payment data.
type PaymentSummary struct {
	TotalPayments    int64 `json:"total_payments"`
	TotalReceived    int64 `json:"total_received"`
	CashReceived     int64 `json:"cash_received"`
	TransferReceived int64 `json:"transfer_received"`
	QRISReceived     int64 `json:"qris_received"`
	OtherReceived    int64 `json:"other_received"`
}

// DailyPayment holds per-day payment data.
type DailyPayment struct {
	Date         string `json:"date"`
	PaymentCount int64  `json:"payment_count"`
	TotalAmount  int64  `json:"total_amount"`
}
