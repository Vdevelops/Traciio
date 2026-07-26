package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupScheduleRoutes sets up schedule routes
func SetupScheduleRoutes(router *gin.RouterGroup, scheduleHandler *handlers.ScheduleHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	schedules := router.Group("/schedules")
	schedules.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// Schedule CRUD
		schedules.GET("", scheduleHandler.List)
		schedules.GET("/:id", scheduleHandler.GetByID)
		schedules.POST("", scheduleHandler.Create)
		schedules.PUT("/:id", scheduleHandler.Update)
		schedules.DELETE("/:id", scheduleHandler.Delete)
	}
}
