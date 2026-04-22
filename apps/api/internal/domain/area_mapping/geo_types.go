package area_mapping

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// GeoPoint is a custom type for orb.Point that implements sql.Scanner and driver.Valuer
type GeoPoint struct {
	orb.Point
}

// Scan implements sql.Scanner interface for reading from database
func (g *GeoPoint) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var wktString string
	switch v := src.(type) {
	case string:
		wktString = v
	case []byte:
		wktString = string(v)
	default:
		return fmt.Errorf("unsupported type for GeoPoint: %T", src)
	}

	// Parse WKT string to orb.Point
	geom, err := wkt.Unmarshal(wktString)
	if err != nil {
		return fmt.Errorf("failed to parse WKT point: %w", err)
	}

	p, ok := geom.(orb.Point)
	if !ok {
		return fmt.Errorf("expected Point, got %T", geom)
	}

	g.Point = p
	return nil
}

// Value implements driver.Valuer interface for writing to database
func (g GeoPoint) Value() (driver.Value, error) {
	if len(g.Point) == 0 {
		return nil, nil
	}
	return wkt.MarshalString(g.Point), nil
}

// GeoPolygon is a custom type for orb.Polygon that implements sql.Scanner and driver.Valuer
type GeoPolygon struct {
	orb.Polygon
}

// Scan implements sql.Scanner interface for reading from database
func (g *GeoPolygon) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var wktString string
	switch v := src.(type) {
	case string:
		wktString = v
	case []byte:
		wktString = string(v)
	default:
		return fmt.Errorf("unsupported type for GeoPolygon: %T", src)
	}

	// Parse WKT string to orb.Polygon
	geom, err := wkt.Unmarshal(wktString)
	if err != nil {
		return fmt.Errorf("failed to parse WKT polygon: %w", err)
	}

	p, ok := geom.(orb.Polygon)
	if !ok {
		return fmt.Errorf("expected Polygon, got %T", geom)
	}

	g.Polygon = p
	return nil
}

// Value implements driver.Valuer interface for writing to database
func (g GeoPolygon) Value() (driver.Value, error) {
	if len(g.Polygon) == 0 {
		return nil, nil
	}
	return wkt.MarshalString(g.Polygon), nil
}

// MarshalJSON implements json.Marshaler interface to convert GeoPolygon to GeoJSON format
func (g GeoPolygon) MarshalJSON() ([]byte, error) {
	if len(g.Polygon) == 0 {
		return json.Marshal(map[string]interface{}{
			"type":        "Polygon",
			"coordinates": [][][]float64{},
		})
	}

	coords := make([][][]float64, len(g.Polygon))
	for i, ring := range g.Polygon {
		coords[i] = make([][]float64, len(ring))
		for j, point := range ring {
			coords[i][j] = []float64{point.Lon(), point.Lat()} // GeoJSON format: [lng, lat]
		}
	}

	return json.Marshal(map[string]interface{}{
		"type":        "Polygon",
		"coordinates": coords,
	})
}
