package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
)

// ProductRepository defines the interface for product repository.
type ProductRepository interface {
	// FindByID finds a product by ID.
	FindByID(id string) (*product.Product, error)

	// List returns a list of products with pagination.
	List(req *product.ListProductsRequest) ([]product.Product, int64, error)

	// Create creates a new product.
	Create(product *product.Product) error

	// Update updates a product.
	Update(product *product.Product) error

	// Delete soft deletes a product.
	Delete(id string) error
}

// ProductCategoryRepository defines the interface for product category repository.
type ProductCategoryRepository interface {
	// FindByID finds a product category by ID.
	FindByID(id string) (*product.ProductCategory, error)

	// List returns a list of product categories.
	List(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error)

	// Create creates a new product category.
	Create(category *product.ProductCategory) error

	// Update updates a product category.
	Update(category *product.ProductCategory) error

	// Delete soft deletes a product category.
	Delete(id string) error
}


