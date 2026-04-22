package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupLeadSourceRoutes sets up lead source routes
func SetupLeadSourceRoutes(router *gin.RouterGroup, leadSourceHandler *handlers.LeadSourceHandler, jwtManager *jwt.JWTManager) {
	leadSources := router.Group("/lead-sources")
	leadSources.Use(middleware.AuthMiddleware(jwtManager))
	leadSources.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		leadSources.GET("", leadSourceHandler.List)        // List with pagination
		leadSources.GET("/all", leadSourceHandler.ListAll) // List all active lead sources
		leadSources.GET("/:id", leadSourceHandler.GetByID) // Get by ID
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/lead-sources*")
		leadSources.POST("", invalidate, leadSourceHandler.Create)       // Create
		leadSources.PUT("/:id", invalidate, leadSourceHandler.Update)    // Update
		leadSources.DELETE("/:id", invalidate, leadSourceHandler.Delete) // Delete
	}
}
