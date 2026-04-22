package route_optimization

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
)

const prefixRouteOptimizationSegment = "route_optimization:segment:"

// PerformanceOptimizer provides optimized route calculation with two-phase approach
// Phase 1: Fast optimization using simplified routing (for initial order)
// Phase 2: Accurate routing using full routing (for final route)
type PerformanceOptimizer struct {
	distanceMatrixCache *DistanceMatrixCache
	routeSegmentCache   *RouteSegmentCache
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer(
	distanceMatrixCache *DistanceMatrixCache,
	routeSegmentCache *RouteSegmentCache,
) *PerformanceOptimizer {
	return &PerformanceOptimizer{
		distanceMatrixCache: distanceMatrixCache,
		routeSegmentCache:   routeSegmentCache,
	}
}

// OptimizeRouteFast performs fast route optimization using two-phase approach
// This is significantly faster than single-phase optimization while maintaining accuracy
func (po *PerformanceOptimizer) OptimizeRouteFast(
	osrmClient *OSRMClient,
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	optimizationType string,
) ([]int, error) {
	if len(waypoints) == 0 {
		return []int{}, nil
	}

	// Normalize optimizationType
	if optimizationType != "duration" {
		optimizationType = "distance"
	}

	// Best-effort: when OSRM table is enabled and waypoint count is within limit,
	// build a (start + waypoints) road cost matrix and optimize against it.
	// This avoids pathological routes where the first stop is far away due to
	// missing start->first cost or inaccurate straight-line estimates.
	if costMatrix, ok := po.tryBuildRoadCostMatrixWithStart(osrmClient, startLat, startLng, waypoints, optimizationType); ok {
		if len(waypoints) <= 8 {
			return exactTSPWithFullMatrix(costMatrix, len(waypoints)), nil
		}

		initialOrder := nearestNeighborFromStartWithFullMatrix(costMatrix, len(waypoints))
		if len(waypoints) >= 4 {
			return twoOptWithFullMatrix(initialOrder, costMatrix), nil
		}
		return initialOrder, nil
	}

	// Phase 1: Fast initial optimization using simplified routing
	// Use simplified OSRM requests (faster) or cached matrix for initial order
	initialOrder, err := po.optimizePhase1Fast(osrmClient, startLat, startLng, waypoints, optimizationType)
	if err != nil {
		// Fallback to nearest neighbor
		initialOrder, _, _ = nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
	}

	// Phase 2: Improve with exact TSP or 2-Opt using cached distance matrix
	distanceMatrix, durationMatrix, err := po.getOrCalculateDistanceMatrixFast(osrmClient, waypoints)
	if err == nil {
		costMatrix := distanceMatrix
		if optimizationType == "duration" {
			costMatrix = floatMatrixFromInt(durationMatrix)
		}

		// Build augmented (start + waypoints) matrix so that start→first-stop cost
		// uses real road distances instead of Haversine straight-line estimates.
		// This prevents pathological first-stop choices in cities with winding roads.
		augmented := po.buildAugmentedCostMatrix(osrmClient, startLat, startLng, waypoints, costMatrix, optimizationType)

		if len(waypoints) <= 8 {
			return exactTSPWithFullMatrix(augmented, len(waypoints)), nil
		}
		initialAugmented := nearestNeighborFromStartWithFullMatrix(augmented, len(waypoints))
		if len(waypoints) >= 4 {
			return twoOptWithFullMatrix(initialAugmented, augmented), nil
		}
		return initialAugmented, nil
	}

	// Fallback if matrix calculation fails
	if len(waypoints) <= 8 {
		costMatrix := po.calculateDistanceMatrixHaversine(waypoints)
		return po.exactTSPWithCostMatrix(startLat, startLng, waypoints, costMatrix, optimizationType), nil
	}

	if len(waypoints) >= 4 {
		optimizedOrder, _ := TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
		return optimizedOrder, nil
	}

	return initialOrder, nil
}

// tryBuildRoadCostMatrixWithStart uses OSRM /table to build an (n+1)x(n+1) cost matrix
// for [start, wp1..wpN]. Costs are either km (distance) or seconds (duration).
func (po *PerformanceOptimizer) tryBuildRoadCostMatrixWithStart(
	osrmClient *OSRMClient,
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	optimizationType string,
) ([][]float64, bool) {
	if osrmClient == nil {
		return nil, false
	}
	if config.AppConfig == nil || !config.AppConfig.OSRM.TableEnabled {
		return nil, false
	}
	maxN := config.AppConfig.OSRM.TableMaxWaypoints
	if maxN <= 0 {
		maxN = 25
	}
	// We need start + waypoints
	if len(waypoints)+1 > maxN {
		return nil, false
	}

	coords := make([][]float64, 0, len(waypoints)+1)
	coords = append(coords, []float64{startLng, startLat})
	for _, wp := range waypoints {
		coords = append(coords, []float64{wp.Lng, wp.Lat})
	}

	tableResp, err := osrmClient.Table(TableRequest{
		Coordinates: coords,
		Profile:     "driving",
		Annotations: []string{"distance", "duration"},
	})
	if err != nil {
		return nil, false
	}
	n := len(coords)
	if len(tableResp.Distances) != n || len(tableResp.Durations) != n {
		return nil, false
	}
	for i := 0; i < n; i++ {
		if len(tableResp.Distances[i]) != n || len(tableResp.Durations[i]) != n {
			return nil, false
		}
	}

	cost := make([][]float64, n)
	for i := 0; i < n; i++ {
		cost[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if optimizationType == "duration" {
				cost[i][j] = tableResp.Durations[i][j]
			} else {
				cost[i][j] = tableResp.Distances[i][j] / 1000.0
			}
		}
	}

	return cost, true
}

// nearestNeighborFromStartWithFullMatrix returns an order of waypoint indices [0..n-1]
// optimized using a full matrix where 0 is start and i+1 corresponds to waypoint i.
func nearestNeighborFromStartWithFullMatrix(costMatrix [][]float64, waypointCount int) []int {
	n := waypointCount
	if n <= 0 {
		return []int{}
	}
	visited := make([]bool, n)
	order := make([]int, 0, n)
	currentNode := 0 // start node

	for len(order) < n {
		bestIdx := -1
		bestCost := 1e18
		for wpIdx := 0; wpIdx < n; wpIdx++ {
			if visited[wpIdx] {
				continue
			}
			node := wpIdx + 1
			c := costMatrix[currentNode][node]
			if c < bestCost {
				bestCost = c
				bestIdx = wpIdx
			}
		}
		if bestIdx < 0 {
			break
		}
		visited[bestIdx] = true
		order = append(order, bestIdx)
		currentNode = bestIdx + 1
	}

	return order
}

// twoOptWithFullMatrix performs 2-opt using a full cost matrix where 0 is start and i+1 is waypoint i.
func twoOptWithFullMatrix(initialOrder []int, costMatrix [][]float64) []int {
	if len(initialOrder) < 4 {
		return initialOrder
	}

	bestOrder := make([]int, len(initialOrder))
	copy(bestOrder, initialOrder)
	bestCost := totalCostWithFullMatrix(bestOrder, costMatrix)

	reverseSegment := func(order []int, i, j int) []int {
		newOrder := make([]int, len(order))
		copy(newOrder, order)
		for k := 0; k <= (j-i-1)/2; k++ {
			newOrder[i+1+k], newOrder[j-k] = newOrder[j-k], newOrder[i+1+k]
		}
		return newOrder
	}

	improved := true
	maxIterations := 80
	iteration := 0

	for improved && iteration < maxIterations {
		improved = false
		iteration++
		for i := 0; i < len(bestOrder)-1; i++ {
			for j := i + 2; j < len(bestOrder); j++ {
				newOrder := reverseSegment(bestOrder, i, j)
				newCost := totalCostWithFullMatrix(newOrder, costMatrix)
				if newCost < bestCost {
					bestOrder = newOrder
					bestCost = newCost
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}
	}

	return bestOrder
}

func totalCostWithFullMatrix(order []int, costMatrix [][]float64) float64 {
	if len(order) == 0 {
		return 0
	}
	// start -> first
	total := costMatrix[0][order[0]+1]
	// between waypoints
	for i := 0; i < len(order)-1; i++ {
		from := order[i] + 1
		to := order[i+1] + 1
		total += costMatrix[from][to]
	}
	return total
}

// optimizePhase1Fast performs fast initial optimization
// Uses simplified OSRM requests or distance matrix for speed
func (po *PerformanceOptimizer) optimizePhase1Fast(
	osrmClient *OSRMClient,
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	optimizationType string,
) ([]int, error) {
	// Try to use cached distance matrix first (fastest)
	waypointHash := HashWaypoints(waypoints)
	if cachedMatrix, found := po.distanceMatrixCache.Get(waypointHash); found {
		// Use cached matrix for nearest neighbor
		costMatrix := cachedMatrix
		if optimizationType == "duration" {
			costMatrix = floatMatrixFromInt(estimateDurationMatrix(cachedMatrix))
		}
		return po.nearestNeighborWithCostMatrix(startLat, startLng, waypoints, costMatrix, optimizationType), nil
	}

	// If no cache, use Haversine for initial optimization (very fast)
	// NOTE: This ignores optimizationType, but still yields a good seed order.
	initialOrder, _, _ := nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
	return initialOrder, nil
}

// getOrCalculateDistanceMatrixFast gets or calculates distance matrix efficiently
func (po *PerformanceOptimizer) getOrCalculateDistanceMatrixFast(
	osrmClient *OSRMClient,
	waypoints []route_optimization.Waypoint,
) ([][]float64, [][]int, error) {
	// Check cache first
	waypointHash := HashWaypoints(waypoints)
	if cachedMatrix, found := po.distanceMatrixCache.Get(waypointHash); found {
		durationMatrix := estimateDurationMatrix(cachedMatrix)
		return cachedMatrix, durationMatrix, nil
	}

	// If not cached and waypoints are many, use simplified approach
	// For small waypoints, calculate full matrix in parallel
	if len(waypoints) <= 10 {
		if config.AppConfig != nil && config.AppConfig.OSRM.TableEnabled {
			maxN := config.AppConfig.OSRM.TableMaxWaypoints
			if maxN <= 0 {
				maxN = 25
			}
			if len(waypoints) <= maxN {
				matrixStart := time.Now()
				coords := make([][]float64, 0, len(waypoints))
				for _, wp := range waypoints {
					coords = append(coords, []float64{wp.Lng, wp.Lat})
				}
				tableResp, err := osrmClient.Table(TableRequest{
					Coordinates: coords,
					Profile:     "driving",
					Annotations: []string{"distance", "duration"},
				})
				// Record matrix duration as table method when successful.
				if err == nil && len(tableResp.Distances) == len(waypoints) {
					routeOptimizationMatrixDurationSeconds.WithLabelValues("table").Observe(time.Since(matrixStart).Seconds())
					n := len(waypoints)
					distanceMatrix := make([][]float64, n)
					durationMatrix := make([][]int, n)
					valid := true
					for i := 0; i < n; i++ {
						if len(tableResp.Distances[i]) != n || len(tableResp.Durations[i]) != n {
							valid = false
							break
						}
						distanceMatrix[i] = make([]float64, n)
						durationMatrix[i] = make([]int, n)
						for j := 0; j < n; j++ {
							distanceMatrix[i][j] = tableResp.Distances[i][j] / 1000.0
							durationMatrix[i][j] = int(tableResp.Durations[i][j])
						}
					}
					if valid {
						waypointHash := HashWaypoints(waypoints)
						po.distanceMatrixCache.Set(waypointHash, distanceMatrix)
						return distanceMatrix, durationMatrix, nil
					}
				}
			}
		}
		// Calculate full matrix in parallel (acceptable for small sets)
		return po.calculateDistanceMatrixParallelFast(osrmClient, waypoints)
	}

	// For large waypoints, use Haversine for speed (good enough for optimization)
	// OSRM will be used for final route anyway
	return po.calculateDistanceMatrixHaversine(waypoints), nil, nil
}

// calculateDistanceMatrixParallelFast calculates distance matrix with optimized settings
func (po *PerformanceOptimizer) calculateDistanceMatrixParallelFast(
	osrmClient *OSRMClient,
	waypoints []route_optimization.Waypoint,
) ([][]float64, [][]int, error) {
	n := len(waypoints)
	distanceMatrix := make([][]float64, n)
	durationMatrix := make([][]int, n)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, n)
		durationMatrix[i] = make([]int, n)
	}

	// Use worker pool with optimized settings
	const maxWorkers = 15 // Increased for faster parallel processing
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	// Calculate distances in parallel with optimized OSRM requests
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// Use simplified OSRM request (faster)
				distance, duration, err := po.calculateDistanceSimplified(
					osrmClient,
					waypoints[i].Lat, waypoints[i].Lng,
					waypoints[j].Lat, waypoints[j].Lng,
				)

				mu.Lock()
				if err != nil && firstErr == nil {
					firstErr = err
				} else if err == nil {
					distanceMatrix[i][j] = distance
					distanceMatrix[j][i] = distance
					durationMatrix[i][j] = duration
					durationMatrix[j][i] = duration
				}
				mu.Unlock()
			}(i, j)
		}
	}

	wg.Wait()

	if firstErr != nil {
		return nil, nil, firstErr
	}

	// Cache the result
	waypointHash := HashWaypoints(waypoints)
	po.distanceMatrixCache.Set(waypointHash, distanceMatrix)

	return distanceMatrix, durationMatrix, nil
}

