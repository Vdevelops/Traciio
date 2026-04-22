package product_analytics

import (
	"time"
)

// Repository defines interface for product analytics data access
type Repository interface {
	// Product Sales
	CreateProductSale(productSale *ProductSales) error
	GetProductSaleByID(id string) (*ProductSales, error)
	GetProductSales(filters map[string]interface{}, page, perPage int) ([]*ProductSales, int64, error)
	DeleteProductSale(id string) error

	// Top Products Analytics
	GetTopProducts(filters map[string]interface{}, limit int) ([]*TopProductResponse, error)
	
	// Product Performance
	GetProductPerformance(productID string, startDate, endDate time.Time) (*ProductPerformanceResponse, error)
	
	// Product Comparison
	GetProductComparison(productIDs []string, startDate, endDate time.Time) ([]*ProductPerformanceResponse, error)
	
	// Product Trends
	GetProductTrends(productID string, startDate, endDate time.Time, groupBy string) (*ProductTrendResponse, error)
}
