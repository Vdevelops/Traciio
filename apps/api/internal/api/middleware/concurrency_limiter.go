package middleware

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// ConcurrencyLimiterMiddleware limits the number of concurrently processed requests.
// This prevents goroutine/memory/DB-connection exhaustion under extreme load.
// Requests that exceed the limit are immediately rejected with 503 Service Unavailable.
//
// Usage:
//
//	router.Use(middleware.ConcurrencyLimiterMiddleware(500))
func ConcurrencyLimiterMiddleware(maxConcurrent int) gin.HandlerFunc {
	sem := make(chan struct{}, maxConcurrent)
	var rejected atomic.Int64

	log.Printf("⚡ Concurrency limiter initialized: max %d concurrent requests", maxConcurrent)

	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}:
			// Acquired slot — process request
			defer func() { <-sem }()
			c.Next()

		default:
			// All slots occupied — reject immediately
			count := rejected.Add(1)
			if count%100 == 1 {
				log.Printf("⚠️  Concurrency limiter: rejected request #%d (max=%d) path=%s",
					count, maxConcurrent, c.Request.URL.Path)
			}
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "SERVER_BUSY",
					"message": "Server is at capacity, please try again later / Server sedang sibuk, silakan coba lagi nanti",
				},
			})
		}
	}
}
