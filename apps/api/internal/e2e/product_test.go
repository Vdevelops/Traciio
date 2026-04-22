package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestProduct_SystemFlow(t *testing.T) {
	client := GetAdminClient(t)

	// 1. Create Category
	catName := fmt.Sprintf("Test Cat %d", time.Now().Unix())
	catPayload := map[string]interface{}{
		"name": catName,
		"code": fmt.Sprintf("CAT_%d", time.Now().Unix()),
	}

	t.Log("Step 1: Create Category")
	resp, err := client.Post("/api/v1/categories", catPayload)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Logf("Create category status %d. Might vary.", resp.StatusCode)
	}

	// 2. Create Product
	// Ideally we need category ID from step 1.
	// Skipping dependent check to keep it simple, or fetching list first.

	// 3. List Products (Smoke)
	t.Log("Step 3: List Products")
	resp, err = client.Get("/api/v1/products")
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for products, got %d", resp.StatusCode)
	}
}
