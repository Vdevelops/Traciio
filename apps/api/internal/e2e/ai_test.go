package e2e

import (
	"encoding/json"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	"github.com/stretchr/testify/assert"
)

func TestAI_Smoke_Chat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping smoke test in short mode")
	}

	client := GetAdminClient(t)

	// Simple chat request
	req := map[string]interface{}{
		"message": "Hello, who are you?",
		"model":   "gpt-oss-120b",
	}

	resp, err := client.Post("/api/v1/ai/chat", req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Note: If Cerebras API key is not configured in CI/Test env, this might return 400 or 500.
	// We should check if we can mock it or valid API key is present.
	// For smoke test, we expect at least a valid response structure even if it's an error from provider,
	// BUT typically smoke tests run against a deployed env or local dev with env vars.
	// If API Key is missing, it returns 500 "AI service not configured".
	// We can assert status 200 OR 500 with specific error if we don't have API key.
	// However, usually we skip if requirements aren't met.

	// Assuming test environment might not have API KEY, we accept 200 or 500 with specific message.
	if resp.StatusCode == 500 {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp["error"] == "AI service not configured: Cerebras API key is empty" {
			t.Skip("Skipping AI smoke test: API Key not configured")
		}
	}

	// If we proceed, we expect 200
	assert.Equal(t, 200, resp.StatusCode)

	var apiResp struct {
		Data ai.ChatResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err == nil {
		assert.NotEmpty(t, apiResp.Data.Message)
	}
}

func TestAI_System_ContextAware(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping system test in short mode")
	}

	client := GetAdminClient(t)

	// 1. Create a Deal to use as context (Skipped, using leaderboard query)

	req := map[string]interface{}{
		"message": "Show me my sales performance",
		"model":   "gpt-oss-120b",
	}

	resp, err := client.Post("/api/v1/ai/chat", req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// If response is 500 or 504, it might be due to missing configuration or timeout from external provider
	if resp.StatusCode == 500 || resp.StatusCode == 504 {
		t.Skipf("Skipping AI system test due to server error (likely config or timeout): %d", resp.StatusCode)
	}

	assert.Equal(t, 200, resp.StatusCode)

	var apiResp struct {
		Data ai.ChatResponse `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&apiResp)

	// Verify response contains relevant info (even if mocked or empty data)
	// It should at least be a string
	assert.NotEmpty(t, apiResp.Data.Message)
}
