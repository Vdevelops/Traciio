package repository

import (
	"context"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"gorm.io/gorm"
)

// BaseRepository provides common repository functionality
// All repositories should embed this for consistent behavior
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// WithTimeout executes a query function with timeout
// This prevents queries from hanging indefinitely
func (r *BaseRepository) WithTimeout(timeout time.Duration, fn func(ctx context.Context, db *gorm.DB) error) error {
	return database.QueryWithTimeout(timeout, func(ctx context.Context) error {
		return fn(ctx, r.db.WithContext(ctx))
	})
}

// DB returns the database connection with context timeout
// Usage: r.DB(ctx).Where(...).Find(...)
func (r *BaseRepository) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		ctx = context.Background()
	}
	return r.db.WithContext(ctx)
}

// DBWithTimeout returns the database connection with timeout context
// Usage: r.DBWithTimeout(30*time.Second).Where(...).Find(...)
func (r *BaseRepository) DBWithTimeout(timeout time.Duration) *gorm.DB {
	ctx, _ := database.WithTimeout(timeout)
	return r.db.WithContext(ctx)
}

// DefaultTimeout is the default timeout for database queries
const DefaultTimeout = 30 * time.Second

