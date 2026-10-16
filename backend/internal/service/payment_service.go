package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"satset-app/internal/audit"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// PaymentService handles payment logic.
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
	auditLogger *audit.Logger
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	orderRepo *repository.OrderRepository,
	auditLogger *audit.Logger,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		auditLogger: auditLogger,
	}
}

// RecordPaymentRequest represents a payment recording request.
type RecordPaymentRequest struct {
	OrderID       uint   `json:"order_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Notes         string `json:"notes"`
}

// RecordPayment records a full payment for an order.
func (s *PaymentService) RecordPayment(ctx context.Context, req RecordPaymentRequest, userID uint, ipAddress string) (*model.Payment, error) {
	if !model.IsValidPaymentMethod(req.PaymentMethod) {
		return nil, fmt.Errorf("metode pembayaran tidak valid: %s", req.PaymentMethod)
	}

	var payment *model.Payment

	err := s.orderRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Retrieve and lock order
		order, err := s.orderRepo.FindByIDForUpdate(ctx, tx, req.OrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("pesanan tidak ditemukan")
			}
			return fmt.Errorf("gagal mengambil data pesanan: %w", err)
		}

		if order.WorkStatus == model.WorkStatusDibatalkan {
			return fmt.Errorf("tidak dapat mencatat pembayaran untuk pesanan yang dibatalkan")
		}

		if order.PaymentStatus == model.PaymentStatusLunas {
			return fmt.Errorf("pesanan sudah lunas")
		}

		// Get current valid payment total
		currentPaid, err := s.paymentRepo.GetValidPaymentTotalWithTx(ctx, tx, req.OrderID)
		if err != nil {
			return fmt.Errorf("gagal menghitung pembayaran: %w", err)
		}

		remaining := order.Total - currentPaid
		if remaining <= 0 {
			return fmt.Errorf("pesanan sudah lunas")
		}

		// Payment must cover the remaining amount (full payment policy)
		if req.Amount != remaining {
			return fmt.Errorf("jumlah pembayaran harus Rp %d (sisa tagihan)", remaining)
		}

		now := time.Now().UTC()
		payment = &model.Payment{
			OrderID:       req.OrderID,
			Amount:        req.Amount,
			PaymentMethod: req.PaymentMethod,
			PaidAt:        now,
			ReceivedBy:    userID,
			Notes:         req.Notes,
			IsCancelled:   false,
		}

		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
			return fmt.Errorf("gagal mencatat pembayaran: %w", err)
		}

		// Update order payment status
		newPaid := currentPaid + req.Amount
		if newPaid >= order.Total {
			order.PaymentStatus = model.PaymentStatusLunas
			if err := s.orderRepo.UpdateWithTx(ctx, tx, order); err != nil {
				return fmt.Errorf("gagal memperbarui status pembayaran: %w", err)
			}
		}

		return s.auditLogger.LogWithTx(tx, userID, model.AuditActionPayment, model.AuditEntityPayment, payment.ID, nil,
			map[string]interface{}{
				"order_id":       order.ID,
				"order_number":   order.OrderNumber,
				"amount":         req.Amount,
				"payment_method": req.PaymentMethod,
				"payment_status": order.PaymentStatus,
			}, ipAddress)
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// CancelPaymentRequest represents a payment cancellation request.
type CancelPaymentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// CancelPayment cancels a payment (Owner only).
func (s *PaymentService) CancelPayment(ctx context.Context, paymentID uint, req CancelPaymentRequest, userID uint, ipAddress string) (*model.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pembayaran tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pembayaran: %w", err)
	}

	if payment.IsCancelled {
		return nil, fmt.Errorf("pembayaran sudah dibatalkan sebelumnya")
	}

	err = s.orderRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		payment.IsCancelled = true
		payment.CancelledBy = &userID
		payment.CancelledAt = &now
		payment.CancelReason = req.Reason

		if err := s.paymentRepo.UpdateWithTx(ctx, tx, payment); err != nil {
			return fmt.Errorf("gagal membatalkan pembayaran: %w", err)
		}

		// Recalculate order payment status with locked order
		validTotal, err := s.paymentRepo.GetValidPaymentTotalWithTx(ctx, tx, payment.OrderID)
		if err != nil {
			return fmt.Errorf("gagal menghitung ulang pembayaran: %w", err)
		}

		order, err := s.orderRepo.FindByIDForUpdate(ctx, tx, payment.OrderID)
		if err != nil {
			return fmt.Errorf("gagal mengambil data pesanan: %w", err)
		}

		if validTotal < order.Total {
			order.PaymentStatus = model.PaymentStatusBelumLunas
		}
		if err := s.orderRepo.UpdateWithTx(ctx, tx, order); err != nil {
			return fmt.Errorf("gagal memperbarui status pembayaran: %w", err)
		}

		return s.auditLogger.LogWithTx(tx, userID, model.AuditActionPaymentCancel, model.AuditEntityPayment, payment.ID,
			map[string]interface{}{"amount": payment.Amount, "method": payment.PaymentMethod},
			map[string]interface{}{"reason": req.Reason, "order_id": payment.OrderID}, ipAddress)
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentsByOrder retrieves all payments for an order.
func (s *PaymentService) GetPaymentsByOrder(ctx context.Context, orderID uint) ([]model.Payment, error) {
	return s.paymentRepo.FindByOrderID(ctx, orderID)
}
