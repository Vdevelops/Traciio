package product

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductCategory represents category for products (e.g. Drug, Medical Device, Supplement)
type ProductCategory struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(150);uniqueIndex" json:"name"`
	Slug        string    `gorm:"type:varchar(150);uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Read-only fields
	ProductCount int64 `gorm:"->" json:"product_count"`
}

// TableName specifies the table name for ProductCategory.
func (ProductCategory) TableName() string {
	return "product_categories"
}

// BeforeCreate hook to generate UUID.
func (c *ProductCategory) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Slug == "" && c.Name != "" {
		c.Slug = generateSlug(c.Name)
	}
	return nil
}

// BeforeUpdate ensures slug stays in sync when name changes (if slug not provided explicitly).
func (c *ProductCategory) BeforeUpdate(tx *gorm.DB) error {
	if c.Slug == "" && c.Name != "" {
		c.Slug = generateSlug(c.Name)
	}
	return nil
}

// Product represents a product in the CRM system.
type Product struct {
	ID          string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string              `gorm:"type:varchar(200);index:idx_products_fts,type:gin,expression:to_tsvector('english'\\, name || ' ' || sku || ' ' || COALESCE(description\\, ''))" json:"name"`
	SKU         string              `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	Barcode     string              `gorm:"type:varchar(100)" json:"barcode"`
	Price       int64               `gorm:"type:bigint;not null;default:0" json:"price"` // Stored in smallest currency unit (sen)
	Cost        int64               `gorm:"type:bigint;not null;default:0" json:"cost"`
	CategoryID  string              `gorm:"type:uuid;not null;index" json:"category_id"`
	Category    *ProductCategoryRef `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Description string              `gorm:"type:text" json:"description"`
	Status      string              `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	ImageURL    string              `gorm:"type:varchar(500)" json:"image_url"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"-"`
}

// TableName specifies the table name for Product.
func (Product) TableName() string {
	return "products"
}

// BeforeCreate hook to generate UUID.
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// ProductCategoryRef represents category reference in product.
type ProductCategoryRef struct {
	ID   string `gorm:"type:uuid;primary_key" json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TableName specifies the table name for ProductCategoryRef.
func (ProductCategoryRef) TableName() string {
	return "product_categories"
}

// generateSlug creates a simple slug from name.
func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}
