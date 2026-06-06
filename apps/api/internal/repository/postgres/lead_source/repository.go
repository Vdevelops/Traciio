package lead_source

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"

	"gorm.io/gorm"
)

// Repository defines lead source repository interface
type Repository interface {
	Create(ls *lead_source.LeadSource) error
	Update(ls *lead_source.LeadSource) error
	Delete(id string) error
	FindByID(id string) (*lead_source.LeadSource, error)
	FindByCode(code string) (*lead_source.LeadSource, error)
	List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error)
	ListAll() ([]*lead_source.LeadSource, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new lead source repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ls *lead_source.LeadSource) error {
	return r.db.Create(ls).Error
}

func (r *repository) Update(ls *lead_source.LeadSource) error {
	return r.db.Save(ls).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&lead_source.LeadSource{}, "id = ?", id).Error
}

func (r *repository) FindByID(id string) (*lead_source.LeadSource, error) {
	var ls lead_source.LeadSource
	err := r.db.Where("id = ?", id).First(&ls).Error
	return &ls, err
}

func (r *repository) FindByCode(code string) (*lead_source.LeadSource, error) {
	var ls lead_source.LeadSource
	err := r.db.Where("code = ?", code).First(&ls).Error
	return &ls, err
}

func (r *repository) List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error) {
	var leadSources []*lead_source.LeadSource
	var total int64

	query := r.db.Model(&lead_source.LeadSource{})

	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Add lead count - fix case sensitivity issue (stored as lowercase/mixed in leads table)
	// Match against either name or code (since seeder uses code for storing, but some might be name)
	query = query.Select("*, (SELECT COUNT(*) FROM leads WHERE LOWER(leads.lead_source) = LOWER(lead_sources.name) OR LOWER(leads.lead_source) = LOWER(lead_sources.code)) as lead_count")

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
		query = query.Order(`lead_sources."order" ` + sortOrder)
	} else {
		query = query.Order("lead_sources." + sortBy + " " + sortOrder)
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

	if err := query.Limit(perPage).Offset(offset).Find(&leadSources).Error; err != nil {
		return nil, 0, err
	}

	return leadSources, total, nil
}

func (r *repository) ListAll() ([]*lead_source.LeadSource, error) {
	var leadSources []*lead_source.LeadSource
	// Order by order field (use double quotes for PostgreSQL reserved keyword)
	// CRITICAL: Add limit for master data (typically small, but safety measure)
	// Master data like lead sources should be < 100, but limit to 500 for safety
	err := r.db.Where("is_active = ?", true).Order("\"order\" ASC").Limit(500).Find(&leadSources).Error
	return leadSources, err
}
