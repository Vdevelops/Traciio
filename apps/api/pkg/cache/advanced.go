// Package cache provides advanced Redis caching layer with high availability support.
//
// This file extends the base cache functionality with advanced patterns for:
// - Consistent cache key management
// - Batch operations for high throughput
// - Cache warming strategies
// - Graceful degradation patterns
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// Enterprise cache key prefixes following naming convention: resource:action:params
const (
	// Dashboard module
	PrefixDashboardOverview   = "dashboard:overview:"
	PrefixDashboardVisitStats = "dashboard:visit_stats:"
	PrefixDashboardSales      = "dashboard:sales:"


	// Account module
	PrefixAccountList   = "account:list:"
	PrefixAccountDetail = "account:detail:"
	PrefixAccountStats  = "account:stats:"

	// Lead module
	PrefixLeadList   = "lead:list:"
	PrefixLeadDetail = "lead:detail:"
	PrefixLeadStats  = "lead:stats:"

	// Contact module
	PrefixContactList   = "contact:list:"
	PrefixContactDetail = "contact:detail:"

	// Activity module
	PrefixActivityList   = "activity:list:"
	PrefixActivityDetail = "activity:detail:"
	PrefixActivityStats  = "activity:stats:"

	// Task module
	PrefixTaskList   = "task:list:"
	PrefixTaskDetail = "task:detail:"

	// Visit Report module
	PrefixVisitReportList   = "visit_report:list:"
	PrefixVisitReportDetail = "visit_report:detail:"
	PrefixVisitReportStats  = "visit_report:stats:"

	// Brick module
	PrefixBrickList    = "brick:list:"
	PrefixBrickDetail  = "brick:detail:"
	PrefixBrickGeoJSON = "brick:geojson:"

	// Schedule module
	PrefixScheduleList   = "schedule:list:"
	PrefixScheduleDetail = "schedule:detail:"

	// Product module (existing)
	PrefixProductList   = "product:list:"
	PrefixProductDetail = "product:detail:"

	// Pipeline/Deal module
	PrefixDealList    = "deal:list:"
	PrefixDealDetail  = "deal:detail:"
	PrefixDealSummary = "deal:summary:"
)

// TTL configurations following enterprise best practice standards
// Based on: Data volatility, access frequency, and tolerance for stale data
const (
	// ============================================================
	// DASHBOARD: 30-60 seconds (high volatility, real-time feel)
	// Used for: KPI, charts, summary cards
	// ============================================================
	TTLDashboardShort = 30 * time.Second  // Real-time dashboards
	TTLDashboardLong  = 60 * time.Second  // Summary dashboards

	// ============================================================
	// LEADS/TRANSACTIONAL: 15-30 seconds (very short - active sales data)
	// Data sering berubah, user expect near real-time
	// ============================================================
	TTLLeadShort  = 15 * time.Second // Lead list yang aktif
	TTLLeadMedium = 30 * time.Second // Lead detail

	// ============================================================
	// ACCOUNTS/CONTACTS: 1-5 minutes (moderate volatility)
	// Master data yang jarang berubah dalam hitungan detik
	// ============================================================
	TTLListShort  = 1 * time.Minute  // Frequently accessed lists
	TTLListMedium = 3 * time.Minute  // Standard list caching
	TTLListLong   = 5 * time.Minute  // Less volatile lists

	// ============================================================
	// DETAIL RECORDS: 3-5 minutes
	// Individual record access, balance freshness vs performance
	// ============================================================
	TTLDetailShort  = 3 * time.Minute
	TTLDetailMedium = 5 * time.Minute
	TTLDetailLong   = 10 * time.Minute

	// ============================================================
	// RBAC/PERMISSIONS: 10-30 minutes (CRITICAL - every request)
	// Role->Permissions, User->Roles mapping
	// Invalidate on role update
	// ============================================================
	TTLRBACShort  = 10 * time.Minute // User->Role mapping
	TTLRBACMedium = 20 * time.Minute // Role->Permissions
	TTLRBACLong   = 30 * time.Minute // Permission list (rarely changes)

	// ============================================================
	// REFERENCE/MASTER DATA: 5-30 minutes (low volatility)
	// Products, categories, divisions, etc.
	// ============================================================
	TTLReferenceShort  = 5 * time.Minute  // Frequently updated reference
	TTLReferenceMedium = 15 * time.Minute // Standard reference data
	TTLReferenceData   = 30 * time.Minute // Static reference data

	// ============================================================
	// ANALYTICS/REPORTS: 1-10 minutes (aggregated data)
	// Pipeline summary, sales analytics, performance reports
	// ============================================================
	TTLAnalyticsShort  = 1 * time.Minute  // Real-time analytics
	TTLAnalyticsMedium = 5 * time.Minute  // Standard analytics
	TTLAnalyticsLong   = 10 * time.Minute // Historical analytics

	// Report data aliases (backward compatibility)
	TTLReportShort  = 5 * time.Minute
	TTLReportMedium = 10 * time.Minute
	TTLReportLong   = 15 * time.Minute

	// ============================================================
	// STATISTICS/AGGREGATIONS: 1-3 minutes
	// Count, sum, avg operations - balance accuracy vs performance
	// ============================================================
	TTLStatsShort  = 1 * time.Minute
	TTLStatsMedium = 2 * time.Minute
	TTLStatsLong   = 3 * time.Minute

	// ============================================================
	// GEOJSON/SPATIAL DATA: 30 minutes (large, expensive to compute)
	// Brick boundaries, map data
	// ============================================================
	TTLGeoJSON = 30 * time.Minute
)

