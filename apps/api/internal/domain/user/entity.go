package user

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity
type User struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"` // Hidden from JSON
	Name       string    `gorm:"type:varchar(255);not null;index:idx_users_fts,type:gin,expression:to_tsvector('english'\\, name || ' ' || email)" json:"name"`
	AvatarURL  string    `gorm:"type:text" json:"avatar_url"`
	RoleID     string    `gorm:"type:uuid;not null;index" json:"role_id"`
	Role       *role.Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	GroupID *string   `gorm:"type:uuid;index" json:"group_id"`
	Group   *group.Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	BrickID    *string   `gorm:"type:uuid;index" json:"brick_id"`
	Status     string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// UserMonthlyTarget represents monthly target in user response (simplified to avoid circular dependency)
type UserMonthlyTarget struct {
	ID                  string    `json:"id"`
	GroupID             *string   `json:"group_id"`
	UserID              *string   `json:"user_id"`
	Year                int       `json:"year"`
	Month               int       `json:"month"`
	TargetAmount        int64     `json:"target_amount"`
	TargetAmountFormatted string  `json:"target_amount_formatted,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// UserResponse represents user response DTO (without sensitive data)
type UserResponse struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Name       string         `json:"name"`
	AvatarURL  string         `json:"avatar_url"`
	RoleID     string         `json:"role_id"`
	Role       *role.RoleResponse  `json:"role,omitempty"`
	GroupID *string         `json:"group_id"`
	Group   *group.GroupResponse `json:"group,omitempty"`
	BrickID    *string         `json:"brick_id"`
	Brick      interface{}    `json:"brick,omitempty"` // BrickResponse - using interface{} to avoid circular dependency
	MonthlyTarget *UserMonthlyTarget `json:"monthly_target,omitempty"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// ToUserResponse converts User to UserResponse
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
		// Brick will be loaded separately in service layer to avoid circular dependency
	}
	return resp
}

// CreateUserRequest represents create user request DTO
type CreateUserRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Name       string `json:"name" binding:"required,min=3"`
	RoleID     string `json:"role_id" binding:"required,uuid"`
	GroupID *string `json:"group_id" binding:"omitempty,uuid"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateUserRequest represents update user request DTO
type UpdateUserRequest struct {
	Email      string  `json:"email" binding:"omitempty,email"`
	Name       string  `json:"name" binding:"omitempty,min=3"`
	RoleID     string  `json:"role_id" binding:"omitempty,uuid"`
	GroupID *string `json:"group_id" binding:"omitempty,uuid"`
	BrickID    *string `json:"brick_id" binding:"omitempty,uuid"`
	Status     string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListUsersRequest represents list users query parameters
type ListUsersRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search        string   `form:"search" binding:"omitempty"`
	Status        string   `form:"status" binding:"omitempty,oneof=active inactive"`
	RoleID        string   `form:"role_id" binding:"omitempty,uuid"`
	GroupID       string   `form:"group_id" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	ScopedUserIDs []string `form:"-" json:"-"` // Injected by scope middleware for RBAC filtering
}

