package seeders

import (
	"fmt"
	"time"

	"github.com/paulmach/orb"
	"gorm.io/gorm"

	"github.com/gilabs/crm-healthcare/api/internal/domain/area_mapping"
)

// SeedAreaMapping seeds area mapping data (territories and area captures)
func SeedAreaMapping(db *gorm.DB) error {
	fmt.Println("🌱 Seeding area mapping data...")

	// Check if data already exists
	var count int64
	db.Model(&area_mapping.Territory{}).Count(&count)
	if count > 0 {
		fmt.Println("✅ Area mapping data already exists, skipping...")
		return nil
	}

	// Create territories using raw SQL (PostGIS requires ST_GeomFromText)
	territories := createTerritories()
	for _, territory := range territories {
		if err := insertTerritory(db, &territory); err != nil {
			return fmt.Errorf("failed to create territory %s: %w", territory.Name, err)
		}
		fmt.Printf("✅ Created territory: %s\n", territory.Name)
	}

	// Create area captures
	captures := createAreaCaptures(db, territories)
	if len(captures) > 0 {
		for _, capture := range captures {
			if err := insertAreaCapture(db, &capture); err != nil {
				return fmt.Errorf("failed to create area capture: %w", err)
			}
		}
		fmt.Printf("✅ Created %d area captures\n", len(captures))
	} else {
		fmt.Println("⚠️  No visit reports found, skipping area captures")
	}

	// Create coverage analysis
	analyses := createCoverageAnalyses(territories)
	for _, analysis := range analyses {
		if err := db.Create(&analysis).Error; err != nil {
			return fmt.Errorf("failed to create coverage analysis: %w", err)
		}
	}
	fmt.Printf("✅ Created %d coverage analyses\n", len(analyses))

	fmt.Println("✅ Area mapping data seeded successfully!")
	return nil
}

// insertTerritory inserts a territory using raw SQL with PostGIS functions
func insertTerritory(db *gorm.DB, territory *area_mapping.Territory) error {
	// Convert GeoPolygon to WKT (Well-Known Text)
	wkt := polygonToWKT(territory.Polygon.Polygon)

	// Generate UUID for the territory
	var id string
	err := db.Raw(`
		INSERT INTO territories (name, description, polygon, assigned_to, color)
		VALUES ($1, $2, ST_GeogFromText($3), $4, $5)
		RETURNING id
	`, territory.Name, territory.Description, wkt, territory.AssignedTo, territory.Color).Scan(&id).Error

	if err != nil {
		return err
	}

	territory.ID = id
	return nil
}

// insertAreaCapture inserts an area capture using raw SQL with PostGIS functions
func insertAreaCapture(db *gorm.DB, capture *area_mapping.AreaCapture) error {
	// Convert GeoPoint to WKT
	wkt := fmt.Sprintf("POINT(%f %f)", capture.Location.Point.Lon(), capture.Location.Point.Lat())

	var id string
	err := db.Raw(`
		INSERT INTO area_captures (visit_report_id, capture_type, location, address, accuracy, captured_at)
		VALUES ($1, $2, ST_GeogFromText($3), $4, $5, $6)
		RETURNING id
	`, capture.VisitReportID, capture.CaptureType, wkt, capture.Address, capture.Accuracy, capture.CapturedAt).Scan(&id).Error

	if err != nil {
		return err
	}

	capture.ID = id
	return nil
}

// polygonToWKT converts orb.Polygon to WKT (Well-Known Text) format
func polygonToWKT(polygon orb.Polygon) string {
	var wkt string
	wkt = "POLYGON(("

	for i, ring := range polygon {
		if i > 0 {
			wkt += "),("
		}
		for j, point := range ring {
			if j > 0 {
				wkt += ","
			}
			wkt += fmt.Sprintf("%f %f", point.Lon(), point.Lat())
		}
	}

	wkt += "))"
	return wkt
}

