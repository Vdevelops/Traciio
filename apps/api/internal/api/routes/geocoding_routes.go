package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
)

// SetupGeocodingRoutes sets up geocoding routes
func SetupGeocodingRoutes(router *gin.RouterGroup, geocodingHandler *handlers.GeocodingHandler, jwtManager *jwt.JWTManager) {
	geocoding := router.Group("/geocoding")
	geocoding.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Geocode address to coordinates
		geocoding.POST("/geocode", geocodingHandler.Geocode)
	}

	// Mobile-specific routes
	mobile := router.Group("/mobile")
	mobile.Use(middleware.AuthMiddleware(jwtManager))
	{
		mobileGeocoding := mobile.Group("/geocoding")
		{
			// Geocode address to coordinates (mobile)
			mobileGeocoding.POST("/geocode", geocodingHandler.Geocode)
		}
	}
}





