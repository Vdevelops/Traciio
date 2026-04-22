package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	lead_source_service "github.com/gilabs/crm-healthcare/api/internal/service/lead_source"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LeadSourceHandler struct {
	service lead_source_service.Service
}

func NewLeadSourceHandler(service lead_source_service.Service) *LeadSourceHandler {
	return &LeadSourceHandler{
		service: service,
	}
}

// List godoc
// @Summary List lead sources
// @Description Get list of lead sources with pagination
// @Tags Lead Source
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by name or code"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort by field" default(order)
// @Param sort_order query string false "Sort order (asc/desc)" default(asc)
// @Success 200 {object} lead_source.ListLeadSourcesResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources [get]
func (h *LeadSourceHandler) List(c *gin.Context) {
	var req lead_source.ListLeadSourcesRequest
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

	leadSources, total, err := h.service.List(&req)
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

	response.SuccessResponse(c, leadSources, meta)
}

// ListAll godoc
// @Summary List all active lead sources
// @Description Get all active lead sources without pagination
// @Tags Lead Source
// @Accept json
// @Produce json
// @Success 200 {array} lead_source.LeadSourceResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources/all [get]
func (h *LeadSourceHandler) ListAll(c *gin.Context) {
	leadSources, err := h.service.ListAll()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, leadSources, nil)
}

// GetByID godoc
// @Summary Get lead source by ID
// @Description Get a specific lead source by ID
// @Tags Lead Source
// @Accept json
// @Produce json
// @Param id path string true "Lead Source ID"
// @Success 200 {object} lead_source.LeadSourceResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources/{id} [get]
func (h *LeadSourceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ls, err := h.service.FindByID(id)
	if err != nil {
		if err == lead_source_service.ErrLeadSourceNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_source",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, ls, nil)
}

// Create godoc
// @Summary Create a new lead source
// @Description Create a new lead source
// @Tags Lead Source
// @Accept json
// @Produce json
// @Param request body lead_source.CreateLeadSourceRequest true "Create lead source request"
// @Success 201 {object} lead_source.LeadSourceResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources [post]
func (h *LeadSourceHandler) Create(c *gin.Context) {
	var req lead_source.CreateLeadSourceRequest
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

	ls, err := h.service.Create(&req, createdBy)
	if err != nil {
		if err == lead_source_service.ErrLeadSourceCodeExists {
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

	response.SuccessResponseCreated(c, ls, meta)
}

// Update godoc
// @Summary Update a lead source
// @Description Update an existing lead source
// @Tags Lead Source
// @Accept json
// @Produce json
// @Param id path string true "Lead Source ID"
// @Param request body lead_source.UpdateLeadSourceRequest true "Update lead source request"
// @Success 200 {object} lead_source.LeadSourceResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources/{id} [put]
func (h *LeadSourceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req lead_source.UpdateLeadSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	ls, err := h.service.Update(id, &req)
	if err != nil {
		if err == lead_source_service.ErrLeadSourceNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_source",
				"resource_id": id,
			}, nil)
			return
		}
		if err == lead_source_service.ErrLeadSourceCodeExists {
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

	response.SuccessResponse(c, ls, meta)
}

// Delete godoc
// @Summary Delete a lead source
// @Description Delete a lead source (soft delete)
// @Tags Lead Source
// @Accept json
// @Produce json
// @Param id path string true "Lead Source ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/lead-sources/{id} [delete]
func (h *LeadSourceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		if err == lead_source_service.ErrLeadSourceNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead_source",
				"resource_id": id,
			}, nil)
			return
		}
		if err == lead_source_service.ErrLeadSourceInUse {
			errors.ErrorResponse(c, "LEAD_SOURCE_IN_USE", nil, nil)
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

	response.SuccessResponseDeleted(c, "lead_source", id, meta)
}

