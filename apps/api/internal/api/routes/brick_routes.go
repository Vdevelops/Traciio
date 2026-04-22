package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupBrickRoutes(router *gin.RouterGroup, brickHandler *handlers.BrickHandler, brickTargetDistributionHandler *handlers.BrickTargetDistributionHandler, brickAnalyticsHandler *handlers.BrickAnalyticsHandler, jwtManager *jwt.JWTManager) {
	bricks := router.Group("/bricks")
	bricks.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Brick CRUD
		bricks.GET("", brickHandler.List)
		bricks.POST("", brickHandler.Create)
		
		// Analytics endpoints (must be before /:id to avoid conflict)
		bricks.GET("/:id/performance", brickAnalyticsHandler.GetBrickPerformance)
		bricks.GET("/performance", brickAnalyticsHandler.ListBrickPerformance)
		
		// Get sales in brick (must be before /:id to avoid conflict)
		bricks.GET("/:id/sales", brickHandler.GetSales)
		bricks.POST("/:id/sales", brickHandler.AssignSales)
		bricks.DELETE("/:id/sales", brickHandler.UnassignSales)
		
		// Brick target distributions (nested routes - must be before /:id to avoid conflict)
		bricks.GET("/:id/targets/:year/:month", brickTargetDistributionHandler.GetBrickTargetWithDistributions)
		bricks.POST("/:id/targets/:target_id/distribute", brickTargetDistributionHandler.DistributeTarget)
		bricks.PATCH("/:id/targets/:target_id/distributions/:distribution_id", brickTargetDistributionHandler.UpdateDistribution)
		bricks.DELETE("/:id/targets/:target_id/distributions/:distribution_id", brickTargetDistributionHandler.DeleteDistribution)
		
		// Brick CRUD by ID (must be last to avoid conflicts)
		bricks.GET("/:id", brickHandler.GetByID)
		bricks.PATCH("/:id", brickHandler.Update)
		bricks.DELETE("/:id", brickHandler.Delete)
	}
}

