package seeders

import (
	"errors"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	lead "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"gorm.io/gorm"
)

// SeedDeals creates pipeline data that is explicitly linked back to seeded
// leads and whose value is derived from synchronized deal product items.
func SeedDeals() error {
	ctx, err := loadSeedFixtureContext()
	if err != nil {
		return err
	}

	records := buildCRMSeedRecords(ctx)
	seededDeals := 0
	for _, record := range records {
		if record.DealStageCode == "" {
			continue
		}

		stage, ok := ctx.StagesByCode[record.DealStageCode]
		if !ok {
			return errors.New("missing pipeline stage for CRM sync seed: " + record.DealStageCode)
		}
		if record.Account == nil {
			return errors.New("missing account for CRM sync deal seed: " + record.Key)
		}

		var sourceLead lead.Lead
		if err := database.DB.Where("email = ?", record.LeadEmail).First(&sourceLead).Error; err != nil {
			return err
		}

		value := int64(0)
		items := make([]pipeline.DealProductItem, 0, len(record.Products))
		purchaseItems := make([]customer_purchase.PurchaseProduct, 0, len(record.Products))
		for index, productItem := range record.Products {
			quantity := 1
			if index < len(record.ProductInterests) && record.ProductInterests[index].Quantity > 0 {
				quantity = record.ProductInterests[index].Quantity
			}
			subtotal := productItem.Price * int64(quantity)
			value += subtotal

			categoryID := stringPtr(productItem.CategoryID)
			categoryName := ""
			if productItem.Category != nil {
				categoryName = productItem.Category.Name
			}

			items = append(items, pipeline.DealProductItem{
				ProductID:           productItem.ID,
				ProductName:         productItem.Name,
				ProductSKU:          productItem.SKU,
				UnitPrice:           productItem.Price,
				UnitCost:            productItem.Cost,
				Quantity:            quantity,
				DiscountAmount:      0,
				Subtotal:            subtotal,
				ProductCategoryID:   categoryID,
				ProductCategoryName: categoryName,
				Notes:               "Synced from lead product interests [" + crmSeedSource + "]",
			})
			purchaseItems = append(purchaseItems, customer_purchase.PurchaseProduct{
				ProductID:           productItem.ID,
				ProductName:         productItem.Name,
				ProductSKU:          productItem.SKU,
				Quantity:            quantity,
				UnitPrice:           productItem.Price,
				Subtotal:            subtotal,
				ProductCategoryName: categoryName,
			})
		}
		if value == 0 {
			value = record.EstimatedValue
		}

		qualificationSnapshot := mustJSON(map[string]any{
			"seed_source":             crmSeedSource,
			"budget_target_amount":    record.BudgetAmount,
			"budget_target_currency":  "IDR",
			"budget_confirmed":        true,
			"authority_target_person": record.AuthorityPerson,
			"authority_target_role":   record.JobTitle,
			"authority_confirmed":     true,
			"need_target_products":    record.NeedProducts,
			"need_priority_level":     record.NeedPriorityLevel,
			"need_confirmed":          len(record.NeedProducts) > 0,
			"timeline_target_date":    record.ExpectedCloseDate,
			"timeline_confirmed":      true,
			"qualification_status":    "qualified",
		})

		var entity pipeline.Deal
		err = database.DB.Where("lead_id = ? AND title = ?", sourceLead.ID, record.DealTitle).First(&entity).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		entity.Title = record.DealTitle
		entity.Description = record.DealDescription
		entity.AccountID = record.Account.ID
		if record.Contact != nil {
			entity.ContactID = stringPtr(record.Contact.ID)
		} else {
			entity.ContactID = nil
		}
		entity.StageID = stage.ID
		entity.Value = value
		entity.Probability = stage.Probability
		entity.ExpectedCloseDate = &record.ExpectedCloseDate
		entity.AssignedTo = stringPtr(record.Owner.ID)
		entity.LeadID = stringPtr(sourceLead.ID)
		entity.BrickID = record.Owner.BrickID
		entity.Status = "open"
		if stage.IsWon {
			entity.Status = "won"
		}
		if stage.IsLost {
			entity.Status = "lost"
		}
		entity.Source = record.LeadSource
		entity.BudgetConfirmed = true
		entity.AuthorityConfirmed = true
		entity.NeedConfirmed = len(record.NeedProducts) > 0
		entity.TimelineConfirmed = true
		entity.QualificationSnapshot = qualificationSnapshot
		entity.CloseReason = record.DealCloseReason
		entity.Notes = "CRM synced pipeline seed [" + crmSeedSource + "]"
		if entity.CreatedBy == "" {
			entity.CreatedBy = ctx.AdminUser.ID
		}

		if stage.IsWon || stage.IsLost {
			closeDate := record.VisitDate.Add(48 * time.Hour)
			entity.ActualCloseDate = &closeDate
		} else {
			entity.ActualCloseDate = nil
		}

		if entity.ID == "" {
			if err := database.DB.Create(&entity).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&entity).Error; err != nil {
				return err
			}
			if err := database.DB.Where("deal_id = ?", entity.ID).Delete(&pipeline.DealProductItem{}).Error; err != nil {
				return err
			}
		}

		for index := range items {
			items[index].DealID = entity.ID
		}
		if len(items) > 0 {
			if err := database.DB.Create(&items).Error; err != nil {
				return err
			}
		}

		sourceLead.AccountID = stringPtr(record.Account.ID)
		if record.Contact != nil {
			sourceLead.ContactID = stringPtr(record.Contact.ID)
		}
		sourceLead.OpportunityID = stringPtr(entity.ID)
		sourceLead.ConvertedPipelineID = stringPtr(entity.ID)
		sourceLead.ConvertedBy = stringPtr(ctx.AdminUser.ID)
		convertedAt := record.VisitDate.Add(24 * time.Hour)
		sourceLead.ConvertedAt = &convertedAt
		sourceLead.LeadStatus = "converted"
		if statusID, ok := ctx.LeadStatusIDs["converted"]; ok {
			sourceLead.LeadStatusID = stringPtr(statusID)
		}
		sourceLead.ConversionMetadata = mustJSON(map[string]any{
			"seed_source": crmSeedSource,
			"deal_id":     entity.ID,
			"stage_code":  record.DealStageCode,
			"deal_value":  entity.Value,
		})
		if err := database.DB.Save(&sourceLead).Error; err != nil {
			return err
		}

		if entity.Status == "won" {
			var purchase customer_purchase.CustomerPurchaseHistory
			err = database.DB.Where("deal_id = ?", entity.ID).First(&purchase).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			purchase.AccountID = entity.AccountID
			purchase.DealID = entity.ID
			if entity.ActualCloseDate != nil {
				purchase.PurchaseDate = *entity.ActualCloseDate
			} else {
				purchase.PurchaseDate = time.Now()
			}
			purchase.TotalAmount = entity.Value
			purchase.TotalItems = len(purchaseItems)
			purchase.Products = mustJSON(purchaseItems)
			purchase.SalesRepID = entity.AssignedTo
			purchase.SalesRepName = record.Owner.Name
			purchase.SourceLeadID = entity.LeadID
			purchase.SourceType = "pipeline"
			purchase.CustomerLifetimeValue = entity.Value

			if purchase.ID == "" {
				if err := database.DB.Create(&purchase).Error; err != nil {
					return err
				}
			} else {
				if err := database.DB.Save(&purchase).Error; err != nil {
					return err
				}
			}
		} else {
			if err := database.DB.Where("deal_id = ?", entity.ID).Delete(&customer_purchase.CustomerPurchaseHistory{}).Error; err != nil {
				return err
			}
		}

		seededDeals++
	}

	log.Printf("Seeded %d CRM-synced deals with linked product items", seededDeals)
	return nil
}
