package route_optimization

import (
	"sort"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
)

// TimeWindowOptimizer optimizes routes considering time window constraints
type TimeWindowOptimizer struct {
	startTime time.Time
}

// OptimizeWithTimeWindows optimizes route order considering time windows and priorities
// Algorithm:
// 1. Sort waypoints by priority (highest first)
// 2. Try to fit waypoints within their time windows
// 3. If time window conflict, prioritize higher priority waypoints
// 4. Use 2-Opt to improve route after time windows are satisfied
func OptimizeWithTimeWindows(
	startLat, startLng float64,
	waypoints []route_optimization.Waypoint,
	startTime time.Time,
	distanceMatrix [][]float64,
	durationMatrix [][]int,
) ([]int, []time.Time, error) {
	if len(waypoints) == 0 {
		return []int{}, []time.Time{}, nil
	}

	// If no time windows specified, use regular optimization
	hasTimeWindows := false
	for _, wp := range waypoints {
		if wp.EarliestArrival != nil || wp.LatestArrival != nil {
			hasTimeWindows = true
			break
		}
	}

	if !hasTimeWindows {
		// No time windows, use regular 2-Opt optimization
		initialOrder, _, _ := nearestNeighborOptimizeFromStart(startLat, startLng, waypoints)
		if len(waypoints) >= 4 {
			optimizedOrder, _ := TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
			return optimizedOrder, nil, nil
		}
		return initialOrder, nil, nil
	}

	// Create waypoint with priority and time window info
	type WaypointInfo struct {
		Index          int
		Waypoint       route_optimization.Waypoint
		Priority       int
		EarliestArrival *time.Time
		LatestArrival   *time.Time
		ServiceDuration int // in minutes
	}

	waypointInfos := make([]WaypointInfo, len(waypoints))
	for i, wp := range waypoints {
		priority := 3 // Default priority
		if wp.Priority != nil {
			priority = *wp.Priority
		}

		serviceDuration := 30 // Default 30 minutes
		if wp.ServiceDuration != nil {
			serviceDuration = *wp.ServiceDuration
		}

		waypointInfos[i] = WaypointInfo{
			Index:           i,
			Waypoint:        wp,
			Priority:        priority,
			EarliestArrival: wp.EarliestArrival,
			LatestArrival:   wp.LatestArrival,
			ServiceDuration: serviceDuration,
		}
	}

	// Sort by priority (highest priority first)
	sort.Slice(waypointInfos, func(i, j int) bool {
		return waypointInfos[i].Priority < waypointInfos[j].Priority
	})

	// Build route considering time windows
	order := make([]int, 0, len(waypoints))
	arrivalTimes := make([]time.Time, 0, len(waypoints))
	currentTime := startTime
	visited := make([]bool, len(waypoints))

	// First, try to schedule high-priority waypoints within their time windows
	for _, info := range waypointInfos {
		if visited[info.Index] {
			continue
		}

		// Calculate travel time from current position
		var travelDuration time.Duration
		if len(order) == 0 {
			// First waypoint: calculate from start location
			distance := haversineDistance(startLat, startLng, info.Waypoint.Lat, info.Waypoint.Lng)
			avgSpeedKmh := 40.0
			travelDuration = time.Duration((distance/avgSpeedKmh)*3600) * time.Second
		} else {
			// Use distance matrix if available
			lastIndex := order[len(order)-1]
			if distanceMatrix != nil && len(distanceMatrix) > lastIndex && len(distanceMatrix[lastIndex]) > info.Index {
				avgSpeedKmh := 40.0
				distance := distanceMatrix[lastIndex][info.Index]
				travelDuration = time.Duration((distance/avgSpeedKmh)*3600) * time.Second
			} else {
				// Fallback: calculate from last waypoint
				lastWp := waypoints[lastIndex]
				distance := haversineDistance(lastWp.Lat, lastWp.Lng, info.Waypoint.Lat, info.Waypoint.Lng)
				avgSpeedKmh := 40.0
				travelDuration = time.Duration((distance/avgSpeedKmh)*3600) * time.Second
			}
		}

		// Calculate arrival time
		arrivalTime := currentTime.Add(travelDuration)

		// Check time window constraints
		canVisit := true
		if info.EarliestArrival != nil && arrivalTime.Before(*info.EarliestArrival) {
			// Arrive too early, wait until earliest arrival
			arrivalTime = *info.EarliestArrival
		}
		if info.LatestArrival != nil && arrivalTime.After(*info.LatestArrival) {
			// Cannot visit within time window
			// For now, we'll still add it but mark as constraint violation
			// In production, you might want to skip or reschedule
			canVisit = false
		}

		if canVisit {
			order = append(order, info.Index)
			arrivalTimes = append(arrivalTimes, arrivalTime)
			visited[info.Index] = true

			// Update current time
			currentTime = arrivalTime.Add(time.Duration(info.ServiceDuration) * time.Minute)
		}
	}

	// Add remaining waypoints (those without time windows or that couldn't be scheduled)
	for i, wp := range waypoints {
		if !visited[i] {
			order = append(order, i)
			// Estimate arrival time
			if len(arrivalTimes) > 0 && len(order) > 0 {
				lastArrival := arrivalTimes[len(arrivalTimes)-1]
				lastIndex := order[len(order)-1]
				lastWp := waypoints[lastIndex]
				serviceDuration := 30
				if wp.ServiceDuration != nil {
					serviceDuration = *wp.ServiceDuration
				}
				distance := haversineDistance(lastWp.Lat, lastWp.Lng, wp.Lat, wp.Lng)
				avgSpeedKmh := 40.0
				travelTime := time.Duration((distance/avgSpeedKmh)*3600) * time.Second
				arrivalTime := lastArrival.Add(time.Duration(serviceDuration) * time.Minute).Add(travelTime)
				arrivalTimes = append(arrivalTimes, arrivalTime)
			} else {
				// First waypoint
				distance := haversineDistance(startLat, startLng, wp.Lat, wp.Lng)
				avgSpeedKmh := 40.0
				travelTime := time.Duration((distance/avgSpeedKmh)*3600) * time.Second
				arrivalTime := startTime.Add(travelTime)
				arrivalTimes = append(arrivalTimes, arrivalTime)
			}
		}
	}

	// Improve route with 2-Opt while trying to maintain time window constraints
	if len(waypoints) >= 4 {
		// Try 2-Opt improvement, but validate time windows after each swap
		improvedOrder, _ := TwoOptOptimize(order, waypoints, startLat, startLng)
		
		// Validate improved order against time windows
		if validateTimeWindows(improvedOrder, waypoints, startTime, distanceMatrix, durationMatrix) {
			order = improvedOrder
		}
	}

	return order, arrivalTimes, nil
}

