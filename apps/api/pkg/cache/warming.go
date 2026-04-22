// Package cache provides advanced Redis caching with cache warming support.
//
// This file implements cache warming strategies for reference/master data.
// Cache warming pre-loads frequently accessed data at startup to avoid
// cold cache issues and improve initial request latency.
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/logger"
)

// WarmingConfig defines cache warming configuration
type WarmingConfig struct {
	// EnableWarmingOnStartup enables automatic cache warming at application startup
	EnableWarmingOnStartup bool
	// WarmingConcurrency controls the number of parallel warming operations
	WarmingConcurrency int
	// WarmingTimeout is the maximum time allowed for warming operations
	WarmingTimeout time.Duration
}

// DefaultWarmingConfig returns sensible defaults for cache warming
func DefaultWarmingConfig() *WarmingConfig {
	return &WarmingConfig{
		EnableWarmingOnStartup: true,
		WarmingConcurrency:     5,
		WarmingTimeout:         30 * time.Second,
	}
}

// WarmingTask represents a cache warming task
type WarmingTask struct {
	Name   string
	Key    string
	TTL    time.Duration
	Loader func(ctx context.Context) (interface{}, error)
}

// CacheWarmer handles cache warming operations
type CacheWarmer struct {
	config *WarmingConfig
	tasks  []WarmingTask
	mu     sync.Mutex
}

// NewCacheWarmer creates a new cache warmer instance
func NewCacheWarmer(config *WarmingConfig) *CacheWarmer {
	if config == nil {
		config = DefaultWarmingConfig()
	}
	return &CacheWarmer{
		config: config,
		tasks:  make([]WarmingTask, 0),
	}
}

// RegisterTask adds a warming task to be executed on startup
func (cw *CacheWarmer) RegisterTask(task WarmingTask) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.tasks = append(cw.tasks, task)
}

// RegisterTasks adds multiple warming tasks
func (cw *CacheWarmer) RegisterTasks(tasks ...WarmingTask) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.tasks = append(cw.tasks, tasks...)
}

// WarmupResult contains the result of a warming operation
type WarmupResult struct {
	TaskName  string
	Success   bool
	Duration  time.Duration
	Error     error
	CacheKey  string
}

// ExecuteWarming runs all registered warming tasks
func (cw *CacheWarmer) ExecuteWarming(ctx context.Context) []WarmupResult {
	if !Advanced().IsEnabled() {
		logger.LogInfo("Cache warming skipped: cache not enabled", nil)
		return nil
	}

	cw.mu.Lock()
	tasks := make([]WarmingTask, len(cw.tasks))
	copy(tasks, cw.tasks)
	cw.mu.Unlock()

	if len(tasks) == 0 {
		logger.LogInfo("Cache warming: no tasks registered", nil)
		return nil
	}

	logger.LogInfo("Starting cache warming", map[string]interface{}{
		"total_tasks":  len(tasks),
		"concurrency":  cw.config.WarmingConcurrency,
		"timeout":      cw.config.WarmingTimeout.String(),
	})

	// Create context with timeout
	warmCtx, cancel := context.WithTimeout(ctx, cw.config.WarmingTimeout)
	defer cancel()

	// Channel for results
	results := make(chan WarmupResult, len(tasks))
	
	// Semaphore for concurrency control
	sem := make(chan struct{}, cw.config.WarmingConcurrency)

	// Execute tasks concurrently
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t WarmingTask) {
			defer wg.Done()
			
			// Acquire semaphore
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-warmCtx.Done():
				results <- WarmupResult{
					TaskName: t.Name,
					CacheKey: t.Key,
					Success:  false,
					Error:    warmCtx.Err(),
				}
				return
			}

			// Execute warming
			start := time.Now()
			data, err := t.Loader(warmCtx)
			if err != nil {
				results <- WarmupResult{
					TaskName: t.Name,
					CacheKey: t.Key,
					Success:  false,
					Duration: time.Since(start),
					Error:    err,
				}
				return
			}

			// Store in cache
			if err := Advanced().Set(t.Key, data, t.TTL); err != nil {
				results <- WarmupResult{
					TaskName: t.Name,
					CacheKey: t.Key,
					Success:  false,
					Duration: time.Since(start),
					Error:    err,
				}
				return
			}

			results <- WarmupResult{
				TaskName: t.Name,
				CacheKey: t.Key,
				Success:  true,
				Duration: time.Since(start),
			}
		}(task)
	}

	// Close results channel when all tasks complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var warmupResults []WarmupResult
	successCount := 0
	failCount := 0
	totalDuration := time.Duration(0)

	for result := range results {
		warmupResults = append(warmupResults, result)
		totalDuration += result.Duration
		if result.Success {
			successCount++
			logger.LogInfo("Cache warming task completed", map[string]interface{}{
				"task":     result.TaskName,
				"key":      result.CacheKey,
				"duration": result.Duration.String(),
			})
		} else {
			failCount++
			errMsg := "unknown"
			if result.Error != nil {
				errMsg = result.Error.Error()
			}
			logger.LogError(result.Error, map[string]interface{}{
				"task":    result.TaskName,
				"key":     result.CacheKey,
				"error":   errMsg,
				"context": "cache_warming_failed",
			})
		}
	}

	logger.LogInfo("Cache warming completed", map[string]interface{}{
		"success":        successCount,
		"failed":         failCount,
		"total":          len(tasks),
		"total_duration": totalDuration.String(),
	})

	return warmupResults
}

// Global cache warmer instance
var (
	globalWarmer     *CacheWarmer
	globalWarmerOnce sync.Once
)

// Warmer returns the global cache warmer instance
func Warmer() *CacheWarmer {
	globalWarmerOnce.Do(func() {
		globalWarmer = NewCacheWarmer(DefaultWarmingConfig())
	})
	return globalWarmer
}

// RegisterReferenceDataWarming registers common reference data warming tasks
// Call this during application initialization to set up warming for master data
func RegisterReferenceDataWarming(
	productLoader func(ctx context.Context) (interface{}, error),
	rolePermissionsLoader func(ctx context.Context) (interface{}, error),
) {
	warmer := Warmer()

	// Products - master data, TTL 30 minutes
	if productLoader != nil {
		warmer.RegisterTask(WarmingTask{
			Name:   "Products",
			Key:    "products:list:default",
			TTL:    TTLReferenceMedium,
			Loader: productLoader,
		})
	}

	// Role Permissions - RBAC data, TTL 20 minutes
	if rolePermissionsLoader != nil {
		warmer.RegisterTask(WarmingTask{
			Name:   "RolePermissions",
			Key:    "rbac:all_roles:permissions",
			TTL:    TTLRBACMedium,
			Loader: rolePermissionsLoader,
		})
	}

}

// WarmOnStartup should be called during application startup
// It executes all registered warming tasks in the background
func WarmOnStartup(ctx context.Context) {
	go func() {
		// Small delay to allow application to fully initialize
		time.Sleep(2 * time.Second)
		Warmer().ExecuteWarming(ctx)
	}()
}
