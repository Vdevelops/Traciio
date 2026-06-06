package route_optimization

import "time"

// OptimizedRouteResponse represents optimized route response DTO.
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

// OptimizeRouteRequest represents route optimization request DTO.
type OptimizeRouteRequest struct {
	RouteName        *string    `json:"route_name" binding:"omitempty,max=255"`
	StartLocation    *Location  `json:"start_location" binding:"required"`
	Waypoints        []Waypoint `json:"waypoints" binding:"required,min=1,max=25,dive"`
	OptimizationType string     `json:"optimization_type" binding:"omitempty,oneof=distance duration"`
	StartTime        *time.Time `json:"start_time,omitempty"`
}

// CalculateDistanceRequest represents distance calculation request DTO.
type CalculateDistanceRequest struct {
	Origin      Location `json:"origin" binding:"required"`
	Destination Location `json:"destination" binding:"required"`
}

// CalculateDistanceResponse represents distance calculation response DTO.
type CalculateDistanceResponse struct {
	Distance          float64 `json:"distance"`
	DistanceFormatted string  `json:"distance_formatted"`
	Duration          int     `json:"duration"`
	DurationFormatted string  `json:"duration_formatted"`
}

// ListRoutesRequest represents list routes query parameters.
type ListRoutesRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	UserID  string `form:"user_id" binding:"omitempty,uuid"`
}

// ExportRouteRequest represents route export request DTO.
type ExportRouteRequest struct {
	Format string `json:"format" binding:"required,oneof=gpx json"`
}
