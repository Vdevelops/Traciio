package customer_purchase

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

// CustomerPurchaseResponse represents purchase history response DTO.
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

// ToResponse converts CustomerPurchaseHistory to response DTO.
func (p *CustomerPurchaseHistory) ToResponse() *CustomerPurchaseResponse {
	return &CustomerPurchaseResponse{
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
}

// CustomerProductAnalytics represents aggregated product analytics per customer.
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

// CustomerPurchaseSummaryResponse represents an account purchase summary.
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

// ListPurchaseHistoryRequest represents query parameters for listing purchase history.
type ListPurchaseHistoryRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	SalesRepID string `form:"sales_rep_id" binding:"omitempty,uuid"`
	StartDate  string `form:"start_date" binding:"omitempty"`
	EndDate    string `form:"end_date" binding:"omitempty"`
}
