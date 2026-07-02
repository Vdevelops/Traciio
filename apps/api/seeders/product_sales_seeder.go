package seeders

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
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
			var existingCount int64
			if err := db.Model(&product_analytics.ProductSales{}).
				Where("deal_id = ? AND product_id = ? AND sales_rep_id = ? AND deleted_at IS NULL", deal.ID, item.ProductID, *deal.AssignedTo).
				Count(&existingCount).Error; err != nil {
				return err
			}
			if existingCount > 0 {
				continue
			}

			record := product_analytics.ProductSales{
				DealID:     deal.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				TotalPrice: item.Subtotal,
				SoldAt:     time.Date(soldAt.Year(), soldAt.Month(), soldAt.Day(), soldAt.Hour(), soldAt.Minute(), soldAt.Second(), soldAt.Nanosecond(), soldAt.Location()),
				SalesRepID: *deal.AssignedTo,
				CreatedAt:  soldAt,
				UpdatedAt:  soldAt,
			}

			if err := db.Create(&record).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
