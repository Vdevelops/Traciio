package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/internal/service/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/google_calendar"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

// ScheduleServiceInterface defines interface for schedule service to avoid circular dependency.
type ScheduleServiceInterface interface {
	CreateScheduleFromTask(taskID string, createdBy string) error
	DeleteScheduleByTaskID(taskID string) error
}

var (
	ErrTaskNotFound                  = errors.New("task not found")
	ErrUserNotFound                  = errors.New("user not found")
	ErrAccountNotFound               = errors.New("account not found")
	ErrContactNotFound               = errors.New("contact not found")
	ErrDealNotFound                  = errors.New("deal not found")
	ErrLeadNotFound                  = errors.New("lead not found")
	ErrTasksRestrictedContext        = errors.New("tasks restricted context")
	ErrReminderNotFound              = errors.New("reminder not found")
	ErrTaskAlreadyCompleted          = errors.New("task already completed")
	ErrCannotMarkCompletedInProgress = errors.New("cannot mark completed task as in progress")
)

const autoReminderPrefix = "AUTO_REMINDER:"

type Service struct {
	taskRepo            interfaces.TaskRepository
	reminderRepo        interfaces.ReminderRepository
	userRepo            interfaces.UserRepository
	accountRepo         interfaces.AccountRepository
	contactRepo         interfaces.ContactRepository
	dealRepo            interfaces.DealRepository
	leadRepo            interfaces.LeadRepository
	scheduleService     ScheduleServiceInterface
	googleCalendarToken *google_calendar_token.Service
	cacheService        *cache.TaskCacheService
	eventHelper         *domainevents.Helper
}

func NewService(
	taskRepo interfaces.TaskRepository,
	reminderRepo interfaces.ReminderRepository,
	userRepo interfaces.UserRepository,
	accountRepo interfaces.AccountRepository,
	contactRepo interfaces.ContactRepository,
	dealRepo interfaces.DealRepository,
	leadRepo interfaces.LeadRepository,
	eventHelper *domainevents.Helper,
) *Service {
	return &Service{
		taskRepo:     taskRepo,
		reminderRepo: reminderRepo,
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		contactRepo:  contactRepo,
		dealRepo:     dealRepo,
		leadRepo:     leadRepo,
		cacheService: cache.NewTaskCacheService(nil),
		eventHelper:  eventHelper,
	}
}

// SetScheduleService sets the schedule service for auto-creating schedules from tasks
func (s *Service) SetScheduleService(scheduleService ScheduleServiceInterface) {
	s.scheduleService = scheduleService
}

// SetGoogleCalendarTokenService sets the Google Calendar token service for auto-syncing tasks
func (s *Service) SetGoogleCalendarTokenService(googleCalendarToken *google_calendar_token.Service) {
	s.googleCalendarToken = googleCalendarToken
}

// cachedTaskListResult for msgpack serialization
type cachedTaskListResult struct {
	Tasks      []task.TaskResponse `msgpack:"tasks"`
	Pagination *PaginationResult   `msgpack:"pagination"`
}

// ListTasks returns a list of tasks with pagination
func (s *Service) ListTasks(req *task.ListTasksRequest) ([]task.TaskResponse, *PaginationResult, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	// Build cache key from request filters (includes ScopedUserIDs to prevent cross-scope cache pollution)
	filterMap := map[string]interface{}{
		"assigned_to": req.AssignedTo,
		"status":      req.Status,
		"priority":    req.Priority,
		"type":        req.Type,
		"account_id":  req.AccountID,
		"contact_id":  req.ContactID,
		"deal_id":     req.DealID,
		"lead_id":     req.LeadID,
		"task_source": req.TaskSource,
	}
	if len(req.ScopedUserIDs) > 0 {
		filterMap["scoped_user_ids"] = fmt.Sprintf("%v", req.ScopedUserIDs)
	}

	// Try cache first
	var cachedResult cachedTaskListResult
	if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && len(cachedResult.Tasks) > 0 {
		return cachedResult.Tasks, cachedResult.Pagination, nil
	}

	tasks, total, err := s.taskRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]task.TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = *t.ToTaskResponse()
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	// Cache the result
	result := cachedTaskListResult{
		Tasks:      responses,
		Pagination: pagination,
	}
	_ = s.cacheService.SetList(page, perPage, filterMap, result)

	return responses, pagination, nil
}

// GetTaskByID returns a task by ID
func (s *Service) GetTaskByID(id string) (*task.TaskResponse, error) {
	// Try cache first
	var cachedResponse task.TaskResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	response := t.ToTaskResponse()

	// Cache the result
	_ = s.cacheService.SetDetail(id, response)

	return response, nil
}

