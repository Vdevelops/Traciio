package database

import (
	"fmt"
	"log"
)

// PoolStats represents database connection pool statistics
type PoolStats struct {
	MaxOpenConnections int           // Maximum number of open connections
	OpenConnections    int           // Current number of open connections
	InUse              int           // Number of connections currently in use
	Idle               int           // Number of idle connections
	WaitCount          int64         // Total number of connections waited for
	WaitDuration       int64         // Total time blocked waiting for a new connection (nanoseconds)
	MaxIdleClosed     int64         // Total number of connections closed due to SetMaxIdleConns
	MaxLifetimeClosed  int64         // Total number of connections closed due to SetConnMaxLifetime
}

// GetPoolStats returns current connection pool statistics
// This is useful for monitoring and alerting on pool saturation
func GetPoolStats() (*PoolStats, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()

	return &PoolStats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:     stats.OpenConnections,
		InUse:              stats.InUse,
		Idle:               stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration.Nanoseconds(),
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}, nil
}

// LogPoolStats logs connection pool statistics
// Useful for debugging connection pool issues
func LogPoolStats() {
	stats, err := GetPoolStats()
	if err != nil {
		log.Printf("Error getting pool stats: %v", err)
		return
	}

	utilizationPercent := float64(0)
	if stats.MaxOpenConnections > 0 {
		utilizationPercent = float64(stats.OpenConnections) / float64(stats.MaxOpenConnections) * 100
	}

	log.Printf("📊 Database Connection Pool Stats:")
	log.Printf("   Max Open: %d", stats.MaxOpenConnections)
	log.Printf("   Open: %d (%.1f%% utilization)", stats.OpenConnections, utilizationPercent)
	log.Printf("   In Use: %d", stats.InUse)
	log.Printf("   Idle: %d", stats.Idle)
	log.Printf("   Wait Count: %d", stats.WaitCount)
	log.Printf("   Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("   Max Lifetime Closed: %d", stats.MaxLifetimeClosed)

	// Alert if pool is getting saturated (>80%)
	if utilizationPercent > 80 {
		log.Printf("⚠️  WARNING: Connection pool utilization is %.1f%% - consider increasing MaxOpenConns", utilizationPercent)
	}
}

// CheckPoolHealth checks if connection pool is healthy
// Returns true if pool is healthy, false otherwise
func CheckPoolHealth() (bool, string) {
	stats, err := GetPoolStats()
	if err != nil {
		return false, fmt.Sprintf("Failed to get pool stats: %v", err)
	}

	// Check if pool is saturated (>90% utilization)
	utilizationPercent := float64(0)
	if stats.MaxOpenConnections > 0 {
		utilizationPercent = float64(stats.OpenConnections) / float64(stats.MaxOpenConnections) * 100
	}

	if utilizationPercent > 90 {
		return false, fmt.Sprintf("Connection pool is saturated: %.1f%% utilization", utilizationPercent)
	}

	// Check if there are many waits (indicates pool might be too small)
	if stats.WaitCount > 100 {
		return false, fmt.Sprintf("High wait count: %d - consider increasing MaxOpenConns", stats.WaitCount)
	}

	return true, "Pool is healthy"
}

