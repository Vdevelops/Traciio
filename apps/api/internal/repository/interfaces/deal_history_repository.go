package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
)

// DealHistoryRepository defines the interface for deal history repository
type DealHistoryRepository interface {
	Create(dealHistory *deal_history.DealHistory) error
	FindByDealID(dealID string) ([]deal_history.DealHistory, error)
	FindByID(id string) (*deal_history.DealHistory, error)
	List(dealID string, limit int) ([]deal_history.DealHistory, error)
	Delete(id string) error
}
