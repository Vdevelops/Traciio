package task

import "time"

type LeadRefResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type TaskResponse struct {
	ID                       string              `json:"id"`
	Title                    string              `json:"title"`
	Description              string              `json:"description"`
	Type                     string              `json:"type"`
	Status                   string              `json:"status"`
	Priority                 string              `json:"priority"`
	DueDate                  *time.Time          `json:"due_date"`
	CompletedAt              *time.Time          `json:"completed_at"`
	AssignedTo               string              `json:"assigned_to"`
	AssignedUser             *UserRefResponse    `json:"assigned_user,omitempty"`
	AssignedFrom             string              `json:"assigned_from"`
	AssignedFromUser         *UserRefResponse    `json:"assigned_from_user,omitempty"`
	AccountID                string              `json:"account_id"`
	Account                  *AccountRefResponse `json:"account,omitempty"`
	ContactID                string              `json:"contact_id"`
	Contact                  *ContactRefResponse `json:"contact,omitempty"`
	DealID                   string              `json:"deal_id"`
	Deal                     *DealRefResponse    `json:"deal,omitempty"`
	LeadID                   string              `json:"lead_id"`
	Lead                     *LeadRefResponse    `json:"lead,omitempty"`
	TaskSource               string              `json:"task_source"`
	ScheduledStartTime       *time.Time          `json:"scheduled_start_time"`
	ScheduledEndTime         *time.Time          `json:"scheduled_end_time"`
	ScheduledLocation        string              `json:"scheduled_location"`
	IsScheduleTask           bool                `json:"is_schedule_task"`
	QuickActionType          string              `json:"quick_action_type,omitempty"`
	GoogleCalendarEventID    *string             `json:"google_calendar_event_id"`
	GoogleCalendarSyncStatus string              `json:"google_calendar_sync_status"`
	GoogleCalendarSyncedAt   *time.Time          `json:"google_calendar_synced_at"`
	CreatedBy                string              `json:"created_by"`
	CreatedAt                time.Time           `json:"created_at"`
	UpdatedAt                time.Time           `json:"updated_at"`
}

type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type DealRefResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (t *Task) ToTaskResponse() *TaskResponse {
	normalizedStatus := NormalizeStatus(t.Status)
	if t.CompletedAt != nil {
		normalizedStatus = "completed"
	}

	resp := &TaskResponse{
		ID:                 t.ID,
		Title:              t.Title,
		Description:        t.Description,
		Type:               t.Type,
		Status:             normalizedStatus,
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
		resp.AssignedUser = &UserRefResponse{ID: t.AssignedUser.ID, Name: t.AssignedUser.Name, Email: t.AssignedUser.Email, AvatarURL: t.AssignedUser.AvatarURL}
	}
	if t.AssignedFrom != nil {
		resp.AssignedFrom = *t.AssignedFrom
	}
	if t.AssignedFromUser != nil {
		resp.AssignedFromUser = &UserRefResponse{ID: t.AssignedFromUser.ID, Name: t.AssignedFromUser.Name, Email: t.AssignedFromUser.Email, AvatarURL: t.AssignedFromUser.AvatarURL}
	}
	if t.Account != nil {
		resp.Account = &AccountRefResponse{ID: t.Account.ID, Name: t.Account.Name}
	}
	if t.Contact != nil {
		resp.Contact = &ContactRefResponse{ID: t.Contact.ID, Name: t.Contact.Name, Email: t.Contact.Email, Phone: t.Contact.Phone}
	}
	if t.Deal != nil {
		resp.Deal = &DealRefResponse{ID: t.Deal.ID, Title: t.Deal.Title}
	}
	if t.Lead != nil {
		resp.Lead = &LeadRefResponse{ID: t.Lead.ID, FirstName: t.Lead.FirstName, LastName: t.Lead.LastName, Email: t.Lead.Email}
	}
	resp.GoogleCalendarEventID = t.GoogleCalendarEventID
	resp.GoogleCalendarSyncStatus = t.GoogleCalendarSyncStatus
	resp.GoogleCalendarSyncedAt = t.GoogleCalendarSyncedAt
	return resp
}

type CreateTaskRequest struct {
	Title              string     `json:"title" binding:"required,min=3,max=255"`
	Description        string     `json:"description" binding:"omitempty"`
	Type               string     `json:"type" binding:"omitempty,oneof=general call email meeting follow_up"`
	Priority           string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate            *time.Time `json:"due_date" binding:"omitempty"`
	AssignedTo         string     `json:"assigned_to" binding:"omitempty,uuid"`
	AccountID          string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID          string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID             string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID             string     `json:"lead_id" binding:"omitempty,uuid"`
	TaskSource         string     `json:"task_source" binding:"omitempty,oneof=manual lead_tab pipeline_tab auto_generated"`
	ScheduledStartTime *time.Time `json:"scheduled_start_time" binding:"omitempty"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time" binding:"omitempty"`
	ScheduledLocation  string     `json:"scheduled_location" binding:"omitempty,max=500"`
	IsScheduleTask     bool       `json:"is_schedule_task" binding:"omitempty"`
}

type UpdateTaskRequest struct {
	Title              string     `json:"title" binding:"omitempty,min=3,max=255"`
	Description        string     `json:"description" binding:"omitempty"`
	Type               string     `json:"type" binding:"omitempty,oneof=general call email meeting follow_up"`
	Status             string     `json:"status" binding:"omitempty,oneof=pending completed in_progress cancelled submitted approved rejected"`
	Priority           string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate            *time.Time `json:"due_date" binding:"omitempty"`
	AssignedTo         string     `json:"assigned_to" binding:"omitempty,uuid"`
	AccountID          string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID          string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID             string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID             string     `json:"lead_id" binding:"omitempty,uuid"`
	ScheduledStartTime *time.Time `json:"scheduled_start_time" binding:"omitempty"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time" binding:"omitempty"`
	ScheduledLocation  string     `json:"scheduled_location" binding:"omitempty,max=500"`
	IsScheduleTask     *bool      `json:"is_schedule_task" binding:"omitempty"`
}

type AssignTaskRequest struct {
	AssignedTo string `json:"assigned_to" binding:"required,uuid"`
}

type ListTasksRequest struct {
	Page          int        `form:"page" binding:"omitempty,min=1"`
	PerPage       int        `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search        string     `form:"search" binding:"omitempty"`
	Status        string     `form:"status" binding:"omitempty"`
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
	ScopedUserIDs []string   `form:"-" json:"-"`
}

type CreateLeadFromTaskRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=1,max=100"`
	LastName    string `json:"last_name" binding:"omitempty,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"omitempty,max=20"`
	CompanyName string `json:"company_name" binding:"omitempty,max=255"`
	LeadSource  string `json:"lead_source" binding:"required"`
	Notes       string `json:"notes" binding:"omitempty"`
}

type CreateLeadFromTaskResponse struct {
	Task *TaskResponse `json:"task"`
	Lead interface{}   `json:"lead"`
}
