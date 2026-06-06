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
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Hidden from JSON
	Name      string         `gorm:"type:varchar(255);not null;index:idx_users_fts,type:gin,expression:to_tsvector('english'\\, name || ' ' || email)" json:"name"`
	AvatarURL string         `gorm:"type:text" json:"avatar_url"`
	RoleID    string         `gorm:"type:uuid;not null;index" json:"role_id"`
	Role      *role.Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	GroupID   *string        `gorm:"type:uuid;index" json:"group_id"`
	Group     *group.Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	BrickID   *string        `gorm:"type:uuid;index" json:"brick_id"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
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
