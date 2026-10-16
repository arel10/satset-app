package service

import (
	"context"
	"fmt"

	"satset-app/internal/audit"
	"satset-app/internal/repository"
)

// SettingService handles business settings logic.
type SettingService struct {
	settingRepo *repository.SettingRepository
	auditLogger *audit.Logger
}

// NewSettingService creates a new SettingService.
func NewSettingService(settingRepo *repository.SettingRepository, auditLogger *audit.Logger) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
		auditLogger: auditLogger,
	}
}

// UpdateSettingsRequest represents a batch update of settings.
type UpdateSettingsRequest struct {
	BusinessName    *string `json:"business_name"`
	BusinessAddress *string `json:"business_address"`
	BusinessPhone   *string `json:"business_phone"`
	ReceiptFooter   *string `json:"receipt_footer"`
	BusinessLogoURL *string `json:"business_logo_url"`
}

// GetAll retrieves all settings as a key-value map.
func (s *SettingService) GetAll(ctx context.Context) (map[string]string, error) {
	return s.settingRepo.GetAllAsMap(ctx)
}

// Update updates business settings.
func (s *SettingService) Update(ctx context.Context, req UpdateSettingsRequest, userID uint, ipAddress string) error {
	// Get old values for audit
	oldSettings, err := s.settingRepo.GetAllAsMap(ctx)
	if err != nil {
		return fmt.Errorf("gagal mengambil pengaturan: %w", err)
	}

	updates := make(map[string]string)
	if req.BusinessName != nil {
		updates["business_name"] = *req.BusinessName
	}
	if req.BusinessAddress != nil {
		updates["business_address"] = *req.BusinessAddress
	}
	if req.BusinessPhone != nil {
		updates["business_phone"] = *req.BusinessPhone
	}
	if req.ReceiptFooter != nil {
		updates["receipt_footer"] = *req.ReceiptFooter
	}
	if req.BusinessLogoURL != nil {
		updates["business_logo_url"] = *req.BusinessLogoURL
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.settingRepo.BulkUpsert(ctx, updates); err != nil {
		return fmt.Errorf("gagal memperbarui pengaturan: %w", err)
	}

	s.auditLogger.Log(userID, "UPDATE", "setting", 0, oldSettings, updates, ipAddress)
	return nil
}

// UpdateLogo updates the business logo URL setting.
func (s *SettingService) UpdateLogo(ctx context.Context, logoURL string, userID uint, ipAddress string) error {
	updates := map[string]string{
		"business_logo_url": logoURL,
	}
	if err := s.settingRepo.BulkUpsert(ctx, updates); err != nil {
		return fmt.Errorf("gagal memperbarui logo: %w", err)
	}
	s.auditLogger.Log(userID, "UPDATE", "setting", 0, nil, updates, ipAddress)
	return nil
}

// InitDefaults initializes default settings if they don't exist.
func (s *SettingService) InitDefaults(ctx context.Context) error {
	return s.settingRepo.InitDefaults(ctx)
}
