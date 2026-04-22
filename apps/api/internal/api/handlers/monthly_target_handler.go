package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	monthlytargetservice "github.com/gilabs/crm-healthcare/api/internal/service/monthly_target"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MonthlyTargetHandler struct {
	monthlyTargetService *monthlytargetservice.Service
}

func NewMonthlyTargetHandler(monthlyTargetService *monthlytargetservice.Service) *MonthlyTargetHandler {
	return &MonthlyTargetHandler{
		monthlyTargetService: monthlyTargetService,
	}
}

// List handles list monthly targets request
func (h *MonthlyTargetHandler) List(c *gin.Context) {
	var req monthly_target.ListMonthlyTargetsRequest

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
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("monthly-targets")
	}

	targets, pagination, err := h.monthlyTargetService.List(&req)
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
		Additional: map[string]interface{}{
			"total_target_amount": pagination.TotalAmount,
		},
	}

	if req.GroupID != nil {
		meta.Filters["group_id"] = *req.GroupID
	}
	if req.UserID != nil {
		meta.Filters["user_id"] = *req.UserID
	}
	if req.BrickID != nil {
		meta.Filters["brick_id"] = *req.BrickID
	}
	if req.Year != nil {
		meta.Filters["year"] = *req.Year
	}
	if req.Month != nil {
		meta.Filters["month"] = *req.Month
	}
	if req.Scope != "" {
		meta.Filters["scope"] = req.Scope
	}

	response.SuccessResponse(c, targets, meta)
}

// GetByID handles get monthly target by ID request
func (h *MonthlyTargetHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	target, err := h.monthlyTargetService.GetByID(id)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "monthly_target",
				"monthly_target_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, target, nil)
}

// Create handles create monthly target request
func (h *MonthlyTargetHandler) Create(c *gin.Context) {
	var req monthly_target.CreateMonthlyTargetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdTarget, err := h.monthlyTargetService.Create(&req)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "monthly_target",
				"reason":    "target already exists for this group/user and period",
			}, nil)
			return
		}
		if err == monthlytargetservice.ErrInvalidTargetScope {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"field":   "group_id/user_id/brick_id",
				"message": "either group_id, user_id, or brick_id must be provided, but only one",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdTarget, meta)
}

// BulkCreate handles bulk create monthly targets request
func (h *MonthlyTargetHandler) BulkCreate(c *gin.Context) {
	var req monthly_target.BulkCreateMonthlyTargetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdTargets, err := h.monthlyTargetService.BulkCreate(&req)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "monthly_target",
				"reason":    "one or more targets already exist for this group/user and period",
			}, nil)
			return
		}
		if err == monthlytargetservice.ErrInvalidTargetScope {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"field":   "group_ids/user_ids",
				"message": "either group_ids or user_ids must be provided, not both",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdTargets, meta)
}

// Update handles update monthly target request
func (h *MonthlyTargetHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req monthly_target.UpdateMonthlyTargetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedTarget, err := h.monthlyTargetService.Update(id, &req)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "monthly_target",
				"monthly_target_id": id,
			}, nil)
			return
		}
		if err == monthlytargetservice.ErrMonthlyTargetAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "monthly_target",
				"reason":    "target already exists for this group/user and period",
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

	response.SuccessResponse(c, updatedTarget, meta)
}

// Delete handles delete monthly target request
func (h *MonthlyTargetHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.monthlyTargetService.Delete(id)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "monthly_target",
				"monthly_target_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Monthly target deleted successfully",
	}, meta)
}

// GetUserEffectiveTarget handles get user effective target request (with group fallback)
func (h *MonthlyTargetHandler) GetUserEffectiveTarget(c *gin.Context) {
	var req monthly_target.GetUserTargetRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	target, err := h.monthlyTargetService.GetUserEffectiveTarget(req.UserID, req.Year, req.Month)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "monthly_target",
				"message":  "no target found for user or group",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, target, nil)
}

// CreateGroupTargetWithUserAssignment handles create group target with auto-assign to all users request
func (h *MonthlyTargetHandler) CreateGroupTargetWithUserAssignment(c *gin.Context) {
	var req monthly_target.CreateGroupTargetWithUserAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	groupTarget, userTargets, err := h.monthlyTargetService.CreateGroupTargetWithUserAssignment(&req)
	if err != nil {
		if err == monthlytargetservice.ErrMonthlyTargetAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "monthly_target",
				"reason":    "group target already exists for this group and period",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, map[string]interface{}{
		"group_target": groupTarget,
		"user_targets":    userTargets,
		"total_users":     len(userTargets),
	}, meta)
}

// BulkSetTarget handles bulk set monthly targets request (e.g., set same target for multiple months)
func (h *MonthlyTargetHandler) BulkSetTarget(c *gin.Context) {
	var req monthly_target.BulkSetTargetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdTargets, err := h.monthlyTargetService.BulkSetTarget(&req)
	if err != nil {
		if err == monthlytargetservice.ErrInvalidTargetScope {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"field":   "group_id/user_id/brick_id",
				"message": "either group_id, user_id, or brick_id must be provided, but only one",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponse(c, createdTargets, meta)
}

