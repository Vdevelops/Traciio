package area_mapping

import (
	"fmt"
	"time"

	"github.com/paulmach/orb"

	"github.com/gilabs/crm-healthcare/api/internal/domain/area_mapping"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
)

type service struct {
	repo area_mapping.Repository
}

// NewService creates a new area mapping service
func NewService(repo area_mapping.Repository) area_mapping.Service {
	return &service{
		repo: repo,
	}
}

// CaptureLocation captures a new location point
func (s *service) CaptureLocation(req area_mapping.CaptureLocationRequest) (*area_mapping.AreaCapture, error) {
	// Validate capture type
	if !area_mapping.ValidateCaptureType(req.CaptureType) {
		return nil, fmt.Errorf("invalid capture type: %s", req.CaptureType)
	}

	// Create capture
	capture := &area_mapping.AreaCapture{
		VisitReportID: req.VisitReportID,
		CaptureType:   req.CaptureType,
		Location:      area_mapping.GeoPoint{Point: orb.Point{req.Longitude, req.Latitude}}, // orb uses [lng, lat]
		CapturedAt:    time.Now(),
	}

	if req.Address != "" {
		capture.Address = &req.Address
	}

	if req.Accuracy > 0 {
		capture.Accuracy = &req.Accuracy
	}

	err := s.repo.CreateAreaCapture(capture)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "create_area_capture"})
		return nil, fmt.Errorf("failed to capture location: %w", err)
	}

	logger.LogInfo("Location captured successfully", map[string]interface{}{
		"capture_id": capture.ID,
		"type":       capture.CaptureType,
	})
	return capture, nil
}

// GetAreaCaptures retrieves area captures with filters
func (s *service) GetAreaCaptures(filter area_mapping.AreaCaptureFilter) ([]area_mapping.AreaCapture, int64, error) {
	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 10
	}

	captures, total, err := s.repo.ListAreaCaptures(filter)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "list_area_captures"})
		return nil, 0, fmt.Errorf("failed to get area captures: %w", err)
	}

	return captures, total, nil
}

// GetAreaCaptureByID retrieves an area capture by ID
func (s *service) GetAreaCaptureByID(id string) (*area_mapping.AreaCapture, error) {
	capture, err := s.repo.GetAreaCaptureByID(id)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "get_area_capture", "id": id})
		return nil, fmt.Errorf("area capture not found: %w", err)
	}

	return capture, nil
}

// CreateTerritory creates a new territory polygon
func (s *service) CreateTerritory(req area_mapping.CreateTerritoryRequest) (*area_mapping.Territory, error) {
	// Validate coordinates (minimum 4 points for closed polygon)
	if len(req.Coordinates) < 4 {
		return nil, fmt.Errorf("polygon must have at least 4 points")
	}

	// Convert coordinates to orb.Polygon
	ring := make(orb.Ring, len(req.Coordinates))
	for i, coord := range req.Coordinates {
		if len(coord) != 2 {
			return nil, fmt.Errorf("invalid coordinate format at index %d", i)
		}
		lng, lat := coord[0], coord[1]

		// Validate coordinates
		if lat < -90 || lat > 90 {
			return nil, fmt.Errorf("latitude must be between -90 and 90 at index %d", i)
		}
		if lng < -180 || lng > 180 {
			return nil, fmt.Errorf("longitude must be between -180 and 180 at index %d", i)
		}

		ring[i] = orb.Point{lng, lat} // orb uses [lng, lat]
	}

	// Ensure polygon is closed (first point == last point)
	if !ring[0].Equal(ring[len(ring)-1]) {
		ring = append(ring, ring[0])
	}

	polygon := orb.Polygon{ring}

	// Create territory
	territory := &area_mapping.Territory{
		Name:    req.Name,
		Polygon: area_mapping.GeoPolygon{Polygon: polygon},
		Color:   req.Color,
	}

	if req.Description != "" {
		territory.Description = &req.Description
	}

	if req.AssignedTo != "" {
		territory.AssignedTo = &req.AssignedTo
	}

	if territory.Color == "" {
		territory.Color = "#3B82F6" // Default blue
	}

	err := s.repo.CreateTerritory(territory)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "create_territory"})
		return nil, fmt.Errorf("failed to create territory: %w", err)
	}

	logger.LogInfo("Territory created successfully", map[string]interface{}{
		"territory_id": territory.ID,
		"name":         territory.Name,
	})
	return territory, nil
}

// UpdateTerritory updates an existing territory
func (s *service) UpdateTerritory(id string, req area_mapping.UpdateTerritoryRequest) (*area_mapping.Territory, error) {
	// Get existing territory
	territory, err := s.repo.GetTerritoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("territory not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		territory.Name = *req.Name
	}
	if req.Description != nil {
		territory.Description = req.Description
	}
	if req.AssignedTo != nil {
		territory.AssignedTo = req.AssignedTo
	}
	if req.Color != nil {
		territory.Color = *req.Color
	}

	// Update polygon if coordinates provided
	if req.Coordinates != nil {
		if len(*req.Coordinates) < 4 {
			return nil, fmt.Errorf("polygon must have at least 4 points")
		}

		ring := make(orb.Ring, len(*req.Coordinates))
		for i, coord := range *req.Coordinates {
			if len(coord) != 2 {
				return nil, fmt.Errorf("invalid coordinate format at index %d", i)
			}
			lng, lat := coord[0], coord[1]

			if lat < -90 || lat > 90 {
				return nil, fmt.Errorf("latitude must be between -90 and 90 at index %d", i)
			}
			if lng < -180 || lng > 180 {
				return nil, fmt.Errorf("longitude must be between -180 and 180 at index %d", i)
			}

			ring[i] = orb.Point{lng, lat}
		}

		if !ring[0].Equal(ring[len(ring)-1]) {
			ring = append(ring, ring[0])
		}

		territory.Polygon = area_mapping.GeoPolygon{Polygon: orb.Polygon{ring}}
	}

	territory.UpdatedAt = time.Now()

	err = s.repo.UpdateTerritory(territory)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "update_territory", "id": id})
		return nil, fmt.Errorf("failed to update territory: %w", err)
	}

	logger.LogInfo("Territory updated successfully", map[string]interface{}{"territory_id": id})
	return territory, nil
}

