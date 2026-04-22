package category

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
)

// MockCategoryRepository
type MockCategoryRepository struct {
	FindByIDFunc   func(id string) (*category.Category, error)
	FindByCodeFunc func(code string) (*category.Category, error)
	ListFunc       func() ([]category.Category, error)
	CreateFunc     func(cat *category.Category) error
	UpdateFunc     func(cat *category.Category) error
	DeleteFunc     func(id string) error
}

func (m *MockCategoryRepository) FindByID(id string) (*category.Category, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockCategoryRepository) FindByCode(code string) (*category.Category, error) {
	if m.FindByCodeFunc != nil { return m.FindByCodeFunc(code) }
	return nil, nil
}
func (m *MockCategoryRepository) List() ([]category.Category, error) {
	if m.ListFunc != nil { return m.ListFunc() }
	return nil, nil
}
func (m *MockCategoryRepository) Create(cat *category.Category) error {
	if m.CreateFunc != nil { return m.CreateFunc(cat) }
	return nil
}
func (m *MockCategoryRepository) Update(cat *category.Category) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(cat) }
	return nil
}
func (m *MockCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}
