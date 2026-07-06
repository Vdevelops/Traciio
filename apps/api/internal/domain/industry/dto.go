package industry

// Request DTOs

// CreateIndustryRequest represents request to create industry
type CreateIndustryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Order       int    `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateIndustryRequest represents request to update industry
type UpdateIndustryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Code        string `json:"code" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Order       *int   `json:"order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// ListIndustriesRequest represents request to list industries
type ListIndustriesRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search" binding:"omitempty,max=255"`
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name code order created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Response DTOs

// IndustryResponse represents industry response
type IndustryResponse struct {
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

// ListIndustriesResponse represents list response
type ListIndustriesResponse struct {
	Success bool                `json:"success"`
	Data    []*IndustryResponse `json:"data"`
	Meta    *Meta               `json:"meta"`
}

// IndustryDetailResponse represents detail response
type IndustryDetailResponse struct {
	Success bool              `json:"success"`
	Data    *IndustryResponse `json:"data"`
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
