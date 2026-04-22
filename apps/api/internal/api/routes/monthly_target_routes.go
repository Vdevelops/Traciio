package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupMonthlyTargetRoutes(router *gin.RouterGroup, monthlyTargetHandler *handlers.MonthlyTargetHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	monthlyTargets := router.Group("/monthly-targets")
	monthlyTargets.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		monthlyTargets.GET("", monthlyTargetHandler.List)
		monthlyTargets.GET("/user-effective", monthlyTargetHandler.GetUserEffectiveTarget)
		monthlyTargets.GET("/:id", monthlyTargetHandler.GetByID)
		monthlyTargets.POST("", monthlyTargetHandler.Create)
		monthlyTargets.POST("/bulk", monthlyTargetHandler.BulkCreate)
		monthlyTargets.POST("/bulk-set", monthlyTargetHandler.BulkSetTarget)
		monthlyTargets.POST("/group-with-users", monthlyTargetHandler.CreateGroupTargetWithUserAssignment)
		monthlyTargets.PUT("/:id", monthlyTargetHandler.Update)
		monthlyTargets.DELETE("/:id", monthlyTargetHandler.Delete)
	}
}

