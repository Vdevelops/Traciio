package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupGeocodingRoutes sets up geocoding routes
func SetupGeocodingRoutes(router *gin.RouterGroup, geocodingHandler *handlers.GeocodingHandler, jwtManager *jwt.JWTManager) {
	geocoding := router.Group("/geocoding")
	geocoding.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Geocode address to coordinates
		geocoding.POST("/geocode", geocodingHandler.Geocode)
	}
}
