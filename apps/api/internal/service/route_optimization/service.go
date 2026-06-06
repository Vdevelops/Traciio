package route_optimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/datatypes"
)

var (
	ErrRouteNotFound         = errors.New("route not found")
	ErrInvalidWaypoints      = errors.New("invalid waypoints")
	ErrWaypointsTooFew       = errors.New("too few waypoints (minimum 1 required)")
	ErrWaypointsTooMany      = errors.New("too many waypoints (maximum 25 allowed)")
	ErrInvalidCoordinates    = errors.New("invalid coordinates")
	ErrOptimizationFailed    = errors.New("failed to optimize route")
	ErrStartLocationRequired = errors.New("start location is required")
)

type Service struct {
	routeRepo            interfaces.RouteOptimizationRepository
	userRepo             interfaces.UserRepository
	osrmClient           *OSRMClient
	distanceMatrixCache  *DistanceMatrixCache
	routeResultCache     *RouteResultCache
	routeSegmentCache    *RouteSegmentCache
	performanceOptimizer *PerformanceOptimizer
}

func NewService(routeRepo interfaces.RouteOptimizationRepository, userRepo interfaces.UserRepository) *Service {
	osrmBaseURL := config.AppConfig.OSRM.BaseURL
	if osrmBaseURL == "" {
		osrmBaseURL = "https://router.project-osrm.org" // Default public OSRM instance
	}

	distanceMatrixCache := NewDistanceMatrixCache(24 * time.Hour) // Cache for 24 hours
	routeResultCache := NewRouteResultCache(1 * time.Hour)        // Cache for 1 hour
	routeSegmentCache := NewRouteSegmentCache(12 * time.Hour)     // Cache route segments for 12 hours

	performanceOptimizer := NewPerformanceOptimizer(distanceMatrixCache, routeSegmentCache)

	return &Service{
		routeRepo:            routeRepo,
		userRepo:             userRepo,
		osrmClient:           NewOSRMClientWithTimeout(osrmBaseURL, 10*time.Second), // 10 s: fast fail for cross-sea full-route calls
		distanceMatrixCache:  distanceMatrixCache,
		routeResultCache:     routeResultCache,
		routeSegmentCache:    routeSegmentCache,
		performanceOptimizer: performanceOptimizer,
	}
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// haversineDistance calculates the distance between two points using Haversine formula
// Returns distance in kilometers
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

// nearestNeighborOptimizeFromStart optimizes route starting from a specific location
// Returns optimized order (indices of destinations) and total distance
func nearestNeighborOptimizeFromStart(startLat, startLng float64, destinations []route_optimization.Waypoint) ([]int, float64, []route_optimization.RouteStep) {
	if len(destinations) == 0 {
		return []int{}, 0, []route_optimization.RouteStep{}
	}

	n := len(destinations)
	visited := make([]bool, n)
	order := make([]int, 0, n)
	steps := make([]route_optimization.RouteStep, 0, n+1)
	totalDistance := 0.0

	currentLat := startLat
	currentLng := startLng
	stepNumber := 1

	// Find nearest neighbor for each destination
	for len(order) < n {
		nearestIndex := -1
		minDistance := math.MaxFloat64

		// Find nearest unvisited destination
		for i := 0; i < n; i++ {
			if !visited[i] {
				distance := haversineDistance(
					currentLat, currentLng,
					destinations[i].Lat, destinations[i].Lng,
				)
				if distance < minDistance {
					minDistance = distance
					nearestIndex = i
				}
			}
		}

		if nearestIndex >= 0 {
			// Calculate estimated duration (40 km/h average speed)
			avgSpeedKmh := 40.0
			estimatedDuration := int((minDistance / avgSpeedKmh) * 3600) // seconds

			// Create route step
			step := route_optimization.RouteStep{
				Step:              stepNumber,
				Distance:          minDistance,
				DistanceFormatted: formatDistance(minDistance),
				Duration:          estimatedDuration,
				DurationFormatted: formatDuration(estimatedDuration),
				Instruction:       fmt.Sprintf("Drive to %s", destinations[nearestIndex].Address),
				Maneuver:          "straight",
				StartLocation: route_optimization.Location{
					Lat: currentLat,
					Lng: currentLng,
				},
				EndLocation: route_optimization.Location{
					Lat: destinations[nearestIndex].Lat,
					Lng: destinations[nearestIndex].Lng,
				},
			}
			steps = append(steps, step)

			order = append(order, nearestIndex)
			visited[nearestIndex] = true
			totalDistance += minDistance

			// Update current position
			currentLat = destinations[nearestIndex].Lat
			currentLng = destinations[nearestIndex].Lng
			stepNumber++
		}
	}

	return order, totalDistance, steps
}

// formatDistance formats distance to string
func formatDistance(km float64) string {
	if km < 1 {
		return fmt.Sprintf("%.0f m", km*1000)
	}
	return fmt.Sprintf("%.2f km", km)
}

// formatDuration formats duration to string
func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%d detik", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		return fmt.Sprintf("%d menit", minutes)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		if minutes > 0 {
			return fmt.Sprintf("%d jam %d menit", hours, minutes)
		}
		return fmt.Sprintf("%d jam", hours)
	}
}

