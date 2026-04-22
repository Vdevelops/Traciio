package routes

import (
	"github.com/gin-gonic/gin"

	areamappinghandler "github.com/gilabs/crm-healthcare/api/internal/api/area_mapping"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
)

// SetupAreaMappingRoutes sets up all area mapping related routes
func SetupAreaMappingRoutes(
	rg *gin.RouterGroup,
	areaMappingHandler *areamappinghandler.Handler,
	jwtManager *jwt.JWTManager,
) {
	areaMapping := rg.Group("/area-mapping")
	areaMapping.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Location capture endpoints
		areaMapping.POST("/capture", areaMappingHandler.CaptureLocation)
		areaMapping.GET("/captures", areaMappingHandler.GetAreaCaptures)

		// Territory management endpoints
		areaMapping.POST("/territories", areaMappingHandler.CreateTerritory)
		areaMapping.GET("/territories", areaMappingHandler.GetTerritories)
		areaMapping.GET("/territories/:id", areaMappingHandler.GetTerritoryByID)
		areaMapping.PUT("/territories/:id", areaMappingHandler.UpdateTerritory)
		areaMapping.DELETE("/territories/:id", areaMappingHandler.DeleteTerritory)

		// Spatial analysis endpoints
		areaMapping.GET("/check-territory", areaMappingHandler.CheckPointInTerritory)
		areaMapping.GET("/coverage", areaMappingHandler.GetCoverageAnalysis)
		areaMapping.GET("/heatmap", areaMappingHandler.GetHeatmap)

		// Territory assignment endpoint
		areaMapping.POST("/assign-territory", areaMappingHandler.AssignTerritory)
	}
}
