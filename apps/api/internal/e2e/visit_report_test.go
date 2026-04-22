package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitReport_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Test list visit reports endpoint
	resp, err := client.Get("/api/v1/visit-reports")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Visit reports list endpoint should return 200")
}
