package google_calendar_token

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoogleCalendarToken represents a user's Google Calendar OAuth2 token
type GoogleCalendarToken struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"` // One token per user
	AccessToken  string         `gorm:"type:text;not null" json:"-"`                   // Encrypted access token (hidden from JSON)
	RefreshToken string         `gorm:"type:text;not null" json:"-"`                   // Encrypted refresh token (hidden from JSON)
	TokenType    string         `gorm:"type:varchar(50);default:'Bearer'" json:"token_type"`
	ExpiresAt    time.Time      `gorm:"type:timestamp;not null;index" json:"expires_at"`
	Scope        string         `gorm:"type:text" json:"scope"` // OAuth2 scopes
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GoogleCalendarToken
func (GoogleCalendarToken) TableName() string {
	return "google_calendar_tokens"
}

// BeforeCreate hook to generate UUID
func (t *GoogleCalendarToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// IsExpired checks if the token is expired
func (t *GoogleCalendarToken) IsExpired() bool {
	// Add 5 minute buffer before actual expiration
	return time.Now().Add(5 * time.Minute).After(t.ExpiresAt)
}

// NeedsRefresh checks if the token needs to be refreshed
func (t *GoogleCalendarToken) NeedsRefresh() bool {
	// Refresh if expires within 10 minutes
	return time.Now().Add(10 * time.Minute).After(t.ExpiresAt)
}
