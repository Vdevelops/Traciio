package account

import (
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new account repository
func NewRepository(db *gorm.DB) interfaces.AccountRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*account.Account, error) {
	var a account.Account
	err := r.db.Preload("Category").Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) List(req *account.ListAccountsRequest) ([]account.Account, int64, error) {
	var accounts []account.Account
	var total int64

	query := r.db.Model(&account.Account{})

	// Apply filters
	if req.Search != "" {
		// Optimized: Use Full Text Search instead of LIKE %...%
		// Uses GIN index on name, address, city
		query = query.Where(
			"to_tsvector('english', name || ' ' || COALESCE(address, '') || ' ' || COALESCE(city, '')) @@ plainto_tsquery('english', ?) OR LOWER(email) = ? OR phone = ?",
			req.Search, strings.ToLower(req.Search), req.Search,
		)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.AssignedTo != "" {
		query = query.Where("assigned_to = ?", req.AssignedTo)
	} else if len(req.ScopedUserIDs) > 0 {
		query = query.Where("assigned_to IN ?", req.ScopedUserIDs)
	}

	if req.BrickID != "" {
		query = query.Where("brick_id = ?", req.BrickID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
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

	// Fetch data with preload and inline contact_count subquery to avoid N+1 queries.
	// The subquery is evaluated per row by the DB engine using the existing FK index on contacts.account_id.
	err := query.
		Preload("Category").
		Select("accounts.*, (SELECT COUNT(*) FROM contacts WHERE contacts.account_id = accounts.id AND contacts.deleted_at IS NULL) AS contact_count").
		Order("accounts.created_at DESC").
		Offset(offset).Limit(perPage).
		Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *repository) ListAll(status string) ([]account.Account, error) {
	var accounts []account.Account

	query := r.db.Model(&account.Account{})

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// CRITICAL: Add limit to prevent memory exhaustion with large datasets
	// For map display, 1000 accounts should be sufficient
	// If more are needed, use paginated List() method instead
	err := query.Preload("Category").Order("created_at DESC").Limit(1000).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *repository) ListByBBox(req *account.BBoxRequest) ([]account.Account, error) {
	var accounts []account.Account

	limit := req.Limit
	if limit < 1 || limit > 5000 {
		limit = 2000
	}

	query := r.db.Model(&account.Account{}).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL").
		Where("latitude BETWEEN ? AND ?", req.South, req.North).
		Where("longitude BETWEEN ? AND ?", req.West, req.East)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.Search != "" {
		query = query.Where(
			"to_tsvector('english', name || ' ' || COALESCE(address, '') || ' ' || COALESCE(city, '')) @@ plainto_tsquery('english', ?) OR LOWER(name) LIKE LOWER(?)",
			req.Search, "%"+strings.ToLower(req.Search)+"%",
		)
	}

	err := query.
		Preload("Category").
		Select("accounts.*, (SELECT COUNT(*) FROM contacts WHERE contacts.account_id = accounts.id AND contacts.deleted_at IS NULL) AS contact_count").
		Order("accounts.name ASC").
		Limit(limit).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *repository) Create(a *account.Account) error {
	return r.db.Create(a).Error
}

func (r *repository) Update(a *account.Account) error {
	return r.db.Save(a).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&account.Account{}).Error
}

// GetStatsByStatus returns account statistics grouped by status using database aggregation
func (r *repository) GetStatsByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}
	err := r.db.Table("accounts").
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}

	return stats, nil
}

// CountByDateRange returns count of accounts created in date range using database aggregation
func (r *repository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	var count int64
	query := r.db.Table("accounts")

	if startDate != nil {
		if start, ok := startDate.(time.Time); ok {
			query = query.Where("created_at >= ?", start)
		}
	}
	if endDate != nil {
		if end, ok := endDate.(time.Time); ok {
			query = query.Where("created_at <= ?", end)
		}
	}

	err := query.Count(&count).Error
	return count, err
}
