package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupGroupRoutes(router *gin.RouterGroup, groupHandler *handlers.GroupHandler, jwtManager *jwt.JWTManager) {
	groups := router.Group("/groups")
	groups.Use(middleware.AuthMiddleware(jwtManager))
	groups.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		groups.GET("", groupHandler.List)
		groups.GET("/:id", groupHandler.GetByID)
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/groups*")
		groups.POST("", invalidate, groupHandler.Create)
		groups.PUT("/:id", invalidate, groupHandler.Update)
		groups.DELETE("/:id", invalidate, groupHandler.Delete)
	}
}
