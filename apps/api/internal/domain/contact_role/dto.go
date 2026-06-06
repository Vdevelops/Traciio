package contact_role

import "time"

// ContactRoleResponse represents contact role response DTO
type ContactRoleResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	BadgeColor   string    `json:"badge_color"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ContactCount int64     `json:"contact_count"`
}

// CreateContactRoleRequest represents create contact role request DTO
type CreateContactRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Code        string `json:"code" binding:"required,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateContactRoleRequest represents update contact role request DTO
type UpdateContactRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3"`
	Code        string `json:"code" binding:"omitempty,min=3"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary outline success warning active"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}
