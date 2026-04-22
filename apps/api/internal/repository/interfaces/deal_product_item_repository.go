package interfaces

import "github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"

// DealProductItemRepository defines persistence for deal product items.
type DealProductItemRepository interface {
  CreateMany(items []pipeline.DealProductItem) error
  DeleteByDealID(dealID string) error
}
