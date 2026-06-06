package route_optimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/pkg/circuitbreaker"
)

var (
	ErrOSRMUnavailable  = errors.New("OSRM routing service unavailable")
	ErrOSRMInvalidRoute = errors.New("OSRM returned invalid route")
)

// OSRMClient handles communication with OSRM routing service
type OSRMClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *circuitbreaker.CircuitBreaker
}

// NewOSRMClient creates a new OSRM client with optimized settings
func NewOSRMClient(baseURL string) *OSRMClient {
	return NewOSRMClientWithTimeout(baseURL, 20*time.Second)
}

// NewOSRMClientWithTimeout creates an OSRM client with a custom HTTP timeout.
// Use a short timeout (e.g., 5 s) for leg-by-leg fallback calls to avoid
// stacking multiple 20 s waits and hitting gateway timeouts.
func NewOSRMClientWithTimeout(baseURL string, timeout time.Duration) *OSRMClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	// Configure circuit breaker for OSRM
	cbConfig := circuitbreaker.Config{
		Name:         "osrm",
		MaxRequests:  100,
		Interval:     0,
		Timeout:      60 * time.Second,
		FailureRatio: 0.5, // Trip at 50% failure rate (more sensitive for external API)
		MinRequests:  5,
		SuccessOn:    3,
	}

	return &OSRMClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		cb: circuitbreaker.New(cbConfig),
	}
}

// Route calculates route between multiple waypoints using OSRM
func (c *OSRMClient) Route(req RouteRequest) (*OSRMRouteResponse, error) {
	if len(req.Coordinates) < 2 {
		return nil, errors.New("at least 2 coordinates required")
	}

	// Build coordinates string: lng,lat;lng,lat;...
	// OSRM uses lng,lat format (longitude first, then latitude)
	coordsStr := ""
	for i, coord := range req.Coordinates {
		if len(coord) != 2 {
			return nil, fmt.Errorf("invalid coordinate at index %d: expected [lng, lat]", i)
		}
		if i > 0 {
			coordsStr += ";"
		}
		coordsStr += fmt.Sprintf("%f,%f", coord[0], coord[1]) // OSRM uses lng,lat format
	}

	// Set defaults
	profile := req.Profile
	if profile == "" {
		profile = "driving"
	}
	overview := req.Overview
	if overview == "" {
		overview = "full"
	}
	geometries := req.Geometries
	if geometries == "" {
		geometries = "polyline"
	}

	// Build URL
	routeURL := fmt.Sprintf("%s/route/v1/%s/%s", c.baseURL, profile, url.PathEscape(coordsStr))

	// Build query parameters
	params := url.Values{}
	params.Set("overview", overview)
	params.Set("geometries", geometries)
	if req.Steps {
		params.Set("steps", "true")
	}
	if req.Alternatives {
		params.Set("alternatives", "true")
	}

	if len(params) > 0 {
		routeURL += "?" + params.Encode()
	}

	// Execute with circuit breaker
	result, err := c.cb.Execute(func() (interface{}, error) {
		resp, err := c.httpClient.Get(routeURL)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrOSRMUnavailable, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("%w: status %d, response: %s", ErrOSRMUnavailable, resp.StatusCode, string(body))
		}

		// Parse response
		var routeResp OSRMRouteResponse
		if err := json.NewDecoder(resp.Body).Decode(&routeResp); err != nil {
			return nil, fmt.Errorf("failed to parse OSRM response: %w", err)
		}

		if routeResp.Code != "Ok" {
			return nil, fmt.Errorf("%w: OSRM code: %s", ErrOSRMInvalidRoute, routeResp.Code)
		}

		if len(routeResp.Routes) == 0 {
			return nil, fmt.Errorf("%w: no routes returned", ErrOSRMInvalidRoute)
		}

		return &routeResp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*OSRMRouteResponse), nil
}

// Table fetches a distance/duration matrix for the given coordinates using OSRM /table.
// This can replace many pairwise /route calls with a single request.
func (c *OSRMClient) Table(req TableRequest) (*OSRMTableResponse, error) {
	if len(req.Coordinates) < 2 {
		return nil, errors.New("at least 2 coordinates required")
	}

	coordsStr := ""
	for i, coord := range req.Coordinates {
		if len(coord) != 2 {
			return nil, fmt.Errorf("invalid coordinate at index %d: expected [lng, lat]", i)
		}
		if i > 0 {
			coordsStr += ";"
		}
		coordsStr += fmt.Sprintf("%f,%f", coord[0], coord[1])
	}

	profile := req.Profile
	if profile == "" {
		profile = "driving"
	}

	annotations := req.Annotations
	if len(annotations) == 0 {
		annotations = []string{"distance", "duration"}
	}

	tableURL := fmt.Sprintf("%s/table/v1/%s/%s", c.baseURL, profile, url.PathEscape(coordsStr))
	params := url.Values{}
	params.Set("annotations", strings.Join(annotations, ","))
	if len(params) > 0 {
		tableURL += "?" + params.Encode()
	}

	// Execute with circuit breaker
	result, err := c.cb.Execute(func() (interface{}, error) {
		resp, err := c.httpClient.Get(tableURL)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrOSRMUnavailable, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("%w: status %d, response: %s", ErrOSRMUnavailable, resp.StatusCode, string(body))
		}

		var tableResp OSRMTableResponse
		if err := json.NewDecoder(resp.Body).Decode(&tableResp); err != nil {
			return nil, fmt.Errorf("failed to parse OSRM table response: %w", err)
		}

		if tableResp.Code != "Ok" {
			return nil, fmt.Errorf("%w: OSRM code: %s", ErrOSRMInvalidRoute, tableResp.Code)
		}

		return &tableResp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*OSRMTableResponse), nil
}