// calculateDistanceSimplified calculates distance using simplified OSRM request
// This is faster than full routing but still follows roads
func (po *PerformanceOptimizer) calculateDistanceSimplified(
	osrmClient *OSRMClient,
	lat1, lng1, lat2, lng2 float64,
) (float64, int, error) {
	coordinates := [][]float64{
		{lng1, lat1},
		{lng2, lat2},
	}

	// Use simplified request (faster response)
	routeReq := RouteRequest{
		Coordinates:  coordinates,
		Profile:      "driving",
		Overview:     "simplified", // Simplified is faster than full
		Geometries:   "polyline",
		Steps:        false, // No steps needed for distance matrix
		Alternatives: false,
	}

	osrmResp, err := osrmClient.Route(routeReq)
	if err != nil {
		// Fallback to Haversine
		distance := haversineDistance(lat1, lng1, lat2, lng2)
		avgSpeedKmh := 40.0
		duration := int((distance / avgSpeedKmh) * 3600)
		return distance, duration, nil
	}

	if len(osrmResp.Routes) == 0 {
		distance := haversineDistance(lat1, lng1, lat2, lng2)
		avgSpeedKmh := 40.0
		duration := int((distance / avgSpeedKmh) * 3600)
		return distance, duration, nil
	}

	osrmRoute := osrmResp.Routes[0]
	distance := osrmRoute.Distance / 1000.0
	duration := int(osrmRoute.Duration)

	return distance, duration, nil
}

