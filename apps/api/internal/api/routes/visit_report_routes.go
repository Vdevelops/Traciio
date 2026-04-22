package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupVisitReportRoutes(router *gin.RouterGroup, visitReportHandler *handlers.VisitReportHandler, activityTypeHandler *handlers.ActivityTypeHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
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
		visitReports.POST("/activity-types", activityTypeHandler.Create)
		visitReports.PUT("/activity-types/:id", activityTypeHandler.Update)
		visitReports.DELETE("/activity-types/:id", activityTypeHandler.Delete)
	}

	// Mobile-specific routes
	mobile := router.Group("/mobile")
	mobile.Use(middleware.AuthMiddleware(jwtManager))
	{
		mobileVisitReports := mobile.Group("/visit-reports")
		{
			// Get form data for visit report creation (accounts, contacts, deals, leads)
			mobileVisitReports.GET("/form-data", visitReportHandler.GetFormData)
			// Get visit reports for logged-in user (sales rep)
			mobileVisitReports.GET("/my-visit-reports", visitReportHandler.GetMyVisitReports)
			// Create visit report for mobile (validates ownership - only for logged-in user)
			mobileVisitReports.POST("", visitReportHandler.CreateMobile)
			// Get visit report by ID for mobile (validates ownership)
			mobileVisitReports.GET("/:id", visitReportHandler.GetByIDMobile)
			// Update visit report for mobile (validates ownership and status)
			mobileVisitReports.PUT("/:id", visitReportHandler.UpdateMobile)
			// Check-in/out endpoints for mobile (validates ownership)
			mobileVisitReports.POST("/:id/check-in", visitReportHandler.CheckInMobile)
			mobileVisitReports.POST("/:id/check-out", visitReportHandler.CheckOutMobile)
			// Submit visit report for mobile (validates ownership and status)
			mobileVisitReports.PATCH("/:id/submit", visitReportHandler.SubmitMobile)
			// Delete visit report for mobile (validates ownership)
			mobileVisitReports.DELETE("/:id", visitReportHandler.DeleteMobile)
			// Upload photo for mobile (validates ownership, only supports file upload)
			mobileVisitReports.POST("/:id/photos", visitReportHandler.UploadPhotoMobile)
		}
	}
}
