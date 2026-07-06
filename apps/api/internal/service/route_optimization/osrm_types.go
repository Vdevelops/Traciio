package route_optimization

// OSRMRouteResponse represents OSRM route API response.
type OSRMRouteResponse struct {
	Code      string         `json:"code"`
	Routes    []OSRMRoute    `json:"routes"`
	Waypoints []OSRMWaypoint `json:"waypoints"`
}

// OSRMRoute represents a route from OSRM.
type OSRMRoute struct {
	Distance float64   `json:"distance"`
	Duration float64   `json:"duration"`
	Geometry string    `json:"geometry"`
	Legs     []OSRMLeg `json:"legs"`
}

// OSRMLeg represents a leg of the route.
type OSRMLeg struct {
	Distance float64    `json:"distance"`
	Duration float64    `json:"duration"`
	Steps    []OSRMStep `json:"steps"`
	Summary  string     `json:"summary"`
}

// OSRMStep represents a step in the route.
type OSRMStep struct {
	Distance    float64      `json:"distance"`
	Duration    float64      `json:"duration"`
	Geometry    string       `json:"geometry"`
	Maneuver    OSRMManeuver `json:"maneuver"`
	Mode        string       `json:"mode"`
	DrivingSide string       `json:"driving_side"`
}

// OSRMManeuver represents maneuver information.
type OSRMManeuver struct {
	BearingAfter  int       `json:"bearing_after"`
	BearingBefore int       `json:"bearing_before"`
	Location      []float64 `json:"location"`
	Modifier      string    `json:"modifier"`
	Type          string    `json:"type"`
	Instruction   string    `json:"instruction,omitempty"`
}

// OSRMWaypoint represents a waypoint in the route.
type OSRMWaypoint struct {
	Hint     string    `json:"hint"`
	Distance float64   `json:"distance"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

// RouteRequest represents a route request.
type RouteRequest struct {
	Coordinates  [][]float64
	Profile      string
	Overview     string
	Geometries   string
	Steps        bool
	Alternatives bool
}

// TableRequest represents an OSRM table request.
type TableRequest struct {
	Coordinates [][]float64
	Profile     string
	Annotations []string
}

// OSRMTableResponse represents OSRM table API response.
type OSRMTableResponse struct {
	Code      string      `json:"code"`
	Distances [][]float64 `json:"distances"`
	Durations [][]float64 `json:"durations"`
}
