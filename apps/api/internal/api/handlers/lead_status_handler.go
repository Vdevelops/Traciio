package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	lead_status_service "github.com/gilabs/crm-healthcare/api/internal/service/lead_status"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LeadStatusHandler struct {
	service *lead_status_service.Service
}

func NewLeadStatusHandler(service *lead_status_service.Service) *LeadStatusHandler {
	return &LeadStatusHandler{
		service: service,
	}
}

// List godoc
// @Summary List lead statuses
// @Description Get list of lead statuses with pagination
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by name or code"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort by field" default(order)
// @Param sort_order query string false "Sort order (asc/desc)" default(asc)
// @Success 200 {object} lead_status.ListLeadStatusesResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses [get]
func (h *LeadStatusHandler) List(c *gin.Context) {
	var req lead_status.ListLeadStatusesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PerPage == 0 {
		req.PerPage = 10
	}
	if req.SortBy == "" {
		req.SortBy = "order"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	statuses, total, err := h.service.List(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      int(total),
			TotalPages: int((total + int64(req.PerPage) - 1) / int64(req.PerPage)),
		},
	}

	response.SuccessResponse(c, statuses, meta)
}

// ListAll godoc
// @Summary List all active lead statuses
// @Description Get all active lead statuses without pagination
// @Tags Lead Status
// @Accept json
// @Produce json
// @Success 200 {array} lead_status.LeadStatusResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses/all [get]
func (h *LeadStatusHandler) ListAll(c *gin.Context) {
	statuses, err := h.service.ListAll()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, statuses, nil)
}

// GetByID godoc
// @Summary Get lead status by ID
// @Description Get a specific lead status by ID
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param id path string true "Lead Status ID"
// @Success 200 {object} lead_status.LeadStatusResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses/{id} [get]
func (h *LeadStatusHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	status, err := h.service.FindByID(id)
	if err != nil {
		if err == lead_status_service.ErrLeadStatusNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_status",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, status, nil)
}

// Create godoc
// @Summary Create a new lead status
// @Description Create a new lead status
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param request body lead_status.CreateLeadStatusRequest true "Create lead status request"
// @Success 201 {object} lead_status.LeadStatusResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses [post]
func (h *LeadStatusHandler) Create(c *gin.Context) {
	var req lead_status.CreateLeadStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	var createdBy string
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			createdBy = id
		}
	}

	status, err := h.service.Create(&req, createdBy)
	if err != nil {
		if err == lead_status_service.ErrLeadStatusCodeExists {
			errors.ErrorResponse(c, "DUPLICATE_CODE", map[string]interface{}{
				"code": req.Code,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		CreatedBy: createdBy,
	}

	response.SuccessResponseCreated(c, status, meta)
}

// Update godoc
// @Summary Update a lead status
// @Description Update an existing lead status
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param id path string true "Lead Status ID"
// @Param request body lead_status.UpdateLeadStatusRequest true "Update lead status request"
// @Success 200 {object} lead_status.LeadStatusResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses/{id] [put]
func (h *LeadStatusHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req lead_status.UpdateLeadStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	status, err := h.service.Update(id, &req)
	if err != nil {
		if err == lead_status_service.ErrLeadStatusNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_status",
				"resource_id": id,
			}, nil)
			return
		}
		if err == lead_status_service.ErrLeadStatusCodeExists {
			errors.ErrorResponse(c, "DUPLICATE_CODE", map[string]interface{}{
				"code": req.Code,
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

	response.SuccessResponse(c, status, meta)
}

// Delete godoc
// @Summary Delete a lead status
// @Description Delete a lead status (soft delete)
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param id path string true "Lead Status ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses/{id} [delete]
func (h *LeadStatusHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		if err == lead_status_service.ErrLeadStatusNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_status",
				"resource_id": id,
			}, nil)
			return
		}
		if err == lead_status_service.ErrCannotDeleteDefault {
			errors.ErrorResponse(c, "CANNOT_DELETE_DEFAULT", nil, nil)
			return
		}
		if err == lead_status_service.ErrCannotDeleteConverted {
			errors.ErrorResponse(c, "CANNOT_DELETE_CONVERTED", nil, nil)
			return
		}
		if err == lead_status_service.ErrLeadStatusInUse {
			errors.ErrorResponse(c, "STATUS_IN_USE", nil, nil)
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

	response.SuccessResponseDeleted(c, "lead_status", id, meta)
}

// SetDefault godoc
// @Summary Set default lead status
// @Description Set a lead status as the default status for new leads
// @Tags Lead Status
// @Accept json
// @Produce json
// @Param id path string true "Lead Status ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-statuses/{id}/set-default [patch]
func (h *LeadStatusHandler) SetDefault(c *gin.Context) {
	id := c.Param("id")

	err := h.service.SetDefault(id)
	if err != nil {
		if err == lead_status_service.ErrLeadStatusNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_status",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "Default lead status updated successfully",
	}, nil)
}
