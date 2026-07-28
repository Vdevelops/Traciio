package ai

import (
	"fmt"
	"sort"
	"strings"

	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

type aiDataAccessRule struct {
	DataType    string
	Resource    string
	Permissions []string
}

var aiDataAccessRules = map[string]aiDataAccessRule{
	"visit_report":       {DataType: "visit_report", Resource: "visit-reports", Permissions: []string{"visit-reports.view"}},
	"account":            {DataType: "account", Resource: "accounts", Permissions: []string{"accounts.view"}},
	"contact":            {DataType: "contact", Resource: "contacts", Permissions: []string{"contacts.view", "accounts.view"}},
	"deal":               {DataType: "deal", Resource: "deals", Permissions: []string{"pipeline.view"}},
	"lead":               {DataType: "lead", Resource: "leads", Permissions: []string{"leads.view"}},
	"activity":           {DataType: "activity", Resource: "activities", Permissions: []string{"activities.view", "visit-reports.view"}},
	"task":               {DataType: "task", Resource: "tasks", Permissions: []string{"tasks.view"}},
	"product":            {DataType: "product", Resource: "products", Permissions: []string{"products.view"}},
	"pipeline":           {DataType: "pipeline", Resource: "deals", Permissions: []string{"pipeline.view"}},
	"schedule":           {DataType: "schedule", Resource: "schedules", Permissions: []string{"schedules.view"}},
	"sales_performance":  {DataType: "sales_performance", Resource: "sales-overview", Permissions: []string{"sales-overview.view", "dashboard.view"}},
	"product_analysis":   {DataType: "product_analysis", Resource: "sales-overview", Permissions: []string{"product-analytics.view"}},
	"report":             {DataType: "report", Resource: "reports", Permissions: []string{"reports.view"}},
	"user":               {DataType: "user", Resource: "users", Permissions: []string{"users.view", "ai-chatbot.view"}},
	"role":               {DataType: "role", Resource: "users", Permissions: []string{"users.roles", "users.permissions", "ai-chatbot.view"}},
	"group":              {DataType: "group", Resource: "groups", Permissions: []string{"groups.view", "ai-chatbot.view"}},
	"brick_management":   {DataType: "brick_management", Resource: "bricks", Permissions: []string{"bricks.view", "area-mapping.view", "area-mapping.territories-view", "ai-chatbot.view"}},
	"target":             {DataType: "target", Resource: "monthly-targets", Permissions: []string{"monthly-targets.view"}},
	"route_optimization": {DataType: "route_optimization", Resource: "route-optimization", Permissions: []string{"route-optimization.view"}},
}

var aiToolPermissions = map[string][]string{
	// Users, groups, settings, profile, and notifications.
	"create_user":            {"users.create"},
	"update_user":            {"users.edit"},
	"delete_user":            {"users.delete"},
	"view_roles":             {"users.roles"},
	"view_permissions":       {"users.permissions"},
	"create_group":           {"groups.create"},
	"update_group":           {"groups.edit"},
	"delete_group":           {"groups.delete"},
	"update_ai_settings":     {"ai-settings.edit"},
	"update_profile":         {"profile.edit"},
	"change_password":        {"profile.change-password"},
	"mark_notification_read": {"notifications.mark-read"},
	"delete_notification":    {"notifications.delete"},

	// Monthly targets and brick target distributions.
	"create_monthly_target":            {"monthly-targets.create"},
	"update_monthly_target":            {"monthly-targets.edit"},
	"delete_monthly_target":            {"monthly-targets.delete"},
	"distribute_brick_targets":         {"bricks.distribute-targets"},
	"update_brick_target_distribution": {"bricks.target-distributions-edit"},
	"delete_brick_target_distribution": {"bricks.target-distributions-delete"},

	// Sales CRM.
	"create_lead":             {"leads.create"},
	"update_lead":             {"leads.edit"},
	"delete_lead":             {"leads.delete"},
	"convert_lead":            {"leads.convert"},
	"update_lead_status":      {"leads.edit"},
	"create_lead_status":      {"leads.status-create"},
	"update_lead_status_meta": {"leads.status-edit"},
	"delete_lead_status":      {"leads.status-delete"},
	"set_default_lead_status": {"leads.status-default"},
	"create_lead_industry":    {"leads.industries-create"},
	"update_lead_industry":    {"leads.industries-edit"},
	"delete_lead_industry":    {"leads.industries-delete"},
	"create_lead_source":      {"leads.sources-create"},
	"update_lead_source":      {"leads.sources-edit"},
	"delete_lead_source":      {"leads.sources-delete"},
	"create_product_interest": {"leads.edit"},
	"upsert_lead_bant":        {"leads.edit"},

	"create_deal":            {"pipeline.create"},
	"update_deal":            {"pipeline.edit"},
	"delete_deal":            {"pipeline.delete"},
	"move_deal":              {"pipeline.move"},
	"update_deal_stage":      {"pipeline.update_stage", "pipeline.move"},
	"convert_quotation":      {"pipeline.convert_quotation"},
	"convert_sales_order":    {"pipeline.convert_sales_order"},
	"create_pipeline_stage":  {"pipeline.stages-create"},
	"update_pipeline_stage":  {"pipeline.stages-edit"},
	"delete_pipeline_stage":  {"pipeline.stages-delete"},
	"reorder_pipeline_stage": {"pipeline.stages-order"},

	"create_task":        {"tasks.create"},
	"update_task":        {"tasks.edit"},
	"delete_task":        {"tasks.delete"},
	"update_task_status": {"tasks.edit", "tasks.complete", "tasks.start", "tasks.cancel"},
	"create_task_lead":   {"tasks.create_lead"},

	"create_schedule": {"schedules.create"},
	"update_schedule": {"schedules.edit"},
	"delete_schedule": {"schedules.delete"},
	"assign_schedule": {"schedules.assign"},

	"create_visit_report":  {"visit-reports.create"},
	"update_visit_report":  {"visit-reports.edit"},
	"delete_visit_report":  {"visit-reports.delete"},
	"approve_visit_report": {"visit-reports.approve"},
	"create_activity":      {"leads.edit", "visit-reports.create"},
	"update_activity_type": {"visit-reports.activity-type"},

	// Route optimization.
	"create_route": {"route-optimization.create"},
	"delete_route": {"route-optimization.delete"},

	// Inventory.
	"create_product":          {"products.create"},
	"update_product":          {"products.edit"},
	"delete_product":          {"products.delete"},
	"update_product_status":   {"products.edit"},
	"create_product_category": {"products.category-create"},
	"update_product_category": {"products.category-edit"},
	"delete_product_category": {"products.category-delete"},

	// Customers.
	"create_account":          {"accounts.create"},
	"update_account":          {"accounts.edit"},
	"delete_account":          {"accounts.delete"},
	"create_account_category": {"accounts.category-create"},
	"update_account_category": {"accounts.category-edit"},
	"delete_account_category": {"accounts.category-delete"},
	"create_contact_role":     {"accounts.role-create"},
	"update_contact_role":     {"accounts.role-edit"},
	"delete_contact_role":     {"accounts.role-delete"},

	// Analytics, reports, area mapping, and bricks.
	"generate_report":     {"reports.generate"},
	"create_brick":        {"bricks.create"},
	"update_brick":        {"bricks.edit"},
	"delete_brick":        {"bricks.delete"},
	"create_territory":    {"area-mapping.territories-create"},
	"update_territory":    {"area-mapping.territories-edit"},
	"delete_territory":    {"area-mapping.territories-delete"},
	"create_area_capture": {"area-mapping.captures-create"},
}

func (s *Service) ensureUserContext(userID string, userCtx *domainauth.UserContext) *domainauth.UserContext {
	if userCtx != nil {
		return userCtx
	}

	permMap := map[string]bool{}
	resolved := &domainauth.UserContext{
		UserID:      userID,
		Permissions: permMap,
		Scopes:      map[string]roledomain.ScopeType{},
	}

	if s.userRepo != nil {
		if userEntity, err := s.userRepo.FindByID(userID); err == nil && userEntity != nil {
			resolved.Email = userEntity.Email
			resolved.RoleID = userEntity.RoleID
			if userEntity.GroupID != nil {
				resolved.GroupID = *userEntity.GroupID
			}
			if userEntity.Role != nil {
				resolved.RoleCode = userEntity.Role.Code
				for _, permission := range userEntity.Role.Permissions {
					if permission.Code != "" {
						permMap[permission.Code] = true
					} else if permission.Resource != "" && permission.Action != "" {
						permMap[fmt.Sprintf("%s.%s", permission.Resource, permission.Action)] = true
					}
				}
			}
		}
	}

	if s.permService != nil {
		var perms []string
		var err error
		if resolved.RoleCode != "" {
			perms, err = s.permService.GetPermissionsByRole(resolved.RoleCode)
		} else {
			perms, err = s.permService.GetUserPermissions(userID)
		}
		if err == nil {
			for _, p := range perms {
				permMap[p] = true
			}
		}
	}

	if s.roleRepo != nil && resolved.RoleID != "" {
		if scopes, err := s.roleRepo.GetScopesByRoleID(resolved.RoleID); err == nil {
			for _, scope := range scopes {
				resolved.Scopes[scope.Resource] = scope.Scope
			}
		}
	}
	resolved.TeamMemberIDs = s.resolveAITeamMemberIDs(userID)
	return resolved
}

func (s *Service) resolveAITeamMemberIDs(userID string) []string {
	seen := map[string]struct{}{userID: {}}
	if s.brickRepo == nil {
		return []string{userID}
	}

	managerID := userID
	if s.userRepo != nil {
		if currentUser, err := s.userRepo.FindByID(userID); err == nil && currentUser != nil && currentUser.BrickID != nil && *currentUser.BrickID != "" {
			if salesReps, salesErr := s.brickRepo.GetSalesByBrickID(*currentUser.BrickID); salesErr == nil {
				for _, rep := range salesReps {
					seen[rep.ID] = struct{}{}
				}
			}
		}
	}

	bricks, _, err := s.brickRepo.List(&brickdomain.ListBricksRequest{ManagerID: &managerID, Page: 1, PerPage: 100})
	if err != nil {
		return []string{userID}
	}

	for _, brick := range bricks {
		salesReps, err := s.brickRepo.GetSalesByBrickID(brick.ID)
		if err != nil {
			continue
		}
		for _, rep := range salesReps {
			seen[rep.ID] = struct{}{}
		}
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func (s *Service) hasAnyPermission(userCtx *domainauth.UserContext, permissions ...string) bool {
	if userCtx == nil {
		return false
	}
	if isAIGlobalRole(userCtx) {
		return true
	}
	for _, p := range permissions {
		if p != "" && userCtx.HasPermission(p) {
			return true
		}
	}
	return false
}

func (s *Service) hasDataPermission(dataType string, userCtx *domainauth.UserContext) bool {
	rule, ok := aiDataAccessRules[dataType]
	if !ok {
		return isAIGlobalRole(userCtx)
	}
	return s.hasAnyPermission(userCtx, rule.Permissions...)
}

func isAIGlobalRole(userCtx *domainauth.UserContext) bool {
	if userCtx == nil {
		return false
	}
	roleCode := strings.ToLower(strings.TrimSpace(userCtx.RoleCode))
	return roleCode == "admin" || roleCode == "super_admin"
}

func (s *Service) scopedUserIDs(userCtx *domainauth.UserContext, dataType string) []string {
	if userCtx == nil {
		return nil
	}
	if isAIGlobalRole(userCtx) {
		return nil
	}
	rule, ok := aiDataAccessRules[dataType]
	if !ok {
		return []string{userCtx.UserID}
	}
	return userCtx.GetScopedUserIDs(rule.Resource)
}

func (s *Service) canAccessOwner(userCtx *domainauth.UserContext, dataType string, ownerID string) bool {
	if userCtx == nil {
		return false
	}
	scopedIDs := s.scopedUserIDs(userCtx, dataType)
	if scopedIDs == nil {
		return true
	}
	if ownerID == "" {
		return false
	}
	for _, id := range scopedIDs {
		if id == ownerID {
			return true
		}
	}
	return false
}

func (s *Service) canRunTool(tool string, userCtx *domainauth.UserContext) bool {
	required, ok := aiToolPermissions[tool]
	if !ok {
		return false
	}
	return s.hasAnyPermission(userCtx, required...)
}

func isRegisteredAITool(tool string) bool {
	_, ok := aiToolPermissions[tool]
	return ok
}

func (s *Service) canRunToolCall(call *ToolCall, userCtx *domainauth.UserContext) bool {
	if call == nil {
		return false
	}
	required, ok := aiToolPermissions[call.Tool]
	if !ok {
		return false
	}

	switch call.Tool {
	case "update_task_status":
		status := strings.ToLower(strings.TrimSpace(paramStr(call.Params, "status")))
		switch status {
		case "completed", "complete", "done":
			required = []string{"tasks.complete"}
		case "in_progress", "in-progress", "in progress", "start", "started":
			required = []string{"tasks.start"}
		case "cancelled", "canceled", "cancel":
			required = []string{"tasks.cancel"}
		default:
			required = []string{"tasks.edit"}
		}
	case "update_lead_status":
		if s.isConvertedLeadStatusCall(call.Params) {
			required = []string{"leads.convert"}
		} else {
			required = []string{"leads.edit"}
		}
	case "update_deal_stage":
		required = []string{"pipeline.update_stage", "pipeline.move"}
	}

	return s.hasAnyPermission(userCtx, required...)
}

func (s *Service) isConvertedLeadStatusCall(params map[string]interface{}) bool {
	for _, key := range []string{"lead_status_code", "lead_status", "status"} {
		value := strings.ToLower(strings.TrimSpace(paramStr(params, key)))
		if value == "converted" || value == "convert" || value == "won" {
			return true
		}
	}
	if s.leadStatusRepo == nil {
		return false
	}
	leadStatusID := strings.TrimSpace(paramStr(params, "lead_status_id"))
	if leadStatusID == "" {
		return false
	}
	status, err := s.leadStatusRepo.FindByID(leadStatusID)
	return err == nil && status != nil && status.IsConverted
}

func (s *Service) buildAIAccessContext(userCtx *domainauth.UserContext) string {
	if userCtx == nil {
		return "\n\n=== USER ACCESS CONTEXT ===\n- AI permission context unavailable. Deny data/action access unless explicit context is provided by backend.\n"
	}

	allowedData := make([]string, 0)
	for dataType := range aiDataAccessRules {
		if s.hasDataPermission(dataType, userCtx) {
			rule := aiDataAccessRules[dataType]
			scope := userCtx.GetScope(rule.Resource)
			allowedData = append(allowedData, fmt.Sprintf("%s(scope=%s)", dataType, scope))
		}
	}
	sort.Strings(allowedData)

	allowedTools := make([]string, 0)
	for tool := range aiToolPermissions {
		if s.canRunTool(tool, userCtx) {
			allowedTools = append(allowedTools, tool)
		}
	}
	sort.Strings(allowedTools)

	if len(allowedData) == 0 {
		allowedData = append(allowedData, "none")
	}
	if len(allowedTools) == 0 {
		allowedTools = append(allowedTools, "none")
	}

	return fmt.Sprintf(`

=== USER ACCESS CONTEXT ===
- User ID: %s
- Role: %s
- Allowed data modules: %s
- Allowed action tools: %s

ACCESS RULES:
- Use only data modules listed above and only data already provided in backend context.
- If a user asks for data outside the allowed modules, say they do not have permission.
- If a user asks for an action outside the allowed tools, do not emit TOOL_CALL.
- Respect scope labels: own means only the user's records, team means team records, global means all permitted records.
`, userCtx.UserID, userCtx.RoleCode, strings.Join(allowedData, ", "), strings.Join(allowedTools, ", "))
}
