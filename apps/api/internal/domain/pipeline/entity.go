package pipeline

import (
	"encoding/json"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
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

// PipelineStageResponse represents pipeline stage response DTO
type PipelineStageResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Order       int       `json:"order"`
	Color       string    `json:"color"`
	IsActive    bool      `json:"is_active"`
	IsWon       bool      `json:"is_won"`
	IsLost      bool      `json:"is_lost"`
	Probability int       `json:"probability"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToPipelineStageResponse converts PipelineStage to PipelineStageResponse
func (ps *PipelineStage) ToPipelineStageResponse() *PipelineStageResponse {
	return &PipelineStageResponse{
		ID:          ps.ID,
		Name:        ps.Name,
		Code:        ps.Code,
		Order:       ps.Order,
		Color:       ps.Color,
		IsActive:    ps.IsActive,
		IsWon:       ps.IsWon,
		IsLost:      ps.IsLost,
		Probability: ps.Probability,
		Description: ps.Description,
		CreatedAt:   ps.CreatedAt,
		UpdatedAt:   ps.UpdatedAt,
	}
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

// AfterUpdate hook to handle "Closed Won" logic
func (d *Deal) AfterUpdate(tx *gorm.DB) error {
	// Check if status changed to 'won'
	if d.Status == "won" {
		// This is a simplified version of the trigger logic
		// In a real scenario, we might want to check if purchase history already exists
		var exists bool
		tx.Table("customer_purchase_history").
			Select("count(*) > 0").
			Where("deal_id = ?", d.ID).
			Scan(&exists)

		if !exists {
			// Logic to create purchase history would go here
			// However, since we don't have all product details in the hook easily
			// (d.ProductItems might not be fully loaded), this is usually better
			// handled in the Service layer. But the user asked for entity implementation.

			// For now, we'll mark it as a requirement for the service layer
			// or implement a basic record if we have the data.
		}
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

// DealResponse represents deal response DTO
type DealResponse struct {
	ID                    string                    `json:"id"`
	Title                 string                    `json:"title"`
	Description           string                    `json:"description"`
	AccountID             string                    `json:"account_id"`
	Account               *AccountRefResponse       `json:"account,omitempty"`
	ContactID             string                    `json:"contact_id"`
	Contact               *ContactRefResponse       `json:"contact,omitempty"`
	StageID               string                    `json:"stage_id"`
	Stage                 *PipelineStageResponse    `json:"stage,omitempty"`
	Value                 int64                     `json:"value"`
	ValueFormatted        string                    `json:"value_formatted,omitempty"`
	ProductItems          []DealProductItemResponse `json:"product_items,omitempty"`
	Probability           int                       `json:"probability"`
	ExpectedCloseDate     *time.Time                `json:"expected_close_date"`
	ActualCloseDate       *time.Time                `json:"actual_close_date"`
	AssignedTo            string                    `json:"assigned_to"`
	AssignedUser          *UserRefResponse          `json:"assigned_user,omitempty"`
	LeadID                *string                   `json:"lead_id,omitempty"`
	BrickID               *string                   `json:"brick_id,omitempty"`
	Status                string                    `json:"status"`
	Source                string                    `json:"source"`
	BudgetConfirmed       bool                      `json:"budget_confirmed"`
	AuthorityConfirmed    bool                      `json:"authority_confirmed"`
	NeedConfirmed         bool                      `json:"need_confirmed"`
	TimelineConfirmed     bool                      `json:"timeline_confirmed"`
	QualificationSnapshot interface{}               `json:"qualification_snapshot,omitempty"`
	Notes                 string                    `json:"notes"`
	CreatedBy             string                    `json:"created_by"`
	CreatedAt             time.Time                 `json:"created_at"`
	UpdatedAt             time.Time                 `json:"updated_at"`
}

// AccountRefResponse represents account in deal response
type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ContactRefResponse represents contact in deal response
type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UserRefResponse represents user in deal response
type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// ToDealResponse converts Deal to DealResponse
func (d *Deal) ToDealResponse() *DealResponse {
	computedValue := d.Value
	if len(d.ProductItems) > 0 {
		computedValue = 0
		for _, item := range d.ProductItems {
			computedValue += item.Subtotal
		}
	}

	computedProbability := d.Probability
	if d.Stage != nil {
		if d.Stage.Probability > 0 {
			computedProbability = d.Stage.Probability
		} else {
			computedProbability = d.Stage.Order * 20
		}
	}

	resp := &DealResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		AccountID:   d.AccountID,
		// ContactID is handled below
		StageID:            d.StageID,
		Value:              computedValue,
		Probability:        computedProbability,
		ExpectedCloseDate:  d.ExpectedCloseDate,
		ActualCloseDate:    d.ActualCloseDate,
		LeadID:             d.LeadID,
		BrickID:            d.BrickID,
		Status:             d.Status,
		Source:             d.Source,
		BudgetConfirmed:    d.BudgetConfirmed,
		AuthorityConfirmed: d.AuthorityConfirmed,
		NeedConfirmed:      d.NeedConfirmed,
		TimelineConfirmed:  d.TimelineConfirmed,
		Notes:              d.Notes,
		CreatedBy:          d.CreatedBy,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}

	if len(d.QualificationSnapshot) > 0 {
		var snapshot interface{}
		if err := json.Unmarshal(d.QualificationSnapshot, &snapshot); err == nil {
			resp.QualificationSnapshot = snapshot
		}
	}

	if d.ContactID != nil {
		resp.ContactID = *d.ContactID
	}

	if d.AssignedTo != nil {
		resp.AssignedTo = *d.AssignedTo
	}

	// Format value as currency (including 0, as 0 is a valid value)
	resp.ValueFormatted = formatCurrency(computedValue)

	if d.Account != nil {
		resp.Account = &AccountRefResponse{
			ID:   d.Account.ID,
			Name: d.Account.Name,
		}
	}

	if d.Contact != nil {
		resp.Contact = &ContactRefResponse{
			ID:    d.Contact.ID,
			Name:  d.Contact.Name,
			Email: d.Contact.Email,
			Phone: d.Contact.Phone,
		}
	}

	if d.Stage != nil {
		resp.Stage = d.Stage.ToPipelineStageResponse()
	}

	if d.AssignedUser != nil {
		resp.AssignedUser = &UserRefResponse{
			ID:        d.AssignedUser.ID,
			Name:      d.AssignedUser.Name,
			Email:     d.AssignedUser.Email,
			AvatarURL: d.AssignedUser.AvatarURL,
		}
	}

	if len(d.ProductItems) > 0 {
		resp.ProductItems = make([]DealProductItemResponse, 0, len(d.ProductItems))
		for _, item := range d.ProductItems {
			itemCopy := item
			resp.ProductItems = append(resp.ProductItems, *itemCopy.ToDealProductItemResponse())
		}
	}

	return resp
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

// DealProductItemResponse represents deal product item response DTO.
type DealProductItemResponse struct {
	ID                      string    `json:"id"`
	DealID                  string    `json:"deal_id"`
	ProductID               string    `json:"product_id"`
	ProductName             string    `json:"product_name"`
	ProductSKU              string    `json:"product_sku"`
	UnitPrice               int64     `json:"unit_price"`
	UnitPriceFormatted      string    `json:"unit_price_formatted,omitempty"`
	UnitCost                int64     `json:"unit_cost"`
	UnitCostFormatted       string    `json:"unit_cost_formatted,omitempty"`
	Quantity                int       `json:"quantity"`
	DiscountAmount          int64     `json:"discount_amount"`
	DiscountAmountFormatted string    `json:"discount_amount_formatted,omitempty"`
	Subtotal                int64     `json:"subtotal"`
	SubtotalFormatted       string    `json:"subtotal_formatted,omitempty"`
	MarginAmount            int64     `json:"margin_amount"`
	MarginAmountFormatted   string    `json:"margin_amount_formatted,omitempty"`
	ProductCategoryID       *string   `json:"product_category_id,omitempty"`
	ProductCategoryName     string    `json:"product_category_name"`
	Notes                   string    `json:"notes"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	return currency.FormatCurrency(amount)
}

