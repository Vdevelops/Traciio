package handlers

import (
	"context"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	googlecalendartokenservice "github.com/gilabs/crm-healthcare/api/internal/service/google_calendar_token"
	scheduleservice "github.com/gilabs/crm-healthcare/api/internal/service/schedule"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ScheduleHandler struct {
	scheduleService *scheduleservice.Service
}

func NewScheduleHandler(scheduleService *scheduleservice.Service) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// List handles list schedules request
func (h *ScheduleHandler) List(c *gin.Context) {
	var req schedule.ListSchedulesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Apply RBAC scope filtering (replaces legacy auto-filter by logged-in user)
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("schedules")
	} else {
		// Fallback: auto-filter by logged-in user if scope middleware not applied
		userID := ""
		if userIDVal, exists := c.Get("user_id"); exists {
			if id, ok := userIDVal.(string); ok {
				userID = id
			}
		}
		if req.UserID == "" && userID != "" {
			req.UserID = userID
		}
	}

	schedules, pagination, err := h.scheduleService.ListSchedules(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.TaskID != "" {
		meta.Filters["task_id"] = req.TaskID
	}
	if req.UserID != "" {
		meta.Filters["user_id"] = req.UserID
	}
	if req.GoogleCalendarSyncStatus != "" {
		meta.Filters["google_calendar_sync_status"] = req.GoogleCalendarSyncStatus
	}

	response.SuccessResponse(c, schedules, meta)
}

// GetByID handles get schedule by ID request
func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	sch, err := h.scheduleService.GetScheduleByID(id)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, sch, nil)
}

// Create handles create schedule request
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req schedule.CreateScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	createdSchedule, err := h.scheduleService.CreateSchedule(&req, userID)
	if err != nil {
		if err == scheduleservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": req.TaskID,
			}, nil)
			return
		}
		if err == scheduleservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", nil, nil)
			return
		}
		if err == scheduleservice.ErrInvalidTaskAssignment {
			errors.ErrorResponse(c, "INVALID_TASK_ASSIGNMENT", map[string]interface{}{
				"message": "Task must be assigned to a user",
			}, nil)
			return
		}
		errors.ErrorResponse(c, "INTERNAL_ERROR", map[string]interface{}{
			"message": "Failed to create schedule",
			"error":   err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, createdSchedule, nil)
}

// Update handles update schedule request
func (h *ScheduleHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req schedule.UpdateScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedSchedule, err := h.scheduleService.UpdateSchedule(id, &req)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, updatedSchedule, nil)
}

// Delete handles delete schedule request
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.scheduleService.DeleteSchedule(id)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, nil, nil)
}

// SyncToGoogleCalendar handles sync schedule to Google Calendar request
func (h *ScheduleHandler) SyncToGoogleCalendar(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	ctx := context.Background()
	syncResponse, err := h.scheduleService.SyncToGoogleCalendar(ctx, id, userID)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}

		errMsg := err.Error()
		if err == googlecalendartokenservice.ErrTokenNotFound || errMsg == "Google Calendar token not found" || strings.Contains(errMsg, "token not found") {
			errors.ErrorResponse(c, "GOOGLE_CALENDAR_NOT_CONNECTED", map[string]interface{}{
				"message": "Google Calendar is not connected. Please connect your Google Calendar first.",
			}, nil)
			return
		}

		if err == googlecalendartokenservice.ErrTokenExpired || strings.Contains(errMsg, "token expired") || strings.Contains(errMsg, "refresh failed") {
			errors.ErrorResponse(c, "GOOGLE_CALENDAR_TOKEN_EXPIRED", map[string]interface{}{
				"message": "Google Calendar token has expired. Please reconnect your Google Calendar.",
			}, nil)
			return
		}

		if strings.Contains(errMsg, "Cannot specify both default reminders and overrides") {
			errors.ErrorResponse(c, "GOOGLE_CALENDAR_REMINDER_ERROR", map[string]interface{}{
				"message": "Google Calendar reminder configuration error. Please try again or contact support.",
				"error":   errMsg,
			}, nil)
			return
		}

		errors.ErrorResponse(c, "INTERNAL_ERROR", map[string]interface{}{
			"message": "Failed to sync schedule to Google Calendar",
			"error":   errMsg,
		}, nil)
		return
	}

	response.SuccessResponse(c, syncResponse, nil)
}

