package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
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
}
