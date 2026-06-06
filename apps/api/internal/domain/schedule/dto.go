package schedule

import "time"

// ScheduleResponse represents schedule response DTO.
type ScheduleResponse struct {
	ID                       string           `json:"id"`
	TaskID                   *string          `json:"task_id"`
	Task                     *TaskRefResponse `json:"task,omitempty"`
	UserID                   string           `json:"user_id"`
	Title                    string           `json:"title"`
	Description              *string          `json:"description"`
	ScheduledAt              time.Time        `json:"scheduled_at"`
	Status                   string           `json:"status"`
	GoogleCalendarEventID    *string          `json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string           `json:"google_calendar_sync_status"`
	GoogleCalendarSyncedAt   *time.Time       `json:"google_calendar_synced_at"`
	GoogleCalendarEventLink  *string          `json:"google_calendar_event_link"`
	ReminderMinutesBefore    *int             `json:"reminder_minutes_before"`
	CreatedBy                string           `json:"created_by"`
	CreatedAt                time.Time        `json:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at"`
}

// TaskRefResponse represents task in schedule response.
type TaskRefResponse struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	DueDate      *time.Time       `json:"due_date"`
	AssignedTo   string           `json:"assigned_to"`
	AssignedUser *UserRefResponse `json:"assigned_user,omitempty"`
}

// UserRefResponse represents user in schedule response.
type UserRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToScheduleResponse converts Schedule to ScheduleResponse.
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

// CreateScheduleRequest represents create schedule request DTO.
type CreateScheduleRequest struct {
	TaskID                string    `json:"task_id" binding:"omitempty,uuid"`
	Title                 string    `json:"title" binding:"required,min=3,max=255"`
	Description           string    `json:"description" binding:"omitempty"`
	ScheduledAt           time.Time `json:"scheduled_at" binding:"required"`
	ReminderMinutesBefore *int      `json:"reminder_minutes_before" binding:"omitempty,min=0,max=10080"`
	SyncToGoogleCalendar  bool      `json:"sync_to_google_calendar" binding:"omitempty"`
}

// UpdateScheduleRequest represents update schedule request DTO.
type UpdateScheduleRequest struct {
	Title                 string     `json:"title" binding:"omitempty,min=3,max=255"`
	Description           string     `json:"description" binding:"omitempty"`
	ScheduledAt           *time.Time `json:"scheduled_at" binding:"omitempty"`
	Status                string     `json:"status" binding:"omitempty,oneof=pending submitted confirmed completed cancelled rejected"`
	ReminderMinutesBefore *int       `json:"reminder_minutes_before" binding:"omitempty,min=0,max=10080"`
	SyncToGoogleCalendar  *bool      `json:"sync_to_google_calendar" binding:"omitempty"`
}

// ListSchedulesRequest represents list schedules query parameters.
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
	ScopedUserIDs            []string   `form:"-" json:"-"`
}

// GoogleCalendarSyncResponse represents Google Calendar sync response.
type GoogleCalendarSyncResponse struct {
	ScheduleID            string `json:"schedule_id"`
	GoogleCalendarEventID string `json:"google_calendar_event_id"`
	EventURL              string `json:"event_url"`
	SyncStatus            string `json:"sync_status"`
}

// ScheduleAssignmentResponse represents schedule assignment response DTO.
type ScheduleAssignmentResponse struct {
	ID         string            `json:"id"`
	ScheduleID string            `json:"schedule_id"`
	Schedule   *ScheduleResponse `json:"schedule,omitempty"`
	UserID     string            `json:"user_id"`
	User       *UserRefResponse  `json:"user,omitempty"`
	AssignedAt time.Time         `json:"assigned_at"`
	Status     string            `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// ToScheduleAssignmentResponse converts ScheduleAssignment to DTO.
func (sa *ScheduleAssignment) ToScheduleAssignmentResponse() *ScheduleAssignmentResponse {
	resp := &ScheduleAssignmentResponse{
		ID:         sa.ID,
		ScheduleID: sa.ScheduleID,
		UserID:     sa.UserID,
		AssignedAt: sa.AssignedAt,
		Status:     sa.Status,
		CreatedAt:  sa.CreatedAt,
		UpdatedAt:  sa.UpdatedAt,
	}

	if sa.Schedule != nil {
		resp.Schedule = sa.Schedule.ToScheduleResponse()
	}

	if sa.User != nil {
		resp.User = &UserRefResponse{
			ID:    sa.User.ID,
			Name:  sa.User.Name,
			Email: sa.User.Email,
		}
	}

	return resp
}
