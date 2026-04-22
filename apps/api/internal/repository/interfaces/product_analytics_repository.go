package interfaces

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
)

// ProductAnalyticsRepository defines the interface for product analytics repository
type ProductAnalyticsRepository interface {
	// Product Sales
	CreateProductSale(productSale *product_analytics.ProductSales) error
	GetProductSaleByID(id string) (*product_analytics.ProductSales, error)
	GetProductSales(filters map[string]interface{}, page, perPage int) ([]*product_analytics.ProductSales, int64, error)
	DeleteProductSale(id string) error
	
	// Product Performance - scopedUserIDs restricts to specific sales reps (nil = global)
	GetProductPerformance(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.ProductPerformanceResponse, error)

	// Product Comparison - scopedUserIDs restricts to specific sales reps (nil = global)
	GetProductComparison(productIDs []string, startDate, endDate time.Time, scopedUserIDs []string) ([]*product_analytics.ProductPerformanceResponse, error)

	// Product Trends - scopedUserIDs restricts to specific sales reps (nil = global)
	GetProductTrends(productID string, startDate, endDate time.Time, groupBy string, scopedUserIDs []string) (*product_analytics.ProductTrendResponse, error)

	// Product List - scopedUserIDs restricts to deals owned by specific sales reps (nil = global)
	GetProductsList(startDate, endDate time.Time, search, sortBy, orderBy string, page, perPage int, scopedUserIDs []string) ([]*product_analytics.ProductListItem, int64, error)

	// User Product Sales - products sold by a specific user
	GetUserProductSales(userID string, startDate, endDate time.Time, sortBy, orderBy string, page, perPage int) ([]*product_analytics.ProductListItem, int64, error)

	// Monthly Sales Chart - scopedUserIDs restricts to deals owned by specific sales reps (nil = global)
	GetMonthlySales(startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error)
	GetProductMonthlySales(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error)
}
