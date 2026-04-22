package route_optimization

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
)

// TwoOptOptimize improves route using 2-opt algorithm
// This algorithm tries to improve a route by swapping two edges
// It iterates until no improvement can be found
// Expected improvement: 15-25% reduction in total distance
func TwoOptOptimize(
	initialOrder []int,
	waypoints []route_optimization.Waypoint,
	startLat, startLng float64,
) ([]int, float64) {
	if len(initialOrder) < 4 {
		// 2-opt requires at least 4 waypoints to swap edges
		return initialOrder, calculateTotalDistance(initialOrder, waypoints, startLat, startLng)
	}

	bestOrder := make([]int, len(initialOrder))
	copy(bestOrder, initialOrder)
	bestDistance := calculateTotalDistance(bestOrder, waypoints, startLat, startLng)

	improved := true
	tuning := getImprovementTuning()
	maxIterations := tuning.twoOptMaxIter // Prevent infinite loops
	iteration := 0
	deadline, hasDeadline := improvementDeadline(len(bestOrder))

	for improved && iteration < maxIterations {
		if hasDeadline && time.Now().After(deadline) {
			break
		}
		improved = false
		iteration++

		// Try all possible 2-opt swaps
		for i := 0; i < len(bestOrder)-1; i++ {
			if hasDeadline && time.Now().After(deadline) {
				break
			}
			for j := i + 2; j < len(bestOrder); j++ {
				if hasDeadline && time.Now().After(deadline) {
					break
				}
				// Try reversing segment between i and j
				newOrder := reverseSegment(bestOrder, i, j)
				newDistance := calculateTotalDistance(newOrder, waypoints, startLat, startLng)

				if newDistance < bestDistance {
					bestOrder = newOrder
					bestDistance = newDistance
					improved = true
					break // Restart search after improvement
				}
			}
			if improved {
				break
			}
		}
	}

	return bestOrder, bestDistance
}

// reverseSegment reverses a segment of the route between indices i and j
// This is the core operation of 2-opt algorithm
func reverseSegment(order []int, i, j int) []int {
	newOrder := make([]int, len(order))
	copy(newOrder, order)

	// Reverse segment from i+1 to j
	// Example: [0, 1, 2, 3, 4] with i=1, j=3 becomes [0, 3, 2, 1, 4]
	for k := 0; k <= (j-i-1)/2; k++ {
		newOrder[i+1+k], newOrder[j-k] = newOrder[j-k], newOrder[i+1+k]
	}

	return newOrder
}

// calculateTotalDistance calculates total distance for a route order
// Uses Haversine formula for distance calculation
func calculateTotalDistance(
	order []int,
	waypoints []route_optimization.Waypoint,
	startLat, startLng float64,
) float64 {
	if len(order) == 0 {
		return 0.0
	}

	total := 0.0
	currentLat := startLat
	currentLng := startLng

	// Calculate distance from start to first waypoint
	firstIdx := order[0]
	firstWp := waypoints[firstIdx]
	total += haversineDistance(currentLat, currentLng, firstWp.Lat, firstWp.Lng)
	currentLat = firstWp.Lat
	currentLng = firstWp.Lng

	// Calculate distances between waypoints
	for i := 1; i < len(order); i++ {
		idx := order[i]
		wp := waypoints[idx]
		distance := haversineDistance(currentLat, currentLng, wp.Lat, wp.Lng)
		total += distance
		currentLat = wp.Lat
		currentLng = wp.Lng
	}

	return total
}

// ThreeOptOptimize improves route using 3-opt algorithm
// More complex than 2-opt but can provide better results
// Expected improvement: 20-30% reduction in total distance
func ThreeOptOptimize(
	initialOrder []int,
	waypoints []route_optimization.Waypoint,
	startLat, startLng float64,
) ([]int, float64) {
	if len(initialOrder) < 5 {
		// 3-opt requires at least 5 waypoints
		return TwoOptOptimize(initialOrder, waypoints, startLat, startLng)
	}

	bestOrder := make([]int, len(initialOrder))
	copy(bestOrder, initialOrder)
	bestDistance := calculateTotalDistance(bestOrder, waypoints, startLat, startLng)

	improved := true
	tuning := getImprovementTuning()
	maxIterations := tuning.threeOptMaxIter // 3-opt is more expensive, limit iterations
	iteration := 0
	deadline, hasDeadline := improvementDeadline(len(bestOrder))

	for improved && iteration < maxIterations {
		if hasDeadline && time.Now().After(deadline) {
			break
		}
		improved = false
		iteration++

		// Try all possible 3-opt swaps
		for i := 0; i < len(bestOrder)-2; i++ {
			if hasDeadline && time.Now().After(deadline) {
				break
			}
			for j := i + 2; j < len(bestOrder)-1; j++ {
				if hasDeadline && time.Now().After(deadline) {
					break
				}
				for k := j + 2; k < len(bestOrder); k++ {
					if hasDeadline && time.Now().After(deadline) {
						break
					}
					// Try different 3-opt reconnections
					orders := generateThreeOptVariants(bestOrder, i, j, k)

					for _, newOrder := range orders {
						newDistance := calculateTotalDistance(newOrder, waypoints, startLat, startLng)
						if newDistance < bestDistance {
							bestOrder = newOrder
							bestDistance = newDistance
							improved = true
							break
						}
					}
					if improved {
						break
					}
				}
				if improved {
					break
				}
			}
			if improved {
				break
			}
		}
	}

	return bestOrder, bestDistance
}

// generateThreeOptVariants generates all possible 3-opt reconnections
// 3-opt can reconnect 3 segments in 7 different ways
func generateThreeOptVariants(order []int, i, j, k int) [][]int {
	variants := make([][]int, 0, 7)

	// Original order (no change)
	original := make([]int, len(order))
	copy(original, order)
	variants = append(variants, original)

	// Variant 1: Reverse segment i+1 to j
	v1 := reverseSegment(order, i, j)
	variants = append(variants, v1)

	// Variant 2: Reverse segment j+1 to k
	v2 := reverseSegment(order, j, k)
	variants = append(variants, v2)

	// Variant 3: Reverse both segments
	v3 := reverseSegment(reverseSegment(order, i, j), j, k)
	variants = append(variants, v3)

	// Variant 4: Swap segments (i+1 to j) and (j+1 to k)
	v4 := swapSegments(order, i, j, k)
	variants = append(variants, v4)

	// Variant 5: Reverse first segment and swap
	v5 := swapSegments(reverseSegment(order, i, j), i, j, k)
	variants = append(variants, v5)

	// Variant 6: Reverse second segment and swap
	v6 := swapSegments(reverseSegment(order, j, k), i, j, k)
	variants = append(variants, v6)

	// Variant 7: Reverse both and swap
	v7 := swapSegments(reverseSegment(reverseSegment(order, i, j), j, k), i, j, k)
	variants = append(variants, v7)

	return variants
}

// swapSegments swaps two segments in the route
func swapSegments(order []int, i, j, k int) []int {
	newOrder := make([]int, 0, len(order))

	// Add segment before i
	newOrder = append(newOrder, order[:i+1]...)

	// Add segment j+1 to k (second segment)
	newOrder = append(newOrder, order[j+1:k+1]...)

	// Add segment i+1 to j (first segment, now second)
	newOrder = append(newOrder, order[i+1:j+1]...)

	// Add segment after k
	if k+1 < len(order) {
		newOrder = append(newOrder, order[k+1:]...)
	}

	return newOrder
}
