package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupKPIRoutes(router *gin.RouterGroup, kpiHandler *handlers.KPIHandler, jwtManager *jwt.JWTManager, scopeMiddleware gin.HandlerFunc, brickRepo interfaces.BrickRepository) {
	kpi := router.Group("/kpi")
	kpi.Use(middleware.AuthMiddleware(jwtManager), scopeMiddleware)
	{
		kpi.GET("/sales-rep", middleware.KPIRepScopeMiddleware(), kpiHandler.GetSalesRepScorecard)
		kpi.GET("/sales-manager", middleware.KPIManagerScopeMiddleware(brickRepo), kpiHandler.GetSalesManagerScorecard)
	}
}