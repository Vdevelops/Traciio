package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
)

// SetupMasterDataRoutes sets up master data routes
// Note: Currently empty as diagnosis, procedures, and categories have been removed
// This can be used for future Sales CRM master data modules
func SetupMasterDataRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	masterData := router.Group("/master-data")
	masterData.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Master data routes will be added here for Sales CRM modules
	}
}

