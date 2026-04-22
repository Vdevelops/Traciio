package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupMenuRoutes(router *gin.RouterGroup, menuHandler *handlers.MenuHandler, jwtManager *jwt.JWTManager) {
	menus := router.Group("/menus")
	menus.Use(middleware.AuthMiddleware(jwtManager))
	{
		menus.GET("", menuHandler.List)
		menus.GET("/:id", menuHandler.GetByID)
	}
}
