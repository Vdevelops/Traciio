package product_category

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new product category repository.
func NewRepository(db *gorm.DB) interfaces.ProductCategoryRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*product.ProductCategory, error) {
	var c product.ProductCategory
	if err := r.db.
		Select("*, (SELECT COUNT(*) FROM products WHERE products.category_id = product_categories.id) as product_count").
		Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error) {
	var categories []product.ProductCategory

	query := r.db.Model(&product.ProductCategory{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.
		Select("*, (SELECT COUNT(*) FROM products WHERE products.category_id = product_categories.id) as product_count").
		Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *repository) Create(c *product.ProductCategory) error {
	return r.db.Create(c).Error
}

func (r *repository) Update(c *product.ProductCategory) error {
	return r.db.Save(c).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&product.ProductCategory{}, "id = ?", id).Error
}


