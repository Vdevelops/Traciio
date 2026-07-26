package schedule

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/internal/service/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/google_calendar"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

var (
	ErrScheduleNotFound        = errors.New("schedule not found")
	ErrTaskNotFound            = errors.New("task not found")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidTaskAssignment   = errors.New("task is not assigned to the user")
	ErrGoogleCalendarNotSynced = errors.New("schedule is not synced to Google Calendar")
)

type Service struct {
	scheduleRepo        interfaces.ScheduleRepository
	taskRepo            interfaces.TaskRepository
	userRepo            interfaces.UserRepository
	googleCalendarToken *google_calendar_token.Service
	cacheService        *cache.ScheduleCacheService
}

func NewService(
	scheduleRepo interfaces.ScheduleRepository,
	taskRepo interfaces.TaskRepository,
	userRepo interfaces.UserRepository,
	googleCalendarToken *google_calendar_token.Service,
) *Service {
	return &Service{
		scheduleRepo:        scheduleRepo,
		taskRepo:            taskRepo,
		userRepo:            userRepo,
		googleCalendarToken: googleCalendarToken,
		cacheService:        cache.NewScheduleCacheService(nil),
	}
}

// cachedScheduleListResult for msgpack serialization
type cachedScheduleListResult struct {
	Schedules  []schedule.ScheduleResponse `msgpack:"schedules"`
	Pagination *PaginationResult           `msgpack:"pagination"`
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// ListSchedules returns a list of schedules with pagination
func (s *Service) ListSchedules(req *schedule.ListSchedulesRequest) ([]schedule.ScheduleResponse, *PaginationResult, error) {
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

	// Build cache key from request filters (avoid pointer values)
	filterMap := map[string]interface{}{
		"search":                      req.Search,
		"status":                      req.Status,
		"task_id":                     req.TaskID,
		"user_id":                     req.UserID,
		"google_calendar_sync_status": req.GoogleCalendarSyncStatus,
		"scoped_user_ids":             scopedUserIDsCacheKey(req.ScopedUserIDs),
	}
	if req.ScheduledAtFrom != nil {
		filterMap["scheduled_at_from"] = req.ScheduledAtFrom.Format(time.RFC3339)
	}
	if req.ScheduledAtTo != nil {
		filterMap["scheduled_at_to"] = req.ScheduledAtTo.Format(time.RFC3339)
	}

	// Try cache first
	var cachedResult cachedScheduleListResult
	if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && cachedResult.Pagination != nil {
		return cachedResult.Schedules, cachedResult.Pagination, nil
	}

	schedules, total, err := s.scheduleRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]schedule.ScheduleResponse, len(schedules))
	for i, sch := range schedules {
		responses[i] = *sch.ToScheduleResponse()
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	// Cache the result
	_ = s.cacheService.SetList(page, perPage, filterMap, cachedScheduleListResult{
		Schedules:  responses,
		Pagination: pagination,
	})

	return responses, pagination, nil
}

func scopedUserIDsCacheKey(ids []string) string {
	if ids == nil {
		return "global"
	}
	if len(ids) == 0 {
		return "none"
	}
	copied := append([]string(nil), ids...)
	sort.Strings(copied)
	return strings.Join(copied, ",")
}

// GetScheduleByID returns a schedule by ID
func (s *Service) GetScheduleByID(id string) (*schedule.ScheduleResponse, error) {
	// Try cache first
	var cachedResponse schedule.ScheduleResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	sch, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	resp := sch.ToScheduleResponse()
	_ = s.cacheService.SetDetail(id, resp)
	return resp, nil
}

// CreateSchedule creates a new schedule
func (s *Service) CreateSchedule(req *schedule.CreateScheduleRequest, createdBy string) (*schedule.ScheduleResponse, error) {
	var userID string
	var taskID string

	// If task_id is provided, validate task and get user_id from task
	if req.TaskID != "" {
		t, err := s.taskRepo.FindByID(req.TaskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrTaskNotFound
			}
			return nil, err
		}

		if t.AssignedTo == nil {
			return nil, ErrInvalidTaskAssignment
		}
		userID = *t.AssignedTo
		taskID = req.TaskID
	} else {
		userID = createdBy
		taskID = ""
	}

	// Validate user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var descriptionPtr *string
	if req.Description != "" {
		descriptionPtr = &req.Description
	}

	sch := &schedule.Schedule{
		UserID:                   userID,
		Title:                    req.Title,
		Description:              descriptionPtr,
		ScheduledAt:              req.ScheduledAt,
		Status:                   "pending",
		ReminderMinutesBefore:    req.ReminderMinutesBefore,
		GoogleCalendarSyncStatus: "not_synced",
		CreatedBy:                createdBy,
	}

	// Set TaskID only if provided
	if taskID != "" {
		sch.TaskID = &taskID
	}

	if err := s.scheduleRepo.Create(sch); err != nil {
		return nil, err
	}

	sch, err = s.scheduleRepo.FindByID(sch.ID)
	if err != nil {
		return nil, err
	}

	resp := sch.ToScheduleResponse()
	_ = s.cacheService.InvalidateOnWrite(sch.ID)
	_ = s.cacheService.SetDetail(sch.ID, resp)
	return resp, nil
}

