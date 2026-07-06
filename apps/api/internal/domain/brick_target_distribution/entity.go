package brick_target_distribution

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BrickTargetDistribution represents a brick target distribution entity
type BrickTargetDistribution struct {
	ID                string                        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BrickID           string                        `gorm:"type:uuid;not null;index" json:"brick_id"`
	Brick             *brick.Brick                  `gorm:"foreignKey:BrickID" json:"brick,omitempty"`
	BrickTargetID     string                        `gorm:"type:uuid;not null;index" json:"brick_target_id"`
	BrickTarget       *monthly_target.MonthlyTarget `gorm:"foreignKey:BrickTargetID" json:"brick_target,omitempty"`
	SalesUserID       string                        `gorm:"type:uuid;not null;index" json:"sales_user_id"`
	SalesUser         *user.User                    `gorm:"foreignKey:SalesUserID" json:"sales_user,omitempty"`
	DistributedAmount int64                         `gorm:"type:bigint;not null;default:0" json:"distributed_amount"`
	DistributedBy     string                        `gorm:"type:uuid;not null" json:"distributed_by"`
	DistributedByUser *user.User                    `gorm:"foreignKey:DistributedBy" json:"distributed_by_user,omitempty"`
	DistributedAt     time.Time                     `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"distributed_at"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
	DeletedAt         gorm.DeletedAt                `gorm:"index" json:"-"`
}

// TableName specifies the table name for BrickTargetDistribution
func (BrickTargetDistribution) TableName() string {
	return "brick_target_distributions"
}

// BeforeCreate hook to generate UUID
func (btd *BrickTargetDistribution) BeforeCreate(tx *gorm.DB) error {
	if btd.ID == "" {
		btd.ID = uuid.New().String()
	}
	return nil
}