// routeLegsSequentially routes each consecutive waypoint pair individually via OSRM.
// When the full-route OSRM call fails (e.g., cross-island routes), this function
// ensures same-island legs still follow actual roads via OSRM while cross-island
// legResult holds the result of routing a single leg
type legResult struct {
	index          int
	distance       float64
	duration       int
	steps          []route_optimization.RouteStep
	polylinePoints [][2]float64
	err            error
}

// routeLegsParallel routes legs in parallel using a semaphore to limit concurrent requests.
// This fixes the bottleneck where sequential calls could take 120s+ for 25 waypoints.
func (s *Service) routeLegsParallel(
	startLat, startLng float64,
	orderedWaypoints []route_optimization.Waypoint,
) (float64, int, []route_optimization.RouteStep, *string) {
	type coord struct{ lat, lng float64 }

	// Short-timeout OSRM client for per-leg calls to prevent 504 timeouts.
	legClient := NewOSRMClientWithTimeout(s.osrmClient.baseURL, 5*time.Second)

	// Build ordered coordinate list: start location + each waypoint in optimized order
	coords := make([]coord, 0, len(orderedWaypoints)+1)
	coords = append(coords, coord{startLat, startLng})
	for _, wp := range orderedWaypoints {
		coords = append(coords, coord{wp.Lat, wp.Lng})
	}

	numLegs := len(coords) - 1
	if numLegs <= 0 {
		return 0, 0, nil, nil
	}

	// Use semaphore to limit concurrent requests (max 5 concurrent OSRM calls)
	const maxConcurrent = 5
	semaphore := make(chan struct{}, maxConcurrent)
	results := make([]legResult, numLegs)
	var wg sync.WaitGroup

	// Process legs in parallel with panic recovery
	for i := 0; i < numLegs; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// CRITICAL: Panic recovery to prevent one failed leg from crashing entire request
			defer func() {
				if r := recover(); r != nil {
					log.Printf("PANIC in routeLegsParallel goroutine %d: %v", idx, r)
					// Set fallback result on panic
					from := coords[idx]
					to := coords[idx+1]
					legStraightLine := haversineDistance(from.lat, from.lng, to.lat, to.lng)
					dur := int((legStraightLine / 40.0) * 3600)
					results[idx] = legResult{
						index:    idx,
						distance: legStraightLine,
						duration: dur,
						steps: []route_optimization.RouteStep{{
							Step:              1,
							Distance:          legStraightLine,
							DistanceFormatted: formatDistance(legStraightLine),
							Duration:          dur,
							DurationFormatted: formatDuration(dur),
							Instruction:       "Menuju tujuan (lintas perairan)",
							Maneuver:          "straight",
							StartLocation:     route_optimization.Location{Lat: from.lat, Lng: from.lng},
							EndLocation:       route_optimization.Location{Lat: to.lat, Lng: to.lng},
						}},
						polylinePoints: [][2]float64{
							{from.lat, from.lng},
							{to.lat, to.lng},
						},
					}
				}
			}()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			from := coords[idx]
			to := coords[idx+1]

			// Pre-compute straight-line distance for fallback
			legStraightLine := haversineDistance(from.lat, from.lng, to.lat, to.lng)

			legReq := RouteRequest{
				Coordinates: [][]float64{
					{from.lng, from.lat},
					{to.lng, to.lat},
				},
				Profile:      "driving",
				Overview:     "full",
				Geometries:   "polyline",
				Steps:        true,
				Alternatives: false,
			}

			legResp, err := legClient.Route(legReq)
			if err != nil || len(legResp.Routes) == 0 {
				// OSRM failed - use fallback
				dur := int((legStraightLine / 40.0) * 3600)

				results[idx] = legResult{
					index:    idx,
					distance: legStraightLine,
					duration: dur,
					steps: []route_optimization.RouteStep{{
						Step:              1, // Will be renumbered later
						Distance:          legStraightLine,
						DistanceFormatted: formatDistance(legStraightLine),
						Duration:          dur,
						DurationFormatted: formatDuration(dur),
						Instruction:       "Menuju tujuan (lintas perairan)",
						Maneuver:          "straight",
						StartLocation:     route_optimization.Location{Lat: from.lat, Lng: from.lng},
						EndLocation:       route_optimization.Location{Lat: to.lat, Lng: to.lng},
					}},
					polylinePoints: [][2]float64{
						{from.lat, from.lng},
						{to.lat, to.lng},
					},
				}
				return
			}

			// OSRM success
			osrmLeg := legResp.Routes[0]

			// Decode polyline
			var polylinePoints [][2]float64
			if osrmLeg.Geometry != "" {
				polylinePoints = DecodePolylineToPoints(osrmLeg.Geometry)
			} else {
				polylinePoints = [][2]float64{
					{from.lat, from.lng},
					{to.lat, to.lng},
				}
			}

			// Convert steps
			legWaypoints := []route_optimization.Waypoint{
				{Lat: from.lat, Lng: from.lng},
				{Lat: to.lat, Lng: to.lng},
			}
			legSteps := ConvertOSRMRouteToRouteSteps(&osrmLeg, legWaypoints)

			results[idx] = legResult{
				index:          idx,
				distance:       osrmLeg.Distance / 1000.0,
				duration:       int(osrmLeg.Duration),
				steps:          legSteps,
				polylinePoints: polylinePoints,
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Combine results in order
	totalDistance := 0.0
	totalDuration := 0
	allSteps := make([]route_optimization.RouteStep, 0)
	allPolylinePoints := make([][2]float64, 0)
	stepNumber := 1

	for i := 0; i < numLegs; i++ {
		result := results[i]

		totalDistance += result.distance
		totalDuration += result.duration

		// Append polyline points (skip first point of subsequent legs to avoid duplication)
		if len(allPolylinePoints) > 0 && len(result.polylinePoints) > 1 {
			allPolylinePoints = append(allPolylinePoints, result.polylinePoints[1:]...)
		} else {
			allPolylinePoints = append(allPolylinePoints, result.polylinePoints...)
		}

		// Renumber steps
		for _, step := range result.steps {
			step.Step = stepNumber
			allSteps = append(allSteps, step)
			stepNumber++
		}
	}

	// Re-encode all collected coordinate points back into a single polyline string
	var combinedPolyline *string
	if len(allPolylinePoints) > 1 {
		encoded := EncodePolylineFromPoints(allPolylinePoints)
		combinedPolyline = &encoded
	}

	return totalDistance, totalDuration, allSteps, combinedPolyline
}

// routeLegsSequentially routes each leg sequentially using OSRM.
// NOTE: Deprecated - use routeLegsParallel for better performance under load.
// OSRM returns NoRoute almost instantly when a leg is truly unreachable, so there
// is no pre-emptive distance guard — every leg is attempted and falls back on failure.
func (s *Service) routeLegsSequentially(
	startLat, startLng float64,
	orderedWaypoints []route_optimization.Waypoint,
) (float64, int, []route_optimization.RouteStep, *string) {
	type coord struct{ lat, lng float64 }

	// Short-timeout OSRM client for per-leg calls to prevent 504 timeouts.
	legClient := NewOSRMClientWithTimeout(s.osrmClient.baseURL, 5*time.Second)

	// Build ordered coordinate list: start location + each waypoint in optimized order
	coords := make([]coord, 0, len(orderedWaypoints)+1)
	coords = append(coords, coord{startLat, startLng})
	for _, wp := range orderedWaypoints {
		coords = append(coords, coord{wp.Lat, wp.Lng})
	}

	totalDistance := 0.0
	totalDuration := 0
	allSteps := make([]route_optimization.RouteStep, 0)
	allPolylinePoints := make([][2]float64, 0)
	stepNumber := 1

	for i := 0; i < len(coords)-1; i++ {
		from := coords[i]
		to := coords[i+1]

		// Pre-compute straight-line distance — used only as the fallback value when
		// OSRM cannot route this leg (NoRoute or timeout within 5 s).
		legStraightLine := haversineDistance(from.lat, from.lng, to.lat, to.lng)

		legReq := RouteRequest{
			Coordinates: [][]float64{
				{from.lng, from.lat}, // OSRM uses lng,lat order
				{to.lng, to.lat},
			},
			Profile:      "driving",
			Overview:     "full",
			Geometries:   "polyline",
			Steps:        true,
			Alternatives: false,
		}

		legResp, err := legClient.Route(legReq)
		if err != nil || len(legResp.Routes) == 0 {
			// This leg cannot be routed by OSRM (e.g., cross-island without ferry data).
			// Reuse the already-computed straight-line distance for this segment.
			dur := int((legStraightLine / 40.0) * 3600)
			totalDistance += legStraightLine
			totalDuration += dur

			allPolylinePoints = append(allPolylinePoints,
				[2]float64{from.lat, from.lng},
				[2]float64{to.lat, to.lng},
			)
			allSteps = append(allSteps, route_optimization.RouteStep{
				Step:              stepNumber,
				Distance:          legStraightLine,
				DistanceFormatted: formatDistance(legStraightLine),
				Duration:          dur,
				DurationFormatted: formatDuration(dur),
				Instruction:       "Menuju tujuan (lintas perairan)",
				Maneuver:          "straight",
				StartLocation:     route_optimization.Location{Lat: from.lat, Lng: from.lng},
				EndLocation:       route_optimization.Location{Lat: to.lat, Lng: to.lng},
			})
			stepNumber++
			continue
		}

		osrmLeg := legResp.Routes[0]
		totalDistance += osrmLeg.Distance / 1000.0
		totalDuration += int(osrmLeg.Duration)

		// Decode and append road-following polyline for this leg
		if osrmLeg.Geometry != "" {
			legPoints := DecodePolylineToPoints(osrmLeg.Geometry)
			if len(allPolylinePoints) > 0 && len(legPoints) > 1 {
				// Skip first point of this leg to avoid duplication with previous leg's last point
				legPoints = legPoints[1:]
			}
			allPolylinePoints = append(allPolylinePoints, legPoints...)
		} else {
			allPolylinePoints = append(allPolylinePoints,
				[2]float64{from.lat, from.lng},
				[2]float64{to.lat, to.lng},
			)
		}

		// Convert OSRM leg steps and renumber sequentially across all legs
		legWaypoints := []route_optimization.Waypoint{
			{Lat: from.lat, Lng: from.lng},
			{Lat: to.lat, Lng: to.lng},
		}
		legSteps := ConvertOSRMRouteToRouteSteps(&osrmLeg, legWaypoints)
		for _, step := range legSteps {
			step.Step = stepNumber
			allSteps = append(allSteps, step)
			stepNumber++
		}
	}

	// Re-encode all collected coordinate points back into a single polyline string
	var combinedPolyline *string
	if len(allPolylinePoints) > 1 {
		encoded := EncodePolylineFromPoints(allPolylinePoints)
		combinedPolyline = &encoded
	}

	return totalDistance, totalDuration, allSteps, combinedPolyline
}

// optimizeRouteWithOSRM optimizes route order and calculates actual routing using OSRM
// Now supports caching, parallel requests, and time windows
func (s *Service) optimizeRouteWithOSRM(startLat, startLng float64, waypoints []route_optimization.Waypoint, startTime *time.Time, optimizationType string) ([]int, float64, int, []route_optimization.RouteStep, *string, error) {
	optimizationType = normalizeOptimizationMode(optimizationType)
	var optimizedOrder []int

	// Check cache first (only for non-time-window routes)
	waypointHash := HashWaypoints(waypoints)
	cacheKey := CreateRouteCacheKeyWithExtras(startLat, startLng, waypointHash, "opt="+optimizationType)

	// Check if time windows are needed
	hasTimeWindows := false
	if startTime != nil {
		for _, wp := range waypoints {
			if wp.EarliestArrival != nil || wp.LatestArrival != nil {
				hasTimeWindows = true
				break
			}
		}
	}

	// Only use cache if no time windows (time windows are dynamic)
	if !hasTimeWindows {
		if cached, found := s.routeResultCache.Get(cacheKey); found {
			routeOptimizationCacheEventsTotal.WithLabelValues("final_route", "hit").Inc()
			// Cache hit!
			// If we have the final route cached (steps/polyline), we can skip OSRM entirely.
			if len(cached.RouteSteps) > 0 {
				routeOptimizationOSRMSkippedTotal.Inc()
				return cached.OptimizedOrder, cached.TotalDistance, cached.TotalDuration, cached.RouteSteps, cached.RoutePolyline, nil
			}
			// Backward-compat: older cache entries may only contain the optimized order.
			optimizedOrder = cached.OptimizedOrder
		}
		if optimizedOrder == nil {
			routeOptimizationCacheEventsTotal.WithLabelValues("final_route", "miss").Inc()
		}
	} else {
		routeOptimizationCacheEventsTotal.WithLabelValues("final_route", "ineligible_time_windows").Inc()
	}

	// If not from cache, optimize
	if optimizedOrder == nil {
		if hasTimeWindows && startTime != nil {
			// Use time windows optimization
			// First, calculate distance matrix in parallel (with caching)
			distanceMatrix, durationMatrix, err := s.CalculateDistanceMatrixParallel(waypoints, s.distanceMatrixCache)
			if err != nil {
				// Fallback to simple optimization
				initialOrder, _, _ := nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
				if len(waypoints) >= 4 {
					optimizedOrder, _ = TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
				} else {
					optimizedOrder = initialOrder
				}
			} else {
				// Optimize with time windows
				var arrivalTimes []time.Time
				optimizedOrder, arrivalTimes, err = OptimizeWithTimeWindows(
					startLat, startLng, waypoints, *startTime,
					distanceMatrix, durationMatrix,
				)
				if err != nil {
					// Fallback to simple optimization
					initialOrder, _, _ := nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
					if len(waypoints) >= 4 {
						optimizedOrder, _ = TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
					} else {
						optimizedOrder = initialOrder
					}
				}
				_ = arrivalTimes // Can be used for ETA display
			}
		} else {
			// Regular optimization without time windows
			// Use performance optimizer for faster processing
			var optErr error
			optimizedOrder, optErr = s.performanceOptimizer.OptimizeRouteFast(
				s.osrmClient,
				startLat, startLng,
				waypoints,
				optimizationType,
			)
			if optErr != nil {
				// Fallback to regular optimization
				initialOrder, _, _ := nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
				if len(waypoints) >= 4 {
					optimizedOrder, _ = TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
				} else {
					optimizedOrder = initialOrder
				}
			}
		}
	}

	// Step 3: Build coordinates array for OSRM (include start location + optimized waypoints)
	coordinates := make([][]float64, 0, len(waypoints)+1)

	// Add start location (OSRM uses lng,lat format)
	coordinates = append(coordinates, []float64{startLng, startLat})

	// Add waypoints in optimized order
	for _, idx := range optimizedOrder {
		wp := waypoints[idx]
		coordinates = append(coordinates, []float64{wp.Lng, wp.Lat})
	}

	// Step 3: Request route from OSRM
	routeReq := RouteRequest{
		Coordinates:  coordinates,
		Profile:      "driving",
		Overview:     "full",
		Geometries:   "polyline",
		Steps:        true,
		Alternatives: false,
	}

	osrmStart := time.Now()
	osrmResp, err := s.osrmClient.Route(routeReq)
	routeOptimizationOSRMRouteDurationSeconds.Observe(time.Since(osrmStart).Seconds())
	if err != nil {
		// Full-route OSRM failed (e.g., cross-island destinations with no ferry data).
		// Try leg-by-leg routing: each consecutive pair is routed individually so that
		// same-island legs still follow roads and cross-island legs degrade gracefully.
		orderedWps := make([]route_optimization.Waypoint, len(optimizedOrder))
		for i, idx := range optimizedOrder {
			orderedWps[i] = waypoints[idx]
		}
		// FIXED: Use parallel routing instead of sequential for better performance under load
		legDist, legDur, legSteps, legPolyline := s.routeLegsParallel(startLat, startLng, orderedWps)
		return optimizedOrder, legDist, legDur, legSteps, legPolyline, nil
	}

	// Step 4: Extract route data from OSRM response
	if len(osrmResp.Routes) == 0 {
		// OSRM returned no routes (NoRoute code). Try leg-by-leg routing so that
		// segments between islands follow roads as much as possible.
		orderedWps := make([]route_optimization.Waypoint, len(optimizedOrder))
		for i, idx := range optimizedOrder {
			orderedWps[i] = waypoints[idx]
		}
		// FIXED: Use parallel routing instead of sequential for better performance under load
		legDist, legDur, legSteps, legPolyline := s.routeLegsParallel(startLat, startLng, orderedWps)
		return optimizedOrder, legDist, legDur, legSteps, legPolyline, nil
	}

	osrmRoute := osrmResp.Routes[0]

	// Convert distance from meters to kilometers
	totalDistance := osrmRoute.Distance / 1000.0
	totalDuration := int(osrmRoute.Duration)

	// Prepare waypoints in optimized order for step conversion
	orderedWaypoints := make([]route_optimization.Waypoint, 0, len(waypoints)+1)
	orderedWaypoints = append(orderedWaypoints, route_optimization.Waypoint{
		Lat: startLat,
		Lng: startLng,
	})
	for _, idx := range optimizedOrder {
		orderedWaypoints = append(orderedWaypoints, waypoints[idx])
	}

	// Convert OSRM route to RouteSteps
	routeSteps := ConvertOSRMRouteToRouteSteps(&osrmRoute, orderedWaypoints)

	// Extract polyline
	var polyline *string
	if osrmRoute.Geometry != "" {
		polyline = &osrmRoute.Geometry
	}

	// Cache the result (only if no time windows)
	if !hasTimeWindows {
		s.routeResultCache.Set(cacheKey, optimizedOrder, totalDistance, totalDuration, routeSteps, polyline)
	}

	return optimizedOrder, totalDistance, totalDuration, routeSteps, polyline, nil
}

// Optimize optimizes a route using TSP algorithm and OSRM routing starting from user's current location
func (s *Service) Optimize(req *route_optimization.OptimizeRouteRequest, userID string) (*route_optimization.OptimizedRouteResponse, error) {
	start := time.Now()
	defer func() {
		routeOptimizationOptimizeDurationSeconds.Observe(time.Since(start).Seconds())
	}()

	// Validate start location (current user location)
	if req.StartLocation == nil {
		return nil, ErrStartLocationRequired
	}
	if req.StartLocation.Lat < -90 || req.StartLocation.Lat > 90 ||
		req.StartLocation.Lng < -180 || req.StartLocation.Lng > 180 {
		return nil, fmt.Errorf("%w: start location has invalid coordinates", ErrInvalidCoordinates)
	}

	// Validate waypoints (destinations)
	if len(req.Waypoints) < 1 {
		return nil, ErrWaypointsTooFew
	}
	if len(req.Waypoints) > 25 {
		return nil, ErrWaypointsTooMany
	}

	// Validate destination coordinates
	for i, wp := range req.Waypoints {
		if wp.Lat < -90 || wp.Lat > 90 || wp.Lng < -180 || wp.Lng > 180 {
			return nil, fmt.Errorf("%w: waypoint %d has invalid coordinates", ErrInvalidCoordinates, i)
		}
	}

	// Optimize route using OSRM (includes TSP optimization + actual routing)
	// Now supports time windows if startTime is provided
	optimizedOrder, totalDistance, totalDuration, routeSteps, routePolyline, err := s.optimizeRouteWithOSRM(
		req.StartLocation.Lat,
		req.StartLocation.Lng,
		req.Waypoints,
		req.StartTime,
		routeOptimizationModeAuto,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOptimizationFailed, err)
	}

	// Reorder waypoints according to optimized order
	optimizedWaypoints := make([]route_optimization.Waypoint, len(req.Waypoints))
	for i, idx := range optimizedOrder {
		wp := req.Waypoints[idx]
		wp.Order = i + 1
		optimizedWaypoints[i] = wp
	}

	// Generate route name if not provided
	routeName := req.RouteName
	if routeName == nil || *routeName == "" {
		name := fmt.Sprintf("Optimized Route - %s", time.Now().Format("2006-01-02"))
		routeName = &name
	}

	// Prepare waypoints JSON (include start location as first waypoint)
	allWaypoints := append([]route_optimization.Waypoint{
		{
			Order:       0,
			Lat:         req.StartLocation.Lat,
			Lng:         req.StartLocation.Lng,
			Address:     req.StartLocation.Address,
			AccountID:   nil,
			AccountName: nil,
			ContactID:   nil,
			ContactName: nil,
		},
	}, optimizedWaypoints...)

	waypointsJSON, err := json.Marshal(allWaypoints)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal waypoints", ErrOptimizationFailed)
	}

	// Prepare optimized order JSON
	orderJSON, err := json.Marshal(optimizedOrder)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal optimized order", ErrOptimizationFailed)
	}

	// Prepare route steps JSON
	stepsJSON, err := json.Marshal(routeSteps)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal route steps", ErrOptimizationFailed)
	}

	// Create route entity
	route := &route_optimization.OptimizedRoute{
		UserID:         userID,
		RouteName:      routeName,
		Waypoints:      datatypes.JSON(waypointsJSON),
		OptimizedOrder: datatypes.JSON(orderJSON),
		TotalDistance:  &totalDistance,
		TotalDuration:  &totalDuration,
		RoutePolyline:  routePolyline, // Set polyline from OSRM
		RouteSteps:     datatypes.JSON(stepsJSON),
	}

	// Save to database
	if err := s.routeRepo.Create(route); err != nil {
		return nil, fmt.Errorf("%w: failed to save route", ErrOptimizationFailed)
	}

	// Load user relation
	if user, err := s.userRepo.FindByID(userID); err == nil {
		route.User = map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		}
	}

	// Convert to response
	response := s.toOptimizedRouteResponse(route, allWaypoints)
	return response, nil
}

