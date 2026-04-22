package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupContactRoleRoutes(router *gin.RouterGroup, contactRoleHandler *handlers.ContactRoleHandler, jwtManager *jwt.JWTManager) {
	contactRoles := router.Group("/contact-roles")
	contactRoles.Use(middleware.AuthMiddleware(jwtManager))
	contactRoles.Use(cachepkg.CacheMiddleware(&cachepkg.CacheMiddlewareConfig{
		TTL:                cachepkg.TTLReferenceData,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}))
	{
		contactRoles.GET("", contactRoleHandler.List)
		contactRoles.GET("/:id", contactRoleHandler.GetByID)
		invalidate := cachepkg.InvalidateCacheMiddleware("api:GET:/api/v1/contact-roles*")
		contactRoles.POST("", invalidate, contactRoleHandler.Create)
		contactRoles.PUT("/:id", invalidate, contactRoleHandler.Update)
		contactRoles.DELETE("/:id", invalidate, contactRoleHandler.Delete)
	}
}
