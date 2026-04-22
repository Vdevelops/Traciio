package area_mapping

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gilabs/crm-healthcare/api/internal/domain/area_mapping"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
)

// buildAreaCaptureFilter builds area capture filter from query params
func buildAreaCaptureFilter(c *gin.Context) area_mapping.AreaCaptureFilter {
	filter := area_mapping.AreaCaptureFilter{
		Page:    getIntQuery(c, "page", 1),
		PerPage: getIntQuery(c, "per_page", 10),
	}

	if visitReportID := c.Query("visit_report_id"); visitReportID != "" {
		filter.VisitReportID = &visitReportID
	}

	if captureType := c.Query("capture_type"); captureType != "" {
		filter.CaptureType = &captureType
	}

	if capturedAfter := c.Query("captured_after"); capturedAfter != "" {
		if date, err := time.Parse("2006-01-02", capturedAfter); err == nil {
			filter.CapturedAfter = &date
		}
	}

	if capturedBefore := c.Query("captured_before"); capturedBefore != "" {
		if date, err := time.Parse("2006-01-02", capturedBefore); err == nil {
			filter.CapturedBefore = &date
		}
	}

	return filter
}

// buildTerritoryFilter builds territory filter from query params
func buildTerritoryFilter(c *gin.Context) area_mapping.TerritoryFilter {
	filter := area_mapping.TerritoryFilter{
		Page:    getIntQuery(c, "page", 1),
		PerPage: getIntQuery(c, "per_page", 10),
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter.AssignedTo = &assignedTo
	}

	return filter
}

// buildHeatmapFilter builds heatmap filter from query params
func buildHeatmapFilter(c *gin.Context) area_mapping.HeatmapFilter {
	filter := area_mapping.HeatmapFilter{}

	if captureType := c.Query("capture_type"); captureType != "" {
		filter.CaptureType = &captureType
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &date
		}
	}

	if territoryID := c.Query("territory_id"); territoryID != "" {
		filter.TerritoryID = &territoryID
	}

	return filter
}

// buildPaginationMeta builds pagination metadata
func buildPaginationMeta(page, perPage, total int, queryParams url.Values) *response.Meta {
	totalPages := (total + perPage - 1) / perPage
	hasNext := page < totalPages
	hasPrev := page > 1

	// Build filters map
	filters := make(map[string]interface{})
	for key, values := range queryParams {
		if key != "page" && key != "per_page" && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	return &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Filters: filters,
	}
}

// getIntQuery gets integer query parameter with default value
func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	value, _ := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(defaultValue)))
	if value < 1 {
		value = defaultValue
	}
	return value
}