// CreateTask creates a new task.
func (s *Service) CreateTask(req *task.CreateTaskRequest, createdBy string) (*task.TaskResponse, error) {
	// CRM Enhancement Phase 1: Task restriction
	// Task must be created within a context (Lead, Deal, Account, or Contact)
	if req.LeadID == "" && req.DealID == "" && req.AccountID == "" && req.ContactID == "" {
		return nil, ErrTasksRestrictedContext
	}

	// Validate assigned user if provided
	if req.AssignedTo != "" {
		_, err := s.userRepo.FindByID(req.AssignedTo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
	}

	// Validate account if provided
	if req.AccountID != "" {
		_, err := s.accountRepo.FindByID(req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}
	}

	// Validate contact if provided
	if req.ContactID != "" {
		_, err := s.contactRepo.FindByID(req.ContactID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrContactNotFound
			}
			return nil, err
		}
	}

	// Validate deal if provided
	if req.DealID != "" {
		_, err := s.dealRepo.FindByID(req.DealID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrDealNotFound
			}
			return nil, err
		}
	}

	// Set defaults
	taskType := req.Type
	if taskType == "" {
		taskType = "general"
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	var assignedToPtr *string
	if req.AssignedTo != "" {
		assignedToPtr = &req.AssignedTo
	}

	var accountIDPtr *string
	if req.AccountID != "" {
		accountIDPtr = &req.AccountID
	}

	var contactIDPtr *string
	if req.ContactID != "" {
		contactIDPtr = &req.ContactID
	}

	var dealIDPtr *string
	if req.DealID != "" {
		dealIDPtr = &req.DealID
	}

	var leadIDPtr *string
	if req.LeadID != "" {
		// Validate lead exists
		if s.leadRepo != nil {
			_, err := s.leadRepo.FindByID(req.LeadID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrLeadNotFound
				}
				return nil, err
			}
		}
		leadIDPtr = &req.LeadID
	}

	// Set assigned_from if task is assigned
	var assignedFromPtr *string
	if req.AssignedTo != "" && createdBy != "" {
		assignedFromPtr = &createdBy
	}

	// Set task source (default to manual)
	taskSource := req.TaskSource
	if taskSource == "" {
		taskSource = "manual"
	}

	t := &task.Task{
		Title:              req.Title,
		Description:        req.Description,
		Type:               taskType,
		Status:             "pending",
		Priority:           priority,
		DueDate:            req.DueDate,
		AssignedTo:         assignedToPtr,
		AssignedFrom:       assignedFromPtr,
		AccountID:          accountIDPtr,
		ContactID:          contactIDPtr,
		DealID:             dealIDPtr,
		LeadID:             leadIDPtr,
		TaskSource:         taskSource,
		ScheduledStartTime: req.ScheduledStartTime,
		ScheduledEndTime:   req.ScheduledEndTime,
		ScheduledLocation:  req.ScheduledLocation,
		IsScheduleTask:     req.IsScheduleTask,
		CreatedBy:          createdBy,
	}

	if err := s.taskRepo.Create(t); err != nil {
		return nil, err
	}

	// Invalidate cache on write
	s.cacheService.InvalidateOnWrite("")

	// Reload to get relations
	t, err := s.taskRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	// Emit task created event
	if s.eventHelper != nil {
		assignedTo := ""
		if t.AssignedTo != nil {
			assignedTo = *t.AssignedTo
		}
		assignedFrom := ""
		if t.AssignedFrom != nil {
			assignedFrom = *t.AssignedFrom
		}
		accountID := ""
		if t.AccountID != nil {
			accountID = *t.AccountID
		}
		dealID := ""
		if t.DealID != nil {
			dealID = *t.DealID
		}

		s.eventHelper.EmitTaskCreated(&domainevents.TaskCreatedEvent{
			TaskID:       t.ID,
			Title:        t.Title,
			Description:  t.Description,
			Status:       t.Status,
			Priority:     t.Priority,
			DueDate:      t.DueDate,
			AssignedTo:   assignedTo,
			AssignedFrom: assignedFrom,
			AccountID:    accountID,
			DealID:       dealID,
			CreatedBy:    createdBy,
			CreatedAt:    t.CreatedAt,
		}, createdBy)

		// Emit task assigned event if task is assigned
		if t.AssignedTo != nil && *t.AssignedTo != "" {
			s.eventHelper.EmitTaskAssigned(&domainevents.TaskAssignedEvent{
				TaskID:      t.ID,
				Title:       t.Title,
				NewAssignee: assignedTo,
				AssignedBy:  createdBy,
				AssignedAt:  t.CreatedAt,
			}, createdBy)
		}
	}

	// Auto-create schedule if task has due_date and assigned_to
	if s.scheduleService != nil && t.DueDate != nil && t.AssignedTo != nil {
		go func() {
			_ = s.scheduleService.CreateScheduleFromTask(t.ID, createdBy)
		}()
	}

	// Auto-sync D-1 and day-H in-app reminders for tasks with due dates.
	s.syncAutoRemindersForTask(t, createdBy)

	return t.ToTaskResponse(), nil
}

