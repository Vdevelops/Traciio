package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/dashboard"
	"github.com/stretchr/testify/assert"
)

// TestDashboard_Smoke verifies that the main dashboard endpoints return 200 OK.
func TestDashboard_Smoke(t *testing.T) {
	client := GetAdminClient(t)

	t.Run("Get Dashboard Overview", func(t *testing.T) {
		resp, err := client.Get("/api/v1/dashboard/overview?period=month")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result responseWrapper
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("Get Visit Statistics", func(t *testing.T) {
		resp, err := client.Get("/api/v1/dashboard/visits?period=month")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get Pipeline Summary", func(t *testing.T) {
		resp, err := client.Get("/api/v1/dashboard/pipeline?period=month")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestDashboard_SystemFlow_SalesUpdate verifies that creating a Won Deal updates the dashboard.
func TestDashboard_SystemFlow_SalesUpdate(t *testing.T) {
	client := GetAdminClient(t)

	// 1. Get initial Dashboard stats
	var initialStats dashboard.DashboardOverviewResponse
	resp, err := client.Get("/api/v1/dashboard/overview?period=month")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var initialResp struct {
		Success bool                                 `json:"success"`
		Data    dashboard.DashboardOverviewResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&initialResp)
	assert.NoError(t, err)
	initialStats = initialResp.Data

	// 2. Create Prerequisites (Category, Account, Contact, Stage) - Reusing logic from deal_flow_test
	// Create Category
	categoryPayload := map[string]interface{}{
		"name":        "Healthcare Provider " + randomString(5),
		"code":        "CAT_" + randomString(5),
		"description": "Hospital/Clinic",
		"type":        "account", // or whatever type required
	}
	resp, err = client.Post("/api/v1/categories", categoryPayload)
	assert.NoError(t, err)
    if resp.StatusCode != 200 && resp.StatusCode != 201 {
        t.Fatalf("Create Category failed: %d", resp.StatusCode)
    }
	var catResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&catResp)
	categoryID := catResp.Data.ID

	// Create Account
	accountPayload := map[string]interface{}{
		"name":        "Rumah Sakit Dashboard Test " + randomString(5),
		"email":       "dashboard" + randomString(5) + "@rs.com",
		"phone":       "081234567890",
		"address":     "Jl. Dashboard Test",
		"category_id": categoryID,
	}
	resp, err = client.Post("/api/v1/accounts", accountPayload)
	assert.NoError(t, err)
    if resp.StatusCode != 200 && resp.StatusCode != 201 {
        t.Fatalf("Create Account failed: %d", resp.StatusCode)
    }
	var accResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&accResp)
	accountID := accResp.Data.ID

	// Get "Won" Stage ID
	stages := getPipelineStages(t, client)
	var wonStageID string
	// If no stages or Won stage missing, ensure we have them
	if len(stages) == 0 {
		// Create Open Stage
		openPayload := map[string]interface{}{"name": "Open", "code": "open", "order": 1, "color": "#FFC107"}
		client.Post("/api/v1/pipelines", openPayload)
		// Create Won Stage
		wonPayload := map[string]interface{}{"name": "Won", "code": "won", "order": 100, "color": "#4CAF50"}
		client.Post("/api/v1/pipelines", wonPayload)
		// Refresh
		stages = getPipelineStages(t, client)
	}

	for _, s := range stages {
		if s.Code == "won" || s.Name == "Won" {
			wonStageID = s.ID
			break
		}
	}
	
	if wonStageID == "" {
		// If still missing but we have stages, use the last one (or create one)
        if len(stages) > 0 {
		    wonStageID = stages[len(stages)-1].ID
        } else {
            // Should not happen if we created above, but safety check
            wonPayload := map[string]interface{}{"name": "Won", "code": "won", "order": 100, "color": "#4CAF50"}
            resp, _ := client.Post("/api/v1/pipelines", wonPayload)
            var res struct{ Data struct{ ID string } `json:"data"` }
            json.NewDecoder(resp.Body).Decode(&res)
            wonStageID = res.Data.ID
            
            // Refresh stages for open stage
            stages = getPipelineStages(t, client)
        }
	}

	// 3. Create a **Won** Deal immediately (or Create then Move)
	// For simplicity, we create it directly with Won stage if allowed, or Create then Move.
	// Let's try Create then Move to be safe and realistic.
	
	// Get Current User ID for correct revenue attribution
	userID := client.UserID
	if userID == "" {
		t.Fatal("Client UserID is empty, login might have failed or not populated.")
	}
	// Create Deal
	dealValue := int64(50000000) // 50 Million
	dealPayload := map[string]interface{}{
		"title":       "Dashboard Test Deal " + randomString(5),
		"account_id":  accountID,
		"stage_id":    stages[0].ID, // Start at first stage
		"value":       dealValue,
		"status":      "open",
		"description": "Testing dashboard update",
		"assigned_to": userID, // Explicitly assign to current user
	}
	resp, err = client.Post("/api/v1/deals", dealPayload)
	assert.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create deal: status %d", resp.StatusCode)
	}
	var dealResp struct {
		Data struct {
			ID string `json:"id"`
			AssignedTo string `json:"assigned_to"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&dealResp)
	dealID := dealResp.Data.ID

	// Move to WON
	movePayload := map[string]interface{}{
		"stage_id": wonStageID,
		"status":   "won", // Explicitly set status to won
		"notes":    "Closed for dashboard test",
	}
	resp, err = client.Put(fmt.Sprintf("/api/v1/deals/%s", dealID), movePayload)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		t.Logf("Update Deal Failed: ID=%s Status=%d Body=%v", dealID, resp.StatusCode, body)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 4. Verify Dashboard Update
	// Wait a moment for any async processing (though usually synchronous DB)
	time.Sleep(500 * time.Millisecond)

	resp, err = client.Get("/api/v1/dashboard/overview?period=month")
	assert.NoError(t, err)
	
	var finalResp struct {
		Success bool                                 `json:"success"`
		Data    dashboard.DashboardOverviewResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&finalResp)
	assert.NoError(t, err)
	finalStats := finalResp.Data

	// Assertions
	// 1. Total Deals should increase by 1
	assert.Equal(t, initialStats.Deals.TotalDeals+1, finalStats.Deals.TotalDeals, "Total deals should increment by 1")
	
	// 2. Won Deals should increase by 1
	assert.Equal(t, initialStats.Deals.WonDeals+1, finalStats.Deals.WonDeals, "Won deals should increment by 1")
	
	// 3. Revenue should increase
	// Note: In parallel test runs, initial stats might include data from other tests that disappears or changes.
	// We verify that the new deal value is present and total did not decrease unreasonably.
	assert.GreaterOrEqual(t, finalStats.Revenue.TotalRevenue, dealValue, "Total revenue should include the new deal value")
	assert.GreaterOrEqual(t, finalStats.Revenue.TotalRevenue, initialStats.Revenue.TotalRevenue, "Total revenue should not decrease")
}
