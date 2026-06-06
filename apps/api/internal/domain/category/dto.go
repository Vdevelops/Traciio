package category

import "time"

// CategoryResponse represents category response DTO
type CategoryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	BadgeColor   string    `json:"badge_color"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	AccountCount int64     `json:"account_count"`
}

// CreateCategoryRequest represents create category request DTO
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Code        string `json:"code" binding:"required,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateCategoryRequest represents update category request DTO
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3"`
	Code        string `json:"code" binding:"omitempty,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}
