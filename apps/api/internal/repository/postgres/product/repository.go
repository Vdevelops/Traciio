package product

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new product repository.
func NewRepository(db *gorm.DB) interfaces.ProductRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*product.Product, error) {
	var p product.Product
	err := r.db.
		Preload("Category").
		Where("id = ?", id).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(req *product.ListProductsRequest) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	query := r.db.Model(&product.Product{}).Preload("Category")

	// Apply filters.
	if req.Search != "" {
		// Optimized: Use Full Text Search instead of LIKE %...%
		// Uses GIN index on name, sku, description
		// Note: Barcode search is usually exact match, but we include it in FTS vector if needed,
		// or keep standard index for it. For now, we focus on Name/SKU/Desc FTS.
		query = query.Where(
			"to_tsvector('english', name || ' ' || sku || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', ?) OR LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR barcode = ?",
			req.Search, "%"+strings.ToLower(req.Search)+"%", "%"+strings.ToLower(req.Search)+"%", req.Search,
		)
	}

	// NOTE:
	// Some environments may still have an older `products` schema without a `status` column.
	// In that case, applying the filter would 500 with "column status does not exist".
	// We guard by checking column existence through information_schema.
	if req.Status != "" {
		var hasStatusColumn bool
		_ = r.db.Raw(
			`SELECT EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = 'public'
				  AND table_name = 'products'
				  AND column_name = 'status'
			)`,
		).Scan(&hasStatusColumn).Error

		if hasStatusColumn {
			query = query.Where("status = ?", req.Status)
		}
	}

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// Count total.
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination.
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	if err := query.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *repository) Create(p *product.Product) error {
	return r.db.Create(p).Error
}

func (r *repository) Update(p *product.Product) error {
	return r.db.Save(p).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&product.Product{}, "id = ?", id).Error
}