// UpdateSchedule updates a schedule
func (s *Service) UpdateSchedule(id string, req *schedule.UpdateScheduleRequest) (*schedule.ScheduleResponse, error) {
	sch, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	if req.Title != "" {
		sch.Title = req.Title
	}
	if req.Description != "" {
		sch.Description = &req.Description
	} else if req.Description == "" && sch.Description != nil {
		sch.Description = nil
	}
	if req.ScheduledAt != nil {
		sch.ScheduledAt = *req.ScheduledAt
	}
	if req.Status != "" {
		sch.Status = req.Status
	}
	if req.ReminderMinutesBefore != nil {
		sch.ReminderMinutesBefore = req.ReminderMinutesBefore
	}

	if err := s.scheduleRepo.Update(sch); err != nil {
		return nil, err
	}

	sch, err = s.scheduleRepo.FindByID(sch.ID)
	if err != nil {
		return nil, err
	}

	resp := sch.ToScheduleResponse()
	_ = s.cacheService.InvalidateOnWrite(sch.ID)
	_ = s.cacheService.SetDetail(sch.ID, resp)
	return resp, nil
}

// DeleteSchedule deletes a schedule.
func (s *Service) DeleteSchedule(id string) error {
	_, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrScheduleNotFound
		}
		return err
	}

	if err := s.scheduleRepo.Delete(id); err != nil {
		return err
	}
	_ = s.cacheService.InvalidateOnWrite(id)
	return nil
}

// SyncToGoogleCalendar syncs a schedule to Google Calendar.
// Creates or updates the event based on whether GoogleCalendarEventID exists.
func (s *Service) SyncToGoogleCalendar(ctx context.Context, id string, userID string) (*schedule.GoogleCalendarSyncResponse, error) {
	sch, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	var t *task.Task
	if sch.TaskID != nil {
		taskEntity, err := s.taskRepo.FindByID(*sch.TaskID)
		if err != nil {
			return nil, err
		}
		t = taskEntity
	}

	calendarService, err := s.googleCalendarToken.GetCalendarService(ctx, userID)
	if err != nil {
		sch.GoogleCalendarSyncStatus = "sync_failed"
		s.scheduleRepo.Update(sch)
		return nil, err
	}

	description := ""
	if sch.Description != nil {
		description = *sch.Description
	}
	if t != nil && t.Description != "" {
		if description != "" {
			description += "\n\n"
		}
		description += "Task: " + t.Description
	}

	event := google_calendar.BuildEventFromSchedule(
		sch.Title,
		description,
		sch.ScheduledAt,
		sch.ReminderMinutesBefore,
	)

	var createdEvent *calendar.Event
	if sch.GoogleCalendarEventID != nil {
		event.Id = *sch.GoogleCalendarEventID

		// Google Calendar API requires UseDefault=false when using Overrides
		if event.Reminders == nil {
			event.Reminders = &calendar.EventReminders{
				UseDefault: true,
				Overrides:  []*calendar.EventReminder{},
			}
		} else {
			event.Reminders.UseDefault = false
			if event.Reminders.Overrides == nil {
				event.Reminders.Overrides = []*calendar.EventReminder{}
			}
		}

		createdEvent, err = calendarService.Events.Update("primary", *sch.GoogleCalendarEventID, event).
			SendUpdates("none").
			Context(ctx).Do()
	} else {
		// Ensure UseDefault=false is sent when using Overrides
		if event.Reminders != nil && len(event.Reminders.Overrides) > 0 {
			event.Reminders.UseDefault = false
			if event.Reminders.ForceSendFields == nil {
				event.Reminders.ForceSendFields = []string{}
			}
			useDefaultFound := false
			for _, field := range event.Reminders.ForceSendFields {
				if field == "UseDefault" {
					useDefaultFound = true
					break
				}
			}
			if !useDefaultFound {
				event.Reminders.ForceSendFields = append(event.Reminders.ForceSendFields, "UseDefault")
			}
		}

		createdEvent, err = calendarService.Events.Insert("primary", event).Context(ctx).Do()
	}

	if err != nil {
		sch.GoogleCalendarSyncStatus = "sync_failed"
		s.scheduleRepo.Update(sch)
		return nil, err
	}

	now := time.Now()
	sch.GoogleCalendarEventID = &createdEvent.Id
	sch.GoogleCalendarSyncStatus = "synced"
	sch.GoogleCalendarSyncedAt = &now
	if createdEvent.HtmlLink != "" {
		sch.GoogleCalendarEventLink = &createdEvent.HtmlLink
	}

	if err := s.scheduleRepo.Update(sch); err != nil {
		return nil, fmt.Errorf("failed to update schedule after sync: %w", err)
	}
	_ = s.cacheService.InvalidateOnWrite(sch.ID)

	eventURL := google_calendar.GetEventURL(createdEvent.Id)
	if createdEvent.HtmlLink != "" {
		eventURL = createdEvent.HtmlLink
	}

	return &schedule.GoogleCalendarSyncResponse{
		ScheduleID:            sch.ID,
		GoogleCalendarEventID: createdEvent.Id,
		EventURL:              eventURL,
		SyncStatus:            "synced",
	}, nil
}

