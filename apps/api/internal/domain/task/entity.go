package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Task represents a task in the CRM system
type Task struct {
	ID               string      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title            string      `gorm:"type:varchar(255);not null;index:idx_tasks_fts,type:gin,expression:to_tsvector('english'\, title || ' ' || COALESCE(description\, ''))" json:"title"`
	Description      string      `gorm:"type:text" json:"description"`
	Type             string      `gorm:"type:varchar(50);not null;default:'general'" json:"type"`   // general, call, email, meeting, follow_up
	Status           string      `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, in_progress, completed, cancelled, submitted, approved, rejected
	Priority         string      `gorm:"type:varchar(20);default:'medium'" json:"priority"`         // low, medium, high, urgent
	DueDate          *time.Time  `gorm:"type:timestamp" json:"due_date"`
	CompletedAt      *time.Time  `gorm:"type:timestamp" json:"completed_at"`
	AssignedTo       *string     `gorm:"type:uuid;index" json:"assigned_to"` // User ID (optional)
	AssignedUser     *UserRef    `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
	AssignedFrom     *string     `gorm:"type:uuid;index" json:"assigned_from"` // User ID who assigned the task (optional)
	AssignedFromUser *UserRef    `gorm:"foreignKey:AssignedFrom" json:"assigned_from_user,omitempty"`
	AccountID        *string     `gorm:"type:uuid;index" json:"account_id"` // Optional: link to account
	Account          *AccountRef `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ContactID        *string     `gorm:"type:uuid;index" json:"contact_id"` // Optional: link to contact
	Contact          *ContactRef `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
	DealID           *string     `gorm:"type:uuid;index" json:"deal_id"` // Optional: link to deal
	Deal             *DealRef    `gorm:"foreignKey:DealID" json:"deal,omitempty"`
	LeadID           *string     `gorm:"type:uuid;index" json:"lead_id"` // Optional: link to lead (context-bound)
	Lead             *LeadRef    `gorm:"foreignKey:LeadID" json:"lead,omitempty"`
	TaskSource       string      `gorm:"type:varchar(50);not null;default:'manual';index" json:"task_source"` // manual, lead_tab, pipeline_tab, auto_generated
	// Schedule unification fields
	ScheduledStartTime *time.Time `gorm:"type:timestamptz" json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time `gorm:"type:timestamptz" json:"scheduled_end_time"`
	ScheduledLocation  string     `gorm:"type:varchar(500);default:''" json:"scheduled_location"`
	IsScheduleTask     bool       `gorm:"type:boolean;not null;default:false" json:"is_schedule_task"`
	// Quick action fields
	QuickActionType    string         `gorm:"type:varchar(50);default:''" json:"quick_action_type"`
	QuickActionPayload datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"quick_action_payload"`
	// Google Calendar integration
	GoogleCalendarEventID    *string        `gorm:"type:varchar(255);index" json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string         `gorm:"type:varchar(20);not null;default:'not_synced'" json:"google_calendar_sync_status"` // not_synced, synced, sync_failed
	GoogleCalendarSyncedAt   *time.Time     `gorm:"type:timestamp" json:"google_calendar_synced_at"`
	CreatedBy                string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Task
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate hook to generate UUID
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// UserRef represents user reference in task
type UserRef struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// TableName specifies the table name for UserRef
func (UserRef) TableName() string {
	return "users"
}

// AccountRef represents account reference in task
type AccountRef struct {
	ID   string `gorm:"type:uuid;primary_key" json:"id"`
	Name string `json:"name"`
}

// TableName specifies the table name for AccountRef
func (AccountRef) TableName() string {
	return "accounts"
}

// ContactRef represents contact reference in task
type ContactRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// TableName specifies the table name for ContactRef
func (ContactRef) TableName() string {
	return "contacts"
}

// DealRef represents deal reference in task
type DealRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Title string `json:"title"`
}

// TableName specifies the table name for DealRef
func (DealRef) TableName() string {
	return "deals"
}

// LeadRef represents lead reference in task
type LeadRef struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// TableName specifies the table name for LeadRef
func (LeadRef) TableName() string {
	return "leads"
}

// LeadRefResponse represents lead in task response
type LeadRefResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// TaskResponse represents task response DTO
type TaskResponse struct {
	ID               string              `json:"id"`
	Title            string              `json:"title"`
	Description      string              `json:"description"`
	Type             string              `json:"type"`
	Status           string              `json:"status"`
	Priority         string              `json:"priority"`
	DueDate          *time.Time          `json:"due_date"`
	CompletedAt      *time.Time          `json:"completed_at"`
	AssignedTo       string              `json:"assigned_to"`
	AssignedUser     *UserRefResponse    `json:"assigned_user,omitempty"`
	AssignedFrom     string              `json:"assigned_from"`
	AssignedFromUser *UserRefResponse    `json:"assigned_from_user,omitempty"`
	AccountID        string              `json:"account_id"`
	Account          *AccountRefResponse `json:"account,omitempty"`
	ContactID        string              `json:"contact_id"`
	Contact          *ContactRefResponse `json:"contact,omitempty"`
	DealID           string              `json:"deal_id"`
	Deal             *DealRefResponse    `json:"deal,omitempty"`
	LeadID           string              `json:"lead_id"`
	Lead             *LeadRefResponse    `json:"lead,omitempty"`
	TaskSource       string              `json:"task_source"`
	// Schedule fields
	ScheduledStartTime *time.Time `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time"`
	ScheduledLocation  string     `json:"scheduled_location"`
	IsScheduleTask     bool       `json:"is_schedule_task"`
	// Quick action fields
	QuickActionType string `json:"quick_action_type,omitempty"`
	// Google Calendar integration
	GoogleCalendarEventID    *string    `json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string     `json:"google_calendar_sync_status"`
	GoogleCalendarSyncedAt   *time.Time `json:"google_calendar_synced_at"`
	CreatedBy                string     `json:"created_by"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

// UserRefResponse represents user in task response
type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// AccountRefResponse represents account in task response
type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ContactRefResponse represents contact in task response
type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// DealRefResponse represents deal in task response
type DealRefResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ToTaskResponse converts Task to TaskResponse
func (t *Task) ToTaskResponse() *TaskResponse {
	resp := &TaskResponse{
		ID:                 t.ID,
		Title:              t.Title,
		Description:        t.Description,
		Type:               t.Type,
		Status:             t.Status,
		Priority:           t.Priority,
		DueDate:            t.DueDate,
		CompletedAt:        t.CompletedAt,
		AssignedTo:         "",
		AssignedFrom:       "",
		AccountID:          "",
		ContactID:          "",
		DealID:             "",
		LeadID:             "",
		TaskSource:         t.TaskSource,
		ScheduledStartTime: t.ScheduledStartTime,
		ScheduledEndTime:   t.ScheduledEndTime,
		ScheduledLocation:  t.ScheduledLocation,
		IsScheduleTask:     t.IsScheduleTask,
		QuickActionType:    t.QuickActionType,
		CreatedBy:          t.CreatedBy,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}

	if t.AssignedTo != nil {
		resp.AssignedTo = *t.AssignedTo
	}
	if t.AccountID != nil {
		resp.AccountID = *t.AccountID
	}
	if t.ContactID != nil {
		resp.ContactID = *t.ContactID
	}
	if t.DealID != nil {
		resp.DealID = *t.DealID
	}
	if t.LeadID != nil {
		resp.LeadID = *t.LeadID
	}

	if t.AssignedUser != nil {
		resp.AssignedUser = &UserRefResponse{
			ID:        t.AssignedUser.ID,
			Name:      t.AssignedUser.Name,
			Email:     t.AssignedUser.Email,
			AvatarURL: t.AssignedUser.AvatarURL,
		}
	}

	if t.AssignedFrom != nil {
		resp.AssignedFrom = *t.AssignedFrom
	}

	if t.AssignedFromUser != nil {
		resp.AssignedFromUser = &UserRefResponse{
			ID:        t.AssignedFromUser.ID,
			Name:      t.AssignedFromUser.Name,
			Email:     t.AssignedFromUser.Email,
			AvatarURL: t.AssignedFromUser.AvatarURL,
		}
	}

	if t.Account != nil {
		resp.Account = &AccountRefResponse{
			ID:   t.Account.ID,
			Name: t.Account.Name,
		}
	}

	if t.Contact != nil {
		resp.Contact = &ContactRefResponse{
			ID:    t.Contact.ID,
			Name:  t.Contact.Name,
			Email: t.Contact.Email,
			Phone: t.Contact.Phone,
		}
	}

	if t.Deal != nil {
		resp.Deal = &DealRefResponse{
			ID:    t.Deal.ID,
			Title: t.Deal.Title,
		}
	}

	if t.Lead != nil {
		resp.Lead = &LeadRefResponse{
			ID:        t.Lead.ID,
			FirstName: t.Lead.FirstName,
			LastName:  t.Lead.LastName,
			Email:     t.Lead.Email,
		}
	}

	resp.GoogleCalendarEventID = t.GoogleCalendarEventID
	resp.GoogleCalendarSyncStatus = t.GoogleCalendarSyncStatus
	resp.GoogleCalendarSyncedAt = t.GoogleCalendarSyncedAt

	return resp
}

// CreateTaskRequest represents create task request DTO
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=3,max=255"`
	Description string     `json:"description" binding:"omitempty"`
	Type        string     `json:"type" binding:"omitempty,oneof=general call email meeting follow_up"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
	AssignedTo  string     `json:"assigned_to" binding:"omitempty,uuid"`
	AccountID   string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID   string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID      string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID      string     `json:"lead_id" binding:"omitempty,uuid"`
	TaskSource  string     `json:"task_source" binding:"omitempty,oneof=manual lead_tab pipeline_tab auto_generated"`
	// Schedule fields
	ScheduledStartTime *time.Time `json:"scheduled_start_time" binding:"omitempty"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time" binding:"omitempty"`
	ScheduledLocation  string     `json:"scheduled_location" binding:"omitempty,max=500"`
	IsScheduleTask     bool       `json:"is_schedule_task" binding:"omitempty"`
}

