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

// NotificationResponse represents notification response DTO
type NotificationResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToNotificationResponse converts Notification to NotificationResponse
func (n *Notification) ToNotificationResponse() *NotificationResponse {
	return &NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		IsRead:    n.IsRead,
		ReadAt:    n.ReadAt,
		Data:      n.Data,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// CreateNotificationRequest represents create notification request DTO
type CreateNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required,uuid"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"omitempty"`
	Type    string `json:"type" binding:"omitempty,oneof=reminder task deal activity"`
	Data    string `json:"data" binding:"omitempty"`
}

// ListNotificationsRequest represents list notifications query parameters
type ListNotificationsRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	UserID  string `form:"user_id" binding:"omitempty,uuid"`
	Type    string `form:"type" binding:"omitempty,oneof=reminder task deal activity"`
	IsRead  *bool  `form:"is_read" binding:"omitempty"`
}

// MarkAsReadRequest represents mark as read request DTO
type MarkAsReadRequest struct {
	ReadAt *time.Time `json:"read_at" binding:"omitempty"`
}