// calculateDistanceMatrixHaversine calculates distance matrix using Haversine
// Very fast but less accurate (straight line distance)
func (po *PerformanceOptimizer) calculateDistanceMatrixHaversine(
	waypoints []route_optimization.Waypoint,
) [][]float64 {
	n := len(waypoints)
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
		for j := range matrix[i] {
			if i != j {
				matrix[i][j] = haversineDistance(
					waypoints[i].Lat, waypoints[i].Lng,
					waypoints[j].Lat, waypoints[j].Lng,
				)
			}
		}
	}
	return matrix
}

// nearestNeighborWithDistanceMatrix performs nearest neighbor using distance matrix
func (po *PerformanceOptimizer) nearestNeighborWithCostMatrix(
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	costMatrix [][]float64,
	optimizationType string,
) []int {
	n := len(waypoints)
	if n == 0 {
		return []int{}
	}

	visited := make([]bool, n)
	order := make([]int, 0, n)
	currentLat := startLat
	currentLng := startLng

	for len(order) < n {
		nearestIndex := -1
		minCost := 1e18

		for i := 0; i < n; i++ {
			if !visited[i] {
				var cost float64
				if len(order) == 0 {
					// First waypoint: include start -> waypoint cost.
					cost = startToWaypointCost(currentLat, currentLng, waypoints[i].Lat, waypoints[i].Lng, optimizationType)
				} else {
					// Use matrix for between-waypoint cost
					lastIndex := order[len(order)-1]
					cost = costMatrix[lastIndex][i]
				}

				if cost < minCost {
					minCost = cost
					nearestIndex = i
				}
			}
		}

		if nearestIndex >= 0 {
			order = append(order, nearestIndex)
			visited[nearestIndex] = true
			currentLat = waypoints[nearestIndex].Lat
			currentLng = waypoints[nearestIndex].Lng
		}
	}

	return order
}

