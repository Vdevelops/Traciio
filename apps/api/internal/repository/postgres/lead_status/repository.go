package lead_status

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"

	"gorm.io/gorm"
)

// Repository defines lead status repository interface
type Repository interface {
	Create(status *lead_status.LeadStatus) error
	Update(status *lead_status.LeadStatus) error
	Delete(id string) error
	FindByID(id string) (*lead_status.LeadStatus, error)
	FindByCode(code string) (*lead_status.LeadStatus, error)
	List(req *lead_status.ListLeadStatusesRequest) ([]*lead_status.LeadStatus, int64, error)
	ListAll() ([]*lead_status.LeadStatus, error)
	SetDefault(id string) error
	FindDefault() (*lead_status.LeadStatus, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new lead status repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(status *lead_status.LeadStatus) error {
	return r.db.Create(status).Error
}

func (r *repository) Update(status *lead_status.LeadStatus) error {
	return r.db.Save(status).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&lead_status.LeadStatus{}, "id = ?", id).Error
}

func (r *repository) FindByID(id string) (*lead_status.LeadStatus, error) {
	var status lead_status.LeadStatus
	err := r.db.Where("id = ?", id).First(&status).Error
	return &status, err
}

func (r *repository) FindByCode(code string) (*lead_status.LeadStatus, error) {
	var status lead_status.LeadStatus
	err := r.db.Where("code = ?", code).First(&status).Error
	return &status, err
}

func (r *repository) List(req *lead_status.ListLeadStatusesRequest) ([]*lead_status.LeadStatus, int64, error) {
	var statuses []*lead_status.LeadStatus
	var total int64

	query := r.db.Model(&lead_status.LeadStatus{})

	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Add lead count
	query = query.Select("*, (SELECT COUNT(*) FROM leads WHERE leads.lead_status_id = lead_statuses.id) as lead_count")

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
	// Use table name prefix to ensure proper quoting
	if sortBy == "order" {
		query = query.Order(`lead_statuses."order" ` + sortOrder)
	} else {
		query = query.Order("lead_statuses." + sortBy + " " + sortOrder)
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

	if err := query.Limit(perPage).Offset(offset).Find(&statuses).Error; err != nil {
		return nil, 0, err
	}

	return statuses, total, nil
}

func (r *repository) ListAll() ([]*lead_status.LeadStatus, error) {
	var statuses []*lead_status.LeadStatus
	// Order by order field (use double quotes for PostgreSQL reserved keyword)
	// CRITICAL: Add limit for master data (typically small, but safety measure)
	// Master data like lead statuses should be < 100, but limit to 500 for safety
	err := r.db.Where("is_active = ?", true).Order("\"order\" ASC").Limit(500).Find(&statuses).Error
	return statuses, err
}

func (r *repository) SetDefault(id string) error {
	// Start transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove default from all statuses
		if err := tx.Model(&lead_status.LeadStatus{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		// Set new default
		if err := tx.Model(&lead_status.LeadStatus{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) FindDefault() (*lead_status.LeadStatus, error) {
	var status lead_status.LeadStatus
	err := r.db.Where("is_default = ?", true).First(&status).Error
	return &status, err
}
