package product

import (
	"errors"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"gorm.io/gorm"
)

func TestService_CreateProduct_Success(t *testing.T) {
	mockProductRepo := &MockProductRepository{}
	mockCategoryRepo := &MockProductCategoryRepository{}
	service := NewService(mockProductRepo, mockCategoryRepo)

	req := &product.CreateProductRequest{
		Name:       "Paracetamol",
		SKU:        "PARA500",
		CategoryID: "cat-1",
		Price:      500000,
		Status:     "active",
	}

	mockCategoryRepo.FindByIDFunc = func(id string) (*product.ProductCategory, error) {
		if id == "cat-1" {
			return &product.ProductCategory{ID: "cat-1", Name: "Medicine"}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockProductRepo.CreateFunc = func(p *product.Product) error {
		p.ID = "prod-1"
		return nil
	}

	mockProductRepo.FindByIDFunc = func(id string) (*product.Product, error) {
		return &product.Product{
			ID:         "prod-1",
			Name:       req.Name,
			CategoryID: req.CategoryID,
			Category:   &product.ProductCategoryRef{ID: "cat-1", Name: "Medicine"},
			CreatedAt:  time.Now(),
		}, nil
	}

	resp, err := service.CreateProduct(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "prod-1" {
		t.Errorf("expected ID prod-1, got %s", resp.ID)
	}
	if resp.Category == nil {
		t.Errorf("expected category to be populated")
	}
}

func TestService_CreateProduct_CategoryNotFound(t *testing.T) {
	mockProductRepo := &MockProductRepository{}
	mockCategoryRepo := &MockProductCategoryRepository{}
	service := NewService(mockProductRepo, mockCategoryRepo)

	req := &product.CreateProductRequest{
		Name:       "Paracetamol",
		SKU:        "PARA500",
		CategoryID: "cat-99",
	}

	mockCategoryRepo.FindByIDFunc = func(id string) (*product.ProductCategory, error) {
		return nil, gorm.ErrRecordNotFound
	}

	_, err := service.CreateProduct(req)
	if !errors.Is(err, ErrProductCategoryNotFound) {
		t.Errorf("expected ErrProductCategoryNotFound, got %v", err)
	}
}

func TestService_ListProducts_Success(t *testing.T) {
	mockProductRepo := &MockProductRepository{}
	mockCategoryRepo := &MockProductCategoryRepository{}
	service := NewService(mockProductRepo, mockCategoryRepo)

	req := &product.ListProductsRequest{Page: 1, PerPage: 10}

	mockProductRepo.ListFunc = func(r *product.ListProductsRequest) ([]product.Product, int64, error) {
		return []product.Product{
			{ID: "p1", Name: "Product 1"},
			{ID: "p2", Name: "Product 2"},
		}, 2, nil
	}

	resp, pagination, err := service.ListProducts(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("expected 2 products, got %d", len(resp))
	}
	if pagination.Total != 2 {
		t.Errorf("expected total 2, got %d", pagination.Total)
	}
}
