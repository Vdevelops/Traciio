package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	productanalyticsservice "github.com/gilabs/crm-healthcare/api/internal/service/product_analytics"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type ProductAnalyticsHandler struct {
	productAnalyticsService *productanalyticsservice.Service
}

func NewProductAnalyticsHandler(productAnalyticsService *productanalyticsservice.Service) *ProductAnalyticsHandler {
	return &ProductAnalyticsHandler{
		productAnalyticsService: productAnalyticsService,
	}
}

// GetProductPerformance handles get product performance request
// GET /api/v1/product-analytics/product/:id/performance
func (h *ProductAnalyticsHandler) GetProductPerformance(c *gin.Context) {
	productID := c.Param("id")

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	endDate := time.Now()
	startDate := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	performance, err := h.productAnalyticsService.GetProductPerformance(productID, startDate, endDate, getScopedUserIDs(c))
	if err != nil {
		if err == productanalyticsservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range. Start date must be before end date.",
			}, nil)
			return
		}
		if err == productanalyticsservice.ErrProductNotFound {
			errors.NotFoundResponse(c, "product", productID)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	}

	response.SuccessResponse(c, performance, meta)
}

// GetProductComparison handles product comparison request
// GET /api/v1/product-analytics/product-comparison
func (h *ProductAnalyticsHandler) GetProductComparison(c *gin.Context) {
	// Get product IDs from query (comma-separated)
	productIDsStr := c.Query("product_ids")
	if productIDsStr == "" {
		errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
			"field":   "product_ids",
			"message": "Product IDs are required (comma-separated)",
		}, nil)
		return
	}

	productIDs := parseCommaSeparated(productIDsStr)
	if len(productIDs) == 0 {
		errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
			"field":   "product_ids",
			"message": "At least one product ID is required",
		}, nil)
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	endDate := time.Now()
	startDate := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	comparison, err := h.productAnalyticsService.GetProductComparison(productIDs, startDate, endDate, getScopedUserIDs(c))
	if err != nil {
		if err == productanalyticsservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range. Start date must be before end date.",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"product_ids": productIDs,
			"start_date":  startDate.Format("2006-01-02"),
			"end_date":    endDate.Format("2006-01-02"),
		},
	}

	response.SuccessResponse(c, comparison, meta)
}

// GetProductTrends handles get product trends request
// GET /api/v1/product-analytics/product-trends/:id
func (h *ProductAnalyticsHandler) GetProductTrends(c *gin.Context) {
	productID := c.Param("id")
	groupBy := c.DefaultQuery("group_by", "month")

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0) // Default to last year

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	trends, err := h.productAnalyticsService.GetProductTrends(productID, startDate, endDate, groupBy, getScopedUserIDs(c))
	if err != nil {
		if err == productanalyticsservice.ErrInvalidDateRange {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid date range. Start date must be before end date.",
			}, nil)
			return
		}
		if err == productanalyticsservice.ErrProductNotFound {
			errors.NotFoundResponse(c, "product", productID)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
			"group_by":   groupBy,
		},
	}

	response.SuccessResponse(c, trends, meta)
}

// GetProductsList handles get products list with analytics request
// GET /api/v1/product-analytics/products-list
// GetProductsList handles get products list with analytics request
// GET /api/v1/product-analytics/products-list
func (h *ProductAnalyticsHandler) GetProductsList(c *gin.Context) {
	// Parse query parameters
	sortBy := c.DefaultQuery("sort_by", "total_sold") // total_sold, revenue, profit, name
	orderBy := c.DefaultQuery("order", "desc")        // asc, desc
	period := c.DefaultQuery("period", "month")        // day, week, month, year, all (default: month for leaderboard)
	search := c.Query("search")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	perPageStr := c.DefaultQuery("per_page", "20")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = 20
	}

	// Support 'limit' for backward compatibility and specific use cases (e.g. top products)
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			perPage = parsedLimit
		}
	}
	if perPage > 100 {
		perPage = 100 // Max limit
	}

	// Calculate date range
	var startDate, endDate time.Time
	now := time.Now()

	// If start_date and end_date are provided, use them directly
	if startDateStr != "" && endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// End of the day for end_date
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	} else {
		// Fallback to period-based calculation
		switch period {
		case "day":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		case "week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startDate.AddDate(0, 0, -(weekday - 1))
			endDate = startDate.AddDate(0, 0, 6)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 1, -1)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, now.Location())
		default:
			// "all" - no date filtering
			startDate = time.Time{}
			endDate = time.Time{}
		}
	}

	productsList, total, err := h.productAnalyticsService.GetProductsList(startDate, endDate, search, sortBy, orderBy, page, perPage, getScopedUserIDs(c))
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"sort_by": sortBy,
			"order":   orderBy,
			"search":  search,
		},
		Pagination: response.NewPaginationMeta(page, perPage, int(total)),
	}
	
	// Include period or date range in filters
	if startDateStr != "" && endDateStr != "" {
		meta.Filters["start_date"] = startDateStr
		meta.Filters["end_date"] = endDateStr
	} else {
		meta.Filters["period"] = period
	}

	response.SuccessResponse(c, productsList, meta)
}

