package ai_settings

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new AI settings repository
func NewRepository(db *gorm.DB) interfaces.AISettingsRepository {
	return &repository{db: db}
}

func (r *repository) GetSettings() (*ai_settings.AISettings, error) {
	var settings ai_settings.AISettings
	err := r.db.First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default settings if not exists
			defaultSettings := &ai_settings.AISettings{
				Enabled:     true,
				Provider:    "cerebras",
				Model:       "llama-3.1-8b",
				DataPrivacy: nil, // Will be set to default in service
			}
			if err := r.db.Create(defaultSettings).Error; err != nil {
				return nil, err
			}
			return defaultSettings, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *repository) UpdateSettings(settings *ai_settings.AISettings) error {
	return r.db.Save(settings).Error
}


