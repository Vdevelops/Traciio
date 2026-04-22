package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

var redisFixedWindowScript = redis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if current == 1 then
  redis.call('EXPIRE', KEYS[1], ARGV[1])
end
local ttl = redis.call('TTL', KEYS[1])
return {current, ttl}
`)

type redisFixedWindowLimiter struct {
	client redis.UniversalClient
}

func getRedisFixedWindowLimiter() *redisFixedWindowLimiter {
	if cachepkg.Client == nil || !cachepkg.Client.IsEnabled() {
		return nil
	}
	client := cachepkg.Client.GetClient()
	if client == nil {
		return nil
	}
	return &redisFixedWindowLimiter{client: client}
}

func (l *redisFixedWindowLimiter) allow(ctx context.Context, key string, limit int, windowSec int) (allowed bool, remaining int, resetUnix int64, err error) {
	if limit <= 0 || windowSec <= 0 {
		return true, 0, time.Now().Unix(), nil
	}

	res, err := redisFixedWindowScript.Run(ctx, l.client, []string{key}, windowSec).Result()
	if err != nil {
		return false, 0, 0, err
	}

	arr, ok := res.([]interface{})
	if !ok || len(arr) < 2 {
		return false, 0, 0, fmt.Errorf("unexpected redis ratelimit result")
	}

	current, ok1 := arr[0].(int64)
	ttl, ok2 := arr[1].(int64)
	if !ok1 || !ok2 {
		return false, 0, 0, fmt.Errorf("unexpected redis ratelimit result values")
	}

	if ttl < 0 {
		ttl = int64(windowSec)
	}

	allowed = current <= int64(limit)
	if allowed {
		remaining = limit - int(current)
		if remaining < 0 {
			remaining = 0
		}
	} else {
		remaining = 0
	}
	resetUnix = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	return allowed, remaining, resetUnix, nil
}

func extractEmailFromRequestBody(c *gin.Context) string {
	var loginReq struct {
		Email string `json:"email"`
	}

	if c.Request.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	// Restore body for handler
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return ""
	}
	if json.Unmarshal(bodyBytes, &loginReq) != nil {
		return ""
	}
	return strings.TrimSpace(strings.ToLower(loginReq.Email))
}

func tryRedisRateLimit(c *gin.Context, limitType string, ip string) (handled bool) {
	rl := getRedisFixedWindowLimiter()
	if rl == nil {
		return false
	}
	ctx := c.Request.Context()
	keyPrefix := "ratelimit:v1:" + limitType + ":"

	// Multi-level login: global -> email -> IP
	if limitType == "login" {
		globalRule := config.AppConfig.RateLimit.LoginGlobal
		allowed, remaining, resetTime, err := rl.allow(ctx, keyPrefix+"global", globalRule.Requests, globalRule.Window)
		if err != nil {
			return false
		}
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", globalRule.Requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
		if !allowed {
			errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
			c.Abort()
			return true
		}

		email := extractEmailFromRequestBody(c)
		if email != "" {
			emailRule := config.AppConfig.RateLimit.LoginByEmail
			allowed, remaining, resetTime, err := rl.allow(ctx, keyPrefix+"email:"+email, emailRule.Requests, emailRule.Window)
			if err != nil {
				return false
			}
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", emailRule.Requests))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
			if !allowed {
				errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
				c.Abort()
				return true
			}
		}

		rule := config.AppConfig.RateLimit.Login
		allowed, remaining, resetTime, err = rl.allow(ctx, keyPrefix+"ip:"+ip, rule.Requests, rule.Window)
		if err != nil {
			return false
		}
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
		if !allowed {
			errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
			c.Abort()
			return true
		}

		c.Next()
		return true
	}

	var rule config.RateLimitRule
	switch limitType {
	case "refresh":
		rule = config.AppConfig.RateLimit.Refresh
	case "upload":
		rule = config.AppConfig.RateLimit.Upload
	case "general":
		rule = config.AppConfig.RateLimit.General
	case "public":
		rule = config.AppConfig.RateLimit.Public
	case "mutation":
		rule = config.AppConfig.RateLimit.Mutation
	case "high_volume":
		rule = config.AppConfig.RateLimit.HighVolume
	default:
		rule = config.AppConfig.RateLimit.General
	}

	allowed, remaining, resetTime, err := rl.allow(ctx, keyPrefix+"ip:"+ip, rule.Requests, rule.Window)
	if err != nil {
		return false
	}
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
	if !allowed {
		errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
		c.Abort()
		return true
	}
	c.Next()
	return true
}

// rateLimiter stores rate limiters per IP address
type rateLimiter struct {
	limiter        *rate.Limiter
	lastSeen       time.Time
	firstLimitTime *time.Time // Time when rate limit was first exceeded (nil if not limited)
	window         int        // Window in seconds for calculating reset time
}

// rateLimiters is a map of IP addresses to their rate limiters
type rateLimiters struct {
	mu       sync.RWMutex
	limiters map[string]*rateLimiter
	cleanup  *time.Ticker
}

var (
	// Redis-based rate limiter
	redisRateLimiter *ratelimit.RedisRateLimiter
	useRedis         bool

	// Global in-memory rate limiters for different endpoint types (fallback)
	loginLimiters        *rateLimiters // Level 1: IP-based
	loginByEmailLimiters *rateLimiters // Level 2: Email-based
	loginGlobalLimiters  *rateLimiters // Level 3: Global limit
	refreshLimiters      *rateLimiters
	uploadLimiters       *rateLimiters
	generalLimiters      *rateLimiters
	publicLimiters       *rateLimiters
	mutationLimiters     *rateLimiters
	highVolumeLimiters   *rateLimiters
)

// InitRateLimiter initializes the rate limiter with optional Redis support.
// If redisClient is nil, falls back to in-memory rate limiting.
func InitRateLimiter(redisClient redis.UniversalClient) {
	if redisClient != nil {
		redisRateLimiter = ratelimit.NewRedisRateLimiter(redisClient, "ratelimit:")
		useRedis = true
		log.Println("[RateLimit] Using Redis-based distributed rate limiting")
	} else {
		useRedis = false
		log.Println("[RateLimit] Using in-memory rate limiting (single-instance only)")
	}
}

func init() {
	// Initialize rate limiters
	loginLimiters = newRateLimiters()        // Level 1: IP-based
	loginByEmailLimiters = newRateLimiters() // Level 2: Email-based
	loginGlobalLimiters = newRateLimiters()  // Level 3: Global (single key "global")
	refreshLimiters = newRateLimiters()
	uploadLimiters = newRateLimiters()
	generalLimiters = newRateLimiters()
	publicLimiters = newRateLimiters()
	mutationLimiters = newRateLimiters()
	highVolumeLimiters = newRateLimiters()
}

// newRateLimiters creates a new rate limiters instance
func newRateLimiters() *rateLimiters {
	rl := &rateLimiters{
		limiters: make(map[string]*rateLimiter),
		cleanup:  time.NewTicker(5 * time.Minute), // Cleanup every 5 minutes
	}

	// Start cleanup goroutine
	go func() {
		for range rl.cleanup.C {
			rl.cleanupLimiters()
		}
	}()

	return rl
}

// getLimiter returns a rate limiter for the given key (IP, email, etc.)
// It also returns the rateLimiter struct to access firstLimitTime
func (rl *rateLimiters) getLimiter(key string, requests int, window int) (*rate.Limiter, *rateLimiter) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		// Create new limiter: rate.Every(window) creates a limiter that allows
		// requests at a rate of 1 request per window/requests seconds
		// So for 5 requests per 15 minutes: rate.Every(180 seconds) with burst of 5
		interval := time.Duration(window) * time.Second / time.Duration(requests)
		limiter = &rateLimiter{
			limiter:        rate.NewLimiter(rate.Every(interval), requests),
			lastSeen:       time.Now(),
			firstLimitTime: nil,
			window:         window,
		}
		rl.limiters[key] = limiter
	} else {
		limiter.lastSeen = time.Now()
		// Update window if changed (shouldn't happen, but just in case)
		limiter.window = window
	}

	return limiter.limiter, limiter
}

// getResetTime calculates the reset time for a rate limiter
// If firstLimitTime is set, use it. Otherwise, calculate from now.
func (rl *rateLimiters) getResetTime(limiter *rateLimiter, window int) int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if limiter.firstLimitTime != nil {
		// Use the time when limit was first exceeded
		resetTime := limiter.firstLimitTime.Add(time.Duration(window) * time.Second)
		return resetTime.Unix()
	}
	// Fallback: calculate from now (shouldn't happen if logic is correct)
	return time.Now().Add(time.Duration(window) * time.Second).Unix()
}

// setFirstLimitTime records when rate limit was first exceeded
func (rl *rateLimiters) setFirstLimitTime(limiter *rateLimiter) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter.firstLimitTime == nil {
		now := time.Now()
		limiter.firstLimitTime = &now
	}
}

// clearFirstLimitTime clears the first limit time when limit resets
func (rl *rateLimiters) clearFirstLimitTime(limiter *rateLimiter) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter.firstLimitTime = nil
}

// cleanupLimiters removes limiters that haven't been used in the last hour
func (rl *rateLimiters) cleanupLimiters() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, limiter := range rl.limiters {
		if now.Sub(limiter.lastSeen) > 1*time.Hour {
			delete(rl.limiters, ip)
		}
	}
}

// isLocalhost checks if the IP is localhost/127.0.0.1/::1
func isLocalhost(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1" || ip == "localhost" || ip == "" || strings.HasSuffix(ip, "127.0.0.1")
}

// checkRateLimitRedis checks rate limit using Redis if available, falls back to in-memory.
// Returns: (allowed bool, remaining int, resetTime int64)
func checkRateLimitRedis(c *gin.Context, key string, limit int, window int, fallbackLimiter *rateLimiters) (bool, int, int64) {
	// Try Redis first if enabled
	if useRedis && redisRateLimiter != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		allowed, remaining, resetTime, err := redisRateLimiter.Allow(
			ctx,
			key,
			limit,
			time.Duration(window)*time.Second,
		)

		if err == nil {
			// Redis success - return results
			return allowed, remaining, resetTime.Unix()
		}

		// Redis failed - log warning and fall back to in-memory
		log.Printf("[RateLimit] Redis error for key=%s, falling back to in-memory: %v", key, err)
	}

	// In-memory fallback
	limiter, limiterStruct := fallbackLimiter.getLimiter(key, limit, window)

	if !limiter.Allow() {
		// Rate limit exceeded
		fallbackLimiter.setFirstLimitTime(limiterStruct)
		resetTime := fallbackLimiter.getResetTime(limiterStruct, window)
		return false, 0, resetTime
	}

	// Request allowed
	fallbackLimiter.clearFirstLimitTime(limiterStruct)
	remaining := int(limiter.Tokens())
	if remaining < 0 {
		remaining = 0
	}
	resetTime := fallbackLimiter.getResetTime(limiterStruct, window)
	return true, remaining, resetTime
}


// rateLimitKeyForType resolves the rate-limit bucket key.
// For authenticated endpoint types (general, mutation, high_volume, upload, refresh)
// we key by user ID so that each user owns their individual quota rather than
// sharing a single IP bucket — critical when 1 000+ users sit behind the same
// corporate NAT or load-balancer.
// Public and login endpoints keep IP-based keying (pre-auth identity).
func rateLimitKeyForType(c *gin.Context, limitType string, ip string) string {
	switch limitType {
	case "login", "public":
		// Pre-authentication: must use IP as no user identity is available yet.
		return ip
	default:
		// Post-authentication: prefer per-user key to avoid shared-IP collisions.
		if raw, exists := c.Get("user_id"); exists {
			if uid, ok := raw.(string); ok && uid != "" {
				return "user:" + uid
			}
		}
		// Fallback to IP when JWT context is absent (should not happen on
		// authenticated routes, but keeps the function safe).
		return ip
	}
}

// RateLimitMiddleware creates a rate limiting middleware based on endpoint type
func RateLimitMiddleware(limitType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Skip rate limiting for localhost in development (to avoid issues during testing)
		// In production, this should be removed or controlled via environment variable
		if isLocalhost(ip) && (config.AppConfig.Server.Env == "development" || config.AppConfig.Server.Env == "dev") {
			c.Next()
			return
		}

		// Resolve the rate-limit key: per-user for authenticated routes, per-IP otherwise.
		rateLimitKey := rateLimitKeyForType(c, limitType, ip)

		// Prefer Redis-based rate limiting for multi-instance consistency.
		// If Redis is disabled/unavailable/errors, fall back to in-memory limiter.
		if handled := tryRedisRateLimit(c, limitType, rateLimitKey); handled {
			return
		}

		var limiter *rate.Limiter
		var limiterStruct *rateLimiter
		var rule config.RateLimitRule

		// For login endpoint, implement multi-level rate limiting
		if limitType == "login" {
			// Level 3: Global rate limit (check first to prevent DOS)
			globalRule := config.AppConfig.RateLimit.LoginGlobal
			globalLimiter, globalLimiterStruct := loginGlobalLimiters.getLimiter("global", globalRule.Requests, globalRule.Window)
			if !globalLimiter.Allow() {
				// Record first limit time if not already set
				loginGlobalLimiters.setFirstLimitTime(globalLimiterStruct)
				// Get reset time based on first limit time (stable, won't change on refresh)
				resetTime := loginGlobalLimiters.getResetTime(globalLimiterStruct, globalRule.Window)
				c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", globalRule.Requests))
				c.Header("X-RateLimit-Remaining", "0")
				c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
				c.Abort()
				return
			} else {
				// Clear first limit time if limit is no longer exceeded
				loginGlobalLimiters.clearFirstLimitTime(globalLimiterStruct)
			}

			// Level 2: Rate limit by email/username (extract from request body)
			// Read request body to get email without consuming it
			var loginReq struct {
				Email string `json:"email"`
			}

			// Read body and restore it so handler can read it too
			if c.Request.Body != nil {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil && len(bodyBytes) > 0 {
					// Restore body for handler
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					// Try to parse JSON to get email
					if json.Unmarshal(bodyBytes, &loginReq) == nil && loginReq.Email != "" {
						emailRule := config.AppConfig.RateLimit.LoginByEmail
						emailLimiter, emailLimiterStruct := loginByEmailLimiters.getLimiter(loginReq.Email, emailRule.Requests, emailRule.Window)
						if !emailLimiter.Allow() {
							// Record first limit time if not already set
							loginByEmailLimiters.setFirstLimitTime(emailLimiterStruct)
							// Get reset time based on first limit time (stable, won't change on refresh)
							resetTime := loginByEmailLimiters.getResetTime(emailLimiterStruct, emailRule.Window)
							c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", emailRule.Requests))
							c.Header("X-RateLimit-Remaining", "0")
							c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
							errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
							c.Abort()
							return
						} else {
							// Clear first limit time if limit is no longer exceeded
							loginByEmailLimiters.clearFirstLimitTime(emailLimiterStruct)
						}
					}
				}
			}

			// Level 1: IP-based rate limit
			rule = config.AppConfig.RateLimit.Login
			limiter, limiterStruct = loginLimiters.getLimiter(ip, rule.Requests, rule.Window)
		} else {
			// Get appropriate rate limit rule based on type for non-login endpoints.
			// Use rateLimitKey (user ID for authenticated types) so each user gets
			// their own in-memory bucket, not a shared IP bucket.
			var limiterStruct *rateLimiter
			switch limitType {
			case "refresh":
				rule = config.AppConfig.RateLimit.Refresh
				limiter, limiterStruct = refreshLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			case "upload":
				rule = config.AppConfig.RateLimit.Upload
				limiter, limiterStruct = uploadLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			case "general":
				rule = config.AppConfig.RateLimit.General
				limiter, limiterStruct = generalLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			case "public":
				rule = config.AppConfig.RateLimit.Public
				limiter, limiterStruct = publicLimiters.getLimiter(ip, rule.Requests, rule.Window)
			case "mutation":
				rule = config.AppConfig.RateLimit.Mutation
				limiter, limiterStruct = mutationLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			case "high_volume":
				rule = config.AppConfig.RateLimit.HighVolume
				limiter, limiterStruct = highVolumeLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			default:
				rule = config.AppConfig.RateLimit.General
				limiter, limiterStruct = generalLimiters.getLimiter(rateLimitKey, rule.Requests, rule.Window)
			}

			// Check if request is allowed
			if !limiter.Allow() {
				// Record first limit time if not already set
				switch limitType {
				case "refresh":
					refreshLimiters.setFirstLimitTime(limiterStruct)
					resetTime := refreshLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				case "upload":
					uploadLimiters.setFirstLimitTime(limiterStruct)
					resetTime := uploadLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				case "general":
					generalLimiters.setFirstLimitTime(limiterStruct)
					resetTime := generalLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				case "public":
					publicLimiters.setFirstLimitTime(limiterStruct)
					resetTime := publicLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				case "mutation":
					mutationLimiters.setFirstLimitTime(limiterStruct)
					resetTime := mutationLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				case "high_volume":
					highVolumeLimiters.setFirstLimitTime(limiterStruct)
					resetTime := highVolumeLimiters.getResetTime(limiterStruct, rule.Window)
					c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
					c.Header("X-RateLimit-Remaining", "0")
					c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				}

				// Return 429 Too Many Requests with proper error format
				errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
				c.Abort()
				return
			} else {
				// Clear first limit time if limit is no longer exceeded
				switch limitType {
				case "refresh":
					refreshLimiters.clearFirstLimitTime(limiterStruct)
				case "upload":
					uploadLimiters.clearFirstLimitTime(limiterStruct)
				case "general":
					generalLimiters.clearFirstLimitTime(limiterStruct)
				case "public":
					publicLimiters.clearFirstLimitTime(limiterStruct)
				case "mutation":
					mutationLimiters.clearFirstLimitTime(limiterStruct)
				case "high_volume":
					highVolumeLimiters.clearFirstLimitTime(limiterStruct)
				}
			}
		}

		// Check if request is allowed (for Level 1 login)
		if limitType == "login" {
			if !limiter.Allow() {
				// Record first limit time if not already set
				loginLimiters.setFirstLimitTime(limiterStruct)
				// Get reset time based on first limit time (stable, won't change on refresh)
				resetTime := loginLimiters.getResetTime(limiterStruct, rule.Window)
				c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
				c.Header("X-RateLimit-Remaining", "0")
				c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
				errors.ErrorResponse(c, "RATE_LIMIT_EXCEEDED", nil, nil)
				c.Abort()
				return
			} else {
				// Clear first limit time if limit is no longer exceeded
				loginLimiters.clearFirstLimitTime(limiterStruct)
			}
		}

		// Calculate remaining requests
		// Note: rate.Limiter doesn't expose remaining directly, so we estimate
		// based on the burst capacity
		remaining := rule.Requests - 1 // Approximate, as limiter doesn't expose exact remaining

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rule.Requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(rule.Window)*time.Second).Unix()))

		c.Next()
	}
}
