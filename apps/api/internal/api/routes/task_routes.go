package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupTaskRoutes sets up task and reminder routes
func SetupTaskRoutes(router *gin.RouterGroup, taskHandler *handlers.TaskHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// Task CRUD
		tasks.GET("", taskHandler.List)
		tasks.GET("/:id", taskHandler.GetByID)
		tasks.POST("", middleware.RateLimitMiddleware("mutation"), taskHandler.Create)
		tasks.PUT("/:id", middleware.RateLimitMiddleware("mutation"), taskHandler.Update)
		tasks.DELETE("/:id", middleware.RateLimitMiddleware("mutation"), taskHandler.Delete)

		// Task actions
		tasks.POST("/:id/assign", middleware.RateLimitMiddleware("mutation"), taskHandler.Assign)
		tasks.POST("/:id/complete", middleware.RateLimitMiddleware("mutation"), taskHandler.Complete)
		tasks.POST("/:id/mark-in-progress", middleware.RateLimitMiddleware("mutation"), taskHandler.MarkInProgress)

		// Quick Actions
		tasks.POST("/:id/create-lead", middleware.RateLimitMiddleware("mutation"), taskHandler.CreateLeadFromTask)

		// Reminder CRUD
		tasks.GET("/reminders", taskHandler.ListReminders)
		tasks.GET("/reminders/:id", taskHandler.GetReminderByID)
		tasks.POST("/reminders", middleware.RateLimitMiddleware("mutation"), taskHandler.CreateReminder)
		tasks.PUT("/reminders/:id", middleware.RateLimitMiddleware("mutation"), taskHandler.UpdateReminder)
		tasks.DELETE("/reminders/:id", middleware.RateLimitMiddleware("mutation"), taskHandler.DeleteReminder)
	}
}