// UpdateTask updates a task
func (s *Service) UpdateTask(id string, req *task.UpdateTaskRequest) (*task.TaskResponse, error) {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		t.Title = req.Title
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if req.Type != "" {
		t.Type = req.Type
	}
	if req.Status != "" {
		t.Status = task.NormalizeStatus(req.Status)
		if t.Status == "completed" && t.CompletedAt == nil {
			now := time.Now()
			t.CompletedAt = &now
		} else if t.Status != "completed" {
			t.CompletedAt = nil
		}
	}
	if req.Priority != "" {
		t.Priority = req.Priority
	}
	// Handle due_date update
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}
	if req.AssignedTo != "" {
		// Validate user exists
		_, err := s.userRepo.FindByID(req.AssignedTo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
		t.AssignedTo = &req.AssignedTo
	}
	if req.AccountID != "" {
		// Validate account exists
		_, err := s.accountRepo.FindByID(req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}
		t.AccountID = &req.AccountID
	}
	if req.ContactID != "" {
		// Validate contact exists
		_, err := s.contactRepo.FindByID(req.ContactID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrContactNotFound
			}
			return nil, err
		}
		t.ContactID = &req.ContactID
	}
	if req.DealID != "" {
		// Validate deal exists
		_, err := s.dealRepo.FindByID(req.DealID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrDealNotFound
			}
			return nil, err
		}
		t.DealID = &req.DealID
	}
	if req.LeadID != "" {
		// Validate lead exists
		if s.leadRepo != nil {
			_, err := s.leadRepo.FindByID(req.LeadID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrLeadNotFound
				}
				return nil, err
			}
		}
		t.LeadID = &req.LeadID
	}
	// Update schedule fields
	if req.ScheduledStartTime != nil {
		t.ScheduledStartTime = req.ScheduledStartTime
	}
	if req.ScheduledEndTime != nil {
		t.ScheduledEndTime = req.ScheduledEndTime
	}
	if req.ScheduledLocation != "" {
		t.ScheduledLocation = req.ScheduledLocation
	}
	if req.IsScheduleTask != nil {
		t.IsScheduleTask = *req.IsScheduleTask
	}

	if err := s.taskRepo.Update(t); err != nil {
		return nil, err
	}

	// Invalidate cache on write
	s.cacheService.InvalidateOnWrite(id)

	// Reload to get relations
	t, err = s.taskRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	// Auto-create/update/delete schedule based on task state
	if s.scheduleService != nil {
		go func() {
			_ = s.scheduleService.CreateScheduleFromTask(t.ID, t.CreatedBy)
		}()
	}

	// Auto-sync task to Google Calendar if task has due_date and assigned_to
	if s.googleCalendarToken != nil && t.DueDate != nil && t.AssignedTo != nil {
		go func() {
			ctx := context.Background()
			_ = s.SyncTaskToGoogleCalendar(ctx, t.ID, *t.AssignedTo)
		}()
	}

	// Keep automatic reminders aligned with the latest task due date.
	s.syncAutoRemindersForTask(t, t.CreatedBy)

	return t.ToTaskResponse(), nil
}

// syncAutoRemindersForTask ensures two automatic in-app reminders exist for a task:
// one at D-1 and one at day-H. Existing auto reminders are replaced to keep it idempotent.
func (s *Service) syncAutoRemindersForTask(t *task.Task, createdBy string) {
	if s.reminderRepo == nil || t == nil || t.DueDate == nil {
		return
	}

	existing, err := s.reminderRepo.FindByTaskID(t.ID)
	if err != nil {
		return
	}

	for i := range existing {
		if strings.HasPrefix(existing[i].Message, autoReminderPrefix) {
			_ = s.reminderRepo.Delete(existing[i].ID)
		}
	}

	now := time.Now()
	targets := []struct {
		when    time.Time
		message string
	}{
		{
			when:    t.DueDate.Add(-24 * time.Hour),
			message: autoReminderPrefix + " D-1 task reminder",
		},
		{
			when:    *t.DueDate,
			message: autoReminderPrefix + " Day-H task reminder",
		},
	}

	for i := range targets {
		if !targets[i].when.After(now) {
			continue
		}

		_ = s.reminderRepo.Create(&reminder.Reminder{
			TaskID:       t.ID,
			RemindAt:     targets[i].when,
			ReminderType: "in_app",
			Message:      targets[i].message,
			CreatedBy:    createdBy,
		})
	}
}

