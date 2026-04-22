package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupProductAnalyticsRoutes sets up product analytics routes
func SetupProductAnalyticsRoutes(router *gin.RouterGroup, productAnalyticsHandler *handlers.ProductAnalyticsHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	productAnalytics := router.Group("/product-analytics")
	productAnalytics.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// Products list with analytics
		productAnalytics.GET("/products-list", productAnalyticsHandler.GetProductsList)
		
		// Monthly sales analytics (all products)
		productAnalytics.GET("/monthly-sales", productAnalyticsHandler.GetMonthlySales)
		
		// Product performance
		productAnalytics.GET("/product/:id/performance", productAnalyticsHandler.GetProductPerformance)
		
		// Product monthly sales (per product)
		productAnalytics.GET("/product/:id/monthly-sales", productAnalyticsHandler.GetProductMonthlySales)
		
		// Product comparison
		productAnalytics.GET("/product-comparison", productAnalyticsHandler.GetProductComparison)
		
		// Product trends
		productAnalytics.GET("/product-trends/:id", productAnalyticsHandler.GetProductTrends)
		
		// User product sales (products sold by a specific user)
		productAnalytics.GET("/user/:userId/products", productAnalyticsHandler.GetUserProductSales)
	}
}
