package contact_role

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContactRole represents a contact role entity (Doctor, PIC, Manager, Other)
type ContactRole struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	BadgeColor  string    `gorm:"type:varchar(50);not null;default:'outline'" json:"badge_color"` // default, secondary, outline, success, warning, active
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Read-only fields
	ContactCount int64 `gorm:"->" json:"contact_count"`
}

// TableName specifies the table name for ContactRole
func (ContactRole) TableName() string {
	return "contact_roles"
}

// BeforeCreate hook to generate UUID
func (cr *ContactRole) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == "" {
		cr.ID = uuid.New().String()
	}
	return nil
}

// ContactRoleResponse represents contact role response DTO
type ContactRoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	BadgeColor  string    `json:"badge_color"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ContactCount int64    `json:"contact_count"`
}

// ToContactRoleResponse converts ContactRole to ContactRoleResponse
func (cr *ContactRole) ToContactRoleResponse() *ContactRoleResponse {
	return &ContactRoleResponse{
		ID:          cr.ID,
		Name:        cr.Name,
		Code:        cr.Code,
		Description: cr.Description,
		BadgeColor:  cr.BadgeColor,
		Status:      cr.Status,
		CreatedAt:   cr.CreatedAt,
		UpdatedAt:   cr.UpdatedAt,
		ContactCount: cr.ContactCount,
	}
}

// CreateContactRoleRequest represents create contact role request DTO
type CreateContactRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Code        string `json:"code" binding:"required,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateContactRoleRequest represents update contact role request DTO
type UpdateContactRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3"`
	Code        string `json:"code" binding:"omitempty,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