// UpdateTaskRequest represents update task request DTO
type UpdateTaskRequest struct {
	Title       string     `json:"title" binding:"omitempty,min=3,max=255"`
	Description string     `json:"description" binding:"omitempty"`
	Type        string     `json:"type" binding:"omitempty,oneof=general call email meeting follow_up"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled submitted approved rejected"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
	AssignedTo  string     `json:"assigned_to" binding:"omitempty,uuid"`
	AccountID   string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID   string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID      string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID      string     `json:"lead_id" binding:"omitempty,uuid"`
	// Schedule fields
	ScheduledStartTime *time.Time `json:"scheduled_start_time" binding:"omitempty"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time" binding:"omitempty"`
	ScheduledLocation  string     `json:"scheduled_location" binding:"omitempty,max=500"`
	IsScheduleTask     *bool      `json:"is_schedule_task" binding:"omitempty"`
}

// AssignTaskRequest represents assign task request DTO
type AssignTaskRequest struct {
	AssignedTo string `json:"assigned_to" binding:"required,uuid"`
}

// ListTasksRequest represents list tasks query parameters
type ListTasksRequest struct {
	Page          int        `form:"page" binding:"omitempty,min=1"`
	PerPage       int        `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search        string     `form:"search" binding:"omitempty"`
	Status        string     `form:"status" binding:"omitempty"` // Can be comma-separated: "pending,in_progress"
	Priority      string     `form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Type          string     `form:"type" binding:"omitempty,oneof=general call email meeting follow_up"`
	AssignedTo    string     `form:"assigned_to" binding:"omitempty,uuid"`
	AccountID     string     `form:"account_id" binding:"omitempty,uuid"`
	ContactID     string     `form:"contact_id" binding:"omitempty,uuid"`
	DealID        string     `form:"deal_id" binding:"omitempty,uuid"`
	LeadID        string     `form:"lead_id" binding:"omitempty,uuid"`
	TaskSource    string     `form:"task_source" binding:"omitempty,oneof=manual lead_tab pipeline_tab auto_generated"`
	IsSchedule    *bool      `form:"is_schedule" binding:"omitempty"`
	DueDateFrom   *time.Time `form:"due_date_from" time_format:"2006-01-02" binding:"omitempty"`
	DueDateTo     *time.Time `form:"due_date_to" time_format:"2006-01-02" binding:"omitempty"`
	CreatedFrom   *time.Time `form:"created_from" time_format:"2006-01-02" binding:"omitempty"`
	CreatedTo     *time.Time `form:"created_to" time_format:"2006-01-02" binding:"omitempty"`
	ScopedUserIDs []string   `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}

// CreateLeadFromTaskRequest represents request to create a new lead from a task (quick action)
type CreateLeadFromTaskRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=1,max=100"`
	LastName    string `json:"last_name" binding:"omitempty,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"omitempty,max=20"`
	CompanyName string `json:"company_name" binding:"omitempty,max=255"`
	LeadSource  string `json:"lead_source" binding:"required"`
	Notes       string `json:"notes" binding:"omitempty"`
}

// CreateLeadFromTaskResponse represents response for lead created from task
type CreateLeadFromTaskResponse struct {
	Task *TaskResponse `json:"task"`
	Lead interface{}   `json:"lead"` // LeadResponse from lead package
}

// formatNumber formats number with thousand separator
func formatNumber(n float64) string {
	// Convert to int64 to remove decimal places
	amount := int64(n)

	// Handle zero case
	if amount == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if amount < 0 {
		negative = true
		amount = -amount
	}

	// Convert to string
	str := fmt.Sprintf("%d", amount)
	length := len(str)

	// Add thousand separators (dot for Indonesian format)
	var parts []string
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{str[start:i]}, parts...)
	}

	result := strings.Join(parts, ".")
	if negative {
		result = "-" + result
	}

	return result
}
