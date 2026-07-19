package product_analytics

import (
	"fmt"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new product analytics repository
func NewRepository(db *gorm.DB) interfaces.ProductAnalyticsRepository {
	return &repository{db: db}
}

func restrictDealAssignedToSalesRole(query *gorm.DB, dealAlias string) *gorm.DB {
	sanitizedAlias := strings.NewReplacer(".", "_").Replace(dealAlias)
	userAlias := sanitizedAlias + "_sales_scope_user"
	roleAlias := sanitizedAlias + "_sales_scope_role"

	return query.
		Joins(fmt.Sprintf(
			"INNER JOIN users %s ON %s.id = %s.assigned_to AND %s.deleted_at IS NULL",
			userAlias, userAlias, dealAlias, userAlias,
		)).
		Joins(fmt.Sprintf(
			"INNER JOIN roles %s ON %s.id = %s.role_id AND %s.deleted_at IS NULL AND %s.code = ?",
			roleAlias, roleAlias, userAlias, roleAlias, roleAlias,
		), "sales")
}

func restrictUserToSalesRole(query *gorm.DB, userColumn string) *gorm.DB {
	sanitizedColumn := strings.NewReplacer(".", "_").Replace(userColumn)
	userAlias := sanitizedColumn + "_sales_scope_user"
	roleAlias := sanitizedColumn + "_sales_scope_role"

	return query.
		Joins(fmt.Sprintf(
			"INNER JOIN users %s ON %s.id = %s AND %s.deleted_at IS NULL",
			userAlias, userAlias, userColumn, userAlias,
		)).
		Joins(fmt.Sprintf(
			"INNER JOIN roles %s ON %s.id = %s.role_id AND %s.deleted_at IS NULL AND %s.code = ?",
			roleAlias, roleAlias, userAlias, roleAlias, roleAlias,
		), "sales")
}

func soldProductRowsSQL() string {
	return `(
		SELECT
			d.id AS deal_id,
			d.assigned_to,
			COALESCE(d.actual_close_date, d.created_at) AS sold_at,
			d.account_id,
			dpi.product_id::text AS product_id,
			dpi.quantity,
			dpi.unit_price,
			dpi.subtotal
		FROM deal_product_items dpi
		INNER JOIN deals d ON dpi.deal_id = d.id AND d.deleted_at IS NULL AND d.status = 'won'
		WHERE dpi.deleted_at IS NULL

		UNION ALL

		SELECT
			d.id AS deal_id,
			d.assigned_to,
			COALESCE(d.actual_close_date, d.created_at) AS sold_at,
			d.account_id,
			item.product_id,
			item.quantity,
			item.unit_price,
			(item.unit_price * item.quantity) AS subtotal
		FROM deals d
		INNER JOIN visit_reports vr ON vr.deleted_at IS NULL AND (
			vr.deal_id = d.id OR (d.lead_id IS NOT NULL AND vr.lead_id = d.lead_id)
		)
		CROSS JOIN LATERAL (
			SELECT
				interest->>'product_id' AS product_id,
				CASE
					WHEN COALESCE(interest->>'quantity', '') ~ '^[0-9]+$' THEN (interest->>'quantity')::int
					ELSE 1
				END AS quantity,
				CASE
					WHEN COALESCE(interest->>'price', '') ~ '^[0-9]+$' THEN (interest->>'price')::bigint
					ELSE COALESCE(p.price, 0)
				END AS unit_price
			FROM jsonb_array_elements(COALESCE(vr.metadata->'product_interests', '[]'::jsonb)) AS interest
			LEFT JOIN products p ON p.id::text = interest->>'product_id' AND p.deleted_at IS NULL
			WHERE COALESCE(interest->>'product_id', '') <> ''
		) item
		WHERE d.deleted_at IS NULL
			AND d.status = 'won'
			AND NOT EXISTS (
				SELECT 1
				FROM deal_product_items existing
				WHERE existing.deal_id = d.id AND existing.deleted_at IS NULL
			)

		UNION ALL

		SELECT
			d.id AS deal_id,
			d.assigned_to,
			COALESCE(d.actual_close_date, d.created_at) AS sold_at,
			d.account_id,
			item.product_id,
			item.quantity,
			item.unit_price,
			(item.unit_price * item.quantity) AS subtotal
		FROM deals d
		INNER JOIN activities a ON a.deleted_at IS NULL AND (
			a.deal_id = d.id OR (d.lead_id IS NOT NULL AND a.lead_id = d.lead_id)
		)
		CROSS JOIN LATERAL (
			SELECT
				interest->>'product_id' AS product_id,
				CASE
					WHEN COALESCE(interest->>'quantity', '') ~ '^[0-9]+$' THEN (interest->>'quantity')::int
					ELSE 1
				END AS quantity,
				CASE
					WHEN COALESCE(interest->>'price', '') ~ '^[0-9]+$' THEN (interest->>'price')::bigint
					ELSE COALESCE(p.price, 0)
				END AS unit_price
			FROM jsonb_array_elements(COALESCE(a.metadata->'product_interests', '[]'::jsonb)) AS interest
			LEFT JOIN products p ON p.id::text = interest->>'product_id' AND p.deleted_at IS NULL
			WHERE COALESCE(interest->>'product_id', '') <> ''
		) item
		WHERE d.deleted_at IS NULL
			AND d.status = 'won'
			AND NOT EXISTS (
				SELECT 1
				FROM deal_product_items existing
				WHERE existing.deal_id = d.id AND existing.deleted_at IS NULL
			)
	)`
}

func (r *repository) CreateProductSale(productSale *product_analytics.ProductSales) error {
	return r.db.Create(productSale).Error
}

func (r *repository) GetProductSaleByID(id string) (*product_analytics.ProductSales, error) {
	var productSale product_analytics.ProductSales
	err := r.db.Preload("Product").Preload("SalesRep").Where("id = ?", id).First(&productSale).Error
	if err != nil {
		return nil, err
	}
	return &productSale, nil
}

func (r *repository) GetProductSales(filters map[string]interface{}, page, perPage int) ([]*product_analytics.ProductSales, int64, error) {
	var productSales []*product_analytics.ProductSales
	var total int64

	query := r.db.Model(&product_analytics.ProductSales{}).Preload("Product").Preload("SalesRep")
	query = restrictUserToSalesRole(query, "product_sales.sales_rep_id")

	// Apply filters
	if productID, ok := filters["product_id"].(string); ok && productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if salesRepID, ok := filters["sales_rep_id"].(string); ok && salesRepID != "" {
		query = query.Where("sales_rep_id = ?", salesRepID)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("sold_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("sold_at <= ?", endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	err := query.Order("sold_at DESC").Offset(offset).Limit(perPage).Find(&productSales).Error
	if err != nil {
		return nil, 0, err
	}

	return productSales, total, nil
}

func (r *repository) DeleteProductSale(id string) error {
	return r.db.Where("id = ?", id).Delete(&product_analytics.ProductSales{}).Error
}

func (r *repository) GetProductPerformance(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.ProductPerformanceResponse, error) {
	var performance product_analytics.ProductPerformanceResponse

	// Get product info
	var product struct {
		ID   string
		Name string
		SKU  string
	}
	err := r.db.Table("products").
		Select("id, name, sku").
		Where("id = ?", productID).
		Scan(&product).Error
	if err != nil {
		return nil, err
	}

	performance.ProductID = product.ID
	performance.ProductName = product.Name
	performance.ProductSKU = product.SKU

	// Get aggregated metrics — source of truth: deal_product_items joined to won deals.
	// This is consistent with GetProductsList / GetMonthlySales and filters by d.assigned_to
	// so the scope (admin / team / own) is applied uniformly.
	var metrics struct {
		TotalQuantity int
		TotalRevenue  int64
		AvgPrice      float64
		TotalSales    int64
		UniqueBuyers  int64
	}

	metricsQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(`
			COALESCE(SUM(si.quantity), 0) as total_quantity,
			COALESCE(SUM(si.subtotal), 0) as total_revenue,
			COALESCE(AVG(si.unit_price), 0) as avg_price,
			COUNT(*) as total_sales,
			COUNT(DISTINCT si.account_id) as unique_buyers
		`).
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate)
	metricsQuery = restrictDealAssignedToSalesRole(metricsQuery, "si")
	if len(scopedUserIDs) > 0 {
		metricsQuery = metricsQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err = metricsQuery.Scan(&metrics).Error
	if err != nil {
		return nil, err
	}

	performance.TotalQuantity = metrics.TotalQuantity
	performance.TotalRevenue = metrics.TotalRevenue
	performance.AvgPrice = metrics.AvgPrice
	performance.TotalSales = int(metrics.TotalSales)
	performance.UniqueBuyers = int(metrics.UniqueBuyers)

	// Calculate growth rate against the equivalent prior period.
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := startDate.Add(-1 * time.Second)

	var prevRevenue struct {
		Total int64
	}
	prevQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select("COALESCE(SUM(si.subtotal), 0) as total").
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", prevStartDate).
		Where("si.sold_at <= ?", prevEndDate)
	prevQuery = restrictDealAssignedToSalesRole(prevQuery, "si")
	if len(scopedUserIDs) > 0 {
		prevQuery = prevQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err = prevQuery.Scan(&prevRevenue).Error
	if err == nil && prevRevenue.Total > 0 {
		performance.GrowthRate = ((float64(performance.TotalRevenue) - float64(prevRevenue.Total)) / float64(prevRevenue.Total)) * 100.0
	} else if performance.TotalRevenue > 0 {
		performance.GrowthRate = 100.0
	}

	// Get sales breakdown by month within the requested period.
	var salesByPeriod []product_analytics.PeriodSalesData
	periodQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(`
			TO_CHAR(si.sold_at, 'YYYY-MM') as period,
			SUM(si.quantity) as quantity,
			SUM(si.subtotal) as revenue
		`).
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate)
	periodQuery = restrictDealAssignedToSalesRole(periodQuery, "si")
	if len(scopedUserIDs) > 0 {
		periodQuery = periodQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err = periodQuery.
		Group("TO_CHAR(si.sold_at, 'YYYY-MM')").
		Order("period ASC").
		Scan(&salesByPeriod).Error
	if err == nil {
		performance.SalesByPeriod = salesByPeriod
	}

	// Get top customer accounts for this product within the scoped dataset.
	var topBuyers []product_analytics.BuyerData
	topBuyersQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Joins("LEFT JOIN accounts a ON si.account_id = a.id AND a.deleted_at IS NULL").
		Select(`
			si.account_id as buyer_id,
			COALESCE(a.name, 'Unknown Customer') as buyer_name,
			SUM(si.quantity) as quantity,
			SUM(si.subtotal) as revenue
		`).
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate)
	topBuyersQuery = restrictDealAssignedToSalesRole(topBuyersQuery, "si")
	if len(scopedUserIDs) > 0 {
		topBuyersQuery = topBuyersQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err = topBuyersQuery.
		Where("si.account_id IS NOT NULL").
		Group("si.account_id, a.name").
		Order("revenue DESC").
		Limit(5).
		Scan(&topBuyers).Error
	if err == nil {
		performance.TopBuyers = topBuyers
	}

	return &performance, nil
}

func (r *repository) GetProductComparison(productIDs []string, startDate, endDate time.Time, scopedUserIDs []string) ([]*product_analytics.ProductPerformanceResponse, error) {
	var performances []*product_analytics.ProductPerformanceResponse

	for _, productID := range productIDs {
		performance, err := r.GetProductPerformance(productID, startDate, endDate, scopedUserIDs)
		if err != nil {
			continue // Skip products with errors
		}
		performances = append(performances, performance)
	}

	return performances, nil
}

func (r *repository) GetProductTrends(productID string, startDate, endDate time.Time, groupBy string, scopedUserIDs []string) (*product_analytics.ProductTrendResponse, error) {
	var trend product_analytics.ProductTrendResponse

	// Get product info
	var product struct {
		ID   string
		Name string
		SKU  string
	}
	err := r.db.Table("products").
		Select("id, name, sku").
		Where("id = ?", productID).
		Scan(&product).Error
	if err != nil {
		return nil, err
	}

	trend.ProductID = product.ID
	trend.ProductName = product.Name
	trend.ProductSKU = product.SKU

	// Determine grouping format
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "YYYY-MM-DD"
	case "week":
		dateFormat = "IYYY-IW" // ISO year and week
	case "month":
		dateFormat = "YYYY-MM"
	case "year":
		dateFormat = "YYYY"
	default:
		dateFormat = "YYYY-MM" // Default to monthly
	}

	// Get trends — use deal_product_items + deals (source of truth) filtered by assigned_to.
	var trends []product_analytics.PeriodSalesData
	trendsQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(fmt.Sprintf(`
			TO_CHAR(si.sold_at, '%s') as period,
			SUM(si.quantity) as quantity,
			SUM(si.subtotal) as revenue
		`, dateFormat)).
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate)
	trendsQuery = restrictDealAssignedToSalesRole(trendsQuery, "si")
	if len(scopedUserIDs) > 0 {
		trendsQuery = trendsQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err = trendsQuery.
		Group(fmt.Sprintf("TO_CHAR(si.sold_at, '%s')", dateFormat)).
		Order("period ASC").
		Scan(&trends).Error
	if err != nil {
		return nil, err
	}

	trend.Trends = trends

	return &trend, nil
}

func (r *repository) GetProductsList(startDate, endDate time.Time, search, sortBy, orderBy string, page, perPage int, scopedUserIDs []string) ([]*product_analytics.ProductListItem, int64, error) {
	var productsList []*product_analytics.ProductListItem
	var total int64

	// Base Query for counting and fetching
	// We need to query products that have been sold (or all products? User requirement implies sold products or just products list)
	// The original query joined with deal_product_items, effectively limiting to sold products.
	// However, usually a "Product List" should show all products, but for "Analytics" it often implies products with sales.
	// Guided by previous implementation: "This ensures we get products that are actually sold through the pipeline".
	// We will stick to that logic but we need to count them first.

	// Common Base Query
	baseQuery := r.db.Table(soldProductRowsSQL() + " AS si").
		Joins("INNER JOIN products p ON p.id::text = si.product_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.deleted_at IS NULL")
	baseQuery = restrictDealAssignedToSalesRole(baseQuery, "si")

	// Apply Filters to Base Query
	if !startDate.IsZero() {
		baseQuery = baseQuery.Where("si.sold_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		baseQuery = baseQuery.Where("si.sold_at <= ?", endDate)
	}
	if search != "" {
		searchLike := "%" + search + "%"
		baseQuery = baseQuery.Where("p.name ILIKE ?", searchLike)
	}
	if len(scopedUserIDs) > 0 {
		baseQuery = baseQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}

	// Count Total Unique Products
	// We group by product ID because the joins produce multiple rows per product (one per sale item)
	// countQuery needs to count distinct product IDs
	if err := baseQuery.Distinct("p.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Select Query
	query := baseQuery.Select(`
			p.id as product_id,
			p.name as product_name,
			p.sku as product_sku,
			p.category_id as category_id,
			pc.name as category_name,
			p.image_url as image_url,
			p.price as unit_price,
			COALESCE(SUM(si.quantity), 0) as total_sold,
			COALESCE(SUM(si.subtotal), 0) as total_revenue,
			COALESCE(CAST(SUM(si.subtotal) - SUM(p.cost * si.quantity) AS BIGINT), 0) as total_profit,
			COALESCE(AVG(si.unit_price), p.price) as avg_unit_price,
			COUNT(*) as sales_count,
			MAX(si.sold_at) as last_sold_at
		`).
		Group("p.id, p.name, p.sku, p.category_id, pc.name, p.image_url, p.price, p.cost")

	// Apply sorting
	var orderClause string
	switch sortBy {
	case "revenue":
		orderClause = "total_revenue"
	case "profit":
		orderClause = "total_profit"
	case "name":
		orderClause = "product_name"
	default: // total_sold
		orderClause = "total_sold"
	}

	if orderBy == "asc" {
		query = query.Order(orderClause + " ASC")
	} else {
		query = query.Order(orderClause + " DESC")
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	query = query.Offset(offset).Limit(perPage)

	err := query.Scan(&productsList).Error
	if err != nil {
		return nil, 0, err
	}

	return productsList, total, nil
}

func (r *repository) GetUserProductSales(userID string, startDate, endDate time.Time, sortBy, orderBy string, page, perPage int) ([]*product_analytics.ProductListItem, int64, error) {
	var productsList []*product_analytics.ProductListItem
	var total int64

	// Get product sales from deal_product_items for won deals (source of truth from pipeline)
	// This ensures we get products that are actually sold through the pipeline
	countQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Joins("INNER JOIN products p ON p.id::text = si.product_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("si.assigned_to = ?", userID). // Deal assigned to this user
		Where("p.deleted_at IS NULL")
	countQuery = restrictDealAssignedToSalesRole(countQuery, "si")

	// Apply date filter based on deal actual_close_date (when deal was won)
	if !startDate.IsZero() {
		countQuery = countQuery.Where("si.sold_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		countQuery = countQuery.Where("si.sold_at <= ?", endDate)
	}

	if err := countQuery.Distinct("p.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(`
			p.id as product_id,
			p.name as product_name,
			p.sku as product_sku,
			p.category_id as category_id,
			pc.name as category_name,
			p.image_url as image_url,
			p.price as unit_price,
			COALESCE(SUM(si.quantity), 0) as total_sold,
			COALESCE(SUM(si.subtotal), 0) as total_revenue,
			COALESCE(CAST(SUM(si.subtotal) - SUM(p.cost * si.quantity) AS BIGINT), 0) as total_profit,
			COALESCE(AVG(si.unit_price), p.price) as avg_unit_price,
			COUNT(*) as sales_count,
			MAX(si.sold_at) as last_sold_at
		`).
		Joins("INNER JOIN products p ON p.id::text = si.product_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.deleted_at IS NULL").
		Where("si.assigned_to = ?", userID) // Deal assigned to this user (connected to pipeline)
	query = restrictDealAssignedToSalesRole(query, "si")

	// Apply date filter based on deal actual_close_date (when deal was won)
	if !startDate.IsZero() {
		query = query.Where("si.sold_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("si.sold_at <= ?", endDate)
	}

	query = query.Group("p.id, p.name, p.sku, p.category_id, pc.name, p.image_url, p.price, p.cost")

	// Apply sorting
	var orderClause string
	switch sortBy {
	case "revenue":
		orderClause = "total_revenue"
	case "profit":
		orderClause = "total_profit"
	case "name":
		orderClause = "product_name"
	default: // total_sold
		orderClause = "total_sold"
	}

	if orderBy == "asc" {
		query = query.Order(orderClause + " ASC")
	} else {
		query = query.Order(orderClause + " DESC")
	}

	// Apply pagination
	if perPage > 0 {
		offset := (page - 1) * perPage
		query = query.Limit(perPage).Offset(offset)
	}

	err := query.Scan(&productsList).Error
	if err != nil {
		return nil, 0, err
	}

	return productsList, total, nil
}

func (r *repository) GetMonthlySales(startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error) {
	var monthlySales []product_analytics.MonthlySalesData

	// Month names mapping
	monthNames := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	// Query for monthly aggregated data
	type MonthlyResult struct {
		Year         int
		Month        int
		TotalSold    int
		TotalRevenue int64
		TotalProfit  int64
		SalesCount   int
	}

	var results []MonthlyResult
	monthlyQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(`
			EXTRACT(YEAR FROM si.sold_at)::int as year,
			EXTRACT(MONTH FROM si.sold_at)::int as month,
			COALESCE(SUM(si.quantity), 0) as total_sold,
			COALESCE(SUM(si.subtotal), 0) as total_revenue,
			COALESCE(CAST(SUM(si.subtotal) - SUM(p.cost * si.quantity) AS BIGINT), 0) as total_profit,
			COUNT(*) as sales_count
		`).
		Joins("INNER JOIN products p ON p.id::text = si.product_id").
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate).
		Where("p.deleted_at IS NULL")
	monthlyQuery = restrictDealAssignedToSalesRole(monthlyQuery, "si")
	if len(scopedUserIDs) > 0 {
		monthlyQuery = monthlyQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err := monthlyQuery.
		Group("EXTRACT(YEAR FROM si.sold_at), EXTRACT(MONTH FROM si.sold_at)").
		Order("year ASC, month ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup: [year][month]
	resultMap := make(map[int]map[int]MonthlyResult)
	for _, result := range results {
		if _, ok := resultMap[result.Year]; !ok {
			resultMap[result.Year] = make(map[int]MonthlyResult)
		}
		resultMap[result.Year][result.Month] = result
	}

	// Build response by iterating through each month in the range
	var totalSold int
	var totalRevenue int64
	var totalProfit int64
	var totalSales int

	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	// We use LastDay of endDate's month to ensure we cover the final month
	limit := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(limit) {
		y := current.Year()
		m := int(current.Month())

		var monthData product_analytics.MonthlySalesData
		monthData.Month = m
		monthData.MonthName = monthNames[m-1]
		monthData.Year = y

		if yrMap, ok := resultMap[y]; ok {
			if res, ok := yrMap[m]; ok {
				monthData.TotalSold = res.TotalSold
				monthData.TotalRevenue = res.TotalRevenue
				monthData.TotalProfit = res.TotalProfit
				monthData.SalesCount = res.SalesCount

				totalSold += res.TotalSold
				totalRevenue += res.TotalRevenue
				totalProfit += res.TotalProfit
				totalSales += res.SalesCount
			}
		}

		monthlySales = append(monthlySales, monthData)
		current = current.AddDate(0, 1, 0)
	}

	response := &product_analytics.MonthlySalesResponse{
		Year:         startDate.Year(), // Primary year or starting year
		MonthlySales: monthlySales,
		TotalSold:    totalSold,
		TotalRevenue: totalRevenue,
		TotalProfit:  totalProfit,
		TotalSales:   totalSales,
	}

	return response, nil
}

func (r *repository) GetProductMonthlySales(productID string, startDate, endDate time.Time, scopedUserIDs []string) (*product_analytics.MonthlySalesResponse, error) {
	var monthlySales []product_analytics.MonthlySalesData

	// Month names mapping
	monthNames := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	// Query for monthly aggregated data for specific product
	type MonthlyResult struct {
		Year         int
		Month        int
		TotalSold    int
		TotalRevenue int64
		TotalProfit  int64
		SalesCount   int
	}

	var results []MonthlyResult
	prodMonthlyQuery := r.db.Table(soldProductRowsSQL()+" AS si").
		Select(`
			EXTRACT(YEAR FROM si.sold_at)::int as year,
			EXTRACT(MONTH FROM si.sold_at)::int as month,
			COALESCE(SUM(si.quantity), 0) as total_sold,
			COALESCE(SUM(si.subtotal), 0) as total_revenue,
			COALESCE(CAST(SUM(si.subtotal) - SUM(p.cost * si.quantity) AS BIGINT), 0) as total_profit,
			COUNT(*) as sales_count
		`).
		Joins("INNER JOIN products p ON p.id::text = si.product_id").
		Where("si.product_id = ?", productID).
		Where("si.sold_at >= ?", startDate).
		Where("si.sold_at <= ?", endDate).
		Where("p.deleted_at IS NULL")
	prodMonthlyQuery = restrictDealAssignedToSalesRole(prodMonthlyQuery, "si")
	if len(scopedUserIDs) > 0 {
		prodMonthlyQuery = prodMonthlyQuery.Where("si.assigned_to IN ?", scopedUserIDs)
	}
	err := prodMonthlyQuery.
		Group("EXTRACT(YEAR FROM si.sold_at), EXTRACT(MONTH FROM si.sold_at)").
		Order("year ASC, month ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup: [year][month]
	resultMap := make(map[int]map[int]MonthlyResult)
	for _, result := range results {
		if _, ok := resultMap[result.Year]; !ok {
			resultMap[result.Year] = make(map[int]MonthlyResult)
		}
		resultMap[result.Year][result.Month] = result
	}

	// Build response by iterating through each month in the range
	var totalSold int
	var totalRevenue int64
	var totalProfit int64
	var totalSales int

	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	limit := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(limit) {
		y := current.Year()
		m := int(current.Month())

		var monthData product_analytics.MonthlySalesData
		monthData.Month = m
		monthData.MonthName = monthNames[m-1]
		monthData.Year = y

		if yrMap, ok := resultMap[y]; ok {
			if res, ok := yrMap[m]; ok {
				monthData.TotalSold = res.TotalSold
				monthData.TotalRevenue = res.TotalRevenue
				monthData.TotalProfit = res.TotalProfit
				monthData.SalesCount = res.SalesCount

				totalSold += res.TotalSold
				totalRevenue += res.TotalRevenue
				totalProfit += res.TotalProfit
				totalSales += res.SalesCount
			}
		}

		monthlySales = append(monthlySales, monthData)
		current = current.AddDate(0, 1, 0)
	}

	response := &product_analytics.MonthlySalesResponse{
		Year:         startDate.Year(),
		MonthlySales: monthlySales,
		TotalSold:    totalSold,
		TotalRevenue: totalRevenue,
		TotalProfit:  totalProfit,
		TotalSales:   totalSales,
	}

	return response, nil
}
