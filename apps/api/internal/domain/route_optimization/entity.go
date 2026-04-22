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
	Order         int        `json:"order"`
	Lat           float64    `json:"lat" binding:"required"`
	Lng           float64    `json:"lng" binding:"required"`
	Address       string     `json:"address,omitempty"`
	AccountID     *string    `json:"account_id,omitempty"`
	AccountName   *string    `json:"account_name,omitempty"`
	ContactID     *string    `json:"contact_id,omitempty"`
	ContactName   *string    `json:"contact_name,omitempty"`
	VisitReportID *string    `json:"visit_report_id,omitempty"`
	// Time window constraints (optional)
	EarliestArrival *time.Time `json:"earliest_arrival,omitempty"` // Earliest time customer is available
	LatestArrival   *time.Time `json:"latest_arrival,omitempty"`   // Latest time customer is available
	ServiceDuration *int       `json:"service_duration,omitempty"` // Service duration in minutes
	Priority        *int       `json:"priority,omitempty"`        // Priority: 1 (highest) to 5 (lowest), default: 3
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

// OptimizedRouteResponse represents optimized route response DTO
type OptimizedRouteResponse struct {
	ID                     string      `json:"id"`
	RouteName              *string     `json:"route_name,omitempty"`
	UserID                 string      `json:"user_id"`
	Waypoints              []Waypoint  `json:"waypoints"`
	OptimizedOrder         []int       `json:"optimized_order"`
	TotalDistance          *float64    `json:"total_distance,omitempty"`
	TotalDistanceFormatted string      `json:"total_distance_formatted,omitempty"`
	TotalDuration          *int        `json:"total_duration,omitempty"`
	TotalDurationFormatted string      `json:"total_duration_formatted,omitempty"`
	RoutePolyline          *string     `json:"route_polyline,omitempty"`
	RouteSteps             []RouteStep `json:"route_steps,omitempty"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
	User                   interface{} `json:"user,omitempty"`
}

// OptimizeRouteRequest represents route optimization request DTO
type OptimizeRouteRequest struct {
	RouteName        *string    `json:"route_name" binding:"omitempty,max=255"`
	StartLocation    *Location  `json:"start_location" binding:"required"`
	Waypoints        []Waypoint `json:"waypoints" binding:"required,min=1,max=25,dive"`
	OptimizationType string     `json:"optimization_type" binding:"omitempty,oneof=distance duration"` // distance or duration
	StartTime        *time.Time `json:"start_time,omitempty"` // When the route starts (for time window calculations)
}

// CalculateDistanceRequest represents distance calculation request DTO
type CalculateDistanceRequest struct {
	Origin      Location `json:"origin" binding:"required"`
	Destination Location `json:"destination" binding:"required"`
}

// CalculateDistanceResponse represents distance calculation response DTO
type CalculateDistanceResponse struct {
	Distance          float64 `json:"distance"`           // in km
	DistanceFormatted string  `json:"distance_formatted"` // e.g., "5.2 km"
	Duration          int     `json:"duration"`           // in seconds
	DurationFormatted string  `json:"duration_formatted"` // e.g., "20 menit"
}

// ListRoutesRequest represents list routes query parameters
type ListRoutesRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	UserID  string `form:"user_id" binding:"omitempty,uuid"`
}

// ExportRouteRequest represents route export request DTO
type ExportRouteRequest struct {
	Format string `json:"format" binding:"required,oneof=gpx json"` // gpx or json
}
