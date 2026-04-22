package area_mapping

import (
	"fmt"
	"strings"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	"gorm.io/gorm"

	"github.com/gilabs/crm-healthcare/api/internal/domain/area_mapping"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new area mapping repository
func NewRepository(db *gorm.DB) area_mapping.Repository {
	return &repository{db: db}
}

// CreateAreaCapture creates a new area capture
func (r *repository) CreateAreaCapture(capture *area_mapping.AreaCapture) error {
	return r.db.Create(capture).Error
}

// GetAreaCaptureByID retrieves an area capture by ID
func (r *repository) GetAreaCaptureByID(id string) (*area_mapping.AreaCapture, error) {
	var capture area_mapping.AreaCapture
	err := r.db.Where("id = ?", id).First(&capture).Error
	if err != nil {
		return nil, err
	}
	return &capture, nil
}

// ListAreaCaptures retrieves area captures with filters
func (r *repository) ListAreaCaptures(filter area_mapping.AreaCaptureFilter) ([]area_mapping.AreaCapture, int64, error) {
	var captures []area_mapping.AreaCapture
	var total int64

	query := r.db.Model(&area_mapping.AreaCapture{})

	// Apply filters
	if filter.VisitReportID != nil {
		query = query.Where("visit_report_id = ?", *filter.VisitReportID)
	}

	if filter.CaptureType != nil {
		query = query.Where("capture_type = ?", *filter.CaptureType)
	}

	if filter.CapturedAfter != nil {
		query = query.Where("captured_at >= ?", *filter.CapturedAfter)
	}

	if filter.CapturedBefore != nil {
		query = query.Where("captured_at <= ?", *filter.CapturedBefore)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build raw SQL query with ST_AsText for geometry column
	sql := `SELECT id, visit_report_id, capture_type, ST_AsText(location) as location_wkt, address, accuracy, captured_at, created_at, updated_at 
			FROM "area_captures"`

	var conditions []string
	var args []interface{}

	if filter.VisitReportID != nil {
		conditions = append(conditions, "visit_report_id = ?")
		args = append(args, *filter.VisitReportID)
	}

	if filter.CaptureType != nil {
		conditions = append(conditions, "capture_type = ?")
		args = append(args, *filter.CaptureType)
	}

	if filter.CapturedAfter != nil {
		conditions = append(conditions, "captured_at >= ?")
		args = append(args, *filter.CapturedAfter)
	}

	if filter.CapturedBefore != nil {
		conditions = append(conditions, "captured_at <= ?")
		args = append(args, *filter.CapturedBefore)
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += " ORDER BY captured_at DESC"

	if filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PerPage, offset)
	}

	// Query rows
	type AreaCaptureRow struct {
		ID            string
		VisitReportID string
		CaptureType   string
		LocationWKT   string
		Address       *string
		Accuracy      *float64
		CapturedAt    time.Time
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	var rows []AreaCaptureRow
	err := r.db.Raw(sql, args...).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert rows to area captures with parsed locations
	captures = make([]area_mapping.AreaCapture, len(rows))
	for i, row := range rows {
		geom, err := wkt.Unmarshal(row.LocationWKT)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse location WKT: %w", err)
		}

		captures[i] = area_mapping.AreaCapture{
			ID:            row.ID,
			VisitReportID: row.VisitReportID,
			CaptureType:   row.CaptureType,
			Location:      area_mapping.GeoPoint{Point: geom.(orb.Point)},
			Address:       row.Address,
			Accuracy:      row.Accuracy,
			CapturedAt:    row.CapturedAt,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
	}

	return captures, total, nil
}

// DeleteAreaCapture deletes an area capture
func (r *repository) DeleteAreaCapture(id string) error {
	return r.db.Where("id = ?", id).Delete(&area_mapping.AreaCapture{}).Error
}

// CreateTerritory creates a new territory
func (r *repository) CreateTerritory(territory *area_mapping.Territory) error {
	return r.db.Create(territory).Error
}

// GetTerritoryByID retrieves a territory by ID
func (r *repository) GetTerritoryByID(id string) (*area_mapping.Territory, error) {
	// Use raw SQL to get territory with ST_AsText for polygon
	// This converts GEOGRAPHY to WKT string that we can parse
	type TerritoryRow struct {
		ID          string
		Name        string
		Description *string
		PolygonWKT  string `gorm:"column:polygon_wkt"`
		AssignedTo  *string
		Color       string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	var row TerritoryRow
	err := r.db.Table("territories").
		Select("id, name, description, ST_AsText(polygon) as polygon_wkt, assigned_to, color, created_at, updated_at").
		Where("id = ?", id).
		First(&row).Error

	if err != nil {
		return nil, err
	}

	// Parse polygon WKT
	geom, err := wkt.Unmarshal(row.PolygonWKT)
	if err != nil {
		return nil, fmt.Errorf("failed to parse polygon WKT: %w", err)
	}

	territory := &area_mapping.Territory{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Polygon:     area_mapping.GeoPolygon{Polygon: geom.(orb.Polygon)},
		AssignedTo:  row.AssignedTo,
		Color:       row.Color,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}

	return territory, nil
}

// ListTerritories retrieves territories with filters
func (r *repository) ListTerritories(filter area_mapping.TerritoryFilter) ([]area_mapping.Territory, int64, error) {
	var territories []area_mapping.Territory
	var total int64

	query := r.db.Model(&area_mapping.Territory{})

	// Apply filters
	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+*filter.Search+"%")
	}

	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset).Limit(filter.PerPage)
	}

	// Use raw SQL to get territories with ST_AsText for polygon
	// This converts GEOGRAPHY to WKT string that we can parse
	type TerritoryRow struct {
		ID          string
		Name        string
		Description *string
		PolygonWKT  string `gorm:"column:polygon_wkt"`
		AssignedTo  *string
		Color       string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	var rows []TerritoryRow
	err := r.db.Table("territories").
		Select("id, name, description, ST_AsText(polygon) as polygon_wkt, assigned_to, color, created_at, updated_at").
		Where(query.Statement.SQL.String()).
		Order("name ASC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert rows to territories with parsed polygons
	territories = make([]area_mapping.Territory, len(rows))
	for i, row := range rows {
		geom, err := wkt.Unmarshal(row.PolygonWKT)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse polygon WKT: %w", err)
		}

		territories[i] = area_mapping.Territory{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Polygon:     area_mapping.GeoPolygon{Polygon: geom.(orb.Polygon)},
			AssignedTo:  row.AssignedTo,
			Color:       row.Color,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}
	}

	return territories, total, nil
}

