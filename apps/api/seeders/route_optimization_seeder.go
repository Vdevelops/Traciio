package seeders

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/datatypes"
)

// SeedRouteOptimization seeds initial optimized routes
func SeedRouteOptimization() error {
	// Check if routes already exist
	var count int64
	database.DB.Model(&route_optimization.OptimizedRoute{}).Count(&count)
	if count > 0 {
		log.Println("Route optimization already seeded, skipping...")
		return nil
	}

	// Get admin user
	var adminUser user.User
	if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, skipping route optimization seeding: %v", err)
		return nil
	}

	// Get accounts with location data
	var accounts []account.Account
	if err := database.DB.Limit(10).Find(&accounts).Error; err != nil {
		log.Printf("Warning: No accounts found, skipping route optimization seeding: %v", err)
		return nil
	}
	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping route optimization seeding")
		return nil
	}

	// Get completed visit reports for reference
	var visitReports []visit_report.VisitReport
	database.DB.Where("status = ?", "completed").Limit(5).Find(&visitReports)

	// Helper function to marshal waypoints
	marshalWaypoints := func(waypoints []route_optimization.Waypoint) datatypes.JSON {
		bytes, _ := json.Marshal(waypoints)
		return bytes
	}

	// Helper function to marshal optimized order
	marshalOrder := func(order []int) datatypes.JSON {
		bytes, _ := json.Marshal(order)
		return bytes
	}

	// Helper function to marshal route steps
	marshalRouteSteps := func(steps []route_optimization.RouteStep) datatypes.JSON {
		bytes, _ := json.Marshal(steps)
		return bytes
	}

	// Get current time
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)
	threeDaysAgo := now.AddDate(0, 0, -3)

	routes := []route_optimization.OptimizedRoute{}

	// Route 1: Jakarta Hospital Circuit (4 hospitals in Jakarta)
	if len(accounts) >= 4 {
		routeName1 := "Jakarta Hospital Circuit - Daily Route"
		totalDistance1 := 18.5
		totalDuration1 := 3600 // 1 hour in seconds

		waypoints1 := []route_optimization.Waypoint{
			{
				Order:     0,
				Lat:       -6.1944,
				Lng:       106.8229,
				Address:   "Jl. Salemba Raya No. 6, Jakarta Pusat",
				AccountID: &accounts[0].ID,
			},
			{
				Order:     1,
				Lat:       -6.2088,
				Lng:       106.8456,
				Address:   "Jl. Diponegoro No. 71, Jakarta Pusat",
				AccountID: &accounts[1].ID,
			},
			{
				Order:     2,
				Lat:       -6.1751,
				Lng:       106.8650,
				Address:   "Jl. Rasuna Said, Jakarta Selatan",
				AccountID: &accounts[2].ID,
			},
			{
				Order:     3,
				Lat:       -6.2297,
				Lng:       106.8371,
				Address:   "Jl. Fatmawati No. 80, Jakarta Selatan",
				AccountID: &accounts[3].ID,
			},
		}

		optimizedOrder1 := []int{0, 1, 2, 3}

		routeSteps1 := []route_optimization.RouteStep{
			{
				Distance:          4.5,
				DistanceFormatted: "4.5 km",
				Duration:          900,
				DurationFormatted: "15 menit",
				Instruction:       "Menuju Jl. Diponegoro",
				Maneuver:          "turn-right",
			},
			{
				Distance:          6.2,
				DistanceFormatted: "6.2 km",
				Duration:          1200,
				DurationFormatted: "20 menit",
				Instruction:       "Menuju Jl. Rasuna Said",
				Maneuver:          "straight",
			},
			{
				Distance:          7.8,
				DistanceFormatted: "7.8 km",
				Duration:          1500,
				DurationFormatted: "25 menit",
				Instruction:       "Menuju Jl. Fatmawati",
				Maneuver:          "turn-left",
			},
		}

		polyline1 := "encodedPolylineString1234567890" // This would be actual encoded polyline from Google Maps

		routes = append(routes, route_optimization.OptimizedRoute{
			UserID:         adminUser.ID,
			RouteName:      &routeName1,
			Waypoints:      marshalWaypoints(waypoints1),
			OptimizedOrder: marshalOrder(optimizedOrder1),
			TotalDistance:  &totalDistance1,
			TotalDuration:  &totalDuration1,
			RoutePolyline:  &polyline1,
			RouteSteps:     marshalRouteSteps(routeSteps1),
			CreatedAt:      threeDaysAgo,
			UpdatedAt:      threeDaysAgo,
		})
	}

	// Route 2: Bandung Clinic Route (3 clinics)
	if len(accounts) >= 7 {
		routeName2 := "Bandung Clinic Weekly Visit"
		totalDistance2 := 12.3
		totalDuration2 := 2700 // 45 minutes

		waypoints2 := []route_optimization.Waypoint{
			{
				Order:     0,
				Lat:       -6.9175,
				Lng:       107.6191,
				Address:   "Jl. Veteran No. 123, Bandung",
				AccountID: &accounts[4].ID,
			},
			{
				Order:     1,
				Lat:       -6.9050,
				Lng:       107.6150,
				Address:   "Jl. Cihampelas No. 45, Bandung",
				AccountID: &accounts[5].ID,
			},
			{
				Order:     2,
				Lat:       -6.8950,
				Lng:       107.6200,
				Address:   "Jl. Dago No. 78, Bandung",
				AccountID: &accounts[6].ID,
			},
		}

		optimizedOrder2 := []int{0, 2, 1} // Optimized different from input order

		routeSteps2 := []route_optimization.RouteStep{
			{
				Distance:          5.5,
				DistanceFormatted: "5.5 km",
				Duration:          1200,
				DurationFormatted: "20 menit",
				Instruction:       "Menuju Jl. Dago",
				Maneuver:          "turn-left",
			},
			{
				Distance:          6.8,
				DistanceFormatted: "6.8 km",
				Duration:          1500,
				DurationFormatted: "25 menit",
				Instruction:       "Menuju Jl. Cihampelas",
				Maneuver:          "turn-right",
			},
		}

		polyline2 := "encodedPolylineString0987654321"

		routes = append(routes, route_optimization.OptimizedRoute{
			UserID:         adminUser.ID,
			RouteName:      &routeName2,
			Waypoints:      marshalWaypoints(waypoints2),
			OptimizedOrder: marshalOrder(optimizedOrder2),
			TotalDistance:  &totalDistance2,
			TotalDuration:  &totalDuration2,
			RoutePolyline:  &polyline2,
			RouteSteps:     marshalRouteSteps(routeSteps2),
			CreatedAt:      twoDaysAgo,
			UpdatedAt:      twoDaysAgo,
		})
	}

	// Route 3: Pharmacy Route - Surabaya (5 pharmacies)
	if len(accounts) >= 10 {
		routeName3 := "Surabaya Pharmacy Distribution"
		totalDistance3 := 22.7
		totalDuration3 := 4500 // 75 minutes

		var accountID7 *string
		var accountID8 *string
		var accountID9 *string
		if len(accounts) > 7 {
			accountID7 = &accounts[7].ID
		}
		if len(accounts) > 8 {
			accountID8 = &accounts[8].ID
		}
		if len(accounts) > 9 {
			accountID9 = &accounts[9].ID
		}

		waypoints3 := []route_optimization.Waypoint{
			{
				Order:     0,
				Lat:       -7.2575,
				Lng:       112.7521,
				Address:   "Jl. Basuki Rahmat No. 100, Surabaya",
				AccountID: accountID7,
			},
			{
				Order:     1,
				Lat:       -7.2819,
				Lng:       112.7950,
				Address:   "Jl. Ahmad Yani No. 200, Surabaya",
				AccountID: accountID8,
			},
			{
				Order:     2,
				Lat:       -7.3061,
				Lng:       112.7378,
				Address:   "Jl. Raya Darmo No. 50, Surabaya",
				AccountID: accountID9,
			},
			{
				Order:     3,
				Lat:       -7.2492,
				Lng:       112.7508,
				Address:   "Jl. Pemuda No. 15, Surabaya",
				AccountID: &accounts[0].ID,
			},
			{
				Order:     4,
				Lat:       -7.2653,
				Lng:       112.7379,
				Address:   "Jl. Gubeng No. 88, Surabaya",
				AccountID: &accounts[1].ID,
			},
		}

		optimizedOrder3 := []int{0, 3, 4, 2, 1} // TSP optimized order

		routeSteps3 := []route_optimization.RouteStep{
			{
				Distance:          4.2,
				DistanceFormatted: "4.2 km",
				Duration:          840,
				DurationFormatted: "14 menit",
				Instruction:       "Menuju Jl. Pemuda",
				Maneuver:          "straight",
			},
			{
				Distance:          5.5,
				DistanceFormatted: "5.5 km",
				Duration:          1100,
				DurationFormatted: "18 menit",
				Instruction:       "Menuju Jl. Gubeng",
				Maneuver:          "turn-right",
			},
			{
				Distance:          6.3,
				DistanceFormatted: "6.3 km",
				Duration:          1260,
				DurationFormatted: "21 menit",
				Instruction:       "Menuju Jl. Raya Darmo",
				Maneuver:          "turn-left",
			},
			{
				Distance:          6.7,
				DistanceFormatted: "6.7 km",
				Duration:          1300,
				DurationFormatted: "22 menit",
				Instruction:       "Menuju Jl. Ahmad Yani",
				Maneuver:          "turn-right",
			},
		}

		polyline3 := "encodedPolylineString1122334455"

		routes = append(routes, route_optimization.OptimizedRoute{
			UserID:         adminUser.ID,
			RouteName:      &routeName3,
			Waypoints:      marshalWaypoints(waypoints3),
			OptimizedOrder: marshalOrder(optimizedOrder3),
			TotalDistance:  &totalDistance3,
			TotalDuration:  &totalDuration3,
			RoutePolyline:  &polyline3,
			RouteSteps:     marshalRouteSteps(routeSteps3),
			CreatedAt:      yesterday,
			UpdatedAt:      yesterday,
		})
	}

	// Route 4: Emergency Route (from visit reports if available)
	if len(visitReports) >= 3 {
		routeName4 := "Follow-up Visit Route"
		totalDistance4 := 15.2
		totalDuration4 := 3300 // 55 minutes

		waypoints4 := []route_optimization.Waypoint{}
		for i, vr := range visitReports[:3] {
			// Parse location from visit report
			var location visit_report.Location
			if err := json.Unmarshal(vr.CheckInLocation, &location); err == nil {
				waypoint := route_optimization.Waypoint{
					Order:         i,
					Lat:           location.Latitude,
					Lng:           location.Longitude,
					Address:       location.Address,
					VisitReportID: &vr.ID,
				}
				if vr.AccountID != nil {
					waypoint.AccountID = vr.AccountID
				}
				waypoints4 = append(waypoints4, waypoint)
			}
		}

		if len(waypoints4) >= 2 {
			optimizedOrder4 := make([]int, len(waypoints4))
			for i := range waypoints4 {
				optimizedOrder4[i] = i
			}

			routeSteps4 := []route_optimization.RouteStep{
				{
					Distance:          7.5,
					DistanceFormatted: "7.5 km",
					Duration:          1500,
					DurationFormatted: "25 menit",
					Instruction:       "Menuju lokasi pertama",
					Maneuver:          "straight",
				},
				{
					Distance:          7.7,
					DistanceFormatted: "7.7 km",
					Duration:          1800,
					DurationFormatted: "30 menit",
					Instruction:       "Menuju lokasi kedua",
					Maneuver:          "turn-left",
				},
			}

			polyline4 := "encodedPolylineString5544332211"

			routes = append(routes, route_optimization.OptimizedRoute{
				UserID:         adminUser.ID,
				RouteName:      &routeName4,
				Waypoints:      marshalWaypoints(waypoints4),
				OptimizedOrder: marshalOrder(optimizedOrder4),
				TotalDistance:  &totalDistance4,
				TotalDuration:  &totalDuration4,
				RoutePolyline:  &polyline4,
				RouteSteps:     marshalRouteSteps(routeSteps4),
				CreatedAt:      now,
				UpdatedAt:      now,
			})
		}
	}

	// Route 5: Simple 2-point route (minimal viable route)
	if len(accounts) >= 2 {
		routeName5 := "Quick Visit - 2 Locations"
		totalDistance5 := 8.3
		totalDuration5 := 1800 // 30 minutes

		waypoints5 := []route_optimization.Waypoint{
			{
				Order:     0,
				Lat:       -6.1944,
				Lng:       106.8229,
				Address:   accounts[0].Address,
				AccountID: &accounts[0].ID,
			},
			{
				Order:     1,
				Lat:       -6.2088,
				Lng:       106.8456,
				Address:   accounts[1].Address,
				AccountID: &accounts[1].ID,
			},
		}

		optimizedOrder5 := []int{0, 1}

		routeSteps5 := []route_optimization.RouteStep{
			{
				Distance:          8.3,
				DistanceFormatted: "8.3 km",
				Duration:          1800,
				DurationFormatted: "30 menit",
				Instruction:       "Menuju lokasi tujuan",
				Maneuver:          "straight",
			},
		}

		polyline5 := "encodedPolylineStringABCDEF"

		routes = append(routes, route_optimization.OptimizedRoute{
			UserID:         adminUser.ID,
			RouteName:      &routeName5,
			Waypoints:      marshalWaypoints(waypoints5),
			OptimizedOrder: marshalOrder(optimizedOrder5),
			TotalDistance:  &totalDistance5,
			TotalDuration:  &totalDuration5,
			RoutePolyline:  &polyline5,
			RouteSteps:     marshalRouteSteps(routeSteps5),
			CreatedAt:      now.Add(-6 * time.Hour),
			UpdatedAt:      now.Add(-6 * time.Hour),
		})
	}

	// Insert routes into database
	for _, route := range routes {
		if err := database.DB.Create(&route).Error; err != nil {
			log.Printf("Error creating route %v: %v", route.RouteName, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d optimized routes", len(routes))
	return nil
}
