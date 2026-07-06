package area_mapping

import (
	"time"
)

// Service defines the interface for area mapping business logic
type Service interface {
	// Area Capture operations
	CaptureLocation(req CaptureLocationRequest) (*AreaCapture, error)
	GetAreaCaptures(filter AreaCaptureFilter) ([]AreaCapture, int64, error)
	GetAreaCaptureByID(id string) (*AreaCapture, error)

	// Territory operations
	CreateTerritory(req CreateTerritoryRequest) (*Territory, error)
	UpdateTerritory(id string, req UpdateTerritoryRequest) (*Territory, error)
	GetTerritories(filter TerritoryFilter) ([]Territory, int64, error)
	GetTerritoryByID(id string) (*Territory, error)
	DeleteTerritory(id string) error

	// Spatial analysis operations
	CheckPointInTerritory(lat, lng float64, territoryID string) (bool, error)
	GetCapturesWithinTerritory(territoryID string, startDate, endDate *time.Time) ([]AreaCapture, error)
	CalculateCoverage(territoryID string, startDate, endDate time.Time) (*CoverageAnalysis, error)
	AssignTerritory(userID, territoryID string) error
	GetHeatmapData(filter HeatmapFilter) ([]HeatmapPoint, error)
}
