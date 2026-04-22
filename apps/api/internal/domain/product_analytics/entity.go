package product_analytics

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductSales represents product sales tracking table
type ProductSales struct {
	ID          string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DealID      string       `gorm:"type:uuid;not null;index" json:"deal_id"`
	ProductID   string       `gorm:"type:uuid;not null;index" json:"product_id"`
	Product     *ProductRef  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity    int          `gorm:"type:integer;not null;default:1" json:"quantity"`
	UnitPrice   int64        `gorm:"type:bigint;not null;default:0" json:"unit_price"` // Stored in smallest currency unit (sen)
	TotalPrice  int64        `gorm:"type:bigint;not null;default:0" json:"total_price"`
	SoldAt      time.Time    `gorm:"type:timestamp;not null;index" json:"sold_at"`
	SalesRepID  string       `gorm:"type:uuid;not null;index" json:"sales_rep_id"`
	SalesRep    *UserRef     `gorm:"foreignKey:SalesRepID" json:"sales_rep,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ProductSales
func (ProductSales) TableName() string {
	return "product_sales"
}

// BeforeCreate hook to generate UUID
func (ps *ProductSales) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = uuid.New().String()
	}
	return nil
}

// ProductRef represents product reference in analytics
type ProductRef struct {
	ID          string `gorm:"type:uuid;primary_key" json:"id"`
	Name        string `gorm:"type:varchar(200)" json:"name"`
	SKU         string `gorm:"type:varchar(100)" json:"sku"`
	CategoryID  string `gorm:"type:uuid" json:"category_id"`
}

// TableName specifies the table name for ProductRef
func (ProductRef) TableName() string {
	return "products"
}

// UserRef represents user reference in analytics
type UserRef struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// TableName specifies the table name for UserRef
func (UserRef) TableName() string {
	return "users"
}

// TopProductResponse represents top products DTO
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

// ProductPerformanceResponse represents product performance DTO
type ProductPerformanceResponse struct {
	ProductID       string               `json:"product_id"`
	ProductName     string               `json:"product_name"`
	ProductSKU      string               `json:"product_sku"`
	TotalQuantity   int                  `json:"total_quantity"`
	TotalRevenue    int64                `json:"total_revenue"`
	AvgPrice        float64              `json:"avg_price"`
	TotalSales      int                  `json:"total_sales"`
	UniqueBuyers    int                  `json:"unique_buyers"`
	GrowthRate      float64              `json:"growth_rate"`
	SalesByPeriod   []PeriodSalesData    `json:"sales_by_period"`
	TopBuyers       []BuyerData          `json:"top_buyers"`
}

// PeriodSalesData represents sales data by period
type PeriodSalesData struct {
	Period   string  `json:"period"`
	Quantity int     `json:"quantity"`
	Revenue  int64   `json:"revenue"`
}

// BuyerData represents buyer information
type BuyerData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Quantity int    `json:"quantity"`
	Revenue  int64  `json:"revenue"`
}

// ProductComparisonResponse represents product comparison DTO
type ProductComparisonResponse struct {
	Products []ProductPerformanceResponse `json:"products"`
}

// ProductTrendResponse represents product trend DTO
type ProductTrendResponse struct {
	ProductID   string            `json:"product_id"`
	ProductName string            `json:"product_name"`
	ProductSKU  string            `json:"product_sku"`
	Trends      []PeriodSalesData `json:"trends"`
}

// ProductListItem represents a product in the products list with analytics
type ProductListItem struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	ProductSKU    string  `json:"product_sku"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	ImageURL      string  `json:"image_url"`
	UnitPrice     int64   `json:"unit_price"`
	TotalSold     int     `json:"total_sold"`
	TotalRevenue  int64   `json:"total_revenue"`
	TotalProfit   int64   `json:"total_profit"`
	AvgUnitPrice  float64 `json:"avg_unit_price"`
	SalesCount    int     `json:"sales_count"`
	LastSoldAt    *time.Time `json:"last_sold_at,omitempty"`
}

// MonthlySalesData represents sales data for a specific month
type MonthlySalesData struct {
	Month        int     `json:"month"`         // 1-12
	MonthName    string  `json:"month_name"`    // January, February, etc.
	Year         int     `json:"year"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue int64   `json:"total_revenue"`
	TotalProfit  int64   `json:"total_profit"`
	SalesCount   int     `json:"sales_count"`
}

// MonthlySalesResponse represents the response for monthly sales chart
type MonthlySalesResponse struct {
	Year          int                 `json:"year"`
	MonthlySales  []MonthlySalesData  `json:"monthly_sales"`
	TotalSold     int                 `json:"total_sold"`
	TotalRevenue  int64               `json:"total_revenue"`
	TotalProfit   int64               `json:"total_profit"`
	TotalSales    int                 `json:"total_sales"`
}
