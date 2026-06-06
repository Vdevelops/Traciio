package account

import "time"

// AccountResponse represents account response DTO
type AccountResponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	CategoryID   string            `json:"category_id"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Address      string            `json:"address"`
	City         string            `json:"city"`
	Province     string            `json:"province"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Latitude     *float64          `json:"latitude"`
	Longitude    *float64          `json:"longitude"`
	PostalCode   string            `json:"postal_code"`
	Country      string            `json:"country"`
	Website      string            `json:"website"`
	Industry     string            `json:"industry"`
	Status       string            `json:"status"`
	AssignedTo   *string           `json:"assigned_to"`
	BrickID      *string           `json:"brick_id"`
	ContactCount int               `json:"contact_count"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// CategoryResponse represents category in account response
type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color"`
	Status      string `json:"status"`
}

// CreateAccountRequest represents create account request DTO
type CreateAccountRequest struct {
	Name       string   `json:"name" binding:"required,min=3"`
	CategoryID string   `json:"category_id" binding:"required,uuid"`
	Address    string   `json:"address" binding:"omitempty"`
	City       string   `json:"city" binding:"omitempty"`
	Province   string   `json:"province" binding:"omitempty"`
	Phone      string   `json:"phone" binding:"omitempty"`
	Email      string   `json:"email" binding:"omitempty,email"`
	Latitude   *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude  *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	PostalCode string   `json:"postal_code" binding:"omitempty,max=20"`
	Country    string   `json:"country" binding:"omitempty,max=100"`
	Website    string   `json:"website" binding:"omitempty,max=255"`
	Industry   string   `json:"industry" binding:"omitempty,max=100"`
	Status     string   `json:"status" binding:"omitempty,oneof=active inactive"`
	AssignedTo string   `json:"assigned_to" binding:"omitempty,uuid"`
	BrickID    string   `json:"brick_id" binding:"omitempty,uuid"`
}

// UpdateAccountRequest represents update account request DTO
type UpdateAccountRequest struct {
	Name       string   `json:"name" binding:"omitempty,min=3"`
	CategoryID string   `json:"category_id" binding:"omitempty,uuid"`
	Address    string   `json:"address" binding:"omitempty"`
	City       string   `json:"city" binding:"omitempty"`
	Province   string   `json:"province" binding:"omitempty"`
	Phone      string   `json:"phone" binding:"omitempty"`
	Email      string   `json:"email" binding:"omitempty,email"`
	Latitude   *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude  *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	PostalCode string   `json:"postal_code" binding:"omitempty,max=20"`
	Country    string   `json:"country" binding:"omitempty,max=100"`
	Website    string   `json:"website" binding:"omitempty,max=255"`
	Industry   string   `json:"industry" binding:"omitempty,max=100"`
	Status     string   `json:"status" binding:"omitempty,oneof=active inactive"`
	AssignedTo string   `json:"assigned_to" binding:"omitempty,uuid"`
}

// ListAccountsRequest represents list accounts query parameters
type ListAccountsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=1000"`
	Search        string   `form:"search" binding:"omitempty"`
	Status        string   `form:"status" binding:"omitempty,oneof=active inactive"`
	CategoryID    string   `form:"category_id" binding:"omitempty,uuid"`
	AssignedTo    string   `form:"assigned_to" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	ScopedUserIDs []string `form:"-" json:"-"`
}

// BBoxRequest represents bounding box query for viewport-based map loading
type BBoxRequest struct {
	North      float64 `form:"north" binding:"required,gte=-90,lte=90"`
	South      float64 `form:"south" binding:"required,gte=-90,lte=90"`
	East       float64 `form:"east" binding:"required,gte=-180,lte=180"`
	West       float64 `form:"west" binding:"required,gte=-180,lte=180"`
	Search     string  `form:"search" binding:"omitempty"`
	Status     string  `form:"status" binding:"omitempty,oneof=active inactive"`
	CategoryID string  `form:"category_id" binding:"omitempty,uuid"`
	Limit      int     `form:"limit" binding:"omitempty,min=1,max=5000"`
}
