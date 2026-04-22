package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	brickanalytics "github.com/gilabs/crm-healthcare/api/internal/service/brick_analytics"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
)

type BrickAnalyticsHandler struct {
	brickAnalyticsService *brickanalytics.Service
}

// NewBrickAnalyticsHandler creates a new brick analytics handler
func NewBrickAnalyticsHandler(brickAnalyticsService *brickanalytics.Service) *BrickAnalyticsHandler {
	return &BrickAnalyticsHandler{
		brickAnalyticsService: brickAnalyticsService,
	}
}

// GetBrickPerformance handles GET /api/v1/bricks/:id/performance
// Query params: period_start (YYYY-MM-DD), period_end (YYYY-MM-DD)
func (h *BrickAnalyticsHandler) GetBrickPerformance(c *gin.Context) {
	brickID := c.Param("id")
	if brickID == "" {
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Check authorization: only admin or manager of this brick can access
	hasAccess, err := checkBrickAccess(c, h.brickAnalyticsService.GetBrickRepo(), brickID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !hasAccess {
		errors.ForbiddenResponse(c, "VIEW_BRICK_ANALYTICS", []string{})
		return
	}

	// Parse period from query params.
	// IMPORTANT: Use WIB timezone and treat end date as end-of-day (inclusive).
	// Semantics match Sales Performance:
	// - No params: all time
	// - Only period_start: from start to today
	// - Only period_end: from beginning to end
	loc := response.GetTimezoneWIB()
	now := time.Now().In(loc)
	periodStart := time.Time{}
	periodEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(24*time.Hour - 1*time.Second)

	startStr := c.Query("period_start")
	endStr := c.Query("period_end")

	if startStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", startStr, loc); err == nil {
			periodStart = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
		}
	}

	if endStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", endStr, loc); err == nil {
			periodEnd = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc).Add(24*time.Hour - 1*time.Second)
		}
	}

	metrics, err := h.brickAnalyticsService.GetBrickPerformance(brickID, periodStart, periodEnd)
	if err != nil {
		if err == brickanalytics.ErrBrickNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": brickID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, metrics, nil)
}

// ListBrickPerformance handles GET /api/v1/bricks/performance
// Query params: brick_ids (comma-separated), period_start, period_end
func (h *BrickAnalyticsHandler) ListBrickPerformance(c *gin.Context) {
	// Parse brick IDs from query params
	brickIDsStr := c.Query("brick_ids")
	
	// If brick_ids not provided, filter by manager for sales_manager
	var brickIDs []string
	if brickIDsStr == "" {
		// Filter bricks by manager if user is sales_manager
		filteredIDs, err := filterBricksByManager(c, h.brickAnalyticsService.GetBrickRepo())
		if err != nil {
			errors.InternalServerErrorResponse(c, "")
			return
		}
		if filteredIDs == nil {
			// Admin - can see all, but need to provide brick_ids
			errors.InvalidQueryParamResponse(c)
			return
		}
		if len(filteredIDs) == 0 {
			// No bricks accessible
			response.SuccessResponse(c, []interface{}{}, nil)
			return
		}
		brickIDs = filteredIDs
	} else {
		// Parse comma-separated brick IDs
		for _, id := range splitCommaSeparated(brickIDsStr) {
			if id != "" {
				brickIDs = append(brickIDs, id)
			}
		}
	}

	if len(brickIDs) == 0 {
		errors.InvalidQueryParamResponse(c)
		return
	}

	// For sales_manager, verify access to each brick
	userRole, exists := c.Get("user_role")
	if exists {
		if roleStr, ok := userRole.(string); ok {
			if roleStr == "sales_manager" {
				// Filter to only bricks they manage
				filteredIDs, err := filterBricksByManager(c, h.brickAnalyticsService.GetBrickRepo())
				if err != nil {
					errors.InternalServerErrorResponse(c, "")
					return
				}
				if filteredIDs != nil {
					// Filter brickIDs to only include accessible ones
					accessibleMap := make(map[string]bool)
					for _, id := range filteredIDs {
						accessibleMap[id] = true
					}
					filtered := []string{}
					for _, id := range brickIDs {
						if accessibleMap[id] {
							filtered = append(filtered, id)
						}
					}
					brickIDs = filtered
					if len(brickIDs) == 0 {
						errors.ForbiddenResponse(c, "VIEW_BRICK_ANALYTICS", []string{})
						return
					}
				}
			}
		}
	}

	// Parse period from query params (same semantics as Sales Performance)
	loc := response.GetTimezoneWIB()
	now := time.Now().In(loc)
	periodStart := time.Time{}
	periodEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(24*time.Hour - 1*time.Second)

	startStr := c.Query("period_start")
	endStr := c.Query("period_end")

	if startStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", startStr, loc); err == nil {
			periodStart = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
		}
	}

	if endStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", endStr, loc); err == nil {
			periodEnd = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc).Add(24*time.Hour - 1*time.Second)
		}
	}

	metrics, err := h.brickAnalyticsService.ListBrickPerformance(brickIDs, periodStart, periodEnd)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, metrics, nil)
}

// Helper function to split comma-separated string
func splitCommaSeparated(s string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if char != ' ' {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

