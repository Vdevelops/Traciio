package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	productservice "github.com/gilabs/crm-healthcare/api/internal/service/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestProductHandler() (*ProductHandler, *MockProductRepository, *MockProductCategoryRepository) {
	mockProductRepo := &MockProductRepository{}
	mockCategoryRepo := &MockProductCategoryRepository{}
	
	svc := productservice.NewService(mockProductRepo, mockCategoryRepo)
	return NewProductHandler(svc), mockProductRepo, mockCategoryRepo
}

func TestProductHandler_Create_Success(t *testing.T) {
	handler, mockProductRepo, mockCategoryRepo := setupTestProductHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/products", handler.Create)
	
	price := int64(500000)
	reqBody := product.CreateProductRequest{
		Name:       "New Product",
		SKU:        "PROD001",
		Price:      price,
		CategoryID: "550e8400-e29b-41d4-a716-446655440001", // Valid UUID
	}
	
	mockCategoryRepo.FindByIDFunc = func(id string) (*product.ProductCategory, error) {
		if id == reqBody.CategoryID {
			return &product.ProductCategory{ID: reqBody.CategoryID, Name: "Category 1"}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	mockProductRepo.CreateFunc = func(p *product.Product) error {
		p.ID = "prod-1"
		return nil
	}
	
	mockProductRepo.FindByIDFunc = func(id string) (*product.Product, error) {
		return &product.Product{
			ID:         id,
			Name:       reqBody.Name,
			SKU:        reqBody.SKU,
			Price:      reqBody.Price,
			CategoryID: reqBody.CategoryID,
			Category:   &product.ProductCategoryRef{ID: reqBody.CategoryID, Name: "Category 1"},
			CreatedAt:  time.Now(),
		}, nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_List_Success(t *testing.T) {
	handler, mockProductRepo, _ := setupTestProductHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/products", handler.List)
	
	mockProductRepo.ListFunc = func(req *product.ListProductsRequest) ([]product.Product, int64, error) {
		return []product.Product{
			{ID: "p1", Name: "Product 1", Price: 1000},
			{ID: "p2", Name: "Product 2", Price: 2000},
		}, 2, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/products?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("expected success true, got %v", response["success"])
	}
}
