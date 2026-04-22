package e2e

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Test list tasks endpoint
	resp, err := client.Get("/api/v1/tasks")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// If not 200, log the response body for debugging
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response body (status %d): %s", resp.StatusCode, string(body))
		
		var errResp map[string]interface{}
		json.Unmarshal(body, &errResp)
		t.Logf("Parsed error: %+v", errResp)
	}

	assert.Equal(t, 200, resp.StatusCode, "Tasks list endpoint should return 200")
}

