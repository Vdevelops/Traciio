package activity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Activity represents an activity entity
type Activity struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type           string         `gorm:"type:varchar(50);not null;index" json:"type"`       // visit, call, email, task, deal (kept for backward compatibility)
	ActivityTypeID *string        `gorm:"type:uuid;index" json:"activity_type_id,omitempty"` // Reference to activity_types table
	AccountID      *string        `gorm:"type:uuid;index" json:"account_id,omitempty"`
	ContactID      *string        `gorm:"type:uuid;index" json:"contact_id,omitempty"`
	DealID         *string        `gorm:"type:uuid;index" json:"deal_id,omitempty"` // Optional link to deal
	LeadID         *string        `gorm:"type:uuid;index" json:"lead_id,omitempty"` // Optional link to lead
	UserID         string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Description    string         `gorm:"type:text;not null" json:"description"`
	Timestamp      time.Time      `gorm:"type:timestamp;not null;index" json:"timestamp"`
	Metadata       datatypes.JSON `gorm:"type:jsonb;index" json:"metadata,omitempty"` // Additional data as JSON
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (for preloading)
	Account      interface{} `gorm:"-" json:"account,omitempty"`
	Contact      interface{} `gorm:"-" json:"contact,omitempty"`
	User         interface{} `gorm:"-" json:"user,omitempty"`
	ActivityType interface{} `gorm:"-" json:"activity_type,omitempty"`
}

// TableName specifies the table name for Activity
func (Activity) TableName() string {
	return "activities"
}

// BeforeCreate hook to generate UUID
func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// ToActivityResponse converts Activity to ActivityResponse
func (a *Activity) ToActivityResponse() *ActivityResponse {
	var metadata interface{}
	if a.Metadata != nil {
		// Parse JSON to interface{}
		// This will be handled in the service layer
	}

	resp := &ActivityResponse{
		ID:             a.ID,
		Type:           a.Type,
		ActivityTypeID: a.ActivityTypeID,
		AccountID:      a.AccountID,
		ContactID:      a.ContactID,
		DealID:         a.DealID,
		LeadID:         a.LeadID,
		UserID:         a.UserID,
		Description:    a.Description,
		Timestamp:      a.Timestamp,
		Metadata:       metadata,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Account:        a.Account,
		Contact:        a.Contact,
		User:           a.User,
		ActivityType:   a.ActivityType,
	}
	return resp
}