// createTerritories creates sample territories with realistic polygons
func createTerritories() []area_mapping.Territory {
	// For now, leave assigned_to as nil (can be assigned later via API)
	// Or query real users from the database if needed

	return []area_mapping.Territory{
		{
			Name:        "Jakarta Pusat - Menteng",
			Description: strPtr("Wilayah Menteng dan sekitarnya"),
			Polygon: area_mapping.GeoPolygon{
				Polygon: orb.Polygon{
					{
						{106.8245, -6.1950}, // Bundaran HI
						{106.8350, -6.1950}, // Thamrin
						{106.8350, -6.2050}, // Menteng
						{106.8245, -6.2050}, // Cikini
						{106.8245, -6.1950}, // Close polygon
					},
				},
			},
			AssignedTo: nil,
			Color:      "#3B82F6",
		},
		{
			Name:        "Jakarta Selatan - Kebayoran",
			Description: strPtr("Wilayah Kebayoran Baru dan sekitarnya"),
			Polygon: area_mapping.GeoPolygon{
				Polygon: orb.Polygon{
					{
						{106.7900, -6.2400}, // Senayan
						{106.8100, -6.2400}, // Blok M
						{106.8100, -6.2600}, // Kebayoran
						{106.7900, -6.2600}, // Panglima Polim
						{106.7900, -6.2400}, // Close polygon
					},
				},
			},
			AssignedTo: nil,
			Color:      "#10B981",
		},
		{
			Name:        "Tangerang - BSD City",
			Description: strPtr("Wilayah BSD City dan sekitarnya"),
			Polygon: area_mapping.GeoPolygon{
				Polygon: orb.Polygon{
					{
						{106.6200, -6.2800}, // BSD
						{106.6500, -6.2800}, // Serpong
						{106.6500, -6.3100}, // Alam Sutera
						{106.6200, -6.3100}, // Bintaro
						{106.6200, -6.2800}, // Close polygon
					},
				},
			},
			AssignedTo: nil,
			Color:      "#F59E0B",
		},
		{
			Name:        "Bekasi - Summarecon",
			Description: strPtr("Wilayah Summarecon Bekasi dan sekitarnya"),
			Polygon: area_mapping.GeoPolygon{
				Polygon: orb.Polygon{
					{
						{106.9900, -6.2200}, // Summarecon
						{107.0200, -6.2200}, // Grand Galaxy
						{107.0200, -6.2500}, // Harapan Indah
						{106.9900, -6.2500}, // Kemang Pratama
						{106.9900, -6.2200}, // Close polygon
					},
				},
			},
			AssignedTo: nil,
			Color:      "#EF4444",
		},
		{
			Name:        "Bandung - Dago",
			Description: strPtr("Wilayah Dago dan sekitarnya"),
			Polygon: area_mapping.GeoPolygon{
				Polygon: orb.Polygon{
					{
						{107.6000, -6.8700}, // Dago Bawah
						{107.6200, -6.8700}, // Dago Atas
						{107.6200, -6.8900}, // Lembang
						{107.6000, -6.8900}, // Setiabudi
						{107.6000, -6.8700}, // Close polygon
					},
				},
			},
			AssignedTo: nil, // Unassigned
			Color:      "#8B5CF6",
		},
	}
}

// createAreaCaptures creates sample area captures within territories
func createAreaCaptures(db *gorm.DB, territories []area_mapping.Territory) []area_mapping.AreaCapture {
	captures := []area_mapping.AreaCapture{}
	now := time.Now()

	// Query for real visit report IDs from the database
	type VisitReportID struct {
		ID string
	}
	var visitReports []VisitReportID
	db.Raw("SELECT id FROM visit_reports LIMIT 10").Scan(&visitReports)

	if len(visitReports) == 0 {
		fmt.Println("⚠️  No visit reports found in database, skipping area captures creation")
		return captures
	}

	// Generate captures for each territory
	for i, territory := range territories {
		// Get center point of territory for realistic captures
		bounds := territory.Polygon.Polygon.Bound()
		centerLng := (bounds.Min.Lon() + bounds.Max.Lon()) / 2
		centerLat := (bounds.Min.Lat() + bounds.Max.Lat()) / 2

		// Use real visit report IDs (cycle through available IDs)
		vrIndex := i % len(visitReports)
		visitReportID := visitReports[vrIndex].ID

		// Create check-in capture
		captures = append(captures, area_mapping.AreaCapture{
			VisitReportID: visitReportID,
			CaptureType:   "check_in",
			Location:      area_mapping.GeoPoint{Point: orb.Point{centerLng, centerLat}},
			Address:       strPtr(fmt.Sprintf("Area %s - Check In", territory.Name)),
			Accuracy:      float64Ptr(5.0),
			CapturedAt:    now.AddDate(0, 0, -i),
		})

		// Create check-out capture (use next visit report ID if available)
		vrIndex2 := (i + 1) % len(visitReports)
		visitReportID2 := visitReports[vrIndex2].ID

		captures = append(captures, area_mapping.AreaCapture{
			VisitReportID: visitReportID2,
			CaptureType:   "check_out",
			Location:      area_mapping.GeoPoint{Point: orb.Point{centerLng + 0.001, centerLat + 0.001}},
			Address:       strPtr(fmt.Sprintf("Area %s - Check Out", territory.Name)),
			Accuracy:      float64Ptr(6.0),
			CapturedAt:    now.AddDate(0, 0, -i).Add(1 * time.Hour),
		})

		// Create area capture
		vrIndex3 := (i + 2) % len(visitReports)
		visitReportID3 := visitReports[vrIndex3].ID

		captures = append(captures, area_mapping.AreaCapture{
			VisitReportID: visitReportID3,
			CaptureType:   "area",
			Location:      area_mapping.GeoPoint{Point: orb.Point{centerLng - 0.001, centerLat - 0.001}},
			Address:       strPtr(fmt.Sprintf("Area %s - Survey", territory.Name)),
			Accuracy:      float64Ptr(10.0),
			CapturedAt:    now.AddDate(0, 0, -i-1),
		})
	}

	return captures
}

// createCoverageAnalyses creates sample coverage analyses
func createCoverageAnalyses(territories []area_mapping.Territory) []area_mapping.CoverageAnalysis {
	analyses := []area_mapping.CoverageAnalysis{}
	now := time.Now()

	for i, territory := range territories {
		if territory.AssignedTo == nil {
			continue
		}

		visitCount := 3 + i
		coveragePercent := float64(65.0 + float64(i*5))

		analyses = append(analyses, area_mapping.CoverageAnalysis{
			TerritoryID:     &territory.ID,
			UserID:          territory.AssignedTo,
			PeriodStart:     now.AddDate(0, -1, 0),
			PeriodEnd:       now,
			VisitCount:      visitCount,
			CoveragePercent: &coveragePercent,
			AnalyzedAt:      now,
		})
	}

	return analyses
}

// Helper functions
func strPtr(s string) *string {
	return &s
}
