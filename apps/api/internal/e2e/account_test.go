package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAccount_SystemFlow(t *testing.T) {
	// Need a logged in user. Admin is fine checking their own profile.
	client := GetAdminClient(t)
	if client.AccessToken == "" {
		t.Fatalf("Client AccessToken is empty after GetAdminClient login")
	}

	// 1. Get Profile (Smoke)
	t.Log("Step 1: Get Profile")
	resp, err := client.Get(fmt.Sprintf("/api/v1/users/%s/profile", client.UserID))
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK for profile, got %d", resp.StatusCode)
	}

	// 2. Update Profile
	t.Log("Step 2: Update Profile")
	// Update name
	newName := fmt.Sprintf("Admin Updated %d", time.Now().Unix())
	updatePayload := map[string]interface{}{
		"name": newName,
		// Assuming we don't need to send everything, or patch logic
	}
	
	resp, err = client.Put(fmt.Sprintf("/api/v1/users/%s/profile", client.UserID), updatePayload)
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Update profile returned %d. Expected 200 OK.", resp.StatusCode)
	} else {
		// Verify update, ignore for now
	}
}
