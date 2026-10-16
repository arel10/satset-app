package model

import (
	"time"
)

// ExpenseCategory defines a category for business expenses.
type ExpenseCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the database table name for ExpenseCategory.
func (ExpenseCategory) TableName() string {
	return "expense_categories"
}

// Expense records a business expense entry.
type Expense struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	CategoryID   uint            `gorm:"index;not null" json:"category_id"`
	Category     ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount       int64           `gorm:"not null" json:"amount"`
	Description  string          `gorm:"type:text" json:"description,omitempty"`
	ExpenseDate  time.Time       `gorm:"index;not null" json:"expense_date"`
	RecordedBy   uint            `gorm:"not null" json:"recorded_by"`
	Recorder     User            `gorm:"foreignKey:RecordedBy" json:"recorder,omitempty"`
	IsCancelled  bool            `gorm:"default:false;not null" json:"is_cancelled"`
	CancelledBy  *uint           `json:"cancelled_by,omitempty"`
	CancelledAt  *time.Time      `json:"cancelled_at,omitempty"`
	CancelReason string          `gorm:"type:text" json:"cancel_reason,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// TableName specifies the database table name for Expense.
func (Expense) TableName() string {
	return "expenses"
}

// DefaultExpenseCategories returns the initial expense categories.
func DefaultExpenseCategories() []ExpenseCategory {
	return []ExpenseCategory{
		{Name: "Deterjen"},
		{Name: "Pewangi"},
		{Name: "Listrik"},
		{Name: "Air"},
		{Name: "Gaji"},
		{Name: "Sewa"},
		{Name: "Perawatan Mesin"},
		{Name: "Perlengkapan"},
		{Name: "Lainnya"},
	}
}
