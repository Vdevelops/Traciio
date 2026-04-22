package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GeocodingService handles geocoding operations
type GeocodingService struct {
	httpClient *http.Client
	provider   string // "nominatim" or "google"
	apiKey     string // For Google Maps API (optional)
}

// GeocodeResult represents the result of a geocoding operation
type GeocodeResult struct {
	Latitude  float64
	Longitude float64
	Address   string
}

// NewGeocodingService creates a new geocoding service
func NewGeocodingService(provider string, apiKey string) *GeocodingService {
	return &GeocodingService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		provider: provider,
		apiKey:   apiKey,
	}
}

// GeocodeAddress converts an address string to latitude and longitude
// Address format: "Street, City, Province, Indonesia"
func (g *GeocodingService) GeocodeAddress(address string) (*GeocodeResult, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	// Clean and format address
	address = strings.TrimSpace(address)
	
	// Ensure address includes Indonesia for better results
	if !strings.Contains(strings.ToLower(address), "indonesia") {
		address = address + ", Indonesia"
	}

	switch g.provider {
	case "google":
		return g.geocodeWithGoogle(address)
	case "nominatim":
		fallthrough
	default:
		return g.geocodeWithNominatim(address)
	}
}

// GeocodeAddressWithFallback tries to geocode with multiple address variations
// Falls back to simpler addresses if detailed address fails
func (g *GeocodingService) GeocodeAddressWithFallback(street, city, province string) (*GeocodeResult, error) {
	// Strategy 1: Try full address (street + city + province)
	if street != "" && city != "" && province != "" {
		fullAddress := BuildFullAddress(street, city, province)
		result, err := g.GeocodeAddress(fullAddress)
		if err == nil {
			return result, nil
		}
	}

	// Strategy 2: Try simplified street (remove RT/RW, postal code, etc.)
	if street != "" && city != "" && province != "" {
		simplifiedStreet := simplifyAddress(street)
		if simplifiedStreet != street {
			fullAddress := BuildFullAddress(simplifiedStreet, city, province)
			result, err := g.GeocodeAddress(fullAddress)
			if err == nil {
				return result, nil
			}
		}
	}

	// Strategy 3: Try just street name + city + province (extract main street name)
	if street != "" && city != "" && province != "" {
		streetName := extractStreetName(street)
		if streetName != "" {
			fullAddress := BuildFullAddress(streetName, city, province)
			result, err := g.GeocodeAddress(fullAddress)
			if err == nil {
				return result, nil
			}
		}
	}

	// Strategy 4: Try city + province only
	if city != "" && province != "" {
		fullAddress := BuildFullAddress("", city, province)
		result, err := g.GeocodeAddress(fullAddress)
		if err == nil {
			return result, nil
		}
	}

	// All strategies failed
	return nil, fmt.Errorf("no results found for address after trying multiple variations")
}

// simplifyAddress removes RT/RW, postal codes, and other detailed parts
func simplifyAddress(address string) string {
	// Remove RT/RW patterns (e.g., "RT.4/RW.2" or "RT 4 RW 2")
	simplified := address
	simplified = strings.ReplaceAll(simplified, "RT.", "RT")
	simplified = strings.ReplaceAll(simplified, "RW.", "RW")
	
	// Remove RT/RW patterns with regex-like logic (simple string replacement)
	parts := strings.Split(simplified, ",")
	var cleanedParts []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Skip RT/RW patterns
		if strings.Contains(strings.ToUpper(part), "RT") && strings.Contains(strings.ToUpper(part), "RW") {
			continue
		}
		// Skip postal codes (usually 5 digits at the end)
		if len(part) == 5 && isNumeric(part) {
			continue
		}
		// Skip "Kec." (kecamatan) and "Kel." (kelurahan) parts if they're too detailed
		if strings.HasPrefix(strings.ToLower(part), "kec.") || strings.HasPrefix(strings.ToLower(part), "kel.") {
			continue
		}
		cleanedParts = append(cleanedParts, part)
	}
	
	return strings.Join(cleanedParts, ", ")
}

// extractStreetName extracts the main street name (e.g., "Jl. Salemba Raya" from "Jl. Salemba Raya No. 6, RT.4/RW.2...")
func extractStreetName(address string) string {
	parts := strings.Split(address, ",")
	if len(parts) > 0 {
		firstPart := strings.TrimSpace(parts[0])
		// Remove "No." and numbers after street name
		noIndex := strings.Index(strings.ToUpper(firstPart), " NO.")
		if noIndex > 0 {
			firstPart = firstPart[:noIndex]
		}
		// Remove trailing numbers
		firstPart = strings.TrimSpace(firstPart)
		return firstPart
	}
	return address
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// geocodeWithNominatim uses OpenStreetMap Nominatim API (free, no API key required)
func (g *GeocodingService) geocodeWithNominatim(address string) (*GeocodeResult, error) {
	// Build URL with proper encoding
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")
	params.Add("limit", "1")
	params.Add("addressdetails", "1")
	params.Add("countrycodes", "id") // Limit to Indonesia

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent (required by Nominatim)
	req.Header.Set("User-Agent", "CRM-Healthcare/1.0")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nominatim API returned status %d: %s", resp.StatusCode, string(body))
	}

	var results []struct {
		Lat     string `json:"lat"`
		Lon     string `json:"lon"`
		Display string `json:"display_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	result := results[0]
	
	var lat, lon float64
	if _, err := fmt.Sscanf(result.Lat, "%f", &lat); err != nil {
		return nil, fmt.Errorf("failed to parse latitude: %w", err)
	}
	if _, err := fmt.Sscanf(result.Lon, "%f", &lon); err != nil {
		return nil, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return &GeocodeResult{
		Latitude:  lat,
		Longitude: lon,
		Address:   result.Display,
	}, nil
}

// geocodeWithGoogle uses Google Maps Geocoding API (requires API key)
func (g *GeocodingService) geocodeWithGoogle(address string) (*GeocodeResult, error) {
	if g.apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key is required")
	}

	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", g.apiKey)
	params.Add("region", "id") // Bias results to Indonesia

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := g.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Status  string `json:"status"`
		Results []struct {
			FormattedAddress string `json:"formatted_address"`
			Geometry         struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("Google API returned status: %s", result.Status)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	firstResult := result.Results[0]
	return &GeocodeResult{
		Latitude:  firstResult.Geometry.Location.Lat,
		Longitude: firstResult.Geometry.Location.Lng,
		Address:   firstResult.FormattedAddress,
	}, nil
}

// BuildFullAddress constructs a full address string from components
func BuildFullAddress(street, city, province string) string {
	parts := []string{}
	if street != "" {
		parts = append(parts, street)
	}
	if city != "" {
		parts = append(parts, city)
	}
	if province != "" {
		parts = append(parts, province)
	}
	return strings.Join(parts, ", ")
}

