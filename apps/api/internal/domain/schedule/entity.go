package schedule

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule represents a user-specific schedule connected to a task (or standalone)
type Schedule struct {
	ID                       string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID                   *string        `gorm:"type:uuid;index" json:"task_id"` // Nullable: schedule can exist without task
	Task                     *TaskRef       `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	UserID                   string         `gorm:"type:uuid;not null;index" json:"user_id"` // User who owns this schedule (from task.assigned_to)
	Title                    string         `gorm:"type:varchar(255);not null;index:idx_schedules_fts,type:gin,expression:to_tsvector('english'\\, title || ' ' || COALESCE(description\\, ''))" json:"title"`
	Description              *string        `gorm:"type:text" json:"description"`
	ScheduledAt              time.Time      `gorm:"type:timestamp;not null;index" json:"scheduled_at"`
	Status                   string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, submitted, confirmed, completed, cancelled, rejected
	GoogleCalendarEventID    *string        `gorm:"type:varchar(255);index" json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string         `gorm:"type:varchar(20);not null;default:'not_synced'" json:"google_calendar_sync_status"` // not_synced, synced, sync_failed
	GoogleCalendarSyncedAt   *time.Time     `gorm:"type:timestamp" json:"google_calendar_synced_at"`
	GoogleCalendarEventLink  *string        `gorm:"type:varchar(500)" json:"google_calendar_event_link"` // Direct URL to view in Google Calendar
	ReminderMinutesBefore    *int           `gorm:"type:integer" json:"reminder_minutes_before"`         // Minutes before task due_date to remind
	CreatedBy                string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Schedule
func (Schedule) TableName() string {
	return "schedules"
}

// BeforeCreate hook to generate UUID
func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// TaskRef represents task reference in schedule
type TaskRef struct {
	ID           string     `gorm:"type:uuid;primary_key" json:"id"`
	Title        string     `json:"title"`
	DueDate      *time.Time `json:"due_date"`
	AssignedTo   *string    `gorm:"type:uuid" json:"assigned_to"`
	AssignedUser *UserRef   `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
}

// TableName specifies the table name for TaskRef
func (TaskRef) TableName() string {
	return "tasks"
}

// UserRef represents user reference in schedule
type UserRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TableName specifies the table name for UserRef
func (UserRef) TableName() string {
	return "users"
}
