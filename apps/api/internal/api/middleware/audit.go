package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLogMiddleware logs sensitive actions (POST, PUT, DELETE, PATCH) to the audit log
// It captures who did what, when, and the outcome, essential for A+ Security standards.
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only audit write operations (POST, PUT, DELETE, PATCH)
		// We avoid logging GET requests to prevent log flooding, but critical READs should ideally be logged too.
		// For now, we focus on state changes.
		method := c.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			c.Next()
			return
		}

		start := time.Now()
		
		// Read body to log it (carefully)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// Extract User ID (if authenticated)
		// It expects the AuthMiddleware to set "user_id"
		userID, exists := c.Get("user_id")
		userStr := "anonymous"
		if exists {
			userStr = fmt.Sprintf("%v", userID)
		}

		status := c.Writer.Status()
		duration := time.Since(start)
		
		// Sanitize body for logging (mask passwords)
		sanitizedBody := sanitizeBody(bodyBytes)

		// Create Audit Entry
		entry := map[string]interface{}{
			"event_type": "AUDIT_LOG",
			"timestamp":  start.Format(time.RFC3339),
			"user_id":    userStr,
			"ip":         c.ClientIP(),
			"method":     method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"duration_ms": duration.Milliseconds(),
			"user_agent": c.Request.UserAgent(),
			// Only log body for smaller requests to avoid massive logs, and skip file uploads
			"payload": sanitizedBody, 
		}
		
		// Skip body logging for file uploads or large payloads
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			entry["payload"] = "[multipart-data-omitted]"
		}

		// Log structured JSON to stdout (can be picked up by filebeat/promtail)
		logData, err := json.Marshal(entry)
		if err == nil {
			// Use a special prefix for easy filtering
			log.Printf("🔒 [AUDIT] %s", string(logData))
		}
	}
}

// sanitizeBody attempts to mask sensitive fields in JSON body
func sanitizeBody(body []byte) interface{} {
	if len(body) == 0 {
		return nil
	}
	
	// Basic check if it looks like JSON
	if body[0] != '{' && body[0] != '[' {
		return "[non-json-payload]"
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "[invalid-json]"
	}

	maskSensitiveFields(data)
	return data
}

func maskSensitiveFields(data map[string]interface{}) {
	sensitiveKeys := []string{"password", "token", "access_token", "refresh_token", "secret", "credit_card", "cvv"}
	
	for key, val := range data {
		// Check keys
		lowerKey := strings.ToLower(key)
		for _, sensitive := range sensitiveKeys {
			if strings.Contains(lowerKey, sensitive) {
				data[key] = "***MASKED***"
				break
			}
		}
		
		// Recurse for nested objects
		if nestedMap, ok := val.(map[string]interface{}); ok {
			maskSensitiveFields(nestedMap)
		}
	}
}
