package database

import (
	"context"
	"time"
)

// DefaultQueryTimeout is the default timeout for database queries
const DefaultQueryTimeout = 30 * time.Second

// WithTimeout creates a context with timeout for database queries
// This prevents queries from hanging indefinitely
func WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = DefaultQueryTimeout
	}
	return context.WithTimeout(context.Background(), timeout)
}

// WithTimeoutContext creates a context with timeout from parent context
func WithTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = DefaultQueryTimeout
	}
	return context.WithTimeout(parent, timeout)
}

// QueryWithTimeout executes a query function with timeout
// Usage:
//   err := QueryWithTimeout(30*time.Second, func(ctx context.Context) error {
//       return DB.WithContext(ctx).Where(...).Find(&results).Error
//   })
func QueryWithTimeout(timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := WithTimeout(timeout)
	defer cancel()
	return fn(ctx)
}