// UnsyncFromGoogleCalendar removes sync from Google Calendar
func (s *Service) UnsyncFromGoogleCalendar(ctx context.Context, id string, userID string) (*schedule.ScheduleResponse, error) {
	sch, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	if sch.GoogleCalendarSyncStatus != "synced" {
		return nil, nil
	}

	if sch.GoogleCalendarEventID != nil {
		calendarService, err := s.googleCalendarToken.GetCalendarService(ctx, userID)
		if err == nil {
			_ = calendarService.Events.Delete("primary", *sch.GoogleCalendarEventID).Context(ctx).Do()
		}
	}

	sch.GoogleCalendarEventID = nil
	sch.GoogleCalendarSyncStatus = "not_synced"
	sch.GoogleCalendarSyncedAt = nil
	sch.GoogleCalendarEventLink = nil

	if err := s.scheduleRepo.Update(sch); err != nil {
		return nil, err
	}

	sch, err = s.scheduleRepo.FindByID(sch.ID)
	if err != nil {
		return nil, err
	}

	resp := sch.ToScheduleResponse()
	_ = s.cacheService.InvalidateOnWrite(sch.ID)
	_ = s.cacheService.SetDetail(sch.ID, resp)
	return resp, nil
}

// CreateScheduleFromTask automatically creates or updates a schedule from a task.
// Called when a task is created or updated with a due date and assigned user.
func (s *Service) CreateScheduleFromTask(taskID string, createdBy string) error {
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	// Delete existing schedules if task no longer has due_date or assigned_to
	if t.DueDate == nil || t.AssignedTo == nil {
		existingSchedules, _, err := s.scheduleRepo.List(&schedule.ListSchedulesRequest{
			TaskID:  taskID,
			PerPage: 1,
		})
		if err == nil && len(existingSchedules) > 0 {
			for _, existingSch := range existingSchedules {
				_ = s.DeleteSchedule(existingSch.ID)
			}
		}
		return nil
	}

	userID := *t.AssignedTo

	_, err = s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if schedule already exists for this task
	existingSchedules, _, err := s.scheduleRepo.List(&schedule.ListSchedulesRequest{
		TaskID:  taskID,
		PerPage: 1,
	})

	var sch *schedule.Schedule
	isUpdate := false

	if err == nil && len(existingSchedules) > 0 {
		sch = &existingSchedules[0]
		isUpdate = true
	} else {
		sch = &schedule.Schedule{
			TaskID:                   &taskID,
			UserID:                   userID,
			GoogleCalendarSyncStatus: "not_synced",
			CreatedBy:                createdBy,
		}
	}

	title := "Reminder: " + t.Title
	if len(title) > 255 {
		title = title[:252] + "..."
	}
	sch.Title = title

	var descriptionPtr *string
	if t.Description != "" {
		desc := "Task: " + t.Description
		descriptionPtr = &desc
	}
	sch.Description = descriptionPtr

	// Use task's due_date directly as scheduled_at
	sch.ScheduledAt = *t.DueDate

	reminderMinutesBefore := 60
	sch.ReminderMinutesBefore = &reminderMinutesBefore

	status := "pending"
	if t.Status == "completed" {
		status = "completed"
	} else if t.Status == "cancelled" {
		status = "cancelled"
	} else if sch.ScheduledAt.Before(time.Now()) {
		status = "confirmed"
	}
	sch.Status = status

	sch.UserID = userID

	if isUpdate {
		if err := s.scheduleRepo.Update(sch); err != nil {
			return err
		}
	} else {
		if err := s.scheduleRepo.Create(sch); err != nil {
			return err
		}
	}
	_ = s.cacheService.InvalidateOnWrite(sch.ID)

	return nil
}

// DeleteScheduleByTaskID deletes all schedules associated with a task
func (s *Service) DeleteScheduleByTaskID(taskID string) error {
	// Use pagination to handle large numbers of schedules efficiently
	const batchSize = 100
	page := 1

	for {
		existingSchedules, total, err := s.scheduleRepo.List(&schedule.ListSchedulesRequest{
			TaskID:  taskID,
			Page:    page,
			PerPage: batchSize,
		})
		if err != nil {
			return err
		}

		// Delete schedules in batch
		for _, existingSch := range existingSchedules {
			// Use direct repository delete instead of service DeleteSchedule
			// to avoid unnecessary Google Calendar API calls for each schedule
			if err := s.scheduleRepo.Delete(existingSch.ID); err != nil {
				// Log error but continue with other schedules
				// In production, consider using a logger here
				continue
			}
			_ = s.cacheService.InvalidateOnWrite(existingSch.ID)
		}

		// Check if there are more schedules to delete
		if len(existingSchedules) < batchSize || int64(page*batchSize) >= total {
			break
		}
		page++
	}

	return nil
}
