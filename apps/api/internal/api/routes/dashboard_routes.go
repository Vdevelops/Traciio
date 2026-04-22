package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(router *gin.RouterGroup, dashboardHandler *handlers.DashboardHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(jwtManager))
	dashboard.Use(middleware.RateLimitMiddleware("high_volume"))
	dashboard.Use(scopeMiddleware)
	{
		dashboard.GET("/overview", dashboardHandler.GetOverview)
		dashboard.GET("/visits", dashboardHandler.GetVisitStatistics)
		dashboard.GET("/activity-trends", dashboardHandler.GetActivityTrends)
		dashboard.GET("/pipeline", dashboardHandler.GetPipelineSummary)
		dashboard.GET("/top-accounts", dashboardHandler.GetTopAccounts)
		dashboard.GET("/top-sales-rep", dashboardHandler.GetTopSalesRep)
		dashboard.GET("/recent-activities", dashboardHandler.GetRecentActivities)

		// Super Admin Dashboard Routes
		superAdmin := dashboard.Group("/super-admin")
		{
			superAdmin.GET("/users-by-role", dashboardHandler.GetSuperAdminUsersByRole)
			superAdmin.GET("/system-activity", dashboardHandler.GetSuperAdminSystemActivity)
			superAdmin.GET("/ai-usage", dashboardHandler.GetSuperAdminAIUsage)
			superAdmin.GET("/data-growth", dashboardHandler.GetSuperAdminDataGrowth)
			superAdmin.GET("/error-summary", dashboardHandler.GetSuperAdminErrorSummary)
		}

		// Admin Dashboard Routes
		admin := dashboard.Group("/admin")
		{
			admin.GET("/total-leads", dashboardHandler.GetAdminTotalLeads)
			admin.GET("/pipeline-value", dashboardHandler.GetAdminPipelineValue)
			admin.GET("/pending-approvals", dashboardHandler.GetAdminPendingApprovals)
			admin.GET("/task-overdue", dashboardHandler.GetAdminTaskOverdue)
		}

		// Sales Manager Dashboard Routes
		salesManager := dashboard.Group("/sales-manager")
		{
			salesManager.GET("/pipeline-funnel", dashboardHandler.GetSalesManagerPipelineFunnel)
			salesManager.GET("/target-vs-actual", dashboardHandler.GetSalesManagerTargetVsActual)
			salesManager.GET("/visit-completion", dashboardHandler.GetSalesManagerVisitCompletion)
			salesManager.GET("/deals-at-risk", dashboardHandler.GetSalesManagerDealsAtRisk)
			salesManager.GET("/team-draft-approvals", dashboardHandler.GetSalesManagerTeamDraftApprovals)
		}

		// Sales Dashboard Routes
		sales := dashboard.Group("/sales")
		{
			sales.GET("/today-tasks", dashboardHandler.GetSalesTodayTasks)
			sales.GET("/assigned-leads", dashboardHandler.GetSalesAssignedLeads)
			sales.GET("/upcoming-visits", dashboardHandler.GetSalesUpcomingVisits)
			sales.GET("/reminders", dashboardHandler.GetSalesReminders)
		}

		// Analyst Dashboard Routes
		analyst := dashboard.Group("/analyst")
		{
			analyst.GET("/revenue-trend", dashboardHandler.GetAnalystRevenueTrend)
			analyst.GET("/conversion-rate", dashboardHandler.GetAnalystConversionRate)
			analyst.GET("/sales-velocity", dashboardHandler.GetAnalystSalesVelocity)
			analyst.GET("/ai-insights", dashboardHandler.GetAnalystAIInsights)
		}

		// Mobile Dashboard Routes
		mobile := dashboard.Group("/mobile")
		{
			mobile.GET("/overview", dashboardHandler.GetMobileOverview)
			mobile.GET("/visits", dashboardHandler.GetMobileVisits)
			mobile.GET("/tasks", dashboardHandler.GetMobileTasks)
		}
	}
}

