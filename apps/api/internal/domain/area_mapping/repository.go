package area_mapping

import (
	"time"

	"github.com/paulmach/orb"
)

// Repository defines the interface for area mapping data access
type Repository interface {
	// Area Capture operations
	CreateAreaCapture(capture *AreaCapture) error
	GetAreaCaptureByID(id string) (*AreaCapture, error)
	ListAreaCaptures(filter AreaCaptureFilter) ([]AreaCapture, int64, error)
	DeleteAreaCapture(id string) error

	// Territory operations
	CreateTerritory(territory *Territory) error
	GetTerritoryByID(id string) (*Territory, error)
	ListTerritories(filter TerritoryFilter) ([]Territory, int64, error)
	UpdateTerritory(territory *Territory) error
	DeleteTerritory(id string) error

	// Coverage Analysis operations
	CreateCoverageAnalysis(analysis *CoverageAnalysis) error
	GetCoverageAnalysisByID(id string) (*CoverageAnalysis, error)
	ListCoverageAnalysis(filter CoverageAnalysisFilter) ([]CoverageAnalysis, int64, error)

	// Spatial queries
	GetCapturesWithinTerritory(territoryID string, startDate, endDate *time.Time) ([]AreaCapture, error)
	CheckPointInTerritory(point orb.Point, territoryID string) (bool, error)
	CalculateDistance(point1, point2 orb.Point) (float64, error)
	GetHeatmapData(filter HeatmapFilter) ([]HeatmapPoint, error)
}

// AreaCaptureFilter represents filter criteria for area captures
type AreaCaptureFilter struct {
	Page           int
	PerPage        int
	VisitReportID  *string
	CaptureType    *string
	CapturedAfter  *time.Time
	CapturedBefore *time.Time
}

// TerritoryFilter represents filter criteria for territories
type TerritoryFilter struct {
	Page       int
	PerPage    int
	Search     *string
	AssignedTo *string
}

// CoverageAnalysisFilter represents filter criteria for coverage analysis
type CoverageAnalysisFilter struct {
	Page        int
	PerPage     int
	TerritoryID *string
	UserID      *string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
}

// HeatmapFilter represents filter criteria for heatmap data
type HeatmapFilter struct {
	TerritoryID *string
	UserID      *string
	StartDate   *time.Time
	EndDate     *time.Time
	CaptureType *string
}

// HeatmapPoint represents a point on the heatmap with visit frequency
type HeatmapPoint struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Intensity int     `json:"intensity"`
}
