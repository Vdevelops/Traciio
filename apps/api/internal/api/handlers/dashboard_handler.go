package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/dashboard"
	dashboardservice "github.com/gilabs/crm-healthcare/api/internal/service/dashboard"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DashboardHandler struct {
	dashboardService *dashboardservice.Service
}

func NewDashboardHandler(dashboardService *dashboardservice.Service) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// applyScopeFromContext extracts RBAC ScopedUserIDs from UserContext and applies them to the request.
func applyScopeFromContext(c *gin.Context, req *dashboard.DashboardRequest) {
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("dashboard")
	}
}

// GetOverview handles dashboard overview request
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	// Get user ID from context for target calculation
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		if id, ok := userID.(string); ok {
			userIDStr = id
		}
	}

	overview, err := h.dashboardService.GetOverview(&req, userIDStr)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, overview, nil)
}

// GetVisitStatistics handles visit statistics request
func (h *DashboardHandler) GetVisitStatistics(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	stats, err := h.dashboardService.GetVisitStatistics(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, stats, nil)
}

// GetPipelineSummary handles pipeline summary request
func (h *DashboardHandler) GetPipelineSummary(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	summary, err := h.dashboardService.GetPipelineSummary(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, summary, nil)
}

// GetTopAccounts handles top accounts request
func (h *DashboardHandler) GetTopAccounts(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	topAccounts, err := h.dashboardService.GetTopAccounts(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, topAccounts, nil)
}

// GetTopSalesRep handles top sales rep request
func (h *DashboardHandler) GetTopSalesRep(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	topSalesRep, err := h.dashboardService.GetTopSalesRep(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, topSalesRep, nil)
}

// GetRecentActivities handles recent activities request
func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	activities, err := h.dashboardService.GetRecentActivities(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, activities, nil)
}

// GetActivityTrends handles activity trends request
func (h *DashboardHandler) GetActivityTrends(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	trends, err := h.dashboardService.GetActivityTrends(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, trends, nil)
}

// ============================================================================
// Super Admin Dashboard Handlers
// ============================================================================

