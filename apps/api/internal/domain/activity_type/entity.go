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
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`                                   // lucide icon name
	BadgeColor  string         `gorm:"type:varchar(50);not null;default:'outline'" json:"badge_color"` // default, secondary, destructive, outline
	Status      string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`       // active, inactive
	Order       int            `gorm:"type:integer;not null;default:0" json:"order"`                   // Display order
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
