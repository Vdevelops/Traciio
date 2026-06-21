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
	"regexp"
	"sort"
	"strings"
	"time"

	aidomain "github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	leaddomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	pipelinedomain "github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	route_optimization_domain "github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	scheduledomain "github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	taskdomain "github.com/gilabs/crm-healthcare/api/internal/domain/task"
)

// toolCallPattern matches <!-- TOOL_CALL:{...} --> markers in LLM responses.
// (?s) enables dotall mode so '.' also matches newlines — the LLM often emits pretty-printed JSON.
var toolCallPattern = regexp.MustCompile(`(?s)<!--\s*TOOL_CALL:(.*?)\s*-->`)

// ToolCall is the parsed representation of an LLM tool invocation.
type ToolCall struct {
	Tool   string                 `json:"tool"`
	Params map[string]interface{} `json:"params"`
}

// toolResult holds the outcome of a single tool execution.
type toolResult struct {
	Success bool
	Entity  string // Human-readable entity name (e.g. "Task", "Lead")
	ID      string // ID of the created/updated entity
	Message string // Short confirmation detail shown under the success header
	PageURL string // CRM page URL for the action card
	Icon    string // Lucide icon name for the action card
}

// processToolCalls scans an LLM response for TOOL_CALL markers, executes each
// tool synchronously, and replaces every marker with its result block.
func (s *Service) processToolCalls(response string, userID string, history []aidomain.ChatMessage, userCtx *domainauth.UserContext) string {
	matches := toolCallPattern.FindAllStringSubmatchIndex(response, -1)
	if len(matches) == 0 {
		return response
	}

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
	if !s.canRunTool(call.Tool, userCtx) {
		return toolResult{
			Success: false,
			Entity:  "AI Action",
			Message: fmt.Sprintf("Anda tidak memiliki permission untuk menjalankan tool '%s'.", call.Tool),
		}
	}

	switch call.Tool {
	case "create_task":
		return s.toolCreateTask(call.Params, userID)
	case "create_lead":
		return s.toolCreateLead(call.Params, userID)
	case "create_deal":
		return s.toolCreateDeal(call.Params, userID, userCtx)
	case "create_schedule":
		return s.toolCreateSchedule(call.Params, userID)
	case "create_route":
		return s.toolCreateRoute(call.Params, userID, history, userCtx)
	case "update_task_status":
		return s.toolUpdateTaskStatus(call.Params, userCtx)
	case "update_lead_status":
		return s.toolUpdateLeadStatus(call.Params, history, userCtx)
	case "update_deal_stage":
		return s.toolUpdateDealStage(call.Params, userID, history, userCtx)
	default:
		return toolResult{Success: false, Entity: "Tool", Message: fmt.Sprintf("Tool '%s' tidak dikenali.", call.Tool)}
	}
}

