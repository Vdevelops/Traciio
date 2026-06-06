package lead_source

// Request DTOs

// CreateLeadSourceRequest represents request to create lead source
type CreateLeadSourceRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Order       int    `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateLeadSourceRequest represents request to update lead source
type UpdateLeadSourceRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Code        string `json:"code" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Order       *int   `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// ListLeadSourcesRequest represents request to list lead sources
type ListLeadSourcesRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search" binding:"omitempty,max=255"`
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name code order created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Response DTOs

// LeadSourceResponse represents lead source response
type LeadSourceResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LeadCount   int64  `json:"lead_count"`
}

// ListLeadSourcesResponse represents list response
type ListLeadSourcesResponse struct {
	Success bool                  `json:"success"`
	Data    []*LeadSourceResponse `json:"data"`
	Meta    *Meta                 `json:"meta"`
}

// LeadSourceDetailResponse represents detail response
type LeadSourceDetailResponse struct {
	Success bool                `json:"success"`
	Data    *LeadSourceResponse `json:"data"`
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
