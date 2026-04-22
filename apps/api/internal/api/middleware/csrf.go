package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gin-gonic/gin"
)

// CSRFToken represents a CSRF token with metadata
type CSRFToken struct {
	Token     string
	CreatedAt time.Time
	UserID    string
}

// CSRFStore stores CSRF tokens in memory with automatic cleanup
type CSRFStore struct {
	mu     sync.RWMutex
	tokens map[string]*CSRFToken
}

var csrfStore = &CSRFStore{
	tokens: make(map[string]*CSRFToken),
}

// CSRFConfig holds CSRF middleware configuration
type CSRFConfig struct {
	TokenLength    int
	TokenLifetime  time.Duration
	CookieName     string
	HeaderName     string
	SkipMethods    []string
	SkipPaths      []string
	ErrorHandler   func(*gin.Context)
}

// DefaultCSRFConfig returns default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:   32,
		TokenLifetime: 24 * time.Hour,
		CookieName:    "csrf_token",
		HeaderName:    "X-CSRF-Token",
		SkipMethods:   []string{"GET", "HEAD", "OPTIONS"},
		SkipPaths:     []string{"/api/v1/auth/login", "/api/v1/auth/refresh"},
		ErrorHandler: func(c *gin.Context) {
			errors.ErrorResponse(c, "CSRF_TOKEN_INVALID", map[string]interface{}{
				"message": "CSRF token validation failed",
			}, nil)
		},
	}
}

// generateToken generates a cryptographically secure random token
func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// cleanupExpiredTokens removes expired tokens from store
func (s *CSRFStore) cleanupExpiredTokens(lifetime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, data := range s.tokens {
		if now.Sub(data.CreatedAt) > lifetime {
			delete(s.tokens, token)
		}
	}
}

// storeToken stores a CSRF token
func (s *CSRFStore) storeToken(token string, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[token] = &CSRFToken{
		Token:     token,
		CreatedAt: time.Now(),
		UserID:    userID,
	}
}

// validateToken validates a CSRF token
func (s *CSRFStore) validateToken(token string, userID string, lifetime time.Duration) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.tokens[token]
	if !exists {
		return false
	}

	// Check if token expired
	if time.Since(data.CreatedAt) > lifetime {
		return false
	}

	// Check if token belongs to the same user
	if data.UserID != userID {
		return false
	}

	return true
}

// removeToken removes a CSRF token from store
func (s *CSRFStore) removeToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}

// CSRFMiddleware provides CSRF protection using double-submit cookie pattern
// This is suitable for SPA applications
func CSRFMiddleware(config ...CSRFConfig) gin.HandlerFunc {
	cfg := DefaultCSRFConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			csrfStore.cleanupExpiredTokens(cfg.TokenLifetime)
		}
	}()

	return func(c *gin.Context) {
		// Skip CSRF check for certain methods
		for _, method := range cfg.SkipMethods {
			if c.Request.Method == method {
				c.Next()
				return
			}
		}

		// Skip CSRF check for certain paths
		for _, path := range cfg.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// Get user ID from context (set by AuthMiddleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			// If no user ID, skip CSRF (public endpoints)
			c.Next()
			return
		}
		userID, ok := userIDInterface.(string)
		if !ok {
			cfg.ErrorHandler(c)
			c.Abort()
			return
		}

		// Get token from header
		headerToken := c.GetHeader(cfg.HeaderName)
		if headerToken == "" {
			cfg.ErrorHandler(c)
			c.Abort()
			return
		}

		// Get token from cookie
		cookieToken, err := c.Cookie(cfg.CookieName)
		if err != nil {
			cfg.ErrorHandler(c)
			c.Abort()
			return
		}

		// Compare tokens using constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
			cfg.ErrorHandler(c)
			c.Abort()
			return
		}

		// Validate token in store
		if !csrfStore.validateToken(headerToken, userID, cfg.TokenLifetime) {
			cfg.ErrorHandler(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateCSRFToken generates and returns a new CSRF token
// This should be called after successful login
func GenerateCSRFToken(c *gin.Context, config ...CSRFConfig) (string, error) {
	cfg := DefaultCSRFConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID type")
	}

	// Generate token
	token, err := generateToken(cfg.TokenLength)
	if err != nil {
		return "", err
	}

	// Store token
	csrfStore.storeToken(token, userID)

	// Set cookie with secure flags
	secure := os.Getenv("ENV") == "production"
	c.SetSameSite(3) // Strict
	c.SetCookie(
		cfg.CookieName,
		token,
		int(cfg.TokenLifetime.Seconds()),
		"/",
		"",
		secure, // Secure (HTTPS only)
		true,   // HttpOnly
	)

	return token, nil
}

// RevokeCSRFToken revokes a CSRF token (useful for logout)
func RevokeCSRFToken(c *gin.Context, config ...CSRFConfig) {
	cfg := DefaultCSRFConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Get token from cookie
	token, err := c.Cookie(cfg.CookieName)
	if err == nil && token != "" {
		csrfStore.removeToken(token)
	}

	// Clear cookie with SameSite protection
	secure := os.Getenv("ENV") == "production"
	c.SetSameSite(3) // Strict
	c.SetCookie(
		cfg.CookieName,
		"",
		-1,
		"/",
		"",
		secure, // Secure (HTTPS only)
		true,   // HttpOnly
	)
}
