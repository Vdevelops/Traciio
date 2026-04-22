package middleware

import (
	"fmt"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gin-gonic/gin"
)

// HSTSMiddleware sets HTTP Strict Transport Security headers
// Only applies to HTTPS connections to prevent MITM attacks
func HSTSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only set HSTS header if connection is HTTPS
		if c.Request.TLS != nil {
			hstsConfig := config.AppConfig.HSTS
			
			// Build HSTS header value
			headerValue := fmt.Sprintf("max-age=%d", hstsConfig.MaxAge)
			
			if hstsConfig.IncludeSubDomains {
				headerValue += "; includeSubDomains"
			}
			
			if hstsConfig.Preload {
				headerValue += "; preload"
			}
			
			c.Header("Strict-Transport-Security", headerValue)
		}
		
		c.Next()
	}
}

