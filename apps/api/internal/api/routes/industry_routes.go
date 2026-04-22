package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupIndustryRoutes sets up industry routes
func SetupIndustryRoutes(router *gin.RouterGroup, industryHandler *handlers.IndustryHandler, jwtManager *jwt.JWTManager) {
	industries := router.Group("/industries")
	industries.Use(middleware.AuthMiddleware(jwtManager))
	industries.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		industries.GET("", industryHandler.List)        // List with pagination
		industries.GET("/all", industryHandler.ListAll) // List all active industries
		industries.GET("/:id", industryHandler.GetByID) // Get by ID
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/industries*")
		industries.POST("", invalidate, industryHandler.Create)       // Create
		industries.PUT("/:id", invalidate, industryHandler.Update)    // Update
		industries.DELETE("/:id", invalidate, industryHandler.Delete) // Delete
	}
}
