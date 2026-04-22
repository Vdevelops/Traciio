package deal_history

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new deal history entry
func (r *Repository) Create(dealHistory *deal_history.DealHistory) error {
	return r.db.Create(dealHistory).Error
}

// FindByDealID finds all history entries for a deal
func (r *Repository) FindByDealID(dealID string) ([]deal_history.DealHistory, error) {
	var histories []deal_history.DealHistory
	err := r.db.Where("deal_id = ?", dealID).
		Order("changed_at DESC").
		Find(&histories).Error
	return histories, err
}

// FindByID finds a history entry by ID
func (r *Repository) FindByID(id string) (*deal_history.DealHistory, error) {
	var history deal_history.DealHistory
	err := r.db.First(&history, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// List returns a limited list of history entries for a deal
func (r *Repository) List(dealID string, limit int) ([]deal_history.DealHistory, error) {
	var histories []deal_history.DealHistory
	query := r.db.Where("deal_id = ?", dealID).
		Order("changed_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// Delete soft deletes a history entry
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&deal_history.DealHistory{}, "id = ?", id).Error
}
