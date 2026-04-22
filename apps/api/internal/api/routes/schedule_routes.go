package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupScheduleRoutes sets up schedule routes
func SetupScheduleRoutes(router *gin.RouterGroup, scheduleHandler *handlers.ScheduleHandler, googleCalendarAuthHandler *handlers.GoogleCalendarAuthHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	schedules := router.Group("/schedules")
	schedules.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		// Schedule CRUD
		schedules.GET("", scheduleHandler.List)
		schedules.GET("/:id", scheduleHandler.GetByID)
		schedules.POST("", scheduleHandler.Create)
		schedules.PUT("/:id", scheduleHandler.Update)
		schedules.DELETE("/:id", scheduleHandler.Delete)

		// Google Calendar sync
		schedules.POST("/:id/sync-google-calendar", scheduleHandler.SyncToGoogleCalendar)
		schedules.POST("/:id/unsync-google-calendar", scheduleHandler.UnsyncFromGoogleCalendar)
	}

	// Mobile Routes - user-owned schedules only
	mobile := router.Group("/mobile/schedules")
	mobile.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		mobile.GET("", scheduleHandler.MobileList)
		mobile.GET("/:id", scheduleHandler.MobileGetByID)
		mobile.POST("", middleware.RateLimitMiddleware("mutation"), scheduleHandler.MobileCreate)
		mobile.PUT("/:id", middleware.RateLimitMiddleware("mutation"), scheduleHandler.MobileUpdate)
		mobile.DELETE("/:id", middleware.RateLimitMiddleware("mutation"), scheduleHandler.MobileDelete)
	}

	googleCalendar := router.Group("/google-calendar")
	googleCalendar.Use(middleware.AuthMiddleware(jwtManager))
	{
		googleCalendar.GET("/status", googleCalendarAuthHandler.GetConnectionStatus)
		googleCalendar.GET("/auth-url", googleCalendarAuthHandler.GetAuthURL)
		googleCalendar.POST("/exchange-code", googleCalendarAuthHandler.ExchangeCode)
		googleCalendar.DELETE("/disconnect", googleCalendarAuthHandler.Disconnect)
	}

	// OAuth2 callback route is registered at root router level in main.go
	// to ensure it's accessible without authentication middleware
}
