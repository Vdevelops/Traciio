package product_analytics

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
)

// logContext is a helper type for structured logging
type logContext map[string]interface{}

// Cache TTLs for analytics following enterprise best practices
const (
	// Analytics cache: 5-10 minutes (aggregated data, expensive queries)
	TTLAnalyticsDefault   = 5 * time.Minute
	TTLAnalyticsPerformance = 5 * time.Minute
	TTLAnalyticsTrends      = 10 * time.Minute
	TTLAnalyticsMonthlySales = 10 * time.Minute

	// Cache key prefixes
	CacheKeyProductPerformance = "analytics:product:performance:%s:%s:%s"
	CacheKeyProductComparison  = "analytics:product:comparison:%s:%s:%s"
	CacheKeyProductTrends      = "analytics:product:trends:%s:%s:%s:%s"
	CacheKeyProductsList       = "analytics:products:list:%s:%s:%s:%s:%d"
	CacheKeyMonthlySales       = "analytics:monthly_sales:%d"
	CacheKeyProductMonthlySales = "analytics:product:monthly_sales:%s:%d"
	CacheKeyUserProductSales   = "analytics:user:product_sales:%s:%s:%s:%s:%s:%d:%d"
)

var (
	ErrProductSaleNotFound = errors.New("product sale not found")
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidDateRange    = errors.New("invalid date range")
	ErrInvalidMetric       = errors.New("invalid metric type")
)

type Service struct {
	productAnalyticsRepo interfaces.ProductAnalyticsRepository
}

func NewService(productAnalyticsRepo interfaces.ProductAnalyticsRepository) *Service {
	return &Service{
		productAnalyticsRepo: productAnalyticsRepo,
	}
}

// CreateProductSale creates a new product sale record
func (s *Service) CreateProductSale(productSale *product_analytics.ProductSales) error {
	return s.productAnalyticsRepo.CreateProductSale(productSale)
}

// GetProductSaleByID returns a product sale by ID
func (s *Service) GetProductSaleByID(id string) (*product_analytics.ProductSales, error) {
	productSale, err := s.productAnalyticsRepo.GetProductSaleByID(id)
	if err != nil {
		return nil, ErrProductSaleNotFound
	}
	return productSale, nil
}

// GetProductSales returns product sales with filters and pagination
func (s *Service) GetProductSales(filters map[string]interface{}, page, perPage int) ([]*product_analytics.ProductSales, int64, error) {
	return s.productAnalyticsRepo.GetProductSales(filters, page, perPage)
}

// DeleteProductSale deletes a product sale
func (s *Service) DeleteProductSale(id string) error {
	return s.productAnalyticsRepo.DeleteProductSale(id)
}

// GetProductPerformance returns detailed performance metrics for a product with caching
func (s *Service) GetProductPerformance(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.ProductPerformanceResponse, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	// Try cache first
	scopeKey := strings.Join(scopedUserIDs, ",")
	cacheKey := fmt.Sprintf(CacheKeyProductPerformance+":scope:%s", productID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), scopeKey)
	if cache.Client != nil && cache.Client.IsEnabled() {
		var cachedResult product_analytics.ProductPerformanceResponse
		if found, _ := cache.Client.Get(cacheKey, &cachedResult); found {
			return &cachedResult, nil
		}
	}

	result, err := s.productAnalyticsRepo.GetProductPerformance(productID, startDate, endDate, scopedUserIDs)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cache.Client != nil && cache.Client.IsEnabled() && result != nil {
		_ = cache.Client.Set(cacheKey, result, TTLAnalyticsPerformance)
	}

	return result, nil
}

// GetProductComparison returns performance comparison for multiple products
func (s *Service) GetProductComparison(productIDs []string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.ProductComparisonResponse, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	if len(productIDs) == 0 {
		return &product_analytics.ProductComparisonResponse{
			Products: []product_analytics.ProductPerformanceResponse{},
		}, nil
	}

	products, err := s.productAnalyticsRepo.GetProductComparison(productIDs, startDate, endDate, scopedUserIDs)
	if err != nil {
		return nil, err
	}

	// Convert []*ProductPerformanceResponse to []ProductPerformanceResponse
	productsList := make([]product_analytics.ProductPerformanceResponse, len(products))
	for i, p := range products {
		productsList[i] = *p
	}

	return &product_analytics.ProductComparisonResponse{
		Products: productsList,
	}, nil
}

// GetProductTrends returns product sales trends over time
func (s *Service) GetProductTrends(productID string, startDate, endDate time.Time, groupBy string, scopedUserIDs []string) (*product_analytics.ProductTrendResponse, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	// Validate groupBy
	validGroupBy := map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"year":  true,
	}
	if groupBy == "" {
		groupBy = "month" // Default
	} else if !validGroupBy[groupBy] {
		groupBy = "month"
	}

	return s.productAnalyticsRepo.GetProductTrends(productID, startDate, endDate, groupBy, scopedUserIDs)
}