// twoOptWithCostMatrix performs 2-Opt optimization using a cost matrix.
// IMPORTANT: Includes start -> first waypoint cost, so the first stop is optimized too.
func (po *PerformanceOptimizer) twoOptWithCostMatrix(
	initialOrder []int,
	waypoints []route_optimization.Waypoint,
	startLat, startLng float64,
	costMatrix [][]float64,
	optimizationType string,
) []int {
	if len(initialOrder) < 4 {
		return initialOrder
	}

	bestOrder := make([]int, len(initialOrder))
	copy(bestOrder, initialOrder)
	bestCost := po.calculateTotalCostWithMatrix(bestOrder, waypoints, startLat, startLng, costMatrix, optimizationType)

	improved := true
	maxIterations := 50 // Reduced for faster optimization
	iteration := 0
	deadline, hasDeadline := improvementDeadline(len(bestOrder))

	// Helper function to reverse segment
	reverseSegment := func(order []int, i, j int) []int {
		newOrder := make([]int, len(order))
		copy(newOrder, order)
		// Reverse segment from i+1 to j
		for k := 0; k <= (j-i-1)/2; k++ {
			newOrder[i+1+k], newOrder[j-k] = newOrder[j-k], newOrder[i+1+k]
		}
		return newOrder
	}

	for improved && iteration < maxIterations {
		if hasDeadline && time.Now().After(deadline) {
			break
		}
		improved = false
		iteration++

		for i := 0; i < len(bestOrder)-1; i++ {
			if hasDeadline && time.Now().After(deadline) {
				break
			}
			for j := i + 2; j < len(bestOrder); j++ {
				if hasDeadline && time.Now().After(deadline) {
					break
				}
				newOrder := reverseSegment(bestOrder, i, j)
				newCost := po.calculateTotalCostWithMatrix(newOrder, waypoints, startLat, startLng, costMatrix, optimizationType)

				if newCost < bestCost {
					bestOrder = newOrder
					bestCost = newCost
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}
	}

	return bestOrder
}

// calculateTotalCostWithMatrix calculates total cost using a cost matrix.
// Includes start -> first waypoint cost.
func (po *PerformanceOptimizer) calculateTotalCostWithMatrix(
	order []int,
	waypoints []route_optimization.Waypoint,
	startLat, startLng float64,
	costMatrix [][]float64,
	optimizationType string,
) float64 {
	if len(order) < 2 {
		if len(order) == 1 {
			wp := waypoints[order[0]]
			return startToWaypointCost(startLat, startLng, wp.Lat, wp.Lng, optimizationType)
		}
		return 0.0
	}

	total := 0.0
	// Start -> first
	first := waypoints[order[0]]
	total += startToWaypointCost(startLat, startLng, first.Lat, first.Lng, optimizationType)

	for i := 0; i < len(order)-1; i++ {
		total += costMatrix[order[i]][order[i+1]]
	}

	return total
}

func floatMatrixFromInt(matrix [][]int) [][]float64 {
	if matrix == nil {
		return nil
	}
	out := make([][]float64, len(matrix))
	for i := range matrix {
		out[i] = make([]float64, len(matrix[i]))
		for j := range matrix[i] {
			out[i][j] = float64(matrix[i][j])
		}
	}
	return out
}

func startToWaypointCost(startLat, startLng, wpLat, wpLng float64, optimizationType string) float64 {
	// Distance cost: kilometers
	distanceKm := haversineDistance(startLat, startLng, wpLat, wpLng)
	if optimizationType != "duration" {
		return distanceKm
	}
	// Duration cost: seconds (rough estimate). Good enough to avoid extreme first-stop detours.
	avgSpeedKmh := 40.0
	seconds := (distanceKm / avgSpeedKmh) * 3600
	return seconds
}

// buildAugmentedCostMatrix builds an (n+1)×(n+1) cost matrix where index 0 is the
// start location and indices 1..n correspond to waypoints. Start→waypoint costs
// are computed via OSRM pairwise calls (falling back to Haversine), while
// waypoint-to-waypoint costs come from the pre-computed waypointCostMatrix.
//
// This eliminates the Haversine-vs-OSRM inconsistency that causes the optimizer
// to pick a sub-optimal first stop in cities with winding road layouts.
func (po *PerformanceOptimizer) buildAugmentedCostMatrix(
	osrmClient *OSRMClient,
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	waypointCostMatrix [][]float64,
	optimizationType string,
) [][]float64 {
	n := len(waypoints)
	augmented := make([][]float64, n+1)
	for i := range augmented {
		augmented[i] = make([]float64, n+1)
	}

	// Fill waypoint→waypoint block (indices 1..n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			augmented[i+1][j+1] = waypointCostMatrix[i][j]
		}
	}

	// Compute start→each waypoint cost via OSRM pairwise calls (parallel for speed)
	type startCost struct {
		idx      int
		distance float64
		duration int
	}
	results := make([]startCost, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			dist, dur, err := po.calculateDistanceSimplified(
				osrmClient,
				startLat, startLng,
				waypoints[idx].Lat, waypoints[idx].Lng,
			)
			if err != nil || dist <= 0 {
				// Fallback to Haversine
				dist = haversineDistance(startLat, startLng, waypoints[idx].Lat, waypoints[idx].Lng)
				dur = int((dist / 40.0) * 3600)
			}
			results[idx] = startCost{idx: idx, distance: dist, duration: dur}
		}(i)
	}
	wg.Wait()

	for _, r := range results {
		var cost float64
		if optimizationType == "duration" {
			cost = float64(r.duration)
		} else {
			cost = r.distance
		}
		augmented[0][r.idx+1] = cost
		augmented[r.idx+1][0] = cost
	}

	return augmented
}

