package product

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
)

// MockProductRepository
type MockProductRepository struct {
	FindByIDFunc func(id string) (*product.Product, error)
	ListFunc     func(req *product.ListProductsRequest) ([]product.Product, int64, error)
	CreateFunc   func(product *product.Product) error
	UpdateFunc   func(product *product.Product) error
	DeleteFunc   func(id string) error
}
func (m *MockProductRepository) FindByID(id string) (*product.Product, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockProductRepository) List(req *product.ListProductsRequest) ([]product.Product, int64, error) {
	if m.ListFunc != nil { return m.ListFunc(req) }
	return nil, 0, nil
}
func (m *MockProductRepository) Create(p *product.Product) error {
	if m.CreateFunc != nil { return m.CreateFunc(p) }
	return nil
}
func (m *MockProductRepository) Update(p *product.Product) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(p) }
	return nil
}
func (m *MockProductRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}

// MockProductCategoryRepository
type MockProductCategoryRepository struct {
	FindByIDFunc func(id string) (*product.ProductCategory, error)
	ListFunc     func(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error)
	CreateFunc   func(category *product.ProductCategory) error
	UpdateFunc   func(category *product.ProductCategory) error
	DeleteFunc   func(id string) error
}
func (m *MockProductCategoryRepository) FindByID(id string) (*product.ProductCategory, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockProductCategoryRepository) List(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error) {
	if m.ListFunc != nil { return m.ListFunc(req) }
	return nil, nil
}
func (m *MockProductCategoryRepository) Create(category *product.ProductCategory) error {
	if m.CreateFunc != nil { return m.CreateFunc(category) }
	return nil
}
func (m *MockProductCategoryRepository) Update(category *product.ProductCategory) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(category) }
	return nil
}
func (m *MockProductCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}