// DeleteTask deletes a task and associated schedule
func (s *Service) DeleteTask(id string) error {
	_, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	// Delete associated schedule if exists
	if s.scheduleService != nil {
		// Delete all schedules associated with this task
		_ = s.scheduleService.DeleteScheduleByTaskID(id)
	}

	// Invalidate cache on write
	s.cacheService.InvalidateOnWrite(id)

	return s.taskRepo.Delete(id)
}

// AssignTask assigns a task to a user
func (s *Service) AssignTask(id string, req *task.AssignTaskRequest, assignedBy string) (*task.TaskResponse, error) {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	// Validate user exists
	_, err = s.userRepo.FindByID(req.AssignedTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	t.AssignedTo = &req.AssignedTo

	if assignedBy != "" {
		t.AssignedFrom = &assignedBy
	}

	if err := s.taskRepo.Update(t); err != nil {
		return nil, err
	}

	// Reload to get relations
	t, err = s.taskRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	return t.ToTaskResponse(), nil
}

// CompleteTask marks a task as completed
func (s *Service) CompleteTask(id string) (*task.TaskResponse, error) {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.NormalizeStatus(t.Status) == "completed" {
		return nil, ErrTaskAlreadyCompleted
	}

	t.Status = "completed"
	now := time.Now()
	t.CompletedAt = &now

	if err := s.taskRepo.Update(t); err != nil {
		return nil, err
	}

	// Reload to get relations
	t, err = s.taskRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	return t.ToTaskResponse(), nil
}

// MarkInProgress marks a task as in progress
func (s *Service) MarkInProgress(id string) (*task.TaskResponse, error) {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.NormalizeStatus(t.Status) == "completed" {
		return nil, ErrCannotMarkCompletedInProgress
	}

	if task.NormalizeStatus(t.Status) == "pending" {
		return t.ToTaskResponse(), nil
	}

	t.Status = "pending"
	t.CompletedAt = nil

	if err := s.taskRepo.Update(t); err != nil {
		return nil, err
	}

	// Reload to get relations
	t, err = s.taskRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	return t.ToTaskResponse(), nil
}

// GetMyTasks returns tasks assigned to the logged-in user.
// GetMyTasks returns tasks assigned to the logged-in user
// This method ensures that search and all filters only apply to tasks assigned to the user
// Search parameter is supported and will only search within the user's assigned tasks
func (s *Service) GetMyTasks(userID string, req *task.ListTasksRequest) ([]task.TaskResponse, *PaginationResult, error) {
	// Force AssignedTo to be the logged-in user to ensure security
	// This ensures that search and all filters only apply to user's tasks
	req.AssignedTo = userID
	return s.ListTasks(req)
}

// ListReminders returns a list of reminders with pagination
func (s *Service) ListReminders(req *reminder.ListRemindersRequest) ([]reminder.ReminderResponse, *PaginationResult, error) {
	reminders, total, err := s.reminderRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]reminder.ReminderResponse, len(reminders))
	for i, rem := range reminders {
		responses[i] = *rem.ToReminderResponse()
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return responses, pagination, nil
}

// GetReminderByID returns a reminder by ID
func (s *Service) GetReminderByID(id string) (*reminder.ReminderResponse, error) {
	rem, err := s.reminderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReminderNotFound
		}
		return nil, err
	}

	return rem.ToReminderResponse(), nil
}

// CreateReminder creates a new reminder
func (s *Service) CreateReminder(req *reminder.CreateReminderRequest, createdBy string) (*reminder.ReminderResponse, error) {
	// Validate task exists
	_, err := s.taskRepo.FindByID(req.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	reminderType := req.ReminderType
	if reminderType == "" {
		reminderType = "in_app"
	}

	rem := &reminder.Reminder{
		TaskID:       req.TaskID,
		RemindAt:     req.RemindAt,
		ReminderType: reminderType,
		Message:      req.Message,
		CreatedBy:    createdBy,
	}

	if err := s.reminderRepo.Create(rem); err != nil {
		return nil, err
	}

	// Reload to get relations
	rem, err = s.reminderRepo.FindByID(rem.ID)
	if err != nil {
		return nil, err
	}

	return rem.ToReminderResponse(), nil
}

// UpdateReminder updates a reminder
func (s *Service) UpdateReminder(id string, req *reminder.UpdateReminderRequest) (*reminder.ReminderResponse, error) {
	rem, err := s.reminderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReminderNotFound
		}
		return nil, err
	}

	if req.RemindAt != nil {
		rem.RemindAt = *req.RemindAt
	}
	if req.ReminderType != "" {
		rem.ReminderType = req.ReminderType
	}
	if req.Message != "" {
		rem.Message = req.Message
	}

	if err := s.reminderRepo.Update(rem); err != nil {
		return nil, err
	}

	// Reload to get relations
	rem, err = s.reminderRepo.FindByID(rem.ID)
	if err != nil {
		return nil, err
	}

	return rem.ToReminderResponse(), nil
}

