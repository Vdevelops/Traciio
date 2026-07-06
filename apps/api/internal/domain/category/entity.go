package category

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents an account category entity (Hospital, Clinic, Pharmacy)
type Category struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	BadgeColor  string    `gorm:"type:varchar(50);not null;default:'outline'" json:"badge_color"` // default, secondary, outline, success, warning, active
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Read-only fields
	AccountCount int64 `gorm:"->" json:"account_count"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate hook to generate UUID
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// ToCategoryResponse converts Category to CategoryResponse
func (c *Category) ToCategoryResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Code:        c.Code,
		Description: c.Description,
		BadgeColor:  c.BadgeColor,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		AccountCount: c.AccountCount,
	}
}
