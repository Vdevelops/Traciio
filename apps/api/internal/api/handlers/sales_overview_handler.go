package handlers

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	salesoverviewservice "github.com/gilabs/crm-healthcare/api/internal/service/sales_overview"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SalesOverviewHandler struct {
	salesOverviewService *salesoverviewservice.Service
	userRepo             interfaces.UserRepository
	brickRepo            interfaces.BrickRepository
}

func NewSalesOverviewHandler(
	salesOverviewService *salesoverviewservice.Service,
	userRepo interfaces.UserRepository,
	brickRepo interfaces.BrickRepository,
) *SalesOverviewHandler {
	return &SalesOverviewHandler{
		salesOverviewService: salesOverviewService,
		userRepo:             userRepo,
		brickRepo:            brickRepo,
	}
}

// GetSalesPerformanceDetail handles get sales performance detail request
func (h *SalesOverviewHandler) GetSalesPerformanceDetail(c *gin.Context) {
	userID := c.Param("userId")
	var req sales_overview.GetSalesPerformanceDetailRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Validate that the requested user is within the caller's scope
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		scoped := userCtx.GetScopedUserIDs("sales-overview")
		if scoped != nil {
			allowed := false
			for _, id := range scoped {
				if id == userID {
					allowed = true
					break
				}
			}
			if !allowed {
				errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
					"message": "You do not have permission to view this user's data",
				}, nil)
				return
			}
		}
	}

	detail, err := h.salesOverviewService.GetSalesPerformanceDetail(userID, &req)
	if err != nil {
		if err == salesoverviewservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range format",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}

	if req.Period != "" {
		meta.Filters["period"] = req.Period
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, detail, meta)
}

// GetMonthlySalesOverview handles get monthly sales overview request for chart
func (h *SalesOverviewHandler) GetMonthlySalesOverview(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	var startDate, endDate interface{}
	
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}
	
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// Set to end of day
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Get scoped user IDs for filtering
	var scopedUserIDs []string
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		scopedUserIDs = userCtx.GetScopedUserIDs("sales-overview")
	}

	result, err := h.salesOverviewService.GetMonthlySalesOverview(startDate, endDate, scopedUserIDs)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}
	if startDateStr != "" {
		meta.Filters["start_date"] = startDateStr
	}
	if endDateStr != "" {
		meta.Filters["end_date"] = endDateStr
	}

	response.SuccessResponse(c, result, meta)
}

// GetSalesRepDetail handles get sales rep detail request
func (h *SalesOverviewHandler) GetSalesRepDetail(c *gin.Context) {
	userID := c.Param("userId")
	var req sales_overview.GetSalesRepDetailRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Validate that the requested user is within the caller's scope
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		scoped := userCtx.GetScopedUserIDs("sales-overview")
		if scoped != nil {
			allowed := false
			for _, id := range scoped {
				if id == userID {
					allowed = true
					break
				}
			}
			if !allowed {
				errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
					"message": "You do not have permission to view this user's data",
				}, nil)
				return
			}
		}
	}

	detail, err := h.salesOverviewService.GetSalesRepDetail(userID, &req)
	if err != nil {
		if err == salesoverviewservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range format",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Load user info if not already loaded
	if detail.User == nil {
		userObj, err := h.userRepo.FindByID(userID)
		if err == nil {
			detail.User = userObj.ToUserResponse()
		}
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}

	if req.Period != "" {
		meta.Filters["period"] = req.Period
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, detail, meta)
}

// ListSalesPerformance handles list sales performance request
func (h *SalesOverviewHandler) ListSalesPerformance(c *gin.Context) {
	var req sales_overview.ListSalesPerformanceRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	if req.Order == "" && req.SortOrder != "" {
		req.Order = req.SortOrder
	}

	// Apply RBAC scope filtering via UserContext
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("deals")
	}

	// Debug logging
	log.Printf("[ListSalesPerformance] Request received: page=%d, per_page=%d, search=%s, start_date=%s, end_date=%s, brick_id=%s",
		req.Page, req.PerPage, req.Search, req.StartDate, req.EndDate, req.BrickID)

	results, total, err := h.salesOverviewService.ListSalesPerformance(&req)
	if err != nil {
		log.Printf("[ListSalesPerformance] Error: %v", err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	log.Printf("[ListSalesPerformance] Success: total=%d, results_count=%d", total, len(results))

	// Calculate pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      int(total),
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}
	if req.SortBy != "" {
		meta.Filters["sort_by"] = req.SortBy
	}
	if req.Order != "" {
		meta.Filters["order"] = req.Order
	}

	response.SuccessResponse(c, results, meta)
}

// GetSalesRepCheckInLocations handles get sales rep check-in locations request
func (h *SalesOverviewHandler) GetSalesRepCheckInLocations(c *gin.Context) {
	userID := c.Param("userId")
	var req sales_overview.GetSalesRepCheckInLocationsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Validate that the requested user is within the caller's scope
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		scoped := userCtx.GetScopedUserIDs("sales-overview")
		if scoped != nil {
			allowed := false
			for _, id := range scoped {
				if id == userID {
					allowed = true
					break
				}
			}
			if !allowed {
				errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
					"message": "You do not have permission to view this user's data",
				}, nil)
				return
			}
		}
	}

	locations, err := h.salesOverviewService.GetSalesRepCheckInLocations(userID, &req)
	if err != nil {
		if err == salesoverviewservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range format",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Load user info
	if locations.SalesRep == nil {
		userObj, err := h.userRepo.FindByID(userID)
		if err == nil {
			locations.SalesRep = userObj.ToUserResponse()
		}
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      int(locations.TotalVisits),
			TotalPages: 0, // Calculate if needed
		},
		Filters: map[string]interface{}{},
	}

	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, locations, meta)
}

