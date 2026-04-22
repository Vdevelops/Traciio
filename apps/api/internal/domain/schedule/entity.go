package schedule

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule represents a user-specific schedule connected to a task (or standalone)
type Schedule struct {
	ID                        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID                    *string        `gorm:"type:uuid;index" json:"task_id"` // Nullable: schedule can exist without task
	Task                      *TaskRef      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	UserID                    string         `gorm:"type:uuid;not null;index" json:"user_id"` // User who owns this schedule (from task.assigned_to)
	Title                     string         `gorm:"type:varchar(255);not null;index:idx_schedules_fts,type:gin,expression:to_tsvector('english'\\, title || ' ' || COALESCE(description\\, ''))" json:"title"`
	Description               *string        `gorm:"type:text" json:"description"`
	ScheduledAt               time.Time      `gorm:"type:timestamp;not null;index" json:"scheduled_at"`
	Status                    string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, submitted, confirmed, completed, cancelled, rejected
	GoogleCalendarEventID     *string        `gorm:"type:varchar(255);index" json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus  string         `gorm:"type:varchar(20);not null;default:'not_synced'" json:"google_calendar_sync_status"` // not_synced, synced, sync_failed
	GoogleCalendarSyncedAt    *time.Time     `gorm:"type:timestamp" json:"google_calendar_synced_at"`
	GoogleCalendarEventLink   *string        `gorm:"type:varchar(500)" json:"google_calendar_event_link"` // Direct URL to view in Google Calendar
	ReminderMinutesBefore     *int           `gorm:"type:integer" json:"reminder_minutes_before"` // Minutes before task due_date to remind
	CreatedBy                 string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:"-"`
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
	ID          string     `gorm:"type:uuid;primary_key" json:"id"`
	Title       string     `json:"title"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *string    `gorm:"type:uuid" json:"assigned_to"`
	AssignedUser *UserRef  `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
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

// ScheduleResponse represents schedule response DTO
type ScheduleResponse struct {
	ID                       string            `json:"id"`
	TaskID                   *string           `json:"task_id"` // Nullable: schedule can exist without task
	Task                     *TaskRefResponse  `json:"task,omitempty"`
	UserID                   string            `json:"user_id"`
	Title                    string            `json:"title"`
	Description              *string           `json:"description"`
	ScheduledAt              time.Time          `json:"scheduled_at"`
	Status                   string            `json:"status"`
	GoogleCalendarEventID    *string           `json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string            `json:"google_calendar_sync_status"`
	GoogleCalendarSyncedAt   *time.Time        `json:"google_calendar_synced_at"`
	GoogleCalendarEventLink  *string           `json:"google_calendar_event_link"`
	ReminderMinutesBefore    *int              `json:"reminder_minutes_before"`
	CreatedBy                string            `json:"created_by"`
	CreatedAt                time.Time         `json:"created_at"`
	UpdatedAt                time.Time         `json:"updated_at"`
}

// TaskRefResponse represents task in schedule response
type TaskRefResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	DueDate     *time.Time        `json:"due_date"`
	AssignedTo  string            `json:"assigned_to"`
	AssignedUser *UserRefResponse `json:"assigned_user,omitempty"`
}

// UserRefResponse represents user in schedule response
type UserRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToScheduleResponse converts Schedule to ScheduleResponse
func (s *Schedule) ToScheduleResponse() *ScheduleResponse {
	var taskID *string
	if s.TaskID != nil {
		taskID = s.TaskID
	}
	
	resp := &ScheduleResponse{
		ID:                       s.ID,
		TaskID:                   taskID,
		UserID:                   s.UserID,
		Title:                    s.Title,
		Description:              s.Description,
		ScheduledAt:              s.ScheduledAt,
		Status:                   s.Status,
		GoogleCalendarEventID:    s.GoogleCalendarEventID,
		GoogleCalendarSyncStatus: s.GoogleCalendarSyncStatus,
		GoogleCalendarSyncedAt:   s.GoogleCalendarSyncedAt,
		GoogleCalendarEventLink:  s.GoogleCalendarEventLink,
		ReminderMinutesBefore:    s.ReminderMinutesBefore,
		CreatedBy:                s.CreatedBy,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
	}

	if s.Task != nil {
		resp.Task = &TaskRefResponse{
			ID:      s.Task.ID,
			Title:   s.Task.Title,
			DueDate: s.Task.DueDate,
		}
		if s.Task.AssignedTo != nil {
			resp.Task.AssignedTo = *s.Task.AssignedTo
		}
		if s.Task.AssignedUser != nil {
			resp.Task.AssignedUser = &UserRefResponse{
				ID:    s.Task.AssignedUser.ID,
				Name:  s.Task.AssignedUser.Name,
				Email: s.Task.AssignedUser.Email,
			}
		}
	}

	return resp
}

// CreateScheduleRequest represents create schedule request DTO
type CreateScheduleRequest struct {
	TaskID                string     `json:"task_id" binding:"omitempty,uuid"` // Optional: schedule can be created without task
	Title                 string     `json:"title" binding:"required,min=3,max=255"`
	Description           string     `json:"description" binding:"omitempty"`
	ScheduledAt           time.Time  `json:"scheduled_at" binding:"required"`
	ReminderMinutesBefore *int       `json:"reminder_minutes_before" binding:"omitempty,min=0,max=10080"`
	SyncToGoogleCalendar  bool       `json:"sync_to_google_calendar" binding:"omitempty"`
}

// UpdateScheduleRequest represents update schedule request DTO
type UpdateScheduleRequest struct {
	Title                 string     `json:"title" binding:"omitempty,min=3,max=255"`
	Description           string     `json:"description" binding:"omitempty"`
	ScheduledAt           *time.Time `json:"scheduled_at" binding:"omitempty"`
	Status                string     `json:"status" binding:"omitempty,oneof=pending submitted confirmed completed cancelled rejected"`
	ReminderMinutesBefore *int       `json:"reminder_minutes_before" binding:"omitempty,min=0,max=10080"`
	SyncToGoogleCalendar  *bool      `json:"sync_to_google_calendar" binding:"omitempty"`
}

// ListSchedulesRequest represents list schedules query parameters
type ListSchedulesRequest struct {
	Page                     int        `form:"page" binding:"omitempty,min=1"`
	PerPage                  int        `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search                   string     `form:"search" binding:"omitempty"`
	Status                   string     `form:"status" binding:"omitempty,oneof=pending submitted confirmed completed cancelled rejected"`
	TaskID                   string     `form:"task_id" binding:"omitempty,uuid"`
	UserID                   string     `form:"user_id" binding:"omitempty,uuid"`
	ScheduledAtFrom          *time.Time `form:"scheduled_at_from" time_format:"2006-01-02T15:04:05Z07:00" binding:"omitempty"`
	ScheduledAtTo            *time.Time `form:"scheduled_at_to" time_format:"2006-01-02T15:04:05Z07:00" binding:"omitempty"`
	GoogleCalendarSyncStatus string     `form:"google_calendar_sync_status" binding:"omitempty,oneof=not_synced synced sync_failed"`
	ScopedUserIDs            []string   `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}

// GoogleCalendarSyncResponse represents Google Calendar sync response
type GoogleCalendarSyncResponse struct {
	ScheduleID            string `json:"schedule_id"`
	GoogleCalendarEventID string `json:"google_calendar_event_id"`
	EventURL              string `json:"event_url"`
	SyncStatus            string `json:"sync_status"`
}
