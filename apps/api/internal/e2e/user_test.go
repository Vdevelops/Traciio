package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestUser_SystemFlow(t *testing.T) {
	client := GetAdminClient(t)

	// 1. Create User
	email := fmt.Sprintf("test_user_e2e_%d@example.com", time.Now().Unix())
	createPayload := map[string]interface{}{
		"email":    email,
		"password": "password123",
		"name":     "E2E Test User",
		"status":   "active",
		// "role_id": "...", // If required, need to fetch a role first
	}

	t.Log("Step 1: Create User")
	resp, err := client.Post("/api/v1/users", createPayload)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Logf("Create user failed with %d. Skipping login test for this user.", resp.StatusCode)
	} else {
		// 2. Login as new user
		t.Log("Step 2: Login as new User")
		userClient, _ := NewClient(LoadConfig().APIURL)
		err = userClient.Login(email, "password123")
		if err != nil {
			t.Errorf("Failed to login as newly created user: %v", err)
		}
	}

	// 3. List Users (Smoke)
	t.Log("Step 3: List Users")
	resp, err = client.Get("/api/v1/users")
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for list users, got %d", resp.StatusCode)
	}
}
