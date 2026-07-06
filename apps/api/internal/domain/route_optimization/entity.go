package route_optimization

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OptimizedRoute represents an optimized route entity
type OptimizedRoute struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         string         `gorm:"type:uuid;not null;index" json:"user_id"`
	RouteName      *string        `gorm:"type:varchar(255)" json:"route_name,omitempty"`
	Waypoints      datatypes.JSON `gorm:"type:jsonb;not null" json:"waypoints"`               // Array of Waypoint
	OptimizedOrder datatypes.JSON `gorm:"type:jsonb;not null" json:"optimized_order"`         // Array of indices
	TotalDistance  *float64       `gorm:"type:decimal(10,2)" json:"total_distance,omitempty"` // in km
	TotalDuration  *int           `gorm:"type:integer" json:"total_duration,omitempty"`       // in seconds
	RoutePolyline  *string        `gorm:"type:text" json:"route_polyline,omitempty"`          // Encoded polyline
	RouteSteps     datatypes.JSON `gorm:"type:jsonb" json:"route_steps,omitempty"`            // Array of RouteStep
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (for preloading)
	User interface{} `gorm:"-" json:"user,omitempty"`
}

// TableName specifies the table name for OptimizedRoute
func (OptimizedRoute) TableName() string {
	return "optimized_routes"
}

// BeforeCreate hook to generate UUID
func (or *OptimizedRoute) BeforeCreate(tx *gorm.DB) error {
	if or.ID == "" {
		or.ID = uuid.New().String()
	}
	return nil
}

// Waypoint represents a waypoint in a route
type Waypoint struct {
	Order         int     `json:"order"`
	Lat           float64 `json:"lat" binding:"required"`
	Lng           float64 `json:"lng" binding:"required"`
	Address       string  `json:"address,omitempty"`
	AccountID     *string `json:"account_id,omitempty"`
	AccountName   *string `json:"account_name,omitempty"`
	ContactID     *string `json:"contact_id,omitempty"`
	ContactName   *string `json:"contact_name,omitempty"`
	VisitReportID *string `json:"visit_report_id,omitempty"`
	// Time window constraints (optional)
	EarliestArrival *time.Time `json:"earliest_arrival,omitempty"` // Earliest time customer is available
	LatestArrival   *time.Time `json:"latest_arrival,omitempty"`   // Latest time customer is available
	ServiceDuration *int       `json:"service_duration,omitempty"` // Service duration in minutes
	Priority        *int       `json:"priority,omitempty"`         // Priority: 1 (highest) to 5 (lowest), default: 3
}

// Location represents a geographic location
type Location struct {
	Lat     float64 `json:"lat" binding:"required"`
	Lng     float64 `json:"lng" binding:"required"`
	Address string  `json:"address,omitempty"`
}

// RouteStep represents a step in the route
type RouteStep struct {
	Step              int      `json:"step"`
	Distance          float64  `json:"distance"`                     // in km
	DistanceFormatted string   `json:"distance_formatted,omitempty"` // e.g., "2.5 km"
	Duration          int      `json:"duration"`                     // in seconds
	DurationFormatted string   `json:"duration_formatted,omitempty"` // e.g., "10 menit"
	Instruction       string   `json:"instruction,omitempty"`
	Polyline          string   `json:"polyline,omitempty"` // Encoded polyline segment
	Maneuver          string   `json:"maneuver,omitempty"` // turn-left, turn-right, straight, etc.
	StartLocation     Location `json:"start_location"`
	EndLocation       Location `json:"end_location"`
}
