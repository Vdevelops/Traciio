package middleware

import (
	"fmt"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds comprehensive security headers to all responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		// Restricts where resources can be loaded from
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " + // Allow inline scripts for React
			"style-src 'self' 'unsafe-inline'; " + // Allow inline styles
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self' https:; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Header("Content-Security-Policy", csp)

		// X-Content-Type-Options
		// Prevents MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options
		// Prevents clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection
		// Enables browser's XSS filter (legacy browsers)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		// Controls how much referrer information is sent
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (formerly Feature-Policy)
		// Controls which browser features can be used
		c.Header("Permissions-Policy", "geolocation=(self), microphone=(), camera=()")

		// Clear potentially dangerous headers
		c.Header("X-Powered-By", "")
		c.Header("Server", "")

		c.Next()
	}
}

// HSTSMiddleware sets HTTP Strict Transport Security headers
// Already exists but ensure it's properly configured
func HSTSMiddlewareEnhanced() gin.HandlerFunc {
	hstsCfg := config.AppConfig.HSTS
	maxAge := hstsCfg.MaxAge
	if maxAge == 0 {
		maxAge = 31536000 // 1 year default
	}

	return func(c *gin.Context) {
		// Only set HSTS on HTTPS connections
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			hstsValue := fmt.Sprintf("max-age=%d", maxAge)
			if hstsCfg.IncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			if hstsCfg.Preload {
				hstsValue += "; preload"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}
		c.Next()
	}
}

// Add fmt import at the top
