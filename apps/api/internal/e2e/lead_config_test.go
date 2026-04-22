package e2e

import (
	"net/http"
	"testing"
)

func TestLeadConfig_Smoke(t *testing.T) {
	client := GetAdminClient(t)

	// List Lead Sources
	t.Log("Step 1: List Lead Sources")
	resp, err := client.Get("/api/v1/lead-sources")
	if err != nil {
		t.Fatalf("Failed to list lead sources: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for lead sources, got %d", resp.StatusCode)
	}
}
