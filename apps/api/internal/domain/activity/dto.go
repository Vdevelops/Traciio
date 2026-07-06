package activity

import "time"

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

// CreateActivityRequest represents create activity request DTO
type CreateActivityRequest struct {
	Type           string      `json:"type" binding:"omitempty,oneof=visit call email task deal"`
	ActivityTypeID *string     `json:"activity_type_id" binding:"omitempty,uuid"`
	AccountID      *string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID      *string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID         *string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID         *string     `json:"lead_id" binding:"omitempty,uuid"`
	UserID         string      `json:"user_id" binding:"omitempty,uuid"`
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
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Type          string   `form:"type" binding:"omitempty,oneof=visit call email task deal"`
	AccountID     string   `form:"account_id" binding:"omitempty,uuid"`
	ContactID     string   `form:"contact_id" binding:"omitempty,uuid"`
	DealID        string   `form:"deal_id" binding:"omitempty,uuid"`
	LeadID        string   `form:"lead_id" binding:"omitempty,uuid"`
	UserID        string   `form:"user_id" binding:"omitempty,uuid"`
	ScopedUserIDs []string `form:"-" json:"-"`
	StartDate     string   `form:"start_date" binding:"omitempty"`
	EndDate       string   `form:"end_date" binding:"omitempty"`
}

// ActivityTimelineRequest represents activity timeline query parameters
type ActivityTimelineRequest struct {
	AccountID     string   `form:"account_id" binding:"omitempty,uuid"`
	ContactID     string   `form:"contact_id" binding:"omitempty,uuid"`
	DealID        string   `form:"deal_id" binding:"omitempty,uuid"`
	LeadID        string   `form:"lead_id" binding:"omitempty,uuid"`
	UserID        string   `form:"user_id" binding:"omitempty,uuid"`
	ScopedUserIDs []string `form:"-" json:"-"`
	StartDate     string   `form:"start_date" binding:"omitempty"`
	EndDate       string   `form:"end_date" binding:"omitempty"`
	Limit         int      `form:"limit" binding:"omitempty,min=1,max=100"`
}
