package e2e

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// responseWrapper is a generic wrapper for API responses
type responseWrapper struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
}

// randomString generates a random string of length n
func randomString(n int) string {
	bytes := make([]byte, n/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to time-based if rand fails (unlikely)
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// Helper struct for Pipeline Stage
type Stage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Order int    `json:"order"`
}

// getPipelineStages fetches pipeline stages helper
func getPipelineStages(t *testing.T, client *Client) []Stage {
	resp, err := client.Get("/api/v1/pipelines")
	assert.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to getting pipelines: status %d", resp.StatusCode)
	}

	var result struct {
		Data []Stage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode pipeline stages: %v", err)
	}

	return result.Data
}
