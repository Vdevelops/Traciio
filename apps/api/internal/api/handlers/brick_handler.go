package handlers

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	brickservice "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	apierrors "github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BrickHandler struct {
	brickService *brickservice.Service
}

func NewBrickHandler(brickService *brickservice.Service) *BrickHandler {
	return &BrickHandler{
		brickService: brickService,
	}
}

// List handles list bricks request
func (h *BrickHandler) List(c *gin.Context) {
	var req brick.ListBricksRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			apierrors.HandleValidationError(c, validationErrors)
			return
		}
		apierrors.InvalidQueryParamResponse(c)
		return
	}

	// Filter by manager if user is sales_manager (not admin)
	userRole, exists := c.Get("user_role")
	if exists {
		if roleStr, ok := userRole.(string); ok {
			if roleStr == "sales_manager" {
				// Sales manager can only see bricks they manage
				userID, exists := c.Get("user_id")
				if exists {
					if userIDStr, ok := userID.(string); ok {
						req.ManagerID = &userIDStr
					}
				}
			}
		}
	}

	bricks, pagination, err := h.brickService.List(&req)
	if err != nil {
		apierrors.InternalServerErrorResponse(c, "")
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
	if req.Province != "" {
		meta.Filters["province"] = req.Province
	}
	if req.Regency != "" {
		meta.Filters["regency"] = req.Regency
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.ManagerID != nil {
		meta.Filters["manager_id"] = *req.ManagerID
	}

	response.SuccessResponse(c, bricks, meta)
}

// GetByID handles get brick by ID request
func (h *BrickHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// Check authorization: only admin or manager of this brick can access
	hasAccess, err := checkBrickAccess(c, h.brickService.GetBrickRepo(), id)
	if err != nil {
		apierrors.InternalServerErrorResponse(c, "")
		return
	}
	if !hasAccess {
		apierrors.ForbiddenResponse(c, "VIEW_BRICKS", []string{})
		return
	}

	brick, err := h.brickService.GetByID(id)
	if err != nil {
		if err == brickservice.ErrBrickNotFound {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": id,
			}, nil)
			return
		}
		apierrors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, brick, nil)
}

// Create handles create brick request
func (h *BrickHandler) Create(c *gin.Context) {
	var req brick.CreateBrickRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			apierrors.HandleValidationError(c, validationErrors)
			return
		}
		apierrors.InvalidRequestBodyResponse(c)
		return
	}

	createdBrick, err := h.brickService.Create(&req)
	if err != nil {
		if errors.Is(err, brickservice.ErrBrickAlreadyExists) {
			apierrors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "brick",
				"field":    "code",
				"value":    req.Code,
			}, nil)
			return
		}
		if errors.Is(err, brickservice.ErrRegencyAlreadyBrick) {
			apierrors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "brick",
				"field":    "regency",
				"value":    req.Regency,
				"reason":   "regency already has a brick",
			}, nil)
			return
		}
		apierrors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdBrick, meta)
}

// Update handles update brick request
func (h *BrickHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req brick.UpdateBrickRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			apierrors.HandleValidationError(c, validationErrors)
			return
		}
		apierrors.InvalidRequestBodyResponse(c)
		return
	}

	updatedBrick, err := h.brickService.Update(id, &req)
	if err != nil {
		if err == brickservice.ErrBrickNotFound {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": id,
			}, nil)
			return
		}
		apierrors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedBrick, meta)
}

// Delete handles delete brick request
func (h *BrickHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.brickService.Delete(id)
	if err != nil {
		if errors.Is(err, brickservice.ErrBrickNotFound) {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": id,
			}, nil)
			return
		}

		// Log the error for debugging
		apierrors.ErrorResponse(c, "INTERNAL_SERVER_ERROR", map[string]interface{}{
			"brick_id": id,
			"error":    err.Error(),
		}, nil)
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Brick deleted successfully",
	}, meta)
}

// GetSales handles get sales users in brick request
func (h *BrickHandler) GetSales(c *gin.Context) {
	brickID := c.Param("id")

	// Check authorization: only admin or manager of this brick can access
	hasAccess, err := checkBrickAccess(c, h.brickService.GetBrickRepo(), brickID)
	if err != nil {
		apierrors.InternalServerErrorResponse(c, "")
		return
	}
	if !hasAccess {
		apierrors.ForbiddenResponse(c, "VIEW_BRICK_SALES", []string{})
		return
	}

	sales, err := h.brickService.GetSalesByBrickID(brickID)
	if err != nil {
		if err == brickservice.ErrBrickNotFound {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": brickID,
			}, nil)
			return
		}
		apierrors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, sales, nil)
}

// AssignSales assigns sales users to a brick
func (h *BrickHandler) AssignSales(c *gin.Context) {
	brickID := c.Param("id")

	hasAccess, err := checkBrickAccess(c, h.brickService.GetBrickRepo(), brickID)
	if err != nil {
		apierrors.InternalServerErrorResponse(c, "")
		return
	}
	if !hasAccess {
		apierrors.ForbiddenResponse(c, "MANAGE_BRICK_SALES", []string{})
		return
	}

	var req brick.AssignSalesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			apierrors.HandleValidationError(c, validationErrors)
			return
		}
		apierrors.InvalidRequestBodyResponse(c)
		return
	}

	sales, err := h.brickService.AssignSales(brickID, req.UserIDs)
	if err != nil {
		if errors.Is(err, brickservice.ErrBrickNotFound) {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": brickID,
			}, nil)
			return
		}
		// User validation errors (not found or not a sales rep)
		apierrors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, sales, nil)
}

// UnassignSales removes sales users from a brick
func (h *BrickHandler) UnassignSales(c *gin.Context) {
	brickID := c.Param("id")

	hasAccess, err := checkBrickAccess(c, h.brickService.GetBrickRepo(), brickID)
	if err != nil {
		apierrors.InternalServerErrorResponse(c, "")
		return
	}
	if !hasAccess {
		apierrors.ForbiddenResponse(c, "MANAGE_BRICK_SALES", []string{})
		return
	}

	var req brick.AssignSalesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			apierrors.HandleValidationError(c, validationErrors)
			return
		}
		apierrors.InvalidRequestBodyResponse(c)
		return
	}

	sales, err := h.brickService.UnassignSales(brickID, req.UserIDs)
	if err != nil {
		if errors.Is(err, brickservice.ErrBrickNotFound) {
			apierrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": brickID,
			}, nil)
			return
		}
		apierrors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, sales, nil)
}

