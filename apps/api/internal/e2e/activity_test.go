package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivity_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Test list activities endpoint
	resp, err := client.Get("/api/v1/activities")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Activities list endpoint should return 200")
}
