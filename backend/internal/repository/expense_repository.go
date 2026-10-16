package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// ExpenseRepository handles database operations for expenses and expense categories.
type ExpenseRepository struct {
	db *gorm.DB
}

// NewExpenseRepository creates a new ExpenseRepository.
func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// --- Expense Categories ---

// CreateCategory inserts a new expense category.
func (r *ExpenseRepository) CreateCategory(ctx context.Context, cat *model.ExpenseCategory) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

// FindCategoryByID retrieves an expense category by ID.
func (r *ExpenseRepository) FindCategoryByID(ctx context.Context, id uint) (*model.ExpenseCategory, error) {
	var cat model.ExpenseCategory
	err := r.db.WithContext(ctx).First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// FindAllCategories retrieves all expense categories.
func (r *ExpenseRepository) FindAllCategories(ctx context.Context, activeOnly bool) ([]model.ExpenseCategory, error) {
	var categories []model.ExpenseCategory
	db := r.db.WithContext(ctx)
	if activeOnly {
		db = db.Where("is_active = ?", true)
	}
	err := db.Order("name ASC").Find(&categories).Error
	return categories, err
}

// UpdateCategory updates an expense category.
func (r *ExpenseRepository) UpdateCategory(ctx context.Context, cat *model.ExpenseCategory) error {
	return r.db.WithContext(ctx).Save(cat).Error
}

// InitDefaultCategories creates default expense categories if they don't exist.
func (r *ExpenseRepository) InitDefaultCategories(ctx context.Context) error {
	defaults := model.DefaultExpenseCategories()
	for _, cat := range defaults {
		var existing model.ExpenseCategory
		err := r.db.WithContext(ctx).Where("name = ?", cat.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(&cat).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

// --- Expenses ---

// Create inserts a new expense.
func (r *ExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

// FindByID retrieves an expense by ID.
func (r *ExpenseRepository) FindByID(ctx context.Context, id uint) (*model.Expense, error) {
	var expense model.Expense
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Recorder").
		First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// FindAll retrieves expenses with filters and pagination.
func (r *ExpenseRepository) FindAll(ctx context.Context, filter ExpenseFilter, offset, limit int) ([]model.Expense, int64, error) {
	var expenses []model.Expense
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Expense{})

	if filter.CategoryID > 0 {
		db = db.Where("category_id = ?", filter.CategoryID)
	}
	if filter.StartDate != "" {
		db = db.Where("expense_date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		db = db.Where("expense_date < ?", filter.EndDate+" 23:59:59")
	}
	if !filter.IncludeCancelled {
		db = db.Where("is_cancelled = ?", false)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("Recorder").
		Where(db).
		Order("expense_date DESC").
		Offset(offset).Limit(limit).
		Find(&expenses).Error
	if err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}

// Update updates an expense record.
func (r *ExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

// GetExpenseSummary returns expense totals for a date range.
func (r *ExpenseRepository) GetExpenseSummary(ctx context.Context, startDate, endDate string) (int64, error) {
	var total *int64
	err := r.db.WithContext(ctx).Model(&model.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("expense_date >= ? AND expense_date < ? AND is_cancelled = ?",
			startDate, endDate+" 23:59:59", false).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

// GetExpenseByCategory returns expense totals grouped by category.
func (r *ExpenseRepository) GetExpenseByCategory(ctx context.Context, startDate, endDate string) ([]ExpenseByCategoryResult, error) {
	var results []ExpenseByCategoryResult
	err := r.db.WithContext(ctx).
		Raw(`SELECT ec.id as category_id, ec.name as category_name,
		     COALESCE(SUM(e.amount), 0) as total_amount,
		     COUNT(e.id) as expense_count
		     FROM expense_categories ec
		     LEFT JOIN expenses e ON ec.id = e.category_id 
		       AND e.expense_date >= ? AND e.expense_date < ?
		       AND e.is_cancelled = false
		     WHERE ec.is_active = true
		     GROUP BY ec.id, ec.name
		     ORDER BY total_amount DESC`, startDate, endDate+" 23:59:59").
		Scan(&results).Error
	return results, err
}

// GetDailyExpenses returns daily expense breakdown.
func (r *ExpenseRepository) GetDailyExpenses(ctx context.Context, startDate, endDate string) ([]DailyExpense, error) {
	var results []DailyExpense
	err := r.db.WithContext(ctx).
		Raw(`SELECT DATE(expense_date) as date,
		     COUNT(*) as expense_count,
		     COALESCE(SUM(amount), 0) as total_amount
		     FROM expenses
		     WHERE expense_date >= ? AND expense_date < ?
		     AND is_cancelled = false
		     GROUP BY DATE(expense_date)
		     ORDER BY date ASC`, startDate, endDate+" 23:59:59").
		Scan(&results).Error
	return results, err
}

// ExpenseFilter defines filter criteria for listing expenses.
type ExpenseFilter struct {
	CategoryID       uint
	StartDate        string
	EndDate          string
	IncludeCancelled bool
}

// ExpenseByCategoryResult holds expense totals grouped by category.
type ExpenseByCategoryResult struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	TotalAmount  int64  `json:"total_amount"`
	ExpenseCount int64  `json:"expense_count"`
}

// DailyExpense holds per-day expense data.
type DailyExpense struct {
	Date         string `json:"date"`
	ExpenseCount int64  `json:"expense_count"`
	TotalAmount  int64  `json:"total_amount"`
}
