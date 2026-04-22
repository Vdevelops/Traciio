// Package cache provides module-specific cache helpers for enterprise caching.
//
// This file contains helper functions for each module to ensure consistent
// caching patterns across the application.
package cache

import (
	"fmt"
	"sort"
	"time"
)

// ModuleCache provides module-specific cache operations
type ModuleCache struct {
	moduleName string
	ec         *AdvancedCache
}

// NewModuleCache creates a new module-specific cache helper
func NewModuleCache(moduleName string) *ModuleCache {
	return &ModuleCache{
		moduleName: moduleName,
		ec:         Advanced(),
	}
}

// ListKey generates a cache key for list operations
func (mc *ModuleCache) ListKey(page, perPage int, filters map[string]interface{}) string {
	builder := NewCacheKeyBuilder(mc.moduleName + ":list:")
	builder.WithPagination(page, perPage)

	if len(filters) > 0 {
		keys := make([]string, 0, len(filters))
		for k := range filters {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			builder.WithParam(k, filters[k])
		}
	}

	return builder.BuildHashed()
}

// DetailKey generates a cache key for detail/single record operations
func (mc *ModuleCache) DetailKey(id string) string {
	return fmt.Sprintf("%s:detail:%s", mc.moduleName, id)
}

// StatsKey generates a cache key for statistics operations
func (mc *ModuleCache) StatsKey(period, userID string, params map[string]interface{}) string {
	builder := NewCacheKeyBuilder(mc.moduleName + ":stats:")
	if period != "" {
		builder.WithParam("period", period)
	}
	if userID != "" {
		builder.WithUser(userID)
	}
	if len(params) > 0 {
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			builder.WithParam(k, params[k])
		}
	}
	return builder.BuildHashed()
}

// GetList retrieves list data from cache
func (mc *ModuleCache) GetList(key string, target interface{}) (bool, error) {
	return mc.ec.Get(key, target)
}

// SetList stores list data in cache
func (mc *ModuleCache) SetList(key string, data interface{}, ttl time.Duration) error {
	return mc.ec.Set(key, data, ttl)
}

// GetDetail retrieves detail data from cache
func (mc *ModuleCache) GetDetail(id string, target interface{}) (bool, error) {
	return mc.ec.Get(mc.DetailKey(id), target)
}

// SetDetail stores detail data in cache
func (mc *ModuleCache) SetDetail(id string, data interface{}, ttl time.Duration) error {
	return mc.ec.Set(mc.DetailKey(id), data, ttl)
}

// GetStats retrieves stats data from cache
func (mc *ModuleCache) GetStats(key string, target interface{}) (bool, error) {
	return mc.ec.Get(key, target)
}

// SetStats stores stats data in cache
func (mc *ModuleCache) SetStats(key string, data interface{}, ttl time.Duration) error {
	return mc.ec.Set(key, data, ttl)
}

// InvalidateList invalidates all list caches for this module
func (mc *ModuleCache) InvalidateList() error {
	return mc.ec.DeletePattern(mc.moduleName + ":list:*")
}

// InvalidateDetail invalidates a specific detail cache
func (mc *ModuleCache) InvalidateDetail(id string) error {
	return mc.ec.Delete(mc.DetailKey(id))
}

// InvalidateStats invalidates all stats caches for this module
func (mc *ModuleCache) InvalidateStats() error {
	return mc.ec.DeletePattern(mc.moduleName + ":stats:*")
}

// InvalidateAll invalidates all caches for this module
func (mc *ModuleCache) InvalidateAll() error {
	return mc.ec.InvalidateModule(mc.moduleName)
}

// IsEnabled checks if cache is enabled
func (mc *ModuleCache) IsEnabled() bool {
	return mc.ec.IsEnabled()
}

// DashboardCache provides dashboard-specific caching
type DashboardCache struct {
	*ModuleCache
}

// NewDashboardCache creates a new dashboard cache helper
func NewDashboardCache() *DashboardCache {
	return &DashboardCache{
		ModuleCache: NewModuleCache("dashboard"),
	}
}

// OverviewKey generates cache key for dashboard overview
func (dc *DashboardCache) OverviewKey(userID, period, startDate, endDate string) string {
	builder := NewCacheKeyBuilder(PrefixDashboardOverview)
	builder.WithUser(userID)
	builder.WithPeriod(period, startDate, endDate)
	return builder.Build()
}

// VisitStatsKey generates cache key for visit statistics
func (dc *DashboardCache) VisitStatsKey(userID, period, startDate, endDate string) string {
	builder := NewCacheKeyBuilder(PrefixDashboardVisitStats)
	builder.WithUser(userID)
	builder.WithPeriod(period, startDate, endDate)
	return builder.Build()
}

// SalesKey generates cache key for sales data
func (dc *DashboardCache) SalesKey(userID, period, startDate, endDate string) string {
	builder := NewCacheKeyBuilder(PrefixDashboardSales)
	builder.WithUser(userID)
	builder.WithPeriod(period, startDate, endDate)
	return builder.Build()
}



// AccountCache provides account-specific caching
type AccountCache struct {
	*ModuleCache
}

// NewAccountCache creates a new account cache helper
func NewAccountCache() *AccountCache {
	return &AccountCache{
		ModuleCache: NewModuleCache("account"),
	}
}

// LeadCache provides lead-specific caching
type LeadCache struct {
	*ModuleCache
}

// NewLeadCache creates a new lead cache helper
func NewLeadCache() *LeadCache {
	return &LeadCache{
		ModuleCache: NewModuleCache("lead"),
	}
}

// ContactCache provides contact-specific caching
type ContactCache struct {
	*ModuleCache
}

