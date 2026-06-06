package product_analytics

import "time"

// TopProductResponse represents top products DTO.
type TopProductResponse struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSKU   string  `json:"product_sku"`
	QuantitySold int     `json:"quantity_sold"`
	TotalRevenue int64   `json:"total_revenue"`
	AvgPrice     float64 `json:"avg_price"`
	GrowthRate   float64 `json:"growth_rate"`
	Rank         int     `json:"rank"`
}

// ProductPerformanceResponse represents product performance DTO.
type ProductPerformanceResponse struct {
	ProductID     string            `json:"product_id"`
	ProductName   string            `json:"product_name"`
	ProductSKU    string            `json:"product_sku"`
	TotalQuantity int               `json:"total_quantity"`
	TotalRevenue  int64             `json:"total_revenue"`
	AvgPrice      float64           `json:"avg_price"`
	TotalSales    int               `json:"total_sales"`
	UniqueBuyers  int               `json:"unique_buyers"`
	GrowthRate    float64           `json:"growth_rate"`
	SalesByPeriod []PeriodSalesData `json:"sales_by_period"`
	TopBuyers     []BuyerData       `json:"top_buyers"`
}

// PeriodSalesData represents sales data by period.
type PeriodSalesData struct {
	Period   string `json:"period"`
	Quantity int    `json:"quantity"`
	Revenue  int64  `json:"revenue"`
}

// BuyerData represents buyer information.
type BuyerData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Quantity int    `json:"quantity"`
	Revenue  int64  `json:"revenue"`
}

// ProductComparisonResponse represents product comparison DTO.
type ProductComparisonResponse struct {
	Products []ProductPerformanceResponse `json:"products"`
}

// ProductTrendResponse represents product trend DTO.
type ProductTrendResponse struct {
	ProductID   string            `json:"product_id"`
	ProductName string            `json:"product_name"`
	ProductSKU  string            `json:"product_sku"`
	Trends      []PeriodSalesData `json:"trends"`
}

// ProductListItem represents a product in the products list with analytics.
type ProductListItem struct {
	ProductID    string     `json:"product_id"`
	ProductName  string     `json:"product_name"`
	ProductSKU   string     `json:"product_sku"`
	CategoryID   string     `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ImageURL     string     `json:"image_url"`
	UnitPrice    int64      `json:"unit_price"`
	TotalSold    int        `json:"total_sold"`
	TotalRevenue int64      `json:"total_revenue"`
	TotalProfit  int64      `json:"total_profit"`
	AvgUnitPrice float64    `json:"avg_unit_price"`
	SalesCount   int        `json:"sales_count"`
	LastSoldAt   *time.Time `json:"last_sold_at,omitempty"`
}

// MonthlySalesData represents sales data for a specific month.
type MonthlySalesData struct {
	Month        int    `json:"month"`
	MonthName    string `json:"month_name"`
	Year         int    `json:"year"`
	TotalSold    int    `json:"total_sold"`
	TotalRevenue int64  `json:"total_revenue"`
	TotalProfit  int64  `json:"total_profit"`
	SalesCount   int    `json:"sales_count"`
}

// MonthlySalesResponse represents the response for monthly sales chart.
type MonthlySalesResponse struct {
	Year         int                `json:"year"`
	MonthlySales []MonthlySalesData `json:"monthly_sales"`
	TotalSold    int                `json:"total_sold"`
	TotalRevenue int64              `json:"total_revenue"`
	TotalProfit  int64              `json:"total_profit"`
	TotalSales   int                `json:"total_sales"`
}
