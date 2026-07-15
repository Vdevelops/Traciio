package ai_settings

import (
	"encoding/json"
	"time"
)

// AISettingsResponse represents AI settings response DTO.
type AISettingsResponse struct {
	ID          string              `json:"id"`
	Enabled     bool                `json:"enabled"`
	Provider    string              `json:"provider"`
	Model       string              `json:"model"`
	BaseURL     string              `json:"base_url,omitempty"`
	DataPrivacy DataPrivacySettings `json:"data_privacy"`
	Timezone    string              `json:"timezone"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// ToAISettingsResponse converts AISettings to AISettingsResponse.
func (a *AISettings) ToAISettingsResponse() *AISettingsResponse {
	dataPrivacy := DefaultDataPrivacySettings()
	if a.DataPrivacy != nil {
		_ = json.Unmarshal(a.DataPrivacy, &dataPrivacy)
	}

	timezone := a.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	return &AISettingsResponse{
		ID:          a.ID,
		Enabled:     a.Enabled,
		Provider:    a.Provider,
		Model:       a.Model,
		BaseURL:     a.BaseURL,
		DataPrivacy: dataPrivacy,
		Timezone:    timezone,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// UpdateAISettingsRequest represents update AI settings request DTO.
type UpdateAISettingsRequest struct {
	Enabled     *bool                `json:"enabled" binding:"omitempty"`
	Provider    string               `json:"provider" binding:"omitempty,oneof=cerebras openai anthropic"`
	APIKey      string               `json:"api_key" binding:"omitempty"`
	Model       string               `json:"model" binding:"omitempty"`
	BaseURL     string               `json:"base_url" binding:"omitempty"`
	DataPrivacy *DataPrivacySettings `json:"data_privacy" binding:"omitempty"`
	Timezone    string               `json:"timezone" binding:"omitempty"`
}
