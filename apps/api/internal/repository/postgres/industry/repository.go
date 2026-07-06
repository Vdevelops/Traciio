package industry

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"

	"gorm.io/gorm"
)

// Repository defines industry repository interface
type Repository interface {
	Create(ind *industry.Industry) error
	Update(ind *industry.Industry) error
	Delete(id string) error
	FindByID(id string) (*industry.Industry, error)
	FindByCode(code string) (*industry.Industry, error)
	List(req *industry.ListIndustriesRequest) ([]*industry.Industry, int64, error)
	ListAll() ([]*industry.Industry, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new industry repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ind *industry.Industry) error {
	return r.db.Create(ind).Error
}

func (r *repository) Update(ind *industry.Industry) error {
	return r.db.Save(ind).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&industry.Industry{}, "id = ?", id).Error
}

func (r *repository) FindByID(id string) (*industry.Industry, error) {
	var ind industry.Industry
	err := r.db.Where("id = ?", id).First(&ind).Error
	return &ind, err
}

func (r *repository) FindByCode(code string) (*industry.Industry, error) {
	var ind industry.Industry
	err := r.db.Where("code = ?", code).First(&ind).Error
	return &ind, err
}

func (r *repository) List(req *industry.ListIndustriesRequest) ([]*industry.Industry, int64, error) {
	var industries []*industry.Industry
	var total int64

	query := r.db.Model(&industry.Industry{})

	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Add lead count
	query = query.Select("*, (SELECT COUNT(*) FROM leads WHERE leads.industry = industries.name) as lead_count")

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "order"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	sortOrder := "ASC"
	if req.SortOrder != "" {
		sortOrder = req.SortOrder
	}
	// Quote column name if it's a reserved keyword (use double quotes for PostgreSQL)
	if sortBy == "order" {
		query = query.Order(`industries."order" ` + sortOrder)
	} else {
		query = query.Order("industries." + sortBy + " " + sortOrder)
	}

	// Apply pagination
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	perPage := 10
	if req.PerPage > 0 {
		perPage = req.PerPage
	}
	offset := (page - 1) * perPage

	if err := query.Limit(perPage).Offset(offset).Find(&industries).Error; err != nil {
		return nil, 0, err
	}

	return industries, total, nil
}

func (r *repository) ListAll() ([]*industry.Industry, error) {
	var industries []*industry.Industry
	// Order by order field (use double quotes for PostgreSQL reserved keyword)
	// CRITICAL: Add limit for master data (typically small, but safety measure)
	// Master data like industries should be < 100, but limit to 500 for safety
	err := r.db.Where("is_active = ?", true).Order("\"order\" ASC").Limit(500).Find(&industries).Error
	return industries, err
}
