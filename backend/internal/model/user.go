package model

import (
	"time"
)

// User represents a system user (Owner, Kasir, Operator).
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	Role         string    `gorm:"size:20;not null;index" json:"role"`
	IsActive     bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the database table name for User.
func (User) TableName() string {
	return "users"
}

// User roles
const (
	RoleOwner    = "owner"
	RoleKasir    = "kasir"
	RoleOperator = "operator"
)

// ValidRoles is a list of valid user roles.
var ValidRoles = []string{RoleOwner, RoleKasir, RoleOperator}

// IsValidRole checks whether a role string is valid.
func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}

// RefreshToken stores active refresh tokens for session management.
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Token     string    `gorm:"uniqueIndex;size:255;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the database table name for RefreshToken.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
