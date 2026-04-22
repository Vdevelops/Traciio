package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(router *gin.RouterGroup, reportHandler *handlers.ReportHandler, jwtManager *jwt.JWTManager) {
	reports := router.Group("/reports")
	reports.Use(middleware.AuthMiddleware(jwtManager))
	{
		reports.GET("/visit-reports",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetVisitReportReport,
		)
		reports.GET("/pipeline",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetPipelineReport,
		)
		reports.GET("/sales-performance",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetSalesPerformanceReport,
		)
		reports.GET("/account-activity",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLReportShort, IncludeQueryParams: true, IncludeUserID: true}),
			reportHandler.GetAccountActivityReport,
		)

		// Export endpoints
		reports.GET("/visit-reports/export", reportHandler.ExportVisitReportReport)
		reports.GET("/pipeline/export", reportHandler.ExportPipelineReport)
		reports.GET("/sales-performance/export", reportHandler.ExportSalesPerformanceReport)
		reports.GET("/account-activity/export", reportHandler.ExportAccountActivityReport)
	}
}