// ConvertOSRMRouteToRouteSteps converts OSRM route response to RouteStep array
func ConvertOSRMRouteToRouteSteps(osrmRoute *OSRMRoute, waypoints []route_optimization.Waypoint) []route_optimization.RouteStep {
	steps := make([]route_optimization.RouteStep, 0)
	stepNumber := 1

	// Process each leg
	for legIdx, leg := range osrmRoute.Legs {
		// Get start location (from waypoint)
		var startLat, startLng float64
		if legIdx < len(waypoints) {
			startLat = waypoints[legIdx].Lat
			startLng = waypoints[legIdx].Lng
		}

		// Process steps in leg
		for stepIdx, step := range leg.Steps {
			// Get end location from maneuver or next step
			// OSRM maneuver location uses [lng, lat] format
			var endLat, endLng float64
			if len(step.Maneuver.Location) >= 2 {
				endLng = step.Maneuver.Location[0] // OSRM uses lng first
				endLat = step.Maneuver.Location[1] // then lat
			} else if stepIdx+1 < len(leg.Steps) && len(leg.Steps[stepIdx+1].Maneuver.Location) >= 2 {
				endLng = leg.Steps[stepIdx+1].Maneuver.Location[0]
				endLat = leg.Steps[stepIdx+1].Maneuver.Location[1]
			} else if legIdx+1 < len(waypoints) {
				endLat = waypoints[legIdx+1].Lat
				endLng = waypoints[legIdx+1].Lng
			} else {
				// Fallback: use start location
				endLat = startLat
				endLng = startLng
			}

			// Convert distance from meters to kilometers
			distanceKm := step.Distance / 1000.0
			durationSec := int(step.Duration)

			// Generate instruction
			instruction := step.Maneuver.Instruction
			if instruction == "" {
				instruction = generateInstruction(step.Maneuver, step.Distance)
			}

			routeStep := route_optimization.RouteStep{
				Step:              stepNumber,
				Distance:          distanceKm,
				DistanceFormatted: formatDistance(distanceKm),
				Duration:          durationSec,
				DurationFormatted: formatDuration(durationSec),
				Instruction:       instruction,
				Polyline:          step.Geometry,
				Maneuver:          step.Maneuver.Type,
				StartLocation: route_optimization.Location{
					Lat: startLat,
					Lng: startLng,
				},
				EndLocation: route_optimization.Location{
					Lat: endLat,
					Lng: endLng,
				},
			}

			steps = append(steps, routeStep)
			stepNumber++

			// Update start location for next step
			startLat = endLat
			startLng = endLng
		}
	}

	return steps
}

// generateInstruction generates a human-readable instruction from maneuver
func generateInstruction(maneuver OSRMManeuver, distance float64) string {
	distanceStr := formatDistance(distance / 1000.0) // Convert meters to km

	modifier := maneuver.Modifier
	if modifier == "" {
		modifier = "straight"
	}

	instruction := fmt.Sprintf("%s for %s", modifier, distanceStr)

	switch maneuver.Type {
	case "turn":
		return fmt.Sprintf("Turn %s and continue for %s", modifier, distanceStr)
	case "new name":
		return fmt.Sprintf("Continue straight on %s for %s", modifier, distanceStr)
	case "depart":
		return fmt.Sprintf("Depart and head %s for %s", modifier, distanceStr)
	case "arrive":
		return "Arrive at destination"
	default:
		return instruction
	}
}

// DecodePolylineToPoints decodes a Google/OSRM encoded polyline string into lat/lng pairs.
// Returns a slice of [2]float64 where each element is [lat, lng].
func DecodePolylineToPoints(encoded string) [][2]float64 {
	points := make([][2]float64, 0)
	index := 0
	lat := 0
	lng := 0

	for index < len(encoded) {
		// Decode latitude
		shift, result := 0, 0
		for {
			if index >= len(encoded) {
				break
			}
			b := int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		if result&1 != 0 {
			lat += ^(result >> 1)
		} else {
			lat += result >> 1
		}

		// Decode longitude
		shift, result = 0, 0
		for {
			if index >= len(encoded) {
				break
			}
			b := int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		if result&1 != 0 {
			lng += ^(result >> 1)
		} else {
			lng += result >> 1
		}

		points = append(points, [2]float64{float64(lat) / 1e5, float64(lng) / 1e5})
	}

	return points
}

// EncodePolylineFromPoints encodes lat/lng coordinate pairs into Google Polyline format.
// Input: [][2]float64 where each element is [lat, lng].
func EncodePolylineFromPoints(points [][2]float64) string {
	var buf strings.Builder
	prevLat, prevLng := 0, 0

	for _, pt := range points {
		lat := int(math.Round(pt[0] * 1e5))
		lng := int(math.Round(pt[1] * 1e5))
		encodePolylineVarint(&buf, lat-prevLat)
		encodePolylineVarint(&buf, lng-prevLng)
		prevLat = lat
		prevLng = lng
	}

	return buf.String()
}

// encodePolylineVarint encodes a single signed integer difference into the Google Polyline varint format.
func encodePolylineVarint(buf *strings.Builder, n int) {
	n <<= 1
	if n < 0 {
		n = ^n
	}
	for n >= 0x20 {
		buf.WriteByte(byte((0x20 | (n & 0x1f)) + 63))
		n >>= 5
	}
	buf.WriteByte(byte(n + 63))
}
