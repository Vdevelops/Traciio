package contact

import (


	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new contact repository
func NewRepository(db *gorm.DB) interfaces.ContactRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*contact.Contact, error) {
	var c contact.Contact
	err := r.db.Preload("Role").Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(req *contact.ListContactsRequest) ([]contact.Contact, int64, error) {
	var contacts []contact.Contact
	var total int64

	query := r.db.Model(&contact.Contact{})

	// Apply filters
	if req.Search != "" {
		// Optimized: Use Full Text Search instead of LIKE %...%
		// Uses GIN index on name, email, position
		query = query.Where(
			"to_tsvector('english', name || ' ' || email || ' ' || COALESCE(position, '')) @@ plainto_tsquery('english', ?) OR phone LIKE ?",
			req.Search, "%"+req.Search+"%", // Keep LIKE for phone partial match as FTS is bad for numbers
		)
	}

	if req.AccountID != "" {
		query = query.Where("account_id = ?", req.AccountID)
	}

	if req.RoleID != "" {
		query = query.Where("role_id = ?", req.RoleID)
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

	// Fetch data with preload
	err := query.Preload("Role").Order("created_at DESC").Offset(offset).Limit(perPage).Find(&contacts).Error
	if err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *repository) Create(c *contact.Contact) error {
	return r.db.Create(c).Error
}

func (r *repository) Update(c *contact.Contact) error {
	return r.db.Save(c).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&contact.Contact{}).Error
}

func (r *repository) FindByAccountID(accountID string) ([]contact.Contact, error) {
	var contacts []contact.Contact
	err := r.db.Preload("Role").Where("account_id = ?", accountID).Order("created_at DESC").Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

