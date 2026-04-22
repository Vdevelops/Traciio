package area_mapping

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/gilabs/crm-healthcare/api/internal/domain/area_mapping"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
)

type Handler struct {
	service area_mapping.Service
}

// NewHandler creates a new area mapping handler
func NewHandler(service area_mapping.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CaptureLocation godoc
// @Summary Capture current location
// @Description Captures the current GPS location with accuracy and metadata
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param body body area_mapping.CaptureLocationRequest true "Location capture data"
// @Success 201 {object} response.APIResponse{data=area_mapping.AreaCapture}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/capture [post]
func (h *Handler) CaptureLocation(c *gin.Context) {
	var req area_mapping.CaptureLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	capture, err := h.service.CaptureLocation(req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to capture location")
		return
	}

	meta := &response.Meta{
		CreatedBy: capture.ID,
	}

	response.SuccessResponseCreated(c, capture, meta)
}

// GetAreaCaptures godoc
// @Summary Get area captures
// @Description Retrieves list of area captures with filters and pagination
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param visit_report_id query string false "Filter by visit report ID"
// @Param capture_type query string false "Filter by type (check_in, check_out, area)"
// @Param captured_after query string false "Start date (YYYY-MM-DD)"
// @Param captured_before query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.APIResponse{data=[]area_mapping.AreaCapture}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/captures [get]
func (h *Handler) GetAreaCaptures(c *gin.Context) {
	filter := buildAreaCaptureFilter(c)

	captures, total, err := h.service.GetAreaCaptures(filter)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to get area captures")
		return
	}

	meta := buildPaginationMeta(filter.Page, filter.PerPage, int(total), c.Request.URL.Query())

	response.SuccessResponse(c, captures, meta)
}

// CreateTerritory godoc
// @Summary Create territory
// @Description Creates a new territory polygon with boundaries
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param body body area_mapping.CreateTerritoryRequest true "Territory data"
// @Success 201 {object} response.APIResponse{data=area_mapping.Territory}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/territories [post]
func (h *Handler) CreateTerritory(c *gin.Context) {
	var req area_mapping.CreateTerritoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	territory, err := h.service.CreateTerritory(req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to create territory")
		return
	}

	meta := &response.Meta{
		CreatedBy: territory.ID,
	}

	response.SuccessResponseCreated(c, territory, meta)
}

// UpdateTerritory godoc
// @Summary Update territory
// @Description Updates an existing territory
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param id path string true "Territory ID"
// @Param body body area_mapping.UpdateTerritoryRequest true "Territory update data"
// @Success 200 {object} response.APIResponse{data=area_mapping.Territory}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/territories/{id} [put]
func (h *Handler) UpdateTerritory(c *gin.Context) {
	id := c.Param("id")

	var req area_mapping.UpdateTerritoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	territory, err := h.service.UpdateTerritory(id, req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to update territory")
		return
	}

	meta := &response.Meta{
		UpdatedBy: territory.ID,
	}

	response.SuccessResponse(c, territory, meta)
}

// GetTerritories godoc
// @Summary Get territories
// @Description Retrieves list of territories with filters and pagination
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param assigned_to query string false "Filter by assigned user ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.APIResponse{data=[]area_mapping.Territory}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/territories [get]
func (h *Handler) GetTerritories(c *gin.Context) {
	filter := buildTerritoryFilter(c)

	territories, total, err := h.service.GetTerritories(filter)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to get territories")
		return
	}

	meta := buildPaginationMeta(filter.Page, filter.PerPage, int(total), c.Request.URL.Query())

	response.SuccessResponse(c, territories, meta)
}

// GetTerritoryByID godoc
// @Summary Get territory by ID
// @Description Retrieves a specific territory by ID
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param id path string true "Territory ID"
// @Success 200 {object} response.APIResponse{data=area_mapping.Territory}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/territories/{id} [get]
func (h *Handler) GetTerritoryByID(c *gin.Context) {
	id := c.Param("id")

	territory, err := h.service.GetTerritoryByID(id)
	if err != nil {
		errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
			"resource":    "territory",
			"resource_id": id,
		}, nil)
		return
	}

	response.SuccessResponse(c, territory, nil)
}

// DeleteTerritory godoc
// @Summary Delete territory
// @Description Deletes a territory
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param id path string true "Territory ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/territories/{id} [delete]
func (h *Handler) DeleteTerritory(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteTerritory(id)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to delete territory")
		return
	}

	meta := &response.Meta{
		DeletedBy: id,
	}

	response.SuccessResponseDeleted(c, "territory", id, meta)
}

// CheckPointInTerritory godoc
// @Summary Check point in territory
// @Description Checks if a geographical point is within a specific territory
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param territory_id query string true "Territory ID"
// @Success 200 {object} response.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/check-territory [get]
func (h *Handler) CheckPointInTerritory(c *gin.Context) {
	var req struct {
		Lat         float64 `form:"lat" binding:"required,min=-90,max=90"`
		Lng         float64 `form:"lng" binding:"required,min=-180,max=180"`
		TerritoryID string  `form:"territory_id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	isWithin, err := h.service.CheckPointInTerritory(req.Lat, req.Lng, req.TerritoryID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to check territory")
		return
	}

	result := map[string]interface{}{
		"is_within":    isWithin,
		"territory_id": req.TerritoryID,
		"latitude":     req.Lat,
		"longitude":    req.Lng,
	}

	response.SuccessResponse(c, result, nil)
}

// GetCoverageAnalysis godoc
// @Summary Get coverage analysis
// @Description Calculates coverage analysis for a territory within a date range
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param territory_id query string true "Territory ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse{data=area_mapping.CoverageAnalysis}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/coverage [get]
func (h *Handler) GetCoverageAnalysis(c *gin.Context) {
	var req struct {
		TerritoryID string `form:"territory_id" binding:"required"`
		StartDate   string `form:"start_date" binding:"required"`
		EndDate     string `form:"end_date" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		errors.InvalidQueryParamResponse(c)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		errors.InvalidQueryParamResponse(c)
		return
	}

	analysis, err := h.service.CalculateCoverage(req.TerritoryID, startDate, endDate)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to calculate coverage")
		return
	}

	response.SuccessResponse(c, analysis, nil)
}

// GetHeatmap godoc
// @Summary Get heatmap data
// @Description Retrieves aggregated location data for heatmap visualization
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param capture_type query string false "Filter by type (check_in, check_out, area)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param territory_id query string false "Filter by territory ID"
// @Success 200 {object} response.APIResponse{data=[]area_mapping.HeatmapPoint}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/heatmap [get]
func (h *Handler) GetHeatmap(c *gin.Context) {
	filter := buildHeatmapFilter(c)

	heatmapData, err := h.service.GetHeatmapData(filter)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to get heatmap data")
		return
	}

	response.SuccessResponse(c, heatmapData, nil)
}

// AssignTerritory godoc
// @Summary Assign territory to user
// @Description Assigns a territory to a specific user
// @Tags Area Mapping
// @Accept json
// @Produce json
// @Param body body map[string]string true "Assignment data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /area-mapping/assign-territory [post]
func (h *Handler) AssignTerritory(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id" binding:"required"`
		TerritoryID string `json:"territory_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	err := h.service.AssignTerritory(req.UserID, req.TerritoryID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "Failed to assign territory")
		return
	}

	result := map[string]interface{}{
		"message":      "Territory assigned successfully",
		"user_id":      req.UserID,
		"territory_id": req.TerritoryID,
	}

	response.SuccessResponse(c, result, nil)
}
