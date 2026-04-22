package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedule_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Test list schedules endpoint
	resp, err := client.Get("/api/v1/schedules")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Schedules list endpoint should return 200")
}
