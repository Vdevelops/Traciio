package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupSalesOverviewRoutes sets up sales overview routes
func SetupSalesOverviewRoutes(router *gin.RouterGroup, salesOverviewHandler *handlers.SalesOverviewHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	salesOverview := router.Group("/sales-overview")
	salesOverview.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// Sales Performance
		salesOverview.GET("/performance",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLStatsShort, IncludeQueryParams: true, IncludeUserID: true}),
			salesOverviewHandler.ListSalesPerformance,
		)
		salesOverview.GET("/monthly-overview",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLStatsShort, IncludeQueryParams: true, IncludeUserID: true}),
			salesOverviewHandler.GetMonthlySalesOverview,
		)
		salesOverview.GET("/performance/:userId",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLStatsShort, IncludeQueryParams: true, IncludeUserID: true}),
			salesOverviewHandler.GetSalesPerformanceDetail,
		)

		// Sales Rep Detail
		salesOverview.GET("/sales-rep/:userId",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLStatsShort, IncludeQueryParams: true, IncludeUserID: true}),
			salesOverviewHandler.GetSalesRepDetail,
		)
		salesOverview.GET("/sales-rep/:userId/check-in-locations",
			cache.CacheMiddleware(&cache.CacheMiddlewareConfig{TTL: cache.TTLStatsShort, IncludeQueryParams: true, IncludeUserID: true}),
			salesOverviewHandler.GetSalesRepCheckInLocations,
		)
	}
}