// GetUserProductSales handles get user product sales request
// GET /api/v1/product-analytics/user/:userId/products
func (h *ProductAnalyticsHandler) GetUserProductSales(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Enforce scope: if the caller has a restricted scope, the requested userID must be within it.
	scopedUserIDs := getScopedUserIDs(c)
	if len(scopedUserIDs) > 0 {
		allowed := false
		for _, id := range scopedUserIDs {
			if id == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			errors.ForbiddenResponse(c, "view:user-product-sales", nil)
			return
		}
	}

	// Parse query parameters
	sortBy := c.DefaultQuery("sort_by", "total_sold") // total_sold, revenue, profit, name
	orderBy := c.DefaultQuery("order", "desc")        // asc, desc
	period := c.DefaultQuery("period", "all")         // day, week, month, year, all
	startDateStr := c.Query("start_date")             // YYYY-MM-DD format
	endDateStr := c.Query("end_date")                 // YYYY-MM-DD format

	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	perPageStr := c.DefaultQuery("per_page", "20")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100 // Max limit
	}

	// Calculate date range: prioritize start_date/end_date over period
	var startDate, endDate time.Time
	now := time.Now()
	
	// If start_date and end_date are provided, use them directly
	if startDateStr != "" && endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	} else {
		// Fallback to period-based calculation
		switch period {
		case "day":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		case "week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startDate.AddDate(0, 0, -(weekday - 1))
			endDate = startDate.AddDate(0, 0, 6)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 1, -1)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, now.Location())
		default:
			// "all" - no date filtering (use zero time for both)
			startDate = time.Time{}
			endDate = time.Time{}
		}
	}

	productsList, total, err := h.productAnalyticsService.GetUserProductSales(userID, startDate, endDate, sortBy, orderBy, page, perPage)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"user_id": userID,
			"sort_by": sortBy,
			"order":   orderBy,
		},
		Pagination: response.NewPaginationMeta(page, perPage, int(total)),
	}

	// Include period or date range in filters
	if startDateStr != "" && endDateStr != "" {
		meta.Filters["start_date"] = startDateStr
		meta.Filters["end_date"] = endDateStr
	} else {
		meta.Filters["period"] = period
	}

	response.SuccessResponse(c, productsList, meta)
}

// GetMonthlySales handles get monthly sales request
// GET /api/v1/product-analytics/monthly-sales
func (h *ProductAnalyticsHandler) GetMonthlySales(c *gin.Context) {
	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	now := time.Now()
	startDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
	
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}
	
	monthlySales, err := h.productAnalyticsService.GetMonthlySales(startDate, endDate, getScopedUserIDs(c))
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}
	
	meta := &response.Meta{
		Filters: map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	}
	
	response.SuccessResponse(c, monthlySales, meta)
}

// GetProductMonthlySales handles get product monthly sales request
// GET /api/v1/product-analytics/product/:id/monthly-sales
func (h *ProductAnalyticsHandler) GetProductMonthlySales(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		errors.InvalidQueryParamResponse(c)
		return
	}
	
	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	now := time.Now()
	startDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
	
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}
	
	monthlySales, err := h.productAnalyticsService.GetProductMonthlySales(productID, startDate, endDate, getScopedUserIDs(c))
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}
	
	meta := &response.Meta{
		Filters: map[string]interface{}{
			"product_id": productID,
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	}
	
	response.SuccessResponse(c, monthlySales, meta)
}

// Helper functions

func parseCommaSeparated(str string) []string {
	if str == "" {
		return []string{}
	}
	
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// getScopedUserIDs extracts the scoped user IDs from the request context.
// Returns nil for global scope (admin), a specific list for team scope (sales_manager),
// or a single-element list for own scope (sales rep).
func getScopedUserIDs(c *gin.Context) []string {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil {
		return nil
	}
	return userCtx.GetScopedUserIDs("sales-overview")
}
