package route_optimization

import (
	"sync"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
)

// DistanceMatrixResult represents a single distance calculation result
type DistanceMatrixResult struct {
	I        int
	J        int
	Distance float64
	Duration int
	Error    error
}

// CalculateDistanceMatrixParallel calculates distance matrix using parallel OSRM requests
// This is 5-10x faster than sequential requests
func (s *Service) CalculateDistanceMatrixParallel(
	waypoints []route_optimization.Waypoint,
	cache *DistanceMatrixCache,
) ([][]float64, [][]int, error) {
	// If OSRM /table is enabled and supported, prefer it (single request).
	// Otherwise fallback to the existing parallel pairwise route calls.
	methodLabel := "parallel"
	start := time.Now()
	defer func() {
		routeOptimizationMatrixDurationSeconds.WithLabelValues(methodLabel).Observe(time.Since(start).Seconds())
	}()

	n := len(waypoints)
	if n == 0 {
		return nil, nil, nil
	}

	if config.AppConfig != nil && config.AppConfig.OSRM.TableEnabled {
		maxN := config.AppConfig.OSRM.TableMaxWaypoints
		if maxN <= 0 {
			maxN = 25
		}
		if n <= maxN {
			coords := make([][]float64, 0, n)
			for _, wp := range waypoints {
				coords = append(coords, []float64{wp.Lng, wp.Lat})
			}
			methodLabel = "table"
			tableResp, err := s.osrmClient.Table(TableRequest{
				Coordinates: coords,
				Profile:     "driving",
				Annotations: []string{"distance", "duration"},
			})
			if err == nil && len(tableResp.Distances) == n && len(tableResp.Durations) == n {
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
					cache.Set(waypointHash, distanceMatrix)
					return distanceMatrix, durationMatrix, nil
				}
			}
			// Fall back to parallel if table is unavailable/invalid.
			methodLabel = "parallel"
		}
	}

	// Check cache first
	waypointHash := HashWaypoints(waypoints)
	if cachedMatrix, found := cache.Get(waypointHash); found {
		// Cache hit! Return cached matrix
		// Also return duration matrix (estimate from distance)
		durationMatrix := estimateDurationMatrix(cachedMatrix)
		return cachedMatrix, durationMatrix, nil
	}

	// Initialize matrices
	distanceMatrix := make([][]float64, n)
	durationMatrix := make([][]int, n)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, n)
		durationMatrix[i] = make([]int, n)
	}

	// Use worker pool to limit concurrent requests
	const maxWorkers = 10 // Limit concurrent OSRM requests
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	// Calculate distances in parallel
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// Calculate distance using OSRM
				distance, duration, err := s.calculateDistanceBetweenPoints(
					waypoints[i].Lat, waypoints[i].Lng,
					waypoints[j].Lat, waypoints[j].Lng,
				)

				mu.Lock()
				if err != nil && firstErr == nil {
					firstErr = err
				} else if err == nil {
					// Distance matrix is symmetric
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
	cache.Set(waypointHash, distanceMatrix)

	return distanceMatrix, durationMatrix, nil
}

// calculateDistanceBetweenPoints calculates distance and duration between two points using OSRM
func (s *Service) calculateDistanceBetweenPoints(
	lat1, lng1, lat2, lng2 float64,
) (float64, int, error) {
	// Build coordinates for OSRM (OSRM uses lng,lat format)
	coordinates := [][]float64{
		{lng1, lat1},
		{lng2, lat2},
	}

	// Request route from OSRM (simplified, no steps for faster response)
	routeReq := RouteRequest{
		Coordinates:  coordinates,
		Profile:      "driving",
		Overview:     "simplified", // Simplified for faster response
		Geometries:   "polyline",
		Steps:        false, // No steps needed for distance matrix
		Alternatives: false,
	}

	osrmResp, err := s.osrmClient.Route(routeReq)
	if err != nil {
		// Fallback to Haversine formula
		distance := haversineDistance(lat1, lng1, lat2, lng2)
		avgSpeedKmh := 40.0
		duration := int((distance / avgSpeedKmh) * 3600)
		return distance, duration, nil
	}

	if len(osrmResp.Routes) == 0 {
		// Fallback to Haversine formula
		distance := haversineDistance(lat1, lng1, lat2, lng2)
		avgSpeedKmh := 40.0
		duration := int((distance / avgSpeedKmh) * 3600)
		return distance, duration, nil
	}

	osrmRoute := osrmResp.Routes[0]

	// Convert distance from meters to kilometers
	distance := osrmRoute.Distance / 1000.0
	duration := int(osrmRoute.Duration)

	return distance, duration, nil
}

// estimateDurationMatrix estimates duration matrix from distance matrix
// Uses average speed of 40 km/h
func estimateDurationMatrix(distanceMatrix [][]float64) [][]int {
	n := len(distanceMatrix)
	durationMatrix := make([][]int, n)
	avgSpeedKmh := 40.0

	for i := range distanceMatrix {
		durationMatrix[i] = make([]int, n)
		for j := range distanceMatrix[i] {
			// Duration in seconds = (distance in km / speed in km/h) * 3600
			durationMatrix[i][j] = int((distanceMatrix[i][j] / avgSpeedKmh) * 3600)
		}
	}

	return durationMatrix
}
