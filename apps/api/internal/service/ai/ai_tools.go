package ai

// AI Tools/Skills layer
// Handles CRUD tool execution when the LLM emits <!-- TOOL_CALL:{...} --> markers.
//
// Flow:
//  1. LLM response is received from Cerebras
//  2. processToolCalls() scans for <!-- TOOL_CALL:{...} --> markers
//  3. Each marker is executed against the real CRM service
//  4. Marker is replaced with a markdown confirmation + action card

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	accountdomain "github.com/gilabs/crm-healthcare/api/internal/domain/account"
	activitydomain "github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	aidomain "github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	contactdomain "github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	leaddomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	leadqualificationdomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	pipelinedomain "github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	productdomain "github.com/gilabs/crm-healthcare/api/internal/domain/product"
	route_optimization_domain "github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	scheduledomain "github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	taskdomain "github.com/gilabs/crm-healthcare/api/internal/domain/task"
	visitreportdomain "github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

// toolCallPattern matches <!-- TOOL_CALL:{...} --> markers in LLM responses.
// (?s) enables dotall mode so '.' also matches newlines — the LLM often emits pretty-printed JSON.
var toolCallPattern = regexp.MustCompile(`(?s)<!--\s*TOOL_CALL:(.*?)\s*-->`)
var aiActionPattern = regexp.MustCompile(`(?s)<!--\s*ACTION\s*:\s*.*?\s*-->`)

// ToolCall is the parsed representation of an LLM tool invocation.
type ToolCall struct {
	Tool   string                 `json:"tool"`
	Params map[string]interface{} `json:"params"`
}

// toolResult holds the outcome of a single tool execution.
type toolResult struct {
	Success      bool
	Confirm      bool   // True when the action needs user confirmation instead of being treated as a failure
	Entity       string // Human-readable entity name (e.g. "Task", "Lead")
	ID           string // ID of the created/updated entity
	Message      string // Short confirmation detail shown under the success header
	Action       string // Verb shown in success header (e.g. "dibuat", "diperbarui")
	PageURL      string // CRM page URL for the action card
	Icon         string // Lucide icon name for the action card
	DetailEntity string // Entity type for opening a detail drawer
	DetailID     string // ID used by the detail drawer
	DetailLabel  string // Label for the detail action card
}

type toolConfirmationError struct {
	Message string
}

func (e *toolConfirmationError) Error() string {
	return e.Message
}

// processToolCalls scans an LLM response for TOOL_CALL markers, executes each
// tool synchronously, and replaces every marker with its result block.
func (s *Service) processToolCalls(response string, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) string {
	matches := toolCallPattern.FindAllStringSubmatchIndex(response, -1)
	if len(matches) == 0 {
		return response
	}

	response = aiActionPattern.ReplaceAllString(response, "")
	matches = toolCallPattern.FindAllStringSubmatchIndex(response, -1)

	// Iterate in reverse so earlier indices stay valid after replacements.
	for i := len(matches) - 1; i >= 0; i-- {
		fullStart, fullEnd := matches[i][0], matches[i][1]
		rawJSON := cleanToolCallJSON(response[matches[i][2]:matches[i][3]])

		var call ToolCall
		if err := json.Unmarshal([]byte(rawJSON), &call); err != nil || call.Tool == "" {
			// Silently remove the broken marker — don't surface raw parse errors to the user.
			response = response[:fullStart] + response[fullEnd:]
			continue
		}

		result := s.executeTool(&call, userID, history, userCtx)
		log.Printf("[AI_TOOL_AUDIT] user_id=%s tool=%s success=%t entity=%s id=%s", userID, call.Tool, result.Success, result.Entity, result.ID)
		if !result.Success && !result.Confirm {
			return strings.TrimSpace(buildToolResultBlock(result))
		}
		response = response[:fullStart] + buildToolResultBlock(result) + response[fullEnd:]
	}
	return response
}

// cleanToolCallJSON normalises the raw text captured between TOOL_CALL markers.
// The LLM (especially small models) frequently wraps JSON in markdown code fences
// or adds preamble text. This function:
//  1. Strips leading/trailing whitespace
//  2. Removes ```json / ``` code fences
//  3. Extracts the first balanced { ... } object it finds
func cleanToolCallJSON(raw string) string {
	raw = strings.TrimSpace(raw)

	// Strip markdown code fences
	for _, fence := range []string{"```json", "```JSON", "```"} {
		if strings.HasPrefix(raw, fence) {
			raw = raw[len(fence):]
			break
		}
	}
	if strings.HasSuffix(raw, "```") {
		raw = raw[:len(raw)-3]
	}
	raw = strings.TrimSpace(raw)

	// Extract the first balanced JSON object
	start := strings.Index(raw, "{")
	if start < 0 {
		return raw
	}
	depth := 0
	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return raw[start : i+1]
			}
		}
	}
	return raw[start:]
}

// executeTool dispatches to the correct handler by tool name.
func (s *Service) executeTool(call *ToolCall, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if call == nil {
		return toolResult{Success: false, Entity: "AI Action", Message: "Tool call tidak valid."}
	}
	userCtx = s.ensureUserContext(userID, userCtx)
	if !isRegisteredAITool(call.Tool) {
		return toolResult{
			Success: false,
			Entity:  "AI Action",
			Message: fmt.Sprintf("Tool '%s' belum didukung oleh AI Assistant.", call.Tool),
		}
	}
	if !s.canRunToolCall(call, userCtx) {
		return toolResult{
			Success: false,
			Entity:  "AI Action",
			Message: fmt.Sprintf("Anda tidak memiliki permission untuk menjalankan tool '%s'.", call.Tool),
		}
	}

	switch call.Tool {
	case "create_task":
		return s.toolCreateTask(call.Params, userID, history, userCtx)
	case "create_lead":
		return s.toolCreateLead(call.Params, userID)
	case "create_activity":
		return s.toolCreateActivity(call.Params, userID, history, userCtx)
	case "create_product_interest":
		return s.toolCreateProductInterest(call.Params, userID, history, userCtx)
	case "create_visit_report":
		return s.toolCreateVisitReport(call.Params, userID, history, userCtx)
	case "upsert_lead_bant":
		if isProductInterestOnlyParams(call.Params) {
			return s.toolCreateProductInterest(call.Params, userID, history, userCtx)
		}
		return s.toolUpsertLeadBANT(call.Params, history, userCtx)
	case "create_deal":
		return s.toolCreateDeal(call.Params, userID, history, userCtx)
	case "create_schedule":
		return s.toolCreateSchedule(call.Params, userID, userCtx)
	case "update_schedule":
		return s.toolUpdateSchedule(call.Params, history, userCtx)
	case "create_route":
		return s.toolCreateRoute(call.Params, userID, history, userCtx)
	case "update_task_status":
		return s.toolUpdateTaskStatus(call.Params, history, userCtx)
	case "update_lead_status":
		return s.toolUpdateLeadStatus(call.Params, history, userCtx)
	case "update_deal_stage":
		return s.toolUpdateDealStage(call.Params, userID, history, userCtx)
	case "update_product_status":
		return s.toolUpdateProductStatus(call.Params)
	case "update_monthly_target":
		return s.toolUpdateMonthlyTarget(call.Params, userCtx)
	default:
		return toolResult{Success: false, Entity: "Tool", Message: fmt.Sprintf("Tool '%s' tidak dikenali.", call.Tool)}
	}
}

// buildToolResultBlock renders a toolResult into the markdown block that
// replaces the TOOL_CALL marker in the final response.
func buildToolResultBlock(r toolResult) string {
	if !r.Success {
		if r.Confirm {
			return "\n\n" + r.Message
		}
		if r.Entity == "AI Action" {
			return fmt.Sprintf("\n\n⚠️ Gagal menjalankan aksi: %s", r.Message)
		}
		return fmt.Sprintf("\n\n⚠️ Gagal membuat %s: %s", r.Entity, r.Message)
	}

	var b strings.Builder
	action := r.Action
	if action == "" {
		action = "dibuat"
	}
	b.WriteString(fmt.Sprintf("\n\n✅ **%s berhasil %s!**", r.Entity, action))
	if r.Message != "" {
		b.WriteString("\n")
		b.WriteString(r.Message)
	}
	if r.PageURL != "" {
		if action, err := json.Marshal(map[string]string{
			"type":        "navigate",
			"label":       fmt.Sprintf("Lihat %s", r.Entity),
			"description": fmt.Sprintf("Buka halaman %s", r.Entity),
			"url":         r.PageURL,
			"icon":        r.Icon,
		}); err == nil {
			b.WriteString(fmt.Sprintf("\n<!-- ACTION:%s -->", action))
		}
	}
	if r.DetailEntity != "" && r.DetailID != "" {
		label := r.DetailLabel
		if label == "" {
			label = fmt.Sprintf("Lihat %s", r.Entity)
		}
		if action, err := json.Marshal(map[string]string{
			"type":        "detail",
			"label":       label,
			"description": fmt.Sprintf("Buka detail %s", r.DetailEntity),
			"entity":      r.DetailEntity,
			"entityId":    r.DetailID,
			"icon":        r.Icon,
		}); err == nil {
			b.WriteString(fmt.Sprintf("\n<!-- ACTION:%s -->", action))
		}
	}
	return b.String()
}

// ----------------------------------------------------------------------------
// Individual Tool Handlers
// ----------------------------------------------------------------------------

func (s *Service) toolCreateTask(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.taskService == nil {
		return toolResult{Success: false, Entity: "Task", Message: "Task service tidak tersedia."}
	}

	title := paramStr(params, "title")
	if title == "" {
		return toolResult{Success: false, Entity: "Task", Message: "Judul task wajib diisi."}
	}

	customerSource := strings.ToLower(strings.TrimSpace(firstNonEmptyParam(params, "customer_source", "source", "customer_type", "entity_source")))
	if customerSource == "" {
		return toolResult{Success: false, Entity: "Task", Message: "Sumber customer wajib diisi. Pilih lead atau account, lalu sebutkan nama perusahaan atau nama contact."}
	}
	if customerSource != "lead" && customerSource != "account" {
		return toolResult{Success: false, Entity: "Task", Message: "Sumber customer tidak valid. Gunakan lead atau account."}
	}

	leadID := ""
	accountID := ""
	contactID := paramStr(params, "contact_id")
	if customerSource == "lead" {
		resolvedLeadID, leadEntity, err := s.resolveLeadForTool(params, history, userCtx)
		if err != nil {
			return toolResult{Success: false, Entity: "Task", Message: err.Error()}
		}
		leadOwner := ""
		if leadEntity.AssignedTo != nil {
			leadOwner = *leadEntity.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "lead", leadOwner) {
			return toolResult{Success: false, Entity: "Task", Message: "Anda tidak memiliki akses ke lead yang dipilih."}
		}
		leadID = resolvedLeadID
	} else {
		resolvedAccountID, resolvedContactID, err := s.resolveAccountTaskContextForTool(params, history, userCtx)
		if err != nil {
			return toolResult{Success: false, Entity: "Task", Message: err.Error()}
		}
		accountID = resolvedAccountID
		if contactID == "" {
			contactID = resolvedContactID
		}
	}

	if err := s.validateTaskLinkedEntityAccess(params, userCtx); err != nil {
		return toolResult{Success: false, Entity: "Task", Message: err.Error()}
	}

	req := &taskdomain.CreateTaskRequest{
		Title:       title,
		Description: paramStr(params, "description"),
		Priority:    paramStrOr(params, "priority", "medium"),
		Type:        paramStrOr(params, "type", "general"),
		AssignedTo:  userID,
		AccountID:   accountID,
		ContactID:   contactID,
		DealID:      paramStr(params, "deal_id"),
		LeadID:      leadID,
	}
	if ds := paramStr(params, "due_date"); ds != "" {
		if t, err := parseFlexibleTime(ds); err == nil {
			req.DueDate = &t
		}
	}

	resp, err := s.taskService.CreateTask(req, userID)
	if err != nil {
		return toolResult{Success: false, Entity: "Task", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Task", ID: resp.ID,
		Message: fmt.Sprintf("**%s** (prioritas: %s)", resp.Title, resp.Priority),
		PageURL: "/tasks", Icon: "clipboard",
	}
}

func (s *Service) toolCreateLead(params map[string]interface{}, userID string) toolResult {
	if s.leadService == nil {
		return toolResult{Success: false, Entity: "Lead", Message: "Lead service tidak tersedia."}
	}

	firstName := paramStr(params, "first_name")
	if firstName == "" {
		firstName = paramStr(params, "name")
	}
	if firstName == "" {
		return toolResult{Success: false, Entity: "Lead", Message: "Nama lead wajib diisi."}
	}
	email := paramStr(params, "email")
	if email == "" {
		return toolResult{Success: false, Entity: "Lead", Message: "Email lead wajib diisi agar lead bisa dibuat."}
	}

	req := &leaddomain.CreateLeadRequest{
		FirstName:   firstName,
		LastName:    paramStr(params, "last_name"),
		Email:       email,
		Phone:       paramStr(params, "phone"),
		CompanyName: paramStr(params, "company_name"),
		JobTitle:    paramStr(params, "job_title"),
		LeadSource:  paramStrOr(params, "lead_source", "other"),
		Notes:       paramStr(params, "notes"),
	}

	resp, err := s.leadService.Create(req, userID, nil)
	if err != nil {
		return toolResult{Success: false, Entity: "Lead", Message: err.Error()}
	}

	name := resp.FirstName
	if resp.LastName != "" {
		name += " " + resp.LastName
	}
	return toolResult{
		Success: true, Entity: "Lead", ID: resp.ID,
		Message: fmt.Sprintf("**%s** dari %s", name, resp.CompanyName),
		PageURL: "/leads", Icon: "user",
	}
}

func (s *Service) toolCreateActivity(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.activityService == nil {
		return toolResult{Success: false, Entity: "Activity", Message: "Activity service tidak tersedia."}
	}

	description := firstNonEmptyParam(params, "description", "notes", "note", "title")
	if description == "" {
		description = "Activity added by AI"
	}

	activityType := normalizeActivityType(paramStrOr(params, "type", "call"))
	timestamp := time.Now()
	if timestampText := paramStr(params, "timestamp"); timestampText != "" {
		if parsed, err := parseFlexibleTime(timestampText); err == nil {
			timestamp = parsed
		}
	}

	var leadID *string
	leadName := ""
	if s.shouldResolveLeadForActivity(params, history) {
		resolvedID, leadEntity, err := s.resolveLeadForTool(params, history, userCtx)
		if err != nil {
			return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
		}
		leadOwner := ""
		if leadEntity.AssignedTo != nil {
			leadOwner = *leadEntity.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "lead", leadOwner) {
			return toolResult{Success: false, Entity: "Activity", Message: "Anda tidak memiliki akses ke lead yang dipilih."}
		}
		leadID = toolStringPtr(resolvedID)
		leadName = strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName))
		if leadName == "" {
			leadName = leadEntity.CompanyName
		}
	}

	accountID := paramStr(params, "account_id")
	contactID := paramStr(params, "contact_id")
	if accountID == "" && hasAccountReferenceForActivity(params) {
		resolvedAccountID, resolvedContactID, err := s.resolveAccountTaskContextForTool(params, history, userCtx)
		if err != nil {
			if confirmationErr, ok := err.(*toolConfirmationError); ok {
				return toolResult{Success: false, Confirm: true, Entity: "Activity", Message: confirmationErr.Message}
			}
			return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
		}
		accountID = resolvedAccountID
		if contactID == "" {
			contactID = resolvedContactID
		}
	}
	if accountID != "" {
		if err := s.validateAccountAccess(accountID, userCtx); err != nil {
			return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
		}
	}
	if contactID != "" {
		if err := s.validateContactAccess(contactID, userCtx); err != nil {
			return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
		}
	}

	dealID := paramStr(params, "deal_id")
	if dealID != "" {
		dealEntity, err := s.dealRepo.FindByID(dealID)
		if err != nil || dealEntity == nil {
			return toolResult{Success: false, Entity: "Activity", Message: "Deal yang ditautkan tidak ditemukan atau tidak dapat diakses."}
		}
		dealOwner := ""
		if dealEntity.AssignedTo != nil {
			dealOwner = *dealEntity.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "deal", dealOwner) {
			return toolResult{Success: false, Entity: "Activity", Message: "Anda tidak memiliki akses ke deal yang ditautkan."}
		}
		if accountID == "" {
			accountID = dealEntity.AccountID
		}
	} else if inferredDealID, inferredAccountID := s.resolveActivityDealContext(params, history, userCtx, accountID); inferredDealID != "" {
		dealID = inferredDealID
		if accountID == "" {
			accountID = inferredAccountID
		}
	}

	if leadID == nil && accountID == "" && dealID == "" {
		return toolResult{Success: false, Entity: "Activity", Message: "Activity wajib ditautkan ke lead, account, atau deal yang bisa diakses."}
	}

	metadata := map[string]interface{}{
		"source": "ai_chatbot",
	}
	if productInterests, err := s.productInterestMetadataFromParams(params); err != nil {
		return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
	} else if len(productInterests) > 0 {
		metadata["product_interests"] = productInterests
	}

	req := &activitydomain.CreateActivityRequest{
		Type:        activityType,
		LeadID:      leadID,
		UserID:      userID,
		Description: description,
		Timestamp:   timestamp.Format(time.RFC3339),
		Metadata:    metadata,
	}
	if accountID != "" {
		req.AccountID = toolStringPtr(accountID)
	}
	if contactID != "" {
		req.ContactID = toolStringPtr(contactID)
	}
	if dealID != "" {
		req.DealID = toolStringPtr(dealID)
	}

	resp, err := s.activityService.Create(req)
	if err != nil {
		return toolResult{Success: false, Entity: "Activity", Message: err.Error()}
	}

	message := fmt.Sprintf("**%s** (%s)", resp.Description, resp.Type)
	if leadName != "" {
		message = fmt.Sprintf("%s untuk **%s**", message, leadName)
	}
	result := toolResult{
		Success: true, Entity: "Activity", ID: resp.ID,
		Message: message,
		PageURL: "/leads", Icon: "clipboard",
	}
	if leadID != nil && *leadID != "" {
		result.DetailEntity = "lead"
		result.DetailID = *leadID
		result.DetailLabel = "Lihat Lead"
	} else if dealID != "" {
		result.PageURL = "/pipeline"
		result.DetailEntity = "deal"
		result.DetailID = dealID
		result.DetailLabel = "Lihat Deal"
	}
	return result
}

