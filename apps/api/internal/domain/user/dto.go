package user

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

// UserMonthlyTarget represents monthly target in user response.
type UserMonthlyTarget struct {
	ID                    string    `json:"id"`
	GroupID               *string   `json:"group_id"`
	UserID                *string   `json:"user_id"`
	Year                  int       `json:"year"`
	Month                 int       `json:"month"`
	TargetAmount          int64     `json:"target_amount"`
	TargetAmountFormatted string    `json:"target_amount_formatted,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// UserResponse represents user response DTO without sensitive data.
type UserResponse struct {
	ID            string               `json:"id"`
	Email         string               `json:"email"`
	Name          string               `json:"name"`
	AvatarURL     string               `json:"avatar_url"`
	RoleID        string               `json:"role_id"`
	Role          *role.RoleResponse   `json:"role,omitempty"`
	GroupID       *string              `json:"group_id"`
	Group         *group.GroupResponse `json:"group,omitempty"`
	BrickID       *string              `json:"brick_id"`
	Brick         interface{}          `json:"brick,omitempty"`
	MonthlyTarget *UserMonthlyTarget   `json:"monthly_target,omitempty"`
	Status        string               `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// ToUserResponse converts User to UserResponse.
func (u *User) ToUserResponse() *UserResponse {
	resp := &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		RoleID:    u.RoleID,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Role != nil {
		resp.Role = u.Role.ToRoleResponse()
	}
	if u.GroupID != nil {
		resp.GroupID = u.GroupID
		if u.Group != nil {
			resp.Group = u.Group.ToGroupResponse()
		}
	}
	if u.BrickID != nil {
		resp.BrickID = u.BrickID
	}
	return resp
}

// CreateUserRequest represents create user request DTO.
type CreateUserRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     string  `json:"name" binding:"required,min=3"`
	RoleID   string  `json:"role_id" binding:"required,uuid"`
	GroupID  *string `json:"group_id" binding:"omitempty,uuid"`
	BrickID  *string `json:"brick_id" binding:"omitempty,uuid"`
	Status   string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateUserRequest represents update user request DTO.
type UpdateUserRequest struct {
	Email   string  `json:"email" binding:"omitempty,email"`
	Name    string  `json:"name" binding:"omitempty,min=3"`
	RoleID  string  `json:"role_id" binding:"omitempty,uuid"`
	GroupID *string `json:"group_id" binding:"omitempty,uuid"`
	BrickID *string `json:"brick_id" binding:"omitempty,uuid"`
	Status  string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListUsersRequest represents list users query parameters.
type ListUsersRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search        string   `form:"search" binding:"omitempty"`
	Status        string   `form:"status" binding:"omitempty,oneof=active inactive"`
	RoleID        string   `form:"role_id" binding:"omitempty,uuid"`
	GroupID       string   `form:"group_id" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	ScopedUserIDs []string `form:"-" json:"-"`
}

// ProfileStats represents user profile statistics.
type ProfileStats struct {
	Visits int `json:"visits"`
	Deals  int `json:"deals"`
	Tasks  int `json:"tasks"`
}

// ProfileActivity represents a profile activity item.
type ProfileActivity struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	DownloadURL string    `json:"download_url,omitempty"`
}

// ProfileTransaction represents a profile transaction.
type ProfileTransaction struct {
	ID      string    `json:"id"`
	Product string    `json:"product"`
	Status  string    `json:"status"`
	Date    time.Time `json:"date"`
	Amount  int64     `json:"amount"`
}

// ProfileConnection represents a user connection.
type ProfileConnection struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	IsConnected bool   `json:"is_connected"`
}

// ProfileResponse represents complete profile response.
type ProfileResponse struct {
	User         *UserResponse        `json:"user"`
	Stats        *ProfileStats        `json:"stats"`
	Activities   []ProfileActivity    `json:"activities"`
	Transactions []ProfileTransaction `json:"transactions"`
}

// UpdateProfileRequest represents update profile request DTO.
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"omitempty,min=3"`
}

// ChangePasswordRequest represents change password request DTO.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=1"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=1"`
}

// GetSettingsSummaryRequest represents request parameters for settings summary.
type GetSettingsSummaryRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

// SettingsStats represents extended statistics for settings page.
type SettingsStats struct {
	Visits                    int     `json:"visits"`
	Deals                     int     `json:"deals"`
	Tasks                     int     `json:"tasks"`
	TotalRevenue              int64   `json:"total_revenue"`
	DealsWon                  int     `json:"deals_won"`
	DealsLost                 int     `json:"deals_lost"`
	DealsOpen                 int     `json:"deals_open"`
	TotalRevenueFormatted     string  `json:"total_revenue_formatted"`
	ConversionRate            float64 `json:"conversion_rate"`
	AverageDealValue          int64   `json:"average_deal_value"`
	AverageDealValueFormatted string  `json:"average_deal_value_formatted"`
}

// SettingsSummaryResponse represents the complete settings summary response.
type SettingsSummaryResponse struct {
	User         *UserResponse        `json:"user"`
	Stats        *SettingsStats       `json:"stats"`
	Activities   []ProfileActivity    `json:"activities"`
	Transactions []ProfileTransaction `json:"transactions"`
}
