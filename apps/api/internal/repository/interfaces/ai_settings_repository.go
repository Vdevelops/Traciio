package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
)

// AISettingsRepository defines the interface for AI settings repository
type AISettingsRepository interface {
	// GetSettings gets the current AI settings (there should only be one)
	GetSettings() (*ai_settings.AISettings, error)
	
	// UpdateSettings updates AI settings
	UpdateSettings(settings *ai_settings.AISettings) error
}


