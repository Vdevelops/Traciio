package reminder

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reminder represents a reminder for a task
type Reminder struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID      string         `gorm:"type:uuid;not null;index" json:"task_id"`
	Task        *TaskRef       `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	RemindAt    time.Time      `gorm:"type:timestamp;not null" json:"remind_at"`
	ReminderType string        `gorm:"type:varchar(50);not null;default:'in_app'" json:"reminder_type"` // in_app, email, sms
	IsSent      bool           `gorm:"type:boolean;default:false" json:"is_sent"`
	SentAt      *time.Time     `gorm:"type:timestamp" json:"sent_at"`
	Message     string         `gorm:"type:text" json:"message"`
	CreatedBy   string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Reminder
func (Reminder) TableName() string {
	return "reminders"
}

// BeforeCreate hook to generate UUID
func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// TaskRef represents task reference in reminder
type TaskRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Title string `json:"title"`
}

// TableName specifies the table name for TaskRef
func (TaskRef) TableName() string {
	return "tasks"
}

// ReminderResponse represents reminder response DTO
type ReminderResponse struct {
	ID          string         `json:"id"`
	TaskID      string         `json:"task_id"`
	Task        *TaskRefResponse `json:"task,omitempty"`
	RemindAt    time.Time      `json:"remind_at"`
	ReminderType string        `json:"reminder_type"`
	IsSent      bool           `json:"is_sent"`
	SentAt      *time.Time     `json:"sent_at"`
	Message     string         `json:"message"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TaskRefResponse represents task in reminder response
type TaskRefResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ToReminderResponse converts Reminder to ReminderResponse
func (r *Reminder) ToReminderResponse() *ReminderResponse {
	resp := &ReminderResponse{
		ID:          r.ID,
		TaskID:      r.TaskID,
		RemindAt:    r.RemindAt,
		ReminderType: r.ReminderType,
		IsSent:      r.IsSent,
		SentAt:      r.SentAt,
		Message:     r.Message,
		CreatedBy:   r.CreatedBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	if r.Task != nil {
		resp.Task = &TaskRefResponse{
			ID:    r.Task.ID,
			Title: r.Task.Title,
		}
	}

	return resp
}

// CreateReminderRequest represents create reminder request DTO
type CreateReminderRequest struct {
	TaskID      string    `json:"task_id" binding:"required,uuid"`
	RemindAt    time.Time `json:"remind_at" binding:"required"`
	ReminderType string   `json:"reminder_type" binding:"omitempty,oneof=in_app email sms"`
	Message     string    `json:"message" binding:"omitempty"`
}

// UpdateReminderRequest represents update reminder request DTO
type UpdateReminderRequest struct {
	RemindAt    *time.Time `json:"remind_at" binding:"omitempty"`
	ReminderType string    `json:"reminder_type" binding:"omitempty,oneof=in_app email sms"`
	Message     string     `json:"message" binding:"omitempty"`
}

// ListRemindersRequest represents list reminders query parameters
type ListRemindersRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	PerPage      int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	TaskID       string `form:"task_id" binding:"omitempty,uuid"`
	ReminderType string `form:"reminder_type" binding:"omitempty,oneof=in_app email sms"`
	IsSent       *bool  `form:"is_sent" binding:"omitempty"`
	RemindAtFrom *time.Time `form:"remind_at_from" binding:"omitempty"`
	RemindAtTo   *time.Time `form:"remind_at_to" binding:"omitempty"`
}



