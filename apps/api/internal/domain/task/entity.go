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
	Status           string      `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, completed
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

func NormalizeStatus(status string) string {
	switch status {
	case "completed", "approved", "cancelled", "rejected":
		return "completed"
	default:
		return "pending"
	}
}
