package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	activitytypeservice "github.com/gilabs/crm-healthcare/api/internal/service/activity_type"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ActivityTypeHandler struct {
	activityTypeService *activitytypeservice.Service
}

func NewActivityTypeHandler(activityTypeService *activitytypeservice.Service) *ActivityTypeHandler {
	return &ActivityTypeHandler{
		activityTypeService: activityTypeService,
	}
}

// List handles list activity types request
func (h *ActivityTypeHandler) List(c *gin.Context) {
	var req activity_type.ListActivityTypesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	types, err := h.activityTypeService.List(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, types, nil)
}

// GetByID handles get activity type by ID request
func (h *ActivityTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	activityType, err := h.activityTypeService.GetByID(id)
	if err != nil {
		if err == activitytypeservice.ErrActivityTypeNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "activity_type",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, activityType, nil)
}

// Create handles create activity type request
func (h *ActivityTypeHandler) Create(c *gin.Context) {
	var req activity_type.CreateActivityTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdType, err := h.activityTypeService.Create(&req)
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

	response.SuccessResponseCreated(c, createdType, meta)
}

// Update handles update activity type request
func (h *ActivityTypeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req activity_type.UpdateActivityTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedType, err := h.activityTypeService.Update(id, &req)
	if err != nil {
		if err == activitytypeservice.ErrActivityTypeNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "activity_type",
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

	response.SuccessResponse(c, updatedType, meta)
}

// Delete handles delete activity type request
func (h *ActivityTypeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.activityTypeService.Delete(id)
	if err != nil {
		if err == activitytypeservice.ErrActivityTypeNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "activity_type",
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

	response.SuccessResponseDeleted(c, "activity_type", id, meta)
}

