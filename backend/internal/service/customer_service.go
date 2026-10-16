package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// CustomerService handles customer management logic.
type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

// CreateCustomerRequest represents the data needed to create a customer.
type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address"`
	Notes   string `json:"notes"`
}

// UpdateCustomerRequest represents updatable customer fields.
type UpdateCustomerRequest struct {
	Name    *string `json:"name" binding:"omitempty,max=100"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
	Notes   *string `json:"notes"`
}

// CreateCustomer creates a new customer with phone normalization.
func (s *CustomerService) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*model.Customer, error) {
	phone := NormalizePhone(req.Phone)
	if phone == "" {
		return nil, fmt.Errorf("nomor telepon tidak valid")
	}

	// Check for existing phone
	_, err := s.customerRepo.FindByPhone(ctx, phone)
	if err == nil {
		return nil, fmt.Errorf("nomor telepon sudah terdaftar")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("gagal memeriksa nomor telepon: %w", err)
	}

	customer := &model.Customer{
		Name:    strings.TrimSpace(req.Name),
		Phone:   phone,
		Address: strings.TrimSpace(req.Address),
		Notes:   strings.TrimSpace(req.Notes),
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("gagal membuat pelanggan: %w", err)
	}

	return customer, nil
}

// GetCustomer retrieves a customer by ID.
func (s *CustomerService) GetCustomer(ctx context.Context, id uint) (*model.Customer, error) {
	customer, err := s.customerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pelanggan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}
	return customer, nil
}

// SearchCustomers searches customers by name or phone.
func (s *CustomerService) SearchCustomers(ctx context.Context, query string, page, perPage int) ([]model.Customer, int64, error) {
	offset := (page - 1) * perPage
	return s.customerRepo.Search(ctx, query, offset, perPage)
}

// UpdateCustomer updates customer information.
func (s *CustomerService) UpdateCustomer(ctx context.Context, id uint, req UpdateCustomerRequest) (*model.Customer, error) {
	customer, err := s.customerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pelanggan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}

	if req.Name != nil {
		customer.Name = strings.TrimSpace(*req.Name)
	}

	if req.Phone != nil {
		phone := NormalizePhone(*req.Phone)
		if phone == "" {
			return nil, fmt.Errorf("nomor telepon tidak valid")
		}
		// Check uniqueness only if phone changed
		if phone != customer.Phone {
			existing, err := s.customerRepo.FindByPhone(ctx, phone)
			if err == nil && existing.ID != customer.ID {
				return nil, fmt.Errorf("nomor telepon sudah digunakan pelanggan lain")
			}
			customer.Phone = phone
		}
	}

	if req.Address != nil {
		customer.Address = strings.TrimSpace(*req.Address)
	}

	if req.Notes != nil {
		customer.Notes = strings.TrimSpace(*req.Notes)
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pelanggan: %w", err)
	}

	return customer, nil
}

// NormalizePhone normalizes an Indonesian phone number.
// Converts 08xx to +628xx, handles +62 prefix, strips non-digits.
func NormalizePhone(phone string) string {
	// Remove all non-digit characters except leading +
	re := regexp.MustCompile(`[^\d+]`)
	phone = re.ReplaceAllString(strings.TrimSpace(phone), "")

	if phone == "" {
		return ""
	}

	// Convert various formats to +62
	if strings.HasPrefix(phone, "+62") {
		// Already in international format
	} else if strings.HasPrefix(phone, "62") {
		phone = "+" + phone
	} else if strings.HasPrefix(phone, "0") {
		phone = "+62" + phone[1:]
	} else {
		// Assume local number, add +62
		phone = "+62" + phone
	}

	// Basic validation: should be at least 10 digits after +62
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) < 10 || len(digits) > 15 {
		return ""
	}

	return phone
}