func (s *Service) toolCreateProductInterest(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.activityService == nil {
		return toolResult{Success: false, Entity: "Product Interest", Message: "Activity service tidak tersedia."}
	}

	productInterests, err := s.productInterestMetadataFromParams(params)
	if err != nil {
		return toolResult{Success: false, Entity: "Product Interest", Message: err.Error()}
	}
	if len(productInterests) == 0 {
		return toolResult{Success: false, Entity: "Product Interest", Message: "Nama produk wajib diisi untuk menambahkan product interest."}
	}

	resolvedLeadID, leadEntity, err := s.resolveLeadForTool(params, history, userCtx)
	if err != nil {
		return toolResult{Success: false, Entity: "Product Interest", Message: err.Error()}
	}
	leadOwner := ""
	if leadEntity.AssignedTo != nil {
		leadOwner = *leadEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "lead", leadOwner) {
		return toolResult{Success: false, Entity: "Product Interest", Message: "Anda tidak memiliki akses ke lead yang dipilih."}
	}

	productNames := make([]string, 0, len(productInterests))
	for _, item := range productInterests {
		if name, ok := item["product_name"].(string); ok && name != "" {
			productNames = append(productNames, name)
		}
	}

	leadName := strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName))
	if leadName == "" {
		leadName = leadEntity.CompanyName
	}

	description := firstNonEmptyParam(params, "description", "notes", "note")
	if description == "" {
		description = "Product interest: " + strings.Join(productNames, ", ")
	}

	req := &activitydomain.CreateActivityRequest{
		Type:        "task",
		LeadID:      toolStringPtr(resolvedLeadID),
		UserID:      userID,
		Description: description,
		Timestamp:   time.Now().Format(time.RFC3339),
		Metadata: map[string]interface{}{
			"source":            "ai_chatbot",
			"activity_category": "product_interest",
			"product_interests": productInterests,
		},
	}

	resp, err := s.activityService.Create(req)
	if err != nil {
		return toolResult{Success: false, Entity: "Product Interest", Message: err.Error()}
	}

	return toolResult{
		Success: true, Entity: "Product Interest", ID: resp.ID,
		Message: fmt.Sprintf("Product interest **%s** berhasil ditambahkan ke lead **%s**.", strings.Join(productNames, ", "), leadName),
		PageURL: "/leads", Icon: "package",
		DetailEntity: "lead", DetailID: resolvedLeadID, DetailLabel: "Lihat Lead",
	}
}

func (s *Service) toolCreateVisitReport(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.visitReportService == nil {
		return toolResult{Success: false, Entity: "Visit Report", Message: "Visit report service tidak tersedia."}
	}

	purpose := firstNonEmptyParam(params, "purpose", "title", "description", "notes")
	if purpose == "" {
		return toolResult{Success: false, Entity: "Visit Report", Message: "Purpose/tujuan visit wajib diisi."}
	}

	visitDate := time.Now()
	if dateText := firstNonEmptyParam(params, "visit_date", "scheduled_at", "date"); dateText != "" {
		parsed, err := parseFlexibleTime(dateText)
		if err != nil {
			return toolResult{Success: false, Entity: "Visit Report", Message: "Tanggal visit tidak valid."}
		}
		visitDate = parsed
	}

	resolvedLeadID, leadEntity, err := s.resolveLeadForTool(params, history, userCtx)
	if err != nil {
		return toolResult{Success: false, Entity: "Visit Report", Message: err.Error()}
	}
	leadOwner := ""
	if leadEntity.AssignedTo != nil {
		leadOwner = *leadEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "lead", leadOwner) {
		return toolResult{Success: false, Entity: "Visit Report", Message: "Anda tidak memiliki akses ke lead yang dipilih."}
	}

	metadata := map[string]interface{}{
		"source": "ai_chatbot",
	}
	if productInterests, err := s.productInterestMetadataFromParams(params); err != nil {
		return toolResult{Success: false, Entity: "Visit Report", Message: err.Error()}
	} else if len(productInterests) > 0 {
		metadata["product_interests"] = productInterests
	}

	req := &visitreportdomain.CreateVisitReportRequest{
		LeadID:     toolStringPtr(resolvedLeadID),
		SalesRepID: userID,
		VisitDate:  visitDate.Format("2006-01-02 15:04"),
		Purpose:    purpose,
		Notes:      paramStr(params, "notes"),
		Metadata:   metadata,
	}

	resp, err := s.visitReportService.Create(req)
	if err != nil {
		return toolResult{Success: false, Entity: "Visit Report", Message: err.Error()}
	}

	leadName := strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName))
	if leadName == "" {
		leadName = leadEntity.CompanyName
	}
	return toolResult{
		Success: true, Entity: "Visit Report", ID: resp.ID,
		Message: fmt.Sprintf("**%s** untuk lead **%s** pada %s", resp.Purpose, leadName, resp.VisitDate.Format("02 Jan 2006")),
		PageURL: "/visit-reports", Icon: "calendar",
		DetailEntity: "lead", DetailID: resolvedLeadID, DetailLabel: "Lihat Lead",
	}
}

func (s *Service) toolUpsertLeadBANT(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.leadQualificationService == nil {
		return toolResult{Success: false, Entity: "BANT", Message: "Lead qualification service tidak tersedia."}
	}

	resolvedLeadID, leadEntity, err := s.resolveLeadForTool(params, history, userCtx)
	if err != nil {
		return toolResult{Success: false, Entity: "BANT", Message: err.Error()}
	}
	leadOwner := ""
	if leadEntity.AssignedTo != nil {
		leadOwner = *leadEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "lead", leadOwner) {
		return toolResult{Success: false, Entity: "BANT", Message: "Anda tidak memiliki akses ke lead yang dipilih."}
	}

	req := &leadqualificationdomain.UpsertLeadQualificationRequest{
		BudgetTargetCurrency:  paramStr(params, "budget_target_currency"),
		BudgetNotes:           paramStr(params, "budget_notes"),
		AuthorityTargetPerson: firstNonEmptyParam(params, "authority_target_person", "authority_person", "decision_maker"),
		AuthorityTargetRole:   firstNonEmptyParam(params, "authority_target_role", "authority_role", "decision_maker_role"),
		AuthorityNotes:        paramStr(params, "authority_notes"),
		NeedPriorityLevel:     normalizeNeedPriority(paramStr(params, "need_priority_level")),
		NeedNotes:             firstNonEmptyParam(params, "need_notes", "product_interest_notes"),
		TimelineFlexibility:   normalizeTimelineFlexibility(paramStr(params, "timeline_flexibility")),
		TimelineNotes:         paramStr(params, "timeline_notes"),
		NeedTargetProducts:    s.needProductsFromParams(params),
	}
	if amount := paramInt64(params, "budget_target_amount"); amount != nil {
		req.BudgetTargetAmount = amount
	}
	if value, ok := paramBoolPtr(params, "budget_confirmed"); ok {
		req.BudgetConfirmed = value
	}
	if value, ok := paramBoolPtr(params, "authority_confirmed"); ok {
		req.AuthorityConfirmed = value
	}
	if value, ok := paramBoolPtr(params, "need_confirmed"); ok {
		req.NeedConfirmed = value
	}
	if value, ok := paramBoolPtr(params, "timeline_confirmed"); ok {
		req.TimelineConfirmed = value
	}
	if dateText := paramStr(params, "timeline_target_date"); dateText != "" {
		if parsed, err := parseFlexibleTime(dateText); err == nil {
			req.TimelineTargetDate = &parsed
		}
	}

	resp, err := s.leadQualificationService.Upsert(resolvedLeadID, req)
	if err != nil {
		return toolResult{Success: false, Entity: "BANT", Message: err.Error()}
	}

	leadName := strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName))
	if leadName == "" {
		leadName = leadEntity.CompanyName
	}
	return toolResult{
		Success: true, Entity: "BANT", ID: resp.ID,
		Message: fmt.Sprintf("BANT lead **%s** diperbarui. Skor: **%d**, status: **%s**.", leadName, resp.QualificationScore, resp.QualificationStatus),
		PageURL: "/leads", Icon: "target",
		DetailEntity: "lead", DetailID: resolvedLeadID, DetailLabel: "Lihat Lead",
	}
}

func (s *Service) toolCreateDeal(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.pipelineService == nil {
		return toolResult{Success: false, Entity: "Deal", Message: "Pipeline service tidak tersedia."}
	}

	title := paramStr(params, "title")
	if title == "" {
		title = paramStr(params, "deal_name")
	}
	if title == "" {
		title = paramStr(params, "name")
	}
	if title == "" {
		return toolResult{Success: false, Entity: "Deal", Message: "Nama deal wajib diisi."}
	}

	accountID, contactID, err := s.resolveAccountTaskContextForTool(params, history, userCtx)
	if err != nil {
		if confirmationErr, ok := err.(*toolConfirmationError); ok {
			return toolResult{
				Success: false,
				Confirm: true,
				Entity:  "Deal",
				Message: confirmationErr.Message,
			}
		}
		return toolResult{
			Success: false, Entity: "Deal",
			Message: err.Error(),
		}
	}
	accountEntity, accountErr := s.accountRepo.FindByID(accountID)
	if accountErr != nil || accountEntity == nil || !s.canAccessAccountForTool(userCtx, *accountEntity) {
		return toolResult{Success: false, Entity: "Deal", Message: "Anda tidak memiliki akses ke account yang dipilih untuk membuat deal."}
	}
	if explicitContactID := paramStr(params, "contact_id"); explicitContactID != "" {
		contactID = explicitContactID
		if err := s.validateContactAccess(contactID, userCtx); err != nil {
			return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
		}
	}

	stageID := ""
	if hasPipelineStageReference(params) {
		resolvedStageID, err := s.resolvePipelineStageID(params)
		if err != nil {
			return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
		}
		stageID = resolvedStageID
	}
	if stageID == "" {
		stages, err := s.pipelineService.ListStages(&pipelinedomain.ListPipelineStagesRequest{})
		if err != nil || len(stages) == 0 {
			return toolResult{Success: false, Entity: "Deal", Message: "Tidak ada pipeline stage yang tersedia."}
		}
		best := stages[0]
		for _, st := range stages[1:] {
			if st.Order < best.Order {
				best = st
			}
		}
		stageID = best.ID
	}

	req := &pipelinedomain.CreateDealRequest{
		Title:      title,
		AccountID:  accountID,
		StageID:    stageID,
		ContactID:  contactID,
		AssignedTo: userID,
		Notes:      paramStr(params, "notes"),
	}
	if valueSen, ok := dealValueSenFromParams(params); ok {
		req.Value = valueSen
	}

	resp, err := s.pipelineService.CreateDeal(req, userID)
	if err != nil {
		return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Deal", ID: resp.ID,
		Message: fmt.Sprintf("**%s**", resp.Title),
		PageURL: "/pipeline", Icon: "trending-up",
		DetailEntity: "deal", DetailID: resp.ID, DetailLabel: "Lihat Deal",
	}
}

