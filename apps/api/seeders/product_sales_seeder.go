package seeders

import (
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	"gorm.io/gorm"
)

// SeedProductSales backfills product_sales from won pipeline deals.
// This keeps the optional analytics table synchronized with pipeline data instead
// of generating synthetic sales volumes that drift from deals/deal_product_items.
func SeedProductSales() error {
	db := database.DB

	var wonDeals []pipeline.Deal
	if err := db.
		Preload("ProductItems").
		Where("status = ? AND deleted_at IS NULL", "won").
		Find(&wonDeals).Error; err != nil {
		return err
	}

	for _, deal := range wonDeals {
		if deal.AssignedTo == nil || len(deal.ProductItems) == 0 {
			continue
		}

		soldAt := deal.CreatedAt
		if deal.ActualCloseDate != nil {
			soldAt = *deal.ActualCloseDate
		}

		for _, item := range deal.ProductItems {
			var record product_analytics.ProductSales
			err := db.Where("deal_id = ? AND product_id = ? AND sales_rep_id = ? AND deleted_at IS NULL", deal.ID, item.ProductID, *deal.AssignedTo).
				First(&record).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			record.DealID = deal.ID
			record.ProductID = item.ProductID
			record.Quantity = item.Quantity
			record.UnitPrice = item.UnitPrice
			record.TotalPrice = item.Subtotal
			record.SoldAt = time.Date(soldAt.Year(), soldAt.Month(), soldAt.Day(), soldAt.Hour(), soldAt.Minute(), soldAt.Second(), soldAt.Nanosecond(), soldAt.Location())
			record.SalesRepID = *deal.AssignedTo
			if record.ID == "" {
				record.CreatedAt = soldAt
			}
			record.UpdatedAt = soldAt

			if err := db.Save(&record).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
