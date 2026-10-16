package model

import (
	"time"
)

// Setting stores key-value business settings.
type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the database table name for Setting.
func (Setting) TableName() string {
	return "settings"
}

// Predefined setting keys
const (
	SettingBusinessName    = "business_name"
	SettingBusinessAddress = "business_address"
	SettingBusinessPhone   = "business_phone"
	SettingBusinessLogo    = "business_logo"
	SettingReceiptFooter   = "receipt_footer"
	SettingCurrency        = "currency"
	SettingTimezone        = "timezone"
)

// DefaultSettings returns the initial settings for a new installation.
func DefaultSettings() []Setting {
	return []Setting{
		{Key: SettingBusinessName, Value: "SatSet Laundry"},
		{Key: SettingBusinessAddress, Value: ""},
		{Key: SettingBusinessPhone, Value: ""},
		{Key: SettingBusinessLogo, Value: ""},
		{Key: SettingReceiptFooter, Value: "Terima kasih atas kepercayaan Anda!"},
		{Key: SettingCurrency, Value: "IDR"},
		{Key: SettingTimezone, Value: "Asia/Jakarta"},
	}
}