// CacheKeyBuilder provides consistent cache key generation
type CacheKeyBuilder struct {
	prefix string
	parts  []string
}

// NewCacheKeyBuilder creates a new cache key builder with the given prefix
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{
		prefix: prefix,
		parts:  make([]string, 0),
	}
}

// WithParam adds a parameter to the cache key
func (b *CacheKeyBuilder) WithParam(key string, value interface{}) *CacheKeyBuilder {
	b.parts = append(b.parts, fmt.Sprintf("%s=%v", key, value))
	return b
}

// WithUser adds user-specific context to the cache key
func (b *CacheKeyBuilder) WithUser(userID string) *CacheKeyBuilder {
	b.parts = append(b.parts, fmt.Sprintf("user=%s", userID))
	return b
}

// WithPagination adds pagination parameters to the cache key
func (b *CacheKeyBuilder) WithPagination(page, perPage int) *CacheKeyBuilder {
	b.parts = append(b.parts, fmt.Sprintf("page=%d", page))
	b.parts = append(b.parts, fmt.Sprintf("per_page=%d", perPage))
	return b
}

// WithPeriod adds period/date range to the cache key
func (b *CacheKeyBuilder) WithPeriod(period, startDate, endDate string) *CacheKeyBuilder {
	if period != "" {
		b.parts = append(b.parts, fmt.Sprintf("period=%s", period))
	}
	if startDate != "" {
		b.parts = append(b.parts, fmt.Sprintf("start=%s", startDate))
	}
	if endDate != "" {
		b.parts = append(b.parts, fmt.Sprintf("end=%s", endDate))
	}
	return b
}

// Build generates the final cache key
func (b *CacheKeyBuilder) Build() string {
	if len(b.parts) == 0 {
		return b.prefix
	}
	return b.prefix + strings.Join(b.parts, ":")
}

// BuildHashed generates a hashed cache key for very long parameter combinations
func (b *CacheKeyBuilder) BuildHashed() string {
	fullKey := b.Build()
	if len(fullKey) > 200 {
		hash := sha256.Sum256([]byte(fullKey))
		return b.prefix + hex.EncodeToString(hash[:16])
	}
	return fullKey
}

// CacheResult wraps cached data with metadata
type CacheResult struct {
	Data      interface{} `msgpack:"d"`
	CachedAt  time.Time   `msgpack:"c"`
	ExpiresAt time.Time   `msgpack:"e"`
	Version   int         `msgpack:"v"`
}

// AdvancedCache provides advanced caching patterns and optimizations
type AdvancedCache struct {
	cache        *Cache
	warmupQueue  chan warmupTask
	warmupWG     sync.WaitGroup
	stopWarmup   chan struct{}
	metricsExt   *AdvancedMetrics
}

type warmupTask struct {
	key    string
	loader func() (interface{}, error)
	ttl    time.Duration
}

// AdvancedMetrics extends basic metrics with advanced feature tracking
type AdvancedMetrics struct {
	WarmupTasks    uint64
	WarmupSuccess  uint64
	WarmupFailures uint64
	BatchOps       uint64
	CacheAsideHit  uint64
	CacheAsideMiss uint64
	Invalidations  uint64
	PatternDeletes uint64
}

var (
	advancedCache     *AdvancedCache
	advancedCacheOnce sync.Once
)