// RouteSegmentCache caches route segments for common routes
type RouteSegmentCache struct {
	cache map[string]*CachedRouteSegment
	mutex sync.RWMutex
	ttl   time.Duration
	ac    *cachepkg.AdvancedCache
}

// CachedRouteSegment represents a cached route segment
type CachedRouteSegment struct {
	Distance  float64
	Duration  int
	ExpiresAt time.Time
}

// NewRouteSegmentCache creates a new route segment cache
func NewRouteSegmentCache(ttl time.Duration) *RouteSegmentCache {
	cache := &RouteSegmentCache{
		cache: make(map[string]*CachedRouteSegment),
		ttl:   ttl,
		ac:    cachepkg.Advanced(),
	}

	// Start background cleanup
	go cache.cleanup()

	return cache
}

// Get retrieves a cached route segment
func (c *RouteSegmentCache) Get(key string) (*CachedRouteSegment, bool) {
	if c.ac != nil && c.ac.IsEnabled() {
		var cached CachedRouteSegment
		if found, _ := c.ac.Get(prefixRouteOptimizationSegment+key, &cached); found {
			if !cached.ExpiresAt.IsZero() && time.Now().After(cached.ExpiresAt) {
				_ = c.ac.Delete(prefixRouteOptimizationSegment + key)
			} else {
				return &cached, true
			}
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cached, exists := c.cache[key]
	if !exists || time.Now().After(cached.ExpiresAt) {
		return nil, false
	}

	return cached, true
}

// Set stores a route segment in cache
func (c *RouteSegmentCache) Set(key string, distance float64, duration int) {
	expiresAt := time.Now().Add(c.ttl)

	if c.ac != nil && c.ac.IsEnabled() {
		_ = c.ac.Set(prefixRouteOptimizationSegment+key, &CachedRouteSegment{
			Distance:  distance,
			Duration:  duration,
			ExpiresAt: expiresAt,
		}, c.ttl)
	}

	c.mutex.Lock()
	c.cache[key] = &CachedRouteSegment{
		Distance:  distance,
		Duration:  duration,
		ExpiresAt: expiresAt,
	}
	c.mutex.Unlock()
}

// cleanup removes expired entries
func (c *RouteSegmentCache) cleanup() {
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

// CreateRouteSegmentKey creates a cache key for a route segment
func CreateRouteSegmentKey(lat1, lng1, lat2, lng2 float64) string {
	// Round coordinates to 4 decimal places (~11 meters precision) for cache efficiency
	// This allows similar routes to share cache
	lat1Rounded := float64(int(lat1*10000)) / 10000
	lng1Rounded := float64(int(lng1*10000)) / 10000
	lat2Rounded := float64(int(lat2*10000)) / 10000
	lng2Rounded := float64(int(lng2*10000)) / 10000

	// Ensure consistent key order
	if lat1Rounded > lat2Rounded || (lat1Rounded == lat2Rounded && lng1Rounded > lng2Rounded) {
		lat1Rounded, lat2Rounded = lat2Rounded, lat1Rounded
		lng1Rounded, lng2Rounded = lng2Rounded, lng1Rounded
	}

	return fmt.Sprintf("segment:%.4f,%.4f:%.4f,%.4f", lat1Rounded, lng1Rounded, lat2Rounded, lng2Rounded)
}

// exactTSPWithFullMatrix calculates exact optimal order using brute force for small N
func exactTSPWithFullMatrix(costMatrix [][]float64, n int) []int {
	bestCost := math.MaxFloat64
	var bestOrder []int

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}

	var generate func(k int)
	generate = func(k int) {
		if k == n {
			cost := costMatrix[0][order[0]+1] // start to first
			for i := 0; i < n-1; i++ {
				cost += costMatrix[order[i]+1][order[i+1]+1]
			}
			if cost < bestCost {
				bestCost = cost
				bestOrder = make([]int, n)
				copy(bestOrder, order)
			}
			return
		}
		for i := k; i < n; i++ {
			order[k], order[i] = order[i], order[k]
			generate(k + 1)
			order[k], order[i] = order[i], order[k]
		}
	}

	generate(0)
	return bestOrder
}

// exactTSPWithCostMatrix calculates exact optimal order using brute force for small N with provided matrix and start location
func (po *PerformanceOptimizer) exactTSPWithCostMatrix(
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	costMatrix [][]float64,
	optimizationType string,
) []int {
	n := len(waypoints)
	if n == 0 {
		return []int{}
	}

	bestCost := math.MaxFloat64
	var bestOrder []int

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}

	var generate func(k int)
	generate = func(k int) {
		if k == n {
			cost := startToWaypointCost(startLat, startLng, waypoints[order[0]].Lat, waypoints[order[0]].Lng, optimizationType)
			for i := 0; i < n-1; i++ {
				cost += costMatrix[order[i]][order[i+1]]
			}
			if cost < bestCost {
				bestCost = cost
				bestOrder = make([]int, n)
				copy(bestOrder, order)
			}
			return
		}
		for i := k; i < n; i++ {
			order[k], order[i] = order[i], order[k]
			generate(k + 1)
			order[k], order[i] = order[i], order[k]
		}
	}

	generate(0)
	return bestOrder
}
