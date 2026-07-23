package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupRoleRoutes(router *gin.RouterGroup, roleHandler *handlers.RoleHandler, jwtManager *jwt.JWTManager) {
	roles := router.Group("/roles")
	roles.Use(middleware.AuthMiddleware(jwtManager))
	roles.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		roles.GET("", roleHandler.List)
		roles.GET("/:id", roleHandler.GetByID)
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/roles*")
		roles.POST("", invalidate, roleHandler.Create)
		roles.PUT("/:id", invalidate, roleHandler.Update)
		roles.DELETE("/:id", invalidate, roleHandler.Delete)
		roles.GET("/:id/permissions", roleHandler.GetRolePermissions)
		roles.PUT("/:id/permissions", invalidate, roleHandler.AssignPermissions)
	}
}
