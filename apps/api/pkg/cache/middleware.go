// Package cache provides cache integration middleware for HTTP handlers.
//
// This file implements caching middleware that can be applied to HTTP endpoints
// for automatic response caching and cache invalidation.
package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheMiddlewareConfig holds configuration for cache middleware
type CacheMiddlewareConfig struct {
	// TTL is the cache duration
	TTL time.Duration
	// KeyPrefix is prepended to all cache keys
	KeyPrefix string
	// ExcludePaths are paths that should not be cached
	ExcludePaths []string
	// ExcludeMethods are HTTP methods that should not be cached (default: POST, PUT, PATCH, DELETE)
	ExcludeMethods []string
	// IncludeQueryParams determines if query params are included in cache key
	IncludeQueryParams bool
	// IncludeUserID determines if user ID is included in cache key
	IncludeUserID bool
}

// DefaultCacheMiddlewareConfig returns default configuration
func DefaultCacheMiddlewareConfig() *CacheMiddlewareConfig {
	return &CacheMiddlewareConfig{
		TTL:                TTLListShort,
		KeyPrefix:          "api:",
		ExcludePaths:       []string{"/health", "/metrics"},
		ExcludeMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		IncludeQueryParams: true,
		IncludeUserID:      true,
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheMiddleware creates a caching middleware for API responses
func CacheMiddleware(config *CacheMiddlewareConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCacheMiddlewareConfig()
	}

	ec := Advanced()

	return func(c *gin.Context) {
		// Skip if cache is not enabled
		if !ec.IsEnabled() {
			c.Next()
			return
		}

		// Skip excluded methods
		for _, method := range config.ExcludeMethods {
			if c.Request.Method == method {
				c.Next()
				return
			}
		}

		// Skip excluded paths
		for _, path := range config.ExcludePaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// Generate cache key
		cacheKey := generateCacheKey(c, config)

		// Try to get from cache
		var cachedResponse CachedHTTPResponse
		if found, _ := ec.Get(cacheKey, &cachedResponse); found {
			// Return cached response
			for key, values := range cachedResponse.Headers {
				for _, value := range values {
					c.Header(key, value)
				}
			}
			c.Header("X-Cache", "HIT")
			c.Data(cachedResponse.StatusCode, cachedResponse.ContentType, cachedResponse.Body)
			c.Abort()
			return
		}

		// Cache miss - capture response
		c.Header("X-Cache", "MISS")
		w := &responseWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		// Only cache successful responses
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			response := CachedHTTPResponse{
				StatusCode:  c.Writer.Status(),
				ContentType: c.Writer.Header().Get("Content-Type"),
				Body:        w.body.Bytes(),
				Headers:     make(map[string][]string),
				CachedAt:    time.Now(),
			}

			// Copy relevant headers
			for key, values := range c.Writer.Header() {
				if shouldCacheHeader(key) {
					response.Headers[key] = values
				}
			}

			// Store in cache asynchronously
			go func() {
				_ = ec.Set(cacheKey, response, config.TTL)
			}()
		}
	}
}

// CachedHTTPResponse represents a cached HTTP response
type CachedHTTPResponse struct {
	StatusCode  int                 `msgpack:"status_code"`
	ContentType string              `msgpack:"content_type"`
	Body        []byte              `msgpack:"body"`
	Headers     map[string][]string `msgpack:"headers"`
	CachedAt    time.Time           `msgpack:"cached_at"`
}

// generateCacheKey creates a unique cache key for the request
func generateCacheKey(c *gin.Context, config *CacheMiddlewareConfig) string {
	var keyParts []string

	keyParts = append(keyParts, config.KeyPrefix)
	keyParts = append(keyParts, c.Request.Method)
	keyParts = append(keyParts, c.Request.URL.Path)

	if config.IncludeQueryParams {
		// Normalize query params so different ordering doesn't fragment cache keys.
		// url.Values.Encode() sorts keys (and values) deterministically.
		if normalized := c.Request.URL.Query().Encode(); normalized != "" {
			keyParts = append(keyParts, normalized)
		}
	}

	if config.IncludeUserID {
		if userID, exists := c.Get("user_id"); exists {
			keyParts = append(keyParts, fmt.Sprintf("user:%v", userID))
		}
	}

	key := strings.Join(keyParts, ":")

	// Hash if key is too long
	if len(key) > 200 {
		hash := sha256.Sum256([]byte(key))
		return config.KeyPrefix + hex.EncodeToString(hash[:])
	}

	return key
}

