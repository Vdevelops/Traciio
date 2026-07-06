package auth

import "time"

// LoginRequest represents login request DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents login response DTO
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	CSRFToken    string        `json:"csrf_token,omitempty"`
	ExpiresIn    int           `json:"expires_in"`
}

// UserResponse represents user response DTO (without sensitive data)
type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
