package monthly_target

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MonthlyTarget represents a monthly target entity
type MonthlyTarget struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupID *string   `gorm:"type:uuid;index" json:"group_id"`
	Group   *group.Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	UserID      *string   `gorm:"type:uuid;index" json:"user_id"`
	User        *user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BrickID     *string   `gorm:"type:uuid;index" json:"brick_id"`
	Brick       *brick.Brick `gorm:"foreignKey:BrickID" json:"brick,omitempty"`
	Year        int       `gorm:"type:integer;not null" json:"year"`
	Month       int       `gorm:"type:integer;not null" json:"month"`
	TargetAmount int64    `gorm:"type:bigint;not null;default:0" json:"target_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for MonthlyTarget
func (MonthlyTarget) TableName() string {
	return "monthly_targets"
}

// BeforeCreate hook to generate UUID
func (mt *MonthlyTarget) BeforeCreate(tx *gorm.DB) error {
	if mt.ID == "" {
		mt.ID = uuid.New().String()
	}
	return nil
}

// ToMonthlyTargetResponse converts MonthlyTarget to MonthlyTargetResponse
func (mt *MonthlyTarget) ToMonthlyTargetResponse() *MonthlyTargetResponse {
	resp := &MonthlyTargetResponse{
		ID:          mt.ID,
		Year:        mt.Year,
		Month:       mt.Month,
		TargetAmount: mt.TargetAmount,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
	}

	if mt.GroupID != nil {
		resp.GroupID = mt.GroupID
		if mt.Group != nil {
			resp.Group = mt.Group.ToGroupResponse()
		}
	}

	if mt.UserID != nil {
		resp.UserID = mt.UserID
		if mt.User != nil {
			resp.User = mt.User.ToUserResponse()
		}
	}

	if mt.BrickID != nil {
		resp.BrickID = mt.BrickID
		if mt.Brick != nil {
			resp.Brick = mt.Brick.ToBrickResponse()
		}
	}

	return resp
}