// GetTerritories retrieves territories with filters
func (s *service) GetTerritories(filter area_mapping.TerritoryFilter) ([]area_mapping.Territory, int64, error) {
	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 10
	}

	territories, total, err := s.repo.ListTerritories(filter)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "list_territories"})
		return nil, 0, fmt.Errorf("failed to get territories: %w", err)
	}

	return territories, total, nil
}

// GetTerritoryByID retrieves a territory by ID
func (s *service) GetTerritoryByID(id string) (*area_mapping.Territory, error) {
	territory, err := s.repo.GetTerritoryByID(id)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "get_territory", "id": id})
		return nil, fmt.Errorf("territory not found: %w", err)
	}

	return territory, nil
}

// DeleteTerritory deletes a territory
func (s *service) DeleteTerritory(id string) error {
	err := s.repo.DeleteTerritory(id)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "delete_territory", "id": id})
		return fmt.Errorf("failed to delete territory: %w", err)
	}

	logger.LogInfo("Territory deleted successfully", map[string]interface{}{"territory_id": id})
	return nil
}

// CheckPointInTerritory checks if a point is within a specific territory
func (s *service) CheckPointInTerritory(lat, lng float64, territoryID string) (bool, error) {
	if lat < -90 || lat > 90 {
		return false, fmt.Errorf("latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return false, fmt.Errorf("longitude must be between -180 and 180")
	}

	point := orb.Point{lng, lat} // orb uses [lng, lat]

	isWithin, err := s.repo.CheckPointInTerritory(point, territoryID)
	if err != nil {
		logger.LogError(err, map[string]interface{}{
			"action": "check_point_in_territory",
			"lat":    lat,
			"lng":    lng,
		})
		return false, fmt.Errorf("failed to check territory: %w", err)
	}

	return isWithin, nil
}

// GetCapturesWithinTerritory retrieves all captures within a territory
func (s *service) GetCapturesWithinTerritory(territoryID string, startDate, endDate *time.Time) ([]area_mapping.AreaCapture, error) {
	captures, err := s.repo.GetCapturesWithinTerritory(territoryID, startDate, endDate)
	if err != nil {
		logger.LogError(err, map[string]interface{}{
			"action":       "get_captures_within_territory",
			"territory_id": territoryID,
		})
		return nil, fmt.Errorf("failed to get captures: %w", err)
	}

	return captures, nil
}

// CalculateCoverage calculates coverage analysis for a territory
func (s *service) CalculateCoverage(territoryID string, startDate, endDate time.Time) (*area_mapping.CoverageAnalysis, error) {
	// Get territory
	territory, err := s.repo.GetTerritoryByID(territoryID)
	if err != nil {
		return nil, fmt.Errorf("territory not found: %w", err)
	}

	// Get captures within territory and date range
	captures, err := s.repo.GetCapturesWithinTerritory(territoryID, &startDate, &endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get captures: %w", err)
	}

	// Count unique visits
	visitMap := make(map[string]bool)
	for _, capture := range captures {
		visitMap[capture.VisitReportID] = true
	}

	visitCount := len(visitMap)
	coveragePercent := 0.0 // This would need actual calculation based on expected visits

	// Create coverage analysis
	analysis := &area_mapping.CoverageAnalysis{
		TerritoryID:     &territory.ID,
		PeriodStart:     startDate,
		PeriodEnd:       endDate,
		VisitCount:      visitCount,
		CoveragePercent: &coveragePercent,
		AnalyzedAt:      time.Now(),
	}

	err = s.repo.CreateCoverageAnalysis(analysis)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "create_coverage_analysis"})
		return nil, fmt.Errorf("failed to save coverage analysis: %w", err)
	}

	logger.LogInfo("Coverage analysis completed", map[string]interface{}{
		"territory_id": territoryID,
		"visits":       visitCount,
	})
	return analysis, nil
}

// AssignTerritory assigns a territory to a user
func (s *service) AssignTerritory(userID, territoryID string) error {
	// Get territory
	territory, err := s.repo.GetTerritoryByID(territoryID)
	if err != nil {
		return fmt.Errorf("territory not found: %w", err)
	}

	// Update assigned user
	territory.AssignedTo = &userID
	territory.UpdatedAt = time.Now()

	err = s.repo.UpdateTerritory(territory)
	if err != nil {
		logger.LogError(err, map[string]interface{}{
			"action":       "assign_territory",
			"territory_id": territoryID,
			"user_id":      userID,
		})
		return fmt.Errorf("failed to assign territory: %w", err)
	}

	logger.LogInfo("Territory assigned successfully", map[string]interface{}{
		"territory_id": territoryID,
		"user_id":      userID,
	})
	return nil
}

// GetHeatmapData retrieves aggregated location data for heatmap
func (s *service) GetHeatmapData(filter area_mapping.HeatmapFilter) ([]area_mapping.HeatmapPoint, error) {
	heatmapPoints, err := s.repo.GetHeatmapData(filter)
	if err != nil {
		logger.LogError(err, map[string]interface{}{"action": "get_heatmap_data"})
		return nil, fmt.Errorf("failed to get heatmap data: %w", err)
	}

	return heatmapPoints, nil
}
