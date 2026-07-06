package notification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents an in-app notification
type Notification struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text" json:"message"`
	Type      string         `gorm:"type:varchar(50);not null;default:'reminder'" json:"type"` // reminder, task, deal, activity
	IsRead    bool           `gorm:"type:boolean;default:false;index" json:"is_read"`
	ReadAt    *time.Time     `gorm:"type:timestamp" json:"read_at"`
	Data      string         `gorm:"type:jsonb" json:"data"` // Additional data as JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook to generate UUID
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
