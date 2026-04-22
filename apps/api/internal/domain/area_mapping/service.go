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

// Request DTOs
type CaptureLocationRequest struct {
	VisitReportID string  `json:"visit_report_id" binding:"required"`
	CaptureType   string  `json:"capture_type" binding:"required,oneof=check_in check_out area"`
	Latitude      float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude     float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Address       string  `json:"address"`
	Accuracy      float64 `json:"accuracy"`
}

type CreateTerritoryRequest struct {
	Name        string      `json:"name" binding:"required,min=3,max=255"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates" binding:"required,min=4"` // [lng, lat] pairs
	AssignedTo  string      `json:"assigned_to"`
	Color       string      `json:"color"`
}

type UpdateTerritoryRequest struct {
	Name        *string      `json:"name" binding:"omitempty,min=3,max=255"`
	Description *string      `json:"description"`
	Coordinates *[][]float64 `json:"coordinates" binding:"omitempty,min=4"`
	AssignedTo  *string      `json:"assigned_to"`
	Color       *string      `json:"color"`
}