func dealValueSenFromParams(params map[string]interface{}) (int64, bool) {
	for _, key := range []string{"value", "target_value", "target_amount", "amount", "nilai", "nilai_target"} {
		raw, ok := params[key]
		if !ok || raw == nil {
			continue
		}
		if valueSen, ok := parseRupiahToolValueToSen(raw); ok {
			return valueSen, true
		}
	}
	return 0, false
}

func parseRupiahToolValueToSen(raw interface{}) (int64, bool) {
	switch val := raw.(type) {
	case float64:
		if val < 0 {
			return 0, false
		}
		return int64(val * 100), true
	case float32:
		if val < 0 {
			return 0, false
		}
		return int64(float64(val) * 100), true
	case int:
		if val < 0 {
			return 0, false
		}
		return int64(val) * 100, true
	case int64:
		if val < 0 {
			return 0, false
		}
		return val * 100, true
	case json.Number:
		if parsed, err := val.Float64(); err == nil && parsed >= 0 {
			return int64(parsed * 100), true
		}
	case string:
		return parseRupiahStringToSen(val)
	}
	return 0, false
}

func parseRupiahStringToSen(value string) (int64, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return 0, false
	}
	normalized = strings.ReplaceAll(normalized, "rupiah", "")
	normalized = strings.ReplaceAll(normalized, "rp", "")
	normalized = strings.ReplaceAll(normalized, "idr", "")

	multiplier := float64(1)
	switch {
	case strings.Contains(normalized, "miliar") || strings.Contains(normalized, "billion"):
		multiplier = 1000000000
	case strings.Contains(normalized, "juta") || containsToken(normalized, "jt") || strings.Contains(normalized, "million"):
		multiplier = 1000000
	case strings.Contains(normalized, "ribu") || containsToken(normalized, "rb") || containsToken(normalized, "k"):
		multiplier = 1000
	}

	numberPattern := regexp.MustCompile(`[-+]?\d[\d.,]*`)
	numberText := numberPattern.FindString(normalized)
	if numberText == "" {
		return 0, false
	}
	number, ok := parseLocalizedNumber(numberText)
	if !ok || number < 0 {
		return 0, false
	}
	return int64(number * multiplier * 100), true
}

func parseLocalizedNumber(value string) (float64, bool) {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	if cleaned == "" {
		return 0, false
	}

	if strings.Contains(cleaned, ".") && strings.Contains(cleaned, ",") {
		lastDot := strings.LastIndex(cleaned, ".")
		lastComma := strings.LastIndex(cleaned, ",")
		if lastComma > lastDot {
			cleaned = strings.ReplaceAll(cleaned, ".", "")
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
		} else {
			cleaned = strings.ReplaceAll(cleaned, ",", "")
		}
	} else if strings.Contains(cleaned, ".") || strings.Contains(cleaned, ",") {
		separator := "."
		if strings.Contains(cleaned, ",") {
			separator = ","
		}
		parts := strings.Split(cleaned, separator)
		allThousands := len(parts) > 1
		for _, part := range parts[1:] {
			if len(part) != 3 {
				allThousands = false
				break
			}
		}
		if allThousands {
			cleaned = strings.Join(parts, "")
		} else if len(parts) == 2 {
			cleaned = strings.Join(parts, ".")
		} else {
			cleaned = strings.Join(parts, "")
		}
	}

	parsed, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func containsToken(value string, token string) bool {
	pattern := regexp.MustCompile(`(^|[^a-z0-9])` + regexp.QuoteMeta(token) + `([^a-z0-9]|$)`)
	return pattern.MatchString(value)
}

func (s *Service) toolUpdateMonthlyTarget(params map[string]interface{}, userCtx *domainauth.UserContext) toolResult {
	if s.monthlyTargetService == nil {
		return toolResult{Success: false, Entity: "Target", Message: "Monthly target service tidak tersedia."}
	}

	targetAmount, ok := monthlyTargetAmountFromParams(params)
	if !ok {
		return toolResult{
			Success: false,
			Confirm: true,
			Entity:  "Target",
			Message: "Nilai target baru belum jelas. Sebutkan nominal target, misalnya **10 juta** atau **Rp 10.000.000**.",
		}
	}

	targetID := firstNonEmptyParam(params, "id", "target_id", "monthly_target_id")
	if targetID != "" {
		target, err := s.monthlyTargetService.GetByID(targetID)
		if err != nil {
			return toolResult{Success: false, Entity: "Target", Message: "Target tidak ditemukan."}
		}
		if !canAccessMonthlyTarget(userCtx, target) {
			return toolResult{Success: false, Entity: "Target", Message: "Target berada di luar scope akses Anda."}
		}
		return s.updateMonthlyTargetByID(targetID, targetAmount)
	}

	now := time.Now()
	year := paramInt(params, "year")
	if year == 0 {
		year = now.Year()
	}
	month := paramInt(params, "month")
	if month == 0 {
		month = int(now.Month())
	}
	scope := strings.ToLower(strings.TrimSpace(paramStrOr(params, "scope", "user")))
	if scope == "" || scope == "all" {
		scope = "user"
	}
	ownerName := normalizeMonthlyTargetOwnerParam(firstNonEmptyParam(params, "owner_name", "owner", "user_name", "sales_name", "name"))
	updateAll := paramBool(params, "update_all") || paramBool(params, "all")

	req := &monthlytargetdomain.ListMonthlyTargetsRequest{
		Page:          1,
		PerPage:       100,
		Year:          &year,
		Month:         &month,
		Scope:         scope,
		ScopedUserIDs: scopedUserIDsFromContext(userCtx, "monthly-targets"),
	}
	targets, _, err := s.monthlyTargetService.List(req)
	if err != nil {
		return toolResult{Success: false, Entity: "Target", Message: "Data target tidak dapat diakses saat ini."}
	}

	matches := make([]monthlytargetdomain.MonthlyTargetResponse, 0, len(targets))
	for _, target := range targets {
		if !canAccessMonthlyTarget(userCtx, &target) {
			continue
		}
		if ownerName != "" && !matchesMonthlyTargetOwner(target, scope, ownerName) {
			continue
		}
		matches = append(matches, target)
	}

	if len(matches) == 0 {
		return toolResult{Success: false, Entity: "Target", Message: "Tidak ada target yang cocok untuk periode dan scope yang diminta."}
	}
	if len(matches) > 1 && !updateAll {
		owners := make([]string, 0, len(matches))
		for _, target := range matches {
			owners = append(owners, targetOwnerName(target, scope))
		}
		sort.Strings(owners)
		return toolResult{
			Success: false,
			Confirm: true,
			Entity:  "Target",
			Message: fmt.Sprintf("Ada beberapa target yang cocok untuk %s %d: **%s**. Sebutkan owner yang ingin diubah, atau gunakan kata **semua target** untuk mengubah semuanya.", time.Month(month).String(), year, strings.Join(owners, ", ")),
		}
	}

	updatedOwners := make([]string, 0, len(matches))
	for _, target := range matches {
		if result := s.updateMonthlyTargetByID(target.ID, targetAmount); !result.Success {
			return result
		}
		updatedOwners = append(updatedOwners, targetOwnerName(target, scope))
	}
	sort.Strings(updatedOwners)

	return toolResult{
		Success: true,
		Entity:  "Target",
		Action:  "diperbarui",
		Message: fmt.Sprintf("Target **%s %d** untuk **%s** menjadi **%s**.", time.Month(month).String(), year, strings.Join(updatedOwners, ", "), formatCurrencyRupiah(targetAmount)),
		PageURL: "/master-data/monthly-targets",
		Icon:    "target",
	}
}

func (s *Service) updateMonthlyTargetByID(targetID string, targetAmount int64) toolResult {
	updated, err := s.monthlyTargetService.Update(targetID, &monthlytargetdomain.UpdateMonthlyTargetRequest{
		TargetAmount: &targetAmount,
	})
	if err != nil {
		return toolResult{Success: false, Entity: "Target", Message: err.Error()}
	}
	return toolResult{
		Success: true,
		Entity:  "Target",
		ID:      updated.ID,
		Action:  "diperbarui",
		Message: fmt.Sprintf("Target **%s %d** menjadi **%s**.", time.Month(updated.Month).String(), updated.Year, formatCurrencyRupiah(updated.TargetAmount)),
		PageURL: "/master-data/monthly-targets",
		Icon:    "target",
	}
}

func monthlyTargetAmountFromParams(params map[string]interface{}) (int64, bool) {
	for _, key := range []string{"target_amount", "amount", "nilai_target", "target", "value"} {
		raw, ok := params[key]
		if !ok || raw == nil {
			continue
		}
		if amount, ok := parseRupiahToolValueToSen(raw); ok {
			return amount, true
		}
	}
	return 0, false
}

func canAccessMonthlyTarget(userCtx *domainauth.UserContext, target *monthlytargetdomain.MonthlyTargetResponse) bool {
	if userCtx == nil || target == nil {
		return false
	}
	scoped := scopedUserIDsFromContext(userCtx, "monthly-targets")
	if scoped == nil {
		return true
	}
	if target.UserID == nil || *target.UserID == "" {
		return false
	}
	for _, id := range scoped {
		if id == *target.UserID {
			return true
		}
	}
	return false
}

func scopedUserIDsFromContext(userCtx *domainauth.UserContext, resource string) []string {
	if userCtx == nil {
		return nil
	}
	if isAIGlobalRole(userCtx) || userCtx.IsGlobalScope(resource) {
		return nil
	}
	return userCtx.GetScopedUserIDs(resource)
}

func matchesMonthlyTargetOwner(target monthlytargetdomain.MonthlyTargetResponse, scope string, ownerName string) bool {
	normalizedOwner := strings.ToLower(strings.TrimSpace(ownerName))
	if normalizedOwner == "" {
		return true
	}
	if normalizedOwner == "sales" || normalizedOwner == "sales rep" || normalizedOwner == "sales representative" {
		return target.User != nil
	}
	targetOwner := strings.ToLower(targetOwnerName(target, scope))
	return targetOwner == normalizedOwner || strings.Contains(targetOwner, normalizedOwner)
}

func normalizeMonthlyTargetOwnerParam(ownerName string) string {
	normalized := strings.ToLower(strings.TrimSpace(ownerName))
	normalized = regexp.MustCompile(`\s+(menjadi|jadi|ke|sebesar|dengan\s+nilai)\b.*$`).ReplaceAllString(normalized, "")
	normalized = strings.Trim(normalized, ".,;:!? ")
	if normalized == "sales" || normalized == "sales rep" || normalized == "sales reps" {
		return ""
	}
	return normalized
}

func paramBool(params map[string]interface{}, key string) bool {
	if v, ok := params[key]; ok {
		switch value := v.(type) {
		case bool:
			return value
		case string:
			normalized := strings.ToLower(strings.TrimSpace(value))
			return normalized == "true" || normalized == "yes" || normalized == "ya" || normalized == "1"
		}
	}
	return false
}

func (s *Service) toolCreateSchedule(params map[string]interface{}, userID string, userCtx *domainauth.UserContext) toolResult {
	if s.scheduleService == nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Schedule service tidak tersedia."}
	}

	title := paramStr(params, "title")
	if title == "" {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Judul jadwal wajib diisi."}
	}
	if taskID := paramStr(params, "task_id"); taskID != "" {
		if err := s.validateScheduleTaskAccess(taskID, userCtx); err != nil {
			return toolResult{Success: false, Entity: "Jadwal", Message: err.Error()}
		}
	}

	// Default to tomorrow at 09:00 WIB if no time is specified.
	scheduledAt := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour).Add(9 * time.Hour)
	if ds := paramStr(params, "scheduled_at"); ds != "" {
		if t, err := parseFlexibleTime(ds); err == nil {
			scheduledAt = t
		}
	}

	req := &scheduledomain.CreateScheduleRequest{
		Title:       title,
		Description: paramStr(params, "description"),
		ScheduledAt: scheduledAt,
		TaskID:      paramStr(params, "task_id"),
	}

	resp, err := s.scheduleService.CreateSchedule(req, userID)
	if err != nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Jadwal", ID: resp.ID,
		Message: fmt.Sprintf("**%s** pada %s", resp.Title, resp.ScheduledAt.Format("02 Jan 2006 15:04")),
		PageURL: "/schedules", Icon: "calendar",
	}
}

