package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// CustomerRepository handles database operations for customers.
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create inserts a new customer.
func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

// FindByID retrieves a customer by ID.
func (r *CustomerRepository) FindByID(ctx context.Context, id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// FindByPhone retrieves a customer by normalized phone number.
func (r *CustomerRepository) FindByPhone(ctx context.Context, phone string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Search searches customers by name or phone with pagination.
func (r *CustomerRepository) Search(ctx context.Context, query string, offset, limit int) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Customer{})
	if query != "" {
		search := "%" + query + "%"
		db = db.Where("name LIKE ? OR phone LIKE ?", search, search)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// Update updates a customer record.
func (r *CustomerRepository) Update(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

// CountAll returns the total number of customers.
func (r *CustomerRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Customer{}).Count(&count).Error
	return count, err
}

// FindTopCustomers returns customers with the most orders.
func (r *CustomerRepository) FindTopCustomers(ctx context.Context, limit int) ([]CustomerOrderCount, error) {
	var results []CustomerOrderCount
	err := r.db.WithContext(ctx).
		Raw(`SELECT c.id, c.name, c.phone, COUNT(o.id) as order_count, 
		     COALESCE(SUM(CASE WHEN o.work_status != 'DIBATALKAN' THEN o.total ELSE 0 END), 0) as total_spent
		     FROM customers c 
		     LEFT JOIN orders o ON c.id = o.customer_id 
		     GROUP BY c.id, c.name, c.phone 
		     ORDER BY order_count DESC 
		     LIMIT ?`, limit).
		Scan(&results).Error
	return results, err
}

// CustomerOrderCount is a view model for customer with order statistics.
type CustomerOrderCount struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	OrderCount int    `json:"order_count"`
	TotalSpent int64  `json:"total_spent"`
}
