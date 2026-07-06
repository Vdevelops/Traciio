package brick

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Brick represents a brick/area entity
type Brick struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Province    string    `gorm:"type:varchar(100);not null" json:"province"`
	Regency     string    `gorm:"type:varchar(100);not null" json:"regency"`
	District    *string   `gorm:"type:varchar(100)" json:"district"`
	ManagerID   *string   `gorm:"type:uuid;index" json:"manager_id"`
	Manager     *user.User `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Brick
func (Brick) TableName() string {
	return "bricks"
}

// BeforeCreate hook to generate UUID
func (b *Brick) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// ToBrickResponse converts Brick to BrickResponse
func (b *Brick) ToBrickResponse() *BrickResponse {
	resp := &BrickResponse{
		ID:          b.ID,
		Name:        b.Name,
		Code:        b.Code,
		Description: b.Description,
		Province:    b.Province,
		Regency:     b.Regency,
		District:    b.District,
		ManagerID:   b.ManagerID,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	if b.Manager != nil {
		resp.Manager = b.Manager.ToUserResponse()
	}
	return resp
}
