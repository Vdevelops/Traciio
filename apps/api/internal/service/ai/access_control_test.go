package ai

import (
	"strings"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

func testUserContext(permissions []string, scopes map[string]roledomain.ScopeType) *domainauth.UserContext {
	permMap := make(map[string]bool, len(permissions))
	for _, permission := range permissions {
		permMap[permission] = true
	}
	return &domainauth.UserContext{
		UserID:        "user-1",
		RoleCode:      "sales",
		Permissions:   permMap,
		Scopes:        scopes,
		TeamMemberIDs: []string{"user-1", "user-2"},
	}
}

func TestCanRunToolCallUsesSpecificTaskStatusPermission(t *testing.T) {
	service := &Service{}

	userCtx := testUserContext([]string{"tasks.start"}, nil)
	call := &ToolCall{
		Tool: "update_task_status",
		Params: map[string]interface{}{
			"status": "in progress",
		},
	}

	if !service.canRunToolCall(call, userCtx) {
		t.Fatal("expected tasks.start permission to allow in progress status update")
	}

	userCtx = testUserContext([]string{"tasks.edit"}, nil)
	if service.canRunToolCall(call, userCtx) {
		t.Fatal("expected tasks.edit alone to be rejected for in progress status update")
	}
}

func TestCanRunToolCallRequiresLeadConvertForConvertedStatus(t *testing.T) {
	service := &Service{}
	call := &ToolCall{
		Tool: "update_lead_status",
		Params: map[string]interface{}{
			"status": "converted",
		},
	}

	userCtx := testUserContext([]string{"leads.edit"}, nil)
	if service.canRunToolCall(call, userCtx) {
		t.Fatal("expected leads.edit alone to be rejected for converted lead status")
	}

	userCtx = testUserContext([]string{"leads.convert"}, nil)
	if !service.canRunToolCall(call, userCtx) {
		t.Fatal("expected leads.convert to allow converted lead status")
	}
}

func TestCanRunToolCallAllowsCreateActivityFromLeadOrVisitPermission(t *testing.T) {
	service := &Service{}
	call := &ToolCall{Tool: "create_activity", Params: map[string]interface{}{}}

	userCtx := testUserContext([]string{"leads.edit"}, nil)
	if !service.canRunToolCall(call, userCtx) {
		t.Fatal("expected leads.edit permission to allow create_activity")
	}

	userCtx = testUserContext([]string{"visit-reports.create"}, nil)
	if !service.canRunToolCall(call, userCtx) {
		t.Fatal("expected visit-reports.create permission to allow create_activity")
	}

	userCtx = testUserContext([]string{"tasks.create"}, nil)
	if service.canRunToolCall(call, userCtx) {
		t.Fatal("expected unrelated task permission to be rejected for create_activity")
	}
}

func TestCanAccessOwnerFollowsResourceScope(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"pipeline.view"}, map[string]roledomain.ScopeType{
		"deals": roledomain.ScopeTeam,
	})

	if !service.canAccessOwner(userCtx, "deal", "user-2") {
		t.Fatal("expected team-scoped user to access team owner's deal")
	}
	if service.canAccessOwner(userCtx, "deal", "user-3") {
		t.Fatal("expected team-scoped user to be denied outside-team deal")
	}

	userCtx = testUserContext([]string{"pipeline.view"}, map[string]roledomain.ScopeType{
		"deals": roledomain.ScopeGlobal,
	})
	if !service.canAccessOwner(userCtx, "deal", "user-3") {
		t.Fatal("expected global-scoped user to access any owner")
	}
}

func TestBrickManagementUsesRealBrickPermission(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"bricks.view"}, map[string]roledomain.ScopeType{
		"bricks": roledomain.ScopeTeam,
	})

	if !service.hasDataPermission("brick_management", userCtx) {
		t.Fatal("expected bricks.view permission to allow AI brick management data")
	}
}

func TestAIDataPrivacyDefaultsUnknownDataTypeToDenied(t *testing.T) {
	if aiDataPrivacyAllows("unknown_module", ai_settings.DefaultDataPrivacySettings()) {
		t.Fatal("expected unknown AI data type to be denied by privacy helper")
	}
}

func TestTargetScopeFromMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "defaults to user scope",
			message: "berikan aku data target monthly sales bulan mei",
			want:    "user",
		},
		{
			name:    "detects brick scope",
			message: "tampilkan target brick bulan mei",
			want:    "brick",
		},
		{
			name:    "detects group scope",
			message: "total target per group untuk bulan mei",
			want:    "group",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := targetScopeFromMessage(tt.message); got != tt.want {
				t.Fatalf("targetScopeFromMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTargetMonthsFromMessage(t *testing.T) {
	now := time.Date(2026, time.July, 21, 0, 0, 0, 0, time.UTC)
	months := targetMonthsFromMessage(
		"berikan aku data target untuk sales representative bulan april dan bulan mei",
		now,
	)

	if len(months) != 2 {
		t.Fatalf("expected 2 target months, got %d", len(months))
	}
	if months[0].Year() != 2026 || months[0].Month() != time.April {
		t.Fatalf("expected first month April 2026, got %s", months[0].Format("2006-01"))
	}
	if months[1].Year() != 2026 || months[1].Month() != time.May {
		t.Fatalf("expected second month May 2026, got %s", months[1].Format("2006-01"))
	}
}

func TestTargetOwnerFilterFromMessage(t *testing.T) {
	filter := targetOwnerFilterFromMessage(
		"berikan aku data target untuk sales representative bulan april dan bulan mei",
	)
	if filter != "sales representative" {
		t.Fatalf("expected sales representative owner filter, got %q", filter)
	}

	filter = targetOwnerFilterFromMessage("berikan target untuk semua sales representative bulan april")
	if filter != "" {
		t.Fatalf("expected empty owner filter for all sales representative query, got %q", filter)
	}
}

func TestTaskSearchTermFromMessageForStatusUpdate(t *testing.T) {
	got := taskSearchTermFromMessage("ubah status task follow up katalog kebutuhan formulary pt healthcare indonesia ke completed")
	want := "follow up katalog kebutuhan formulary pt healthcare indonesia"
	if got != want {
		t.Fatalf("taskSearchTermFromMessage() = %q, want %q", got, want)
	}
}

func TestNormalizeTaskStatusForTool(t *testing.T) {
	if got := normalizeTaskStatusForTool("done"); got != "completed" {
		t.Fatalf("normalizeTaskStatusForTool(done) = %q, want completed", got)
	}
	if got := normalizeTaskStatusForTool("selesai"); got != "completed" {
		t.Fatalf("normalizeTaskStatusForTool(selesai) = %q, want completed", got)
	}
}

func TestIsMyProductSalesIntentHandlesIndonesianSalesWording(t *testing.T) {
	message := "berikan saya data pejuaalan saya pada bulan juli dari product paling sering saya jual hingga product yang saya jarang jual"
	if !isMyProductSalesIntent(message) {
		t.Fatal("expected Indonesian product sales wording with typo to be detected as user product sales intent")
	}
}

func TestProductSalesAnalyticsIntentHandlesGenericSalesQuery(t *testing.T) {
	message := "berikan data penjualan pada bulan ini"
	if !isProductSalesAnalyticsIntent(message) {
		t.Fatal("expected generic sales query to be detected as product sales analytics intent")
	}
}

func TestProductSalesScopeKindFollowsRBACScope(t *testing.T) {
	service := &Service{}

	adminCtx := testUserContext([]string{"product-analytics.view"}, map[string]roledomain.ScopeType{
		"sales-overview": roledomain.ScopeGlobal,
	})
	if got := service.productSalesScopeKind(adminCtx); got != "global" {
		t.Fatalf("expected global product sales scope for admin-like context, got %q", got)
	}

	managerCtx := testUserContext([]string{"product-analytics.view"}, map[string]roledomain.ScopeType{
		"sales-overview": roledomain.ScopeTeam,
	})
	if got := service.productSalesScopeKind(managerCtx); got != "team" {
		t.Fatalf("expected team product sales scope for manager-like context, got %q", got)
	}
}

func TestBuildNoUserProductSalesMessageDoesNotSayNoAccess(t *testing.T) {
	message := buildNoUserProductSalesMessage("Juli 2026")
	if !strings.Contains(message, "belum menjual") {
		t.Fatalf("expected no-sales message to explain empty sales result, got %q", message)
	}
	if strings.Contains(strings.ToLower(message), "tidak memiliki akses") {
		t.Fatalf("expected no-sales message not to mention access denial, got %q", message)
	}
}
