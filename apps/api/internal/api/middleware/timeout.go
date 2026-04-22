package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware adds a strict timeout to the request context.
// Critical for A+ Availability and preventing DOS via slow queries.
//
// IMPORTANT: This middleware sets context.WithTimeout so that downstream
// handlers/repositories using db.WithContext(ctx) will have their DB queries
// cancelled when the deadline expires. The handler chain runs synchronously
// (no goroutine spawned) to avoid goroutine leaks under high load.
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update the request with the new context
		c.Request = c.Request.WithContext(ctx)

		// Run handler chain synchronously — no goroutine leak.
		// Context cancellation propagates to DB queries via db.WithContext(ctx).
		c.Next()

		// After handler completes, check if context deadline was exceeded
		if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timed out",
				"code":  "REQUEST_TIMEOUT",
			})
		}
	}
}

// TimeoutMiddlewareByPath applies different timeouts per URL path prefix.
// Note: a longer timeout cannot override a shorter parent context deadline.
// Use this instead of stacking multiple TimeoutMiddleware instances.
//
// Runs handler chain synchronously to prevent goroutine leaks under high load.
func TimeoutMiddlewareByPath(defaultTimeout time.Duration, prefixTimeouts map[string]time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		selectedTimeout := defaultTimeout
		path := c.Request.URL.Path
		bestPrefixLen := -1
		for prefix, timeout := range prefixTimeouts {
			if strings.HasPrefix(path, prefix) {
				if len(prefix) > bestPrefixLen {
					bestPrefixLen = len(prefix)
					selectedTimeout = timeout
				}
			}
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), selectedTimeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		// Run handler chain synchronously — no goroutine leak.
		c.Next()

		// After handler completes, check if context deadline was exceeded
		if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timed out",
				"code":  "REQUEST_TIMEOUT",
			})
		}
	}
}
