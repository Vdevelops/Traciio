package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type GeocodingHandler struct{}

func NewGeocodingHandler() *GeocodingHandler {
	return &GeocodingHandler{}
}

// Geocode handles geocoding request (address to coordinates)
func (h *GeocodingHandler) Geocode(c *gin.Context) {
	var req GeocodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.InvalidRequestBodyResponse(c)
		return
	}

	if req.Address == "" {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Address is required",
		}, nil)
		return
	}

	// Use Nominatim for geocoding with timeout and retry
	coords, err := geocodeWithNominatim(req.Address)
	if err != nil {
		// Check if error is "no results found" - try with simplified address
		errMsg := err.Error()
		if strings.Contains(errMsg, "no results found") {
			// Try geocoding with just city and province if available
			// Extract city/province from address for fallback
			simplifiedAddress := trySimplifiedAddress(req.Address)
			if simplifiedAddress != req.Address {
				// Retry with simplified address
				coords, retryErr := geocodeWithNominatim(simplifiedAddress)
				if retryErr == nil {
					// Success with simplified address
					response.SuccessResponse(c, GeocodeResponse{
						Latitude:  coords.Latitude,
						Longitude: coords.Longitude,
						Address:   req.Address, // Keep original address
					}, nil)
					return
				}
			}

			// Still no results - return error
			errors.ErrorResponse(c, "GEOCODING_NO_RESULTS", map[string]interface{}{
				"message": fmt.Sprintf("No geocoding results found for address: %s. Please verify the address is correct.", req.Address),
			}, nil)
			return
		}
		// For other errors (network, timeout, etc.), return 500
		errors.InternalServerErrorResponse(c, fmt.Sprintf("Failed to geocode address: %v", err))
		return
	}

	response.SuccessResponse(c, GeocodeResponse{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
		Address:   req.Address,
	}, nil)
}

// geocodeWithNominatim performs geocoding using Nominatim API
func geocodeWithNominatim(address string) (*GeocodeResponse, error) {
	// Check if address already contains "Indonesia" to avoid duplication
	query := address
	if !strings.Contains(strings.ToLower(address), "indonesia") {
		query = fmt.Sprintf("%s, Indonesia", address)
	}
	encodedQuery := url.QueryEscape(query)

	// Build Nominatim URL
	apiURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1&countrycodes=id",
		encodedQuery,
	)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request with proper headers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "CRM-Healthcare-Backend/1.0")
	req.Header.Set("Accept-Language", "id,en")

	// Make request with retry logic
	var resp *http.Response
	maxRetries := 2
	for i := 0; i <= maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				break // Success, exit retry loop
			}
			// Non-200 status code - close body and retry if not last attempt
			if resp != nil {
				resp.Body.Close()
			}
			// If rate limited (429) or server error (5xx), retry
			if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
				if i < maxRetries {
					// Longer delay for rate limiting
					delay := time.Duration(i+1) * 2 * time.Second
					time.Sleep(delay)
					continue
				}
			}
			// For other non-200 status codes, don't retry
			if i < maxRetries {
				time.Sleep(time.Duration(i+1) * time.Second)
			}
		} else {
			// Network error - retry if not last attempt
			if resp != nil {
				resp.Body.Close()
			}
			if i < maxRetries {
				time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries+1, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("no response received after %d retries", maxRetries+1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("nominatim API returned status %d after %d retries", resp.StatusCode, maxRetries+1)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	// Extract coordinates
	result := results[0]
	latStr, ok := result["lat"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid latitude in response")
	}

	lonStr, ok := result["lon"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid longitude in response")
	}

	var lat, lon float64
	if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
		return nil, fmt.Errorf("failed to parse latitude: %w", err)
	}

	if _, err := fmt.Sscanf(lonStr, "%f", &lon); err != nil {
		return nil, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return &GeocodeResponse{
		Latitude:  lat,
		Longitude: lon,
		Address:   address,
	}, nil
}

// trySimplifiedAddress attempts to extract city/province from full address for fallback geocoding
func trySimplifiedAddress(fullAddress string) string {
	// Remove "Indonesia" if present to avoid duplication
	address := strings.TrimSpace(fullAddress)
	address = strings.TrimSuffix(address, ", Indonesia")
	address = strings.TrimSuffix(address, ",Indonesia")
	address = strings.TrimSpace(address)

	// Try to extract last two parts (city, province) if address has multiple parts
	parts := strings.Split(address, ",")
	if len(parts) >= 2 {
		// Take last two parts (usually city and province)
		simplified := strings.TrimSpace(parts[len(parts)-2]) + ", " + strings.TrimSpace(parts[len(parts)-1])
		return simplified + ", Indonesia"
	}

	// If only one part, return as is with Indonesia
	return address + ", Indonesia"
}