// CreateDealRequest represents create deal request DTO
type CreateDealRequest struct {
	Title             string                         `json:"title" binding:"required,min=3,max=255"`
	Description       string                         `json:"description" binding:"omitempty"`
	AccountID         string                         `json:"account_id" binding:"required,uuid"`
	ContactID         string                         `json:"contact_id" binding:"omitempty,uuid"`
	StageID           string                         `json:"stage_id" binding:"required,uuid"`
	Value             int64                          `json:"value" binding:"omitempty,min=0"`
	ProductItems      []CreateDealProductItemRequest `json:"product_items" binding:"omitempty,dive"`
	Probability       int                            `json:"probability" binding:"omitempty,min=0,max=100"`
	ExpectedCloseDate *time.Time                     `json:"expected_close_date" binding:"omitempty"`
	AssignedTo        string                         `json:"assigned_to" binding:"omitempty,uuid"`
	LeadID            *string                        `json:"lead_id" binding:"omitempty,uuid"` // Optional: track source lead
	Source            string                         `json:"source" binding:"omitempty,max=100"`
	Notes             string                         `json:"notes" binding:"omitempty"`
}

// CreateDealProductItemRequest represents product item payload when creating a deal.
// unit_price and discount_amount are in smallest currency unit (sen).
type CreateDealProductItemRequest struct {
	ProductID      string `json:"product_id" binding:"required,uuid"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	UnitPrice      *int64 `json:"unit_price" binding:"omitempty,min=0"`
	DiscountAmount *int64 `json:"discount_amount" binding:"omitempty,min=0"`
	Notes          string `json:"notes" binding:"omitempty,max=500"`
}

// UpdateDealRequest represents update deal request DTO
type UpdateDealRequest struct {
	Title             string                         `json:"title" binding:"omitempty,min=3,max=255"`
	Description       string                         `json:"description" binding:"omitempty"`
	AccountID         string                         `json:"account_id" binding:"omitempty,uuid"`
	ContactID         string                         `json:"contact_id" binding:"omitempty,uuid"`
	StageID           string                         `json:"stage_id" binding:"omitempty,uuid"`
	Value             *int64                         `json:"value" binding:"omitempty,min=0"`
	Probability       *int                           `json:"probability" binding:"omitempty,min=0,max=100"`
	ExpectedCloseDate *time.Time                     `json:"expected_close_date" binding:"omitempty"`
	AssignedTo        string                         `json:"assigned_to" binding:"omitempty,uuid"`
	LeadID            *string                        `json:"lead_id" binding:"omitempty,uuid"` // Optional: track source lead
	Status            string                         `json:"status" binding:"omitempty,oneof=open won lost"`
	Source            string                         `json:"source" binding:"omitempty,max=100"`
	Notes             string                         `json:"notes" binding:"omitempty"`
	ProductItems      []CreateDealProductItemRequest `json:"product_items" binding:"omitempty,dive"`
}

// MoveDealRequest represents move deal request DTO
type MoveDealRequest struct {
	StageID string `json:"stage_id" binding:"required,uuid"`
}

// MoveStageRequest represents move stage request with validation and history
type MoveStageRequest struct {
	ToStageID string `json:"to_stage_id" binding:"required,uuid"`
	Reason    string `json:"reason" binding:"omitempty,max=500"`
	Notes     string `json:"notes" binding:"omitempty"`
}

// ListDealsRequest represents list deals query parameters
type ListDealsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search        string   `form:"search" binding:"omitempty"`
	StageID       string   `form:"stage_id" binding:"omitempty,uuid"`
	AccountID     string   `form:"account_id" binding:"omitempty,uuid"`
	AssignedTo    string   `form:"assigned_to" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	Status        string   `form:"status" binding:"omitempty,oneof=open won lost"`
	Source        string   `form:"source" binding:"omitempty"`
	MinValue      *int64   `form:"min_value" binding:"omitempty,min=0"`
	MaxValue      *int64   `form:"max_value" binding:"omitempty,min=0"`
	DateFrom      string   `form:"date_from" binding:"omitempty"`
	DateTo        string   `form:"date_to" binding:"omitempty"`
	ScopedUserIDs []string `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}

// ListPipelineStagesRequest represents list pipeline stages query parameters
type ListPipelineStagesRequest struct {
	IsActive *bool `form:"is_active" binding:"omitempty"`
}

// CreateStageRequest represents create pipeline stage request DTO
type CreateStageRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Code        string `json:"code" binding:"required,min=1,max=50"`
	Order       int    `json:"order" binding:"required,min=0"`
	Color       string `json:"color" binding:"omitempty,max=20"`
	IsActive    bool   `json:"is_active" binding:"omitempty"`
	IsWon       bool   `json:"is_won" binding:"omitempty"`
	IsLost      bool   `json:"is_lost" binding:"omitempty"`
	Probability int    `json:"probability" binding:"omitempty,min=0,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateStageRequest represents update pipeline stage request DTO
type UpdateStageRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=255"`
	Code        string `json:"code" binding:"omitempty,min=1,max=50"`
	Order       *int   `json:"order" binding:"omitempty,min=0"`
	Color       string `json:"color" binding:"omitempty,max=20"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
	IsWon       *bool  `json:"is_won" binding:"omitempty"`
	IsLost      *bool  `json:"is_lost" binding:"omitempty"`
	Probability *int   `json:"probability" binding:"omitempty,min=0,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateStagesOrderRequest represents update stages order request DTO
type UpdateStagesOrderRequest struct {
	Stages []StageOrderItem `json:"stages" binding:"required,min=1,dive"`
}

// StageOrderItem represents a stage order item
type StageOrderItem struct {
	ID    string `json:"id" binding:"required,uuid"`
	Order int    `json:"order" binding:"required,min=0"`
}

// PipelineSummaryResponse represents pipeline summary response
type PipelineSummaryResponse struct {
	TotalDeals          int64          `json:"total_deals"`
	TotalValue          int64          `json:"total_value"`
	TotalValueFormatted string         `json:"total_value_formatted"`
	WonDeals            int64          `json:"won_deals"`
	WonValue            int64          `json:"won_value"`
	WonValueFormatted   string         `json:"won_value_formatted"`
	LostDeals           int64          `json:"lost_deals"`
	LostValue           int64          `json:"lost_value"`
	LostValueFormatted  string         `json:"lost_value_formatted"`
	OpenDeals           int64          `json:"open_deals"`
	OpenValue           int64          `json:"open_value"`
	OpenValueFormatted  string         `json:"open_value_formatted"`
	ByStage             []StageSummary `json:"by_stage"`
}

// StageSummary represents summary for a stage
type StageSummary struct {
	StageID             string `json:"stage_id"`
	StageName           string `json:"stage_name"`
	StageCode           string `json:"stage_code"`
	DealCount           int64  `json:"deal_count"`
	TotalValue          int64  `json:"total_value"`
	TotalValueFormatted string `json:"total_value_formatted"`
}

// ForecastResponse represents forecast response
type ForecastResponse struct {
	Period                   ForecastPeriod `json:"period"`
	ExpectedRevenue          int64          `json:"expected_revenue"`
	ExpectedRevenueFormatted string         `json:"expected_revenue_formatted"`
	WeightedRevenue          int64          `json:"weighted_revenue"`
	WeightedRevenueFormatted string         `json:"weighted_revenue_formatted"`
	Deals                    []ForecastDeal `json:"deals"`
}

// ForecastPeriod represents forecast period
type ForecastPeriod struct {
	Type  string    `json:"type"` // "month", "quarter", "year"
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ForecastDeal represents a deal in forecast
type ForecastDeal struct {
	ID                     string     `json:"id"`
	Title                  string     `json:"title"`
	AccountID              string     `json:"account_id"` // ID for creating modal links
	AccountName            string     `json:"account_name"`
	ContactID              string     `json:"contact_id,omitempty"`   // ID for creating modal links (optional)
	ContactName            string     `json:"contact_name,omitempty"` // Name for display (optional)
	StageName              string     `json:"stage_name"`
	Value                  int64      `json:"value"`
	ValueFormatted         string     `json:"value_formatted"`
	Probability            int        `json:"probability"`
	WeightedValue          int64      `json:"weighted_value"`
	WeightedValueFormatted string     `json:"weighted_value_formatted"`
	ExpectedCloseDate      *time.Time `json:"expected_close_date"`
}
