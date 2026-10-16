package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"satset-app/internal/audit"
	"satset-app/internal/config"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// AuthService handles authentication logic.
type AuthService struct {
	userRepo    *repository.UserRepository
	cfg         *config.Config
	auditLogger *audit.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config, auditLogger *audit.Logger) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		cfg:         cfg,
		auditLogger: auditLogger,
	}
}

// LoginRequest represents login credentials.
type LoginRequest struct {
	Identifier string `json:"identifier"` // username or email
	Username   string `json:"username"`   // alternative field for username
	Password   string `json:"password" binding:"required"`
}

// TokenResponse represents the JWT token pair response.
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"` // seconds
	User         UserInfo `json:"user"`
}

// UserInfo is a safe user representation for API responses.
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// AuthClaims are the JWT claims for access tokens.
type AuthClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns JWT tokens.
func (s *AuthService) Login(ctx context.Context, req LoginRequest, ipAddress string) (*TokenResponse, error) {
	targetIdentifier := req.Identifier
	if targetIdentifier == "" {
		targetIdentifier = req.Username
	}
	if targetIdentifier == "" {
		return nil, fmt.Errorf("username atau email diperlukan")
	}

	user, err := s.userRepo.FindByUsernameOrEmail(ctx, targetIdentifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.auditLogger.Log(0, model.AuditActionLoginFailed, model.AuditEntityUser, 0, nil,
				map[string]string{"identifier": targetIdentifier, "reason": "user not found"}, ipAddress)
			return nil, fmt.Errorf("username/email atau password salah")
		}
		return nil, fmt.Errorf("gagal memproses login: %w", err)
	}

	if !user.IsActive {
		s.auditLogger.Log(user.ID, model.AuditActionLoginFailed, model.AuditEntityUser, user.ID, nil,
			map[string]string{"reason": "account inactive"}, ipAddress)
		return nil, fmt.Errorf("akun telah dinonaktifkan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.auditLogger.Log(user.ID, model.AuditActionLoginFailed, model.AuditEntityUser, user.ID, nil,
			map[string]string{"reason": "wrong password"}, ipAddress)
		return nil, fmt.Errorf("username/email atau password salah")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token: %w", err)
	}

	// Store refresh token
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpiryDays) * 24 * time.Hour)
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.userRepo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("gagal menyimpan refresh token: %w", err)
	}

	s.auditLogger.Log(user.ID, model.AuditActionLogin, model.AuditEntityUser, user.ID, nil, nil, ipAddress)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.AccessExpiryMinutes * 60,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token.
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (*TokenResponse, error) {
	// Validate the refresh token JWT
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWT.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("refresh token tidak valid")
	}

	// Verify token exists in database
	rt, err := s.userRepo.FindRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("refresh token tidak ditemukan")
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token sudah kedaluwarsa")
	}

	// Get current user data
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("pengguna tidak ditemukan")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("akun telah dinonaktifkan")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token baru: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr, // reuse same refresh token
		ExpiresIn:    s.cfg.JWT.AccessExpiryMinutes * 60,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}

// Logout invalidates a refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string, userID uint, ipAddress string) error {
	if err := s.userRepo.DeleteRefreshToken(ctx, refreshTokenStr); err != nil {
		return fmt.Errorf("gagal menghapus sesi: %w", err)
	}
	s.auditLogger.Log(userID, model.AuditActionLogout, model.AuditEntityUser, userID, nil, nil, ipAddress)
	return nil
}

// ValidateAccessToken parses and validates an access token, returning its claims.
func (s *AuthService) ValidateAccessToken(tokenStr string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWT.AccessSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}
	return claims, nil
}

// HashPassword generates a bcrypt hash for a password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// generateAccessToken creates a short-lived JWT access token.
func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.JWT.AccessExpiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "satset",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.AccessSecret))
}

// generateRefreshToken creates a long-lived JWT refresh token.
func (s *AuthService) generateRefreshToken(user *model.User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpiryDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "satset",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.RefreshSecret))
}
