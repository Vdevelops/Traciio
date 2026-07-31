package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets up CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	
	// Get allowed origins from environment variable or use defaults
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		// Production origins (add more if needed)
		"https://api.gilabs.id",
		"https://crm-demo.gilabs.id",
	}
	
	// Add production origins from environment variable
	// Format: comma-separated list, e.g., "https://gilabs.id,https://www.gilabs.id"
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		origins := strings.Split(envOrigins, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				// Check if already exists to avoid duplicates
				exists := false
				for _, existing := range allowedOrigins {
					if existing == trimmed {
						exists = true
						break
					}
				}
				if !exists {
					allowedOrigins = append(allowedOrigins, trimmed)
				}
			}
		}
	}
	
	// Use an AllowOriginFunc that accepts configured origins and
	// localhost/127.0.0.1 variants during development so preflight
	// requests from local frontends succeed even when the exact
	// host/port combination may vary.
	config.AllowOriginFunc = func(origin string) bool {
		if origin == "" {
			return false
		}
		// Exact match against configured origins
		for _, a := range allowedOrigins {
			if a == origin {
				return true
			}
		}
		// Allow common local development hosts
		if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") {
			return true
		}
		return false
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Accept",
		"X-Requested-With",
		"X-Request-ID",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{
		"X-Request-ID",
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
	}
	// Set max age for preflight requests (24 hours)
	config.MaxAge = 86400

	return cors.New(config)
}

