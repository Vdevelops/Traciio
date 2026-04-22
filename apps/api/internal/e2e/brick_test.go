package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestBrick_SystemFlow(t *testing.T) {
	client := GetAdminClient(t)

	// 1. Create Brick
	timestamp := time.Now().Unix()
	brickCode := fmt.Sprintf("BRICK_%d", timestamp)
	brickPayload := map[string]interface{}{
		"name": fmt.Sprintf("Test Brick %d", timestamp),
		"code": brickCode,
		"description": "E2E Test Brick",
		"province": fmt.Sprintf("Province %d", timestamp),
		"regency": fmt.Sprintf("Regency %d", timestamp),
	}

	t.Log("Step 1: Create Brick")
	resp, err := client.Post("/api/v1/bricks", brickPayload)
	if err != nil {
		t.Fatalf("Failed to create brick: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
	}

	// 2. List Bricks (Smoke)
	t.Log("Step 2: List Bricks")
	resp, err = client.Get("/api/v1/bricks")
	if err != nil {
		t.Fatalf("Failed to list bricks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for bricks, got %d", resp.StatusCode)
	}
}
