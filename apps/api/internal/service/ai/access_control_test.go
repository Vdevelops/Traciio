package ai

import (
	"strings"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	aidomain "github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	permissiondomain "github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	productdomain "github.com/gilabs/crm-healthcare/api/internal/domain/product"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
	userdomain "github.com/gilabs/crm-healthcare/api/internal/domain/user"
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

type stubAIUserRepo struct {
	user *userdomain.User
}

func (r stubAIUserRepo) FindByID(string) (*userdomain.User, error) {
	return r.user, nil
}

func (r stubAIUserRepo) FindByEmail(string) (*userdomain.User, error) {
	return nil, nil
}

func (r stubAIUserRepo) List(*userdomain.ListUsersRequest) ([]userdomain.User, int64, error) {
	return nil, 0, nil
}

func (r stubAIUserRepo) Create(*userdomain.User) error {
	return nil
}

func (r stubAIUserRepo) Update(*userdomain.User) error {
	return nil
}

func (r stubAIUserRepo) Delete(string) error {
	return nil
}

func (r stubAIUserRepo) CountUsersByRoleID(string) (int64, error) {
	return 0, nil
}

func (r stubAIUserRepo) GetUsersByGroupID(string) ([]userdomain.User, error) {
	return nil, nil
}

func (r stubAIUserRepo) GetUsersByBrickID(string) ([]userdomain.User, error) {
	return nil, nil
}

func (r stubAIUserRepo) GetUsersByRoleID(string) ([]string, error) {
	return nil, nil
}

type stubAIRoleRepo struct {
	roles  []roledomain.Role
	scopes []roledomain.RoleScope
}

func (r stubAIRoleRepo) FindByID(id string) (*roledomain.Role, error) {
	for _, roleEntity := range r.roles {
		if roleEntity.ID == id {
			return &roleEntity, nil
		}
	}
	return nil, nil
}

func (r stubAIRoleRepo) FindByCode(code string) (*roledomain.Role, error) {
	for _, roleEntity := range r.roles {
		if roleEntity.Code == code {
			return &roleEntity, nil
		}
	}
	return nil, nil
}

func (r stubAIRoleRepo) List() ([]roledomain.Role, error) {
	return r.roles, nil
}

func (r stubAIRoleRepo) Create(*roledomain.Role) error {
	return nil
}

func (r stubAIRoleRepo) Update(*roledomain.Role) error {
	return nil
}

func (r stubAIRoleRepo) Delete(string) error {
	return nil
}

func (r stubAIRoleRepo) AssignPermissions(string, []string) error {
	return nil
}

func (r stubAIRoleRepo) GetPermissions(string) ([]string, error) {
	return nil, nil
}

func (r stubAIRoleRepo) GetScopesByRoleID(string) ([]roledomain.RoleScope, error) {
	return r.scopes, nil
}

func (r stubAIRoleRepo) UpsertScopes(string, []roledomain.RoleScopeItem) error {
	return nil
}

func TestEnsureUserContextHydratesCustomRolePermissionsAndScopes(t *testing.T) {
	service := &Service{
		userRepo: stubAIUserRepo{user: &userdomain.User{
			ID:     "user-1",
			Email:  "manager@example.com",
			RoleID: "role-custom-manager",
			Role: &roledomain.Role{
				ID:   "role-custom-manager",
				Code: "custom_sales_manager",
				Permissions: []permissiondomain.Permission{
					{Code: "ai-chatbot.view"},
					{Code: "leads.view"},
				},
			},
		}},
		roleRepo: stubAIRoleRepo{scopes: []roledomain.RoleScope{
			{RoleID: "role-custom-manager", Resource: "leads", Scope: roledomain.ScopeTeam},
		}},
	}

	userCtx := service.ensureUserContext("user-1", nil)

	if userCtx.RoleCode != "custom_sales_manager" {
		t.Fatalf("expected current DB role code, got %q", userCtx.RoleCode)
	}
	if !userCtx.HasPermission("leads.view") || !userCtx.HasPermission("ai-chatbot.view") {
		t.Fatalf("expected role permissions to be hydrated, got %#v", userCtx.Permissions)
	}
	if !userCtx.IsTeamScope("leads") {
		t.Fatalf("expected team scope for leads, got %s", userCtx.GetScope("leads"))
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

func TestCanRunToolCallAllowsScheduleEdit(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"schedules.edit"}, nil)
	call := &ToolCall{Tool: "update_schedule", Params: map[string]interface{}{"schedule_id": "schedule-1", "scheduled_at": "2026-07-24"}}

	if !service.canRunToolCall(call, userCtx) {
		t.Fatal("expected schedules.edit permission to allow update_schedule")
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

func TestAllPromptedAIToolsAreRegisteredAndPermissionGated(t *testing.T) {
	service := &Service{}
	cases := []struct {
		tool       string
		permission string
	}{
		{tool: "create_task", permission: "tasks.create"},
		{tool: "create_lead", permission: "leads.create"},
		{tool: "create_activity", permission: "leads.edit"},
		{tool: "create_product_interest", permission: "leads.edit"},
		{tool: "create_visit_report", permission: "visit-reports.create"},
		{tool: "upsert_lead_bant", permission: "leads.edit"},
		{tool: "create_deal", permission: "pipeline.create"},
		{tool: "create_schedule", permission: "schedules.create"},
		{tool: "update_schedule", permission: "schedules.edit"},
		{tool: "create_route", permission: "route-optimization.create"},
		{tool: "update_task_status", permission: "tasks.edit"},
		{tool: "update_lead_status", permission: "leads.edit"},
		{tool: "update_deal_stage", permission: "pipeline.update_stage"},
		{tool: "update_product_status", permission: "products.edit"},
		{tool: "update_monthly_target", permission: "monthly-targets.edit"},
	}

	for _, tc := range cases {
		t.Run(tc.tool, func(t *testing.T) {
			if _, ok := aiToolPermissions[tc.tool]; !ok {
				t.Fatalf("expected %s to be registered in aiToolPermissions", tc.tool)
			}
			if !service.canRunToolCall(&ToolCall{Tool: tc.tool, Params: map[string]interface{}{}}, testUserContext([]string{tc.permission}, nil)) {
				t.Fatalf("expected %s permission to allow %s", tc.permission, tc.tool)
			}
			if service.canRunToolCall(&ToolCall{Tool: tc.tool, Params: map[string]interface{}{}}, testUserContext(nil, nil)) {
				t.Fatalf("expected %s to be blocked without permission", tc.tool)
			}
		})
	}
}

func TestComprehensiveCRUDAIToolsAreRegistered(t *testing.T) {
	service := &Service{}
	expected := map[string]string{
		"create_user":                      "users.create",
		"update_user":                      "users.edit",
		"delete_user":                      "users.delete",
		"create_group":                     "groups.create",
		"update_group":                     "groups.edit",
		"delete_group":                     "groups.delete",
		"update_ai_settings":               "ai-settings.edit",
		"update_profile":                   "profile.edit",
		"change_password":                  "profile.change-password",
		"mark_notification_read":           "notifications.mark-read",
		"delete_notification":              "notifications.delete",
		"create_monthly_target":            "monthly-targets.create",
		"update_monthly_target":            "monthly-targets.edit",
		"delete_monthly_target":            "monthly-targets.delete",
		"distribute_brick_targets":         "bricks.distribute-targets",
		"update_brick_target_distribution": "bricks.target-distributions-edit",
		"delete_brick_target_distribution": "bricks.target-distributions-delete",
		"create_lead":                      "leads.create",
		"update_lead":                      "leads.edit",
		"delete_lead":                      "leads.delete",
		"convert_lead":                     "leads.convert",
		"create_lead_status":               "leads.status-create",
		"update_lead_status_meta":          "leads.status-edit",
		"delete_lead_status":               "leads.status-delete",
		"set_default_lead_status":          "leads.status-default",
		"create_lead_industry":             "leads.industries-create",
		"update_lead_industry":             "leads.industries-edit",
		"delete_lead_industry":             "leads.industries-delete",
		"create_lead_source":               "leads.sources-create",
		"update_lead_source":               "leads.sources-edit",
		"delete_lead_source":               "leads.sources-delete",
		"create_deal":                      "pipeline.create",
		"update_deal":                      "pipeline.edit",
		"delete_deal":                      "pipeline.delete",
		"move_deal":                        "pipeline.move",
		"convert_quotation":                "pipeline.convert_quotation",
		"convert_sales_order":              "pipeline.convert_sales_order",
		"create_pipeline_stage":            "pipeline.stages-create",
		"update_pipeline_stage":            "pipeline.stages-edit",
		"delete_pipeline_stage":            "pipeline.stages-delete",
		"reorder_pipeline_stage":           "pipeline.stages-order",
		"create_task":                      "tasks.create",
		"update_task":                      "tasks.edit",
		"delete_task":                      "tasks.delete",
		"create_task_lead":                 "tasks.create_lead",
		"create_schedule":                  "schedules.create",
		"update_schedule":                  "schedules.edit",
		"delete_schedule":                  "schedules.delete",
		"assign_schedule":                  "schedules.assign",
		"create_visit_report":              "visit-reports.create",
		"update_visit_report":              "visit-reports.edit",
		"delete_visit_report":              "visit-reports.delete",
		"approve_visit_report":             "visit-reports.approve",
		"update_activity_type":             "visit-reports.activity-type",
		"create_route":                     "route-optimization.create",
		"delete_route":                     "route-optimization.delete",
		"create_product":                   "products.create",
		"update_product":                   "products.edit",
		"delete_product":                   "products.delete",
		"create_product_category":          "products.category-create",
		"update_product_category":          "products.category-edit",
		"delete_product_category":          "products.category-delete",
		"create_account":                   "accounts.create",
		"update_account":                   "accounts.edit",
		"delete_account":                   "accounts.delete",
		"create_account_category":          "accounts.category-create",
		"update_account_category":          "accounts.category-edit",
		"delete_account_category":          "accounts.category-delete",
		"create_contact_role":              "accounts.role-create",
		"update_contact_role":              "accounts.role-edit",
		"delete_contact_role":              "accounts.role-delete",
		"generate_report":                  "reports.generate",
		"create_brick":                     "bricks.create",
		"update_brick":                     "bricks.edit",
		"delete_brick":                     "bricks.delete",
		"create_territory":                 "area-mapping.territories-create",
		"update_territory":                 "area-mapping.territories-edit",
		"delete_territory":                 "area-mapping.territories-delete",
		"create_area_capture":              "area-mapping.captures-create",
	}

	for tool, permission := range expected {
		t.Run(tool, func(t *testing.T) {
			if !isRegisteredAITool(tool) {
				t.Fatalf("expected %s to be registered", tool)
			}
			if !service.canRunToolCall(&ToolCall{Tool: tool, Params: map[string]interface{}{}}, testUserContext([]string{permission}, nil)) {
				t.Fatalf("expected %s to allow %s", permission, tool)
			}
			if service.canRunToolCall(&ToolCall{Tool: tool, Params: map[string]interface{}{}}, testUserContext(nil, nil)) {
				t.Fatalf("expected %s to be blocked without permission", tool)
			}
		})
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

func TestUserDataPermissionAllowsChatbotScopedUserAccess(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"ai-chatbot.view"}, map[string]roledomain.ScopeType{
		"users": roledomain.ScopeOwn,
	})

	if !service.hasDataPermission("user", userCtx) {
		t.Fatal("expected chatbot user to access users data within RBAC scope")
	}
}

func TestRoleDataPermissionAllowsChatbotAccess(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"ai-chatbot.view"}, nil)

	if !service.hasDataPermission("role", userCtx) {
		t.Fatal("expected chatbot user to access basic role data")
	}
}

func TestRoleManagementSearchTermIgnoresGenericRoleRequest(t *testing.T) {
	if got := roleManagementSearchTerm("berikan data roles"); got != "" {
		t.Fatalf("expected generic roles request to have no search filter, got %q", got)
	}
	if got := roleManagementSearchTerm("berikan data role sales"); got != "sales" {
		t.Fatalf("expected role search term sales, got %q", got)
	}
}

func TestGroupAndBrickDataPermissionAllowsChatbotAccess(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"ai-chatbot.view"}, nil)

	if !service.hasDataPermission("group", userCtx) {
		t.Fatal("expected chatbot user to access basic group data")
	}
	if !service.hasDataPermission("brick_management", userCtx) {
		t.Fatal("expected chatbot user to access brick data within RBAC scope")
	}
}