// validateTimeWindows validates if an order satisfies all time window constraints
func validateTimeWindows(
	order []int,
	waypoints []route_optimization.Waypoint,
	startTime time.Time,
	distanceMatrix [][]float64,
	durationMatrix [][]int,
) bool {
	currentTime := startTime

	for i, idx := range order {
		wp := waypoints[idx]

		// Calculate travel time
		var travelDuration time.Duration
		if i == 0 {
			// First waypoint - would need start location
			// For validation, we'll use approximate
			travelDuration = 30 * time.Minute // Approximate
		} else {
			prevIdx := order[i-1]
			if durationMatrix != nil && len(durationMatrix) > prevIdx && len(durationMatrix[prevIdx]) > idx {
				travelDuration = time.Duration(durationMatrix[prevIdx][idx]) * time.Second
			} else {
				travelDuration = 30 * time.Minute // Approximate
			}
		}

		arrivalTime := currentTime.Add(travelDuration)

		// Check time window
		if wp.EarliestArrival != nil && arrivalTime.Before(*wp.EarliestArrival) {
			// OK, can wait
			arrivalTime = *wp.EarliestArrival
		}
		if wp.LatestArrival != nil && arrivalTime.After(*wp.LatestArrival) {
			// Violates time window
			return false
		}

		// Update current time
		serviceDuration := 30
		if wp.ServiceDuration != nil {
			serviceDuration = *wp.ServiceDuration
		}
		currentTime = arrivalTime.Add(time.Duration(serviceDuration) * time.Minute)
	}

	return true
}

