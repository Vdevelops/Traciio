package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupPipelineRoutes sets up pipeline routes
func SetupPipelineRoutes(router *gin.RouterGroup, pipelineHandler *handlers.PipelineHandler, dealHandler *handlers.DealHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc) {
	// Standard Web Routes
	pipelineGroup := router.Group("/pipelines")
	pipelineGroup.Use(middleware.AuthMiddleware(jwtManager))
	{
		pipelineGroup.GET("", pipelineHandler.ListStages)
		pipelineGroup.GET("/stages", pipelineHandler.ListStages)
		pipelineGroup.GET("/stages/:id", pipelineHandler.GetStageByID)
		pipelineGroup.POST("/stages", middleware.RateLimitMiddleware("mutation"), pipelineHandler.CreateStage)
		pipelineGroup.PUT("/stages/:id", middleware.RateLimitMiddleware("mutation"), pipelineHandler.UpdateStage)
		pipelineGroup.DELETE("/stages/:id", middleware.RateLimitMiddleware("mutation"), pipelineHandler.DeleteStage)
		pipelineGroup.PUT("/stages/order", middleware.RateLimitMiddleware("mutation"), pipelineHandler.UpdateStagesOrder)

		pipelineGroup.GET("/summary", pipelineHandler.GetSummary)
		pipelineGroup.GET("/forecast", pipelineHandler.GetForecast)
	}

	dealsGroup := router.Group("/deals")
	dealsGroup.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		dealsGroup.GET("/by-stage", dealHandler.ListByStage)
		dealsGroup.GET("", dealHandler.List)
		dealsGroup.GET("/:id", dealHandler.GetByID)
		dealsGroup.POST("", middleware.RateLimitMiddleware("mutation"), dealHandler.Create)
		dealsGroup.PUT("/:id", middleware.RateLimitMiddleware("mutation"), dealHandler.Update)
		dealsGroup.DELETE("/:id", middleware.RateLimitMiddleware("mutation"), dealHandler.Delete)
		dealsGroup.PATCH("/:id/move", middleware.RateLimitMiddleware("mutation"), dealHandler.Move)
		dealsGroup.POST("/:id/move-stage", middleware.RateLimitMiddleware("mutation"), pipelineHandler.MoveStage)
		dealsGroup.GET("/:id/history", pipelineHandler.GetDealHistory)
		dealsGroup.GET("/:id/visit-reports", dealHandler.GetVisitReportsByDeal)
		dealsGroup.GET("/:id/activities", dealHandler.GetActivitiesByDeal)
	}

	// Mobile Routes
	mobile := router.Group("/mobile")
	{
		pipelineMob := mobile.Group("/pipelines")
		pipelineMob.Use(middleware.AuthMiddleware(jwtManager))
		{
			pipelineMob.GET("/stages", pipelineHandler.ListStages)
			pipelineMob.GET("/stages/:id", pipelineHandler.GetStageByID)
			pipelineMob.GET("/summary", pipelineHandler.GetSummary)
			pipelineMob.GET("/forecast", pipelineHandler.GetForecast)
			pipelineMob.GET("/form-data", dealHandler.GetFormData)
		}

		dealsMob := mobile.Group("/deals")
		dealsMob.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
		{
			dealsMob.GET("/by-stage", dealHandler.ListByStage)
			dealsMob.GET("", dealHandler.List)
			dealsMob.GET("/:id", dealHandler.GetByID)
			dealsMob.POST("", middleware.RateLimitMiddleware("mutation"), dealHandler.Create)
			dealsMob.PUT("/:id", middleware.RateLimitMiddleware("mutation"), dealHandler.Update)
			dealsMob.DELETE("/:id", middleware.RateLimitMiddleware("mutation"), dealHandler.Delete)
			dealsMob.PATCH("/:id/move", middleware.RateLimitMiddleware("mutation"), dealHandler.Move)
			dealsMob.POST("/:id/move-stage", middleware.RateLimitMiddleware("mutation"), pipelineHandler.MoveStage)
			dealsMob.GET("/:id/history", pipelineHandler.GetDealHistory)
			dealsMob.GET("/:id/visit-reports", dealHandler.GetVisitReportsByDeal)
			dealsMob.GET("/:id/activities", dealHandler.GetActivitiesByDeal)
		}
	}
}
