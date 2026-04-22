package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Logger.Printf(
			"[%d] %s %s %s %s",
			statusCode,
			method,
			path,
			clientIP,
			latency,
		)
	}
}