func (s *Service) toolUpdateSchedule(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.scheduleService == nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Schedule service tidak tersedia."}
	}

	scheduleResp, resolveErr := s.resolveScheduleForTool(params, history, userCtx)
	if resolveErr != nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: resolveErr.Error()}
	}
	if scheduleResp == nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Jadwal yang ingin diubah tidak ditemukan."}
	}
	if !s.canAccessOwner(userCtx, "schedule", scheduleResp.UserID) {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Anda tidak memiliki akses untuk mengubah jadwal tersebut."}
	}

	req := &scheduledomain.UpdateScheduleRequest{
		Title:       paramStr(params, "title"),
		Description: paramStr(params, "description"),
		Status:      normalizeScheduleStatus(paramStr(params, "status")),
	}
	if reminder := paramInt64(params, "reminder_minutes_before"); reminder != nil {
		value := int(*reminder)
		req.ReminderMinutesBefore = &value
	}
	if scheduledText := firstNonEmptyParam(params, "scheduled_at", "schedule_at", "date", "tanggal", "time", "waktu"); scheduledText != "" {
		parsed, err := parseFlexibleTime(scheduledText)
		if err != nil {
			return toolResult{Success: false, Entity: "Jadwal", Message: "Format tanggal/jam tidak dikenali. Gunakan format ISO seperti 2026-07-24T10:00:00+07:00."}
		}
		if isDateOnlyText(scheduledText) {
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), scheduleResp.ScheduledAt.Hour(), scheduleResp.ScheduledAt.Minute(), scheduleResp.ScheduledAt.Second(), scheduleResp.ScheduledAt.Nanosecond(), scheduleResp.ScheduledAt.Location())
		}
		req.ScheduledAt = &parsed
	}
	if req.Title == "" && req.Description == "" && req.Status == "" && req.ScheduledAt == nil && req.ReminderMinutesBefore == nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Tidak ada perubahan jadwal yang diberikan."}
	}

	updated, err := s.scheduleService.UpdateSchedule(scheduleResp.ID, req)
	if err != nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Jadwal", ID: updated.ID, Action: "diperbarui",
		Message: fmt.Sprintf("**%s** menjadi %s", updated.Title, updated.ScheduledAt.Format("02 Jan 2006 15:04")),
		PageURL: "/schedules", Icon: "calendar",
	}
}

func (s *Service) toolCreateRoute(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.routeOptimizationService == nil {
		return toolResult{Success: false, Entity: "Rute", Message: "Route Optimization service tidak tersedia."}
	}

	// Start location: prefer explicit params, then fall back to conversation history.
	startLat := paramFloat(params, "start_lat")
	startLng := paramFloat(params, "start_lng")
	if startLat == 0 && startLng == 0 {
		lat, lng, ok := extractLocationFromHistory(history)
		if !ok {
			return toolResult{
				Success: false, Entity: "Rute",
				Message: "Lokasi awal tidak ditemukan. Mohon bagikan lokasi Anda terlebih dahulu.",
			}
		}
		startLat, startLng = lat, lng
	}

	// Account IDs: prefer explicit params, then scan conversation history.
	accountIDs := paramStringSlice(params, "account_ids")
	if len(accountIDs) == 0 {
		accountIDs = extractAccountIDsFromHistory(history)
	}
	if len(accountIDs) == 0 {
		return toolResult{
			Success: false, Entity: "Rute",
			Message: "Tidak ada akun ditemukan untuk membuat rute. Pastikan percakapan sebelumnya berisi daftar akun.",
		}
	}

	waypoints := s.buildWaypointsFromAccountIDs(accountIDs, userCtx)
	if len(waypoints) == 0 {
		return toolResult{
			Success: false, Entity: "Rute",
			Message: "Akun-akun yang dipilih tidak memiliki data koordinat GPS. Silakan lengkapi koordinat akun di halaman Accounts.",
		}
	}

	routeName := paramStr(params, "route_name")
	if routeName == "" {
		routeName = fmt.Sprintf("Rute AI - %s", time.Now().Format("02 Jan 2006"))
	}

	req := &route_optimization_domain.OptimizeRouteRequest{
		RouteName:     &routeName,
		StartLocation: &route_optimization_domain.Location{Lat: startLat, Lng: startLng},
		Waypoints:     waypoints,
	}

	result, err := s.routeOptimizationService.Optimize(req, userID)
	if err != nil {
		return toolResult{
			Success: false, Entity: "Rute",
			Message: fmt.Sprintf("Gagal mengoptimalkan rute: %v. Silakan buat rute manual di halaman Route Optimization.", err),
		}
	}

	detail := fmt.Sprintf("**%s** — %d titik kunjungan", routeName, len(result.Waypoints))
	if result.TotalDistanceFormatted != "" {
		detail += fmt.Sprintf(", jarak total: %s", result.TotalDistanceFormatted)
	}
	if result.TotalDurationFormatted != "" {
		detail += fmt.Sprintf(", estimasi: %s", result.TotalDurationFormatted)
	}
	return toolResult{
		Success: true, Entity: "Rute", ID: result.ID,
		Message: detail,
		PageURL: "/route-optimization", Icon: "map",
	}
}

func (s *Service) toolUpdateTaskStatus(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.taskService == nil {
		return toolResult{Success: false, Entity: "Task", Message: "Task service tidak tersedia."}
	}
	status := paramStr(params, "status")
	id, taskEntity, resolveErr := s.resolveTaskForTool(params, history, userCtx)
	if resolveErr != nil {
		return toolResult{Success: false, Entity: "Task", Message: resolveErr.Error()}
	}
	if status == "" {
		return toolResult{Success: false, Entity: "Task", Message: "Status task wajib diisi."}
	}
	status = normalizeTaskStatusForTool(status)
	taskOwner := ""
	if taskEntity.AssignedTo != nil {
		taskOwner = *taskEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "task", taskOwner) {
		return toolResult{Success: false, Entity: "Task", Message: "Anda tidak memiliki akses untuk mengubah task tersebut."}
	}
	resp, err := s.taskService.UpdateTask(id, &taskdomain.UpdateTaskRequest{Status: status})
	if err != nil {
		return toolResult{Success: false, Entity: "Task", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Task", ID: resp.ID, Action: "diperbarui",
		Message: fmt.Sprintf("**%s** → status: **%s**", resp.Title, resp.Status),
		PageURL: "/tasks", Icon: "clipboard",
	}
}

func normalizeTaskStatusForTool(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "complete", "done", "selesai":
		return "completed"
	case "in_progress", "in-progress", "in progress", "start", "started":
		return "in_progress"
	case "cancelled", "canceled", "cancel", "batal":
		return "cancelled"
	case "pending", "todo", "open":
		return "pending"
	default:
		return strings.ToLower(strings.TrimSpace(status))
	}
}

func normalizeScheduleStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "complete", "done", "selesai":
		return "completed"
	case "confirmed", "confirm", "konfirmasi":
		return "confirmed"
	case "submitted", "submit":
		return "submitted"
	case "cancelled", "canceled", "cancel", "batal":
		return "cancelled"
	case "rejected", "reject", "ditolak":
		return "rejected"
	case "pending", "todo", "open":
		return "pending"
	default:
		return strings.ToLower(strings.TrimSpace(status))
	}
}

func (s *Service) toolUpdateLeadStatus(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.leadService == nil {
		return toolResult{Success: false, Entity: "Lead", Message: "Lead service tidak tersedia."}
	}
	id, leadEntity, resolveErr := s.resolveLeadForTool(params, history, userCtx)
	leadStatusID, statusErr := s.resolveLeadStatusID(params)
	if resolveErr != nil {
		return toolResult{Success: false, Entity: "Lead", Message: resolveErr.Error()}
	}
	if leadStatusID == "" || statusErr != nil {
		return toolResult{Success: false, Entity: "Lead", Message: "Status lead wajib diisi. Gunakan lead_status_id, lead_status_code, atau status seperti new/contacted/interested/qualified/proposal_sent/converted/lost."}
	}
	leadOwner := ""
	if leadEntity.AssignedTo != nil {
		leadOwner = *leadEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "lead", leadOwner) {
		return toolResult{Success: false, Entity: "Lead", Message: "Anda tidak memiliki akses untuk mengubah lead tersebut."}
	}
	resp, err := s.leadService.Update(id, &leaddomain.UpdateLeadRequest{
		LeadStatusID: leadStatusID,
		StatusReason: paramStr(params, "reason"),
	}, nil)
	if err != nil {
		return toolResult{Success: false, Entity: "Lead", Message: err.Error()}
	}
	name := resp.FirstName
	if resp.LastName != "" {
		name += " " + resp.LastName
	}
	return toolResult{
		Success: true, Entity: "Lead", ID: resp.ID, Action: "diperbarui",
		Message: fmt.Sprintf("**%s** status diperbarui.", name),
		PageURL: "/leads", Icon: "user",
	}
}

func (s *Service) toolUpdateDealStage(params map[string]interface{}, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) toolResult {
	if s.pipelineService == nil {
		return toolResult{Success: false, Entity: "Deal", Message: "Pipeline service tidak tersedia."}
	}
	id, dealEntity, resolveErr := s.resolveDealForTool(params, history, userCtx)
	stageID, stageErr := s.resolvePipelineStageID(params)
	if resolveErr != nil {
		return toolResult{Success: false, Entity: "Deal", Message: resolveErr.Error()}
	}
	if stageID == "" || stageErr != nil {
		return toolResult{Success: false, Entity: "Deal", Message: "Status/stage deal wajib diisi. Gunakan stage_id, stage_code, stage_name, atau status seperti negotiation/won/lost."}
	}
	dealOwner := ""
	if dealEntity.AssignedTo != nil {
		dealOwner = *dealEntity.AssignedTo
	}
	if !s.canAccessOwner(userCtx, "deal", dealOwner) {
		return toolResult{Success: false, Entity: "Deal", Message: "Anda tidak memiliki akses untuk memindahkan deal tersebut."}
	}
	productItems, productErr := s.dealProductItemsFromParams(params)
	if productErr != nil {
		return toolResult{Success: false, Entity: "Deal", Message: productErr.Error()}
	}
	resp, err := s.pipelineService.MoveStageWithValidation(id, stageID, userID, "Dipindahkan oleh AI", productItems)
	if err != nil {
		return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Deal", ID: resp.ID, Action: "diperbarui",
		Message: fmt.Sprintf("**%s** stage diperbarui.", resp.Title),
		PageURL: "/pipeline", Icon: "trending-up",
	}
}

func (s *Service) validateTaskLinkedEntityAccess(params map[string]interface{}, userCtx *domainauth.UserContext) error {
	if accountID := paramStr(params, "account_id"); accountID != "" {
		if err := s.validateAccountAccess(accountID, userCtx); err != nil {
			return err
		}
	}
	if contactID := paramStr(params, "contact_id"); contactID != "" {
		if err := s.validateContactAccess(contactID, userCtx); err != nil {
			return err
		}
	}
	if dealID := paramStr(params, "deal_id"); dealID != "" {
		dealEntity, err := s.dealRepo.FindByID(dealID)
		if err != nil || dealEntity == nil {
			return fmt.Errorf("deal yang ditautkan tidak ditemukan atau tidak dapat diakses")
		}
		dealOwner := ""
		if dealEntity.AssignedTo != nil {
			dealOwner = *dealEntity.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "deal", dealOwner) {
			return fmt.Errorf("anda tidak memiliki akses ke deal yang ditautkan")
		}
	}
	return nil
}

func (s *Service) validateScheduleTaskAccess(taskID string, userCtx *domainauth.UserContext) error {
	if s.taskRepo == nil {
		return fmt.Errorf("task repository tidak tersedia")
	}
	taskEntity, err := s.taskRepo.FindByID(taskID)
	if err != nil || taskEntity == nil {
		return fmt.Errorf("task untuk jadwal tidak ditemukan atau tidak dapat diakses")
	}
	if taskEntity.AssignedTo == nil || *taskEntity.AssignedTo == "" {
		return fmt.Errorf("task untuk jadwal belum memiliki assignee")
	}
	taskOwner := *taskEntity.AssignedTo
	if !s.canAccessOwner(userCtx, "task", taskOwner) {
		return fmt.Errorf("anda tidak memiliki akses ke task yang dipilih")
	}
	if !s.canAccessOwner(userCtx, "schedule", taskOwner) {
		return fmt.Errorf("anda tidak memiliki akses untuk membuat jadwal untuk user task tersebut")
	}
	return nil
}

func (s *Service) validateAccountAccess(accountID string, userCtx *domainauth.UserContext) error {
	if s.accountRepo == nil {
		return fmt.Errorf("account repository tidak tersedia")
	}
	accountEntity, err := s.accountRepo.FindByID(accountID)
	if err != nil || accountEntity == nil {
		return fmt.Errorf("account yang dipilih tidak ditemukan atau tidak dapat diakses")
	}
	if !s.canAccessAccountForTool(userCtx, *accountEntity) {
		return fmt.Errorf("anda tidak memiliki akses ke account yang dipilih")
	}
	return nil
}

func (s *Service) validateContactAccess(contactID string, userCtx *domainauth.UserContext) error {
	if s.contactRepo == nil {
		return fmt.Errorf("contact repository tidak tersedia")
	}
	contactEntity, err := s.contactRepo.FindByID(contactID)
	if err != nil || contactEntity == nil {
		return fmt.Errorf("contact yang dipilih tidak ditemukan atau tidak dapat diakses")
	}
	return s.validateAccountAccess(contactEntity.AccountID, userCtx)
}