// Advanced returns the singleton advanced cache instance
func Advanced() *AdvancedCache {
	advancedCacheOnce.Do(func() {
		if Client == nil {
			advancedCache = &AdvancedCache{
				metricsExt: &AdvancedMetrics{},
			}
			return
		}
		
		advancedCache = &AdvancedCache{
			cache:       Client,
			warmupQueue: make(chan warmupTask, 1000),
			stopWarmup:  make(chan struct{}),
			metricsExt:  &AdvancedMetrics{},
		}

		// Start warmup workers
		for i := 0; i < 5; i++ {
			advancedCache.warmupWG.Add(1)
			go advancedCache.warmupWorker()
		}
	})
	return advancedCache
}

// warmupWorker processes cache warmup tasks
func (ac *AdvancedCache) warmupWorker() {
	defer ac.warmupWG.Done()
	for {
		select {
		case task := <-ac.warmupQueue:
			atomic.AddUint64(&ac.metricsExt.WarmupTasks, 1)
			data, err := task.loader()
			if err != nil {
				atomic.AddUint64(&ac.metricsExt.WarmupFailures, 1)
				continue
			}
			if err := ac.Set(task.key, data, task.ttl); err == nil {
				atomic.AddUint64(&ac.metricsExt.WarmupSuccess, 1)
			} else {
				atomic.AddUint64(&ac.metricsExt.WarmupFailures, 1)
			}
		case <-ac.stopWarmup:
			return
		}
	}
}

// IsEnabled checks if advanced cache is active
func (ac *AdvancedCache) IsEnabled() bool {
	return ac != nil && ac.cache != nil && ac.cache.IsEnabled()
}

// Get retrieves a value from cache with advanced patterns
func (ac *AdvancedCache) Get(key string, target interface{}) (bool, error) {
	if !ac.IsEnabled() {
		return false, nil
	}
	return ac.cache.Get(key, target)
}

// Set stores a value in cache with advanced patterns
func (ac *AdvancedCache) Set(key string, value interface{}, ttl time.Duration) error {
	if !ac.IsEnabled() {
		return nil
	}
	return ac.cache.Set(key, value, ttl)
}

// Delete removes a key from cache
func (ac *AdvancedCache) Delete(key string) error {
	if !ac.IsEnabled() {
		return nil
	}
	atomic.AddUint64(&ac.metricsExt.Invalidations, 1)
	return ac.cache.Delete(key)
}

// DeletePattern removes all keys matching the pattern
func (ac *AdvancedCache) DeletePattern(pattern string) error {
	if !ac.IsEnabled() {
		return nil
	}
	atomic.AddUint64(&ac.metricsExt.PatternDeletes, 1)
	return ac.cache.DeletePattern(pattern)
}

// GetOrLoad implements cache-aside pattern with atomic operations
// This is the primary method for enterprise caching: try cache first, load if miss
func (ac *AdvancedCache) GetOrLoad(key string, target interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	if !ac.IsEnabled() {
		// Cache disabled, just load from source
		data, err := loader()
		if err != nil {
			return err
		}
		return copyData(data, target)
	}

	// Try to get from cache first
	found, err := ac.Get(key, target)
	if err == nil && found {
		atomic.AddUint64(&ac.metricsExt.CacheAsideHit, 1)
		return nil
	}

	// Cache miss - load from source
	atomic.AddUint64(&ac.metricsExt.CacheAsideMiss, 1)
	data, err := loader()
	if err != nil {
		return err
	}

	// Store in cache (async, don't block on cache write)
	go func() {
		_ = ac.Set(key, data, ttl)
	}()

	return copyData(data, target)
}

// MGet retrieves multiple keys in a single round trip
func (ac *AdvancedCache) MGet(keys []string, targets []interface{}) ([]bool, error) {
	if !ac.IsEnabled() || len(keys) == 0 {
		return make([]bool, len(keys)), nil
	}
	atomic.AddUint64(&ac.metricsExt.BatchOps, 1)
	return ac.cache.MGet(keys, targets)
}

// MSet stores multiple key-value pairs in a single round trip
func (ac *AdvancedCache) MSet(items map[string]interface{}, ttl time.Duration) error {
	if !ac.IsEnabled() || len(items) == 0 {
		return nil
	}
	atomic.AddUint64(&ac.metricsExt.BatchOps, 1)
	return ac.cache.MSet(items, ttl)
}

// MDelete removes multiple keys in a single round trip
func (ac *AdvancedCache) MDelete(keys []string) error {
	if !ac.IsEnabled() || len(keys) == 0 {
		return nil
	}
	atomic.AddUint64(&ac.metricsExt.BatchOps, 1)
	atomic.AddUint64(&ac.metricsExt.Invalidations, uint64(len(keys)))
	return ac.cache.MDelete(keys)
}

