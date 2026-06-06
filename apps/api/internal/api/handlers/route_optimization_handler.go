package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	routeoptimizationservice "github.com/gilabs/crm-healthcare/api/internal/service/route_optimization"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RouteOptimizationHandler struct {
	routeOptimizationService *routeoptimizationservice.Service
}

func NewRouteOptimizationHandler(routeOptimizationService *routeoptimizationservice.Service) *RouteOptimizationHandler {
	return &RouteOptimizationHandler{
		routeOptimizationService: routeOptimizationService,
	}
}

// Optimize handles route optimization request
func (h *RouteOptimizationHandler) Optimize(c *gin.Context) {
	var req route_optimization.OptimizeRouteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	optimizedRoute, err := h.routeOptimizationService.Optimize(&req, userIDStr)
	if err != nil {
		if err == routeoptimizationservice.ErrWaypointsTooFew {
			errors.ErrorResponse(c, "WAYPOINTS_TOO_FEW", map[string]interface{}{
				"message": "Waypoints terlalu sedikit (minimum 2 required)",
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrWaypointsTooMany {
			errors.ErrorResponse(c, "WAYPOINTS_TOO_MANY", map[string]interface{}{
				"message": "Waypoints terlalu banyak (maximum 25 allowed)",
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrInvalidCoordinates {
			errors.ErrorResponse(c, "INVALID_COORDINATES", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrOptimizationFailed {
			errors.ErrorResponse(c, "ROUTE_OPTIMIZATION_FAILED", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Additional: map[string]interface{}{
			"optimization_type": "auto",
			"waypoints_count":   len(req.Waypoints),
		},
	}

	response.SuccessResponseCreated(c, optimizedRoute, meta)
}

// GetByID handles get optimized route by ID request
func (h *RouteOptimizationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	optimizedRoute, err := h.routeOptimizationService.GetByID(id)
	if err != nil {
		if err == routeoptimizationservice.ErrRouteNotFound {
			errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
				"resource":    "route",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, optimizedRoute, nil)
}

// List handles list optimized routes request
func (h *RouteOptimizationHandler) List(c *gin.Context) {
	var req route_optimization.ListRoutesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context for filtering (optional - can filter by other users if admin)
	userID, exists := c.Get("user_id")
	if exists {
		if id, ok := userID.(string); ok {
			// If no user_id in query, use current user's ID
			if req.UserID == "" {
				req.UserID = id
			}
		}
	}

	routes, pagination, err := h.routeOptimizationService.List(&req)
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
	}

	if req.UserID != "" {
		meta.Filters["user_id"] = req.UserID
	}

	if pagination.Page < pagination.TotalPages {
		nextPage := pagination.Page + 1
		meta.Pagination.NextPage = &nextPage
	}

	if pagination.Page > 1 {
		prevPage := pagination.Page - 1
		meta.Pagination.PrevPage = &prevPage
	}

	response.SuccessResponse(c, routes, meta)
}

// CalculateDistance handles distance calculation request
func (h *RouteOptimizationHandler) CalculateDistance(c *gin.Context) {
	var req route_optimization.CalculateDistanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	distance, err := h.routeOptimizationService.CalculateDistance(&req)
	if err != nil {
		if err == routeoptimizationservice.ErrInvalidCoordinates {
			errors.ErrorResponse(c, "INVALID_COORDINATES", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		errors.ErrorResponse(c, "DISTANCE_CALCULATION_FAILED", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, distance, nil)
}

// Delete handles delete optimized route request
func (h *RouteOptimizationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.routeOptimizationService.Delete(id)
	if err != nil {
		if err == routeoptimizationservice.ErrRouteNotFound {
			errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
				"resource":    "route",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "route", id, meta)
}

// MobileOptimize handles mobile route optimization request (sales only)
// This endpoint is specifically for mobile app and only returns routes for the logged-in sales user
func (h *RouteOptimizationHandler) MobileOptimize(c *gin.Context) {
	var req route_optimization.OptimizeRouteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	optimizedRoute, err := h.routeOptimizationService.Optimize(&req, userIDStr)
	if err != nil {
		if err == routeoptimizationservice.ErrWaypointsTooFew {
			errors.ErrorResponse(c, "WAYPOINTS_TOO_FEW", map[string]interface{}{
				"message": "Waypoints terlalu sedikit (minimum 2 required)",
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrWaypointsTooMany {
			errors.ErrorResponse(c, "WAYPOINTS_TOO_MANY", map[string]interface{}{
				"message": "Waypoints terlalu banyak (maximum 25 allowed)",
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrInvalidCoordinates {
			errors.ErrorResponse(c, "INVALID_COORDINATES", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		if err == routeoptimizationservice.ErrOptimizationFailed {
			errors.ErrorResponse(c, "ROUTE_OPTIMIZATION_FAILED", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Additional: map[string]interface{}{
			"optimization_type": "auto",
			"waypoints_count":   len(req.Waypoints),
		},
	}

	response.SuccessResponseCreated(c, optimizedRoute, meta)
}

// MobileGetMyRoutes handles get my routes request (sales only)
// Returns only routes created by the logged-in user
func (h *RouteOptimizationHandler) MobileGetMyRoutes(c *gin.Context) {
	var req route_optimization.ListRoutesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context - mobile always filters by current user
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	// Force filter by current user for mobile
	req.UserID = userIDStr

	routes, pagination, err := h.routeOptimizationService.List(&req)
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
		Filters: map[string]interface{}{
			"user_id": userIDStr,
		},
	}

	if pagination.Page < pagination.TotalPages {
		nextPage := pagination.Page + 1
		meta.Pagination.NextPage = &nextPage
	}

	if pagination.Page > 1 {
		prevPage := pagination.Page - 1
		meta.Pagination.PrevPage = &prevPage
	}

	response.SuccessResponse(c, routes, meta)
}

// MobileGetRouteByID handles get route by ID request (sales only)
// Only returns route if it belongs to the logged-in user
func (h *RouteOptimizationHandler) MobileGetRouteByID(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	optimizedRoute, err := h.routeOptimizationService.GetByID(id)
	if err != nil {
		if err == routeoptimizationservice.ErrRouteNotFound {
			errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
				"resource":    "route",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Verify route belongs to current user
	if optimizedRoute.UserID != userIDStr {
		errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
			"resource":    "route",
			"resource_id": id,
		}, nil)
		return
	}

	response.SuccessResponse(c, optimizedRoute, nil)
}

// MobileCalculateDistance handles mobile distance calculation request
func (h *RouteOptimizationHandler) MobileCalculateDistance(c *gin.Context) {
	var req route_optimization.CalculateDistanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	distance, err := h.routeOptimizationService.CalculateDistance(&req)
	if err != nil {
		if err == routeoptimizationservice.ErrInvalidCoordinates {
			errors.ErrorResponse(c, "INVALID_COORDINATES", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		errors.ErrorResponse(c, "DISTANCE_CALCULATION_FAILED", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, distance, nil)
}

// MobileDelete handles mobile delete route request (sales only)
// Only allows deletion of routes owned by the logged-in user
func (h *RouteOptimizationHandler) MobileDelete(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	// Verify route belongs to current user before deletion
	optimizedRoute, err := h.routeOptimizationService.GetByID(id)
	if err != nil {
		if err == routeoptimizationservice.ErrRouteNotFound {
			errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
				"resource":    "route",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Verify ownership
	if optimizedRoute.UserID != userIDStr {
		errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
			"resource":    "route",
			"resource_id": id,
		}, nil)
		return
	}

	// Delete route
	err = h.routeOptimizationService.Delete(id)
	if err != nil {
		if err == routeoptimizationservice.ErrRouteNotFound {
			errors.ErrorResponse(c, "ROUTE_NOT_FOUND", map[string]interface{}{
				"resource":    "route",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "route", id, meta)
}
