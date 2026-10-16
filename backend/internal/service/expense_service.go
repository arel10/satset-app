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

// ExpenseService handles expense management logic.
type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
	auditLogger *audit.Logger
}

// NewExpenseService creates a new ExpenseService.
func NewExpenseService(expenseRepo *repository.ExpenseRepository, auditLogger *audit.Logger) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
		auditLogger: auditLogger,
	}
}

// CreateExpenseRequest represents the data needed to create an expense.
type CreateExpenseRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
	ExpenseDate string `json:"expense_date" binding:"required"` // YYYY-MM-DD
}

// UpdateExpenseRequest represents updatable expense fields.
type UpdateExpenseRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Amount      *int64  `json:"amount" binding:"omitempty,gt=0"`
	Description *string `json:"description"`
	ExpenseDate *string `json:"expense_date"` // YYYY-MM-DD
}

// CreateCategoryRequest represents the data needed to create an expense category.
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateCategoryRequest represents updatable category fields.
type UpdateCategoryRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=50"`
	IsActive *bool   `json:"is_active"`
}

// --- Expense Categories ---

// CreateCategory creates a new expense category.
func (s *ExpenseService) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*model.ExpenseCategory, error) {
	cat := &model.ExpenseCategory{
		Name:     req.Name,
		IsActive: true,
	}
	if err := s.expenseRepo.CreateCategory(ctx, cat); err != nil {
		return nil, fmt.Errorf("gagal membuat kategori: %w", err)
	}
	return cat, nil
}

// GetCategories retrieves expense categories.
func (s *ExpenseService) GetCategories(ctx context.Context, activeOnly bool) ([]model.ExpenseCategory, error) {
	return s.expenseRepo.FindAllCategories(ctx, activeOnly)
}

// UpdateCategory updates an expense category.
func (s *ExpenseService) UpdateCategory(ctx context.Context, id uint, req UpdateCategoryRequest) (*model.ExpenseCategory, error) {
	cat, err := s.expenseRepo.FindCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("kategori tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data kategori: %w", err)
	}

	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}

	if err := s.expenseRepo.UpdateCategory(ctx, cat); err != nil {
		return nil, fmt.Errorf("gagal memperbarui kategori: %w", err)
	}

	return cat, nil
}

// --- Expenses ---

// CreateExpense creates a new expense entry.
func (s *ExpenseService) CreateExpense(ctx context.Context, req CreateExpenseRequest, userID uint, ipAddress string) (*model.Expense, error) {
	// Validate category
	cat, err := s.expenseRepo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("kategori tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa kategori: %w", err)
	}
	if !cat.IsActive {
		return nil, fmt.Errorf("kategori '%s' tidak aktif", cat.Name)
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("nominal harus lebih dari 0")
	}

	// Parse expense date
	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}

	expense := &model.Expense{
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		ExpenseDate: expenseDate.UTC(),
		RecordedBy:  userID,
		IsCancelled: false,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("gagal mencatat pengeluaran: %w", err)
	}

	s.auditLogger.Log(userID, model.AuditActionCreate, model.AuditEntityExpense, expense.ID, nil,
		map[string]interface{}{
			"category_id": req.CategoryID,
			"amount":      req.Amount,
			"date":        req.ExpenseDate,
		}, ipAddress)

	return expense, nil
}

// GetExpense retrieves an expense by ID.
func (s *ExpenseService) GetExpense(ctx context.Context, id uint) (*model.Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengeluaran tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pengeluaran: %w", err)
	}
	return expense, nil
}

// ListExpenses retrieves expenses with filters and pagination.
func (s *ExpenseService) ListExpenses(ctx context.Context, filter repository.ExpenseFilter, page, perPage int) ([]model.Expense, int64, error) {
	offset := (page - 1) * perPage
	return s.expenseRepo.FindAll(ctx, filter, offset, perPage)
}

// UpdateExpense updates an expense.
func (s *ExpenseService) UpdateExpense(ctx context.Context, id uint, req UpdateExpenseRequest, userID uint, ipAddress string) (*model.Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengeluaran tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pengeluaran: %w", err)
	}

	if expense.IsCancelled {
		return nil, fmt.Errorf("pengeluaran yang dibatalkan tidak dapat diubah")
	}

	oldAmount := expense.Amount

	if req.CategoryID != nil {
		cat, err := s.expenseRepo.FindCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("kategori tidak ditemukan")
		}
		if !cat.IsActive {
			return nil, fmt.Errorf("kategori tidak aktif")
		}
		expense.CategoryID = *req.CategoryID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, fmt.Errorf("nominal harus lebih dari 0")
		}
		expense.Amount = *req.Amount
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.ExpenseDate != nil {
		expDate, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal tidak valid")
		}
		expense.ExpenseDate = expDate.UTC()
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengeluaran: %w", err)
	}

	s.auditLogger.Log(userID, model.AuditActionUpdate, model.AuditEntityExpense, expense.ID,
		map[string]interface{}{"amount": oldAmount},
		map[string]interface{}{"amount": expense.Amount}, ipAddress)

	return expense, nil
}

// CancelExpense cancels an expense (Owner only).
func (s *ExpenseService) CancelExpense(ctx context.Context, id uint, reason string, userID uint, ipAddress string) (*model.Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengeluaran tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pengeluaran: %w", err)
	}

	if expense.IsCancelled {
		return nil, fmt.Errorf("pengeluaran sudah dibatalkan")
	}

	now := time.Now().UTC()
	expense.IsCancelled = true
	expense.CancelledBy = &userID
	expense.CancelledAt = &now
	expense.CancelReason = reason

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("gagal membatalkan pengeluaran: %w", err)
	}

	s.auditLogger.Log(userID, model.AuditActionExpenseCancel, model.AuditEntityExpense, expense.ID,
		map[string]interface{}{"amount": expense.Amount},
		map[string]interface{}{"reason": reason}, ipAddress)

	return expense, nil
}

// InitDefaultCategories creates default expense categories if they don't exist.
func (s *ExpenseService) InitDefaultCategories(ctx context.Context) error {
	return s.expenseRepo.InitDefaultCategories(ctx)
}