// NewContactCache creates a new contact cache helper
func NewContactCache() *ContactCache {
	return &ContactCache{
		ModuleCache: NewModuleCache("contact"),
	}
}

// ActivityCache provides activity-specific caching
type ActivityCache struct {
	*ModuleCache
}

// NewActivityCache creates a new activity cache helper
func NewActivityCache() *ActivityCache {
	return &ActivityCache{
		ModuleCache: NewModuleCache("activity"),
	}
}

// TaskCache provides task-specific caching
type TaskCache struct {
	*ModuleCache
}

// NewTaskCache creates a new task cache helper
func NewTaskCache() *TaskCache {
	return &TaskCache{
		ModuleCache: NewModuleCache("task"),
	}
}

// VisitReportCache provides visit report-specific caching
type VisitReportCache struct {
	*ModuleCache
}

// NewVisitReportCache creates a new visit report cache helper
func NewVisitReportCache() *VisitReportCache {
	return &VisitReportCache{
		ModuleCache: NewModuleCache("visit_report"),
	}
}

// BrickCache provides brick-specific caching
type BrickCache struct {
	*ModuleCache
}

// NewBrickCache creates a new brick cache helper
func NewBrickCache() *BrickCache {
	return &BrickCache{
		ModuleCache: NewModuleCache("brick"),
	}
}

// GeoJSONKey generates cache key for brick GeoJSON data
func (bc *BrickCache) GeoJSONKey(brickID string) string {
	return PrefixBrickGeoJSON + brickID
}

// ScheduleCache provides schedule-specific caching
type ScheduleCache struct {
	*ModuleCache
}

// NewScheduleCache creates a new schedule cache helper
func NewScheduleCache() *ScheduleCache {
	return &ScheduleCache{
		ModuleCache: NewModuleCache("schedule"),
	}
}

// ProductCache provides product-specific caching
type ProductCache struct {
	*ModuleCache
}

// NewProductCache creates a new product cache helper
func NewProductCache() *ProductCache {
	return &ProductCache{
		ModuleCache: NewModuleCache("product"),
	}
}

// DealCache provides deal/pipeline-specific caching
type DealCache struct {
	*ModuleCache
}

// NewDealCache creates a new deal cache helper
func NewDealCache() *DealCache {
	return &DealCache{
		ModuleCache: NewModuleCache("deal"),
	}
}

// SummaryKey generates cache key for deal summary
func (dc *DealCache) SummaryKey(userID string) string {
	if userID == "" {
		return PrefixDealSummary + "all"
	}
	return PrefixDealSummary + userID
}

// UserCache provides user-specific caching
type UserCache struct {
	*ModuleCache
}

// NewUserCache creates a new user cache helper
func NewUserCache() *UserCache {
	return &UserCache{
		ModuleCache: NewModuleCache("user"),
	}
}

// CachedListResult represents a cached list with pagination metadata
type CachedListResult struct {
	Data       interface{} `msgpack:"data"`
	Total      int64       `msgpack:"total"`
	Page       int         `msgpack:"page"`
	PerPage    int         `msgpack:"per_page"`
	TotalPages int         `msgpack:"total_pages"`
	CachedAt   time.Time   `msgpack:"cached_at"`
}

// NewCachedListResult creates a new cached list result
func NewCachedListResult(data interface{}, total int64, page, perPage int) *CachedListResult {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return &CachedListResult{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		CachedAt:   time.Now(),
	}
}

// CachedStatsResult represents cached statistics data
type CachedStatsResult struct {
	Data     interface{} `msgpack:"data"`
	CachedAt time.Time   `msgpack:"cached_at"`
}

// NewCachedStatsResult creates a new cached stats result
func NewCachedStatsResult(data interface{}) *CachedStatsResult {
	return &CachedStatsResult{
		Data:     data,
		CachedAt: time.Now(),
	}
}

// BatchCacheLoader provides utilities for batch loading data with caching
type BatchCacheLoader struct {
	ec        *AdvancedCache
	keyPrefix string
	ttl       time.Duration
}

// NewBatchCacheLoader creates a new batch cache loader
func NewBatchCacheLoader(keyPrefix string, ttl time.Duration) *BatchCacheLoader {
	return &BatchCacheLoader{
		ec:        Advanced(),
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

// LoadMany loads multiple items, using cache where available
// Returns map of id -> data for found items, and slice of missing ids
func (bcl *BatchCacheLoader) LoadMany(ids []string, targetFactory func() interface{}) (map[string]interface{}, []string, error) {
	if !bcl.ec.IsEnabled() || len(ids) == 0 {
		return nil, ids, nil
	}

	// Generate cache keys
	keys := make([]string, len(ids))
	keyToID := make(map[string]string)
	for i, id := range ids {
		key := bcl.keyPrefix + id
		keys[i] = key
		keyToID[key] = id
	}

	// Create targets for each key
	targets := make([]interface{}, len(keys))
	for i := range targets {
		targets[i] = targetFactory()
	}

	// Batch get from cache
	found, err := bcl.ec.MGet(keys, targets)
	if err != nil {
		return nil, ids, err
	}

	// Separate found and missing
	results := make(map[string]interface{})
	missing := make([]string, 0)

	for i, isFound := range found {
		id := keyToID[keys[i]]
		if isFound {
			results[id] = targets[i]
		} else {
			missing = append(missing, id)
		}
	}

	return results, missing, nil
}

// StoreMany stores multiple items in cache
func (bcl *BatchCacheLoader) StoreMany(items map[string]interface{}) error {
	if !bcl.ec.IsEnabled() || len(items) == 0 {
		return nil
	}

	cacheItems := make(map[string]interface{})
	for id, data := range items {
		cacheItems[bcl.keyPrefix+id] = data
	}

	return bcl.ec.MSet(cacheItems, bcl.ttl)
}
