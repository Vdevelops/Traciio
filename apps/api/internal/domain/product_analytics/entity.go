package product_analytics

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductSales represents product sales tracking table
type ProductSales struct {
	ID         string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DealID     string         `gorm:"type:uuid;not null;index" json:"deal_id"`
	ProductID  string         `gorm:"type:uuid;not null;index" json:"product_id"`
	Product    *ProductRef    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int            `gorm:"type:integer;not null;default:1" json:"quantity"`
	UnitPrice  int64          `gorm:"type:bigint;not null;default:0" json:"unit_price"` // Stored in smallest currency unit (sen)
	TotalPrice int64          `gorm:"type:bigint;not null;default:0" json:"total_price"`
	SoldAt     time.Time      `gorm:"type:timestamp;not null;index" json:"sold_at"`
	SalesRepID string         `gorm:"type:uuid;not null;index" json:"sales_rep_id"`
	SalesRep   *UserRef       `gorm:"foreignKey:SalesRepID" json:"sales_rep,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ProductSales
func (ProductSales) TableName() string {
	return "product_sales"
}

// BeforeCreate hook to generate UUID
func (ps *ProductSales) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = uuid.New().String()
	}
	return nil
}

// ProductRef represents product reference in analytics
type ProductRef struct {
	ID         string `gorm:"type:uuid;primary_key" json:"id"`
	Name       string `gorm:"type:varchar(200)" json:"name"`
	SKU        string `gorm:"type:varchar(100)" json:"sku"`
	CategoryID string `gorm:"type:uuid" json:"category_id"`
}

// TableName specifies the table name for ProductRef
func (ProductRef) TableName() string {
	return "products"
}

// UserRef represents user reference in analytics
type UserRef struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// TableName specifies the table name for UserRef
func (UserRef) TableName() string {
	return "users"
}