// UnsyncFromGoogleCalendar handles unsync schedule from Google Calendar request
func (h *ScheduleHandler) UnsyncFromGoogleCalendar(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	ctx := context.Background()
	unsyncResponse, err := h.scheduleService.UnsyncFromGoogleCalendar(ctx, id, userID)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		if err == scheduleservice.ErrGoogleCalendarNotSynced {
			errors.ErrorResponse(c, "GOOGLE_CALENDAR_NOT_SYNCED", map[string]interface{}{
				"message": "Schedule is not synced to Google Calendar",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, unsyncResponse, nil)
}

// ============================================================================
// Mobile-specific handlers — always scoped to the logged-in user
// ============================================================================

// MobileList handles list schedules request scoped to the logged-in user
func (h *ScheduleHandler) MobileList(c *gin.Context) {
	var req schedule.ListSchedulesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Force scope to logged-in user only
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}
	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}
	req.UserID = userID

	schedules, pagination, err := h.scheduleService.ListSchedules(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.TaskID != "" {
		meta.Filters["task_id"] = req.TaskID
	}

	response.SuccessResponse(c, schedules, meta)
}

// MobileGetByID handles get schedule by ID — verifies ownership
func (h *ScheduleHandler) MobileGetByID(c *gin.Context) {
	id := c.Param("id")

	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(string); ok {
			userID = uid
		}
	}
	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	sch, err := h.scheduleService.GetScheduleByID(id)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Ownership check: only return if schedule belongs to this user
	if sch.UserID != userID {
		errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
			"resource":    "schedule",
			"resource_id": id,
		}, nil)
		return
	}

	response.SuccessResponse(c, sch, nil)
}

// MobileCreate handles create schedule — sets user_id to the logged-in user
func (h *ScheduleHandler) MobileCreate(c *gin.Context) {
	var req schedule.CreateScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}
	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	createdSchedule, err := h.scheduleService.CreateSchedule(&req, userID)
	if err != nil {
		if err == scheduleservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": req.TaskID,
			}, nil)
			return
		}
		if err == scheduleservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", nil, nil)
			return
		}
		if err == scheduleservice.ErrInvalidTaskAssignment {
			errors.ErrorResponse(c, "INVALID_TASK_ASSIGNMENT", map[string]interface{}{
				"message": "Task must be assigned to a user",
			}, nil)
			return
		}
		errors.ErrorResponse(c, "INTERNAL_ERROR", map[string]interface{}{
			"message": "Failed to create schedule",
			"error":   err.Error(),
		}, nil)
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, createdSchedule, meta)
}

// MobileUpdate handles update schedule — verifies ownership
func (h *ScheduleHandler) MobileUpdate(c *gin.Context) {
	id := c.Param("id")

	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(string); ok {
			userID = uid
		}
	}
	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	// Verify ownership first
	existing, err := h.scheduleService.GetScheduleByID(id)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if existing.UserID != userID {
		errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
			"resource":    "schedule",
			"resource_id": id,
		}, nil)
		return
	}

	var req schedule.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedSchedule, err := h.scheduleService.UpdateSchedule(id, &req)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, updatedSchedule, nil)
}

// MobileDelete handles delete schedule — verifies ownership
func (h *ScheduleHandler) MobileDelete(c *gin.Context) {
	id := c.Param("id")

	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(string); ok {
			userID = uid
		}
	}
	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	// Verify ownership first
	existing, err := h.scheduleService.GetScheduleByID(id)
	if err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if existing.UserID != userID {
		errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
			"resource":    "schedule",
			"resource_id": id,
		}, nil)
		return
	}

	if err := h.scheduleService.DeleteSchedule(id); err != nil {
		if err == scheduleservice.ErrScheduleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "schedule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponseDeleted(c, "schedule", id, nil)
}
