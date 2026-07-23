package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupVisitReportRoutes(
	router *gin.RouterGroup,
	visitReportHandler *handlers.VisitReportHandler,
	activityTypeHandler *handlers.ActivityTypeHandler,
	jwtManager *jwt.JWTManager,
	scopeMiddleware gin.HandlerFunc,
	permissionService *permissionservice.Service,
) {
	visitReports := router.Group("/visit-reports")
	visitReports.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		visitReports.GET("", visitReportHandler.List)
		visitReports.GET("/:id", visitReportHandler.GetByID)
		visitReports.POST("", visitReportHandler.Create)
		visitReports.PUT("/:id", visitReportHandler.Update)
		visitReports.DELETE("/:id", visitReportHandler.Delete)
		visitReports.POST("/:id/check-in", visitReportHandler.CheckIn)
		visitReports.POST("/:id/check-out", visitReportHandler.CheckOut)
		visitReports.PATCH("/:id/submit", visitReportHandler.Submit) // NEW: Submit for approval
		visitReports.POST("/:id/approve", visitReportHandler.Approve)
		visitReports.POST("/:id/reject", visitReportHandler.Reject)
		visitReports.POST("/:id/photos", visitReportHandler.UploadPhoto)

		// Activity Types management
		visitReports.GET("/activity-types", activityTypeHandler.List)
		visitReports.GET("/activity-types/:id", activityTypeHandler.GetByID)

		// Gated activity types mutations (CRUD)
		activityTypeMutations := visitReports.Group("/activity-types")
		activityTypeMutations.Use(middleware.PermissionMiddleware(permissionService, "visit-reports.activity-type"))
		{
			activityTypeMutations.POST("", activityTypeHandler.Create)
			activityTypeMutations.PUT("/:id", activityTypeHandler.Update)
			activityTypeMutations.DELETE("/:id", activityTypeHandler.Delete)
		}
	}
}