// GetProductsList returns a list of all products with their sales analytics
func (s *Service) GetProductsList(startDate, endDate time.Time, search, sortBy, orderBy string, page, perPage int, scopedUserIDs []string) ([]*product_analytics.ProductListItem, int64, error) {
	// Validate sortBy
	validSortBy := map[string]bool{
		"total_sold": true,
		"revenue":    true,
		"profit":     true,
		"name":       true,
	}
	if sortBy == "" || !validSortBy[sortBy] {
		sortBy = "total_sold" // Default
	}

	// Validate orderBy
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "desc" // Default
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	result, total, err := s.productAnalyticsRepo.GetProductsList(startDate, endDate, search, sortBy, orderBy, page, perPage, scopedUserIDs)
	if err != nil {
		return nil, 0, err
	}

	for _, item := range result {
		if item == nil {
			continue
		}

		performance, perfErr := s.productAnalyticsRepo.GetProductPerformance(
			item.ProductID,
			startDate,
			endDate,
			scopedUserIDs,
		)
		if perfErr != nil || performance == nil {
			continue
		}

		item.GrowthRate = performance.GrowthRate
	}
	
	// Debug logging for data verification
	if len(result) > 0 {
		logger.LogInfo("ProductsList Query Result", logContext{
			"period_start":        startDate,
			"period_end":          endDate,
			"search":              search,
			"page":                page,
			"per_page":            perPage,
			"total_products":      total,
			"fetched_products":    len(result),
			"first_product":       result[0].ProductName,
			"first_total_sold":    result[0].TotalSold,
			"first_total_revenue": result[0].TotalRevenue,
			"first_sales_count":   result[0].SalesCount,
		})
	}
	
	return result, total, nil
}

// GetUserProductSales returns a list of products sold by a specific user with their sales analytics
func (s *Service) GetUserProductSales(userID string, startDate, endDate time.Time, sortBy, orderBy string, page, perPage int) ([]*product_analytics.ProductListItem, int64, error) {
	// Validate sortBy
	validSortBy := map[string]bool{
		"total_sold": true,
		"revenue":    true,
		"profit":     true,
		"name":       true,
	}
	if sortBy == "" || !validSortBy[sortBy] {
		sortBy = "total_sold" // Default
	}

	// Validate orderBy
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "desc" // Default
	}

	// Validate pagination
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100 // Max limit
	}

	return s.productAnalyticsRepo.GetUserProductSales(userID, startDate, endDate, sortBy, orderBy, page, perPage)
}

// GetMonthlySales returns monthly sales data for a specific range with caching
func (s *Service) GetMonthlySales(startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error) {
	// Try cache first
	scopeKey := strings.Join(scopedUserIDs, ",")
	cacheKey := fmt.Sprintf("analytics:monthly_sales:%s:%s:scope:%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), scopeKey)
	if cache.Client != nil && cache.Client.IsEnabled() {
		var cachedResult product_analytics.MonthlySalesResponse
		if found, _ := cache.Client.Get(cacheKey, &cachedResult); found {
			return &cachedResult, nil
		}
	}

	result, err := s.productAnalyticsRepo.GetMonthlySales(startDate, endDate, scopedUserIDs)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cache.Client != nil && cache.Client.IsEnabled() && result != nil {
		_ = cache.Client.Set(cacheKey, result, TTLAnalyticsMonthlySales)
	}

	return result, nil
}

// GetProductMonthlySales returns monthly sales data for a specific product and range with caching
func (s *Service) GetProductMonthlySales(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error) {
	// Try cache first
	scopeKey := strings.Join(scopedUserIDs, ",")
	cacheKey := fmt.Sprintf("analytics:product:monthly_sales:%s:%s:%s:scope:%s", productID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), scopeKey)
	if cache.Client != nil && cache.Client.IsEnabled() {
		var cachedResult product_analytics.MonthlySalesResponse
		if found, _ := cache.Client.Get(cacheKey, &cachedResult); found {
			return &cachedResult, nil
		}
	}

	result, err := s.productAnalyticsRepo.GetProductMonthlySales(productID, startDate, endDate, scopedUserIDs)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cache.Client != nil && cache.Client.IsEnabled() && result != nil {
		_ = cache.Client.Set(cacheKey, result, TTLAnalyticsMonthlySales)
	}

	// Debug logging for data verification
	logger.LogInfo("ProductMonthlySales Query Result", logContext{
		"product_id":    productID,
		"start_date":    startDate,
		"end_date":      endDate,
		"total_sold":    result.TotalSold,
		"total_revenue": result.TotalRevenue,
		"total_profit":  result.TotalProfit,
		"total_sales":   result.TotalSales,
		"cache_status":  "miss",
	})

	return result, nil
}
