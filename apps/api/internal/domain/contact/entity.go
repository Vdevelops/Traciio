package contact

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Contact represents a contact entity (Doctor, PIC, Manager)
type Contact struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID string         `gorm:"type:uuid;not null;index" json:"account_id"`
	Name      string         `gorm:"type:varchar(255);not null;index;index:idx_contacts_fts,type:gin,expression:to_tsvector('english'\\, name || ' ' || email || ' ' || COALESCE(position\\, ''))" json:"name"`
	RoleID    string         `gorm:"type:uuid;not null;index" json:"role_id"`
	Role      *ContactRole   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Email     string         `gorm:"type:varchar(255);index" json:"email"`
	Position  string         `gorm:"type:varchar(255)" json:"position"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ContactRole represents contact role (imported from contact_role package)
type ContactRole struct {
	ID          string `gorm:"type:uuid;primary_key" json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color"`
	Status      string `json:"status"`
}

// TableName specifies the table name for Contact
func (Contact) TableName() string {
	return "contacts"
}

// BeforeCreate hook to generate UUID
func (c *Contact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
