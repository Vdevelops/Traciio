package ai_settings

import (
	"encoding/json"
	"fmt"

	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
)

// Service represents AI settings service
type Service struct {
	settingsRepo interfaces.AISettingsRepository
}

// NewService creates a new AI settings service
func NewService(settingsRepo interfaces.AISettingsRepository) *Service {
	return &Service{
		settingsRepo: settingsRepo,
	}
}

// GetSettings gets current AI settings
func (s *Service) GetSettings() (*ai_settings.AISettingsResponse, error) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// Set default data privacy if nil
	if settings.DataPrivacy == nil {
		defaultPrivacy := ai_settings.DataPrivacySettings{
			AllowVisitReports: true,
			AllowAccounts:     true,
			AllowContacts:     true,
			AllowDeals:        true,
			AllowLeads:        true,
			AllowActivities:   true,
			AllowTasks:         true,
			AllowProducts:      true,
		}
		privacyJSON, _ := json.Marshal(defaultPrivacy)
		settings.DataPrivacy = privacyJSON
	}

	return settings.ToAISettingsResponse(), nil
}

// UpdateSettings updates AI settings
func (s *Service) UpdateSettings(req *ai_settings.UpdateAISettingsRequest) (*ai_settings.AISettingsResponse, error) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// Update fields if provided
	if req.Enabled != nil {
		settings.Enabled = *req.Enabled
	}
	if req.Provider != "" {
		settings.Provider = req.Provider
	}
	if req.APIKey != "" {
		settings.APIKey = req.APIKey // In production, this should be encrypted
	}
	if req.Model != "" {
		settings.Model = req.Model
	}
	if req.BaseURL != "" {
		settings.BaseURL = req.BaseURL
	}
	if req.DataPrivacy != nil {
		privacyJSON, err := json.Marshal(req.DataPrivacy)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data privacy: %w", err)
		}
		settings.DataPrivacy = privacyJSON
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}

	if err := s.settingsRepo.UpdateSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return settings.ToAISettingsResponse(), nil
}



