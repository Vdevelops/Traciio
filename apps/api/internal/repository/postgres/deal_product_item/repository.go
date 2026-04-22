package deal_product_item

import (
  "github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
  "github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
  "gorm.io/gorm"
)

type repository struct {
  db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.DealProductItemRepository {
  return &repository{db: db}
}

func (r *repository) CreateMany(items []pipeline.DealProductItem) error {
  if len(items) == 0 {
    return nil
  }
  return r.db.Create(&items).Error
}

func (r *repository) DeleteByDealID(dealID string) error {
  return r.db.Where("deal_id = ?", dealID).Delete(&pipeline.DealProductItem{}).Error
}
