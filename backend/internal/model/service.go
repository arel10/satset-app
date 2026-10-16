package model

import (
	"time"
)

// Service represents a laundry service offering.
type Service struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:100;not null" json:"name"`
	Category          string    `gorm:"size:50" json:"category,omitempty"`
	Unit              string    `gorm:"size:20;not null" json:"unit"`
	PricePerUnit      int64     `gorm:"not null" json:"price_per_unit"`
	EstimatedDuration int       `gorm:"default:0" json:"estimated_duration"` // in hours
	IsActive          bool      `gorm:"default:true;not null" json:"is_active"`
	Description       string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName specifies the database table name for Service.
func (Service) TableName() string {
	return "services"
}

// Service unit types
const (
	UnitKg     = "kg"
	UnitPcs    = "pcs"
	UnitPasang = "pasang"
)

// ValidUnits is a list of valid service units.
var ValidUnits = []string{UnitKg, UnitPcs, UnitPasang}

// IsValidUnit checks whether a unit string is valid.
func IsValidUnit(unit string) bool {
	for _, u := range ValidUnits {
		if u == unit {
			return true
		}
	}
	return false
}
