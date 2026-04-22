package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
)

type CustomerPurchaseHistoryRepository interface {
	Create(history *customer_purchase.CustomerPurchaseHistory) error
	FindByAccountID(accountID string) ([]customer_purchase.CustomerPurchaseHistory, error)
	GetAnalytics(accountID string) (map[string]interface{}, error)
	GetProductAnalytics(accountID string) ([]customer_purchase.CustomerProductAnalytics, error)
	GetSummary(accountID string) (*customer_purchase.CustomerPurchaseSummaryResponse, error)
}
