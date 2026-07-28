package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupMonthlyTargetRoutes(
	router *gin.RouterGroup,
	monthlyTargetHandler *handlers.MonthlyTargetHandler,
	jwtManager *jwt.JWTManager,
	scopeMiddleware gin.HandlerFunc,
	permissionService *permissionservice.Service,
) {
	monthlyTargets := router.Group("/monthly-targets")
	monthlyTargets.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		monthlyTargets.GET("", middleware.PermissionMiddleware(permissionService, "monthly-targets.view"), monthlyTargetHandler.List)
		monthlyTargets.GET("/user-effective", middleware.PermissionMiddleware(permissionService, "monthly-targets.view"), monthlyTargetHandler.GetUserEffectiveTarget)
		monthlyTargets.GET("/:id", middleware.PermissionMiddleware(permissionService, "monthly-targets.view"), monthlyTargetHandler.GetByID)
		
		monthlyTargets.POST("", middleware.PermissionMiddleware(permissionService, "monthly-targets.create"), monthlyTargetHandler.Create)
		monthlyTargets.POST("/bulk", middleware.PermissionMiddleware(permissionService, "monthly-targets.create"), monthlyTargetHandler.BulkCreate)
		monthlyTargets.POST("/bulk-set", middleware.PermissionMiddleware(permissionService, "monthly-targets.edit"), monthlyTargetHandler.BulkSetTarget)
		monthlyTargets.POST("/group-with-users", middleware.PermissionMiddleware(permissionService, "monthly-targets.create"), monthlyTargetHandler.CreateGroupTargetWithUserAssignment)
		monthlyTargets.PUT("/:id", middleware.PermissionMiddleware(permissionService, "monthly-targets.edit"), monthlyTargetHandler.Update)
		monthlyTargets.DELETE("/:id", middleware.PermissionMiddleware(permissionService, "monthly-targets.delete"), monthlyTargetHandler.Delete)
	}
}