// GetByID returns an optimized route by ID
func (s *Service) GetByID(id string) (*route_optimization.OptimizedRouteResponse, error) {
	route, err := s.routeRepo.FindByID(id)
	if err != nil {
		return nil, ErrRouteNotFound
	}

	// Parse waypoints
	var waypoints []route_optimization.Waypoint
	if route.Waypoints != nil {
		if err := json.Unmarshal(route.Waypoints, &waypoints); err != nil {
			return nil, fmt.Errorf("failed to parse waypoints: %w", err)
		}
	}

	// Load user relation
	if user, err := s.userRepo.FindByID(route.UserID); err == nil {
		route.User = map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		}
	}

	return s.toOptimizedRouteResponse(route, waypoints), nil
}

// List returns a list of optimized routes with pagination
func (s *Service) List(req *route_optimization.ListRoutesRequest) ([]route_optimization.OptimizedRouteResponse, *PaginationResult, error) {
	routes, total, err := s.routeRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]route_optimization.OptimizedRouteResponse, len(routes))
	for i, route := range routes {
		// Parse waypoints
		var waypoints []route_optimization.Waypoint
		if route.Waypoints != nil {
			if err := json.Unmarshal(route.Waypoints, &waypoints); err == nil {
				// Load user relation
				if user, err := s.userRepo.FindByID(route.UserID); err == nil {
					route.User = map[string]interface{}{
						"id":   user.ID,
						"name": user.Name,
					}
				}
				responses[i] = *s.toOptimizedRouteResponse(&route, waypoints)
			}
		}
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	totalPages := (int(total) + perPage - 1) / perPage

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return responses, pagination, nil
}