func (s *Service) resolveLeadStatusID(params map[string]interface{}) (string, error) {
	if leadStatusID := paramStr(params, "lead_status_id"); leadStatusID != "" {
		return leadStatusID, nil
	}
	if s.leadStatusRepo == nil {
		return "", fmt.Errorf("lead status repository tidak tersedia")
	}

	candidates := []string{
		paramStr(params, "lead_status_code"),
		paramStr(params, "lead_status"),
		paramStr(params, "status"),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		normalized := strings.ToLower(strings.TrimSpace(candidate))
		if status, err := s.leadStatusRepo.FindByCode(normalized); err == nil && status != nil {
			return status.ID, nil
		}
		if statuses, err := s.leadStatusRepo.ListAll(); err == nil {
			for _, status := range statuses {
				if status == nil {
					continue
				}
				if strings.EqualFold(status.Name, candidate) || strings.EqualFold(status.Code, normalized) {
					return status.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("lead status tidak ditemukan")
}

func (s *Service) resolvePipelineStageID(params map[string]interface{}) (string, error) {
	if stageID := paramStr(params, "stage_id"); stageID != "" {
		return stageID, nil
	}
	if s.pipelineService == nil {
		return "", fmt.Errorf("pipeline service tidak tersedia")
	}

	if stageCode := paramStr(params, "stage_code"); stageCode != "" {
		if stage, err := s.pipelineRepo.FindStageByCode(strings.ToLower(stageCode)); err == nil && stage != nil {
			return stage.ID, nil
		}
	}

	statusValue := strings.ToLower(strings.TrimSpace(paramStr(params, "status")))
	switch statusValue {
	case "won", "closed_won", "closed won":
		if stage, err := s.pipelineRepo.FindStageByCode("closed_won"); err == nil && stage != nil {
			return stage.ID, nil
		}
	case "lost", "closed_lost", "closed lost":
		if stage, err := s.pipelineRepo.FindStageByCode("closed_lost"); err == nil && stage != nil {
			return stage.ID, nil
		}
	}

	stageName := strings.Trim(strings.ToLower(strings.TrimSpace(paramStr(params, "stage_name"))), ".,;:!?")
	stageNameCode := strings.ReplaceAll(stageName, " ", "_")
	stages, err := s.pipelineService.ListStages(&pipelinedomain.ListPipelineStagesRequest{})
	if err != nil {
		return "", err
	}
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Order < stages[j].Order
	})

	for _, stage := range stages {
		if stageName != "" &&
			(strings.EqualFold(stage.Name, stageName) ||
				strings.EqualFold(stage.Code, stageNameCode)) {
			return stage.ID, nil
		}
		if statusValue == "" {
			continue
		}
		if strings.EqualFold(stage.Code, strings.ReplaceAll(statusValue, " ", "_")) || strings.EqualFold(stage.Name, statusValue) {
			return stage.ID, nil
		}
	}

	return "", fmt.Errorf("pipeline stage tidak ditemukan")
}

func hasPipelineStageReference(params map[string]interface{}) bool {
	return firstNonEmptyParam(params, "stage_id", "stage_code", "stage_name", "status") != ""
}

func (s *Service) resolveAccountTaskContextForTool(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) (string, string, error) {
	if s.accountRepo == nil {
		return "", "", fmt.Errorf("account repository tidak tersedia.")
	}

	contactID := paramStr(params, "contact_id")
	if contactID != "" {
		if s.contactRepo == nil {
			return "", "", fmt.Errorf("contact repository tidak tersedia.")
		}
		contactEntity, err := s.contactRepo.FindByID(contactID)
		if err != nil || contactEntity == nil {
			return "", "", fmt.Errorf("contact yang dipilih tidak ditemukan atau tidak dapat diakses.")
		}
		if err := s.validateAccountAccess(contactEntity.AccountID, userCtx); err != nil {
			return "", "", err
		}
		return contactEntity.AccountID, contactEntity.ID, nil
	}

	for _, key := range []string{"id", "account_id"} {
		if id := paramStr(params, key); id != "" {
			accountEntity, err := s.accountRepo.FindByID(id)
			if err != nil || accountEntity == nil {
				return "", "", fmt.Errorf("account yang dipilih tidak ditemukan atau tidak dapat diakses.")
			}
			if err := s.validateAccountAccess(accountEntity.ID, userCtx); err != nil {
				return "", "", err
			}
			return accountEntity.ID, "", nil
		}
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	if candidateIDs := entityIDs["account"]; len(candidateIDs) == 1 {
		accountEntity, err := s.accountRepo.FindByID(candidateIDs[0])
		if err == nil && accountEntity != nil {
			if err := s.validateAccountAccess(accountEntity.ID, userCtx); err != nil {
				return "", "", err
			}
			return accountEntity.ID, "", nil
		}
	}

	contactTerms := collectEntityHints(params, "contact_name", "contact", "pic_name")
	if len(contactTerms) > 0 {
		contactEntity, err := s.resolveContactForAccountTask(contactTerms, userCtx)
		if err != nil {
			return "", "", err
		}
		if err := s.validateAccountAccess(contactEntity.AccountID, userCtx); err != nil {
			return "", "", err
		}
		return contactEntity.AccountID, contactEntity.ID, nil
	}

	terms := collectEntityHints(params, "account_name", "company_name", "company", "name")
	if len(terms) == 0 {
		return "", "", fmt.Errorf("Account tidak ditemukan. Sebutkan nama perusahaan/account atau nama contact.")
	}

	results, err := s.searchAccountsForTool(terms, userCtx)
	if err != nil {
		return "", "", fmt.Errorf("Account tidak dapat dicari saat ini.")
	}

	filtered := make([]accountdomain.Account, 0, len(results))
	inaccessibleMatches := 0
	for _, accountEntity := range results {
		if s.canAccessAccountForTool(userCtx, accountEntity) {
			filtered = append(filtered, accountEntity)
			continue
		}
		if scoreAccountMatch(accountEntity, terms) > 0 {
			inaccessibleMatches++
		}
	}
	if len(filtered) == 0 && inaccessibleMatches > 0 {
		return "", "", &toolConfirmationError{Message: buildAccountOutOfScopeMessage(terms, userCtx)}
	}

	bestMatches := selectBestAccountMatches(filtered, terms)
	if len(bestMatches) == 1 {
		if !isExactAccountMatch(bestMatches[0], terms) {
			return "", "", &toolConfirmationError{Message: buildAccountConfirmationMessage(terms, bestMatches)}
		}
		return bestMatches[0].ID, "", nil
	}
	if len(bestMatches) > 1 {
		return "", "", &toolConfirmationError{Message: buildAccountConfirmationMessage(terms, bestMatches)}
	}

	suggestions := topAccountSuggestions(filtered, terms, 5)
	if len(suggestions) > 0 {
		return "", "", &toolConfirmationError{Message: buildAccountConfirmationMessage(terms, suggestions)}
	}

	return "", "", &toolConfirmationError{Message: fmt.Sprintf("Saya belum menemukan account yang cocok untuk **%s** sesuai akses Anda.\n\nSebutkan nama account yang lebih spesifik atau nama contact yang terhubung.", strings.Join(terms, ", "))}
}

func (s *Service) searchAccountsForTool(terms []string, userCtx *domainauth.UserContext) ([]accountdomain.Account, error) {
	searchTerms := accountSearchTermCandidates(terms)
	resultsByID := make(map[string]accountdomain.Account)

	for _, searchTerm := range searchTerms {
		results, _, err := s.accountRepo.List(&accountdomain.ListAccountsRequest{
			Page:    1,
			PerPage: 10,
			Search:  searchTerm,
		})
		if err != nil {
			return nil, err
		}
		for _, accountEntity := range results {
			resultsByID[accountEntity.ID] = accountEntity
		}
	}

	if len(resultsByID) == 0 || !hasPositiveAccountMatch(resultsByID, terms) {
		results, _, err := s.accountRepo.List(&accountdomain.ListAccountsRequest{
			Page:    1,
			PerPage: 1000,
		})
		if err != nil {
			return nil, err
		}
		for _, accountEntity := range results {
			resultsByID[accountEntity.ID] = accountEntity
		}
	}

	results := make([]accountdomain.Account, 0, len(resultsByID))
	for _, accountEntity := range resultsByID {
		results = append(results, accountEntity)
	}
	return results, nil
}

func hasPositiveAccountMatch(accounts map[string]accountdomain.Account, terms []string) bool {
	for _, accountEntity := range accounts {
		if scoreAccountMatch(accountEntity, terms) > 0 {
			return true
		}
	}
	return false
}

func (s *Service) canAccessAccountForTool(userCtx *domainauth.UserContext, accountEntity accountdomain.Account) bool {
	if userCtx == nil {
		return false
	}
	if userCtx.IsGlobalScope("accounts") {
		return true
	}

	accountOwner := ""
	if accountEntity.AssignedTo != nil {
		accountOwner = *accountEntity.AssignedTo
	}
	if s.canAccessOwner(userCtx, "account", accountOwner) {
		return true
	}

	if accountEntity.BrickID == nil || *accountEntity.BrickID == "" {
		return false
	}

	switch {
	case userCtx.IsOwnScope("accounts"):
		if s.userRepo == nil {
			return false
		}
		currentUser, err := s.userRepo.FindByID(userCtx.UserID)
		return err == nil && currentUser != nil && currentUser.BrickID != nil && *currentUser.BrickID == *accountEntity.BrickID
	case userCtx.IsTeamScope("accounts"):
		if s.brickRepo == nil {
			return false
		}
		brickEntity, err := s.brickRepo.FindByID(*accountEntity.BrickID)
		return err == nil && brickEntity != nil && brickEntity.ManagerID != nil && *brickEntity.ManagerID == userCtx.UserID
	default:
		return false
	}
}

func accountSearchTermCandidates(terms []string) []string {
	candidates := make([]string, 0)
	for _, term := range terms {
		normalized := strings.ToLower(strings.Trim(strings.TrimSpace(term), ".,;:!?"))
		if normalized == "" {
			continue
		}
		candidates = appendUnique(candidates, normalized)
		for _, token := range strings.Fields(normalized) {
			token = strings.Trim(token, ".,;:!?")
			if len(token) < 3 || isAccountNameNoiseToken(token) {
				continue
			}
			candidates = appendUnique(candidates, token)
		}
	}
	return candidates
}

func isAccountNameNoiseToken(token string) bool {
	switch token {
	case "rs", "rsu", "rsup", "rsud", "dr", "dokter", "rumah", "sakit", "prospect", "propect", "prospek":
		return true
	default:
		return isStopWord(token)
	}
}

func isExactAccountMatch(accountEntity accountdomain.Account, terms []string) bool {
	accountName := normalizeExactAccountText(accountEntity.Name)
	for _, term := range terms {
		if normalizeExactAccountText(term) == accountName {
			return true
		}
	}
	return false
}

func buildAccountConfirmationMessage(terms []string, accounts []accountdomain.Account) string {
	var b strings.Builder
	input := strings.Join(terms, ", ")
	if input == "" {
		input = "account yang diminta"
	}
	b.WriteString(fmt.Sprintf("Saya menemukan beberapa account yang mirip dengan **%s**. Mohon konfirmasi account yang dimaksud:\n\n", input))

	limit := min(len(accounts), 5)
	for i := 0; i < limit; i++ {
		accountEntity := accounts[i]
		location := accountEntity.City
		if accountEntity.Province != "" {
			if location != "" {
				location += ", "
			}
			location += accountEntity.Province
		}
		if location == "" {
			location = "-"
		}
		b.WriteString(fmt.Sprintf("%d. **%s** — %s\n", i+1, accountEntity.Name, location))
	}

	b.WriteString("\nBalas dengan nama account yang benar, misalnya: **")
	b.WriteString(accounts[0].Name)
	b.WriteString("**.")
	return b.String()
}

func buildAccountOutOfScopeMessage(terms []string, userCtx *domainauth.UserContext) string {
	input := strings.Join(terms, ", ")
	if input == "" {
		input = "account yang diminta"
	}

	scopeNote := "Account tersebut tidak berada dalam scope data yang dapat Anda gunakan."
	if userCtx != nil {
		switch {
		case userCtx.IsGlobalScope("accounts"):
			scopeNote = "Scope accounts user login adalah global, tetapi account tetap tidak lolos validasi. Periksa assignment account atau cache RBAC."
		case userCtx.IsTeamScope("accounts"):
			scopeNote = "Scope accounts user login adalah team. Account hanya dapat digunakan jika berada pada team/brick yang dikelola atau dimiliki sales bawahan."
		case userCtx.IsOwnScope("accounts"):
			scopeNote = "Scope accounts user login adalah own. Account hanya dapat digunakan jika di-assign ke user login atau berada pada brick user login."
		}
	}

	return fmt.Sprintf("Account **%s** ditemukan di database, tetapi tidak dapat digunakan berdasarkan RBAC/scope user login.\n\n%s\n\nGunakan account yang berada dalam akses Anda, atau minta admin mengubah assignment/scope account tersebut.", input, scopeNote)
}

func topAccountSuggestions(accounts []accountdomain.Account, terms []string, limit int) []accountdomain.Account {
	type scoredAccount struct {
		Account accountdomain.Account
		Score   int
	}
	scored := make([]scoredAccount, 0, len(accounts))
	for _, accountEntity := range accounts {
		score := scoreAccountMatch(accountEntity, terms)
		if score <= 0 {
			continue
		}
		scored = append(scored, scoredAccount{Account: accountEntity, Score: score})
	}
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Account.Name < scored[j].Account.Name
		}
		return scored[i].Score > scored[j].Score
	})

	if limit <= 0 || limit > len(scored) {
		limit = len(scored)
	}
	result := make([]accountdomain.Account, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scored[i].Account)
	}
	return result
}

func normalizeAccountMatchText(value string) string {
	tokens := strings.Fields(strings.ToLower(strings.Trim(value, ".,;:!?")))
	filtered := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.Trim(token, ".,;:!?")
		if token == "" || isAccountNameNoiseToken(token) {
			continue
		}
		filtered = append(filtered, token)
	}
	return strings.Join(filtered, " ")
}

func normalizeExactAccountText(value string) string {
	value = strings.ToLower(strings.Trim(value, ".,;:!?"))
	parts := strings.Fields(value)
	return strings.Join(parts, " ")
}

func (s *Service) resolveContactForAccountTask(terms []string, userCtx *domainauth.UserContext) (*contactdomain.Contact, error) {
	if s.contactRepo == nil {
		return nil, fmt.Errorf("contact repository tidak tersedia.")
	}

	results, _, err := s.contactRepo.List(&contactdomain.ListContactsRequest{
		Page:          1,
		PerPage:       10,
		Search:        strings.Join(terms, " "),
		ScopedUserIDs: s.scopedUserIDs(userCtx, "account"),
	})
	if err != nil {
		return nil, fmt.Errorf("Contact tidak dapat dicari saat ini.")
	}

	bestMatches := selectBestContactMatches(results, terms)
	if len(bestMatches) == 1 {
		contactEntity := bestMatches[0]
		return &contactEntity, nil
	}
	if len(bestMatches) > 1 {
		return nil, fmt.Errorf("Ditemukan beberapa contact yang mirip. Mohon sebutkan nama contact yang lebih spesifik.")
	}

	return nil, fmt.Errorf("Contact tidak ditemukan. Sebutkan nama contact yang lebih spesifik.")
}

