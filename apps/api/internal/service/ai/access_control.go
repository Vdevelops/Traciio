package ai

import (
	"fmt"
	"sort"
	"strings"

	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
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
	"create_task":             {"tasks.create"},
	"create_lead":             {"leads.create"},
	"create_activity":         {"leads.edit", "visit-reports.create"},
	"create_product_interest": {"leads.edit"},
	"create_visit_report":     {"visit-reports.create"},
	"upsert_lead_bant":        {"leads.edit"},
	"create_deal":             {"pipeline.create"},
	"create_schedule":         {"schedules.create"},
	"update_schedule":         {"schedules.edit"},
	"create_route":            {"route-optimization.create"},
	"update_task_status":      {"tasks.edit", "tasks.complete", "tasks.start", "tasks.cancel"},
	"update_lead_status":      {"leads.edit"},
	"update_deal_stage":       {"pipeline.update_stage", "pipeline.move"},
	"update_product_status":   {"products.edit"},
}

func (s *Service) ensureUserContext(userID string, userCtx *domainauth.UserContext) *domainauth.UserContext {
	if userCtx != nil {
		return userCtx
	}

	permMap := map[string]bool{}
	if s.permService != nil {
		if perms, err := s.permService.GetUserPermissions(userID); err == nil {
			for _, p := range perms {
				permMap[p] = true
			}
		}
	}

	return &domainauth.UserContext{
		UserID:      userID,
		Permissions: permMap,
	}
}

func (s *Service) hasAnyPermission(userCtx *domainauth.UserContext, permissions ...string) bool {
	if userCtx == nil {
		return false
	}
	if userCtx.RoleCode == "super_admin" {
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
		return userCtx != nil && userCtx.RoleCode == "super_admin"
	}
	return s.hasAnyPermission(userCtx, rule.Permissions...)
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

func isAIGlobalRole(userCtx *domainauth.UserContext) bool {
	if userCtx == nil {
		return false
	}
	return userCtx.RoleCode == "admin" || userCtx.RoleCode == "super_admin"
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
