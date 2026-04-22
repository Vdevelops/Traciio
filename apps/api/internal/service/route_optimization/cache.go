package route_optimization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
)

const (
	prefixRouteOptimizationDistanceMatrix = "route_optimization:distance_matrix:"
	prefixRouteOptimizationResult         = "route_optimization:result:"
)

// DistanceMatrixCache provides in-memory caching for distance matrices
// This reduces OSRM API calls significantly (80%+ reduction expected)
type DistanceMatrixCache struct {
	cache map[string]*CachedDistanceMatrix
	mutex sync.RWMutex
	ttl   time.Duration
	ac    *cachepkg.AdvancedCache
}

// CachedDistanceMatrix represents a cached distance matrix
type CachedDistanceMatrix struct {
	Matrix    [][]float64
	ExpiresAt time.Time
}

// NewDistanceMatrixCache creates a new distance matrix cache
func NewDistanceMatrixCache(ttl time.Duration) *DistanceMatrixCache {
	cache := &DistanceMatrixCache{
		cache: make(map[string]*CachedDistanceMatrix),
		ttl:   ttl,
		ac:    cachepkg.Advanced(),
	}

	// Start background cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached distance matrix
func (c *DistanceMatrixCache) Get(waypointHash string) ([][]float64, bool) {
	if c.ac != nil && c.ac.IsEnabled() {
		var cached CachedDistanceMatrix
		if found, _ := c.ac.Get(prefixRouteOptimizationDistanceMatrix+waypointHash, &cached); found {
			if !cached.ExpiresAt.IsZero() && time.Now().After(cached.ExpiresAt) {
				_ = c.ac.Delete(prefixRouteOptimizationDistanceMatrix + waypointHash)
			} else {
				// Return a copy to prevent external modification
				matrix := make([][]float64, len(cached.Matrix))
				for i := range cached.Matrix {
					matrix[i] = make([]float64, len(cached.Matrix[i]))
					copy(matrix[i], cached.Matrix[i])
				}
				return matrix, true
			}
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cached, exists := c.cache[waypointHash]
	if !exists || time.Now().After(cached.ExpiresAt) {
		return nil, false
	}

	// Return a copy to prevent external modification
	matrix := make([][]float64, len(cached.Matrix))
	for i := range cached.Matrix {
		matrix[i] = make([]float64, len(cached.Matrix[i]))
		copy(matrix[i], cached.Matrix[i])
	}

	return matrix, true
}

// Set stores a distance matrix in cache
func (c *DistanceMatrixCache) Set(waypointHash string, matrix [][]float64) {
	// Create a copy to prevent external modification
	cachedMatrix := make([][]float64, len(matrix))
	for i := range matrix {
		cachedMatrix[i] = make([]float64, len(matrix[i]))
		copy(cachedMatrix[i], matrix[i])
	}

	expiresAt := time.Now().Add(c.ttl)

	if c.ac != nil && c.ac.IsEnabled() {
		_ = c.ac.Set(prefixRouteOptimizationDistanceMatrix+waypointHash, &CachedDistanceMatrix{
			Matrix:    cachedMatrix,
			ExpiresAt: expiresAt,
		}, c.ttl)
	}

	c.mutex.Lock()
	c.cache[waypointHash] = &CachedDistanceMatrix{
		Matrix:    cachedMatrix,
		ExpiresAt: expiresAt,
	}
	c.mutex.Unlock()
}

// cleanup removes expired entries periodically
func (c *DistanceMatrixCache) cleanup() {
	ticker := time.NewTicker(1 * time.Hour) // Cleanup every hour
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, cached := range c.cache {
			if now.After(cached.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// HashWaypoints creates a consistent hash for waypoints
// This is used as cache key to identify same waypoint combinations
func HashWaypoints(waypoints []route_optimization.Waypoint) string {
	// Sort waypoints by lat/lng for consistent hash
	sorted := make([]route_optimization.Waypoint, len(waypoints))
	copy(sorted, waypoints)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Lat != sorted[j].Lat {
			return sorted[i].Lat < sorted[j].Lat
		}
		return sorted[i].Lng < sorted[j].Lng
	})

	// Create hash from sorted waypoints
	h := sha256.New()
	for _, wp := range sorted {
		h.Write([]byte(fmt.Sprintf("%.6f,%.6f", wp.Lat, wp.Lng)))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// RouteResultCache provides caching for complete route optimization results
type RouteResultCache struct {
	cache map[string]*CachedRouteResult
	mutex sync.RWMutex
	ttl   time.Duration
	ac    *cachepkg.AdvancedCache
}

// CachedRouteResult represents a cached route optimization result
type CachedRouteResult struct {
	OptimizedOrder []int
	TotalDistance  float64
	TotalDuration  int
	RouteSteps     []route_optimization.RouteStep
	RoutePolyline  *string
	ExpiresAt      time.Time
}

// NewRouteResultCache creates a new route result cache
func NewRouteResultCache(ttl time.Duration) *RouteResultCache {
	cache := &RouteResultCache{
		cache: make(map[string]*CachedRouteResult),
		ttl:   ttl,
		ac:    cachepkg.Advanced(),
	}

	// Start background cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached route result
func (c *RouteResultCache) Get(cacheKey string) (*CachedRouteResult, bool) {
	if c.ac != nil && c.ac.IsEnabled() {
		var cached CachedRouteResult
		if found, _ := c.ac.Get(prefixRouteOptimizationResult+cacheKey, &cached); found {
			if !cached.ExpiresAt.IsZero() && time.Now().After(cached.ExpiresAt) {
				_ = c.ac.Delete(prefixRouteOptimizationResult + cacheKey)
			} else {
				result := &CachedRouteResult{
					OptimizedOrder: make([]int, len(cached.OptimizedOrder)),
					TotalDistance:  cached.TotalDistance,
					TotalDuration:  cached.TotalDuration,
					ExpiresAt:      cached.ExpiresAt,
				}
				copy(result.OptimizedOrder, cached.OptimizedOrder)
				if len(cached.RouteSteps) > 0 {
					result.RouteSteps = make([]route_optimization.RouteStep, len(cached.RouteSteps))
					copy(result.RouteSteps, cached.RouteSteps)
				}
				if cached.RoutePolyline != nil {
					poly := *cached.RoutePolyline
					result.RoutePolyline = &poly
				}
				return result, true
			}
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cached, exists := c.cache[cacheKey]
	if !exists || time.Now().After(cached.ExpiresAt) {
		return nil, false
	}

	// Return a copy
	result := &CachedRouteResult{
		OptimizedOrder: make([]int, len(cached.OptimizedOrder)),
		TotalDistance:  cached.TotalDistance,
		TotalDuration:  cached.TotalDuration,
		ExpiresAt:      cached.ExpiresAt,
	}
	copy(result.OptimizedOrder, cached.OptimizedOrder)
	if len(cached.RouteSteps) > 0 {
		result.RouteSteps = make([]route_optimization.RouteStep, len(cached.RouteSteps))
		copy(result.RouteSteps, cached.RouteSteps)
	}
	if cached.RoutePolyline != nil {
		poly := *cached.RoutePolyline
		result.RoutePolyline = &poly
	}

	return result, true
}

// Set stores a route result in cache
func (c *RouteResultCache) Set(
	cacheKey string,
	order []int,
	distance float64,
	duration int,
	steps []route_optimization.RouteStep,
	polyline *string,
) {
	expiresAt := time.Now().Add(c.ttl)
	var polyCopy *string
	if polyline != nil {
		p := *polyline
		polyCopy = &p
	}

	stepsCopy := []route_optimization.RouteStep(nil)
	if len(steps) > 0 {
		stepsCopy = make([]route_optimization.RouteStep, len(steps))
		copy(stepsCopy, steps)
	}

	if c.ac != nil && c.ac.IsEnabled() {
		_ = c.ac.Set(prefixRouteOptimizationResult+cacheKey, &CachedRouteResult{
			OptimizedOrder: append([]int(nil), order...),
			TotalDistance:  distance,
			TotalDuration:  duration,
			RouteSteps:     stepsCopy,
			RoutePolyline:  polyCopy,
			ExpiresAt:      expiresAt,
		}, c.ttl)
	}

	c.mutex.Lock()
	c.cache[cacheKey] = &CachedRouteResult{
		OptimizedOrder: append([]int(nil), order...), // Copy
		TotalDistance:  distance,
		TotalDuration:  duration,
		RouteSteps:     stepsCopy,
		RoutePolyline:  polyCopy,
		ExpiresAt:      expiresAt,
	}
	c.mutex.Unlock()
}

// cleanup removes expired entries periodically
func (c *RouteResultCache) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, cached := range c.cache {
			if now.After(cached.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// CreateRouteCacheKey creates a cache key for route optimization
func CreateRouteCacheKey(startLat, startLng float64, waypointHash string) string {
	h := sha256.New()
	base := fmt.Sprintf("%.6f,%.6f:%s", startLat, startLng, waypointHash)
	h.Write([]byte(base))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateRouteCacheKeyWithExtras creates a cache key with additional discriminators.
// Use this when route results can differ based on request parameters (e.g. optimization_type).
// Note: this keeps the hashing stable and avoids accidental cache collisions.
func CreateRouteCacheKeyWithExtras(startLat, startLng float64, waypointHash string, extras ...string) string {
	if len(extras) == 0 {
		return CreateRouteCacheKey(startLat, startLng, waypointHash)
	}
	h := sha256.New()
	base := fmt.Sprintf("%.6f,%.6f:%s", startLat, startLng, waypointHash)
	h.Write([]byte(base))
	h.Write([]byte(":"))
	h.Write([]byte(strings.Join(extras, "|")))
	return hex.EncodeToString(h.Sum(nil))
}
