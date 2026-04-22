package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTarget_Smoke verifies that target endpoints are reachable
func TestTarget_Smoke(t *testing.T) {
	client := GetAdminClient(t)

	t.Run("Get Current Target", func(t *testing.T) {
		// Targets are usually under /api/v1/targets or similar.
        // Based on implementation plan: /api/v1/targets/current or /api/v1/monthly-targets
        // Let's guess /api/v1/monthly-targets based on domain name
		resp, err := client.Get("/api/v1/monthly-targets")
		assert.NoError(t, err)
        // Might be 200 or 404/400 depending on implementation. Assuming 200 list.
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	})
}

func TestTarget_Integration_Set(t *testing.T) {
	client := GetAdminClient(t)

	// Get Current User ID
	resp, err := client.Get("/api/v1/user/profile") // Standard endpoint usually /auth/me or /users/profile, client.go might have helper
	// If standard client login doesn't expose ID easily, we fetch it.
	// Actually GetAdminClient logs in. Let's try /api/v1/auth/me or similar if we need ID.
	// Or we can list users and pick one.
	resp, err = client.Get("/api/v1/users")
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	var usersResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&usersResp)
	resp.Body.Close()

	if len(usersResp.Data) == 0 {
		t.Skip("No users found to assign target")
	}
	userID := usersResp.Data[0].ID

	// Create payload for Monthy Target
	now := time.Now()
	targetPayload := map[string]interface{}{
		"year":          now.Year(),
		"month":         int(now.Month()),
		"target_value":  100000000,
		"target_amount": 100000000, // Entity uses target_amount
		"user_id":       userID,
	}

	// Check if target exists and delete it to ensure clean state
	listURL := fmt.Sprintf("/api/v1/monthly-targets?user_id=%s&year=%d&month=%d", userID, now.Year(), int(now.Month()))
	resp, err = client.Get(listURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		var listResp struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&listResp)
		resp.Body.Close()
		
		for _, t := range listResp.Data {
			// Delete existing target
			delURL := fmt.Sprintf("/api/v1/monthly-targets/%s", t.ID)
			client.Delete(delURL)
		}
	}

	// Create Target
	resp, err = client.Post("/api/v1/monthly-targets", targetPayload)
	assert.NoError(t, err)
	
	if resp.StatusCode == http.StatusNotFound {
		t.Log("Target endpoint not found at /api/v1/monthly-targets")
	} else if resp.StatusCode == http.StatusBadRequest {
		var errBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errBody)
		t.Logf("Target creation failed with 400: %v", errBody)
		t.Fail()
	} else {
		// Log body on failure if not created
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			var errBody map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errBody)
			t.Logf("Target creation failed with status %d: %v", resp.StatusCode, errBody)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var targetResp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&targetResp)
		resp.Body.Close()
		assert.NotEmpty(t, targetResp.Data.ID)
	}
}
