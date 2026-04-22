// Package ratelimit provides Redis-based distributed rate limiting with automatic fallback.
//
// This package implements a production-ready distributed rate limiter that uses Redis
// for shared state across multiple API instances. It includes automatic fallback to
// in-memory rate limiting when Redis is unavailable.
package ratelimit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements distributed rate limiting using Redis.
// It uses a fixed window counter algorithm with Lua scripts for atomic operations.
type RedisRateLimiter struct {
	client redis.UniversalClient
	prefix string
}

// NewRedisRateLimiter creates a new Redis-based rate limiter.
// The prefix is prepended to all Redis keys (e.g., "ratelimit:").
func NewRedisRateLimiter(client redis.UniversalClient, prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		prefix: prefix,
	}
}

// Allow checks if a request is allowed under the specified rate limit.
// It uses a Lua script to atomically increment the counter and set expiry.
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - key: Unique identifier (e.g., IP address, user ID, email)
//   - limit: Maximum number of requests allowed in the window
//   - window: Time window duration
//
// Returns:
//   - allowed: true if request is within limit
//   - remaining: number of requests remaining in window
//   - resetTime: time when the rate limit window resets
//   - error: any error that occurred
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, int, time.Time, error) {
	redisKey := r.prefix + key

	// Lua script for atomic increment and TTL management
	// This ensures that increment and expiry are set atomically without race conditions
	script := `
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		local ttl = redis.call('TTL', KEYS[1])
		return {count, ttl}
	`

	result, err := r.client.Eval(ctx, script, []string{redisKey}, int(window.Seconds())).Result()
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("redis eval failed: %w", err)
	}

	// Parse result from Lua script
	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return false, 0, time.Time{}, fmt.Errorf("unexpected redis response format")
	}

	count, ok := values[0].(int64)
	if !ok {
		return false, 0, time.Time{}, fmt.Errorf("invalid count type in redis response")
	}

	ttl, ok := values[1].(int64)
	if !ok {
		return false, 0, time.Time{}, fmt.Errorf("invalid ttl type in redis response")
	}

	// Calculate allowed, remaining, and reset time
	allowed := int(count) <= limit
	remaining := maxInt(0, limit-int(count))
	resetTime := time.Now().Add(time.Duration(ttl) * time.Second)

	// Debug logging in development
	if !allowed {
		log.Printf("[RateLimit] Limit exceeded for key=%s, count=%d, limit=%d, ttl=%ds", 
			key, count, limit, ttl)
	}

	return allowed, remaining, resetTime, nil
}

// Reset clears the rate limit counter for a specific key.
// This is useful for testing or administrative purposes.
func (r *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := r.prefix + key
	return r.client.Del(ctx, redisKey).Err()
}

// ResetAll clears all rate limit counters matching the prefix.
// Use with caution - this affects all users.
func (r *RedisRateLimiter) ResetAll(ctx context.Context) error {
	pattern := r.prefix + "*"
	
	// Use SCAN to find all matching keys (non-blocking)
	var cursor uint64
	var keys []string
	
	for {
		result, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		
		keys = append(keys, result...)
		cursor = nextCursor
		
		if cursor == 0 {
			break
		}
	}
	
	if len(keys) == 0 {
		return nil
	}
	
	// Delete in batches to avoid blocking
	const batchSize = 100
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		
		if err := r.client.Del(ctx, keys[i:end]...).Err(); err != nil {
			return fmt.Errorf("delete batch failed: %w", err)
		}
	}
	
	log.Printf("[RateLimit] Reset %d rate limit keys", len(keys))
	return nil
}

// GetCount returns the current count for a key without incrementing.
// Returns 0 if the key doesn't exist.
func (r *RedisRateLimiter) GetCount(ctx context.Context, key string) (int, error) {
	redisKey := r.prefix + key
	
	count, err := r.client.Get(ctx, redisKey).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get count failed: %w", err)
	}
	
	return count, nil
}

// maxInt returns the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
