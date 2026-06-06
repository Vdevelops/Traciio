package area_mapping

// CaptureLocationRequest represents area capture input DTO.
type CaptureLocationRequest struct {
	VisitReportID string  `json:"visit_report_id" binding:"required"`
	CaptureType   string  `json:"capture_type" binding:"required,oneof=check_in check_out area"`
	Latitude      float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude     float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Address       string  `json:"address"`
	Accuracy      float64 `json:"accuracy"`
}

// CreateTerritoryRequest represents territory creation input DTO.
type CreateTerritoryRequest struct {
	Name        string      `json:"name" binding:"required,min=3,max=255"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates" binding:"required,min=4"`
	AssignedTo  string      `json:"assigned_to"`
	Color       string      `json:"color"`
}

// UpdateTerritoryRequest represents territory update input DTO.
type UpdateTerritoryRequest struct {
	Name        *string      `json:"name" binding:"omitempty,min=3,max=255"`
	Description *string      `json:"description"`
	Coordinates *[][]float64 `json:"coordinates" binding:"omitempty,min=4"`
	AssignedTo  *string      `json:"assigned_to"`
	Color       *string      `json:"color"`
}
