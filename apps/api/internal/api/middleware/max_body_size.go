package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
)

const (
	// DefaultMaxBodySize is the default maximum request body size (10MB)
	DefaultMaxBodySize = 10 * 1024 * 1024 // 10MB

	// MaxFileUploadSize is the maximum file upload size (50MB)
	// This is already configured in main.go via router.MaxMultipartMemory
	MaxFileUploadSize = 50 * 1024 * 1024 // 50MB
)

// MaxBodySizeMiddleware limits the size of request bodies to prevent OOM attacks.
// This is critical for high-throughput APIs (100k+ req/s) to prevent memory exhaustion.
//
// Performance optimizations:
// 1. Checks Content-Length header first (O(1), no body reading)
// 2. Uses io.LimitReader to prevent reading beyond limit
// 3. Returns early before expensive body parsing
// 4. Zero-copy when size is acceptable
//
// For file uploads, use router.MaxMultipartMemory instead (configured in main.go).
func MaxBodySizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for GET, HEAD, OPTIONS (no body)
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Skip for multipart form data (file uploads)
		// These are handled by router.MaxMultipartMemory
		contentType := c.GetHeader("Content-Type")
		if len(contentType) > 19 && contentType[:19] == "multipart/form-data" {
			c.Next()
			return
		}

		// Fast path: Check Content-Length header (O(1))
		// This prevents reading the body if size is already known
		if contentLength := c.GetHeader("Content-Length"); contentLength != "" {
			size, err := strconv.ParseInt(contentLength, 10, 64)
			if err == nil && size > maxSize {
				errors.ErrorResponse(c, "REQUEST_BODY_TOO_LARGE", map[string]interface{}{
					"max_size":      maxSize,
					"content_length": size,
					"max_size_mb":   maxSize / 1024 / 1024,
				}, nil)
				c.Abort()
				return
			}
		}

		// Slow path: Wrap body with LimitReader to enforce limit during reading
		// This handles cases where Content-Length is not set or is incorrect
		// io.LimitReader is zero-overhead if body size <= maxSize
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		}

		c.Next()
	}
}