func (s *Service) resolveLeadForTool(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) (string, *leaddomain.Lead, error) {
	for _, key := range []string{"id", "lead_id"} {
		if id := paramStr(params, key); id != "" {
			if leadEntity, err := s.leadRepo.FindByID(id); err == nil && leadEntity != nil {
				return leadEntity.ID, leadEntity, nil
			}
		}
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	if candidateIDs := entityIDs["lead"]; len(candidateIDs) == 1 {
		if leadEntity, err := s.leadRepo.FindByID(candidateIDs[0]); err == nil && leadEntity != nil {
			return leadEntity.ID, leadEntity, nil
		}
	}

	terms := collectEntityHints(params, "lead_name", "name", "full_name", "email", "phone", "company_name")
	if len(terms) == 0 {
		return "", nil, fmt.Errorf("Lead tidak ditemukan. Sebutkan nama lead, email, telepon, atau perusahaan tanpa perlu menyebut ID internal.")
	}

	results, _, err := s.leadRepo.List(&leaddomain.ListLeadsRequest{
		Page:          1,
		PerPage:       10,
		Search:        strings.Join(terms, " "),
		ScopedUserIDs: s.scopedUserIDs(userCtx, "lead"),
	})
	if err != nil {
		return "", nil, fmt.Errorf("Lead tidak dapat dicari saat ini.")
	}

	filtered := make([]leaddomain.Lead, 0, len(results))
	for _, leadEntity := range results {
		leadOwner := ""
		if leadEntity.AssignedTo != nil {
			leadOwner = *leadEntity.AssignedTo
		}
		if s.canAccessOwner(userCtx, "lead", leadOwner) {
			filtered = append(filtered, leadEntity)
		}
	}

	bestMatches := selectBestLeadMatches(filtered, terms)
	if len(bestMatches) == 1 {
		leadEntity := bestMatches[0]
		return leadEntity.ID, &leadEntity, nil
	}
	if len(bestMatches) > 1 {
		return "", nil, fmt.Errorf("Ditemukan beberapa lead yang mirip. Mohon pilih salah satu opsi berikut:\n%s", formatLeadDisambiguationOptions(bestMatches))
	}

	return "", nil, fmt.Errorf("Lead tidak ditemukan. Sebutkan nama lead, email, telepon, atau perusahaan yang lebih spesifik.")
}

func (s *Service) resolveDealForTool(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) (string, *pipelinedomain.Deal, error) {
	for _, key := range []string{"id", "deal_id"} {
		if id := paramStr(params, key); id != "" {
			if dealEntity, err := s.dealRepo.FindByID(id); err == nil && dealEntity != nil {
				return dealEntity.ID, dealEntity, nil
			}
		}
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	if candidateIDs := entityIDs["deal"]; len(candidateIDs) == 1 {
		if dealEntity, err := s.dealRepo.FindByID(candidateIDs[0]); err == nil && dealEntity != nil {
			return dealEntity.ID, dealEntity, nil
		}
	}
	if len(entityIDs["deal"]) == 0 && hasRecentDealCreationInHistory(history) {
		results, _, err := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
			Page:          1,
			PerPage:       1,
			Status:        "open",
			ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
		})
		if err == nil && len(results) > 0 {
			dealEntity := results[0]
			return dealEntity.ID, &dealEntity, nil
		}
	}

	terms := collectEntityHints(params, "deal_name", "name", "title", "account_name")
	if len(terms) == 0 {
		return "", nil, fmt.Errorf("Deal tidak ditemukan. Sebutkan nama deal tanpa perlu menyebut ID internal.")
	}

	results, _, err := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
		Page:          1,
		PerPage:       10,
		Search:        strings.Join(terms, " "),
		ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
	})
	if err != nil {
		return "", nil, fmt.Errorf("Deal tidak dapat dicari saat ini.")
	}

	filtered := make([]pipelinedomain.Deal, 0, len(results))
	for _, dealEntity := range results {
		dealOwner := ""
		if dealEntity.AssignedTo != nil {
			dealOwner = *dealEntity.AssignedTo
		}
		if s.canAccessOwner(userCtx, "deal", dealOwner) {
			filtered = append(filtered, dealEntity)
		}
	}

	bestMatches := selectBestDealMatches(filtered, terms)
	if len(bestMatches) == 1 {
		dealEntity := bestMatches[0]
		return dealEntity.ID, &dealEntity, nil
	}
	if len(bestMatches) > 1 {
		return "", nil, fmt.Errorf("Ditemukan beberapa deal yang mirip. Mohon pilih salah satu opsi berikut:\n%s", formatDealDisambiguationOptions(bestMatches))
	}

	return "", nil, fmt.Errorf("Deal tidak ditemukan. Sebutkan nama deal yang lebih spesifik.")
}

func (s *Service) resolveActivityDealContext(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext, accountID string) (string, string) {
	if s.dealRepo == nil {
		return "", ""
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	for _, candidateID := range entityIDs["deal"] {
		dealEntity, err := s.dealRepo.FindByID(candidateID)
		if err != nil || dealEntity == nil {
			continue
		}
		dealOwner := ""
		if dealEntity.AssignedTo != nil {
			dealOwner = *dealEntity.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "deal", dealOwner) {
			continue
		}
		if accountID != "" && dealEntity.AccountID != accountID {
			continue
		}
		return dealEntity.ID, dealEntity.AccountID
	}

	if accountID == "" {
		if hasRecentDealCreationInHistory(history) {
			deals, _, err := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
				Page:          1,
				PerPage:       1,
				Status:        "open",
				ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
			})
			if err == nil && len(deals) > 0 {
				dealEntity := deals[0]
				dealOwner := ""
				if dealEntity.AssignedTo != nil {
					dealOwner = *dealEntity.AssignedTo
				}
				if s.canAccessOwner(userCtx, "deal", dealOwner) {
					return dealEntity.ID, dealEntity.AccountID
				}
			}
		}
		return "", ""
	}
	deals, _, err := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
		Page:          1,
		PerPage:       10,
		AccountID:     accountID,
		Status:        "open",
		ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
	})
	if err != nil {
		return "", ""
	}
	for _, dealEntity := range deals {
		dealOwner := ""
		if dealEntity.AssignedTo != nil {
			dealOwner = *dealEntity.AssignedTo
		}
		if s.canAccessOwner(userCtx, "deal", dealOwner) {
			return dealEntity.ID, dealEntity.AccountID
		}
	}
	return "", ""
}

func hasRecentDealCreationInHistory(history []aidomain.ChatMessage) bool {
	checkedAssistantMessages := 0
	for i := len(history) - 1; i >= 0 && checkedAssistantMessages < 5; i-- {
		msg := history[i]
		if msg.Role != "assistant" {
			continue
		}
		checkedAssistantMessages++
		content := strings.ToLower(msg.Content)
		if strings.Contains(content, "deal berhasil dibuat") || strings.Contains(content, "deal baru") {
			return true
		}
	}
	return false
}

func (s *Service) resolveTaskForTool(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) (string, *taskdomain.Task, error) {
	if s.taskRepo == nil {
		return "", nil, fmt.Errorf("Task repository tidak tersedia.")
	}

	for _, key := range []string{"id", "task_id"} {
		if id := paramStr(params, key); id != "" {
			if taskEntity, err := s.taskRepo.FindByID(id); err == nil && taskEntity != nil {
				return taskEntity.ID, taskEntity, nil
			}
		}
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	if candidateIDs := entityIDs["task"]; len(candidateIDs) == 1 {
		if taskEntity, err := s.taskRepo.FindByID(candidateIDs[0]); err == nil && taskEntity != nil {
			return taskEntity.ID, taskEntity, nil
		}
	}

	terms := collectEntityHints(params, "task_name", "title", "name", "description")
	if len(terms) == 0 {
		return "", nil, fmt.Errorf("Task tidak ditemukan. Sebutkan judul task tanpa perlu menyebut ID internal.")
	}

	results, _, err := s.taskRepo.List(&taskdomain.ListTasksRequest{
		Page:          1,
		PerPage:       20,
		Search:        strings.Join(terms, " "),
		ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
	})
	if err != nil || len(results) == 0 {
		results, _, err = s.taskRepo.List(&taskdomain.ListTasksRequest{
			Page:          1,
			PerPage:       100,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("Task tidak dapat dicari saat ini.")
		}
	}

	filtered := make([]taskdomain.Task, 0, len(results))
	for _, taskEntity := range results {
		taskOwner := ""
		if taskEntity.AssignedTo != nil {
			taskOwner = *taskEntity.AssignedTo
		}
		if s.canAccessOwner(userCtx, "task", taskOwner) {
			filtered = append(filtered, taskEntity)
		}
	}

	bestMatches := selectBestTaskMatches(filtered, terms)
	if len(bestMatches) == 1 {
		taskEntity := bestMatches[0]
		return taskEntity.ID, &taskEntity, nil
	}
	if len(bestMatches) > 1 {
		return "", nil, fmt.Errorf("Ditemukan beberapa task yang mirip. Mohon pilih salah satu opsi berikut:\n%s", formatTaskDisambiguationOptions(bestMatches))
	}

	return "", nil, fmt.Errorf("Task tidak ditemukan. Sebutkan judul task yang lebih spesifik.")
}

func (s *Service) resolveScheduleForTool(params map[string]interface{}, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) (*scheduledomain.ScheduleResponse, error) {
	if s.scheduleService == nil {
		return nil, fmt.Errorf("Schedule service tidak tersedia.")
	}

	for _, key := range []string{"id", "schedule_id"} {
		if id := paramStr(params, key); id != "" {
			resp, err := s.scheduleService.GetScheduleByID(id)
			if err != nil || resp == nil {
				return nil, fmt.Errorf("Jadwal tidak ditemukan.")
			}
			if !s.canAccessOwner(userCtx, "schedule", resp.UserID) {
				return nil, fmt.Errorf("Anda tidak memiliki akses ke jadwal tersebut.")
			}
			return resp, nil
		}
	}

	entityIDs := s.extractEntityIDsFromHistory(history)
	if candidateIDs := entityIDs["schedule"]; len(candidateIDs) == 1 {
		resp, err := s.scheduleService.GetScheduleByID(candidateIDs[0])
		if err == nil && resp != nil && s.canAccessOwner(userCtx, "schedule", resp.UserID) {
			return resp, nil
		}
	}

	terms := collectEntityHints(params, "schedule_title", "schedule_name", "title", "name", "description")
	if len(terms) == 0 {
		terms = []string{"meeting"}
	}

	results, _, err := s.scheduleService.ListSchedules(&scheduledomain.ListSchedulesRequest{
		Page:          1,
		PerPage:       20,
		Search:        strings.Join(terms, " "),
		ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"),
	})
	if err != nil || len(results) == 0 {
		results, _, err = s.scheduleService.ListSchedules(&scheduledomain.ListSchedulesRequest{
			Page:          1,
			PerPage:       100,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"),
		})
		if err != nil {
			return nil, fmt.Errorf("Jadwal tidak dapat dicari saat ini.")
		}
	}

	filtered := make([]scheduledomain.ScheduleResponse, 0, len(results))
	for _, scheduleEntity := range results {
		if s.canAccessOwner(userCtx, "schedule", scheduleEntity.UserID) {
			filtered = append(filtered, scheduleEntity)
		}
	}

	bestMatches := selectBestScheduleMatches(filtered, terms)
	if len(bestMatches) == 1 {
		return &bestMatches[0], nil
	}
	if len(bestMatches) > 1 {
		return nil, fmt.Errorf("Ditemukan beberapa jadwal yang mirip. Mohon pilih salah satu opsi berikut:\n%s", formatScheduleDisambiguationOptions(bestMatches))
	}

	return nil, fmt.Errorf("Jadwal tidak ditemukan. Sebutkan judul jadwal yang lebih spesifik atau tampilkan daftar schedule terlebih dahulu.")
}

func collectEntityHints(params map[string]interface{}, keys ...string) []string {
	seen := make(map[string]struct{})
	hints := make([]string, 0, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(paramStr(params, key))
		if value == "" {
			continue
		}
		normalized := strings.ToLower(value)
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		hints = append(hints, value)
	}
	return hints
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		normalized := strings.ToLower(value)
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, value)
	}
	return result
}

func firstNonEmptyParam(params map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value := paramStr(params, key); value != "" {
			return value
		}
	}
	return ""
}

func hasLeadReference(params map[string]interface{}) bool {
	for _, key := range []string{"id", "lead_id", "lead_name", "name", "full_name", "email", "phone", "company_name"} {
		if paramStr(params, key) != "" {
			return true
		}
	}
	return false
}

func hasAccountReferenceForActivity(params map[string]interface{}) bool {
	for _, key := range []string{"account_name", "company_name", "company", "contact_id", "contact_name", "contact", "pic_name"} {
		if paramStr(params, key) != "" {
			return true
		}
	}
	return false
}

func (s *Service) shouldResolveLeadForActivity(params map[string]interface{}, history []aidomain.ChatMessage) bool {
	for _, key := range []string{"lead_id", "lead_name", "full_name", "email", "phone"} {
		if paramStr(params, key) != "" {
			return true
		}
	}
	if paramStr(params, "deal_id") != "" || paramStr(params, "deal_name") != "" || paramStr(params, "title") != "" || hasAccountReferenceForActivity(params) {
		return false
	}
	entityIDs := s.extractEntityIDsFromHistory(history)
	if len(entityIDs["deal"]) > 0 || len(entityIDs["account"]) > 0 {
		return false
	}
	return paramStr(params, "company_name") != "" || paramStr(params, "name") != ""
}

func normalizeActivityType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "visit", "call", "email", "task", "deal":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "call"
	}
}

func toolStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (s *Service) productInterestsFromParams(params map[string]interface{}) []map[string]string {
	names := paramStringSlice(params, "product_interests")
	names = append(names, paramStringSlice(params, "product_names")...)
	names = append(names, paramStringSlice(params, "products")...)
	if single := firstNonEmptyParam(params, "product_interest", "product_name", "product"); single != "" {
		names = append(names, splitCSVLike(single)...)
	}

	names = uniqueNonEmpty(names)
	if len(names) == 0 {
		return nil
	}

	interests := make([]map[string]string, 0, len(names))
	for _, name := range names {
		item := map[string]string{"product_name": name}
		if s.productRepo != nil {
			if productID, productName := s.resolveProductByName(name); productName != "" {
				item["product_name"] = productName
				if productID != "" {
					item["product_id"] = productID
				}
			}
		}
		interests = append(interests, item)
	}
	return interests
}

