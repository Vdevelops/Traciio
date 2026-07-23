package pipeline

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PipelineStage represents a pipeline stage (e.g., Lead, Qualification, Proposal, Negotiation, Closed Won, Closed Lost)
type PipelineStage struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Order       int            `gorm:"type:integer;not null;default:0" json:"order"`
	Color       string         `gorm:"type:varchar(20);default:'#3B82F6'" json:"color"`
	IsActive    bool           `gorm:"type:boolean;default:true" json:"is_active"`
	IsWon       bool           `gorm:"type:boolean;default:false" json:"is_won"`  // True for "Closed Won"
	IsLost      bool           `gorm:"type:boolean;default:false" json:"is_lost"` // True for "Closed Lost"
	Probability int            `gorm:"type:integer;default:0" json:"probability"` // 0-100 percentage
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for PipelineStage
func (PipelineStage) TableName() string {
	return "pipeline_stages"
}

// BeforeCreate hook to generate UUID
func (ps *PipelineStage) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = uuid.New().String()
	}
	return nil
}

// Deal represents a sales deal/opportunity
type Deal struct {
	ID                string            `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title             string            `gorm:"type:varchar(255);not null;index:idx_deals_fts,type:gin,expression:to_tsvector('english'\, title || ' ' || COALESCE(description\, '') || ' ' || COALESCE(notes\, ''))" json:"title"`
	Description       string            `gorm:"type:text" json:"description"`
	AccountID         string            `gorm:"type:uuid;not null;index" json:"account_id"`
	Account           *AccountRef       `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ContactID         *string           `gorm:"type:uuid;index" json:"contact_id"` // Optional contact
	Contact           *ContactRef       `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
	StageID           string            `gorm:"type:uuid;not null;index" json:"stage_id"`
	Stage             *PipelineStage    `gorm:"foreignKey:StageID" json:"stage,omitempty"`
	Value             int64             `gorm:"type:bigint;not null;default:0" json:"value"` // Deal value in smallest currency unit (sen)
	ProductItems      []DealProductItem `gorm:"foreignKey:DealID" json:"product_items,omitempty"`
	Probability       int               `gorm:"type:integer;default:0" json:"probability"` // 0-100 percentage
	ExpectedCloseDate *time.Time        `gorm:"type:date" json:"expected_close_date"`
	ActualCloseDate   *time.Time        `gorm:"type:date;index" json:"actual_close_date"`
	AssignedTo        *string           `gorm:"type:uuid;index" json:"assigned_to"` // Sales rep ID
	AssignedUser      *UserRef          `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
	LeadID            *string           `gorm:"type:uuid;index" json:"lead_id,omitempty"`                     // Optional: track source lead
	BrickID           *string           `gorm:"type:uuid;index" json:"brick_id,omitempty"`                    // Brick/Area assignment
	Status            string            `gorm:"type:varchar(20);not null;default:'open';index" json:"status"` // open, won, lost
	Source            string            `gorm:"type:varchar(100);index" json:"source"`                        // e.g., "website", "referral", "cold_call"
	// BANT fields carried from Lead
	BudgetConfirmed       bool           `gorm:"type:boolean;default:false;index" json:"budget_confirmed"`
	AuthorityConfirmed    bool           `gorm:"type:boolean;default:false;index" json:"authority_confirmed"`
	NeedConfirmed         bool           `gorm:"type:boolean;default:false;index" json:"need_confirmed"`
	TimelineConfirmed     bool           `gorm:"type:boolean;default:false;index" json:"timeline_confirmed"`
	QualificationSnapshot datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"qualification_snapshot,omitempty"`
	CloseReason           string         `gorm:"type:text" json:"close_reason,omitempty"` // Reason for won/lost
	Notes                 string         `gorm:"type:text" json:"notes"`
	CreatedBy             string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Deal
func (Deal) TableName() string {
	return "deals"
}

// BeforeCreate hook to generate UUID
func (d *Deal) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	// Default status if not set
	if d.Status == "" {
		d.Status = "open"
	}
	return nil
}

// AfterCreate hook to handle lead conversion data mapping
func (d *Deal) AfterCreate(tx *gorm.DB) error {
	// If this deal was created from a lead, link all existing visit reports and activities
	if d.LeadID != nil && *d.LeadID != "" {
		// Update visit reports linked to this lead to also point to this deal
		tx.Table("visit_reports").
			Where("lead_id = ? AND deal_id IS NULL", *d.LeadID).
			Update("deal_id", d.ID)

		// Update activities linked to this lead to also point to this deal
		tx.Table("activities").
			Where("lead_id = ? AND deal_id IS NULL", *d.LeadID).
			Update("deal_id", d.ID)

		tx.Table("tasks").
			Where("lead_id = ? AND deal_id IS NULL", *d.LeadID).
			Update("deal_id", d.ID)
	}
	return nil
}

// AccountRef represents account reference in deal
type AccountRef struct {
	ID   string `gorm:"type:uuid;primary_key" json:"id"`
	Name string `json:"name"`
}

// TableName specifies the table name for AccountRef
func (AccountRef) TableName() string {
	return "accounts"
}

// ContactRef represents contact reference in deal
type ContactRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// TableName specifies the table name for ContactRef
func (ContactRef) TableName() string {
	return "contacts"
}

// UserRef represents user reference in deal
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

// DealProductItem represents product items attached to a deal/opportunity.
// NOTE: we keep a snapshot of product fields (name, sku, unit_price) for historical accuracy.
type DealProductItem struct {
	ID                      string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DealID                  string         `gorm:"type:uuid;not null;index" json:"deal_id"`
	Deal                    *Deal          `gorm:"foreignKey:DealID" json:"-"`
	ProductID               string         `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName             string         `gorm:"type:varchar(200);not null" json:"product_name"`
	ProductSKU              string         `gorm:"type:varchar(100);not null" json:"product_sku"`
	UnitPrice               int64          `gorm:"type:bigint;not null;default:0" json:"unit_price"`
	UnitPriceFormatted      string         `gorm:"-" json:"unit_price_formatted,omitempty"`
	UnitCost                int64          `gorm:"type:bigint;not null;default:0" json:"unit_cost"`
	UnitCostFormatted       string         `gorm:"-" json:"unit_cost_formatted,omitempty"`
	Quantity                int            `gorm:"type:integer;not null;default:1" json:"quantity"`
	DiscountAmount          int64          `gorm:"type:bigint;not null;default:0" json:"discount_amount"`
	DiscountAmountFormatted string         `gorm:"-" json:"discount_amount_formatted,omitempty"`
	Subtotal                int64          `gorm:"type:bigint;not null;default:0" json:"subtotal"`
	SubtotalFormatted       string         `gorm:"-" json:"subtotal_formatted,omitempty"`
	ProductCategoryID       *string        `gorm:"type:uuid" json:"product_category_id,omitempty"`
	ProductCategoryName     string         `gorm:"type:varchar(100);default:''" json:"product_category_name"`
	MarginAmount            int64          `gorm:"type:bigint;not null;default:0" json:"margin_amount"`
	MarginPercentage        float64        `gorm:"type:decimal(5,2);default:0" json:"margin_percentage"`
	Notes                   string         `gorm:"type:text" json:"notes"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DealProductItem) TableName() string {
	return "deal_product_items"
}

func (i *DealProductItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.CalculateMargins()
	return nil
}

func (i *DealProductItem) BeforeSave(tx *gorm.DB) error {
	i.CalculateMargins()
	return nil
}

// CalculateMargins computes margin amount and percentage
func (i *DealProductItem) CalculateMargins() {
	i.MarginAmount = i.Subtotal - (int64(i.Quantity) * i.UnitCost)
	if i.Subtotal > 0 {
		i.MarginPercentage = (float64(i.MarginAmount) / float64(i.Subtotal)) * 100
	} else {
		i.MarginPercentage = 0
	}
}

// ToDealProductItemResponse converts DealProductItem to response DTO.
func (i *DealProductItem) ToDealProductItemResponse() *DealProductItemResponse {
	marginAmount := i.Subtotal - (int64(i.Quantity) * i.UnitCost)
	resp := &DealProductItemResponse{
		ID:                  i.ID,
		DealID:              i.DealID,
		ProductID:           i.ProductID,
		ProductName:         i.ProductName,
		ProductSKU:          i.ProductSKU,
		UnitPrice:           i.UnitPrice,
		UnitCost:            i.UnitCost,
		Quantity:            i.Quantity,
		DiscountAmount:      i.DiscountAmount,
		Subtotal:            i.Subtotal,
		MarginAmount:        marginAmount,
		ProductCategoryID:   i.ProductCategoryID,
		ProductCategoryName: i.ProductCategoryName,
		Notes:               i.Notes,
		CreatedAt:           i.CreatedAt,
		UpdatedAt:           i.UpdatedAt,
	}
	resp.UnitPriceFormatted = formatCurrency(i.UnitPrice)
	resp.UnitCostFormatted = formatCurrency(i.UnitCost)
	resp.DiscountAmountFormatted = formatCurrency(i.DiscountAmount)
	resp.SubtotalFormatted = formatCurrency(i.Subtotal)
	resp.MarginAmountFormatted = formatCurrency(marginAmount)
	return resp
}
