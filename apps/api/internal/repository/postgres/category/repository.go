package category

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new category repository
func NewRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*category.Category, error) {
	var cat category.Category
	err := r.db.
		Select("*, (SELECT COUNT(*) FROM accounts WHERE accounts.category_id = categories.id) as account_count").
		Where("id = ?", id).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *repository) FindByCode(code string) (*category.Category, error) {
	var cat category.Category
	err := r.db.
		Select("*, (SELECT COUNT(*) FROM accounts WHERE accounts.category_id = categories.id) as account_count").
		Where("code = ?", code).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *repository) List() ([]category.Category, error) {
	var categories []category.Category
	err := r.db.
		Select("*, (SELECT COUNT(*) FROM accounts WHERE accounts.category_id = categories.id) as account_count").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) Create(cat *category.Category) error {
	return r.db.Create(cat).Error
}

func (r *repository) Update(cat *category.Category) error {
	return r.db.Save(cat).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&category.Category{}).Error
}