// GetSuperAdminUsersByRole handles super admin users by role request
func (h *DashboardHandler) GetSuperAdminUsersByRole(c *gin.Context) {
	data, err := h.dashboardService.GetSuperAdminUsersByRole()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSuperAdminSystemActivity handles super admin system activity request
func (h *DashboardHandler) GetSuperAdminSystemActivity(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSuperAdminSystemActivity(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSuperAdminAIUsage handles super admin AI usage request
func (h *DashboardHandler) GetSuperAdminAIUsage(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSuperAdminAIUsage(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSuperAdminDataGrowth handles super admin data growth request
func (h *DashboardHandler) GetSuperAdminDataGrowth(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSuperAdminDataGrowth(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSuperAdminErrorSummary handles super admin error summary request
func (h *DashboardHandler) GetSuperAdminErrorSummary(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSuperAdminErrorSummary(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// ============================================================================
// Admin Dashboard Handlers
// ============================================================================

// GetAdminTotalLeads handles admin total leads request
func (h *DashboardHandler) GetAdminTotalLeads(c *gin.Context) {
	data, err := h.dashboardService.GetAdminTotalLeads()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAdminPipelineValue handles admin pipeline value request
func (h *DashboardHandler) GetAdminPipelineValue(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAdminPipelineValue(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAdminPendingApprovals handles admin pending approvals request
func (h *DashboardHandler) GetAdminPendingApprovals(c *gin.Context) {
	data, err := h.dashboardService.GetAdminPendingApprovals()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAdminTaskOverdue handles admin task overdue request
func (h *DashboardHandler) GetAdminTaskOverdue(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAdminTaskOverdue(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// ============================================================================
// Sales Manager Dashboard Handlers
// ============================================================================

// GetSalesManagerPipelineFunnel handles sales manager pipeline funnel request
func (h *DashboardHandler) GetSalesManagerPipelineFunnel(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	data, err := h.dashboardService.GetSalesManagerPipelineFunnel(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesManagerTargetVsActual handles sales manager target vs actual request
func (h *DashboardHandler) GetSalesManagerTargetVsActual(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	data, err := h.dashboardService.GetSalesManagerTargetVsActual(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesManagerVisitCompletion handles sales manager visit completion request
func (h *DashboardHandler) GetSalesManagerVisitCompletion(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	data, err := h.dashboardService.GetSalesManagerVisitCompletion(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesManagerDealsAtRisk handles sales manager deals at risk request
func (h *DashboardHandler) GetSalesManagerDealsAtRisk(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	applyScopeFromContext(c, &req)

	data, err := h.dashboardService.GetSalesManagerDealsAtRisk(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesManagerTeamDraftApprovals handles the team draft-approvals request.
// It returns submitted items (visit reports, tasks, schedules, leads, pipeline)
// from all sales reps under the requesting manager.
func (h *DashboardHandler) GetSalesManagerTeamDraftApprovals(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.UnauthorizedResponse(c, "invalid user ID in context")
		return
	}

	data, err := h.dashboardService.GetSalesManagerTeamDraftApprovals(userIDStr)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// ============================================================================
// Sales Dashboard Handlers
// ============================================================================

// GetSalesTodayTasks handles sales today tasks request
func (h *DashboardHandler) GetSalesTodayTasks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	data, err := h.dashboardService.GetSalesTodayTasks(userIDStr)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesAssignedLeads handles sales assigned leads request
func (h *DashboardHandler) GetSalesAssignedLeads(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSalesAssignedLeads(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesUpcomingVisits handles sales upcoming visits request
func (h *DashboardHandler) GetSalesUpcomingVisits(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSalesUpcomingVisits(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetSalesReminders handles sales reminders request
func (h *DashboardHandler) GetSalesReminders(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetSalesReminders(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// ============================================================================
// Analyst Dashboard Handlers
// ============================================================================

// GetAnalystRevenueTrend handles analyst revenue trend request
func (h *DashboardHandler) GetAnalystRevenueTrend(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAnalystRevenueTrend(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAnalystConversionRate handles analyst conversion rate request
func (h *DashboardHandler) GetAnalystConversionRate(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAnalystConversionRate(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAnalystSalesVelocity handles analyst sales velocity request
func (h *DashboardHandler) GetAnalystSalesVelocity(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAnalystSalesVelocity(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// GetAnalystAIInsights handles analyst AI insights request
func (h *DashboardHandler) GetAnalystAIInsights(c *gin.Context) {
	var req dashboard.DashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	data, err := h.dashboardService.GetAnalystAIInsights(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, data, nil)
}

// ============================================================================
// Mobile Dashboard Handlers
// ============================================================================

// GetMobileOverview handles mobile dashboard overview request
func (h *DashboardHandler) GetMobileOverview(c *gin.Context) {
	var req dashboard.MobileDashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		if id, ok := userID.(string); ok {
			userIDStr = id
		}
	}

	if userIDStr == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"message": "User ID not found in context",
		}, nil)
		return
	}

	overview, err := h.dashboardService.GetMobileOverview(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, overview, nil)
}

// GetMobileVisits handles mobile dashboard visits list request
func (h *DashboardHandler) GetMobileVisits(c *gin.Context) {
	var req dashboard.MobileDashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		if id, ok := userID.(string); ok {
			userIDStr = id
		}
	}

	if userIDStr == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"message": "User ID not found in context",
		}, nil)
		return
	}

	visits, err := h.dashboardService.GetMobileVisits(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visits, nil)
}

// GetMobileTasks handles mobile dashboard tasks list request
func (h *DashboardHandler) GetMobileTasks(c *gin.Context) {
	var req dashboard.MobileDashboardRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		if id, ok := userID.(string); ok {
			userIDStr = id
		}
	}

	if userIDStr == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"message": "User ID not found in context",
		}, nil)
		return
	}

	tasks, err := h.dashboardService.GetMobileTasks(userIDStr, &req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, tasks, nil)
}

