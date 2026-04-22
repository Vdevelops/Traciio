package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/google_calendar_token"
)

// GoogleCalendarTokenRepository defines the interface for Google Calendar token repository
type GoogleCalendarTokenRepository interface {
	// FindByUserID finds a token by user ID
	FindByUserID(userID string) (*google_calendar_token.GoogleCalendarToken, error)

	// Create creates a new token
	Create(token *google_calendar_token.GoogleCalendarToken) error

	// Update updates an existing token
	Update(token *google_calendar_token.GoogleCalendarToken) error

	// Delete deletes a token by user ID
	DeleteByUserID(userID string) error
}

