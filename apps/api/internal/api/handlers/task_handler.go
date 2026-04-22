package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	taskservice "github.com/gilabs/crm-healthcare/api/internal/service/task"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	taskService *taskservice.Service
}

func NewTaskHandler(taskService *taskservice.Service) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// List handles list tasks request
func (h *TaskHandler) List(c *gin.Context) {
	var req task.ListTasksRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Apply RBAC scope filtering
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("tasks")
	}

	tasks, pagination, err := h.taskService.ListTasks(&req)
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
	if req.Priority != "" {
		meta.Filters["priority"] = req.Priority
	}
	if req.Type != "" {
		meta.Filters["type"] = req.Type
	}
	if req.AssignedTo != "" {
		meta.Filters["assigned_to"] = req.AssignedTo
	}
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.ContactID != "" {
		meta.Filters["contact_id"] = req.ContactID
	}
	if req.DealID != "" {
		meta.Filters["deal_id"] = req.DealID
	}
	if req.DueDateFrom != nil {
		meta.Filters["due_date_from"] = req.DueDateFrom.Format("2006-01-02")
	}
	if req.DueDateTo != nil {
		meta.Filters["due_date_to"] = req.DueDateTo.Format("2006-01-02")
	}

	response.SuccessResponse(c, tasks, meta)
}

// GetByID handles get task by ID request
func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	t, err := h.taskService.GetTaskByID(id)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, t, nil)
}

// Create handles create task request
func (h *TaskHandler) Create(c *gin.Context) {
	// Check permission to create tasks
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !userCtx.HasPermission("tasks.create") {
			errors.ForbiddenResponse(c, "tasks.create", nil)
			return
		}
	}

	var req task.CreateTaskRequest

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

	createdTask, err := h.taskService.CreateTask(&req, userID)
	if err != nil {
		if err == taskservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": req.AssignedTo,
			}, nil)
			return
		}
		if err == taskservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
			}, nil)
			return
		}
		if err == taskservice.ErrContactNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "contact",
				"resource_id": req.ContactID,
			}, nil)
			return
		}
		if err == taskservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": req.DealID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, createdTask, meta)
}

// Update handles update task request
func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req task.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedTask, err := h.taskService.UpdateTask(id, &req)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		if err == taskservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": req.AssignedTo,
			}, nil)
			return
		}
		if err == taskservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
			}, nil)
			return
		}
		if err == taskservice.ErrContactNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "contact",
				"resource_id": req.ContactID,
			}, nil)
			return
		}
		if err == taskservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": req.DealID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedTask, meta)
}

// Delete handles delete task request
func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.taskService.DeleteTask(id)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "task", id, meta)
}

// Assign handles assign task request
func (h *TaskHandler) Assign(c *gin.Context) {
	id := c.Param("id")
	var req task.AssignTaskRequest

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

	assignedTask, err := h.taskService.AssignTask(id, &req, userID)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		if err == taskservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": req.AssignedTo,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.UpdatedBy = userID
	}

	response.SuccessResponse(c, assignedTask, meta)
}

// Complete handles complete task request
func (h *TaskHandler) Complete(c *gin.Context) {
	id := c.Param("id")

	completedTask, err := h.taskService.CompleteTask(id)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		if err == taskservice.ErrTaskAlreadyCompleted {
			errors.ErrorResponse(c, "CONFLICT", map[string]interface{}{
				"message": "Task is already completed",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, completedTask, meta)
}

// ListReminders handles list reminders request
func (h *TaskHandler) ListReminders(c *gin.Context) {
	var req reminder.ListRemindersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	reminders, pagination, err := h.taskService.ListReminders(&req)
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

	if req.TaskID != "" {
		meta.Filters["task_id"] = req.TaskID
	}
	if req.ReminderType != "" {
		meta.Filters["reminder_type"] = req.ReminderType
	}
	if req.IsSent != nil {
		meta.Filters["is_sent"] = *req.IsSent
	}

	response.SuccessResponse(c, reminders, meta)
}

// GetReminderByID handles get reminder by ID request
func (h *TaskHandler) GetReminderByID(c *gin.Context) {
	id := c.Param("id")

	rem, err := h.taskService.GetReminderByID(id)
	if err != nil {
		if err == taskservice.ErrReminderNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "reminder",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, rem, nil)
}

// CreateReminder handles create reminder request
func (h *TaskHandler) CreateReminder(c *gin.Context) {
	var req reminder.CreateReminderRequest

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

	createdReminder, err := h.taskService.CreateReminder(&req, userID)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": req.TaskID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, createdReminder, meta)
}

// UpdateReminder handles update reminder request
func (h *TaskHandler) UpdateReminder(c *gin.Context) {
	id := c.Param("id")
	var req reminder.UpdateReminderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedReminder, err := h.taskService.UpdateReminder(id, &req)
	if err != nil {
		if err == taskservice.ErrReminderNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "reminder",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedReminder, meta)
}

// DeleteReminder handles delete reminder request
func (h *TaskHandler) DeleteReminder(c *gin.Context) {
	id := c.Param("id")

	err := h.taskService.DeleteReminder(id)
	if err != nil {
		if err == taskservice.ErrReminderNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "reminder",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "reminder", id, meta)
}

// MarkInProgress handles mark task as in progress request
func (h *TaskHandler) MarkInProgress(c *gin.Context) {
	id := c.Param("id")

	updatedTask, err := h.taskService.MarkInProgress(id)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": id,
			}, nil)
			return
		}
		if err == taskservice.ErrCannotMarkCompletedInProgress {
			errors.ErrorResponse(c, "CONFLICT", map[string]interface{}{
				"message": "Cannot mark completed task as in progress",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedTask, meta)
}

// GetMyTasks handles get tasks for logged-in user request (mobile endpoint)
// Supports search parameter: ?search=keyword
// Search only applies to tasks assigned to the logged-in user (security enforced)
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	var req task.ListTasksRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"message": "User ID not found in context",
		}, nil)
		return
	}

	tasks, pagination, err := h.taskService.GetMyTasks(userID, &req)
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
	if req.Priority != "" {
		meta.Filters["priority"] = req.Priority
	}
	if req.Type != "" {
		meta.Filters["type"] = req.Type
	}
	if req.DueDateFrom != nil {
		meta.Filters["due_date_from"] = req.DueDateFrom.Format("2006-01-02")
	}
	if req.DueDateTo != nil {
		meta.Filters["due_date_to"] = req.DueDateTo.Format("2006-01-02")
	}

	response.SuccessResponse(c, tasks, meta)
}

// CreateLeadFromTask handles the quick action: create a lead and link it to a task.
// Requires "tasks.create_lead" permission.
func (h *TaskHandler) CreateLeadFromTask(c *gin.Context) {
	// Check RBAC permission
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !userCtx.HasPermission("tasks.create_lead") {
			errors.ForbiddenResponse(c, "tasks.create_lead", nil)
			return
		}
	}

	taskID := c.Param("id")
	var req task.CreateLeadFromTaskRequest

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

	result, err := h.taskService.CreateLeadFromTask(taskID, &req, userID)
	if err != nil {
		if err == taskservice.ErrTaskNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "task",
				"resource_id": taskID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, err.Error())
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, result, meta)
}
