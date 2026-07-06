package google_calendar_token

import "time"

// GoogleCalendarTokenResponse represents token response DTO without sensitive data.
type GoogleCalendarTokenResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	Scope     string    `json:"scope"`
	IsExpired bool      `json:"is_expired"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts GoogleCalendarToken to GoogleCalendarTokenResponse.
func (t *GoogleCalendarToken) ToResponse() *GoogleCalendarTokenResponse {
	return &GoogleCalendarTokenResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenType: t.TokenType,
		ExpiresAt: t.ExpiresAt,
		Scope:     t.Scope,
		IsExpired: t.IsExpired(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
