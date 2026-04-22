package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupLeadStatusRoutes sets up lead status routes
func SetupLeadStatusRoutes(router *gin.RouterGroup, leadStatusHandler *handlers.LeadStatusHandler, jwtManager *jwt.JWTManager) {
	leadStatuses := router.Group("/lead-statuses")
	leadStatuses.Use(middleware.AuthMiddleware(jwtManager))
	leadStatuses.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		leadStatuses.GET("", leadStatusHandler.List)        // List with pagination
		leadStatuses.GET("/all", leadStatusHandler.ListAll) // List all active statuses
		leadStatuses.GET("/:id", leadStatusHandler.GetByID) // Get by ID
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/lead-statuses*")
		leadStatuses.POST("", invalidate, leadStatusHandler.Create)                      // Create
		leadStatuses.PUT("/:id", invalidate, leadStatusHandler.Update)                   // Update
		leadStatuses.DELETE("/:id", invalidate, leadStatusHandler.Delete)                // Delete
		leadStatuses.PATCH("/:id/set-default", invalidate, leadStatusHandler.SetDefault) // Set as default
	}
}
