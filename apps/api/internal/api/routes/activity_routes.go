package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupActivityRoutes(router *gin.RouterGroup, activityHandler *handlers.ActivityHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	activities := router.Group("/activities")
	activities.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		activities.GET("", activityHandler.List)
		activities.GET("/:id", activityHandler.GetByID)
		activities.POST("", middleware.RateLimitMiddleware("mutation"), activityHandler.Create)
		activities.PUT("/:id", middleware.RateLimitMiddleware("mutation"), activityHandler.Update)
		activities.GET("/timeline", activityHandler.GetTimeline)
	}
}
