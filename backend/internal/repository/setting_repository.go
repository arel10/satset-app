package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// SettingRepository handles database operations for business settings.
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository creates a new SettingRepository.
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetAll retrieves all settings.
func (r *SettingRepository) GetAll(ctx context.Context) ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

// GetByKey retrieves a single setting by its key.
func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*model.Setting, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Upsert creates or updates a setting by key.
func (r *SettingRepository) Upsert(ctx context.Context, key, value string) error {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(&model.Setting{Key: key, Value: value}).Error
	}
	if err != nil {
		return err
	}
	setting.Value = value
	return r.db.WithContext(ctx).Save(&setting).Error
}

// BulkUpsert creates or updates multiple settings.
func (r *SettingRepository) BulkUpsert(ctx context.Context, settings map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range settings {
			var setting model.Setting
			err := tx.Where("`key` = ?", key).First(&setting).Error
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&model.Setting{Key: key, Value: value}).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			setting.Value = value
			if err := tx.Save(&setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// InitDefaults creates default settings if they don't exist.
func (r *SettingRepository) InitDefaults(ctx context.Context) error {
	defaults := model.DefaultSettings()
	for _, s := range defaults {
		var existing model.Setting
		err := r.db.WithContext(ctx).Where("`key` = ?", s.Key).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(&s).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

// GetAllAsMap returns all settings as a key-value map.
func (r *SettingRepository) GetAllAsMap(ctx context.Context) (map[string]string, error) {
	settings, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}
