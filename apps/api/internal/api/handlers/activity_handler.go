package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	activityservice "github.com/gilabs/crm-healthcare/api/internal/service/activity"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ActivityHandler struct {
	activityService *activityservice.Service
}

func NewActivityHandler(activityService *activityservice.Service) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

// List handles list activities request
func (h *ActivityHandler) List(c *gin.Context) {
	var req activity.ListActivitiesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	activities, pagination, err := h.activityService.List(&req)
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

	if req.Type != "" {
		meta.Filters["type"] = req.Type
	}
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.ContactID != "" {
		meta.Filters["contact_id"] = req.ContactID
	}
	if req.UserID != "" {
		meta.Filters["user_id"] = req.UserID
	}

	response.SuccessResponse(c, activities, meta)
}

// GetByID handles get activity by ID request
func (h *ActivityHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	activity, err := h.activityService.GetByID(id)
	if err != nil {
		if err == activityservice.ErrActivityNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "activity",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, activity, nil)
}

// Create handles create activity request
func (h *ActivityHandler) Create(c *gin.Context) {
	var req activity.CreateActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	// Set UserID from authenticated user
	req.UserID = userIDStr

	createdActivity, err := h.activityService.Create(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdActivity, meta)
}

// GetTimeline handles get activity timeline request
func (h *ActivityHandler) GetTimeline(c *gin.Context) {
	var req activity.ActivityTimelineRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	activities, err := h.activityService.GetTimeline(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}

	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.ContactID != "" {
		meta.Filters["contact_id"] = req.ContactID
	}
	if req.UserID != "" {
		meta.Filters["user_id"] = req.UserID
	}

	response.SuccessResponse(c, activities, meta)
}

