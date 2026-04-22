package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupCustomerPurchaseRoutes sets up customer purchase history routes
func SetupCustomerPurchaseRoutes(router *gin.RouterGroup, handler *handlers.CustomerPurchaseHandler, jwtManager *jwt.JWTManager) {
	purchases := router.Group("/accounts/:id/purchases")
	purchases.Use(middleware.AuthMiddleware(jwtManager))
	{
		purchases.GET("", handler.GetByAccount)
		purchases.GET("/analytics", handler.GetProductAnalytics)
		purchases.GET("/summary", handler.GetSummary)
	}

	// Mobile Routes
	mobile := router.Group("/mobile/accounts/:id/purchases")
	mobile.Use(middleware.AuthMiddleware(jwtManager))
	{
		mobile.GET("", handler.GetByAccount)
		mobile.GET("/analytics", handler.GetProductAnalytics)
		mobile.GET("/summary", handler.GetSummary)
	}
}