// UpdateTerritory updates a territory
func (r *repository) UpdateTerritory(territory *area_mapping.Territory) error {
	return r.db.Save(territory).Error
}

// DeleteTerritory deletes a territory
func (r *repository) DeleteTerritory(id string) error {
	return r.db.Where("id = ?", id).Delete(&area_mapping.Territory{}).Error
}

// CreateCoverageAnalysis creates a new coverage analysis
func (r *repository) CreateCoverageAnalysis(analysis *area_mapping.CoverageAnalysis) error {
	return r.db.Create(analysis).Error
}

// GetCoverageAnalysisByID retrieves a coverage analysis by ID
func (r *repository) GetCoverageAnalysisByID(id string) (*area_mapping.CoverageAnalysis, error) {
	var analysis area_mapping.CoverageAnalysis
	err := r.db.Where("id = ?", id).First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// ListCoverageAnalysis retrieves coverage analyses with filters
func (r *repository) ListCoverageAnalysis(filter area_mapping.CoverageAnalysisFilter) ([]area_mapping.CoverageAnalysis, int64, error) {
	var analyses []area_mapping.CoverageAnalysis
	var total int64

	query := r.db.Model(&area_mapping.CoverageAnalysis{})

	// Apply filters
	if filter.TerritoryID != nil {
		query = query.Where("territory_id = ?", *filter.TerritoryID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.PeriodStart != nil {
		query = query.Where("period_start >= ?", *filter.PeriodStart)
	}

	if filter.PeriodEnd != nil {
		query = query.Where("period_end <= ?", *filter.PeriodEnd)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset).Limit(filter.PerPage)
	}

	// Execute query
	if err := query.Order("analyzed_at DESC").Find(&analyses).Error; err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// GetCapturesWithinTerritory retrieves all captures within a territory polygon
func (r *repository) GetCapturesWithinTerritory(territoryID string, startDate, endDate *time.Time) ([]area_mapping.AreaCapture, error) {
	// Get territory first
	territory, err := r.GetTerritoryByID(territoryID)
	if err != nil {
		return nil, fmt.Errorf("territory not found: %w", err)
	}

	// `territory.Polygon` is a wrapper type (`GeoPolygon`) around an `orb.Polygon`.
	// The orb/wkb encoder only supports orb geometry types directly; passing the wrapper
	// causes a runtime panic: "unsupported type".
	poly := territory.Polygon.Polygon

	// IMPORTANT:
	// `AreaCapture.Location` uses `GeoPoint.Scan`, which expects WKT text.
	// When we SELECT `ac.*`, Postgres returns geography in a binary form, and Scan fails.
	// So we select `ST_AsText(ac.location)` and parse it ourselves.

	baseSQL := `
		SELECT ac.id,
			ac.visit_report_id,
			ac.capture_type,
			ST_AsText(ac.location) as location_wkt,
			ac.address,
			ac.accuracy,
			ac.captured_at,
			ac.created_at,
			ac.updated_at
		FROM area_captures ac
		WHERE ST_Within(ac.location::geometry, ST_GeogFromWKB(?::bytea)::geometry)
	`

	args := []interface{}{wkb.Value(poly)}
	clauses := []string{}

	if startDate != nil {
		clauses = append(clauses, "ac.captured_at >= ?")
		args = append(args, *startDate)
	}

	if endDate != nil {
		clauses = append(clauses, "ac.captured_at <= ?")
		args = append(args, *endDate)
	}

	if len(clauses) > 0 {
		baseSQL += " AND " + strings.Join(clauses, " AND ")
	}

	// Keep consistent ordering with other capture listings
	baseSQL += " ORDER BY ac.captured_at DESC"

	type AreaCaptureRow struct {
		ID            string
		VisitReportID string
		CaptureType   string
		LocationWKT   string `gorm:"column:location_wkt"`
		Address       *string
		Accuracy      *float64
		CapturedAt    time.Time
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	var rows []AreaCaptureRow
	if err := r.db.Raw(baseSQL, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get captures within territory: %w", err)
	}

	captures := make([]area_mapping.AreaCapture, len(rows))
	for i, row := range rows {
		geom, err := wkt.Unmarshal(row.LocationWKT)
		if err != nil {
			return nil, fmt.Errorf("failed to parse location WKT: %w", err)
		}

		p, ok := geom.(orb.Point)
		if !ok {
			return nil, fmt.Errorf("expected Point, got %T", geom)
		}

		captures[i] = area_mapping.AreaCapture{
			ID:            row.ID,
			VisitReportID: row.VisitReportID,
			CaptureType:   row.CaptureType,
			Location:      area_mapping.GeoPoint{Point: p},
			Address:       row.Address,
			Accuracy:      row.Accuracy,
			CapturedAt:    row.CapturedAt,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
	}

	return captures, nil
}

// CheckPointInTerritory checks if a point is within a specific territory
func (r *repository) CheckPointInTerritory(point orb.Point, territoryID string) (bool, error) {
	// Get territory
	territory, err := r.GetTerritoryByID(territoryID)
	if err != nil {
		return false, fmt.Errorf("territory not found: %w", err)
	}

	var result struct {
		IsWithin bool
	}

	// Use PostGIS ST_Within to check if point is inside polygon
	err = r.db.Raw(`
		SELECT ST_Within(
			ST_GeogFromWKB(?::bytea)::geometry,
			ST_GeogFromWKB(?::bytea)::geometry
		) AS is_within
	`, wkb.Value(point), wkb.Value(territory.Polygon)).Scan(&result).Error

	if err != nil {
		return false, fmt.Errorf("failed to check point in territory: %w", err)
	}

	return result.IsWithin, nil
}

// CalculateDistance calculates distance between two points in meters
func (r *repository) CalculateDistance(point1, point2 orb.Point) (float64, error) {
	var result struct {
		Distance float64
	}

	// Use PostGIS ST_Distance to calculate distance
	err := r.db.Raw(`
		SELECT ST_Distance(
			ST_GeogFromWKB(?::bytea),
			ST_GeogFromWKB(?::bytea)
		) AS distance
	`, wkb.Value(point1), wkb.Value(point2)).Scan(&result).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate distance: %w", err)
	}

	return result.Distance, nil
}

// GetHeatmapData gets aggregated location data for heatmap visualization
func (r *repository) GetHeatmapData(filter area_mapping.HeatmapFilter) ([]area_mapping.HeatmapPoint, error) {
	var heatmapPoints []area_mapping.HeatmapPoint

	// Build base query
	query := r.db.Table("area_captures")

	// Apply filters
	if filter.CaptureType != nil {
		query = query.Where("capture_type = ?", *filter.CaptureType)
	}

	if filter.StartDate != nil {
		query = query.Where("captured_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("captured_at <= ?", *filter.EndDate)
	}

	// If territory filter is provided, filter by territory
	if filter.TerritoryID != nil {
		territory, err := r.GetTerritoryByID(*filter.TerritoryID)
		if err != nil {
			return nil, fmt.Errorf("territory not found: %w", err)
		}

		// `territory.Polygon` is a wrapper type. The wkb encoder supports orb geometries,
		// so we need to pass the underlying orb.Polygon to avoid a runtime panic.
		poly := territory.Polygon.Polygon
		query = query.Where(
			"ST_Within(location::geometry, ST_GeogFromWKB(?::bytea)::geometry)",
			wkb.Value(poly),
		)
	}

	// Aggregate by lat/lng and count intensity
	err := query.Select(`
		ST_Y(location::geometry) AS lat,
		ST_X(location::geometry) AS lng,
		COUNT(*) AS intensity
	`).Group("ST_Y(location::geometry), ST_X(location::geometry)").
		Order("intensity DESC").
		Limit(1000).
		Scan(&heatmapPoints).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get heatmap data: %w", err)
	}

	return heatmapPoints, nil
}
