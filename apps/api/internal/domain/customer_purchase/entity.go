package customer_purchase

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
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

// CustomerPurchaseResponse represents purchase history response DTO
type CustomerPurchaseResponse struct {
	ID                             string            `json:"id"`
	AccountID                      string            `json:"account_id"`
	DealID                         string            `json:"deal_id"`
	PurchaseDate                   time.Time         `json:"purchase_date"`
	PurchaseNumber                 int               `json:"purchase_number"`
	TotalAmount                    int64             `json:"total_amount"`
	TotalAmountFormatted           string            `json:"total_amount_formatted"`
	TotalItems                     int               `json:"total_items"`
	Products                       []PurchaseProduct `json:"products"`
	SalesRepID                     *string           `json:"sales_rep_id"`
	SalesRepName                   string            `json:"sales_rep_name"`
	SourceLeadID                   *string           `json:"source_lead_id,omitempty"`
	SourceType                     string            `json:"source_type"`
	CustomerLifetimeValue          int64             `json:"customer_lifetime_value"`
	CustomerLifetimeValueFormatted string            `json:"customer_lifetime_value_formatted"`
	CreatedAt                      time.Time         `json:"created_at"`
	UpdatedAt                      time.Time         `json:"updated_at"`
}

// ToResponse converts CustomerPurchaseHistory to response DTO
func (p *CustomerPurchaseHistory) ToResponse() *CustomerPurchaseResponse {
	resp := &CustomerPurchaseResponse{
		ID:                             p.ID,
		AccountID:                      p.AccountID,
		DealID:                         p.DealID,
		PurchaseDate:                   p.PurchaseDate,
		PurchaseNumber:                 p.PurchaseNumber,
		TotalAmount:                    p.TotalAmount,
		TotalAmountFormatted:           currency.FormatCurrency(p.TotalAmount),
		TotalItems:                     p.TotalItems,
		SalesRepID:                     p.SalesRepID,
		SalesRepName:                   p.SalesRepName,
		SourceLeadID:                   p.SourceLeadID,
		SourceType:                     p.SourceType,
		CustomerLifetimeValue:          p.CustomerLifetimeValue,
		CustomerLifetimeValueFormatted: currency.FormatCurrency(p.CustomerLifetimeValue),
		CreatedAt:                      p.CreatedAt,
		UpdatedAt:                      p.UpdatedAt,
	}
	return resp
}

// CustomerProductAnalytics represents aggregated product analytics per customer
type CustomerProductAnalytics struct {
	AccountID              string    `json:"account_id"`
	ProductID              string    `json:"product_id"`
	ProductName            string    `json:"product_name"`
	ProductCategoryID      *string   `json:"product_category_id,omitempty"`
	ProductCategoryName    string    `json:"product_category_name"`
	TotalQuantityPurchased int64     `json:"total_quantity_purchased"`
	TotalAmountPurchased   int64     `json:"total_amount_purchased"`
	TotalAmountFormatted   string    `json:"total_amount_formatted"`
	PurchaseCount          int64     `json:"purchase_count"`
	FirstPurchaseDate      time.Time `json:"first_purchase_date"`
	LastPurchaseDate       time.Time `json:"last_purchase_date"`
}

// CustomerPurchaseSummaryResponse represents an account's purchase summary
type CustomerPurchaseSummaryResponse struct {
	AccountID                      string     `json:"account_id"`
	TotalPurchases                 int64      `json:"total_purchases"`
	TotalAmount                    int64      `json:"total_amount"`
	TotalAmountFormatted           string     `json:"total_amount_formatted"`
	TotalItems                     int64      `json:"total_items"`
	AveragePurchaseAmount          int64      `json:"average_purchase_amount"`
	AveragePurchaseAmountFormatted string     `json:"average_purchase_amount_formatted"`
	FirstPurchaseDate              *time.Time `json:"first_purchase_date,omitempty"`
	LastPurchaseDate               *time.Time `json:"last_purchase_date,omitempty"`
}

// ListPurchaseHistoryRequest represents query parameters for listing purchase history
type ListPurchaseHistoryRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	SalesRepID string `form:"sales_rep_id" binding:"omitempty,uuid"`
	StartDate  string `form:"start_date" binding:"omitempty"`
	EndDate    string `form:"end_date" binding:"omitempty"`
}
