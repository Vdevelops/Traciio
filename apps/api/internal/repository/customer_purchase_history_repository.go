package repository

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type customerPurchaseHistoryRepository struct {
	db *gorm.DB
}

func NewCustomerPurchaseHistoryRepository(db *gorm.DB) interfaces.CustomerPurchaseHistoryRepository {
	return &customerPurchaseHistoryRepository{db: db}
}

func (r *customerPurchaseHistoryRepository) Create(history *customer_purchase.CustomerPurchaseHistory) error {
	return r.db.Create(history).Error
}

func (r *customerPurchaseHistoryRepository) FindByAccountID(accountID string) ([]customer_purchase.CustomerPurchaseHistory, error) {
	var results []customer_purchase.CustomerPurchaseHistory
	err := r.db.Where("account_id = ?", accountID).Order("purchase_date desc").Find(&results).Error
	return results, err
}

func (r *customerPurchaseHistoryRepository) GetAnalytics(accountID string) (map[string]interface{}, error) {
	var stats struct {
		TotalSpent int64   `json:"total_spent"`
		OrderCount int64   `json:"order_count"`
		AvgOrder   float64 `json:"avg_order"`
	}

	err := r.db.Table("customer_purchase_history").
		Select("SUM(total_amount) as total_spent, COUNT(*) as order_count, AVG(total_amount) as avg_order").
		Where("account_id = ?", accountID).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_spent": stats.TotalSpent,
		"order_count": stats.OrderCount,
		"avg_order":   stats.AvgOrder,
	}, nil
}

func (r *customerPurchaseHistoryRepository) GetProductAnalytics(accountID string) ([]customer_purchase.CustomerProductAnalytics, error) {
	var results []customer_purchase.CustomerProductAnalytics

	// Aggregate from JSONB products array
	err := r.db.Raw(`
        SELECT 
            account_id,
            p->>'product_id' as product_id,
            p->>'product_name' as product_name,
            p->>'product_category_id' as product_category_id,
            p->>'product_category_name' as product_category_name,
            SUM((p->>'quantity')::numeric)::bigint as total_quantity_purchased,
            SUM((p->>'subtotal')::numeric)::bigint as total_amount_purchased,
            COUNT(*) as purchase_count,
            MIN(purchase_date) as first_purchase_date,
            MAX(purchase_date) as last_purchase_date
        FROM 
            customer_purchase_history,
            jsonb_array_elements(products) p
        WHERE 
            account_id = ?
        GROUP BY 
            account_id, product_id, product_name, product_category_id, product_category_name
        ORDER BY 
            total_amount_purchased DESC
    `, accountID).Scan(&results).Error

	return results, err
}

func (r *customerPurchaseHistoryRepository) GetSummary(accountID string) (*customer_purchase.CustomerPurchaseSummaryResponse, error) {
	var summary customer_purchase.CustomerPurchaseSummaryResponse

	err := r.db.Table("customer_purchase_history").
		Select("account_id, COUNT(*) as total_purchases, SUM(total_amount) as total_amount, SUM(total_items) as total_items, AVG(total_amount) as average_purchase_amount, MIN(purchase_date) as first_purchase_date, MAX(purchase_date) as last_purchase_date").
		Where("account_id = ?", accountID).
		Group("account_id").
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
