package area_mapping

import (
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// AreaCapture represents a GPS location capture from a visit
type AreaCapture struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VisitReportID string    `gorm:"type:uuid;not null" json:"visit_report_id"`
	CaptureType   string    `gorm:"type:varchar(20);not null" json:"capture_type"` // check_in, check_out, area
	Location      GeoPoint  `gorm:"type:geography(POINT,4326);not null" json:"location"`
	Address       *string   `gorm:"type:text" json:"address,omitempty"`
	Accuracy      *float64  `gorm:"type:decimal(10,2)" json:"accuracy,omitempty"`
	CapturedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"captured_at"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for AreaCapture
func (AreaCapture) TableName() string {
	return "area_captures"
}

// LocationJSON represents location in JSON format for API responses
type LocationJSON struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// GetLocationJSON converts GeoPoint to LocationJSON
func (a *AreaCapture) GetLocationJSON() LocationJSON {
	return LocationJSON{
		Lat: a.Location.Point.Lat(),
		Lng: a.Location.Point.Lon(),
	}
}

// SetLocationFromJSON sets location from LocationJSON
func (a *AreaCapture) SetLocationFromJSON(loc LocationJSON) {
	a.Location = GeoPoint{Point: orb.Point{loc.Lng, loc.Lat}}
}

// Territory represents a geographic territory assigned to a user
type Territory struct {
	ID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Polygon     GeoPolygon `gorm:"type:geography(POLYGON,4326);not null" json:"polygon"`
	AssignedTo  *string    `gorm:"type:uuid" json:"assigned_to,omitempty"`
	Color       string     `gorm:"type:varchar(50);default:'#3B82F6'" json:"color"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for Territory
func (Territory) TableName() string {
	return "territories"
}

// PolygonJSON represents polygon coordinates in GeoJSON format
type PolygonJSON struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// GetPolygonJSON converts GeoPolygon to GeoJSON format
func (t *Territory) GetPolygonJSON() PolygonJSON {
	coords := make([][][]float64, len(t.Polygon.Polygon))
	for i, ring := range t.Polygon.Polygon {
		coords[i] = make([][]float64, len(ring))
		for j, point := range ring {
			coords[i][j] = []float64{point.Lon(), point.Lat()}
		}
	}
	return PolygonJSON{
		Type:        "Polygon",
		Coordinates: coords,
	}
}

// SetPolygonFromJSON sets polygon from GeoJSON coordinates
func (t *Territory) SetPolygonFromJSON(coords [][][]float64) error {
	polygon := make(orb.Polygon, len(coords))
	for i, ring := range coords {
		polygon[i] = make(orb.Ring, len(ring))
		for j, point := range ring {
			if len(point) < 2 {
				continue
			}
			polygon[i][j] = orb.Point{point[0], point[1]}
		}
	}
	t.Polygon = GeoPolygon{Polygon: polygon}
	return nil
}

// GetPolygonWKT returns the polygon as WKT (Well-Known Text)
func (t *Territory) GetPolygonWKT() string {
	return wkt.MarshalString(t.Polygon.Polygon)
}

// CoverageAnalysis represents coverage analysis for a territory
type CoverageAnalysis struct {
	ID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TerritoryID     *string   `gorm:"type:uuid" json:"territory_id,omitempty"`
	UserID          *string   `gorm:"type:uuid" json:"user_id,omitempty"`
	PeriodStart     time.Time `gorm:"type:date;not null" json:"period_start"`
	PeriodEnd       time.Time `gorm:"type:date;not null" json:"period_end"`
	VisitCount      int       `gorm:"not null;default:0" json:"visit_count"`
	CoveragePercent *float64  `gorm:"type:decimal(5,2)" json:"coverage_percent,omitempty"`
	AnalyzedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"analyzed_at"`
	CreatedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for CoverageAnalysis
func (CoverageAnalysis) TableName() string {
	return "coverage_analysis"
}

// CaptureTypeValues lists valid capture types
var CaptureTypeValues = []string{"check_in", "check_out", "area"}

// ValidateCaptureType validates if capture type is valid
func ValidateCaptureType(captureType string) bool {
	for _, valid := range CaptureTypeValues {
		if captureType == valid {
			return true
		}
	}
	return false
}