// CalculateDistance calculates distance between two points using OSRM routing
func (s *Service) CalculateDistance(req *route_optimization.CalculateDistanceRequest) (*route_optimization.CalculateDistanceResponse, error) {
	// Validate coordinates
	if req.Origin.Lat < -90 || req.Origin.Lat > 90 || req.Origin.Lng < -180 || req.Origin.Lng > 180 {
		return nil, fmt.Errorf("%w: origin has invalid coordinates", ErrInvalidCoordinates)
	}
	if req.Destination.Lat < -90 || req.Destination.Lat > 90 || req.Destination.Lng < -180 || req.Destination.Lng > 180 {
		return nil, fmt.Errorf("%w: destination has invalid coordinates", ErrInvalidCoordinates)
	}

	// Build coordinates for OSRM (OSRM uses lng,lat format)
	coordinates := [][]float64{
		{req.Origin.Lng, req.Origin.Lat},
		{req.Destination.Lng, req.Destination.Lat},
	}

	// Request route from OSRM
	routeReq := RouteRequest{
		Coordinates:  coordinates,
		Profile:      "driving",
		Overview:     "simplified",
		Geometries:   "polyline",
		Steps:        false,
		Alternatives: false,
	}

	osrmResp, err := s.osrmClient.Route(routeReq)
	if err != nil {
		// Fallback to Haversine formula if OSRM fails
		distance := haversineDistance(
			req.Origin.Lat, req.Origin.Lng,
			req.Destination.Lat, req.Destination.Lng,
		)
		avgSpeedKmh := 40.0
		estimatedDuration := int((distance / avgSpeedKmh) * 3600)

		return &route_optimization.CalculateDistanceResponse{
			Distance:          distance,
			DistanceFormatted: formatDistance(distance),
			Duration:          estimatedDuration,
			DurationFormatted: formatDuration(estimatedDuration),
		}, nil
	}

	// Extract distance and duration from OSRM response
	if len(osrmResp.Routes) == 0 {
		// Fallback to Haversine formula
		distance := haversineDistance(
			req.Origin.Lat, req.Origin.Lng,
			req.Destination.Lat, req.Destination.Lng,
		)
		avgSpeedKmh := 40.0
		estimatedDuration := int((distance / avgSpeedKmh) * 3600)

		return &route_optimization.CalculateDistanceResponse{
			Distance:          distance,
			DistanceFormatted: formatDistance(distance),
			Duration:          estimatedDuration,
			DurationFormatted: formatDuration(estimatedDuration),
		}, nil
	}

	osrmRoute := osrmResp.Routes[0]

	// Convert distance from meters to kilometers
	distance := osrmRoute.Distance / 1000.0
	duration := int(osrmRoute.Duration)

	return &route_optimization.CalculateDistanceResponse{
		Distance:          distance,
		DistanceFormatted: formatDistance(distance),
		Duration:          duration,
		DurationFormatted: formatDuration(duration),
	}, nil
}

