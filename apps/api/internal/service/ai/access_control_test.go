package ai

import (
	"testing"

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
