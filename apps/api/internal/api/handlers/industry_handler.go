package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	industry_service "github.com/gilabs/crm-healthcare/api/internal/service/industry"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IndustryHandler struct {
	service industry_service.Service
}

func NewIndustryHandler(service industry_service.Service) *IndustryHandler {
	return &IndustryHandler{
		service: service,
	}
}

// List godoc
// @Summary List industries
// @Description Get list of industries with pagination
// @Tags Industry
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by name or code"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort by field" default(order)
// @Param sort_order query string false "Sort order (asc/desc)" default(asc)
// @Success 200 {object} industry.ListIndustriesResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries [get]
func (h *IndustryHandler) List(c *gin.Context) {
	var req industry.ListIndustriesRequest
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

	industries, total, err := h.service.List(&req)
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

	response.SuccessResponse(c, industries, meta)
}

// ListAll godoc
// @Summary List all active industries
// @Description Get all active industries without pagination
// @Tags Industry
// @Accept json
// @Produce json
// @Success 200 {array} industry.IndustryResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries/all [get]
func (h *IndustryHandler) ListAll(c *gin.Context) {
	industries, err := h.service.ListAll()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, industries, nil)
}

// GetByID godoc
// @Summary Get industry by ID
// @Description Get a specific industry by ID
// @Tags Industry
// @Accept json
// @Produce json
// @Param id path string true "Industry ID"
// @Success 200 {object} industry.IndustryResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries/{id} [get]
func (h *IndustryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ind, err := h.service.FindByID(id)
	if err != nil {
		if err == industry_service.ErrIndustryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "industry",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, ind, nil)
}

// Create godoc
// @Summary Create a new industry
// @Description Create a new industry
// @Tags Industry
// @Accept json
// @Produce json
// @Param request body industry.CreateIndustryRequest true "Create industry request"
// @Success 201 {object} industry.IndustryResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries [post]
func (h *IndustryHandler) Create(c *gin.Context) {
	var req industry.CreateIndustryRequest
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

	ind, err := h.service.Create(&req, createdBy)
	if err != nil {
		if err == industry_service.ErrIndustryCodeExists {
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

	response.SuccessResponseCreated(c, ind, meta)
}

// Update godoc
// @Summary Update an industry
// @Description Update an existing industry
// @Tags Industry
// @Accept json
// @Produce json
// @Param id path string true "Industry ID"
// @Param request body industry.UpdateIndustryRequest true "Update industry request"
// @Success 200 {object} industry.IndustryResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries/{id} [put]
func (h *IndustryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req industry.UpdateIndustryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	ind, err := h.service.Update(id, &req)
	if err != nil {
		if err == industry_service.ErrIndustryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "industry",
				"resource_id": id,
			}, nil)
			return
		}
		if err == industry_service.ErrIndustryCodeExists {
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

	response.SuccessResponse(c, ind, meta)
}

// Delete godoc
// @Summary Delete an industry
// @Description Delete an industry (soft delete)
// @Tags Industry
// @Accept json
// @Produce json
// @Param id path string true "Industry ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/industries/{id} [delete]
func (h *IndustryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		if err == industry_service.ErrIndustryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "industry",
				"resource_id": id,
			}, nil)
			return
		}
		if err == industry_service.ErrIndustryInUse {
			errors.ErrorResponse(c, "INDUSTRY_IN_USE", nil, nil)
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

	response.SuccessResponseDeleted(c, "industry", id, meta)
}

