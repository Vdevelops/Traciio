package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(
	router *gin.RouterGroup,
	reportHandler *handlers.ReportHandler,
	jwtManager *jwt.JWTManager,
	permissionService *permissionservice.Service,
) {
	reports := router.Group("/reports")
	reports.Use(middleware.AuthMiddleware(jwtManager))
	{
		reports.GET("/visit-reports",
			middleware.PermissionMiddleware(permissionService, "reports.view"),
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetVisitReportReport,
		)
		reports.GET("/pipeline",
			middleware.PermissionMiddleware(permissionService, "reports.view"),
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetPipelineReport,
		)
		reports.GET("/sales-performance",
			middleware.PermissionMiddleware(permissionService, "reports.view"),
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetSalesPerformanceReport,
		)
		reports.GET("/account-activity",
			middleware.PermissionMiddleware(permissionService, "reports.view"),
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetAccountActivityReport,
		)

		// Export endpoints
		reports.GET("/visit-reports/export",
			middleware.PermissionMiddleware(permissionService, "reports.generate"),
			reportHandler.ExportVisitReportReport,
		)
		reports.GET("/pipeline/export",
			middleware.PermissionMiddleware(permissionService, "reports.generate"),
			reportHandler.ExportPipelineReport,
		)
		reports.GET("/sales-performance/export",
			middleware.PermissionMiddleware(permissionService, "reports.generate"),
			reportHandler.ExportSalesPerformanceReport,
		)
		reports.GET("/account-activity/export",
			middleware.PermissionMiddleware(permissionService, "reports.generate"),
			reportHandler.ExportAccountActivityReport,
		)
	}
}
