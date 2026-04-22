package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, permissionHandler *handlers.PermissionHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// User-scoped "me" routes (MUST be before /:id to avoid conflicts)
		users.GET("/me/settings-summary", userHandler.GetMySettingsSummary)
		
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.GetByID)
		users.POST("", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.GET("/:id/permissions", permissionHandler.GetUserPermissions)
		// Profile routes
		users.GET("/:id/profile", userHandler.GetProfile)
		users.PUT("/:id/profile", userHandler.UpdateProfile)
		users.PUT("/:id/password", userHandler.ChangePassword)
	}
}