func (s *Service) dealProductItemsFromParams(params map[string]interface{}) ([]pipelinedomain.CreateDealProductItemRequest, error) {
	interests := s.productInterestsFromParams(params)
	if len(interests) == 0 {
		return nil, nil
	}

	quantity := paramInt(params, "quantity")
	if quantity < 1 {
		quantity = 1
	}

	items := make([]pipelinedomain.CreateDealProductItemRequest, 0, len(interests))
	unresolved := make([]string, 0)
	for _, interest := range interests {
		productID := strings.TrimSpace(interest["product_id"])
		productName := strings.TrimSpace(interest["product_name"])
		if productID == "" {
			unresolved = append(unresolved, productName)
			continue
		}
		item := pipelinedomain.CreateDealProductItemRequest{
			ProductID: productID,
			Quantity:  quantity,
		}
		if price := paramInt64(params, "unit_price"); price != nil {
			item.UnitPrice = price
		} else if price := paramInt64(params, "price"); price != nil {
			item.UnitPrice = price
		}
		if discount := paramInt64(params, "discount_amount"); discount != nil {
			item.DiscountAmount = discount
		}
		item.Notes = paramStr(params, "notes")
		items = append(items, item)
	}

	if len(unresolved) > 0 {
		return nil, fmt.Errorf("produk tidak ditemukan: %s. Sebutkan nama produk sesuai katalog.", strings.Join(unresolved, ", "))
	}
	return items, nil
}

func (s *Service) productInterestMetadataFromParams(params map[string]interface{}) ([]map[string]interface{}, error) {
	basicInterests, err := s.strictProductInterestsFromParams(params)
	if err != nil {
		return nil, err
	}
	if len(basicInterests) == 0 {
		return nil, nil
	}

	interestLevel := clampInt(paramInt(params, "interest_level"), 0, 5)
	if interestLevel == 0 {
		interestLevel = clampInt(paramInt(params, "rating"), 0, 5)
	}
	if interestLevel == 0 {
		interestLevel = clampInt(extractStarRating(firstNonEmptyParam(params, "notes", "note", "description")), 0, 5)
	}
	quantity := paramInt(params, "quantity")
	price := paramInt(params, "price")

	result := make([]map[string]interface{}, 0, len(basicInterests))
	for _, interest := range basicInterests {
		item := map[string]interface{}{
			"product_name": interest["product_name"],
		}
		if productID := interest["product_id"]; productID != "" {
			item["product_id"] = productID
		}
		if interestLevel > 0 {
			item["interest_level"] = interestLevel
		}
		if quantity > 0 {
			item["quantity"] = quantity
		}
		if price > 0 {
			item["price"] = price
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *Service) strictProductInterestsFromParams(params map[string]interface{}) ([]map[string]string, error) {
	names := paramStringSlice(params, "product_interests")
	names = append(names, paramStringSlice(params, "product_names")...)
	names = append(names, paramStringSlice(params, "products")...)
	if single := firstNonEmptyParam(params, "product_interest", "product_name", "product"); single != "" {
		names = append(names, splitCSVLike(single)...)
	}

	names = uniqueNonEmpty(names)
	if len(names) == 0 {
		return nil, nil
	}
	if s.productRepo == nil {
		return nil, fmt.Errorf("Product repository tidak tersedia. Product interest harus mengacu pada produk yang ada di master Products.")
	}

	resolved := make([]map[string]string, 0, len(names))
	unresolved := make([]string, 0)
	ambiguous := make([]string, 0)
	for _, name := range names {
		products, err := s.searchProductsForTool([]string{name})
		if err != nil {
			return nil, fmt.Errorf("Produk **%s** tidak dapat divalidasi saat ini.", name)
		}
		matches := selectBestProductMatches(products, []string{name})
		switch len(matches) {
		case 0:
			unresolved = append(unresolved, name)
		case 1:
			resolved = append(resolved, map[string]string{
				"product_id":   matches[0].ID,
				"product_name": matches[0].Name,
			})
		default:
			ambiguous = append(ambiguous, name)
		}
	}

	if len(unresolved) > 0 || len(ambiguous) > 0 {
		parts := make([]string, 0, 2)
		if len(unresolved) > 0 {
			parts = append(parts, fmt.Sprintf("produk tidak ditemukan di master Products: **%s**", strings.Join(unresolved, ", ")))
		}
		if len(ambiguous) > 0 {
			parts = append(parts, fmt.Sprintf("produk masih ambigu: **%s**", strings.Join(ambiguous, ", ")))
		}
		return nil, fmt.Errorf("%s. Product interest tidak disimpan. Gunakan nama produk atau SKU yang tepat dari menu Products.", strings.Join(parts, "; "))
	}

	return resolved, nil
}

func (s *Service) needProductsFromParams(params map[string]interface{}) []leadqualificationdomain.NeedProduct {
	interests := s.productInterestsFromParams(params)
	if len(interests) == 0 {
		return nil
	}

	products := make([]leadqualificationdomain.NeedProduct, 0, len(interests))
	for _, item := range interests {
		products = append(products, leadqualificationdomain.NeedProduct{
			ProductID:   item["product_id"],
			ProductName: item["product_name"],
		})
	}
	return products
}

func isProductInterestOnlyParams(params map[string]interface{}) bool {
	hasProduct := len(paramStringSlice(params, "product_interests")) > 0 ||
		len(paramStringSlice(params, "product_names")) > 0 ||
		len(paramStringSlice(params, "products")) > 0 ||
		firstNonEmptyParam(params, "product_interest", "product_name", "product") != ""
	if !hasProduct {
		return false
	}

	for _, key := range []string{
		"budget_target_amount", "budget_target_currency", "budget_confirmed", "budget_notes",
		"authority_target_person", "authority_person", "decision_maker", "authority_target_role", "authority_role", "decision_maker_role", "authority_confirmed", "authority_notes",
		"need_priority_level", "need_confirmed", "need_notes",
		"timeline_target_date", "timeline_flexibility", "timeline_confirmed", "timeline_notes",
	} {
		if value, ok := params[key]; ok && value != nil {
			if s, ok := value.(string); ok && strings.TrimSpace(s) == "" {
				continue
			}
			return false
		}
	}

	return true
}

func (s *Service) resolveProductByName(name string) (string, string) {
	if s.productRepo == nil || strings.TrimSpace(name) == "" {
		return "", ""
	}
	results, _, err := s.productRepo.List(&productdomain.ListProductsRequest{
		Page:    1,
		PerPage: 5,
		Search:  name,
		Status:  "active",
	})
	if err != nil || len(results) == 0 {
		return "", ""
	}

	normalized := strings.ToLower(strings.TrimSpace(name))
	for _, productEntity := range results {
		if strings.EqualFold(productEntity.Name, normalized) || strings.EqualFold(productEntity.SKU, normalized) {
			return productEntity.ID, productEntity.Name
		}
	}
	return results[0].ID, results[0].Name
}

func (s *Service) toolUpdateProductStatus(params map[string]interface{}) toolResult {
	if s.productRepo == nil {
		return toolResult{Success: false, Entity: "Produk", Message: "Product repository tidak tersedia."}
	}

	status, ok := normalizeProductStatus(paramStr(params, "status"))
	if !ok {
		return toolResult{
			Success: false,
			Confirm: true,
			Entity:  "Produk",
			Message: "Status tujuan produk belum jelas. Sebutkan status baru: **active/aktif** atau **inactive/nonaktif**.",
		}
	}

	productEntity, err := s.resolveProductForTool(params)
	if err != nil {
		if confirmationErr, ok := err.(*toolConfirmationError); ok {
			return toolResult{Success: false, Confirm: true, Entity: "Produk", Message: confirmationErr.Message}
		}
		return toolResult{Success: false, Entity: "Produk", Message: err.Error()}
	}

	if productEntity.Status == status {
		return toolResult{
			Success: true,
			Entity:  "Produk",
			ID:      productEntity.ID,
			Action:  "diperbarui",
			Message: fmt.Sprintf("Produk **%s** sudah berstatus **%s**.", productEntity.Name, status),
			PageURL: "/products",
			Icon:    "package",
		}
	}

	productEntity.Status = status
	if err := s.productRepo.Update(productEntity); err != nil {
		return toolResult{Success: false, Entity: "Produk", Message: "Status produk gagal diperbarui."}
	}

	return toolResult{
		Success: true,
		Entity:  "Produk",
		ID:      productEntity.ID,
		Action:  "diperbarui",
		Message: fmt.Sprintf("Status produk **%s** berhasil diubah menjadi **%s**.", productEntity.Name, status),
		PageURL: "/products",
		Icon:    "package",
	}
}

func normalizeProductStatus(status string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active", "aktif", "activate", "activated", "enabled", "enable":
		return "active", true
	case "inactive", "nonaktif", "non-aktif", "deactivate", "deactivated", "disabled", "disable":
		return "inactive", true
	default:
		return "", false
	}
}

func (s *Service) resolveProductForTool(params map[string]interface{}) (*productdomain.Product, error) {
	for _, key := range []string{"id", "product_id"} {
		if id := paramStr(params, key); id != "" {
			productEntity, err := s.productRepo.FindByID(id)
			if err != nil || productEntity == nil {
				return nil, fmt.Errorf("produk yang dipilih tidak ditemukan.")
			}
			return productEntity, nil
		}
	}

	terms := collectEntityHints(params, "product_name", "name", "sku")
	if len(terms) == 0 {
		return nil, &toolConfirmationError{Message: "Nama produk wajib diisi. Sebutkan nama produk atau SKU yang ingin diubah statusnya."}
	}

	products, err := s.searchProductsForTool(terms)
	if err != nil {
		return nil, fmt.Errorf("produk tidak dapat dicari saat ini.")
	}
	matches := selectBestProductMatches(products, terms)
	if len(matches) == 1 {
		return &matches[0], nil
	}
	if len(matches) > 1 {
		return nil, &toolConfirmationError{Message: buildProductConfirmationMessage(terms, matches)}
	}
	return nil, &toolConfirmationError{Message: fmt.Sprintf("Saya belum menemukan produk yang cocok untuk **%s**. Sebutkan nama produk atau SKU yang lebih spesifik.", strings.Join(terms, ", "))}
}

func (s *Service) searchProductsForTool(terms []string) ([]productdomain.Product, error) {
	resultsByID := make(map[string]productdomain.Product)
	for _, searchTerm := range productSearchTermCandidates(terms) {
		results, _, err := s.productRepo.List(&productdomain.ListProductsRequest{
			Page:    1,
			PerPage: 10,
			Search:  searchTerm,
		})
		if err != nil {
			return nil, err
		}
		for _, productEntity := range results {
			resultsByID[productEntity.ID] = productEntity
		}
	}

	if len(resultsByID) == 0 || !hasPositiveProductMatch(resultsByID, terms) {
		results, _, err := s.productRepo.List(&productdomain.ListProductsRequest{
			Page:    1,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		for _, productEntity := range results {
			resultsByID[productEntity.ID] = productEntity
		}
	}

	products := make([]productdomain.Product, 0, len(resultsByID))
	for _, productEntity := range resultsByID {
		products = append(products, productEntity)
	}
	return products, nil
}

func productSearchTermCandidates(terms []string) []string {
	candidates := make([]string, 0)
	for _, term := range terms {
		normalized := strings.ToLower(strings.Trim(strings.TrimSpace(term), ".,;:!?"))
		if normalized == "" {
			continue
		}
		candidates = appendUnique(candidates, normalized)
		for _, token := range strings.Fields(normalized) {
			token = strings.Trim(token, ".,;:!?")
			if len(token) < 3 || isStopWord(token) {
				continue
			}
			candidates = appendUnique(candidates, token)
		}
	}
	return candidates
}

func hasPositiveProductMatch(products map[string]productdomain.Product, terms []string) bool {
	for _, productEntity := range products {
		if scoreProductMatch(productEntity, terms) > 0 {
			return true
		}
	}
	return false
}

func selectBestProductMatches(products []productdomain.Product, terms []string) []productdomain.Product {
	bestScore := 0
	bestMatches := make([]productdomain.Product, 0, 1)
	for _, productEntity := range products {
		score := scoreProductMatch(productEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []productdomain.Product{productEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, productEntity)
		}
	}
	return bestMatches
}

func scoreProductMatch(productEntity productdomain.Product, terms []string) int {
	name := strings.ToLower(strings.TrimSpace(productEntity.Name))
	sku := strings.ToLower(strings.TrimSpace(productEntity.SKU))
	barcode := strings.ToLower(strings.TrimSpace(productEntity.Barcode))

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case sku != "" && normalized == sku:
			score += 120
		case barcode != "" && normalized == barcode:
			score += 120
		case name != "" && normalized == name:
			score += 100
		case name != "" && strings.Contains(name, normalized):
			score += 70
		case sku != "" && strings.Contains(sku, normalized):
			score += 60
		}

		for _, token := range productSearchTermCandidates([]string{normalized}) {
			if token != "" && name != "" && strings.Contains(name, token) {
				score += 25
			}
			if token != "" && sku != "" && strings.Contains(sku, token) {
				score += 20
			}
		}
	}
	return score
}

func buildProductConfirmationMessage(terms []string, products []productdomain.Product) string {
	var b strings.Builder
	input := strings.Join(terms, ", ")
	if input == "" {
		input = "produk yang diminta"
	}
	b.WriteString(fmt.Sprintf("Saya menemukan beberapa produk yang mirip dengan **%s**. Mohon konfirmasi produk yang dimaksud:\n\n", input))
	limit := min(len(products), 5)
	for i := 0; i < limit; i++ {
		productEntity := products[i]
		b.WriteString(fmt.Sprintf("%d. **%s** — SKU: %s, status: %s\n", i+1, productEntity.Name, productEntity.SKU, productEntity.Status))
	}
	b.WriteString("\nBalas dengan nama produk atau SKU yang benar.")
	return b.String()
}

func splitCSVLike(value string) []string {
	value = strings.ReplaceAll(value, ";", ",")
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeNeedPriority(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low", "medium", "high", "critical":
		return strings.ToLower(strings.TrimSpace(value))
	case "urgent", "mendesak", "kritis":
		return "critical"
	default:
		return ""
	}
}

func normalizeTimelineFlexibility(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "fixed", "flexible", "urgent":
		return strings.ToLower(strings.TrimSpace(value))
	case "tetap", "pasti":
		return "fixed"
	case "fleksibel":
		return "flexible"
	case "mendesak":
		return "urgent"
	default:
		return ""
	}
}

