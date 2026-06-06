package customer_purchase

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CustomerPurchaseHistory represents a purchase record created when a deal is Closed Won
type CustomerPurchaseHistory struct {
	ID        string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID string `gorm:"type:uuid;not null;index" json:"account_id"`
	DealID    string `gorm:"type:uuid;not null;uniqueIndex:uq_customer_purchase_deal" json:"deal_id"`

	// Purchase summary
	PurchaseDate   time.Time `gorm:"type:date;not null;default:CURRENT_DATE" json:"purchase_date"`
	PurchaseNumber int       `gorm:"type:integer;not null;default:1" json:"purchase_number"`
	TotalAmount    int64     `gorm:"type:bigint;not null;default:0" json:"total_amount"`
	TotalItems     int       `gorm:"type:integer;not null;default:0" json:"total_items"`

	// Product breakdown (JSONB for flexibility)
	Products datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"products"`

	// Sales rep snapshot
	SalesRepID   *string `gorm:"type:uuid" json:"sales_rep_id"`
	SalesRepName string  `gorm:"type:varchar(255);default:''" json:"sales_rep_name"`

	// Lead source tracking
	SourceLeadID *string `gorm:"type:uuid" json:"source_lead_id"`
	SourceType   string  `gorm:"type:varchar(50);not null;default:'pipeline'" json:"source_type"` // pipeline, direct, referral

	// Analytics
	CustomerLifetimeValue int64 `gorm:"type:bigint;not null;default:0" json:"customer_lifetime_value"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for CustomerPurchaseHistory
func (CustomerPurchaseHistory) TableName() string {
	return "customer_purchase_history"
}

// BeforeCreate hook to generate UUID
func (p *CustomerPurchaseHistory) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// PurchaseProduct represents a product item within a purchase
type PurchaseProduct struct {
	ProductID           string `json:"product_id"`
	ProductName         string `json:"product_name"`
	ProductSKU          string `json:"product_sku"`
	Quantity            int    `json:"quantity"`
	UnitPrice           int64  `json:"unit_price"`
	Subtotal            int64  `json:"subtotal"`
	ProductCategoryID   string `json:"product_category_id,omitempty"`
	ProductCategoryName string `json:"product_category_name,omitempty"`
}
