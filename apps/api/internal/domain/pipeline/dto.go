package pipeline

import (
	"encoding/json"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

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

func (ps *PipelineStage) ToPipelineStageResponse() *PipelineStageResponse {
	return &PipelineStageResponse{ID: ps.ID, Name: ps.Name, Code: ps.Code, Order: ps.Order, Color: ps.Color, IsActive: ps.IsActive, IsWon: ps.IsWon, IsLost: ps.IsLost, Probability: ps.Probability, Description: ps.Description, CreatedAt: ps.CreatedAt, UpdatedAt: ps.UpdatedAt}
}

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

type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

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
	resp := &DealResponse{ID: d.ID, Title: d.Title, Description: d.Description, AccountID: d.AccountID, StageID: d.StageID, Value: computedValue, Probability: computedProbability, ExpectedCloseDate: d.ExpectedCloseDate, ActualCloseDate: d.ActualCloseDate, LeadID: d.LeadID, BrickID: d.BrickID, Status: d.Status, Source: d.Source, BudgetConfirmed: d.BudgetConfirmed, AuthorityConfirmed: d.AuthorityConfirmed, NeedConfirmed: d.NeedConfirmed, TimelineConfirmed: d.TimelineConfirmed, Notes: d.Notes, CreatedBy: d.CreatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
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
	resp.ValueFormatted = formatCurrency(computedValue)
	if d.Account != nil {
		resp.Account = &AccountRefResponse{ID: d.Account.ID, Name: d.Account.Name}
	}
	if d.Contact != nil {
		resp.Contact = &ContactRefResponse{ID: d.Contact.ID, Name: d.Contact.Name, Email: d.Contact.Email, Phone: d.Contact.Phone}
	}
	if d.Stage != nil {
		resp.Stage = d.Stage.ToPipelineStageResponse()
	}
	if d.AssignedUser != nil {
		resp.AssignedUser = &UserRefResponse{ID: d.AssignedUser.ID, Name: d.AssignedUser.Name, Email: d.AssignedUser.Email, AvatarURL: d.AssignedUser.AvatarURL}
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

func formatCurrency(amount int64) string { return currency.FormatCurrency(amount) }

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
	LeadID            *string                        `json:"lead_id" binding:"omitempty,uuid"`
	Source            string                         `json:"source" binding:"omitempty,max=100"`
	Notes             string                         `json:"notes" binding:"omitempty"`
}

type CreateDealProductItemRequest struct {
	ProductID      string `json:"product_id" binding:"required,uuid"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	UnitPrice      *int64 `json:"unit_price" binding:"omitempty,min=0"`
	DiscountAmount *int64 `json:"discount_amount" binding:"omitempty,min=0"`
	Notes          string `json:"notes" binding:"omitempty,max=500"`
}

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
	LeadID            *string                        `json:"lead_id" binding:"omitempty,uuid"`
	Status            string                         `json:"status" binding:"omitempty,oneof=open won lost"`
	Source            string                         `json:"source" binding:"omitempty,max=100"`
	Notes             string                         `json:"notes" binding:"omitempty"`
	ProductItems      []CreateDealProductItemRequest `json:"product_items" binding:"omitempty,dive"`
}

type MoveDealRequest struct {
	StageID string `json:"stage_id" binding:"required,uuid"`
}
type MoveStageRequest struct {
	ToStageID string `json:"to_stage_id" binding:"required,uuid"`
	Reason    string `json:"reason" binding:"omitempty,max=500"`
	Notes     string `json:"notes" binding:"omitempty"`
}
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
	ScopedUserIDs []string `form:"-" json:"-"`
}
type ListPipelineStagesRequest struct {
	IsActive *bool `form:"is_active" binding:"omitempty"`
}
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
type UpdateStagesOrderRequest struct {
	Stages []StageOrderItem `json:"stages" binding:"required,min=1,dive"`
}
type StageOrderItem struct {
	ID    string `json:"id" binding:"required,uuid"`
	Order int    `json:"order" binding:"required,min=0"`
}
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
type StageSummary struct {
	StageID, StageName, StageCode string
	DealCount, TotalValue         int64
	TotalValueFormatted           string
}
type ForecastResponse struct {
	Period                   ForecastPeriod `json:"period"`
	ExpectedRevenue          int64          `json:"expected_revenue"`
	ExpectedRevenueFormatted string         `json:"expected_revenue_formatted"`
	WeightedRevenue          int64          `json:"weighted_revenue"`
	WeightedRevenueFormatted string         `json:"weighted_revenue_formatted"`
	Deals                    []ForecastDeal `json:"deals"`
}
type ForecastPeriod struct {
	Type       string `json:"type"`
	Start, End time.Time
}
type ForecastDeal struct {
	ID                     string     `json:"id"`
	Title                  string     `json:"title"`
	AccountID              string     `json:"account_id"`
	AccountName            string     `json:"account_name"`
	ContactID              string     `json:"contact_id,omitempty"`
	ContactName            string     `json:"contact_name,omitempty"`
	StageName              string     `json:"stage_name"`
	Value                  int64      `json:"value"`
	ValueFormatted         string     `json:"value_formatted"`
	Probability            int        `json:"probability"`
	WeightedValue          int64      `json:"weighted_value"`
	WeightedValueFormatted string     `json:"weighted_value_formatted"`
	ExpectedCloseDate      *time.Time `json:"expected_close_date"`
}
