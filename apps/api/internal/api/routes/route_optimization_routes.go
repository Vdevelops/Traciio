package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
)

// SetupRouteOptimizationRoutes sets up route optimization routes
func SetupRouteOptimizationRoutes(router *gin.RouterGroup, routeOptimizationHandler *handlers.RouteOptimizationHandler, jwtManager *jwt.JWTManager) {
	routes := router.Group("/route-optimization")
	routes.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Optimize route
		routes.POST("/optimize", routeOptimizationHandler.Optimize)

		// Calculate distance
		routes.POST("/calculate-distance", routeOptimizationHandler.CalculateDistance)

		// List routes (history)
		routes.GET("/history", routeOptimizationHandler.List)

		// Get route by ID
		routes.GET("/route/:id", routeOptimizationHandler.GetByID)

		// Delete route
		routes.DELETE("/route/:id", routeOptimizationHandler.Delete)
	}

	// Mobile-specific routes (sales only)
	mobile := router.Group("/mobile")
	mobile.Use(middleware.AuthMiddleware(jwtManager))
	{
		mobileRoutes := mobile.Group("/route-optimization")
		{
			// Optimize route (mobile)
			mobileRoutes.POST("/optimize", routeOptimizationHandler.MobileOptimize)

			// Get my routes (mobile - only returns routes for logged-in user)
			mobileRoutes.GET("/my-routes", routeOptimizationHandler.MobileGetMyRoutes)

			// Get route by ID (mobile - only if belongs to logged-in user)
			mobileRoutes.GET("/route/:id", routeOptimizationHandler.MobileGetRouteByID)

			// Calculate distance (mobile)
			mobileRoutes.POST("/calculate-distance", routeOptimizationHandler.MobileCalculateDistance)

			// Delete route (mobile - only if belongs to logged-in user)
			mobileRoutes.DELETE("/route/:id", routeOptimizationHandler.MobileDelete)
		}
	}
}