func TestAdminUsesGlobalAIScopeEvenWithoutResourceScope(t *testing.T) {
	service := &Service{}
	userCtx := testUserContext([]string{"ai-chatbot.view"}, map[string]roledomain.ScopeType{
		"users": roledomain.ScopeGlobal,
	})
	userCtx.RoleCode = "admin"

	if scoped := service.scopedUserIDs(userCtx, "user"); scoped != nil {
		t.Fatalf("expected admin user data scope to be global, got %#v", scoped)
	}
	if scoped := service.scopedUserIDs(userCtx, "brick_management"); scoped != nil {
		t.Fatalf("expected admin brick data scope to be global, got %#v", scoped)
	}
}

func TestManagementSearchTermsIgnoreGenericRequests(t *testing.T) {
	if got := groupManagementSearchTerm("berikan data groups"); got != "" {
		t.Fatalf("expected generic groups request to have no search filter, got %q", got)
	}
	if got := brickManagementSearchTerm("berikan data area"); got != "" {
		t.Fatalf("expected generic area request to have no search filter, got %q", got)
	}
	if got := userManagementSearchTerm("berikan dara users"); got != "" {
		t.Fatalf("expected typo users request to have no search filter, got %q", got)
	}
	if got := userManagementSearchTerm("jadi berika data data users yang ada"); got != "" {
		t.Fatalf("expected generic users request with filler words to have no search filter, got %q", got)
	}
	if got := brickManagementSearchTerm("tampilkan data bricks yang ada"); got != "" {
		t.Fatalf("expected generic bricks request with filler words to have no search filter, got %q", got)
	}
	if got := roleManagementSearchTerm("berikan data roles yang ada"); got != "" {
		t.Fatalf("expected generic roles request with filler words to have no search filter, got %q", got)
	}
	if got := brickManagementSearchTerm("tampilkan data bricks Semarang dan siapa manager serta penghasilan bricks / area tersebut berapa"); got != "semarang" {
		t.Fatalf("expected brick search term semarang, got %q", got)
	}
}