// ScheduleWarmup queues a cache warmup task (non-blocking)
func (ac *AdvancedCache) ScheduleWarmup(key string, loader func() (interface{}, error), ttl time.Duration) {
	if !ac.IsEnabled() {
		return
	}
	select {
	case ac.warmupQueue <- warmupTask{key: key, loader: loader, ttl: ttl}:
	default:
		// Queue full, skip warmup
	}
}

// InvalidateByPrefix invalidates all cache entries with the given prefix
func (ac *AdvancedCache) InvalidateByPrefix(prefix string) error {
	return ac.DeletePattern(prefix + "*")
}

// InvalidateModule invalidates all cache entries for a module
func (ac *AdvancedCache) InvalidateModule(module string) error {
	return ac.DeletePattern(module + ":*")
}

// GetAdvancedMetrics returns advanced feature metrics
func (ac *AdvancedCache) GetAdvancedMetrics() *AdvancedMetrics {
	if ac.metricsExt == nil {
		return &AdvancedMetrics{}
	}
	return &AdvancedMetrics{
		WarmupTasks:    atomic.LoadUint64(&ac.metricsExt.WarmupTasks),
		WarmupSuccess:  atomic.LoadUint64(&ac.metricsExt.WarmupSuccess),
		WarmupFailures: atomic.LoadUint64(&ac.metricsExt.WarmupFailures),
		BatchOps:       atomic.LoadUint64(&ac.metricsExt.BatchOps),
		CacheAsideHit:  atomic.LoadUint64(&ac.metricsExt.CacheAsideHit),
		CacheAsideMiss: atomic.LoadUint64(&ac.metricsExt.CacheAsideMiss),
		Invalidations:  atomic.LoadUint64(&ac.metricsExt.Invalidations),
		PatternDeletes: atomic.LoadUint64(&ac.metricsExt.PatternDeletes),
	}
}

// Close gracefully shuts down the advanced cache
func (ac *AdvancedCache) Close() {
	if ac.stopWarmup != nil {
		close(ac.stopWarmup)
		ac.warmupWG.Wait()
	}
}

// copyData performs a deep copy using msgpack serialization
func copyData(src, dst interface{}) error {
	data, err := msgpack.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}
	return msgpack.Unmarshal(data, dst)
}

// CacheAside is a helper function implementing cache-aside pattern
// Usage: cache.CacheAside(key, &result, ttl, loaderFunc)
func CacheAside(key string, target interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	return Advanced().GetOrLoad(key, target, ttl, loader)
}

// InvalidateByEntity invalidates all cache entries related to an entity
// This is useful after CREATE, UPDATE, or DELETE operations
func InvalidateByEntity(entityType, entityID string) error {
	ac := Advanced()
	if !ac.IsEnabled() {
		return nil
	}
	
	// Common invalidation patterns
	patterns := []string{
		fmt.Sprintf("%s:list:*", entityType),
		fmt.Sprintf("%s:detail:%s*", entityType, entityID),
		fmt.Sprintf("%s:stats:*", entityType),
	}
	
	for _, pattern := range patterns {
		if err := ac.DeletePattern(pattern); err != nil {
			return err
		}
	}
	
	return nil
}

// InvalidateRelatedEntities invalidates cache for related entities
// Use this when changes to one entity affect others
func InvalidateRelatedEntities(entities ...string) error {
	ac := Advanced()
	if !ac.IsEnabled() {
		return nil
	}
	
	for _, entity := range entities {
		if err := ac.InvalidateModule(entity); err != nil {
			return err
		}
	}
	
	return nil
}

// WithContext creates a context-aware cache operation
func (ac *AdvancedCache) WithContext(ctx context.Context) *ContextCache {
	return &ContextCache{
		ac:  ac,
		ctx: ctx,
	}
}

// ContextCache provides context-aware cache operations
type ContextCache struct {
	ac  *AdvancedCache
	ctx context.Context
}

// Get retrieves a value with context awareness
func (cc *ContextCache) Get(key string, target interface{}) (bool, error) {
	select {
	case <-cc.ctx.Done():
		return false, cc.ctx.Err()
	default:
		return cc.ac.Get(key, target)
	}
}

// Set stores a value with context awareness
func (cc *ContextCache) Set(key string, value interface{}, ttl time.Duration) error {
	select {
	case <-cc.ctx.Done():
		return cc.ctx.Err()
	default:
		return cc.ac.Set(key, value, ttl)
	}
}
