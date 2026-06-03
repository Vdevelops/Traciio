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

// ActivityResponse represents activity response DTO
type ActivityResponse struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"`
	ActivityTypeID *string     `json:"activity_type_id,omitempty"`
	AccountID      *string     `json:"account_id,omitempty"`
	ContactID      *string     `json:"contact_id,omitempty"`
	DealID         *string     `json:"deal_id,omitempty"`
	LeadID         *string     `json:"lead_id,omitempty"`
	UserID         string      `json:"user_id"`
	Description    string      `json:"description"`
	Timestamp      time.Time   `json:"timestamp"`
	Metadata       interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Account        interface{} `json:"account,omitempty"`
	Contact        interface{} `json:"contact,omitempty"`
	User           interface{} `json:"user,omitempty"`
	ActivityType   interface{} `json:"activity_type,omitempty"`
	Deal           interface{} `json:"deal,omitempty"`
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

// CreateActivityRequest represents create activity request DTO
type CreateActivityRequest struct {
	Type           string      `json:"type" binding:"omitempty,oneof=visit call email task deal"` // Kept for backward compatibility
	ActivityTypeID *string     `json:"activity_type_id" binding:"omitempty,uuid"`                 // New field for dynamic activity types
	AccountID      *string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID      *string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID         *string     `json:"deal_id" binding:"omitempty,uuid"` // Optional link to deal
	LeadID         *string     `json:"lead_id" binding:"omitempty,uuid"` // Optional link to lead
	UserID         string      `json:"user_id" binding:"omitempty,uuid"` // Will be set from context
	Description    string      `json:"description" binding:"required,min=3"`
	Timestamp      string      `json:"timestamp" binding:"required"`
	Metadata       interface{} `json:"metadata" binding:"omitempty"`
}

// UpdateActivityRequest represents update activity request DTO.
type UpdateActivityRequest struct {
	ActivityTypeID *string     `json:"activity_type_id" binding:"omitempty,uuid"`
	Description    string      `json:"description" binding:"omitempty,min=3"`
	Timestamp      string      `json:"timestamp" binding:"omitempty"`
	Metadata       interface{} `json:"metadata" binding:"omitempty"`
}

// ListActivitiesRequest represents list activities query parameters
type ListActivitiesRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Type      string `form:"type" binding:"omitempty,oneof=visit call email task deal"`
	AccountID string `form:"account_id" binding:"omitempty,uuid"`
	ContactID string `form:"contact_id" binding:"omitempty,uuid"`
	DealID    string `form:"deal_id" binding:"omitempty,uuid"` // Filter by deal
	LeadID    string `form:"lead_id" binding:"omitempty,uuid"` // Filter by lead
	UserID    string `form:"user_id" binding:"omitempty,uuid"`
	// ScopedUserIDs is injected by RBAC scope middleware and not user-provided.
	ScopedUserIDs []string `form:"-" json:"-"`
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
}

// ActivityTimelineRequest represents activity timeline query parameters
type ActivityTimelineRequest struct {
	AccountID string `form:"account_id" binding:"omitempty,uuid"`
	ContactID string `form:"contact_id" binding:"omitempty,uuid"`
	DealID    string `form:"deal_id" binding:"omitempty,uuid"` // Filter by deal
	LeadID    string `form:"lead_id" binding:"omitempty,uuid"` // Filter by lead
	UserID    string `form:"user_id" binding:"omitempty,uuid"`
	// ScopedUserIDs is injected by RBAC scope middleware and not user-provided.
	ScopedUserIDs []string `form:"-" json:"-"`
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
