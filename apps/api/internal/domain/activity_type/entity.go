package activity_type

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActivityType represents activity type entity (visit, call, email, task, deal, etc.)
type ActivityType struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"` // lucide icon name
	BadgeColor  string         `gorm:"type:varchar(50);not null;default:'outline'" json:"badge_color"` // default, secondary, destructive, outline
	Status      string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, inactive
	Order       int            `gorm:"type:integer;not null;default:0" json:"order"` // Display order
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// Read-only fields
	ActivityCount int64 `gorm:"->" json:"activity_count"`
}

// TableName specifies the table name for ActivityType
func (ActivityType) TableName() string {
	return "activity_types"
}

// BeforeCreate hook to generate UUID
func (at *ActivityType) BeforeCreate(tx *gorm.DB) error {
	if at.ID == "" {
		at.ID = uuid.New().String()
	}
	return nil
}

// ActivityTypeResponse represents activity type response DTO
type ActivityTypeResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	BadgeColor  string    `json:"badge_color"`
	Status      string    `json:"status"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ActivityCount int64    `json:"activity_count"`
}

// ToActivityTypeResponse converts ActivityType to ActivityTypeResponse
func (at *ActivityType) ToActivityTypeResponse() *ActivityTypeResponse {
	return &ActivityTypeResponse{
		ID:          at.ID,
		Name:        at.Name,
		Code:        at.Code,
		Description: at.Description,
		Icon:        at.Icon,
		BadgeColor:  at.BadgeColor,
		Status:      at.Status,
		Order:       at.Order,
		CreatedAt:   at.CreatedAt,
		UpdatedAt:   at.UpdatedAt,
		ActivityCount: at.ActivityCount,
	}
}

// CreateActivityTypeRequest represents create activity type request DTO
type CreateActivityTypeRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary destructive outline"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	Order       int    `json:"order" binding:"omitempty"`
}

// UpdateActivityTypeRequest represents update activity type request DTO
type UpdateActivityTypeRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary destructive outline"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	Order       *int   `json:"order" binding:"omitempty"`
}

// ListActivityTypesRequest represents list activity types query parameters
type ListActivityTypesRequest struct {
	Status string `form:"status" binding:"omitempty,oneof=active inactive"`
}