// Delete deletes an optimized route
func (s *Service) Delete(id string) error {
	_, err := s.routeRepo.FindByID(id)
	if err != nil {
		return ErrRouteNotFound
	}

	return s.routeRepo.Delete(id)
}

// toOptimizedRouteResponse converts OptimizedRoute to OptimizedRouteResponse
func (s *Service) toOptimizedRouteResponse(route *route_optimization.OptimizedRoute, waypoints []route_optimization.Waypoint) *route_optimization.OptimizedRouteResponse {
	response := &route_optimization.OptimizedRouteResponse{
		ID:        route.ID,
		RouteName: route.RouteName,
		UserID:    route.UserID,
		Waypoints: waypoints,
		CreatedAt: route.CreatedAt,
		UpdatedAt: route.UpdatedAt,
		User:      route.User,
	}

	// Parse optimized order
	if route.OptimizedOrder != nil {
		json.Unmarshal(route.OptimizedOrder, &response.OptimizedOrder)
	}

	// Set distance
	if route.TotalDistance != nil {
		response.TotalDistance = route.TotalDistance
		response.TotalDistanceFormatted = formatDistance(*route.TotalDistance)
	}

	// Set duration
	if route.TotalDuration != nil {
		response.TotalDuration = route.TotalDuration
		response.TotalDurationFormatted = formatDuration(*route.TotalDuration)
	}

	// Set polyline
	if route.RoutePolyline != nil {
		response.RoutePolyline = route.RoutePolyline
	}

	// Parse route steps
	if route.RouteSteps != nil {
		var steps []route_optimization.RouteStep
		if err := json.Unmarshal(route.RouteSteps, &steps); err == nil {
			response.RouteSteps = steps
		}
	}

	return response
}