func TestUserRoleFilterTermFromNaturalLanguage(t *testing.T) {
	if got := userRoleFilterTerm("berikan data users yang memiliki role sales"); got != "sales" {
		t.Fatalf("expected sales role filter, got %q", got)
	}
	if got := userRoleFilterTerm("tampilkan users dengan role sales manager aktif"); got != "sales_manager" {
		t.Fatalf("expected sales_manager role filter, got %q", got)
	}
	if got := userRoleFilterTerm("berikan data users yang ada"); got != "" {
		t.Fatalf("expected no role filter, got %q", got)
	}
}

func TestResolveRoleIDForUserFilter(t *testing.T) {
	service := &Service{
		roleRepo: stubAIRoleRepo{roles: []roledomain.Role{
			{ID: "role-sales", Name: "Sales Representative", Code: "sales"},
			{ID: "role-manager", Name: "Sales Manager", Code: "sales_manager"},
		}},
	}

	roleID, ok := service.resolveRoleIDForUserFilter("sales")
	if !ok || roleID != "role-sales" {
		t.Fatalf("expected sales role id, got id=%q ok=%v", roleID, ok)
	}
	roleID, ok = service.resolveRoleIDForUserFilter("sales_manager")
	if !ok || roleID != "role-manager" {
		t.Fatalf("expected sales manager role id, got id=%q ok=%v", roleID, ok)
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

	filter = targetOwnerFilterFromMessage("update targets bulan ini untuk sales menjadi 10 juta")
	if filter != "" {
		t.Fatalf("expected empty owner filter for generic sales update, got %q", filter)
	}

	filter = targetOwnerFilterFromMessage("update targets bulan juli untuk sales menjadi 10 juta")
	if filter != "" {
		t.Fatalf("expected empty owner filter for generic sales July update, got %q", filter)
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

func TestProspectPredictionIntentHandlesPotentialDealWording(t *testing.T) {
	message := "prediksi prospect yang paling berpotensi deal bulan ini"
	if !isProspectPredictionIntent(message) {
		t.Fatal("expected prospect potential deal wording to be detected as prediction intent")
	}
}

func TestVisitRecommendationIntentHandlesPotentialLeadAccountWording(t *testing.T) {
	message := "berikan rekomendasi kunjungan berdasarkan lead/account yang paling potensial"
	if !isVisitRecommendationPlannerIntent(message) {
		t.Fatal("expected visit recommendation wording to be detected")
	}
}

func TestProposalDraftIntentDoesNotRequireCRMDeal(t *testing.T) {
	message := "buatkan saya draft proposal dengan nilai dan timeline yang telah diprediksi"
	if !isProposalDraftIntent(message) {
		t.Fatal("expected proposal draft request to be handled as content generation")
	}

	saveMessage := "buat deal dan simpan proposal ini sebagai opportunity"
	if isProposalDraftIntent(saveMessage) {
		t.Fatal("expected explicit save/create deal request to stay available for CRM write flow")
	}
}

func TestScheduleCRUDIntentHandlesRescheduleWording(t *testing.T) {
	message := "ubah jadwal meeting langsung menjadi tanggal 24 juli"
	if !isScheduleCRUDIntent(message) {
		t.Fatal("expected schedule change wording to be detected as schedule CRUD intent")
	}
}

func TestDealValueTargetIntentDoesNotTriggerMonthlyTarget(t *testing.T) {
	message := "propect rs kariadi, memiliki stage desire, dengan target deal 50 juta"
	if !isDealValueTargetIntent(message) {
		t.Fatal("expected deal target value wording to be treated as deal value, not monthly target data")
	}

	targetMessage := "tampilkan target sales bulan ini"
	if isDealValueTargetIntent(targetMessage) {
		t.Fatal("expected monthly target wording to remain target data intent")
	}
}

func TestPendingDealAccountConfirmationParsesReplyAndContext(t *testing.T) {
	history := []aidomain.ChatMessage{
		{Role: "user", Content: "Propect RS Kariadi, memiliki stage desire, dengan target deal 50 juta"},
		{Role: "assistant", Content: "Saya menemukan beberapa account yang mirip dengan **RS Kariadi**. Mohon konfirmasi account yang dimaksud:\n\n1. **RSUP Dr Kariadi** — Semarang, Jawa Tengah\n\nBalas dengan nama account yang benar, misalnya: **RSUP Dr Kariadi**."},
	}

	if !isAccountConfirmationReply("RSUP Dr Kariadi.") {
		t.Fatal("expected RSUP Dr Kariadi reply to be treated as account confirmation")
	}
	if !lastAssistantAskedAccountConfirmation(history) {
		t.Fatal("expected last assistant message to be recognized as account confirmation prompt")
	}
	pending := latestPendingCreateDealUserMessage(history)
	if pending == "" {
		t.Fatal("expected pending deal user message to be found")
	}
	if got := extractDealStageName(pending); got != "Desire" {
		t.Fatalf("expected Desire stage, got %q", got)
	}
	if got := extractDealValueText(pending); got != "50 juta" {
		t.Fatalf("expected 50 juta value, got %q", got)
	}
}

func TestExtractNamesFromProspectHospitalWording(t *testing.T) {
	names := extractNamesFromHistory("propect rs kariadi, memiliki stage desire, dengan target deal 50 juta")
	joined := strings.Join(names, "|")
	if !strings.Contains(joined, "rs kariadi") && !strings.Contains(joined, "kariadi") {
		t.Fatalf("expected account name hint for RS Kariadi, got %#v", names)
	}
}

func TestAccountSearchTermCandidatesIncludesDistinctiveHospitalName(t *testing.T) {
	candidates := accountSearchTermCandidates([]string{"rsup dr kariadi"})
	joined := strings.Join(candidates, "|")
	if !strings.Contains(joined, "kariadi") {
		t.Fatalf("expected account search candidates to include distinctive token, got %#v", candidates)
	}
	if strings.Contains(joined, "|rsup|") || strings.Contains(joined, "|dr|") {
		t.Fatalf("expected account search candidates to skip noisy hospital prefixes, got %#v", candidates)
	}
}

func TestAccountMatchDistinguishesSimilarFromExactHospitalName(t *testing.T) {
	accountEntity := account.Account{Name: "RSUP Dr Kariadi"}
	if scoreAccountMatch(accountEntity, []string{"rs kariadi"}) == 0 {
		t.Fatal("expected RS Kariadi to produce a similar account match for RSUP Dr Kariadi")
	}
	if isExactAccountMatch(accountEntity, []string{"rs kariadi"}) {
		t.Fatal("expected RS Kariadi not to be treated as an exact match for RSUP Dr Kariadi")
	}
	if !isExactAccountMatch(accountEntity, []string{"rsup dr kariadi"}) {
		t.Fatal("expected RSUP Dr Kariadi to be treated as an exact account match")
	}
}

func TestHasPositiveAccountMatchFindsContainedHospitalToken(t *testing.T) {
	accounts := map[string]account.Account{
		"account-1": {Name: "RSUP Dr Kariadi"},
	}
	if !hasPositiveAccountMatch(accounts, []string{"RS Kariadi"}) {
		t.Fatal("expected account containing Kariadi to be treated as a positive match")
	}
}

func TestBuildAccountOutOfScopeMessageExplainsOwnScope(t *testing.T) {
	message := buildAccountOutOfScopeMessage([]string{"RSUP Dr Kariadi"}, &domainauth.UserContext{
		RoleCode: "custom_field_role",
		Scopes: map[string]roledomain.ScopeType{
			"accounts": roledomain.ScopeOwn,
		},
	})
	if !strings.Contains(message, "ditemukan di database") {
		t.Fatalf("expected message to distinguish out-of-scope from not-found, got %q", message)
	}
	if !strings.Contains(message, "Scope accounts user login adalah own") {
		t.Fatalf("expected message to explain own scope, got %q", message)
	}
}

func TestCanAccessAccountForToolAllowsOwnScopeSameBrick(t *testing.T) {
	brickID := "brick-semarang"
	otherUserID := "other-sales"
	service := &Service{
		userRepo: stubAIUserRepo{
			user: &userdomain.User{ID: "sales-semarang", BrickID: &brickID},
		},
	}
	userCtx := &domainauth.UserContext{
		UserID:   "sales-semarang",
		RoleCode: "custom_field_role",
		Scopes: map[string]roledomain.ScopeType{
			"accounts": roledomain.ScopeOwn,
		},
	}
	accountEntity := account.Account{
		Name:       "RSUP Dr Kariadi",
		AssignedTo: &otherUserID,
		BrickID:    &brickID,
	}

	if !service.canAccessAccountForTool(userCtx, accountEntity) {
		t.Fatal("expected own-scoped user to access account in the same brick")
	}
}

func TestNormalizeProductStatusHandlesIndonesianStatus(t *testing.T) {
	status, ok := normalizeProductStatus("nonaktif")
	if !ok {
		t.Fatal("expected nonaktif to be parsed")
	}
	if status != "inactive" {
		t.Fatalf("expected nonaktif to map to inactive, got %q", status)
	}
}

func TestProcessToolCallsReturnsOnlyFailureMessage(t *testing.T) {
	service := &Service{}
	response := "Berhasil memperbarui status produk.\n\nInsight singkat\n- Produk inactive.\n<!-- TOOL_CALL:{\"tool\":\"update_product_status\",\"params\":{\"product_name\":\"Thermometer Digital\",\"status\":\"inactive\"}} -->"

	finalMessage := service.processToolCalls(response, "user-1", nil, testUserContext(nil, nil))
	if strings.Contains(finalMessage, "Berhasil memperbarui") || strings.Contains(finalMessage, "Insight singkat") {
		t.Fatalf("expected failed tool response to suppress optimistic LLM text, got %q", finalMessage)
	}
	if !strings.Contains(finalMessage, "Gagal menjalankan aksi") {
		t.Fatalf("expected failure message to be shown, got %q", finalMessage)
	}
}

func TestBuildToolResultBlockIncludesDealDetailAction(t *testing.T) {
	dealID := "733722e7-c4a5-46ec-b9b6-010f85833f94"
	message := buildToolResultBlock(toolResult{
		Success:      true,
		Entity:       "Deal",
		ID:           dealID,
		Message:      "**Penawaran RSUP Dr Kariadi - Desire**",
		DetailEntity: "deal",
		DetailID:     dealID,
		DetailLabel:  "Lihat Deal",
	})

	service := &Service{}
	entityIDs := service.extractEntityIDsFromHistory([]aidomain.ChatMessage{
		{Role: "assistant", Content: message},
	})
	if len(entityIDs["deal"]) != 1 || entityIDs["deal"][0] != dealID {
		t.Fatalf("expected deal id to be extracted from detail action, got %#v", entityIDs["deal"])
	}
}

func TestExtractProductNamesFromStageUpdate(t *testing.T) {
	products := extractProductNamesFromStageUpdate("ubah stages nya ke closed won, produk Cetirizine 10mg Tablet 100, Vitamin C 500mg 200")
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %#v", products)
	}
	if products[0] != "Cetirizine 10mg Tablet 100" || products[1] != "Vitamin C 500mg 200" {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestScoreProductMatchFindsThermometerToken(t *testing.T) {
	productEntity := productdomain.Product{Name: "Thermometer Digital", SKU: "THERMO-001"}
	if scoreProductMatch(productEntity, []string{"thermometer"}) == 0 {
		t.Fatal("expected thermometer token to match Thermometer Digital")
	}
}

func TestDealValueToolParamsConvertRupiahToSen(t *testing.T) {
	valueSen, ok := dealValueSenFromParams(map[string]interface{}{
		"value": float64(50000000),
	})
	if !ok {
		t.Fatal("expected numeric deal value to be parsed")
	}
	if valueSen != 5000000000 {
		t.Fatalf("expected Rp 50.000.000 to be stored as 5.000.000.000 sen, got %d", valueSen)
	}
}

func TestDealValueToolParamsParseJutaTextToSen(t *testing.T) {
	valueSen, ok := dealValueSenFromParams(map[string]interface{}{
		"value": "50 juta",
	})
	if !ok {
		t.Fatal("expected text deal value to be parsed")
	}
	if valueSen != 5000000000 {
		t.Fatalf("expected 50 juta to be stored as 5.000.000.000 sen, got %d", valueSen)
	}
}

func TestDealValueToolParamsParseFormattedRupiahToSen(t *testing.T) {
	valueSen, ok := dealValueSenFromParams(map[string]interface{}{
		"value": "Rp 50.000.000",
	})
	if !ok {
		t.Fatal("expected formatted Rupiah deal value to be parsed")
	}
	if valueSen != 5000000000 {
		t.Fatalf("expected Rp 50.000.000 to be stored as 5.000.000.000 sen, got %d", valueSen)
	}
}

func TestIsDateOnlyText(t *testing.T) {
	if !isDateOnlyText("2026-07-24") {
		t.Fatal("expected ISO date without time to be date-only")
	}
	if isDateOnlyText("2026-07-24T10:00:00+07:00") {
		t.Fatal("expected RFC3339 datetime not to be date-only")
	}
}

func TestBuildLeadPredictionItemScoresQualifiedBANTLead(t *testing.T) {
	now := time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC)
	closeDate := now.AddDate(0, 0, 14)
	item := buildLeadPredictionItem(lead.Lead{
		ID:                 "lead-1",
		FirstName:          "Sari",
		CompanyName:        "RS Sehat",
		LeadStatus:         "qualified",
		Probability:        40,
		EstimatedValue:     500000000,
		BudgetConfirmed:    true,
		AuthorityConfirmed: true,
		NeedConfirmed:      true,
		TimelineConfirmed:  true,
		ExpectedCloseDate:  &closeDate,
		UpdatedAt:          now,
	}, now)

	if item.Score < 80 {
		t.Fatalf("expected qualified BANT lead score >= 80, got %d", item.Score)
	}
	if len(item.ScoreBreakdown) == 0 || !strings.Contains(item.ScoreBreakdown[len(item.ScoreBreakdown)-1], "final=") {
		t.Fatalf("expected lead score breakdown with final score, got %#v", item.ScoreBreakdown)
	}
	if item.NextBestAction == "" {
		t.Fatal("expected next best action to be populated")
	}
}

func TestBuildDealPredictionItemScoresNegotiationDeal(t *testing.T) {
	now := time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC)
	closeDate := now.AddDate(0, 0, 7)
	item := buildDealPredictionItem(pipeline.Deal{
		ID:                 "deal-1",
		Title:              "Formulary Expansion",
		Status:             "open",
		Probability:        50,
		Value:              1200000000,
		BudgetConfirmed:    true,
		AuthorityConfirmed: true,
		NeedConfirmed:      true,
		ExpectedCloseDate:  &closeDate,
		UpdatedAt:          now,
		Stage: &pipeline.PipelineStage{
			Name:        "Negotiation",
			Code:        "negotiation",
			Probability: 70,
		},
	}, now)

	if item.Score < 80 {
		t.Fatalf("expected negotiation deal score >= 80, got %d", item.Score)
	}
	if len(item.ScoreBreakdown) == 0 || !strings.Contains(item.ScoreBreakdown[len(item.ScoreBreakdown)-1], "final=") {
		t.Fatalf("expected deal score breakdown with final score, got %#v", item.ScoreBreakdown)
	}
	if !strings.Contains(strings.ToLower(item.NextBestAction), "negosiasi") {
		t.Fatalf("expected negotiation next best action, got %q", item.NextBestAction)
	}
}
