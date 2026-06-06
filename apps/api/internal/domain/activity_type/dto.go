package activity_type

import "time"

// ActivityTypeResponse represents activity type response DTO.
type ActivityTypeResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	Icon          string    `json:"icon"`
	BadgeColor    string    `json:"badge_color"`
	Status        string    `json:"status"`
	Order         int       `json:"order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ActivityCount int64     `json:"activity_count"`
}

// ToActivityTypeResponse converts ActivityType to ActivityTypeResponse.
func (at *ActivityType) ToActivityTypeResponse() *ActivityTypeResponse {
	return &ActivityTypeResponse{
		ID:            at.ID,
		Name:          at.Name,
		Code:          at.Code,
		Description:   at.Description,
		Icon:          at.Icon,
		BadgeColor:    at.BadgeColor,
		Status:        at.Status,
		Order:         at.Order,
		CreatedAt:     at.CreatedAt,
		UpdatedAt:     at.UpdatedAt,
		ActivityCount: at.ActivityCount,
	}
}

// CreateActivityTypeRequest represents create activity type request DTO.
type CreateActivityTypeRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary destructive outline"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	Order       int    `json:"order" binding:"omitempty"`
}

// UpdateActivityTypeRequest represents update activity type request DTO.
type UpdateActivityTypeRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	BadgeColor  string `json:"badge_color" binding:"omitempty,oneof=default secondary destructive outline"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	Order       *int   `json:"order" binding:"omitempty"`
}

// ListActivityTypesRequest represents list activity types query parameters.
type ListActivityTypesRequest struct {
	Status string `form:"status" binding:"omitempty,oneof=active inactive"`
}
