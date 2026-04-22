package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupPermissionRoutes(router *gin.RouterGroup, permissionHandler *handlers.PermissionHandler, jwtManager *jwt.JWTManager) {
	permissions := router.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(jwtManager))
	{
		permissions.GET("", permissionHandler.List)
		permissions.GET("/:id", permissionHandler.GetByID)
	}
}

