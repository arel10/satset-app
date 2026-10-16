package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"satset-app/internal/audit"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// UserService handles user management logic.
type UserService struct {
	userRepo    *repository.UserRepository
	auditLogger *audit.Logger
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repository.UserRepository, auditLogger *audit.Logger) *UserService {
	return &UserService{
		userRepo:    userRepo,
		auditLogger: auditLogger,
	}
}

// CreateUserRequest represents the data needed to create a user.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,max=100"`
	Role     string `json:"role" binding:"required"`
}

// UpdateUserRequest represents the data that can be updated for a user.
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	FullName string `json:"full_name" binding:"omitempty,max=100"`
	Role     string `json:"role" binding:"omitempty"`
	IsActive *bool  `json:"is_active"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// CreateUser creates a new user account. Only Owner can do this.
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest, creatorID uint, ipAddress string) (*model.User, error) {
	if !model.IsValidRole(req.Role) {
		return nil, fmt.Errorf("peran tidak valid: %s", req.Role)
	}

	// Check if username already exists
	_, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username sudah digunakan")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("gagal memeriksa username: %w", err)
	}

	// Check if email already exists
	_, err = s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email sudah digunakan")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("gagal memeriksa email: %w", err)
	}

	// Hash password
	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal membuat pengguna: %w", err)
	}

	s.auditLogger.Log(creatorID, model.AuditActionCreate, model.AuditEntityUser, user.ID, nil,
		map[string]string{"username": user.Username, "role": user.Role}, ipAddress)

	return user, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengguna tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}
	return user, nil
}

// ListUsers retrieves all users with pagination.
func (s *UserService) ListUsers(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	offset := (page - 1) * perPage
	return s.userRepo.FindAll(ctx, offset, perPage)
}

// UpdateUser updates user information. Only Owner can do this.
func (s *UserService) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest, updaterID uint, ipAddress string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengguna tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}

	// Prevent deactivating own account
	if updaterID == id && req.IsActive != nil && !*req.IsActive {
		return nil, fmt.Errorf("anda tidak dapat menonaktifkan akun anda sendiri")
	}

	// Prevent demoting or deactivating the last active owner
	if user.Role == model.RoleOwner && user.IsActive {
		isDeactivating := req.IsActive != nil && !*req.IsActive
		isDemoting := req.Role != "" && req.Role != model.RoleOwner
		if isDeactivating || isDemoting {
			activeOwners, err := s.userRepo.CountActiveByRole(ctx, model.RoleOwner)
			if err != nil {
				return nil, fmt.Errorf("gagal memeriksa jumlah owner aktif: %w", err)
			}
			if activeOwners <= 1 {
				if isDeactivating {
					return nil, fmt.Errorf("tidak dapat menonaktifkan owner terakhir yang aktif")
				}
				return nil, fmt.Errorf("tidak dapat mengubah peran karena merupakan satu-satunya owner aktif")
			}
		}
	}

	oldRole := user.Role
	oldActive := user.IsActive

	if req.Email != "" && req.Email != user.Email {
		existing, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existing.ID != user.ID {
			return nil, fmt.Errorf("email sudah digunakan")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Role != "" {
		if !model.IsValidRole(req.Role) {
			return nil, fmt.Errorf("peran tidak valid: %s", req.Role)
		}
		user.Role = req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
		// If deactivating, remove all refresh tokens
		if !*req.IsActive {
			_ = s.userRepo.DeleteUserRefreshTokens(ctx, user.ID)
		}
	}

	if req.Password != "" {
		hash, err := HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
		}
		user.PasswordHash = hash
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengguna: %w", err)
	}

	s.auditLogger.Log(updaterID, model.AuditActionUpdate, model.AuditEntityUser, user.ID,
		map[string]interface{}{"role": oldRole, "is_active": oldActive},
		map[string]interface{}{"role": user.Role, "is_active": user.IsActive}, ipAddress)

	return user, nil
}

// SeedDefaultOwner creates the default owner account if no owner exists.
func (s *UserService) SeedDefaultOwner(ctx context.Context) error {
	return s.SeedDefaultUsers(ctx)
}

// SeedDefaultUsers creates default accounts for owner, kasir, and operator if they do not exist.
func (s *UserService) SeedDefaultUsers(ctx context.Context) error {
	// 1. Owner
	ownerCount, err := s.userRepo.CountByRole(ctx, model.RoleOwner)
	if err == nil && ownerCount == 0 {
		hash, err := HashPassword("owner123")
		if err == nil {
			_ = s.userRepo.Create(ctx, &model.User{
				Username:     "owner",
				Email:        "owner@satset.local",
				PasswordHash: hash,
				FullName:     "Pemilik SatSet",
				Role:         model.RoleOwner,
				IsActive:     true,
			})
		}
	}

	// 2. Kasir
	kasirCount, err := s.userRepo.CountByRole(ctx, model.RoleKasir)
	if err == nil && kasirCount == 0 {
		hash, err := HashPassword("kasir123")
		if err == nil {
			_ = s.userRepo.Create(ctx, &model.User{
				Username:     "kasir",
				Email:        "kasir@satset.local",
				PasswordHash: hash,
				FullName:     "Kasir SatSet",
				Role:         model.RoleKasir,
				IsActive:     true,
			})
		}
	}

	// 3. Operator
	operatorCount, err := s.userRepo.CountByRole(ctx, model.RoleOperator)
	if err == nil && operatorCount == 0 {
		hash, err := HashPassword("operator123")
		if err == nil {
			_ = s.userRepo.Create(ctx, &model.User{
				Username:     "operator",
				Email:        "operator@satset.local",
				PasswordHash: hash,
				FullName:     "Operator SatSet",
				Role:         model.RoleOperator,
				IsActive:     true,
			})
		}
	}

	return nil
}

// DeleteUser soft deletes a user and revokes refresh tokens.
func (s *UserService) DeleteUser(ctx context.Context, id uint, actorID uint, ipAddress string) error {
	if actorID == id {
		return fmt.Errorf("anda tidak dapat menghapus akun anda sendiri")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pengguna tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}

	// Prevent deleting the last active owner
	if user.Role == model.RoleOwner && user.IsActive {
		activeOwners, err := s.userRepo.CountActiveByRole(ctx, model.RoleOwner)
		if err != nil {
			return fmt.Errorf("gagal memeriksa jumlah owner aktif: %w", err)
		}
		if activeOwners <= 1 {
			return fmt.Errorf("tidak dapat menghapus owner terakhir yang aktif")
		}
	}

	user.IsActive = false

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("gagal menghapus pengguna: %w", err)
	}

	// Revoke all refresh tokens
	_ = s.userRepo.DeleteUserRefreshTokens(ctx, user.ID)

	s.auditLogger.Log(actorID, model.AuditActionDelete, model.AuditEntityUser, user.ID,
		map[string]interface{}{"username": user.Username, "role": user.Role}, nil, ipAddress)

	return nil
}
