package ai

import (
	"strings"
	"testing"

	productdomain "github.com/gilabs/crm-healthcare/api/internal/domain/product"
)

type stubAIProductRepo struct {
	products []productdomain.Product
}

func (r stubAIProductRepo) FindByID(id string) (*productdomain.Product, error) {
	for _, productEntity := range r.products {
		if productEntity.ID == id {
			return &productEntity, nil
		}
	}
	return nil, nil
}

func (r stubAIProductRepo) List(req *productdomain.ListProductsRequest) ([]productdomain.Product, int64, error) {
	if req == nil || strings.TrimSpace(req.Search) == "" {
		return r.products, int64(len(r.products)), nil
	}

	search := strings.ToLower(strings.TrimSpace(req.Search))
	results := make([]productdomain.Product, 0)
	for _, productEntity := range r.products {
		name := strings.ToLower(productEntity.Name)
		sku := strings.ToLower(productEntity.SKU)
		if strings.Contains(name, search) || strings.Contains(sku, search) {
			results = append(results, productEntity)
		}
	}
	return results, int64(len(results)), nil
}

func (r stubAIProductRepo) Create(*productdomain.Product) error {
	return nil
}

func (r stubAIProductRepo) Update(*productdomain.Product) error {
	return nil
}

func (r stubAIProductRepo) Delete(string) error {
	return nil
}

func TestStrictProductInterestsRejectsUnknownProducts(t *testing.T) {
	service := &Service{
		productRepo: stubAIProductRepo{
			products: []productdomain.Product{
				{ID: "prod-1", Name: "Amoxicillin 500mg Capsule", SKU: "AMX-500", Status: "active"},
			},
		},
	}

	_, err := service.strictProductInterestsFromParams(map[string]interface{}{
		"product_names": []interface{}{"Mixagrip", "Amoxilin"},
	})

	if err == nil {
		t.Fatal("expected unknown products to be rejected")
	}
	errMsg := err.Error()
	if !strings.Contains(errMsg, "Mixagrip") || !strings.Contains(errMsg, "Amoxilin") {
		t.Fatalf("expected error to mention unresolved product names, got %s", errMsg)
	}
	if !strings.Contains(errMsg, "Product interest tidak disimpan") {
		t.Fatalf("expected error to prevent product interest persistence, got %s", errMsg)
	}
}

func TestStrictProductInterestsResolvesExistingProduct(t *testing.T) {
	service := &Service{
		productRepo: stubAIProductRepo{
			products: []productdomain.Product{
				{ID: "prod-1", Name: "Amoxicillin 500mg Capsule", SKU: "AMX-500", Status: "active"},
			},
		},
	}

	interests, err := service.strictProductInterestsFromParams(map[string]interface{}{
		"product_name": "Amoxicillin",
	})

	if err != nil {
		t.Fatalf("expected existing product to resolve, got %v", err)
	}
	if len(interests) != 1 {
		t.Fatalf("expected one resolved interest, got %d", len(interests))
	}
	if interests[0]["product_id"] != "prod-1" {
		t.Fatalf("expected product ID prod-1, got %s", interests[0]["product_id"])
	}
	if interests[0]["product_name"] != "Amoxicillin 500mg Capsule" {
		t.Fatalf("expected canonical product name, got %s", interests[0]["product_name"])
	}
}
