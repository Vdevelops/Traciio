package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestOrg_SystemFlow(t *testing.T) {
	client := GetAdminClient(t)

	// 1. Create Division
	divCode := fmt.Sprintf("DIV_%d", time.Now().Unix())
	divPayload := map[string]interface{}{
		"name": "Test Division",
		"code": divCode,
	}

	t.Log("Step 1: Create Division")
	resp, err := client.Post("/api/v1/divisions", divPayload)
	if err != nil {
		t.Fatalf("Failed to create division: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Logf("Create division status %d.", resp.StatusCode)
	}

	// 2. List Divisions (Smoke)
	t.Log("Step 2: List Divisions")
	resp, err = client.Get("/api/v1/divisions")
	if err != nil {
		t.Fatalf("Failed to list divisions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for divisions, got %d", resp.StatusCode)
	}
}
