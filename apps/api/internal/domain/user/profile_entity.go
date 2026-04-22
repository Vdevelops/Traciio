package user

import (
	"time"
)

// ProfileStats represents user profile statistics
type ProfileStats struct {
	Visits int `json:"visits"` // Visit reports count
	Deals  int `json:"deals"`  // Deals count
	Tasks  int `json:"tasks"`  // Tasks count
}

// ProfileActivity represents a profile activity item
type ProfileActivity struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // visit, call, email, task, deal
	Date        time.Time `json:"date"`
	DownloadURL string    `json:"download_url,omitempty"`
}

// ProfileTransaction represents a profile transaction
type ProfileTransaction struct {
	ID      string    `json:"id"`
	Product string    `json:"product"`
	Status  string    `json:"status"` // pending, paid, failed
	Date    time.Time `json:"date"`
	Amount  int64     `json:"amount"`
}

// ProfileConnection represents a user connection
type ProfileConnection struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	IsConnected bool `json:"is_connected"`
}

// ProfileResponse represents complete profile response
type ProfileResponse struct {
	User         *UserResponse        `json:"user"`
	Stats        *ProfileStats        `json:"stats"`
	Activities   []ProfileActivity    `json:"activities"`
	Transactions []ProfileTransaction `json:"transactions"`
}

// UpdateProfileRequest represents update profile request DTO
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"omitempty,min=3"`
}

// ChangePasswordRequest represents change password request DTO
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=1"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=1"`
}