// buildToolResultBlock renders a toolResult into the markdown block that
// replaces the TOOL_CALL marker in the final response.
func buildToolResultBlock(r toolResult) string {
	if !r.Success {
		return fmt.Sprintf("\n\n⚠️ Gagal membuat %s: %s", r.Entity, r.Message)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n\n✅ **%s berhasil dibuat!**", r.Entity))
	if r.Message != "" {
		b.WriteString("\n")
		b.WriteString(r.Message)
	}
	if r.PageURL != "" {
		b.WriteString(fmt.Sprintf(
			"\n<!-- ACTION:{\"type\":\"navigate\",\"label\":\"Lihat %s\",\"description\":\"Buka halaman %s\",\"url\":\"%s\",\"icon\":\"%s\"} -->",
			r.Entity, r.Entity, r.PageURL, r.Icon,
		))
	}
	return b.String()
}

// ----------------------------------------------------------------------------
// Individual Tool Handlers
// ----------------------------------------------------------------------------

func (s *Service) toolCreateTask(params map[string]interface{}, userID string) toolResult {
	if s.taskService == nil {
		return toolResult{Success: false, Entity: "Task", Message: "Task service tidak tersedia."}
	}

	title := paramStr(params, "title")
	if title == "" {
		return toolResult{Success: false, Entity: "Task", Message: "Judul task wajib diisi."}
	}

	req := &taskdomain.CreateTaskRequest{
		Title:       title,
		Description: paramStr(params, "description"),
		Priority:    paramStrOr(params, "priority", "medium"),
		Type:        paramStrOr(params, "type", "general"),
		AssignedTo:  userID,
		AccountID:   paramStr(params, "account_id"),
		ContactID:   paramStr(params, "contact_id"),
		DealID:      paramStr(params, "deal_id"),
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

func (s *Service) toolCreateDeal(params map[string]interface{}, userID string, userCtx *domainauth.UserContext) toolResult {
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

	accountID := paramStr(params, "account_id")
	if accountID == "" {
		return toolResult{
			Success: false, Entity: "Deal",
			Message: "Account ID wajib diisi untuk membuat deal. Sebutkan nama akun yang bersangkutan.",
		}
	}
	accountEntity, accountErr := s.accountRepo.FindByID(accountID)
	accountOwner := ""
	if accountErr == nil && accountEntity != nil && accountEntity.AssignedTo != nil {
		accountOwner = *accountEntity.AssignedTo
	}
	if accountErr != nil || accountEntity == nil || !s.canAccessOwner(userCtx, "account", accountOwner) {
		return toolResult{Success: false, Entity: "Deal", Message: "Anda tidak memiliki akses ke account yang dipilih untuk membuat deal."}
	}

	// Resolve stage: use provided stage_id or pick the stage with the lowest order.
	stageID := paramStr(params, "stage_id")
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
		Title:     title,
		AccountID: accountID,
		StageID:   stageID,
		ContactID: paramStr(params, "contact_id"),
		Notes:     paramStr(params, "notes"),
	}
	if v, ok := params["value"]; ok {
		switch val := v.(type) {
		case float64:
			req.Value = int64(val)
		case int64:
			req.Value = val
		}
	}

	resp, err := s.pipelineService.CreateDeal(req, userID)
	if err != nil {
		return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Deal", ID: resp.ID,
		Message: fmt.Sprintf("**%s**", resp.Title),
		PageURL: "/pipeline", Icon: "trending-up",
	}
}

func (s *Service) toolCreateSchedule(params map[string]interface{}, userID string) toolResult {
	if s.scheduleService == nil {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Schedule service tidak tersedia."}
	}

	title := paramStr(params, "title")
	if title == "" {
		return toolResult{Success: false, Entity: "Jadwal", Message: "Judul jadwal wajib diisi."}
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

func (s *Service) toolUpdateTaskStatus(params map[string]interface{}, userCtx *domainauth.UserContext) toolResult {
	if s.taskService == nil {
		return toolResult{Success: false, Entity: "Task", Message: "Task service tidak tersedia."}
	}
	id := paramStr(params, "id")
	status := paramStr(params, "status")
	if id == "" || status == "" {
		return toolResult{Success: false, Entity: "Task", Message: "ID dan status task wajib diisi."}
	}
	taskEntity, err := s.taskRepo.FindByID(id)
	if err != nil || taskEntity == nil {
		return toolResult{Success: false, Entity: "Task", Message: "Task tidak ditemukan atau tidak dapat diakses."}
	}
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
		Success: true, Entity: "Task", ID: resp.ID,
		Message: fmt.Sprintf("**%s** → status: **%s**", resp.Title, resp.Status),
		PageURL: "/tasks", Icon: "clipboard",
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
		Success: true, Entity: "Lead", ID: resp.ID,
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
	resp, err := s.pipelineService.MoveStageWithValidation(id, stageID, userID, "Dipindahkan oleh AI")
	if err != nil {
		return toolResult{Success: false, Entity: "Deal", Message: err.Error()}
	}
	return toolResult{
		Success: true, Entity: "Deal", ID: resp.ID,
		Message: fmt.Sprintf("**%s** stage diperbarui.", resp.Title),
		PageURL: "/pipeline", Icon: "trending-up",
	}
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

	stageName := paramStr(params, "stage_name")
	stages, err := s.pipelineService.ListStages(&pipelinedomain.ListPipelineStagesRequest{})
	if err != nil {
		return "", err
	}
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Order < stages[j].Order
	})

	for _, stage := range stages {
		if stageName != "" && strings.EqualFold(stage.Name, stageName) {
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

func paramStringSlice(params map[string]interface{}, key string) []string {
	if v, ok := params[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

// parseFlexibleTime tries several common date/time formats.
func parseFlexibleTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised date format: %s", s)
}
