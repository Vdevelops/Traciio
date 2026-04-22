package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(router *gin.RouterGroup, notificationHandler *handlers.NotificationHandler, wsHandler *handlers.WebSocketHandler, jwtManager *jwt.JWTManager) {
	notifications := router.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware(jwtManager))
	{
		notifications.GET("", notificationHandler.List)
		notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
		notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		notifications.DELETE("/:id", notificationHandler.Delete)
	}

	// WebSocket route
	ws := router.Group("/ws")
	ws.GET("/notifications", wsHandler.HandleWebSocket)
}

