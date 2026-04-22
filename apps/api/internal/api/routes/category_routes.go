package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.RouterGroup, categoryHandler *handlers.CategoryHandler, jwtManager *jwt.JWTManager) {
	categories := router.Group("/categories")
	categories.Use(middleware.AuthMiddleware(jwtManager))
	categories.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		categories.GET("", categoryHandler.List)
		categories.GET("/:id", categoryHandler.GetByID)
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/categories*")
		categories.POST("", invalidate, categoryHandler.Create)
		categories.PUT("/:id", invalidate, categoryHandler.Update)
		categories.DELETE("/:id", invalidate, categoryHandler.Delete)
	}
}
