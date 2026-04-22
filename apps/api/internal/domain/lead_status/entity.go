package lead_status

// Request DTOs

// CreateLeadStatusRequest represents request to create lead status
type CreateLeadStatusRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Score       *int   `json:"score" binding:"required,min=0"`
	Color       string `json:"color" binding:"omitempty,max=20"`
	Order       int    `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
	IsDefault   *bool  `json:"is_default" binding:"omitempty"`
	IsConverted *bool  `json:"is_converted" binding:"omitempty"`
}

// UpdateLeadStatusRequest represents request to update lead status
type UpdateLeadStatusRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Code        string `json:"code" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Score       *int   `json:"score" binding:"omitempty,min=0"`
	Color       string `json:"color" binding:"omitempty,max=20"`
	Order       *int   `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
	IsDefault   *bool  `json:"is_default" binding:"omitempty"`
	IsConverted *bool  `json:"is_converted" binding:"omitempty"`
}

// ListLeadStatusesRequest represents request to list lead statuses
type ListLeadStatusesRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search" binding:"omitempty,max=255"`
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name code score order created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Response DTOs

// LeadStatusResponse represents lead status response
type LeadStatusResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Color       string `json:"color"`
	Order       int    `json:"order"`
	IsActive    bool   `json:"is_active"`
	IsDefault   bool   `json:"is_default"`
	IsConverted bool   `json:"is_converted"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LeadCount   int64  `json:"lead_count"`
}

// ListLeadStatusesResponse represents list response
type ListLeadStatusesResponse struct {
	Success bool                  `json:"success"`
	Data    []*LeadStatusResponse `json:"data"`
	Meta    *Meta                 `json:"meta"`
}

// LeadStatusDetailResponse represents detail response
type LeadStatusDetailResponse struct {
	Success bool                `json:"success"`
	Data    *LeadStatusResponse `json:"data"`
}

// Meta represents pagination metadata
type Meta struct {
	Pagination *Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}
