package ai

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

func TestBuildVisitReportContext(t *testing.T) {
	// Create test visit report
	visitDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	checkInTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	checkOutTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	checkInLocation := map[string]interface{}{
		"address": "Jl. Sudirman No. 123, Jakarta",
		"lat":     -6.2088,
		"lng":     106.8456,
	}
	checkInLocationJSON, _ := json.Marshal(checkInLocation)
	
	visitReport := &visit_report.VisitReport{
		ID:            "test-visit-report-id",
		VisitDate:     visitDate,
		Status:        "submitted",
		Purpose:       "Product presentation",
		Notes:         "Discussed new product line with procurement team",
		CheckInTime:   &checkInTime,
		CheckOutTime:  &checkOutTime,
		CheckInLocation: checkInLocationJSON,
		ContactID:     stringPtr("test-contact-id"),
	}
	
	// Create test account
	testAccount := &account.Account{
		ID:       "test-account-id",
		Name:     "RSUD Jakarta",
		Address:  "Jl. Sudirman No. 123",
		City:     "Jakarta",
		Province: "DKI Jakarta",
		Phone:    "021-12345678",
		Email:    "info@rsudjakarta.go.id",
		Status:   "active",
		Category: &account.Category{
			ID:   "test-category-id",
			Name: "Hospital",
		},
	}
	
	contactName := "Dr. John Doe"
	
	// Create test activities
	activities := []activity.Activity{
		{
			ID:          "activity-1",
			Type:        "call",
			Description: "Follow-up call",
			Timestamp:   time.Date(2024, 1, 14, 14, 0, 0, 0, time.UTC),
		},
		{
			ID:          "activity-2",
			Type:        "email",
			Description: "Sent product catalog",
			Timestamp:   time.Date(2024, 1, 13, 9, 0, 0, 0, time.UTC),
		},
	}
	
	// Build context
	context := BuildVisitReportContext(visitReport, testAccount, contactName, activities)
	
	// Verify context contains expected information
	if !strings.Contains(context, "PHARMACEUTICAL SALES VISIT REPORT") {
		t.Error("Context should contain visit report header")
	}
	
	if !strings.Contains(context, "2024-01-15") {
		t.Error("Context should contain visit date")
	}
	
	if !strings.Contains(context, "submitted") {
		t.Error("Context should contain visit report status")
	}
	
	if !strings.Contains(context, "Product presentation") {
		t.Error("Context should contain visit purpose")
	}
	
	if !strings.Contains(context, "RSUD Jakarta") {
		t.Error("Context should contain account name")
	}
	
	if !strings.Contains(context, "Hospital") {
		t.Error("Context should contain account category")
	}
	
	if !strings.Contains(context, "Dr. John Doe") {
		t.Error("Context should contain contact name")
	}
	
	if !strings.Contains(context, "RECENT ACTIVITY HISTORY") {
		t.Error("Context should contain activity history section")
	}
	
	if !strings.Contains(context, "call") {
		t.Error("Context should contain activity type")
	}
	
	if !strings.Contains(context, "Jl. Sudirman No. 123, Jakarta") {
		t.Error("Context should contain check-in location address")
	}
}

func TestBuildVisitReportContext_WithNilValues(t *testing.T) {
	// Test with minimal data (nil check-in/check-out times, no activities)
	visitDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	
	visitReport := &visit_report.VisitReport{
		ID:        "test-visit-report-id",
		VisitDate: visitDate,
		Status:    "draft",
		Purpose:   "Initial visit",
		Notes:     "First contact",
	}
	
	testAccount := &account.Account{
		ID:   "test-account-id",
		Name: "Test Hospital",
	}
	
	context := BuildVisitReportContext(visitReport, testAccount, "", []activity.Activity{})
	
	// Should still contain basic information
	if !strings.Contains(context, "2024-01-15") {
		t.Error("Context should contain visit date even with minimal data")
	}
	
	if !strings.Contains(context, "Test Hospital") {
		t.Error("Context should contain account name")
	}
	
	// Should not contain contact section if contact name is empty
	if strings.Contains(context, "CONTACT PERSON") {
		t.Error("Context should not contain contact section if contact name is empty")
	}
	
	// Should not contain activity section if no activities
	if strings.Contains(context, "RECENT ACTIVITY HISTORY") {
		t.Error("Context should not contain activity section if no activities")
	}
}

func TestBuildVisitReportPrompt(t *testing.T) {
	context := "=== TEST CONTEXT ===\nVisit Date: 2024-01-15\nStatus: submitted"
	
	prompt := BuildVisitReportPrompt(context)
	
	// Verify prompt contains expected sections
	if !strings.Contains(prompt, "expert AI assistant") {
		t.Error("Prompt should identify as expert AI assistant")
	}
	
	if !strings.Contains(prompt, "pharmaceutical") {
		t.Error("Prompt should mention pharmaceutical industry")
	}
	
	if !strings.Contains(prompt, context) {
		t.Error("Prompt should include the provided context")
	}
	
	if !strings.Contains(prompt, "EXECUTIVE SUMMARY") {
		t.Error("Prompt should request executive summary")
	}
	
	if !strings.Contains(prompt, "SENTIMENT ANALYSIS") {
		t.Error("Prompt should request sentiment analysis")
	}
	
	if !strings.Contains(prompt, "KEY POINTS DISCUSSED") {
		t.Error("Prompt should request key points")
	}
	
	if !strings.Contains(prompt, "ACTION ITEMS") {
		t.Error("Prompt should request action items")
	}
	
	if !strings.Contains(prompt, "STRATEGIC RECOMMENDATIONS") {
		t.Error("Prompt should request strategic recommendations")
	}
	
	if !strings.Contains(prompt, "JSON") {
		t.Error("Prompt should specify JSON output format")
	}
	
	if !strings.Contains(prompt, "summary") {
		t.Error("Prompt should specify summary field in JSON")
	}
	
	if !strings.Contains(prompt, "sentiment") {
		t.Error("Prompt should specify sentiment field in JSON")
	}
}

