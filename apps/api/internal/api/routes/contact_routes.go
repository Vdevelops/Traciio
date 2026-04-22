package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupContactRoutes(router *gin.RouterGroup, contactHandler *handlers.ContactHandler, jwtManager *jwt.JWTManager) {
	contacts := router.Group("/contacts")
	contacts.Use(middleware.AuthMiddleware(jwtManager))
	{
		contacts.GET("", contactHandler.List)
		contacts.GET("/:id", contactHandler.GetByID)
		contacts.POST("", contactHandler.Create)
		contacts.PUT("/:id", contactHandler.Update)
		contacts.DELETE("/:id", contactHandler.Delete)
	}
}