func selectBestLeadMatches(leads []leaddomain.Lead, terms []string) []leaddomain.Lead {
	bestScore := 0
	bestMatches := make([]leaddomain.Lead, 0, 1)

	for _, leadEntity := range leads {
		score := scoreLeadMatch(leadEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []leaddomain.Lead{leadEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, leadEntity)
		}
	}

	return bestMatches
}

func scoreLeadMatch(leadEntity leaddomain.Lead, terms []string) int {
	fullName := strings.ToLower(strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName)))
	email := strings.ToLower(strings.TrimSpace(leadEntity.Email))
	phone := strings.ToLower(strings.TrimSpace(leadEntity.Phone))
	company := strings.ToLower(strings.TrimSpace(leadEntity.CompanyName))

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case email != "" && normalized == email:
			score += 120
		case phone != "" && normalized == phone:
			score += 120
		case fullName != "" && normalized == fullName:
			score += 100
		case company != "" && normalized == company:
			score += 90
		case fullName != "" && strings.Contains(fullName, normalized):
			score += 60
		case company != "" && strings.Contains(company, normalized):
			score += 50
		case email != "" && strings.Contains(email, normalized):
			score += 40
		case phone != "" && strings.Contains(phone, normalized):
			score += 40
		}
	}

	return score
}

func selectBestDealMatches(deals []pipelinedomain.Deal, terms []string) []pipelinedomain.Deal {
	bestScore := 0
	bestMatches := make([]pipelinedomain.Deal, 0, 1)

	for _, dealEntity := range deals {
		score := scoreDealMatch(dealEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []pipelinedomain.Deal{dealEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, dealEntity)
		}
	}

	return bestMatches
}

func selectBestAccountMatches(accounts []accountdomain.Account, terms []string) []accountdomain.Account {
	bestScore := 0
	bestMatches := make([]accountdomain.Account, 0, 1)

	for _, accountEntity := range accounts {
		score := scoreAccountMatch(accountEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []accountdomain.Account{accountEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, accountEntity)
		}
	}

	return bestMatches
}

func selectBestContactMatches(contacts []contactdomain.Contact, terms []string) []contactdomain.Contact {
	bestScore := 0
	bestMatches := make([]contactdomain.Contact, 0, 1)

	for _, contactEntity := range contacts {
		score := scoreContactMatch(contactEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []contactdomain.Contact{contactEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, contactEntity)
		}
	}

	return bestMatches
}

func selectBestTaskMatches(tasks []taskdomain.Task, terms []string) []taskdomain.Task {
	bestScore := 0
	bestMatches := make([]taskdomain.Task, 0, 1)

	for _, taskEntity := range tasks {
		score := scoreTaskMatch(taskEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []taskdomain.Task{taskEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, taskEntity)
		}
	}

	return bestMatches
}

func selectBestScheduleMatches(schedules []scheduledomain.ScheduleResponse, terms []string) []scheduledomain.ScheduleResponse {
	bestScore := 0
	bestMatches := make([]scheduledomain.ScheduleResponse, 0, 1)

	for _, scheduleEntity := range schedules {
		score := scoreScheduleMatch(scheduleEntity, terms)
		if score == 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			bestMatches = []scheduledomain.ScheduleResponse{scheduleEntity}
			continue
		}
		if score == bestScore {
			bestMatches = append(bestMatches, scheduleEntity)
		}
	}

	return bestMatches
}

func scoreAccountMatch(accountEntity accountdomain.Account, terms []string) int {
	name := strings.ToLower(strings.TrimSpace(accountEntity.Name))
	email := strings.ToLower(strings.TrimSpace(accountEntity.Email))
	phone := strings.ToLower(strings.TrimSpace(accountEntity.Phone))

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case email != "" && normalized == email:
			score += 120
		case phone != "" && normalized == phone:
			score += 120
		case name != "" && normalized == name:
			score += 100
		case name != "" && strings.Contains(name, normalized):
			score += 70
		case email != "" && strings.Contains(email, normalized):
			score += 40
		case phone != "" && strings.Contains(phone, normalized):
			score += 40
		}

		for _, token := range accountSearchTermCandidates([]string{normalized}) {
			if token != "" && name != "" && strings.Contains(name, token) {
				score += 25
			}
		}
	}

	return score
}

func scoreContactMatch(contactEntity contactdomain.Contact, terms []string) int {
	name := strings.ToLower(strings.TrimSpace(contactEntity.Name))
	email := strings.ToLower(strings.TrimSpace(contactEntity.Email))
	phone := strings.ToLower(strings.TrimSpace(contactEntity.Phone))
	position := strings.ToLower(strings.TrimSpace(contactEntity.Position))

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case email != "" && normalized == email:
			score += 120
		case phone != "" && normalized == phone:
			score += 120
		case name != "" && normalized == name:
			score += 100
		case name != "" && strings.Contains(name, normalized):
			score += 70
		case email != "" && strings.Contains(email, normalized):
			score += 40
		case phone != "" && strings.Contains(phone, normalized):
			score += 40
		case position != "" && strings.Contains(position, normalized):
			score += 20
		}
	}

	return score
}

func scoreDealMatch(dealEntity pipelinedomain.Deal, terms []string) int {
	title := strings.ToLower(strings.TrimSpace(dealEntity.Title))
	accountName := ""
	if dealEntity.Account != nil {
		accountName = strings.ToLower(strings.TrimSpace(dealEntity.Account.Name))
	}

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case title != "" && normalized == title:
			score += 120
		case accountName != "" && normalized == accountName:
			score += 80
		case title != "" && strings.Contains(title, normalized):
			score += 70
		case accountName != "" && strings.Contains(accountName, normalized):
			score += 40
		}
	}

	return score
}

func scoreTaskMatch(taskEntity taskdomain.Task, terms []string) int {
	title := strings.ToLower(strings.TrimSpace(taskEntity.Title))
	description := strings.ToLower(strings.TrimSpace(taskEntity.Description))
	accountName := ""
	if taskEntity.Account != nil {
		accountName = strings.ToLower(strings.TrimSpace(taskEntity.Account.Name))
	}
	dealTitle := ""
	if taskEntity.Deal != nil {
		dealTitle = strings.ToLower(strings.TrimSpace(taskEntity.Deal.Title))
	}
	leadName := ""
	if taskEntity.Lead != nil {
		leadName = strings.ToLower(strings.TrimSpace(strings.TrimSpace(taskEntity.Lead.FirstName + " " + taskEntity.Lead.LastName)))
	}

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case title != "" && normalized == title:
			score += 140
		case title != "" && strings.Contains(title, normalized):
			score += 100
		case normalized != "" && strings.Contains(normalized, title):
			score += 90
		case accountName != "" && strings.Contains(normalized, accountName):
			score += 50
		case dealTitle != "" && strings.Contains(normalized, dealTitle):
			score += 50
		case leadName != "" && strings.Contains(normalized, leadName):
			score += 40
		case description != "" && strings.Contains(description, normalized):
			score += 30
		}
	}

	return score
}

func scoreScheduleMatch(scheduleEntity scheduledomain.ScheduleResponse, terms []string) int {
	title := strings.ToLower(strings.TrimSpace(scheduleEntity.Title))
	description := ""
	if scheduleEntity.Description != nil {
		description = strings.ToLower(strings.TrimSpace(*scheduleEntity.Description))
	}
	taskTitle := ""
	if scheduleEntity.Task != nil {
		taskTitle = strings.ToLower(strings.TrimSpace(scheduleEntity.Task.Title))
	}

	score := 0
	for _, term := range terms {
		normalized := strings.ToLower(strings.TrimSpace(term))
		switch {
		case normalized == "":
			continue
		case title != "" && normalized == title:
			score += 140
		case title != "" && strings.Contains(title, normalized):
			score += 100
		case title != "" && strings.Contains(normalized, title):
			score += 90
		case taskTitle != "" && strings.Contains(taskTitle, normalized):
			score += 60
		case taskTitle != "" && strings.Contains(normalized, taskTitle):
			score += 50
		case description != "" && strings.Contains(description, normalized):
			score += 30
		}
	}
	return score
}

func formatLeadDisambiguationOptions(leads []leaddomain.Lead) string {
	limit := min(len(leads), 5)
	options := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		leadEntity := leads[i]
		fullName := strings.TrimSpace(strings.TrimSpace(leadEntity.FirstName + " " + leadEntity.LastName))
		if fullName == "" {
			fullName = "Tanpa nama"
		}

		parts := []string{fmt.Sprintf("%d. %s", i+1, fullName)}
		if leadEntity.CompanyName != "" {
			parts = append(parts, leadEntity.CompanyName)
		}
		if leadEntity.Email != "" {
			parts = append(parts, leadEntity.Email)
		}
		if leadEntity.Phone != "" {
			parts = append(parts, leadEntity.Phone)
		}
		if leadEntity.LeadStatus != "" {
			parts = append(parts, "status "+leadEntity.LeadStatus)
		}
		if leadEntity.City != "" {
			parts = append(parts, leadEntity.City)
		}
		options = append(options, "- "+strings.Join(parts, " | "))
	}

	return strings.Join(options, "\n")
}

func formatScheduleDisambiguationOptions(schedules []scheduledomain.ScheduleResponse) string {
	limit := min(len(schedules), 5)
	options := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		scheduleEntity := schedules[i]
		parts := []string{fmt.Sprintf("%d. %s", i+1, scheduleEntity.Title)}
		if scheduleEntity.Status != "" {
			parts = append(parts, "status "+scheduleEntity.Status)
		}
		parts = append(parts, "tanggal "+scheduleEntity.ScheduledAt.Format("2006-01-02 15:04"))
		if scheduleEntity.Task != nil && scheduleEntity.Task.Title != "" {
			parts = append(parts, "task "+scheduleEntity.Task.Title)
		}
		options = append(options, "- "+strings.Join(parts, " | "))
	}

	return strings.Join(options, "\n")
}

func formatTaskDisambiguationOptions(tasks []taskdomain.Task) string {
	limit := min(len(tasks), 5)
	options := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		taskEntity := tasks[i]
		parts := []string{fmt.Sprintf("%d. %s", i+1, taskEntity.Title)}
		if taskEntity.Status != "" {
			parts = append(parts, "status "+taskEntity.Status)
		}
		if taskEntity.Account != nil && taskEntity.Account.Name != "" {
			parts = append(parts, taskEntity.Account.Name)
		}
		if taskEntity.DueDate != nil {
			parts = append(parts, "due "+taskEntity.DueDate.Format("2006-01-02"))
		}
		options = append(options, "- "+strings.Join(parts, " | "))
	}

	return strings.Join(options, "\n")
}

func formatDealDisambiguationOptions(deals []pipelinedomain.Deal) string {
	limit := min(len(deals), 5)
	options := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		dealEntity := deals[i]
		parts := []string{fmt.Sprintf("%d. %s", i+1, dealEntity.Title)}
		if dealEntity.Account != nil && dealEntity.Account.Name != "" {
			parts = append(parts, dealEntity.Account.Name)
		}
		if dealEntity.Status != "" {
			parts = append(parts, "status "+dealEntity.Status)
		}
		if dealEntity.Stage != nil && dealEntity.Stage.Name != "" {
			parts = append(parts, "stage "+dealEntity.Stage.Name)
		}
		if !dealEntity.UpdatedAt.IsZero() {
			parts = append(parts, "updated "+dealEntity.UpdatedAt.Format("2006-01-02 15:04"))
		}
		options = append(options, "- "+strings.Join(parts, " | "))
	}

	return strings.Join(options, "\n")
}

// ----------------------------------------------------------------------------
// Parameter helpers
// ----------------------------------------------------------------------------

func paramStr(params map[string]interface{}, key string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func paramStrOr(params map[string]interface{}, key, def string) string {
	if v := paramStr(params, key); v != "" {
		return v
	}
	return def
}

func paramFloat(params map[string]interface{}, key string) float64 {
	if v, ok := params[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func paramInt(params map[string]interface{}, key string) int {
	if value := paramInt64(params, key); value != nil {
		return int(*value)
	}
	return 0
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func extractStarRating(text string) int {
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "bintang") && !strings.Contains(lower, "star") {
		return 0
	}
	re := regexp.MustCompile(`([1-5])\s*(?:bintang|star)`)
	match := re.FindStringSubmatch(lower)
	if len(match) < 2 {
		return 0
	}
	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0
	}
	return value
}

func paramInt64(params map[string]interface{}, key string) *int64 {
	if v, ok := params[key]; ok {
		switch value := v.(type) {
		case float64:
			result := int64(value)
			return &result
		case int64:
			return &value
		case int:
			result := int64(value)
			return &result
		case string:
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				return nil
			}
			parsed, err := strconv.ParseInt(strings.ReplaceAll(trimmed, ".", ""), 10, 64)
			if err != nil {
				return nil
			}
			return &parsed
		}
	}
	return nil
}

func paramBoolPtr(params map[string]interface{}, key string) (*bool, bool) {
	if v, ok := params[key]; ok {
		switch value := v.(type) {
		case bool:
			return &value, true
		case string:
			normalized := strings.ToLower(strings.TrimSpace(value))
			switch normalized {
			case "true", "yes", "ya", "y", "confirmed", "terkonfirmasi", "sudah":
				result := true
				return &result, true
			case "false", "no", "tidak", "n", "unconfirmed", "belum":
				result := false
				return &result, true
			}
		}
	}
	return nil, false
}

func paramStringSlice(params map[string]interface{}, key string) []string {
	if v, ok := params[key]; ok {
		if arr, ok := v.([]string); ok {
			return arr
		}
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
		if s, ok := v.(string); ok {
			return splitCSVLike(s)
		}
	}
	return nil
}

// parseFlexibleTime tries several common date/time formats.
func parseFlexibleTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, s, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised date format: %s", s)
}

func isDateOnlyText(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if strings.Contains(value, "T") || strings.Contains(value, ":") {
		return false
	}
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return true
	}
	return false
}