// shouldCacheHeader determines if a header should be included in cached response
func shouldCacheHeader(key string) bool {
	// Headers to include in cache
	includedHeaders := map[string]bool{
		"Content-Type":     true,
		"Content-Language": true,
		"Cache-Control":    true,
		"ETag":             true,
		"Last-Modified":    true,
	}
	return includedHeaders[key]
}

// InvalidateCacheMiddleware creates middleware that invalidates cache on write operations
func InvalidateCacheMiddleware(patterns ...string) gin.HandlerFunc {
	ec := Advanced()

	return func(c *gin.Context) {
		c.Next()

		// Only invalidate on successful write operations
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			for _, pattern := range patterns {
				go func(p string) {
					_ = ec.DeletePattern(p)
				}(pattern)
			}
		}
	}
}

// CacheControlMiddleware adds cache-control headers for browser caching
func CacheControlMiddleware(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only add for GET requests
		if c.Request.Method == http.MethodGet {
			c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		}
		c.Next()
	}
}

// NoCacheMiddleware adds no-cache headers
func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// CacheableRequest checks if request body can be used for cache key
type CacheableRequest interface {
	CacheKey() string
}

// ReadCacheableBody reads and caches request body for potential cache key generation
func ReadCacheableBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// WarmupCacheHandler provides an endpoint to trigger cache warmup for specific resources
func WarmupCacheHandler(warmupFuncs map[string]func() error) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource := c.Param("resource")

		if warmupFunc, exists := warmupFuncs[resource]; exists {
			if err := warmupFunc(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   fmt.Sprintf("warmup failed: %v", err),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": fmt.Sprintf("cache warmup completed for %s", resource),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("unknown resource: %s", resource),
		})
	}
}

// CacheStatsHandler provides an endpoint to view cache statistics
func CacheStatsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ec := Advanced()

		if !ec.IsEnabled() {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"enabled": false,
					"message": "cache is disabled",
				},
			})
			return
		}

		baseMetrics := Client.GetMetrics()
		advancedMetrics := ec.GetAdvancedMetrics()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"enabled":        true,
				"hit_rate":       Client.GetHitRate(),
				"avg_latency_us": Client.GetAvgLatency(),
				"base_metrics": gin.H{
					"hits":          baseMetrics.Hits,
					"misses":        baseMetrics.Misses,
					"errors":        baseMetrics.Errors,
					"timeouts":      baseMetrics.Timeouts,
					"cb_open":       baseMetrics.CBOpen,
					"cb_half_open":  baseMetrics.CBHalfOpen,
					"pipeline_ops":  baseMetrics.PipelineOps,
					"compressions":  baseMetrics.Compressions,
					"pool_hits":     baseMetrics.PoolHits,
					"pool_misses":   baseMetrics.PoolMisses,
					"pool_timeouts": baseMetrics.PoolTimeouts,
					"total_conns":   baseMetrics.TotalConns,
					"idle_conns":    baseMetrics.IdleConns,
					"stale_conns":   baseMetrics.StaleConns,
				},
				"enterprise_metrics": gin.H{
					"warmup_tasks":     advancedMetrics.WarmupTasks,
					"warmup_success":   advancedMetrics.WarmupSuccess,
					"warmup_failures":  advancedMetrics.WarmupFailures,
					"batch_ops":        advancedMetrics.BatchOps,
					"cache_aside_hit":  advancedMetrics.CacheAsideHit,
					"cache_aside_miss": advancedMetrics.CacheAsideMiss,
					"invalidations":    advancedMetrics.Invalidations,
					"pattern_deletes":  advancedMetrics.PatternDeletes,
				},
			},
		})
	}
}
