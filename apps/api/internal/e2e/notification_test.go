package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotification_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Test list notifications endpoint
	resp, err := client.Get("/api/v1/notifications")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Notifications list endpoint should return 200")

	// Test unread count endpoint
	resp2, err := client.Get("/api/v1/notifications/unread-count")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp2.Body.Close()

	assert.Equal(t, 200, resp2.StatusCode, "Notifications unread count endpoint should return 200")
}
