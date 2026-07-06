package reminder

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reminder represents a reminder for a task
type Reminder struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID       string         `gorm:"type:uuid;not null;index" json:"task_id"`
	Task         *TaskRef       `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	RemindAt     time.Time      `gorm:"type:timestamp;not null" json:"remind_at"`
	ReminderType string         `gorm:"type:varchar(50);not null;default:'in_app'" json:"reminder_type"` // in_app, email, sms
	IsSent       bool           `gorm:"type:boolean;default:false" json:"is_sent"`
	SentAt       *time.Time     `gorm:"type:timestamp" json:"sent_at"`
	Message      string         `gorm:"type:text" json:"message"`
	CreatedBy    string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
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