// DeleteReminder deletes a reminder
func (s *Service) DeleteReminder(id string) error {
	_, err := s.reminderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReminderNotFound
		}
		return err
	}

	return s.reminderRepo.Delete(id)
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// SyncTaskToGoogleCalendar syncs a task to Google Calendar
func (s *Service) SyncTaskToGoogleCalendar(ctx context.Context, taskID string, userID string) error {
	// Get task
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	if t.DueDate == nil {
		return nil
	}

	calendarService, err := s.googleCalendarToken.GetCalendarService(ctx, userID)
	if err != nil {
		t.GoogleCalendarSyncStatus = "sync_failed"
		s.taskRepo.Update(t)
		return err
	}

	// Build event title
	title := t.Title
	if t.Type != "" {
		title = "[" + t.Type + "] " + title
	}

	// Build event description
	description := ""
	if t.Description != "" {
		description = t.Description
	}
	if t.Priority != "" {
		if description != "" {
			description += "\n\n"
		}
		description += "Priority: " + t.Priority
	}
	if t.Status != "" {
		if description != "" {
			description += "\n"
		}
		description += "Status: " + t.Status
	}

	// Build Google Calendar event
	event := google_calendar.BuildEventFromSchedule(
		title,
		description,
		*t.DueDate,
		nil, // Use default reminder (15 minutes)
	)

	// Create or update event in Google Calendar
	var createdEvent *calendar.Event
	if t.GoogleCalendarEventID != nil {
		// Update existing event
		event.Id = *t.GoogleCalendarEventID
		createdEvent, err = calendarService.Events.Update("primary", *t.GoogleCalendarEventID, event).Context(ctx).Do()
	} else {
		// Create new event
		createdEvent, err = calendarService.Events.Insert("primary", event).Context(ctx).Do()
	}

	if err != nil {
		t.GoogleCalendarSyncStatus = "sync_failed"
		s.taskRepo.Update(t)
		return err
	}

	// Update task with Google Calendar event ID
	now := time.Now()
	t.GoogleCalendarEventID = &createdEvent.Id
	t.GoogleCalendarSyncStatus = "synced"
	t.GoogleCalendarSyncedAt = &now

	if err := s.taskRepo.Update(t); err != nil {
		return err
	}

	return nil
}

// CreateLeadFromTask creates a new lead and links it to the specified task (quick action).
// This is a transactional operation: create lead → update task.lead_id.
func (s *Service) CreateLeadFromTask(taskID string, req *task.CreateLeadFromTaskRequest, createdBy string) (*task.CreateLeadFromTaskResponse, error) {
	// 1. Find the task
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// 2. Build lead entity from request + task context
	newLead := &lead.Lead{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		CompanyName: req.CompanyName,
		LeadSource:  req.LeadSource,
		LeadStatus:  "new",
		Notes:       req.Notes,
		Country:     "Indonesia",
		CreatedBy:   createdBy,
	}

	// Pre-populate assigned_to from task if available
	if t.AssignedTo != nil && *t.AssignedTo != "" {
		newLead.AssignedTo = t.AssignedTo
	}

	// Pre-populate account from task if available
	if t.AccountID != nil && *t.AccountID != "" {
		newLead.AccountID = t.AccountID
	}

	// 3. Create lead
	if err := s.leadRepo.Create(newLead); err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}

	// 4. Link lead to task
	t.LeadID = &newLead.ID
	if err := s.taskRepo.Update(t); err != nil {
		return nil, fmt.Errorf("failed to link lead to task: %w", err)
	}

	// 5. Reload task to get updated relations
	t, err = s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// 6. Reload lead to get relations
	createdLead, err := s.leadRepo.FindByID(newLead.ID)
	if err != nil {
		return nil, err
	}

	return &task.CreateLeadFromTaskResponse{
		Task: t.ToTaskResponse(),
		Lead: createdLead.ToLeadResponse(),
	}, nil
}