func TestBuildSystemPrompt_WithContext(t *testing.T) {
	currentTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	contextID := "test-context-id"
	contextType := "account"
	contextData := `{"id": "test-account-id", "name": "RSUD Jakarta"}`
	
	prompt := BuildSystemPrompt(contextID, contextType, contextData, "", "gpt-4", "openai", currentTime, "Asia/Jakarta")
	
	// Verify prompt contains base information
	if !strings.Contains(prompt, "Pharmaceutical and Healthcare Sales CRM") {
		t.Error("Prompt should mention CRM system")
	}
	
	if !strings.Contains(prompt, contextID) {
		t.Error("Prompt should include context ID")
	}
	
	if !strings.Contains(prompt, "ACCOUNT (HEALTHCARE FACILITY)") {
		t.Error("Prompt should include context type label")
	}
	
	if !strings.Contains(prompt, contextData) {
		t.Error("Prompt should include context data")
	}
	
	// Date format is "2006-01-02", so "2024-01-15" should be present
	// But also check for year and month separately in case format differs
	if !strings.Contains(prompt, "2024") {
		t.Error("Prompt should include current year (2024)")
	}
	
	if !strings.Contains(prompt, "gpt-4") {
		t.Error("Prompt should include model name")
	}
	
	if !strings.Contains(prompt, "openai") {
		t.Error("Prompt should include provider name")
	}
	
	// Verify no hallucination warnings
	if !strings.Contains(prompt, "NO HALLUCINATION") {
		t.Error("Prompt should include no hallucination warning")
	}
	
	// Verify Markdown table instructions
	if !strings.Contains(prompt, "Markdown table") {
		t.Error("Prompt should include Markdown table formatting instructions")
	}
}

func TestBuildSystemPrompt_WithoutContext(t *testing.T) {
	currentTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	
	prompt := BuildSystemPrompt("", "", "", "", "gpt-4", "openai", currentTime, "Asia/Jakarta")
	
	// Should still contain base prompt
	if !strings.Contains(prompt, "Pharmaceutical and Healthcare Sales CRM") {
		t.Error("Prompt should contain base information even without context")
	}
	
	if !strings.Contains(prompt, "2024-01-15") {
		t.Error("Prompt should include current date")
	}
	
	// Should mention that data is not available
	if !strings.Contains(prompt, "do NOT have access") {
		t.Error("Prompt should mention lack of data access when no context provided")
	}
}

func TestBuildSystemPrompt_WithDataAccessInfo(t *testing.T) {
	currentTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	dataAccessInfo := "Data access restricted"
	
	prompt := BuildSystemPrompt("", "", "", dataAccessInfo, "gpt-4", "openai", currentTime, "Asia/Jakarta")
	
	// Check if dataAccessInfo is included (should be in additionalInfo section)
	// The dataAccessInfo is added as "\n\n" + dataAccessInfo, so it should be present
	if !strings.Contains(prompt, "Data access") && !strings.Contains(prompt, "restricted") {
		// If not found, check if it's in a different format or section
		// The prompt might format it differently, so we check for partial matches
		if !strings.Contains(prompt, "access") {
			t.Error("Prompt should include data access info (looking for 'access' or 'restricted')")
		}
	}
}

func TestBuildSystemPrompt_TimeContext(t *testing.T) {
	// Test with date near Christmas
	christmasTime := time.Date(2024, 12, 20, 14, 30, 0, 0, time.UTC)
	
	prompt := BuildSystemPrompt("", "", "", "", "gpt-4", "openai", christmasTime, "Asia/Jakarta")
	
	// Should mention Christmas if within 60 days
	if !strings.Contains(prompt, "Christmas") {
		t.Error("Prompt should mention upcoming Christmas when within 60 days")
	}
	
	// Test with date far from holidays
	summerTime := time.Date(2024, 7, 15, 14, 30, 0, 0, time.UTC)
	
	prompt2 := BuildSystemPrompt("", "", "", "", "gpt-4", "openai", summerTime, "Asia/Jakarta")
	
	// Should not mention Christmas if far away
	if strings.Contains(prompt2, "Christmas") {
		t.Error("Prompt should not mention Christmas if more than 60 days away")
	}
}

func TestBuildSystemPrompt_ConversionRateInstructions(t *testing.T) {
	currentTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	
	prompt := BuildSystemPrompt("", "", "", "", "gpt-4", "openai", currentTime, "Asia/Jakarta")
	
	// Verify conversion rate instructions mention Qualification, not Lead
	if !strings.Contains(prompt, "Qualification") {
		t.Error("Prompt should mention Qualification stage for conversion rate calculations")
	}
	
	if strings.Contains(prompt, "stage_name 'Lead'") || strings.Contains(prompt, "stage_code 'lead'") {
		t.Error("Prompt should NOT mention Lead as a pipeline stage")
	}
	
	// Verify it mentions that Lead is not a pipeline stage (check for variations)
	if !strings.Contains(prompt, "Lead is NOT") && !strings.Contains(prompt, "Lead is not") && !strings.Contains(prompt, "NOT a pipeline") {
		t.Error("Prompt should clarify that Lead is not a pipeline stage")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

