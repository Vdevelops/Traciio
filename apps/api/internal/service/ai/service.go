package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	groupdomain "github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	productanalyticsdomain "github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	reportdomain "github.com/gilabs/crm-healthcare/api/internal/domain/report"
	route_optimization_domain "github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	scheduledomain "github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	userdomain "github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	activityservice "github.com/gilabs/crm-healthcare/api/internal/service/activity"
	dashboardservice "github.com/gilabs/crm-healthcare/api/internal/service/dashboard"
	leadservice "github.com/gilabs/crm-healthcare/api/internal/service/lead"
	leadqualificationservice "github.com/gilabs/crm-healthcare/api/internal/service/lead_qualification"
	monthlytargetservice "github.com/gilabs/crm-healthcare/api/internal/service/monthly_target"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	productanalyticsservice "github.com/gilabs/crm-healthcare/api/internal/service/product_analytics"
	reportservice "github.com/gilabs/crm-healthcare/api/internal/service/report"
	routeoptimizationservice "github.com/gilabs/crm-healthcare/api/internal/service/route_optimization"
	salesoverviewservice "github.com/gilabs/crm-healthcare/api/internal/service/sales_overview"
	scheduleservice "github.com/gilabs/crm-healthcare/api/internal/service/schedule"
	taskservice "github.com/gilabs/crm-healthcare/api/internal/service/task"
	visitreportservice "github.com/gilabs/crm-healthcare/api/internal/service/visit_report"
	"github.com/gilabs/crm-healthcare/api/pkg/cerebras"
)

var (
	ErrAPINotConfigured = errors.New("Cerebras API key is not configured")
	ErrAPIKeyEmpty      = errors.New("Cerebras API key is empty")
)

// Service represents AI service
type Service struct {
	cerebrasClient           *cerebras.Client
	visitReportRepo          interfaces.VisitReportRepository
	accountRepo              interfaces.AccountRepository
	contactRepo              interfaces.ContactRepository
	dealRepo                 interfaces.DealRepository
	leadRepo                 interfaces.LeadRepository
	leadStatusRepo           interfaces.LeadStatusRepository
	activityRepo             interfaces.ActivityRepository
	taskRepo                 interfaces.TaskRepository
	productRepo              interfaces.ProductRepository
	pipelineRepo             interfaces.PipelineRepository
	userRepo                 interfaces.UserRepository
	roleRepo                 interfaces.RoleRepository
	groupRepo                interfaces.GroupRepository
	brickRepo                interfaces.BrickRepository
	settingsRepo             interfaces.AISettingsRepository
	permService              *permissionservice.Service
	dashboardService         *dashboardservice.Service         // For analytics data
	routeOptimizationService *routeoptimizationservice.Service // For creating real routes
	salesOverviewService     *salesoverviewservice.Service
	productAnalyticsService  *productanalyticsservice.Service
	reportService            *reportservice.Service
	monthlyTargetService     *monthlytargetservice.Service
	// CRUD tool services
	leadService              *leadservice.Service
	activityService          *activityservice.Service
	taskService              *taskservice.Service
	pipelineService          *pipelineservice.Service
	scheduleService          *scheduleservice.Service
	visitReportService       *visitreportservice.Service
	leadQualificationService *leadqualificationservice.Service
	externalIntelligence     *externalIntelligenceService
	apiKey                   string
}

// NewService creates a new AI service
func NewService(
	cerebrasClient *cerebras.Client,
	visitReportRepo interfaces.VisitReportRepository,
	accountRepo interfaces.AccountRepository,
	contactRepo interfaces.ContactRepository,
	dealRepo interfaces.DealRepository,
	leadRepo interfaces.LeadRepository,
	leadStatusRepo interfaces.LeadStatusRepository,
	activityRepo interfaces.ActivityRepository,
	taskRepo interfaces.TaskRepository,
	productRepo interfaces.ProductRepository,
	pipelineRepo interfaces.PipelineRepository,
	userRepo interfaces.UserRepository,
	roleRepo interfaces.RoleRepository,
	groupRepo interfaces.GroupRepository,
	brickRepo interfaces.BrickRepository,
	settingsRepo interfaces.AISettingsRepository,
	permService *permissionservice.Service,
	dashboardService *dashboardservice.Service,
	routeOptimizationService *routeoptimizationservice.Service,
	salesOverviewService *salesoverviewservice.Service,
	productAnalyticsService *productanalyticsservice.Service,
	reportService *reportservice.Service,
	monthlyTargetService *monthlytargetservice.Service,
	leadSvc *leadservice.Service,
	activitySvc *activityservice.Service,
	taskSvc *taskservice.Service,
	pipelineSvc *pipelineservice.Service,
	scheduleSvc *scheduleservice.Service,
	visitReportSvc *visitreportservice.Service,
	leadQualificationSvc *leadqualificationservice.Service,
	apiKey string,
) *Service {
	return &Service{
		cerebrasClient:           cerebrasClient,
		visitReportRepo:          visitReportRepo,
		accountRepo:              accountRepo,
		contactRepo:              contactRepo,
		dealRepo:                 dealRepo,
		leadRepo:                 leadRepo,
		leadStatusRepo:           leadStatusRepo,
		activityRepo:             activityRepo,
		taskRepo:                 taskRepo,
		productRepo:              productRepo,
		pipelineRepo:             pipelineRepo,
		userRepo:                 userRepo,
		roleRepo:                 roleRepo,
		groupRepo:                groupRepo,
		brickRepo:                brickRepo,
		settingsRepo:             settingsRepo,
		permService:              permService,
		dashboardService:         dashboardService,
		routeOptimizationService: routeOptimizationService,
		salesOverviewService:     salesOverviewService,
		productAnalyticsService:  productAnalyticsService,
		reportService:            reportService,
		monthlyTargetService:     monthlyTargetService,
		leadService:              leadSvc,
		activityService:          activitySvc,
		taskService:              taskSvc,
		pipelineService:          pipelineSvc,
		scheduleService:          scheduleSvc,
		visitReportService:       visitReportSvc,
		leadQualificationService: leadQualificationSvc,
		externalIntelligence:     newExternalIntelligenceServiceFromEnv(),
		apiKey:                   apiKey,
	}
}

// validateAPIKey checks if API key is configured
func (s *Service) validateAPIKey() error {
	if s.apiKey == "" {
		return ErrAPIKeyEmpty
	}
	return nil
}

// AnalyzeVisitReport analyzes visit report and returns AI insights
func (s *Service) AnalyzeVisitReport(visitReportID string, userID string, userCtx *domainauth.UserContext) (*ai.VisitReportInsight, int, error) {
	userCtx = s.ensureUserContext(userID, userCtx)
	allowed, err := s.checkDataPrivacy("visit_report", userID, userCtx)
	if err != nil {
		return nil, 0, err
	}
	if !allowed {
		return nil, 0, fmt.Errorf("you do not have permission to access visit report data")
	}

	// Get visit report
	visitReport, err := s.visitReportRepo.FindByID(visitReportID)
	if err != nil {
		return nil, 0, fmt.Errorf("visit report not found: %w", err)
	}
	if !s.canAccessOwner(userCtx, "visit_report", visitReport.SalesRepID) {
		return nil, 0, fmt.Errorf("you do not have permission to access this visit report")
	}

	// Get account (if AccountID is provided)
	var accountEntity *account.Account
	if visitReport.AccountID != nil && *visitReport.AccountID != "" {
		acc, err := s.accountRepo.FindByID(*visitReport.AccountID)
		if err != nil {
			return nil, 0, fmt.Errorf("account not found: %w", err)
		}
		accountAllowed := s.scopedUserIDs(userCtx, "account") == nil
		if acc.AssignedTo != nil {
			accountAllowed = s.canAccessOwner(userCtx, "account", *acc.AssignedTo)
		}
		if accountAllowed {
			accountEntity = acc
		}
	}

	// Get contact if exists
	var contactName string
	if visitReport.ContactID != nil && accountEntity != nil {
		contactEntity, err := s.contactRepo.FindByID(*visitReport.ContactID)
		if err == nil && contactEntity.AccountID == accountEntity.ID {
			contactName = contactEntity.Name
		}
	}

	// Get recent activities for context (if AccountID is provided)
	var activities []activity.Activity
	if visitReport.AccountID != nil && *visitReport.AccountID != "" && accountEntity != nil {
		activities, _ = s.activityRepo.FindByAccountID(*visitReport.AccountID)
	}
	// Limit to 5 most recent
	if len(activities) > 5 {
		activities = activities[:5]
	}

	// Validate API key
	if err := s.validateAPIKey(); err != nil {
		return nil, 0, fmt.Errorf("AI service not configured: %w", err)
	}

	// Build context for AI
	context := BuildVisitReportContext(visitReport, accountEntity, contactName, activities)

	// Build prompt
	prompt := BuildVisitReportPrompt(context)

	// Call Cerebras API
	response, err := s.cerebrasClient.Generate(&cerebras.GenerateRequest{
		Prompt:      prompt,
		MaxTokens:   800,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate insight: %w", err)
	}

	// Parse AI response
	insight, err := s.parseVisitReportInsight(response.Text)
	if err != nil {
		// If parsing fails, return raw response as summary
		insight = &ai.VisitReportInsight{
			Summary:         response.Text,
			ActionItems:     []string{},
			Sentiment:       "neutral",
			KeyPoints:       []string{},
			Recommendations: []string{},
		}
	}

	return insight, response.Tokens, nil
}

// checkDataPrivacy checks if data type is allowed based on settings AND user permissions
// First checks data privacy settings (global), then checks user's role-based permissions
func (s *Service) checkDataPrivacy(dataType string, userID string, userCtx *domainauth.UserContext) (bool, error) {
	userCtx = s.ensureUserContext(userID, userCtx)

	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return s.hasDataPermission(dataType, userCtx), nil
	}

	if !settings.Enabled {
		return false, fmt.Errorf("AI service is disabled")
	}

	// Parse data privacy settings
	dataPrivacy := ai_settings.DefaultDataPrivacySettings()
	if settings.DataPrivacy != nil {
		if err := json.Unmarshal(settings.DataPrivacy, &dataPrivacy); err != nil {
			return true, nil // Default to allow if parsing fails
		}
	}

	// Check data privacy setting first (global setting)
	var privacyAllowed bool
	switch dataType {
	case "visit_report":
		privacyAllowed = dataPrivacy.AllowVisitReports
	case "account":
		privacyAllowed = dataPrivacy.AllowAccounts
	case "contact":
		privacyAllowed = dataPrivacy.AllowContacts
	case "deal":
		privacyAllowed = dataPrivacy.AllowDeals
	case "lead":
		privacyAllowed = dataPrivacy.AllowLeads
	case "activity":
		privacyAllowed = dataPrivacy.AllowActivities
	case "task":
		privacyAllowed = dataPrivacy.AllowTasks
	case "product":
		privacyAllowed = dataPrivacy.AllowProducts
	case "pipeline":
		privacyAllowed = dataPrivacy.AllowPipelines
	case "schedule":
		privacyAllowed = dataPrivacy.AllowSchedule
	// Analytics modules
	case "sales_performance":
		privacyAllowed = dataPrivacy.AllowSalesPerformance
	case "product_analysis":
		privacyAllowed = dataPrivacy.AllowProductAnalysis
	case "report":
		privacyAllowed = dataPrivacy.AllowReports
	// Management modules
	case "user":
		privacyAllowed = dataPrivacy.AllowUsers
	case "role":
		privacyAllowed = dataPrivacy.AllowRoles
	case "group":
		privacyAllowed = dataPrivacy.AllowGroups
	case "brick_management":
		privacyAllowed = dataPrivacy.AllowBrickManagement
	case "target":
		privacyAllowed = dataPrivacy.AllowTarget
	// Route Optimization
	case "route_optimization":
		privacyAllowed = dataPrivacy.AllowRouteOptimization
	default:
		privacyAllowed = true // Default to allow for unknown types
	}

	// If data privacy setting disallows, return false immediately
	if !privacyAllowed {
		return false, nil
	}

	return s.hasDataPermission(dataType, userCtx), nil
}

func (s *Service) buildPersonalAccessInfo(settings *ai_settings.AISettings, userCtx *domainauth.UserContext) string {
	if userCtx == nil {
		return "Konteks permission user tidak tersedia. AI tidak dapat memastikan akses data Anda."
	}

	dataPrivacy := ai_settings.DefaultDataPrivacySettings()
	if settings != nil && settings.DataPrivacy != nil {
		_ = json.Unmarshal(settings.DataPrivacy, &dataPrivacy)
	}

	labels := map[string]string{
		"visit_report":       "Visit Reports",
		"account":            "Accounts",
		"contact":            "Contacts",
		"deal":               "Pipeline/Deals",
		"lead":               "Leads",
		"activity":           "Activities",
		"task":               "Tasks",
		"product":            "Products",
		"pipeline":           "Pipeline",
		"schedule":           "Schedules",
		"sales_performance":  "Sales Performance",
		"product_analysis":   "Product Analytics",
		"report":             "Reports",
		"user":               "Users",
		"role":               "Roles",
		"group":              "Groups",
		"brick_management":   "Bricks/Territories",
		"target":             "Monthly Targets",
		"route_optimization": "Route Optimization",
	}

	allowedData := make([]string, 0)
	blockedData := make([]string, 0)
	for dataType, rule := range aiDataAccessRules {
		label := labels[dataType]
		if label == "" {
			label = dataType
		}
		if !aiDataPrivacyAllows(dataType, dataPrivacy) {
			blockedData = append(blockedData, fmt.Sprintf("%s (disabled by AI privacy)", label))
			continue
		}
		if !s.hasDataPermission(dataType, userCtx) {
			blockedData = append(blockedData, fmt.Sprintf("%s (missing permission: %s)", label, strings.Join(rule.Permissions, " or ")))
			continue
		}
		allowedData = append(allowedData, fmt.Sprintf("%s (scope: %s)", label, userCtx.GetScope(rule.Resource)))
	}
	sort.Strings(allowedData)
	sort.Strings(blockedData)

	allowedTools := make([]string, 0)
	for tool := range aiToolPermissions {
		if s.canRunTool(tool, userCtx) {
			allowedTools = append(allowedTools, tool)
		}
	}
	sort.Strings(allowedTools)

	if len(allowedData) == 0 {
		allowedData = append(allowedData, "Tidak ada modul data yang bisa diakses.")
	}
	if len(allowedTools) == 0 {
		allowedTools = append(allowedTools, "Tidak ada action tool yang bisa dijalankan.")
	}

	var sb strings.Builder
	sb.WriteString("### Akses AI Anda\n\n")
	sb.WriteString("AI akan menjawab dan mengambil data hanya berdasarkan permission serta scope user login.\n\n")
	sb.WriteString("**Data yang bisa diakses:**\n")
	for _, item := range allowedData {
		sb.WriteString("- ")
		sb.WriteString(item)
		sb.WriteString("\n")
	}
	sb.WriteString("\n**Action tool yang bisa dijalankan:**\n")
	for _, item := range allowedTools {
		sb.WriteString("- ")
		sb.WriteString(item)
		sb.WriteString("\n")
	}
	if len(blockedData) > 0 {
		sb.WriteString("\n**Data yang tidak bisa diakses saat ini:**\n")
		limit := min(len(blockedData), 8)
		for _, item := range blockedData[:limit] {
			sb.WriteString("- ")
			sb.WriteString(item)
			sb.WriteString("\n")
		}
		if len(blockedData) > limit {
			sb.WriteString(fmt.Sprintf("- %d modul lainnya tidak tersedia.\n", len(blockedData)-limit))
		}
	}
	return strings.TrimSpace(sb.String())
}

func aiDataPrivacyAllows(dataType string, dataPrivacy ai_settings.DataPrivacySettings) bool {
	switch dataType {
	case "visit_report":
		return dataPrivacy.AllowVisitReports
	case "account":
		return dataPrivacy.AllowAccounts
	case "contact":
		return dataPrivacy.AllowContacts
	case "deal":
		return dataPrivacy.AllowDeals
	case "lead":
		return dataPrivacy.AllowLeads
	case "activity":
		return dataPrivacy.AllowActivities
	case "task":
		return dataPrivacy.AllowTasks
	case "product":
		return dataPrivacy.AllowProducts
	case "pipeline":
		return dataPrivacy.AllowPipelines
	case "schedule":
		return dataPrivacy.AllowSchedule
	case "sales_performance":
		return dataPrivacy.AllowSalesPerformance
	case "product_analysis":
		return dataPrivacy.AllowProductAnalysis
	case "report":
		return dataPrivacy.AllowReports
	case "user":
		return dataPrivacy.AllowUsers
	case "role":
		return dataPrivacy.AllowRoles
	case "group":
		return dataPrivacy.AllowGroups
	case "brick_management":
		return dataPrivacy.AllowBrickManagement
	case "target":
		return dataPrivacy.AllowTarget
	case "route_optimization":
		return dataPrivacy.AllowRouteOptimization
	default:
		return false
	}
}

// Chat handles chat conversation with AI
// userID is required to check user permissions for data access
// domain is an optional hint from the frontend about which CRM domain the user is interacting with
func (s *Service) Chat(message string, contextID string, contextType string, conversationHistory []ai.ChatMessage, model string, userID string, domain string, userCtx *domainauth.UserContext) (*ai.ChatResponse, error) {
	userCtx = s.ensureUserContext(userID, userCtx)

	// Get AI settings
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get AI settings: %w", err)
	}

	if !settings.Enabled {
		return nil, fmt.Errorf("AI service is disabled")
	}

	// Use model from request or settings
	selectedModel := model
	if selectedModel == "" {
		selectedModel = settings.Model
	}

	// Get API key from settings or fallback to env
	apiKey := settings.APIKey
	if apiKey == "" {
		apiKey = s.apiKey // Fallback to env
	}

	if apiKey == "" {
		return nil, fmt.Errorf("AI service not configured: API key is empty")
	}

	// Handle specific query about data privacy settings
	messageLower := strings.ToLower(message)
	if response, handled := s.tryHandlePendingDealAccountConfirmation(message, conversationHistory, userID, userCtx); handled {
		return response, nil
	}
	if response, handled := s.tryHandleCreateActivity(message, conversationHistory, userID, userCtx); handled {
		return response, nil
	}
	if response, handled := s.tryHandleUpdateDealStageWithProducts(message, conversationHistory, userID, userCtx); handled {
		return response, nil
	}
	if strings.Contains(messageLower, "data privacy") || strings.Contains(messageLower, "privacy") ||
		strings.Contains(messageLower, "data privasi") || strings.Contains(messageLower, "privasi") ||
		strings.Contains(messageLower, "akses data") || strings.Contains(messageLower, "data yang bisa") {
		personalAccess := s.buildPersonalAccessInfo(settings, userCtx)

		// Get data privacy settings
		var dataPrivacy ai_settings.DataPrivacySettings

		if settings.DataPrivacy != nil && s.hasAnyPermission(userCtx, "ai-settings.view") {
			if err := json.Unmarshal(settings.DataPrivacy, &dataPrivacy); err == nil {
				privacyInfo := buildPrivacyInfoMessage(dataPrivacy, false)

				// Return direct response about data privacy settings
				return &ai.ChatResponse{
					Message: privacyInfo + "\n\n" + personalAccess + "\n\nIni adalah pengaturan data privacy dan modul analytics yang aktif di sistem. Perubahan setting tetap memerlukan permission AI Settings.",
					Tokens:  0, // No tokens consumed for this internal response
				}, nil
			}
		}

		if !s.hasAnyPermission(userCtx, "ai-settings.view") {
			return &ai.ChatResponse{
				Message: personalAccess,
				Tokens:  0,
			}, nil
		}

		// Use default settings (all allowed) if DataPrivacy is nil
		defaultPrivacy := ai_settings.DataPrivacySettings{
			AllowVisitReports:      true,
			AllowAccounts:          true,
			AllowContacts:          true,
			AllowDeals:             true,
			AllowLeads:             true,
			AllowActivities:        true,
			AllowTasks:             true,
			AllowProducts:          true,
			AllowSalesPerformance:  true,
			AllowBrickManagement:   true,
			AllowProductAnalysis:   true,
			AllowGroups:            true,
			AllowTarget:            true,
			AllowPipelines:         true,
			AllowReports:           true,
			AllowUsers:             true,
			AllowRoles:             true,
			AllowRouteOptimization: true,
		}
		privacyInfo := buildPrivacyInfoMessage(defaultPrivacy, true)

		return &ai.ChatResponse{
			Message: privacyInfo + "\n\n" + personalAccess + "\n\nAnda dapat mengatur data privacy dan enable/disable modul analytics melalui halaman AI Settings jika memiliki permission edit.",
			Tokens:  0,
		}, nil
	}

	// Load context data if provided
	var contextData string
	var dataAccessInfo string
	var queryPlan aiQueryPlan

	if contextID != "" && contextType != "" {
		// Load specific context data
		switch contextType {
		case "visit_report":
			allowed, _ := s.checkDataPrivacy("visit_report", userID, userCtx)
			visitReport, err := s.visitReportRepo.FindByID(contextID)
			if allowed && err == nil && s.canAccessOwner(userCtx, "visit_report", visitReport.SalesRepID) {
				visitReportJSON, _ := json.Marshal(visitReport)
				contextData = string(visitReportJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data visit report dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		case "deal":
			allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
			deal, err := s.dealRepo.FindByID(contextID)
			dealOwner := ""
			if err == nil && deal.AssignedTo != nil {
				dealOwner = *deal.AssignedTo
			}
			if allowed && err == nil && s.canAccessOwner(userCtx, "deal", dealOwner) {
				dealJSON, _ := json.Marshal(deal)
				contextData = string(dealJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data deal dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		case "contact":
			allowed, _ := s.checkDataPrivacy("contact", userID, userCtx)
			contactEntity, err := s.contactRepo.FindByID(contextID)
			contactOwner := ""
			if err == nil && contactEntity != nil {
				if acc, accErr := s.accountRepo.FindByID(contactEntity.AccountID); accErr == nil && acc != nil && acc.AssignedTo != nil {
					contactOwner = *acc.AssignedTo
				}
			}
			if allowed && err == nil && s.canAccessOwner(userCtx, "contact", contactOwner) {
				contactJSON, _ := json.Marshal(contactEntity)
				contextData = string(contactJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data contact dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		case "account":
			allowed, _ := s.checkDataPrivacy("account", userID, userCtx)
			accountEntity, err := s.accountRepo.FindByID(contextID)
			accountOwner := ""
			if err == nil && accountEntity.AssignedTo != nil {
				accountOwner = *accountEntity.AssignedTo
			}
			if allowed && err == nil && s.canAccessOwner(userCtx, "account", accountOwner) {
				accountJSON, _ := json.Marshal(accountEntity)
				contextData = string(accountJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data account dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		case "lead":
			allowed, _ := s.checkDataPrivacy("lead", userID, userCtx)
			leadEntity, err := s.leadRepo.FindByID(contextID)
			leadOwner := ""
			if err == nil && leadEntity.AssignedTo != nil {
				leadOwner = *leadEntity.AssignedTo
			}
			if allowed && err == nil && s.canAccessOwner(userCtx, "lead", leadOwner) {
				leadJSON, _ := json.Marshal(leadEntity)
				contextData = string(leadJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data lead dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		case "product":
			allowed, _ := s.checkDataPrivacy("product", userID, userCtx)
			var productEntity *product.Product
			var err error
			if s.productRepo != nil {
				productEntity, err = s.productRepo.FindByID(contextID)
			}
			if allowed && err == nil && productEntity != nil {
				productJSON, _ := json.Marshal(productEntity)
				contextData = string(productJSON)
			} else {
				dataAccessInfo = "⚠️ Tidak dapat mengakses data product dengan ID tersebut. Data mungkin tidak ditemukan atau tidak memiliki akses."
			}
		}
	} else {
		// Try to extract data from user message - ALWAYS try to get data
		messageLower := strings.ToLower(message)
		queryPlan = s.planAIQuery(message, domain, time.Now())
		if queryPlan.hasDataTypes() && !queryPlan.HandledByLegacyFlow {
			retrievedContext := s.retrieveAIContext(queryPlan, message, userID, userCtx, time.Now())
			if retrievedContext.hasData() {
				contextData = composeAIContext(queryPlan, retrievedContext)
				contextType = retrievedContext.ContextType
				if domain == "" || domain == "auto" {
					domain = retrievedContext.Domain
				}
			} else if retrievedContext.AccessInfo != "" {
				return &ai.ChatResponse{
					Message: noDataAIMessage(queryPlan, retrievedContext.AccessInfo),
					Tokens:  0,
				}, nil
			}
		}

		// ── CRUD INTENT DETECTION (HIGHEST PRIORITY) ──────────────────────────
		// When the user explicitly asks to create/update an entity, we must
		// enrich the context with relevant entity data so the LLM can emit a
		// proper TOOL_CALL with real IDs instead of asking the user for them.
		isDraftProposalIntent := isProposalDraftIntent(messageLower)
		if isDraftProposalIntent && contextData == "" {
			contextData = buildProposalDraftContext(conversationHistory)
			contextType = "proposal_draft"
		}

		isCRUDIntent := strings.Contains(messageLower, "buat") ||
			strings.Contains(messageLower, "buatkan") ||
			strings.Contains(messageLower, "tambah") ||
			strings.Contains(messageLower, "create") ||
			strings.Contains(messageLower, "add") ||
			strings.Contains(messageLower, "log ") ||
			strings.Contains(messageLower, "catat") ||
			strings.Contains(messageLower, "update") ||
			strings.Contains(messageLower, "ubah") ||
			strings.Contains(messageLower, "bant") ||
			strings.Contains(messageLower, "qualification") ||
			strings.Contains(messageLower, "kualifikasi") ||
			strings.Contains(messageLower, "pindah") ||
			strings.Contains(messageLower, "jadwalkan") ||
			strings.Contains(messageLower, "ya") && (strings.Contains(messageLower, "buat") || strings.Contains(messageLower, "follow")) ||
			strings.Contains(messageLower, "follow up") ||
			strings.Contains(messageLower, "follow-up") ||
			strings.Contains(messageLower, "followup")

		if isCRUDIntent && !isDraftProposalIntent && contextData == "" {
			// Extract entity names/references from the current message and
			// conversation history so we can load their real IDs from the DB.
			crudContext := s.buildCRUDContext(messageLower, conversationHistory, userID, userCtx)
			if crudContext != "" {
				contextData = crudContext
				contextType = "crud_context"
			}
		}

		// Check for ROUTE OPTIMIZATION queries (HIGHEST PRIORITY)
		isRouteQuery := contextData == "" && (strings.Contains(messageLower, "rute") ||
			strings.Contains(messageLower, "route") ||
			strings.Contains(messageLower, "optimasi") ||
			strings.Contains(messageLower, "optimize") ||
			strings.Contains(messageLower, "navigasi") ||
			strings.Contains(messageLower, "navigation") ||
			strings.Contains(messageLower, "kunjungan optimal") ||
			strings.Contains(messageLower, "visit route") ||
			strings.Contains(messageLower, "jalur") ||
			strings.Contains(messageLower, "perjalanan"))

		// Detect explicit CREATE intent in current message
		isCreateRouteIntent := isRouteQuery && (strings.Contains(messageLower, "buat") ||
			strings.Contains(messageLower, "create") ||
			strings.Contains(messageLower, "mulai") ||
			strings.Contains(messageLower, "sekarang") ||
			strings.Contains(messageLower, "langsung") ||
			strings.Contains(messageLower, "implementasi") ||
			strings.Contains(messageLower, "jalankan") ||
			strings.Contains(messageLower, "eksekusi") ||
			strings.Contains(messageLower, "random"))

		if isCreateRouteIntent && s.routeOptimizationService != nil {
			// Try to extract starting location from conversation history
			startLat, startLng, hasLocation := extractLocationFromHistory(conversationHistory)
			// Extract account IDs from conversation history (previous AI list responses)
			accountIDs := extractAccountIDsFromHistory(conversationHistory)

			if hasLocation && len(accountIDs) > 0 {
				allowed, _ := s.checkDataPrivacy("route_optimization", userID, userCtx)
				if allowed {
					// Build waypoints from the real accounts found in history
					waypoints := s.buildWaypointsFromAccountIDs(accountIDs, userCtx)
					if len(waypoints) > 0 {
						req := &route_optimization_domain.OptimizeRouteRequest{
							StartLocation: &route_optimization_domain.Location{
								Lat: startLat,
								Lng: startLng,
							},
							Waypoints: waypoints,
						}
						result, err := s.routeOptimizationService.Optimize(req, userID)
						if err == nil && result != nil {
							// Return a structured response with the real created route
							resultJSON, _ := json.Marshal(result)
							contextData = fmt.Sprintf("ROUTE_CREATED_SUCCESSFULLY:\n%s", string(resultJSON))
							contextType = "route_created"
						} else {
							dataAccessInfo = fmt.Sprintf("⚠️ Gagal membuat rute: %v. Silakan coba melalui halaman Route Optimization.", err)
						}
					} else {
						// Accounts found in history but none have GPS coordinates.
						// Set dataAccessInfo so the LLM explains this clearly and we
						// do NOT fall through to the account list re-fetch (which would
						// cause the LLM to loop and emit another failing TOOL_CALL).
						dataAccessInfo = fmt.Sprintf(
							"INFO_NO_GPS_COORDINATES: Ditemukan %d akun dari percakapan, namun semua akun tersebut belum memiliki data koordinat GPS (Latitude/Longitude). "+
								"Rute otomatis tidak dapat dibuat tanpa koordinat. "+
								"Sampaikan kepada pengguna bahwa mereka perlu: (1) menambahkan koordinat GPS di halaman Accounts untuk masing-masing akun, "+
								"ATAU (2) membuat rute secara manual di halaman Route Optimization dengan memilih akun dan memasukkan koordinat secara langsung. "+
								"Sertakan action card untuk /accounts dan /route-optimization.",
							len(accountIDs),
						)
					}
				}
			}
		}

		if isRouteQuery && contextData == "" && dataAccessInfo == "" {
			allowed, _ := s.checkDataPrivacy("route_optimization", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Modul Route Optimization tidak diaktifkan di pengaturan AI atau Anda tidak memiliki akses."
			} else {
				// Fetch real accounts with addresses for route planning
				accounts, total, err := s.accountRepo.List(&account.ListAccountsRequest{
					Page:          1,
					PerPage:       20,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "account"),
				})
				if err == nil && len(accounts) > 0 {
					accountsFormatted := s.formatAccountsForAI(accounts)
					accountsJSON, _ := json.Marshal(accountsFormatted)
					contextData = fmt.Sprintf("REAL ACCOUNTS DATA FOR ROUTE OPTIMIZATION (%d of %d total accounts):\n%s\n\nIMPORTANT: Use ONLY these real accounts to suggest routes. Do NOT invent addresses, distances, or travel times. If address data is incomplete, inform the user.", len(accounts), total, string(accountsJSON))
					contextType = "account"
				} else {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data akun untuk optimasi rute. Data mungkin tidak tersedia."
				}
			}
		}

		isWonLostDealsQuery := contextData == "" &&
			(strings.Contains(messageLower, "deal") || strings.Contains(messageLower, "deals") || strings.Contains(messageLower, "pipeline")) &&
			(strings.Contains(messageLower, "won") || strings.Contains(messageLower, "lost") || strings.Contains(messageLower, "win"))

		if isWonLostDealsQuery {
			if wonLostContext, wonLostInfo := s.buildWonLostDealsContext(messageLower, userID, userCtx); wonLostContext != "" {
				contextData = wonLostContext
				contextType = "deal"
			} else if wonLostInfo != "" {
				dataAccessInfo = wonLostInfo
			}
		}

		// Check for SALES PERFORMANCE queries
		isSalesPerformanceQuery := contextData == "" && (strings.Contains(messageLower, "sales performance") ||
			strings.Contains(messageLower, "performa penjualan") ||
			strings.Contains(messageLower, "performa sales") ||
			strings.Contains(messageLower, "target vs actual") ||
			strings.Contains(messageLower, "target pencapaian") ||
			strings.Contains(messageLower, "quota") ||
			(strings.Contains(messageLower, "report") && strings.Contains(messageLower, "sales")) ||
			(strings.Contains(messageLower, "laporan") && strings.Contains(messageLower, "sales")) ||
			(strings.Contains(messageLower, "report") && strings.Contains(messageLower, "penjualan")) ||
			(strings.Contains(messageLower, "laporan") && strings.Contains(messageLower, "penjualan")) ||
			(strings.Contains(messageLower, "performa") && strings.Contains(messageLower, "brick")) ||
			(strings.Contains(messageLower, "performance") && strings.Contains(messageLower, "brick")) ||
			(strings.Contains(messageLower, "penjualan") && hasTrendOrChartTerm(messageLower)) ||
			(strings.Contains(messageLower, "revenue") && strings.Contains(messageLower, "target")))

		if isSalesPerformanceQuery {
			if performanceContext, performanceInfo := s.buildPerformanceContext(messageLower, userID, userCtx); performanceContext != "" {
				contextData = performanceContext
				contextType = "sales_performance"
			} else if performanceInfo != "" {
				dataAccessInfo = performanceInfo
			}
		}

		isTargetQuery := contextData == "" && !isDealValueTargetIntent(messageLower) && (strings.Contains(messageLower, "target") ||
			strings.Contains(messageLower, "quota") ||
			strings.Contains(messageLower, "kuota"))

		if isTargetQuery {
			if targetContext, targetInfo := s.buildTargetContext(messageLower, userID, userCtx); targetContext != "" {
				contextData = targetContext
				contextType = "target"
			} else if targetInfo != "" {
				dataAccessInfo = targetInfo
			}
		}

		isReportQuery := contextData == "" && (strings.Contains(messageLower, "laporan") ||
			strings.Contains(messageLower, "report generator") ||
			strings.Contains(messageLower, "buat report") ||
			strings.Contains(messageLower, "generate report") ||
			strings.Contains(messageLower, "ringkasan report") ||
			strings.Contains(messageLower, "summary report") ||
			strings.Contains(messageLower, "rekap report"))

		if isReportQuery {
			if reportContext, reportInfo := s.buildReportContext(messageLower, userID, userCtx); reportContext != "" {
				contextData = reportContext
				contextType = "report"
			} else if reportInfo != "" {
				dataAccessInfo = reportInfo
			}
		}

		isUserManagementQuery := contextData == "" && (strings.Contains(messageLower, "user") ||
			strings.Contains(messageLower, "pengguna") ||
			strings.Contains(messageLower, "sales rep") ||
			strings.Contains(messageLower, "sales representative") ||
			strings.Contains(messageLower, "admin"))

		if isUserManagementQuery {
			if userContext, userInfo := s.buildUserManagementContext(messageLower, userID, userCtx); userContext != "" {
				contextData = userContext
				contextType = "user"
			} else if userInfo != "" {
				dataAccessInfo = userInfo
			}
		}

		isRoleManagementQuery := contextData == "" && (strings.Contains(messageLower, "role") ||
			strings.Contains(messageLower, "roles") ||
			strings.Contains(messageLower, "peran"))

		if isRoleManagementQuery {
			if roleContext, roleInfo := s.buildRoleManagementContext(messageLower, userID, userCtx); roleContext != "" {
				contextData = roleContext
				contextType = "role"
			} else if roleInfo != "" {
				dataAccessInfo = roleInfo
			}
		}

		isGroupManagementQuery := contextData == "" && (strings.Contains(messageLower, "group") ||
			strings.Contains(messageLower, "groups") ||
			strings.Contains(messageLower, "grup"))

		if isGroupManagementQuery {
			if groupContext, groupInfo := s.buildGroupManagementContext(messageLower, userID, userCtx); groupContext != "" {
				contextData = groupContext
				contextType = "group"
			} else if groupInfo != "" {
				dataAccessInfo = groupInfo
			}
		}

		isBrickManagementQuery := contextData == "" && (strings.Contains(messageLower, "brick") ||
			strings.Contains(messageLower, "bricks") ||
			strings.Contains(messageLower, "territory") ||
			strings.Contains(messageLower, "territories") ||
			strings.Contains(messageLower, "wilayah") ||
			strings.Contains(messageLower, "area") ||
			strings.Contains(messageLower, "area mapping"))

		if isBrickManagementQuery {
			if brickContext, brickInfo := s.buildBrickManagementContext(messageLower, userID, userCtx); brickContext != "" {
				contextData = brickContext
				contextType = "brick_management"
			} else if brickInfo != "" {
				dataAccessInfo = brickInfo
			}
		}

		if contextData == "" && isProspectPredictionIntent(messageLower) {
			if predictionContext, predictionInfo := s.buildProspectPredictionContext(userID, userCtx); predictionContext != "" {
				contextData = predictionContext
				contextType = "prospect_prediction"
			} else if predictionInfo != "" {
				dataAccessInfo = predictionInfo
			}
		}

		// Check for analytics/statistics queries (HIGHEST PRIORITY - needs all deals data)
		// These queries require comprehensive deals data for calculations
		isAnalyticsQuery := contextData == "" && (strings.Contains(messageLower, "conversion rate") ||
			strings.Contains(messageLower, "conversion") ||
			strings.Contains(messageLower, "rate konversi") ||
			strings.Contains(messageLower, "statistik") ||
			strings.Contains(messageLower, "statistics") ||
			strings.Contains(messageLower, "analisis") ||
			strings.Contains(messageLower, "analysis") ||
			strings.Contains(messageLower, "rata-rata") ||
			strings.Contains(messageLower, "average") ||
			strings.Contains(messageLower, "trend") ||
			strings.Contains(messageLower, "perbandingan") ||
			strings.Contains(messageLower, "comparison") ||
			strings.Contains(messageLower, "breakdown") ||
			(strings.Contains(messageLower, "berapa") && (strings.Contains(messageLower, "lead") || strings.Contains(messageLower, "closed won") || strings.Contains(messageLower, "closed lost"))))

		if isAnalyticsQuery {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Akses ke data deals/pipeline tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
			} else {
				// For analytics queries, fetch ALL deals (or at least a large sample) without stage filter
				// This allows AI to calculate conversion rates, averages, etc.
				req := &pipeline.ListDealsRequest{
					Page:          1,
					PerPage:       100, // Fetch more deals for analytics
					ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
				}

				deals, _, err := s.dealRepo.List(req)
				if err == nil && len(deals) > 0 {
					// For analytics, limit to 50 deals max to prevent token overflow
					maxDeals := 50
					if len(deals) > maxDeals {
						deals = deals[:maxDeals]
					}

					// Transform deals to user-friendly format with names
					dealsFormatted := s.formatDealsForAI(deals)
					dealsJSON, _ := json.Marshal(dealsFormatted)

					// Build concise instruction for analytics
					instruction := fmt.Sprintf("REAL DEALS DATA (%d deals for analytics):\n%s\n\nCalculate statistics using ONLY this data. Present results clearly with actual numbers.", len(deals), string(dealsJSON))
					contextData = instruction
					contextType = "deal"
				} else if err == nil {
					dataAccessInfo = "Tidak ada data pipeline/deals sesuai akses dan filter yang diminta."
				} else {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data pipeline/deals dari database untuk perhitungan statistik. Data mungkin tidak tersedia."
				}
			}
		}

		// Check priority: pipeline/deals first (more specific), then accounts
		// Check if user is asking for pipeline/deals/sales funnel (HIGHEST PRIORITY)
		if contextData == "" && (strings.Contains(messageLower, "pipeline") || strings.Contains(messageLower, "sales funnel") ||
			strings.Contains(messageLower, "funnel") || strings.Contains(messageLower, "deal") ||
			strings.Contains(messageLower, "opportunity") || strings.Contains(messageLower, "kesempatan")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Akses ke data deals/pipeline tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
			} else {
				// Build request with optional stage filter
				req := &pipeline.ListDealsRequest{
					Page:          1,
					PerPage:       20,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
				}

				// Extract stage filter from message if mentioned
				var stageID string
				if strings.Contains(messageLower, "lead") {
					if stage, err := s.pipelineRepo.FindStageByCode("lead"); err == nil {
						stageID = stage.ID
					}
				} else if strings.Contains(messageLower, "qualification") {
					if stage, err := s.pipelineRepo.FindStageByCode("qualification"); err == nil {
						stageID = stage.ID
					}
				} else if strings.Contains(messageLower, "proposal") {
					if stage, err := s.pipelineRepo.FindStageByCode("proposal"); err == nil {
						stageID = stage.ID
					}
				} else if strings.Contains(messageLower, "negotiation") {
					if stage, err := s.pipelineRepo.FindStageByCode("negotiation"); err == nil {
						stageID = stage.ID
					}
				} else if strings.Contains(messageLower, "closed won") || strings.Contains(messageLower, "won") {
					if stage, err := s.pipelineRepo.FindStageByCode("closed_won"); err == nil {
						stageID = stage.ID
					}
				} else if strings.Contains(messageLower, "closed lost") || strings.Contains(messageLower, "lost") {
					if stage, err := s.pipelineRepo.FindStageByCode("closed_lost"); err == nil {
						stageID = stage.ID
					}
				}

				if stageID != "" {
					req.StageID = stageID
				}

				deals, _, err := s.dealRepo.List(req)
				if err == nil && len(deals) > 0 {
					// Limit number of deals to prevent token overflow (max 15 deals for large responses)
					maxDeals := 15
					if len(deals) > maxDeals {
						deals = deals[:maxDeals]
					}

					// Transform deals to user-friendly format with names
					dealsFormatted := s.formatDealsForAI(deals)
					dealsJSON, _ := json.Marshal(dealsFormatted)

					// Build concise instruction
					instruction := fmt.Sprintf("REAL PIPELINE/DEALS DATA (showing %d deals):\n%s\n\nPresent in Markdown table. Use [Title](deal://id) for clickable links. Show only names, not IDs.", len(deals), string(dealsJSON))
					contextData = instruction
					contextType = "deal"
					fmt.Printf("Context data set with %d deals (context size: %d chars)\n", len(deals), len(contextData))
				} else if err == nil {
					dataAccessInfo = "Tidak ada data pipeline/deals sesuai akses dan filter yang diminta."
				} else {
					fmt.Printf("Error fetching deals: %v\n", err)
					dataAccessInfo = "⚠️ Tidak dapat mengakses data pipeline/deals dari database. Data mungkin tidak tersedia."
				}
				fmt.Printf("========================\n")
			}
		}

		// Check if user is asking for leads/lead management (HIGH PRIORITY - check before general data)
		// This should be checked early to avoid being caught by "general data" logic
		if contextData == "" && (strings.Contains(messageLower, "lead") || strings.Contains(messageLower, "lead management") ||
			strings.Contains(messageLower, "prospek") || strings.Contains(messageLower, "calon pelanggan") ||
			(strings.Contains(messageLower, "tampilkan") && strings.Contains(messageLower, "lead")) ||
			(strings.Contains(messageLower, "data") && strings.Contains(messageLower, "lead"))) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("lead", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Akses ke data leads tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
			} else {
				// Build request with optional status filter
				req := &lead.ListLeadsRequest{
					Page:          1,
					PerPage:       20, // No limit - AI will handle overflow automatically
					ScopedUserIDs: s.scopedUserIDs(userCtx, "lead"),
				}
				dateRange := parseAIDateRange(messageLower, time.Now())
				start, end, _ := normalizeDateRangeForRequest(dateRange)
				req.StartDate = start
				req.EndDate = end

				// Extract status filter from message if mentioned
				if strings.Contains(messageLower, "new") {
					req.Status = "new"
				} else if strings.Contains(messageLower, "contacted") {
					req.Status = "contacted"
				} else if strings.Contains(messageLower, "proposal sent") || strings.Contains(messageLower, "proposal_sent") {
					req.Status = "proposal_sent"
				} else if strings.Contains(messageLower, "interested") {
					req.Status = "interested"
				} else if strings.Contains(messageLower, "qualified") {
					req.Status = "qualified"
				} else if strings.Contains(messageLower, "converted") {
					req.Status = "converted"
				} else if strings.Contains(messageLower, "lost") {
					req.Status = "lost"
				}

				leads, total, err := s.leadRepo.List(req)
				fmt.Printf("=== DATA FETCH DEBUG ===\n")
				fmt.Printf("Fetching leads - Error: %v, Count: %d, Total: %d, Status: %s, StartDate: %s, EndDate: %s\n", err, len(leads), total, req.Status, req.StartDate, req.EndDate)
				if err == nil && len(leads) > 0 {
					// Limit number of leads to prevent token overflow (max 15 leads for large responses)
					maxLeads := 15
					if len(leads) > maxLeads {
						leads = leads[:maxLeads]
						fmt.Printf("Limited leads to %d to prevent token overflow\n", maxLeads)
					}

					// Transform leads to user-friendly format
					leadsFormatted := s.formatLeadsForAI(leads)
					leadsJSON, _ := json.Marshal(leadsFormatted)

					// Build concise instruction
					instruction := fmt.Sprintf("REAL LEADS DATA (showing %d of %d total):\n%s\n\nPresent in Markdown table. Use [Name](lead://id) for clickable links. Show only names, not IDs.", len(leads), total, string(leadsJSON))
					if wantsLeadFullDetail(messageLower) {
						if details := s.buildLeadFullDetailsForAI(leads, userCtx); details != "" {
							instruction += "\n\n" + details
						}
					}
					contextData = instruction
					contextType = "lead"
				} else if err == nil {
					periodLabel := "periode yang diminta"
					if dateRange.Label != "" {
						periodLabel = dateRange.Label
					}
					statusLabel := req.Status
					if statusLabel == "" {
						statusLabel = "semua status"
					}
					return &ai.ChatResponse{
						Message: fmt.Sprintf("Tidak ada data lead dengan status `%s` pada %s sesuai akses data Anda.", statusLabel, periodLabel),
						Tokens:  0,
					}, nil
				} else {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data leads dari database. Data mungkin tidak tersedia."
				}
			}
		}

		// Check if user is asking for accounts (only if not pipeline)
		if contextData == "" && (strings.Contains(messageLower, "account") || strings.Contains(messageLower, "akun") ||
			strings.Contains(messageLower, "rumah sakit") || strings.Contains(messageLower, "klinik") ||
			strings.Contains(messageLower, "apotek") || strings.Contains(messageLower, "facility")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("account", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Akses ke data accounts tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
			} else {
				accounts, total, err := s.accountRepo.List(&account.ListAccountsRequest{
					Page:          1,
					PerPage:       10,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "account"),
				})
				if err == nil && len(accounts) > 0 {
					// Transform accounts to user-friendly format with names
					accountsFormatted := s.formatAccountsForAI(accounts)
					accountsJSON, _ := json.Marshal(accountsFormatted)
					contextData = fmt.Sprintf("REAL ACCOUNTS DATA FROM DATABASE (showing %d of %d total accounts):\n%s\n\nCRITICAL INSTRUCTION: You MUST use ONLY the data above. Present it in a Markdown table format. CRITICAL: NEVER show IDs as separate columns - IDs are ONLY used in clickable links. ALWAYS show ONLY NAMES (name, category, city, province) in tables. For clickable actions that trigger detail components, the 'Nama Akun' (Name) column MUST be formatted as [Name](account://id) to create clickable links. Example: [RSUD Jakarta](account://ab868b77-e9b3-429f-ad8c-d55ac1f6561b). Use the EXACT id from the data above - DO NOT create or invent IDs. DO NOT create columns like 'ID', 'Account ID', etc. - these should NOT appear in tables. DO NOT create, invent, or make up any data. DO NOT add columns that don't exist in the data. AFTER presenting the table, ALWAYS provide 1-2 insights, ask 2-3 follow-up questions to understand what the user wants, and offer actionable recommendations.", len(accounts), total, string(accountsJSON))
					contextType = "account" // Set context type for proper prompt
				} else if err == nil {
					dataAccessInfo = "Tidak ada data accounts sesuai akses dan filter yang diminta."
				} else {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data accounts dari database. Data mungkin tidak tersedia."
				}
			}
		}

		// Check if user is asking for contacts (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "contact") || strings.Contains(messageLower, "kontak") ||
			strings.Contains(messageLower, "dokter") || strings.Contains(messageLower, "apoteker")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("contact", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data contacts tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else {
				contacts, _, err := s.contactRepo.List(&contact.ListContactsRequest{
					Page:          1,
					PerPage:       10,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "contact"),
				})
				if err == nil && len(contacts) > 0 {
					contactsJSON, _ := json.Marshal(contacts)
					if contextData != "" {
						contextData += "\n\n"
					}
					contextData += fmt.Sprintf("REAL CONTACTS DATA FROM DATABASE (showing %d contacts):\n%s\n\nCRITICAL INSTRUCTION: You MUST use ONLY the data above. Present it in a Markdown table format. CRITICAL: NEVER show IDs as separate columns - IDs are ONLY used in clickable links. ALWAYS show ONLY NAMES (name, email, phone, job_title) in tables. For clickable actions that trigger detail components, the Contact Name column MUST be formatted as [Name](contact://id) to allow clicking and opening contact detail. DO NOT create columns like 'ID', 'Contact ID', etc. - these should NOT appear in tables. DO NOT create, invent, or make up any data. AFTER presenting the table, ALWAYS provide 1-2 insights, ask 2-3 follow-up questions to understand what the user wants, and offer actionable recommendations.", len(contacts), string(contactsJSON))
					if contextType == "" {
						contextType = "contact"
					}
				} else {
					if dataAccessInfo == "" && err == nil {
						dataAccessInfo = "Tidak ada data contacts sesuai akses dan filter yang diminta."
					} else if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data contacts dari database. Data mungkin tidak tersedia."
					}
				}
			}
		}

		// Check if user is asking for visit reports (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "visit") || strings.Contains(messageLower, "kunjungan") ||
			strings.Contains(messageLower, "laporan kunjungan")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("visit_report", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data visit reports tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else {
				// Build request with optional status filter
				req := &visit_report.ListVisitReportsRequest{
					Page:          1,
					PerPage:       10,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "visit_report"),
				}

				// Extract status filter from message if mentioned
				if strings.Contains(messageLower, "approved") {
					req.Status = "approved"
				} else if strings.Contains(messageLower, "submitted") {
					req.Status = "submitted"
				} else if strings.Contains(messageLower, "draft") {
					req.Status = "draft"
				} else if strings.Contains(messageLower, "rejected") {
					req.Status = "rejected"
				}

				visitReports, _, err := s.visitReportRepo.List(req)
				if err == nil && len(visitReports) > 0 {
					// Transform visit reports to user-friendly format with names
					visitReportsFormatted := s.formatVisitReportsForAI(visitReports)
					visitReportsJSON, _ := json.Marshal(visitReportsFormatted)
					if contextData != "" {
						contextData += "\n\n"
					}
					contextData += fmt.Sprintf("REAL VISIT REPORTS DATA FROM DATABASE (showing %d visit reports):\n%s\n\nCRITICAL INSTRUCTION: You MUST use ONLY the data above. Present it in a Markdown table format. CRITICAL: NEVER show IDs as separate columns - IDs are ONLY used in clickable links. ALWAYS show ONLY NAMES (account_name, contact_name, purpose, status) in tables. For clickable actions that trigger detail components, use format [Name](type://ID) where type is 'visit', 'account', or 'contact'. DO NOT create columns like 'ID', 'Visit Report ID', 'Account ID', etc. - these should NOT appear in tables. DO NOT create, invent, or make up any data. DO NOT add columns that don't exist in the data. AFTER presenting the table, ALWAYS provide 1-2 insights, ask 2-3 follow-up questions to understand what the user wants, and offer actionable recommendations.", len(visitReports), string(visitReportsJSON))
					if contextType == "" {
						contextType = "visit_report"
					}
				} else {
					if dataAccessInfo == "" && err == nil {
						dataAccessInfo = "Tidak ada data visit reports sesuai akses dan filter yang diminta."
					} else if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data visit reports dari database. Data mungkin tidak tersedia."
					}
				}
			}
		}

		// Check if user is asking for schedules (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "schedule") || strings.Contains(messageLower, "schedules") ||
			strings.Contains(messageLower, "jadwal")) {
			allowed, _ := s.checkDataPrivacy("schedule", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data schedules tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else if s.scheduleService != nil {
				req := &scheduledomain.ListSchedulesRequest{
					Page:          1,
					PerPage:       20,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"),
				}
				dateRange := parseAIDateRange(messageLower, time.Now())
				start, end, _ := normalizeDateRangeForRequest(dateRange)
				if start != "" {
					if parsed, err := time.Parse(dateFormat, start); err == nil {
						req.ScheduledAtFrom = &parsed
					}
				}
				if end != "" {
					if parsed, err := time.Parse(dateFormat, end); err == nil {
						endOfDay := parsed.Add(24*time.Hour - time.Nanosecond)
						req.ScheduledAtTo = &endOfDay
					}
				}
				if strings.Contains(messageLower, "pending") {
					req.Status = "pending"
				} else if strings.Contains(messageLower, "submitted") {
					req.Status = "submitted"
				} else if strings.Contains(messageLower, "confirmed") {
					req.Status = "confirmed"
				} else if strings.Contains(messageLower, "completed") || strings.Contains(messageLower, "selesai") {
					req.Status = "completed"
				} else if strings.Contains(messageLower, "cancelled") || strings.Contains(messageLower, "batal") {
					req.Status = "cancelled"
				} else if strings.Contains(messageLower, "rejected") {
					req.Status = "rejected"
				}

				schedules, pagination, err := s.scheduleService.ListSchedules(req)
				if err == nil && len(schedules) > 0 {
					schedulesJSON, _ := json.Marshal(schedules)
					total := len(schedules)
					if pagination != nil {
						total = pagination.Total
					}
					contextData = fmt.Sprintf("REAL SCHEDULES DATA (showing %d of %d total schedules):\n%s\n\nPresent in Markdown table. Use [Title](schedule://id) for clickable schedule links. Show only title, scheduled_at, status, and task title. Never show IDs as columns. Do not invent schedule data. Include a navigate action card to /schedules when useful.", len(schedules), total, string(schedulesJSON))
					contextType = "schedule"
				} else if dataAccessInfo == "" && err == nil {
					dataAccessInfo = "Tidak ada data schedules sesuai akses dan filter yang diminta."
				} else if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data schedules dari database. Data mungkin tidak tersedia."
				}
			} else if dataAccessInfo == "" {
				dataAccessInfo = "⚠️ Layanan schedule belum tersedia untuk AI Assistant."
			}
		}

		// Check if user is asking for tasks (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "task") || strings.Contains(messageLower, "tugas")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("task", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data tasks tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else {
				// Build request with optional filters. For update intents, the status
				// mentioned by the user is usually the target status, not a list filter.
				req := &task.ListTasksRequest{
					Page:          1,
					PerPage:       20,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
				}
				if search := taskSearchTermFromMessage(messageLower); search != "" {
					req.Search = search
				}

				if !isTaskStatusUpdateIntent(messageLower) {
					if status := taskStatusFilterFromMessage(messageLower); status != "" {
						req.Status = status
					}
				}

				tasks, _, err := s.taskRepo.List(req)
				if err == nil && len(tasks) > 0 {
					tasksJSON, _ := json.Marshal(tasks)
					if contextData != "" {
						contextData += "\n\n"
					}
					contextData += fmt.Sprintf("REAL TASKS DATA FROM DATABASE (showing %d tasks):\n%s\n\nCRITICAL INSTRUCTION: You MUST use ONLY the data above. Present it in a Markdown table format. CRITICAL: NEVER show IDs as separate columns - IDs are ONLY used in clickable links. ALWAYS show ONLY NAMES (title, account_name, contact_name, status, priority) in tables. For clickable actions that trigger detail components, use format [Name](type://ID) where type is 'task', 'account', or 'contact'. DO NOT create columns like 'ID', 'Task ID', 'Account ID', etc. - these should NOT appear in tables. DO NOT create, invent, or make up any data. AFTER presenting the table, ALWAYS provide 1-2 insights, ask 2-3 follow-up questions to understand what the user wants, and offer actionable recommendations.", len(tasks), string(tasksJSON))
					if contextType == "" {
						contextType = "task"
					}
				} else {
					if dataAccessInfo == "" && err == nil {
						dataAccessInfo = "Tidak ada data tasks sesuai akses dan filter yang diminta."
					} else if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data tasks dari database. Data mungkin tidak tersedia."
					}
				}
			}
		}

		// Check if user is asking for product sales analytics (only if no data fetched yet).
		// Data visibility follows RBAC scope: own for sales, team for sales manager, global for admin.
		if contextData == "" && isProductSalesAnalyticsIntent(messageLower) {
			allowed, _ := s.checkDataPrivacy("product_analysis", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data product analytics tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else if s.productAnalyticsService != nil {
				dateRange := parseAIDateRange(messageLower, time.Now())
				startDate, endDate := aiDateRangeToTimes(dateRange, time.Now())
				scopeKind := s.productSalesScopeKind(userCtx)
				scopedUserIDs := s.scopedUserIDs(userCtx, "product_analysis")
				sortBy := productSalesSortByFromMessage(messageLower)
				var products []*productanalyticsdomain.ProductListItem
				var total int64
				var err error
				if scopeKind == "own" {
					products, total, err = s.productAnalyticsService.GetUserProductSales(userID, startDate, endDate, sortBy, "desc", 1, 20)
				} else {
					products, total, err = s.productAnalyticsService.GetProductsList(startDate, endDate, "", sortBy, "desc", 1, 20, scopedUserIDs)
				}
				if err == nil && len(products) > 0 {
					productsJSON, _ := json.Marshal(products)
					periodLabel := "semua periode"
					if dateRange.HasFilter && dateRange.Label != "" {
						periodLabel = dateRange.Label
					}
					contextData = fmt.Sprintf("REAL PRODUCT SALES DATA FROM DATABASE (scope: %s, period: %s, showing %d of %d sold products, sorted by total_sold desc):\n%s\n\nCRITICAL INSTRUCTION: You MUST answer the user's request using ONLY the data above. Present products in a Markdown table sorted from highest quantity sold to lowest. Show product name, SKU, category, total_sold, total_revenue, total_profit, sales_count, and last_sold_at when present. Never say sales data is unavailable when this context is present. Do not invent products or quantities. Include 1-2 concise insights and a navigate action card to /product-analytics when useful.", scopeKind, periodLabel, len(products), total, string(productsJSON))
					contextType = "product_analysis"
				} else if dataAccessInfo == "" && err == nil {
					periodLabel := "periode yang diminta"
					if dateRange.HasFilter && dateRange.Label != "" {
						periodLabel = dateRange.Label
					}
					return &ai.ChatResponse{
						Message: buildNoProductSalesMessage(periodLabel, scopeKind),
						Tokens:  0,
					}, nil
				} else if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data product analytics dari database. Data mungkin tidak tersedia."
				}
			} else if dataAccessInfo == "" {
				dataAccessInfo = "⚠️ Layanan product analytics belum tersedia untuk AI Assistant."
			}
		}

		// Check if user is asking for products/inventory (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "product") ||
			strings.Contains(messageLower, "produk") ||
			strings.Contains(messageLower, "inventory") ||
			strings.Contains(messageLower, "inventaris") ||
			strings.Contains(messageLower, "obat") ||
			strings.Contains(messageLower, "sku")) {
			allowed, _ := s.checkDataPrivacy("product", userID, userCtx)
			if !allowed {
				if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Akses ke data products tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
				}
			} else if s.productRepo != nil {
				req := &product.ListProductsRequest{
					Page:    1,
					PerPage: 20,
				}
				if strings.Contains(messageLower, "active") || strings.Contains(messageLower, "aktif") {
					req.Status = "active"
				} else if strings.Contains(messageLower, "inactive") || strings.Contains(messageLower, "nonaktif") {
					req.Status = "inactive"
				}

				products, total, err := s.productRepo.List(req)
				if err == nil && len(products) > 0 {
					productsFormatted := s.formatProductsForAI(products)
					productsJSON, _ := json.Marshal(productsFormatted)
					contextData = fmt.Sprintf("REAL PRODUCTS DATA (showing %d of %d total products):\n%s\n\nPresent in Markdown table. Show product name as plain text, plus SKU, category, price, status, and concise recommendations. Include a navigate action card to /products when useful. Never invent stock, margin, or sales numbers unless present in context.", len(products), total, string(productsJSON))
					contextType = "product"
				} else if dataAccessInfo == "" && err == nil {
					dataAccessInfo = "Tidak ada data products sesuai akses dan filter yang diminta."
				} else if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data products dari database. Data mungkin tidak tersedia."
				}
			} else if dataAccessInfo == "" {
				dataAccessInfo = "⚠️ Layanan product belum tersedia untuk AI Assistant."
			}
		}

		// Check if user is asking for forecast data (only if no data fetched yet)
		if contextData == "" && (strings.Contains(messageLower, "forecast") || strings.Contains(messageLower, "grafik forecast") ||
			strings.Contains(messageLower, "prediksi") || strings.Contains(messageLower, "ramalan")) {
			now := time.Now()

			// Check for specific forecast queries
			isNextMonthQuery := strings.Contains(messageLower, "bulan depan") || strings.Contains(messageLower, "next month") ||
				strings.Contains(messageLower, "month depan") || strings.Contains(messageLower, "bulan berikutnya")
			isThreeMonthsQuery := strings.Contains(messageLower, "3 bulan") || strings.Contains(messageLower, "tiga bulan") ||
				strings.Contains(messageLower, "three months") || strings.Contains(messageLower, "3 months")
			isQuarterQuery := strings.Contains(messageLower, "kuartal") || strings.Contains(messageLower, "quarter") ||
				strings.Contains(messageLower, "triwulan")
			isYearQuery := strings.Contains(messageLower, "tahun") && (strings.Contains(messageLower, "ini") ||
				strings.Contains(messageLower, "depan") || strings.Contains(messageLower, "year"))

			var forecastStart, forecastEnd time.Time
			var periodType string

			if isNextMonthQuery {
				// Next month forecast
				nextMonth := now.AddDate(0, 1, 0)
				forecastStart = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
				forecastEnd = forecastStart.AddDate(0, 1, 0).Add(-time.Second)
				periodType = "month"
			} else if isThreeMonthsQuery {
				// 3 months ahead forecast
				forecastStart = now
				forecastEnd = now.AddDate(0, 3, 0)
				periodType = "quarter"
			} else if isQuarterQuery {
				// Quarter forecast - check if next quarter or current
				if strings.Contains(messageLower, "depan") || strings.Contains(messageLower, "berikutnya") || strings.Contains(messageLower, "next") {
					// Next quarter
					quarter := (now.Month() - 1) / 3
					nextQuarter := quarter + 1
					if nextQuarter >= 4 {
						nextQuarter = 0
						forecastStart = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
					} else {
						forecastStart = time.Date(now.Year(), nextQuarter*3+1, 1, 0, 0, 0, 0, now.Location())
					}
					forecastEnd = forecastStart.AddDate(0, 3, 0).Add(-time.Second)
				} else {
					// Current quarter
					quarter := (now.Month() - 1) / 3
					forecastStart = time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
					forecastEnd = forecastStart.AddDate(0, 3, 0).Add(-time.Second)
				}
				periodType = "quarter"
			} else if isYearQuery {
				// Year forecast
				if strings.Contains(messageLower, "depan") || strings.Contains(messageLower, "berikutnya") || strings.Contains(messageLower, "next") {
					// Next year
					forecastStart = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
				} else {
					// Current year
					forecastStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
				}
				forecastEnd = forecastStart.AddDate(1, 0, 0).Add(-time.Second)
				periodType = "year"
			} else {
				// Default: current month forecast
				forecastStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				forecastEnd = forecastStart.AddDate(0, 1, 0).Add(-time.Second)
				periodType = "month"
			}

			periodDesc := "Current Month"
			if isNextMonthQuery {
				periodDesc = "Next Month"
			} else if isThreeMonthsQuery {
				periodDesc = "Next 3 Months"
			} else if isQuarterQuery {
				if strings.Contains(messageLower, "depan") || strings.Contains(messageLower, "berikutnya") || strings.Contains(messageLower, "next") {
					periodDesc = "Next Quarter"
				} else {
					periodDesc = "Current Quarter"
				}
			} else if isYearQuery {
				if strings.Contains(messageLower, "depan") || strings.Contains(messageLower, "berikutnya") || strings.Contains(messageLower, "next") {
					periodDesc = "Next Year"
				} else {
					periodDesc = "Current Year"
				}
			}

			scopedDealIDs := s.scopedUserIDs(userCtx, "deal")
			if scopedDealIDs != nil {
				deals, _, err := s.dealRepo.List(&pipeline.ListDealsRequest{
					Page:          1,
					PerPage:       100,
					DateFrom:      forecastStart.Format(dateFormat),
					DateTo:        forecastEnd.Format(dateFormat),
					ScopedUserIDs: scopedDealIDs,
				})
				if err == nil {
					dealsFormatted := s.formatDealsForAI(deals)
					dealsJSON, _ := json.Marshal(dealsFormatted)
					if contextData != "" {
						contextData += "\n\n"
					}
					contextData += fmt.Sprintf("SCOPED DEALS DATA FOR FORECAST (%s, %d deals):\n%s\n\nCalculate forecast insights ONLY from these scoped deals. If there are no deals, explain that no scoped forecast data is available.", periodDesc, len(deals), string(dealsJSON))
				} else if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data forecast sesuai scope Anda."
				}
			} else {
				forecast, err := s.dealRepo.GetForecast(periodType, forecastStart, forecastEnd)
				if err == nil && forecast != nil {
					forecastJSON, _ := json.Marshal(forecast)
					if contextData != "" {
						contextData += "\n\n"
					}
					contextData += fmt.Sprintf("REAL FORECAST DATA FROM DATABASE (%s):\n%s\n\nCRITICAL: You MUST use ONLY this forecast data. DO NOT create, invent, or make up any forecast data. If forecast data is empty or incomplete, inform the user that forecast data is not available.", periodDesc, string(forecastJSON))
				} else if dataAccessInfo == "" {
					dataAccessInfo = "⚠️ Tidak dapat mengakses data forecast dari database. Data mungkin tidak tersedia."
				}
			}
		}

		// If no specific data type detected but user is asking for general data, default to accounts
		// BUT: Skip if message contains specific data type keywords (lead, deal, account, etc.)
		hasSpecificDataType := strings.Contains(messageLower, "lead") ||
			strings.Contains(messageLower, "deal") ||
			strings.Contains(messageLower, "pipeline") ||
			strings.Contains(messageLower, "account") ||
			strings.Contains(messageLower, "contact") ||
			strings.Contains(messageLower, "visit") ||
			strings.Contains(messageLower, "task") ||
			strings.Contains(messageLower, "product")

		if contextData == "" && !hasSpecificDataType && (strings.Contains(messageLower, "data") || strings.Contains(messageLower, "paparkan") ||
			strings.Contains(messageLower, "tampilkan") || strings.Contains(messageLower, "lihat") ||
			strings.Contains(messageLower, "sistem") || strings.Contains(messageLower, "database")) {
			// Check data privacy and user permissions
			allowed, _ := s.checkDataPrivacy("account", userID, userCtx)
			if !allowed {
				dataAccessInfo = "⚠️ Akses ke data accounts tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
			} else {
				accounts, total, err := s.accountRepo.List(&account.ListAccountsRequest{
					Page:          1,
					PerPage:       10,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "account"),
				})
				if err == nil && len(accounts) > 0 {
					// Transform accounts to user-friendly format with names
					accountsFormatted := s.formatAccountsForAI(accounts)
					accountsJSON, _ := json.Marshal(accountsFormatted)
					contextData = fmt.Sprintf("REAL ACCOUNTS DATA FROM DATABASE (showing %d of %d total accounts):\n%s\n\nCRITICAL INSTRUCTION: You MUST use ONLY the data above. Present it in a Markdown table format. CRITICAL: NEVER show IDs as separate columns - IDs are ONLY used in clickable links. ALWAYS show ONLY NAMES (name, category, city, province) in tables. For clickable actions that trigger detail components, the 'Nama Akun' (Name) column MUST be formatted as [Name](account://id) to create clickable links. Example: [RSUD Jakarta](account://ab868b77-e9b3-429f-ad8c-d55ac1f6561b). Use the EXACT id from the data above - DO NOT create or invent IDs. DO NOT create columns like 'ID', 'Account ID', etc. - these should NOT appear in tables. DO NOT create, invent, or make up any data. DO NOT add columns that don't exist in the data. AFTER presenting the table, ALWAYS provide 1-2 insights, ask 2-3 follow-up questions to understand what the user wants, and offer actionable recommendations.", len(accounts), total, string(accountsJSON))
					contextType = "account"
				} else {
					if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data dari database. Data mungkin tidak tersedia."
					}
				}
			}
		}
	}

	// Get current time in configured timezone
	timezone := settings.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta" // Default to Jakarta timezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		if timezone == "Asia/Jakarta" {
			loc = time.FixedZone("WIB", 7*60*60)
		} else {
			// If timezone is invalid, use UTC
			loc = time.UTC
			fmt.Printf("Warning: Invalid timezone '%s', using UTC instead\n", timezone)
		}
	}

	currentTime := time.Now().In(loc)

	// Build modular system prompt: Core + Domain-specific + Context
	// The domain hint from the frontend is used to select the right domain prompt.
	// If no domain is provided, the router detects it from the user message.
	accessContext := s.buildAIAccessContext(userCtx)
	systemPrompt := BuildModularSystemPrompt(domain, message, contextID, contextType, contextData, dataAccessInfo, accessContext, selectedModel, settings.Provider, currentTime, timezone)

	// Build messages with conversation history
	messages := []cerebras.ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add conversation history (limit to last 10 messages to avoid token limit)
	// Note: No dynamic reduction based on context size - AI will handle overflow automatically
	historyLimit := 10
	startIdx := 0
	if len(conversationHistory) > historyLimit {
		startIdx = len(conversationHistory) - historyLimit
	}

	for i := startIdx; i < len(conversationHistory); i++ {
		msg := conversationHistory[i]
		// Skip system messages from history (only include user and assistant)
		if msg.Role != "user" && msg.Role != "assistant" {
			continue
		}
		content := msg.Content
		// Trim large assistant messages to avoid context-window overflow and
		// REQUEST_TIMEOUT errors. These are usually big data tables or JSON dumps
		// that the LLM doesn't need verbatim in subsequent turns.
		// Use a higher limit for recent messages to preserve entity links/IDs.
		truncLimit := 1500
		if i >= len(conversationHistory)-2 {
			truncLimit = 3000 // Keep more context from recent messages
		}
		if msg.Role == "assistant" && len(content) > truncLimit {
			content = content[:truncLimit] + "\n\n[...pesan terpotong untuk efisiensi token...]"
		}
		messages = append(messages, cerebras.ChatMessage{
			Role:    msg.Role,
			Content: content,
		})
	}

	// Add current user message
	messages = append(messages, cerebras.ChatMessage{
		Role:    "user",
		Content: message,
	})

	// Normalize model name to lowercase for consistent matching
	originalModel := selectedModel
	selectedModel = strings.ToLower(selectedModel)

	// Keep the app on a single supported model. Legacy aliases are accepted so
	// older persisted settings do not break existing users after deployment.
	availableModels := map[string]string{
		"gpt-oss-120b":                   "gpt-oss-120b",
		"gpt-oss":                        "gpt-oss-120b",
		"llama-3.1-8b":                   "gpt-oss-120b",
		"llama3-8b":                      "gpt-oss-120b",
		"llama3.1-8b":                    "gpt-oss-120b",
		"qwen3-235b":                     "gpt-oss-120b",
		"qwen-3-235b-a22b-instruct-2507": "gpt-oss-120b",
		"qwen3-235b-a22b-instruct-2507":  "gpt-oss-120b",
		"qwen-3-235b":                    "gpt-oss-120b",
		"zai-glm-4.7":                    "gpt-oss-120b",
		"zai-glm-4.6":                    "gpt-oss-120b",
		"zai-glm":                        "gpt-oss-120b",
		"zai glm 4.7":                    "gpt-oss-120b",
		"zai glm 4.6":                    "gpt-oss-120b",
		"zai_glm_4.7":                    "gpt-oss-120b",
		"zai_glm_4.6":                    "gpt-oss-120b",
	}

	// Check if model is in the available models map
	if normalizedModel, exists := availableModels[selectedModel]; exists {
		selectedModel = normalizedModel
	} else {
		// Model not found in available models
		// Check if it's a GPT model (not GPT-OSS)
		if strings.HasPrefix(selectedModel, "gpt-") && selectedModel != "gpt-oss-120b" {
			return nil, fmt.Errorf("model '%s' tidak didukung. Model yang tersedia hanya: gpt-oss-120b.", originalModel)
		}
		// For other unknown models, let the API handle it (might be valid but not in our map)
	}

	// Calculate optimal MaxTokens based on context size.
	// More context data means less room for the response — reduce accordingly.
	maxTokens := 4000
	if len(contextData) > 100000 { // Very large context (>100 KB)
		maxTokens = 2000
	} else if len(contextData) > 50000 { // Large context (50–100 KB)
		maxTokens = 3000
	}

	// Call Cerebras API with error handling and panic recovery
	var response *cerebras.ChatResponse
	var apiErr error

	// Add panic recovery for API calls
	func() {
		defer func() {
			if r := recover(); r != nil {
				apiErr = fmt.Errorf("internal error: panic recovered: %v", r)
			}
		}()

		response, apiErr = s.cerebrasClient.Chat(&cerebras.ChatRequest{
			Model:       selectedModel, // Pass the selected model
			Messages:    messages,
			MaxTokens:   maxTokens,
			Temperature: 0.7,
		})
	}()

	if apiErr != nil {
		fmt.Printf("=== CEREBRAS API ERROR ===\n")
		fmt.Printf("Error: %v\n", apiErr)
		fmt.Printf("Error type: %T\n", apiErr)
		fmt.Printf("User message: %s\n", message)
		fmt.Printf("Selected model: %s\n", selectedModel)
		fmt.Printf("==========================\n")

		// Check if error is model not found
		errorStr := apiErr.Error()
		if strings.Contains(errorStr, "model_not_found") ||
			strings.Contains(errorStr, "does not exist") ||
			strings.Contains(errorStr, "not found") {
			return nil, fmt.Errorf("model '%s' tidak ditemukan atau tidak tersedia. Model yang tersedia hanya: gpt-oss-120b.", selectedModel)
		}

		// Check if error is about GPT models
		if strings.Contains(errorStr, "gpt-") {
			return nil, fmt.Errorf("model '%s' tidak didukung. Model Cerebras yang dipakai aplikasi ini hanya gpt-oss-120b.", selectedModel)
		}

		return nil, fmt.Errorf("gagal menghasilkan respons: %w", apiErr)
	}

	// Validate response
	if response == nil {
		return nil, fmt.Errorf("empty response from AI service")
	}

	// Check if Message is nil (defensive check)
	if response.Message.Content == "" {
		return nil, fmt.Errorf("empty message content from AI service")
	}

	// Add data access info to response if needed
	finalMessage := response.Message.Content
	finalMessage = validateGroundedAIAnswer(finalMessage, contextData, dataAccessInfo)
	if dataAccessInfo != "" && !strings.Contains(finalMessage, dataAccessInfo) {
		finalMessage = dataAccessInfo + "\n\n" + finalMessage
	}

	// Execute any TOOL_CALL markers the LLM emitted (real CRUD operations).
	finalMessage = s.processToolCalls(finalMessage, userID, conversationHistory, userCtx)

	return &ai.ChatResponse{
		Message: finalMessage,
		Tokens:  response.Tokens,
	}, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *Service) tryHandlePendingDealAccountConfirmation(message string, history []ai.ChatMessage, userID string, userCtx *domainauth.UserContext) (*ai.ChatResponse, bool) {
	if !isAccountConfirmationReply(message) || !lastAssistantAskedAccountConfirmation(history) {
		return nil, false
	}

	pendingMessage := latestPendingCreateDealUserMessage(history)
	if pendingMessage == "" {
		return nil, false
	}
	if !s.canRunTool("create_deal", userCtx) {
		return &ai.ChatResponse{
			Message: "Anda tidak memiliki permission untuk membuat deal.",
			Tokens:  0,
		}, true
	}

	accountName := cleanConfirmationAccountName(message)
	params := map[string]interface{}{
		"title":        buildConfirmedDealTitle(accountName, pendingMessage),
		"account_name": accountName,
	}
	if stageName := extractDealStageName(pendingMessage); stageName != "" {
		params["stage_name"] = stageName
	}
	if valueText := extractDealValueText(pendingMessage); valueText != "" {
		params["value"] = valueText
	}

	result := s.toolCreateDeal(params, userID, history, userCtx)
	finalMessage := "Saya lanjutkan pembuatan deal menggunakan account yang Anda konfirmasi."
	finalMessage += buildToolResultBlock(result)
	return &ai.ChatResponse{Message: finalMessage, Tokens: 0}, true
}

func (s *Service) tryHandleCreateActivity(message string, history []ai.ChatMessage, userID string, userCtx *domainauth.UserContext) (*ai.ChatResponse, bool) {
	if !isCreateActivityIntent(message) {
		return nil, false
	}

	params := map[string]interface{}{
		"description": cleanActivityDescription(message),
		"type":        inferActivityType(message),
	}
	call := &ToolCall{Tool: "create_activity", Params: params}
	if !s.canRunToolCall(call, userCtx) {
		return &ai.ChatResponse{
			Message: fmt.Sprintf("⚠️ Gagal menjalankan aksi: Anda tidak memiliki permission untuk menjalankan tool '%s'.", call.Tool),
			Tokens:  0,
		}, true
	}

	result := s.executeTool(call, userID, history, userCtx)
	if !result.Success && !result.Confirm {
		return &ai.ChatResponse{Message: strings.TrimSpace(buildToolResultBlock(result)), Tokens: 0}, true
	}
	if result.Confirm {
		return &ai.ChatResponse{Message: strings.TrimSpace(buildToolResultBlock(result)), Tokens: 0}, true
	}

	finalMessage := "Saya mencatat aktivitas sesuai konteks terakhir."
	finalMessage += buildToolResultBlock(result)
	return &ai.ChatResponse{Message: finalMessage, Tokens: 0}, true
}

func isCreateActivityIntent(message string) bool {
	lower := strings.ToLower(message)
	return strings.Contains(lower, "tambahkan aktivitas") ||
		strings.Contains(lower, "tambah aktivitas") ||
		strings.Contains(lower, "catat aktivitas") ||
		strings.Contains(lower, "log activity") ||
		strings.Contains(lower, "add activity")
}

func cleanActivityDescription(message string) string {
	description := strings.TrimSpace(message)
	lower := strings.ToLower(description)
	for _, prefix := range []string{
		"tambahkan aktivitas",
		"tambah aktivitas",
		"catat aktivitas",
		"log activity",
		"add activity",
	} {
		if strings.HasPrefix(lower, prefix) {
			description = strings.TrimSpace(description[len(prefix):])
			break
		}
	}
	description = strings.Trim(description, " .,:;-")
	if description == "" {
		return "Activity added by AI"
	}
	return description
}

func inferActivityType(message string) string {
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "email"):
		return "email"
	case strings.Contains(lower, "task"):
		return "task"
	case strings.Contains(lower, "visit") || strings.Contains(lower, "kunjungan"):
		return "visit"
	case strings.Contains(lower, "deal"):
		return "deal"
	default:
		return "call"
	}
}

func (s *Service) tryHandleUpdateDealStageWithProducts(message string, history []ai.ChatMessage, userID string, userCtx *domainauth.UserContext) (*ai.ChatResponse, bool) {
	lower := strings.ToLower(message)
	if !strings.Contains(lower, "stage") && !strings.Contains(lower, "stages") && !strings.Contains(lower, "status") {
		return nil, false
	}
	if !strings.Contains(lower, "closed won") && !strings.Contains(lower, "won") {
		return nil, false
	}

	params := map[string]interface{}{
		"status": "won",
	}
	if productNames := extractProductNamesFromStageUpdate(message); len(productNames) > 0 {
		params["product_names"] = productNames
	}

	call := &ToolCall{Tool: "update_deal_stage", Params: params}
	if !s.canRunToolCall(call, userCtx) {
		return &ai.ChatResponse{
			Message: fmt.Sprintf("⚠️ Gagal menjalankan aksi: Anda tidak memiliki permission untuk menjalankan tool '%s'.", call.Tool),
			Tokens:  0,
		}, true
	}

	result := s.executeTool(call, userID, history, userCtx)
	if !result.Success && !result.Confirm {
		return &ai.ChatResponse{Message: strings.TrimSpace(buildToolResultBlock(result)), Tokens: 0}, true
	}
	if result.Confirm {
		return &ai.ChatResponse{Message: strings.TrimSpace(buildToolResultBlock(result)), Tokens: 0}, true
	}

	finalMessage := "Saya memperbarui stage deal sesuai konteks terakhir."
	finalMessage += buildToolResultBlock(result)
	return &ai.ChatResponse{Message: finalMessage, Tokens: 0}, true
}

func extractProductNamesFromStageUpdate(message string) []string {
	lower := strings.ToLower(message)
	idx := strings.Index(lower, "produk")
	if idx < 0 {
		idx = strings.Index(lower, "product")
	}
	if idx < 0 {
		return nil
	}

	productsText := strings.TrimSpace(message[idx:])
	fields := strings.Fields(productsText)
	if len(fields) > 0 {
		productsText = strings.TrimSpace(strings.TrimPrefix(productsText, fields[0]))
	}
	productsText = strings.Trim(productsText, " .,:;-")
	if productsText == "" {
		return nil
	}
	return uniqueNonEmpty(splitCSVLike(productsText))
}

func isAccountConfirmationReply(message string) bool {
	cleaned := cleanConfirmationAccountName(message)
	if cleaned == "" {
		return false
	}
	words := strings.Fields(cleaned)
	if len(words) > 8 {
		return false
	}
	lower := strings.ToLower(cleaned)
	return strings.Contains(lower, "rs") ||
		strings.Contains(lower, "hospital") ||
		strings.Contains(lower, "klinik") ||
		strings.Contains(lower, "apotek") ||
		strings.Contains(lower, "kariadi")
}

func lastAssistantAskedAccountConfirmation(history []ai.ChatMessage) bool {
	for i := len(history) - 1; i >= 0; i-- {
		msg := history[i]
		if msg.Role != "assistant" {
			continue
		}
		lower := strings.ToLower(msg.Content)
		return strings.Contains(lower, "mohon konfirmasi account") ||
			strings.Contains(lower, "balas dengan nama account") ||
			strings.Contains(lower, "account yang dimaksud")
	}
	return false
}

func latestPendingCreateDealUserMessage(history []ai.ChatMessage) string {
	for i := len(history) - 1; i >= 0; i-- {
		msg := history[i]
		if msg.Role != "user" {
			continue
		}
		lower := strings.ToLower(msg.Content)
		if strings.Contains(lower, "deal") &&
			(strings.Contains(lower, "stage") ||
				strings.Contains(lower, "tahap") ||
				strings.Contains(lower, "target") ||
				strings.Contains(lower, "nilai")) {
			return msg.Content
		}
	}
	return ""
}

func cleanConfirmationAccountName(message string) string {
	cleaned := strings.TrimSpace(message)
	cleaned = strings.Trim(cleaned, " \t\r\n.,;:!?")
	return strings.Join(strings.Fields(cleaned), " ")
}

func buildConfirmedDealTitle(accountName string, pendingMessage string) string {
	stageName := extractDealStageName(pendingMessage)
	if stageName == "" {
		return "Penawaran " + accountName
	}
	return fmt.Sprintf("Penawaran %s - %s", accountName, stageName)
}

func extractDealStageName(message string) string {
	lower := strings.ToLower(message)
	stageAliases := []struct {
		Needle string
		Name   string
	}{
		{"desire", "Desire"},
		{"qualification", "Qualification"},
		{"proposal sent", "Proposal Sent"},
		{"proposal", "Proposal"},
		{"negotiation", "Negotiation"},
		{"closed won", "Closed Won"},
		{"won", "Closed Won"},
		{"closed lost", "Closed Lost"},
		{"lost", "Closed Lost"},
	}
	for _, alias := range stageAliases {
		if strings.Contains(lower, alias.Needle) {
			return alias.Name
		}
	}
	return ""
}

func extractDealValueText(message string) string {
	valuePattern := regexp.MustCompile(`(?i)(?:rp\s*)?[\d][\d.,]*(?:\s*(?:juta|jt|ribu|rb|miliar|million|billion))?`)
	matches := valuePattern.FindAllString(message, -1)
	for _, match := range matches {
		cleaned := strings.TrimSpace(match)
		if cleaned == "" {
			continue
		}
		if strings.Contains(strings.ToLower(cleaned), "juta") ||
			strings.Contains(strings.ToLower(cleaned), "jt") ||
			strings.Contains(strings.ToLower(cleaned), "ribu") ||
			strings.Contains(strings.ToLower(cleaned), "rb") ||
			strings.Contains(strings.ToLower(cleaned), "miliar") ||
			strings.Contains(strings.ToLower(cleaned), "rp") {
			return cleaned
		}
	}
	if len(matches) > 0 {
		return strings.TrimSpace(matches[len(matches)-1])
	}
	return ""
}

// DealFormatted represents a user-friendly deal format for AI
type DealFormatted struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	AccountID         string  `json:"account_id"`           // ID for creating links
	AccountName       string  `json:"account_name"`         // Name instead of ID
	ContactID         string  `json:"contact_id,omitempty"` // ID for creating links (optional)
	ContactName       string  `json:"contact_name"`         // Name instead of ID
	StageName         string  `json:"stage_name"`           // Name instead of ID
	Value             int64   `json:"value"`
	ValueFormatted    string  `json:"value_formatted"` // Human-readable format
	Status            string  `json:"status"`
	Probability       int     `json:"probability"`
	ExpectedCloseDate *string `json:"expected_close_date,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

// AccountFormatted represents a user-friendly account format for AI
type AccountFormatted struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
	Status   string `json:"status"`
}

type ProductFormatted struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	SKU            string `json:"sku"`
	Category       string `json:"category"`
	Price          int64  `json:"price"`
	PriceFormatted string `json:"price_formatted"`
	Status         string `json:"status"`
	Description    string `json:"description,omitempty"`
}

type UserManagementFormatted struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Group     string `json:"group,omitempty"`
	BrickID   string `json:"brick_id,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type BrickManagementFormatted struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Code                  string `json:"code"`
	Province              string `json:"province"`
	Regency               string `json:"regency"`
	ManagerName           string `json:"manager_name,omitempty"`
	ManagerEmail          string `json:"manager_email,omitempty"`
	Status                string `json:"status"`
	TotalRevenue          int64  `json:"total_revenue"`
	TotalRevenueFormatted string `json:"total_revenue_formatted"`
	DealsClosed           int    `json:"deals_closed"`
	VisitsCompleted       int    `json:"visits_completed"`
	TasksCompleted        int    `json:"tasks_completed"`
	TargetAmount          int64  `json:"target_amount"`
	TargetAmountFormatted string `json:"target_amount_formatted"`
	RevenuePeriodLabel    string `json:"revenue_period_label"`
}

func (s *Service) formatProductsForAI(products []product.Product) []ProductFormatted {
	formatted := make([]ProductFormatted, 0, len(products))
	for _, p := range products {
		categoryName := "N/A"
		if p.Category != nil && p.Category.Name != "" {
			categoryName = p.Category.Name
		}
		formatted = append(formatted, ProductFormatted{
			ID:             p.ID,
			Name:           p.Name,
			SKU:            p.SKU,
			Category:       categoryName,
			Price:          p.Price,
			PriceFormatted: formatCurrencyRupiah(p.Price),
			Status:         p.Status,
			Description:    p.Description,
		})
	}
	return formatted
}

func (s *Service) formatUsersForAI(users []userdomain.User) []UserManagementFormatted {
	formatted := make([]UserManagementFormatted, 0, len(users))
	for _, userEntity := range users {
		roleName := ""
		if userEntity.Role != nil {
			roleName = userEntity.Role.Name
			if roleName == "" {
				roleName = userEntity.Role.Code
			}
		}
		groupName := ""
		if userEntity.Group != nil {
			groupName = userEntity.Group.Name
		}
		brickID := ""
		if userEntity.BrickID != nil {
			brickID = *userEntity.BrickID
		}
		formatted = append(formatted, UserManagementFormatted{
			ID:        userEntity.ID,
			Name:      userEntity.Name,
			Email:     userEntity.Email,
			Role:      roleName,
			Group:     groupName,
			BrickID:   brickID,
			Status:    userEntity.Status,
			CreatedAt: userEntity.CreatedAt.Format("2006-01-02"),
		})
	}
	return formatted
}

func (s *Service) formatBricksForAI(bricks []brickdomain.Brick, revenueByBrickID map[string]brickRevenueSummary, revenuePeriodLabel string) []BrickManagementFormatted {
	formatted := make([]BrickManagementFormatted, 0, len(bricks))
	for _, brickEntity := range bricks {
		managerName := ""
		managerEmail := ""
		if brickEntity.Manager != nil {
			managerName = brickEntity.Manager.Name
			managerEmail = brickEntity.Manager.Email
		}
		revenue := revenueByBrickID[brickEntity.ID]
		formatted = append(formatted, BrickManagementFormatted{
			ID:                    brickEntity.ID,
			Name:                  brickEntity.Name,
			Code:                  brickEntity.Code,
			Province:              brickEntity.Province,
			Regency:               brickEntity.Regency,
			ManagerName:           managerName,
			ManagerEmail:          managerEmail,
			Status:                brickEntity.Status,
			TotalRevenue:          revenue.TotalRevenue,
			TotalRevenueFormatted: revenue.TotalRevenueFormatted,
			DealsClosed:           revenue.DealsClosed,
			VisitsCompleted:       revenue.VisitsCompleted,
			TasksCompleted:        revenue.TasksCompleted,
			TargetAmount:          revenue.TargetAmount,
			TargetAmountFormatted: revenue.TargetAmountFormatted,
			RevenuePeriodLabel:    revenuePeriodLabel,
		})
	}
	return formatted
}

// formatDealsForAI transforms deals to user-friendly format with names
func (s *Service) formatDealsForAI(deals []pipeline.Deal) []DealFormatted {
	formatted := make([]DealFormatted, 0, len(deals))

	for _, deal := range deals {
		accountName := ""
		if deal.Account != nil {
			accountName = deal.Account.Name
		}
		if accountName == "" {
			accountName = "N/A"
		}

		contactName := ""
		if deal.Contact != nil {
			contactName = deal.Contact.Name
		}
		if contactName == "" {
			contactName = "N/A"
		}

		stageName := ""
		if deal.Stage != nil {
			stageName = deal.Stage.Name
		}
		if stageName == "" {
			stageName = "N/A"
		}

		// Format value to Rupiah
		valueFormatted := formatCurrencyRupiah(deal.Value)

		expectedCloseDate := ""
		if deal.ExpectedCloseDate != nil {
			expectedCloseDate = deal.ExpectedCloseDate.Format("2006-01-02")
		}

		// Get account and contact IDs
		accountID := deal.AccountID
		contactID := ""
		if deal.ContactID != nil {
			contactID = *deal.ContactID
		}

		formatted = append(formatted, DealFormatted{
			ID:                deal.ID,
			Title:             deal.Title,
			AccountID:         accountID,
			AccountName:       accountName,
			ContactID:         contactID,
			ContactName:       contactName,
			StageName:         stageName,
			Value:             deal.Value,
			ValueFormatted:    valueFormatted,
			Status:            deal.Status,
			Probability:       deal.Probability,
			ExpectedCloseDate: &expectedCloseDate,
			CreatedAt:         deal.CreatedAt.Format("2006-01-02"),
		})
	}

	return formatted
}

// VisitReportFormatted represents a user-friendly visit report format for AI
type VisitReportFormatted struct {
	ID          string `json:"id"`
	AccountName string `json:"account_name"` // Name instead of ID
	ContactName string `json:"contact_name"` // Name instead of ID
	VisitDate   string `json:"visit_date"`
	Purpose     string `json:"purpose"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// formatVisitReportsForAI transforms visit reports to user-friendly format with names
func (s *Service) formatVisitReportsForAI(visitReports []visit_report.VisitReport) []VisitReportFormatted {
	formatted := make([]VisitReportFormatted, 0, len(visitReports))

	for _, vr := range visitReports {
		// Fetch account name
		accountName := "N/A"
		if vr.AccountID != nil && *vr.AccountID != "" {
			if account, err := s.accountRepo.FindByID(*vr.AccountID); err == nil && account != nil {
				accountName = account.Name
			}
		}

		// Fetch contact name if available
		contactName := "N/A"
		if vr.ContactID != nil {
			if contact, err := s.contactRepo.FindByID(*vr.ContactID); err == nil && contact != nil {
				contactName = contact.Name
			}
		}

		formatted = append(formatted, VisitReportFormatted{
			ID:          vr.ID,
			AccountName: accountName,
			ContactName: contactName,
			VisitDate:   vr.VisitDate.Format("2006-01-02"),
			Purpose:     vr.Purpose,
			Status:      vr.Status,
			CreatedAt:   vr.CreatedAt.Format("2006-01-02"),
		})
	}

	return formatted
}

// formatCurrencyRupiah formats integer (sen) to formatted currency string
func formatCurrencyRupiah(amount int64) string {
	// Convert to Rupiah (divide by 100 if stored in sen)
	rupiah := float64(amount) / 100.0
	// Format with thousand separator
	formatted := formatNumberRupiah(rupiah)
	return "Rp " + formatted
}

// formatNumberRupiah formats number with thousand separator
func formatNumberRupiah(n float64) string {
	// Convert to int64 to remove decimal places
	amount := int64(n)

	// Handle zero case
	if amount == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if amount < 0 {
		negative = true
		amount = -amount
	}

	// Convert to string
	str := fmt.Sprintf("%d", amount)
	length := len(str)

	// Add thousand separators (dot for Indonesian format)
	// Split into chunks of 3 digits from right
	var parts []string
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{str[start:i]}, parts...)
	}

	result := strings.Join(parts, ".")
	if negative {
		result = "-" + result
	}

	return result
}

// TaskFormatted represents a user-friendly task format for AI
type TaskFormatted struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	AccountName string `json:"account_name"` // Name instead of ID
	ContactName string `json:"contact_name"` // Name instead of ID
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// formatTasksForAI transforms tasks to user-friendly format with names
func (s *Service) formatTasksForAI(tasks []task.Task) []TaskFormatted {
	formatted := make([]TaskFormatted, 0, len(tasks))

	for _, t := range tasks {
		// Fetch account name
		accountName := "N/A"
		if t.AccountID != nil {
			if account, err := s.accountRepo.FindByID(*t.AccountID); err == nil && account != nil {
				accountName = account.Name
			}
		}

		// Fetch contact name if available
		contactName := "N/A"
		if t.ContactID != nil {
			if contact, err := s.contactRepo.FindByID(*t.ContactID); err == nil && contact != nil {
				contactName = contact.Name
			}
		}

		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format("2006-01-02")
		}

		formatted = append(formatted, TaskFormatted{
			ID:          t.ID,
			Title:       t.Title,
			AccountName: accountName,
			ContactName: contactName,
			Status:      t.Status,
			Priority:    t.Priority,
			DueDate:     dueDate,
			CreatedAt:   t.CreatedAt.Format("2006-01-02"),
		})
	}

	return formatted
}

// formatAccountsForAI transforms accounts to user-friendly format with names
func (s *Service) formatAccountsForAI(accounts []account.Account) []AccountFormatted {
	formatted := make([]AccountFormatted, 0, len(accounts))

	for _, acc := range accounts {
		categoryName := ""
		if acc.Category != nil {
			categoryName = acc.Category.Name
		}
		if categoryName == "" {
			categoryName = "N/A"
		}

		formatted = append(formatted, AccountFormatted{
			ID:       acc.ID,
			Name:     acc.Name,
			Category: categoryName,
			Address:  acc.Address,
			City:     acc.City,
			Province: acc.Province,
			Status:   acc.Status,
		})
	}

	return formatted
}

// LeadFormatted represents a user-friendly lead format for AI
type LeadFormatted struct {
	ID               string `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	FullName         string `json:"full_name"`
	CompanyName      string `json:"company_name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	JobTitle         string `json:"job_title"`
	LeadSource       string `json:"lead_source"`
	LeadStatus       string `json:"lead_status"`
	LeadScore        int    `json:"lead_score"`
	AccountID        string `json:"account_id"`
	AccountName      string `json:"account_name"` // Name instead of ID
	ContactID        string `json:"contact_id"`
	ContactName      string `json:"contact_name"` // Name instead of ID
	AssignedTo       string `json:"assigned_to"`
	AssignedUserName string `json:"assigned_user_name"` // Name instead of ID
	City             string `json:"city"`
	Province         string `json:"province"`
	CreatedAt        string `json:"created_at"`
}

// formatLeadsForAI transforms leads to user-friendly format with names
func (s *Service) formatLeadsForAI(leads []lead.Lead) []LeadFormatted {
	formatted := make([]LeadFormatted, 0, len(leads))

	for _, l := range leads {
		// Build full name
		fullName := strings.TrimSpace(l.FirstName + " " + l.LastName)
		if fullName == "" {
			fullName = "N/A"
		}

		// Get account name
		accountName := "N/A"
		accountID := ""
		if l.AccountID != nil && *l.AccountID != "" {
			accountID = *l.AccountID
			if l.Account != nil {
				accountName = l.Account.Name
			} else {
				// Try to fetch if not preloaded
				if account, err := s.accountRepo.FindByID(*l.AccountID); err == nil && account != nil {
					accountName = account.Name
				}
			}
		}

		// Get contact name
		contactName := "N/A"
		contactID := ""
		if l.ContactID != nil && *l.ContactID != "" {
			contactID = *l.ContactID
			if l.Contact != nil {
				contactName = l.Contact.Name
			} else {
				// Try to fetch if not preloaded
				if contact, err := s.contactRepo.FindByID(*l.ContactID); err == nil && contact != nil {
					contactName = contact.Name
				}
			}
		}

		// Get assigned user name
		assignedUserName := "N/A"
		assignedTo := ""
		if l.AssignedTo != nil && *l.AssignedTo != "" {
			assignedTo = *l.AssignedTo
			if l.AssignedUser != nil {
				assignedUserName = l.AssignedUser.Name
			}
		}

		formatted = append(formatted, LeadFormatted{
			ID:               l.ID,
			FirstName:        l.FirstName,
			LastName:         l.LastName,
			FullName:         fullName,
			CompanyName:      l.CompanyName,
			Email:            l.Email,
			Phone:            l.Phone,
			JobTitle:         l.JobTitle,
			LeadSource:       l.LeadSource,
			LeadStatus:       l.LeadStatus,
			LeadScore:        l.LeadScore,
			AccountID:        accountID,
			AccountName:      accountName,
			ContactID:        contactID,
			ContactName:      contactName,
			AssignedTo:       assignedTo,
			AssignedUserName: assignedUserName,
			City:             l.City,
			Province:         l.Province,
			CreatedAt:        l.CreatedAt.Format("2006-01-02"),
		})
	}

	return formatted
}

type ProspectPredictionItem struct {
	EntityType         string   `json:"entity_type"`
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	CompanyName        string   `json:"company_name,omitempty"`
	AccountName        string   `json:"account_name,omitempty"`
	Status             string   `json:"status"`
	StageName          string   `json:"stage_name,omitempty"`
	Score              int      `json:"score"`
	Probability        int      `json:"probability"`
	Value              int64    `json:"value"`
	ValueFormatted     string   `json:"value_formatted"`
	ExpectedCloseDate  string   `json:"expected_close_date,omitempty"`
	AssignedUserName   string   `json:"assigned_user_name,omitempty"`
	QualificationScore int      `json:"qualification_score"`
	ScoreBreakdown     []string `json:"score_breakdown"`
	Reasons            []string `json:"reasons"`
	Risks              []string `json:"risks"`
	NextBestAction     string   `json:"next_best_action"`
}

func isProspectPredictionIntent(messageLower string) bool {
	predictionTerms := []string{
		"prediksi", "prediction", "predict", "forecast prospect", "forecast prospek",
		"berpotensi deal", "potensi deal", "potential deal", "peluang deal",
		"prospek potensial", "prospect potensial", "prospect potential",
		"kemungkinan deal", "kemungkinan closing", "siapa yang bisa deal",
		"next best action", "aksi terbaik", "rekomendasi follow", "prioritas follow",
	}
	for _, term := range predictionTerms {
		if strings.Contains(messageLower, term) {
			return true
		}
	}
	return (strings.Contains(messageLower, "prospect") || strings.Contains(messageLower, "prospek")) &&
		(strings.Contains(messageLower, "potensi") || strings.Contains(messageLower, "peluang") || strings.Contains(messageLower, "deal") || strings.Contains(messageLower, "closing"))
}

func isProposalDraftIntent(messageLower string) bool {
	hasProposal := strings.Contains(messageLower, "proposal") ||
		strings.Contains(messageLower, "penawaran") ||
		strings.Contains(messageLower, "quotation") ||
		strings.Contains(messageLower, "quote")
	hasDraftAction := strings.Contains(messageLower, "draft") ||
		strings.Contains(messageLower, "draf") ||
		strings.Contains(messageLower, "buat") ||
		strings.Contains(messageLower, "buatkan") ||
		strings.Contains(messageLower, "susun") ||
		strings.Contains(messageLower, "siapkan")
	wantsCRMWrite := strings.Contains(messageLower, "buat deal") ||
		strings.Contains(messageLower, "create deal") ||
		strings.Contains(messageLower, "simpan") ||
		strings.Contains(messageLower, "masukkan ke crm") ||
		strings.Contains(messageLower, "jadikan opportunity") ||
		strings.Contains(messageLower, "buat opportunity")
	return hasProposal && hasDraftAction && !wantsCRMWrite
}

func buildProposalDraftContext(history []ai.ChatMessage) string {
	var b strings.Builder
	b.WriteString("PROPOSAL DRAFT MODE:\n")
	b.WriteString("- The user is asking for a text/document draft, not to create a CRM deal/opportunity.\n")
	b.WriteString("- Do NOT emit TOOL_CALL for create_deal unless the user explicitly asks to save/create a deal in CRM.\n")
	b.WriteString("- Use prospect name, company, predicted value, timeline, reasons, risks, and next best action from recent conversation when available.\n")
	b.WriteString("- If account_id is missing, do NOT ask for account ID. Use the lead/company name from context and mark unavailable commercial details as placeholders.\n")
	b.WriteString("- Produce a complete draft proposal in Indonesian with sections: title, recipient, background, objective, proposed solution, value/timeline, assumptions, next steps, and closing.\n")

	if len(history) == 0 {
		b.WriteString("\nRECENT CONVERSATION CONTEXT: none. Ask only for business details that are truly missing for the proposal content, not internal IDs.\n")
		return b.String()
	}

	start := len(history) - 6
	if start < 0 {
		start = 0
	}
	b.WriteString("\nRECENT CONVERSATION CONTEXT:\n")
	for i := start; i < len(history); i++ {
		msg := history[i]
		if msg.Role != "user" && msg.Role != "assistant" {
			continue
		}
		content := strings.TrimSpace(msg.Content)
		if len(content) > 1200 {
			content = content[:1200] + "\n[...truncated...]"
		}
		b.WriteString(fmt.Sprintf("\n[%s]\n%s\n", msg.Role, content))
	}
	return b.String()
}

func (s *Service) buildProspectPredictionContext(userID string, userCtx *domainauth.UserContext) (string, string) {
	leadAllowed, _ := s.checkDataPrivacy("lead", userID, userCtx)
	dealAllowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
	if !leadAllowed && !dealAllowed {
		return "", "⚠️ Akses ke data leads/deals tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}

	now := time.Now()
	items := make([]ProspectPredictionItem, 0, 40)
	var leadTotal, dealTotal int64

	if leadAllowed && s.leadRepo != nil {
		leads, total, err := s.leadRepo.List(&lead.ListLeadsRequest{
			Page:          1,
			PerPage:       50,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "lead"),
			Order:         "desc",
		})
		if err == nil {
			leadTotal = total
			for _, l := range leads {
				if isClosedLeadStatus(l.LeadStatus) {
					continue
				}
				items = append(items, buildLeadPredictionItem(l, now))
			}
		}
	}

	if dealAllowed && s.dealRepo != nil {
		deals, total, err := s.dealRepo.List(&pipeline.ListDealsRequest{
			Page:          1,
			PerPage:       50,
			Status:        "open",
			ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
		})
		if err == nil {
			dealTotal = total
			for _, d := range deals {
				items = append(items, buildDealPredictionItem(d, now))
			}
		}
	}

	if len(items) == 0 {
		return "", "Dari hasil data leads dan deals yang dapat Anda akses, belum ada prospect terbuka yang bisa diprediksi saat ini."
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].Value > items[j].Value
		}
		return items[i].Score > items[j].Score
	})
	if len(items) > 15 {
		items = items[:15]
	}

	payload, _ := json.Marshal(items)
	return fmt.Sprintf("REAL PROSPECT PREDICTION DATA (directional CRM scoring, top %d; accessible_leads=%d, accessible_open_deals=%d):\n%s\n\nUse ONLY this scoring context. Present a ranked Markdown table from highest score to lowest. Include clickable [Name](lead://id) or [Name](deal://id), score, probability, value, reasons, risks, and next best action. When the user asks where a score comes from, explain the score_breakdown exactly. Explain this is a directional CRM prediction based on available data, not a guarantee. Do not invent prospects.", len(items), leadTotal, dealTotal, string(payload)), ""
}

func buildLeadPredictionItem(l lead.Lead, now time.Time) ProspectPredictionItem {
	score := maxInt(l.Probability, l.LeadScore)
	breakdown := []string{fmt.Sprintf("base=max(probability %d, lead_score %d)=%d", l.Probability, l.LeadScore, score)}
	reasons := make([]string, 0, 5)
	risks := make([]string, 0, 3)

	status := strings.ToLower(strings.TrimSpace(l.LeadStatus))
	switch status {
	case "qualified":
		score += 25
		breakdown = append(breakdown, "status qualified +25")
		reasons = append(reasons, "status lead sudah qualified")
	case "proposal_sent", "proposal sent":
		score += 20
		breakdown = append(breakdown, "status proposal_sent +20")
		reasons = append(reasons, "proposal sudah dikirim")
	case "interested":
		score += 15
		breakdown = append(breakdown, "status interested +15")
		reasons = append(reasons, "lead menunjukkan minat")
	case "contacted":
		score += 5
		breakdown = append(breakdown, "status contacted +5")
		reasons = append(reasons, "lead sudah dihubungi")
	case "new", "":
		risks = append(risks, "lead masih tahap awal")
	}

	qualificationScore := 0
	if l.BudgetConfirmed {
		score += 8
		qualificationScore += 25
		breakdown = append(breakdown, "budget_confirmed +8")
		reasons = append(reasons, "budget terkonfirmasi")
	}
	if l.AuthorityConfirmed {
		score += 8
		qualificationScore += 25
		breakdown = append(breakdown, "authority_confirmed +8")
		reasons = append(reasons, "decision maker terkonfirmasi")
	}
	if l.NeedConfirmed {
		score += 8
		qualificationScore += 25
		breakdown = append(breakdown, "need_confirmed +8")
		reasons = append(reasons, "kebutuhan terkonfirmasi")
	}
	if l.TimelineConfirmed {
		score += 8
		qualificationScore += 25
		breakdown = append(breakdown, "timeline_confirmed +8")
		reasons = append(reasons, "timeline terkonfirmasi")
	}
	if l.EstimatedValue > 0 {
		score += 5
		breakdown = append(breakdown, "estimated_value > 0 +5")
		reasons = append(reasons, "memiliki estimasi nilai")
	}
	if l.ExpectedCloseDate != nil {
		days := int(l.ExpectedCloseDate.Sub(now).Hours() / 24)
		switch {
		case days >= 0 && days <= 30:
			score += 10
			breakdown = append(breakdown, "expected_close_date within 30 days +10")
			reasons = append(reasons, "target closing dalam 30 hari")
		case days < 0:
			score -= 10
			breakdown = append(breakdown, "expected_close_date overdue -10")
			risks = append(risks, "target closing sudah lewat")
		}
	}
	if now.Sub(l.UpdatedAt) > 14*24*time.Hour {
		score -= 10
		breakdown = append(breakdown, "updated_at older than 14 days -10")
		risks = append(risks, "belum ada update lebih dari 14 hari")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "skor berasal dari lead score/probability yang tersedia")
	}
	if len(risks) == 0 {
		risks = append(risks, "belum ada risiko utama pada data yang tersedia")
	}

	name := strings.TrimSpace(l.FirstName + " " + l.LastName)
	if name == "" {
		name = l.CompanyName
	}
	if name == "" {
		name = "Lead"
	}
	expectedCloseDate := ""
	if l.ExpectedCloseDate != nil {
		expectedCloseDate = l.ExpectedCloseDate.Format(dateFormat)
	}
	assignedUserName := ""
	if l.AssignedUser != nil {
		assignedUserName = l.AssignedUser.Name
	}

	return ProspectPredictionItem{
		EntityType:         "lead",
		ID:                 l.ID,
		Name:               name,
		CompanyName:        l.CompanyName,
		Status:             l.LeadStatus,
		Score:              clampInt(score, 0, 100),
		Probability:        clampInt(l.Probability, 0, 100),
		Value:              l.EstimatedValue,
		ValueFormatted:     formatCurrencyRupiah(l.EstimatedValue),
		ExpectedCloseDate:  expectedCloseDate,
		AssignedUserName:   assignedUserName,
		QualificationScore: qualificationScore,
		ScoreBreakdown:     append(breakdown, fmt.Sprintf("final=%d", clampInt(score, 0, 100))),
		Reasons:            reasons,
		Risks:              risks,
		NextBestAction:     nextBestActionForLead(l, qualificationScore),
	}
}

func buildDealPredictionItem(d pipeline.Deal, now time.Time) ProspectPredictionItem {
	score := d.Probability
	breakdown := []string{fmt.Sprintf("base=probability %d", d.Probability)}
	reasons := make([]string, 0, 5)
	risks := make([]string, 0, 3)

	stageName := ""
	stageCode := ""
	if d.Stage != nil {
		stageName = d.Stage.Name
		stageCode = strings.ToLower(d.Stage.Code)
		if d.Stage.Probability > score {
			score = d.Stage.Probability
			breakdown = append(breakdown, fmt.Sprintf("stage_probability %d overrides base", d.Stage.Probability))
		}
	}
	switch stageCode {
	case "negotiation":
		score += 15
		breakdown = append(breakdown, "stage negotiation +15")
		reasons = append(reasons, "deal berada di tahap negotiation")
	case "proposal", "proposal_sent":
		score += 10
		breakdown = append(breakdown, "stage proposal +10")
		reasons = append(reasons, "deal sudah masuk tahap proposal")
	case "qualification":
		score += 5
		breakdown = append(breakdown, "stage qualification +5")
		reasons = append(reasons, "deal sudah masuk tahap qualification")
	}

	qualificationScore := 0
	if d.BudgetConfirmed {
		score += 5
		qualificationScore += 25
		breakdown = append(breakdown, "budget_confirmed +5")
		reasons = append(reasons, "budget terkonfirmasi")
	}
	if d.AuthorityConfirmed {
		score += 5
		qualificationScore += 25
		breakdown = append(breakdown, "authority_confirmed +5")
		reasons = append(reasons, "decision maker terkonfirmasi")
	}
	if d.NeedConfirmed {
		score += 5
		qualificationScore += 25
		breakdown = append(breakdown, "need_confirmed +5")
		reasons = append(reasons, "kebutuhan terkonfirmasi")
	}
	if d.TimelineConfirmed {
		score += 5
		qualificationScore += 25
		breakdown = append(breakdown, "timeline_confirmed +5")
		reasons = append(reasons, "timeline terkonfirmasi")
	}
	if d.Value > 0 {
		score += 5
		breakdown = append(breakdown, "value > 0 +5")
		reasons = append(reasons, "memiliki nilai deal")
	}
	if d.ExpectedCloseDate != nil {
		days := int(d.ExpectedCloseDate.Sub(now).Hours() / 24)
		switch {
		case days >= 0 && days <= 30:
			score += 10
			breakdown = append(breakdown, "expected_close_date within 30 days +10")
			reasons = append(reasons, "target closing dalam 30 hari")
		case days < 0:
			score -= 15
			breakdown = append(breakdown, "expected_close_date overdue -15")
			risks = append(risks, "expected close date sudah lewat")
		}
	}
	if now.Sub(d.UpdatedAt) > 14*24*time.Hour {
		score -= 10
		breakdown = append(breakdown, "updated_at older than 14 days -10")
		risks = append(risks, "deal belum di-update lebih dari 14 hari")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "skor berasal dari probability/stage yang tersedia")
	}
	if len(risks) == 0 {
		risks = append(risks, "belum ada risiko utama pada data yang tersedia")
	}

	accountName := ""
	if d.Account != nil {
		accountName = d.Account.Name
	}
	expectedCloseDate := ""
	if d.ExpectedCloseDate != nil {
		expectedCloseDate = d.ExpectedCloseDate.Format(dateFormat)
	}
	assignedUserName := ""
	if d.AssignedUser != nil {
		assignedUserName = d.AssignedUser.Name
	}

	return ProspectPredictionItem{
		EntityType:         "deal",
		ID:                 d.ID,
		Name:               d.Title,
		AccountName:        accountName,
		Status:             d.Status,
		StageName:          stageName,
		Score:              clampInt(score, 0, 100),
		Probability:        clampInt(d.Probability, 0, 100),
		Value:              d.Value,
		ValueFormatted:     formatCurrencyRupiah(d.Value),
		ExpectedCloseDate:  expectedCloseDate,
		AssignedUserName:   assignedUserName,
		QualificationScore: qualificationScore,
		ScoreBreakdown:     append(breakdown, fmt.Sprintf("final=%d", clampInt(score, 0, 100))),
		Reasons:            reasons,
		Risks:              risks,
		NextBestAction:     nextBestActionForDeal(d, stageCode, qualificationScore),
	}
}

func isClosedLeadStatus(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return status == "converted" || status == "lost" || status == "unqualified"
}

func nextBestActionForLead(l lead.Lead, qualificationScore int) string {
	switch {
	case !l.BudgetConfirmed:
		return "Validasi budget dan estimasi nilai kebutuhan."
	case !l.AuthorityConfirmed:
		return "Identifikasi decision maker dan jadwalkan follow-up."
	case !l.NeedConfirmed:
		return "Konfirmasi kebutuhan produk dan pain point utama."
	case !l.TimelineConfirmed:
		return "Kunci timeline pembelian atau jadwal evaluasi."
	case strings.EqualFold(l.LeadStatus, "qualified"):
		return "Konversi menjadi deal/opportunity jika account dan contact sudah jelas."
	case qualificationScore >= 75:
		return "Siapkan proposal dan jadwalkan meeting closing."
	default:
		return "Lakukan follow-up terarah untuk melengkapi kualifikasi BANT."
	}
}

func nextBestActionForDeal(d pipeline.Deal, stageCode string, qualificationScore int) string {
	switch {
	case d.ExpectedCloseDate != nil && d.ExpectedCloseDate.Before(time.Now()):
		return "Update expected close date dan klarifikasi blocker closing."
	case stageCode == "negotiation":
		return "Follow-up negosiasi, konfirmasi blocker, dan minta komitmen next step."
	case stageCode == "proposal" || stageCode == "proposal_sent":
		return "Follow-up proposal dan validasi approval internal customer."
	case qualificationScore < 75:
		return "Lengkapi BANT sebelum mendorong closing."
	default:
		return "Jadwalkan closing call dan siapkan ringkasan value proposition."
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func wantsLeadFullDetail(messageLower string) bool {
	return strings.Contains(messageLower, "lengkap") ||
		strings.Contains(messageLower, "detail") ||
		strings.Contains(messageLower, "semua data") ||
		strings.Contains(messageLower, "full") ||
		strings.Contains(messageLower, "log visit") ||
		strings.Contains(messageLower, "activity") ||
		strings.Contains(messageLower, "aktivitas") ||
		strings.Contains(messageLower, "product interest") ||
		strings.Contains(messageLower, "minat produk") ||
		strings.Contains(messageLower, "bant") ||
		strings.Contains(messageLower, "task")
}

func (s *Service) buildLeadFullDetailsForAI(leads []lead.Lead, userCtx *domainauth.UserContext) string {
	type leadFullDetail struct {
		Lead             LeadFormatted `json:"lead"`
		VisitReports     interface{}   `json:"visit_reports"`
		Activities       interface{}   `json:"activities"`
		Tasks            interface{}   `json:"tasks"`
		BANT             interface{}   `json:"bant"`
		ProductInterests []interface{} `json:"product_interests"`
	}

	maxLeads := len(leads)
	if maxLeads > 5 {
		maxLeads = 5
	}
	details := make([]leadFullDetail, 0, maxLeads)

	for i := 0; i < maxLeads; i++ {
		leadEntity := leads[i]
		item := leadFullDetail{
			Lead: s.formatLeadsForAI([]lead.Lead{leadEntity})[0],
		}

		productInterests := make([]interface{}, 0)

		if s.visitReportRepo != nil {
			visitReports, _, err := s.visitReportRepo.List(&visit_report.ListVisitReportsRequest{
				Page:          1,
				PerPage:       20,
				LeadID:        leadEntity.ID,
				ScopedUserIDs: s.scopedUserIDs(userCtx, "visit_report"),
			})
			if err == nil {
				item.VisitReports = s.formatVisitReportsForAI(visitReports)
				productInterests = append(productInterests, productInterestsFromVisitReports(visitReports)...)
			}
		}

		if s.activityRepo != nil {
			activities, _, err := s.activityRepo.List(&activity.ListActivitiesRequest{
				Page:          1,
				PerPage:       20,
				LeadID:        leadEntity.ID,
				ScopedUserIDs: s.scopedUserIDs(userCtx, "activity"),
			})
			if err == nil {
				item.Activities = activities
				productInterests = append(productInterests, productInterestsFromActivities(activities)...)
			}
		}

		if s.taskRepo != nil {
			tasks, _, err := s.taskRepo.List(&task.ListTasksRequest{
				Page:          1,
				PerPage:       20,
				LeadID:        leadEntity.ID,
				ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
			})
			if err == nil {
				item.Tasks = s.formatTasksForAI(tasks)
			}
		}

		if s.leadQualificationService != nil {
			if qualification, err := s.leadQualificationService.GetByLeadID(leadEntity.ID); err == nil {
				item.BANT = qualification
				for _, product := range qualification.NeedTargetProducts {
					productInterests = append(productInterests, map[string]string{
						"product_id":   product.ProductID,
						"product_name": product.ProductName,
						"source":       "bant",
					})
				}
			}
		}

		item.ProductInterests = productInterests
		details = append(details, item)
	}

	if len(details) == 0 {
		return ""
	}
	detailsJSON, _ := json.Marshal(details)
	return fmt.Sprintf("LEAD COMPLETE DETAILS (include log visit, activity, product interest, task, and BANT when asked):\n%s\n\nWhen presenting complete lead data, group the answer by lead and include sections for Log Visit, Activity, Product Interest, Task, and BANT. If a section is empty, say it is empty.", string(detailsJSON))
}

func productInterestsFromVisitReports(visitReports []visit_report.VisitReport) []interface{} {
	result := make([]interface{}, 0)
	for _, visitReport := range visitReports {
		result = append(result, productInterestsFromMetadata(visitReport.Metadata, "visit_report", visitReport.ID)...)
	}
	return result
}

func productInterestsFromActivities(activities []activity.Activity) []interface{} {
	result := make([]interface{}, 0)
	for _, activityEntity := range activities {
		result = append(result, productInterestsFromMetadata(activityEntity.Metadata, "activity", activityEntity.ID)...)
	}
	return result
}

func productInterestsFromMetadata(raw []byte, source string, sourceID string) []interface{} {
	if len(raw) == 0 {
		return nil
	}
	var metadata map[string]interface{}
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil
	}
	value, ok := metadata["product_interests"]
	if !ok {
		return nil
	}
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			itemMap["source"] = source
			itemMap["source_id"] = sourceID
			result = append(result, itemMap)
			continue
		}
		result = append(result, map[string]interface{}{
			"product_name": item,
			"source":       source,
			"source_id":    sourceID,
		})
	}
	return result
}

// parseVisitReportInsight parses AI response into VisitReportInsight
func (s *Service) parseVisitReportInsight(text string) (*ai.VisitReportInsight, error) {
	// Clean up the text: remove comment markers and extra whitespace
	cleaned := strings.TrimSpace(text)

	// Remove common comment markers that might appear before/after JSON
	cleaned = strings.TrimPrefix(cleaned, "*/")
	cleaned = strings.TrimPrefix(cleaned, "/*")
	cleaned = strings.TrimPrefix(cleaned, "//")
	cleaned = strings.TrimSpace(cleaned)

	// Remove trailing comments (like "// Output format in JSON format.")
	// Find the last closing brace and remove everything after it that's not part of JSON
	jsonStart := strings.Index(cleaned, "{")
	jsonEnd := strings.LastIndex(cleaned, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	// Extract JSON portion
	jsonStr := cleaned[jsonStart : jsonEnd+1]

	// Note: We don't remove comment markers from inside the JSON string
	// because they might be valid parts of JSON string values (e.g., URLs).
	// Comment markers should only appear outside the JSON boundaries, which
	// we've already handled by extracting the JSON portion.
	jsonStr = strings.TrimSpace(jsonStr)

	// Try to parse as JSON
	var rawInsight map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawInsight); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Build the insight struct, handling different data types
	insight := &ai.VisitReportInsight{
		Summary:         "",
		ActionItems:     []string{},
		Sentiment:       "neutral",
		KeyPoints:       []string{},
		Recommendations: []string{},
	}

	// Extract summary
	if summary, ok := rawInsight["summary"].(string); ok {
		insight.Summary = summary
	}

	// Extract sentiment
	if sentiment, ok := rawInsight["sentiment"].(string); ok {
		insight.Sentiment = sentiment
	}

	// Extract key_points (array of strings)
	if keyPoints, ok := rawInsight["key_points"].([]interface{}); ok {
		for _, point := range keyPoints {
			if str, ok := point.(string); ok {
				insight.KeyPoints = append(insight.KeyPoints, str)
			}
		}
	}

	// Extract action_items (can be array of strings or array of objects)
	if actionItems, ok := rawInsight["action_items"].([]interface{}); ok {
		for _, item := range actionItems {
			if str, ok := item.(string); ok {
				// Simple string
				insight.ActionItems = append(insight.ActionItems, str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				// Object with description and urgency - convert to string
				if desc, ok := obj["description"].(string); ok {
					insight.ActionItems = append(insight.ActionItems, desc)
				}
			}
		}
	}

	// Extract recommendations (array of strings)
	if recommendations, ok := rawInsight["recommendations"].([]interface{}); ok {
		for _, rec := range recommendations {
			if str, ok := rec.(string); ok {
				insight.Recommendations = append(insight.Recommendations, str)
			}
		}
	}

	return insight, nil
}

// extractLocationFromHistory scans conversation history for a user message that
// looks like "Saya berada di {lat}, {lng}" and returns the parsed coordinates.
func extractLocationFromHistory(history []ai.ChatMessage) (lat, lng float64, ok bool) {
	// Regex-free approach: look for user messages containing coordinate patterns
	for i := len(history) - 1; i >= 0; i-- {
		msg := history[i]
		if msg.Role != "user" {
			continue
		}
		lower := strings.ToLower(msg.Content)
		// Match patterns like "1.352083, 103.819836" or "berada di lat,lng"
		if !strings.Contains(lower, ",") {
			continue
		}
		// Try to parse two floats separated by comma (with optional spaces)
		parts := strings.SplitN(msg.Content, ",", 2)
		if len(parts) != 2 {
			continue
		}
		// Strip non-numeric prefix from first part (e.g. "Saya berada di 1.352083")
		latStr := strings.TrimSpace(parts[0])
		for len(latStr) > 0 && latStr[0] != '-' && (latStr[0] < '0' || latStr[0] > '9') {
			latStr = latStr[1:]
		}
		lngStr := strings.TrimSpace(parts[1])
		// Remove any trailing text after the number
		for j, c := range lngStr {
			if c != '.' && c != '-' && (c < '0' || c > '9') {
				lngStr = lngStr[:j]
				break
			}
		}
		var parsedLat, parsedLng float64
		if _, err := fmt.Sscanf(latStr, "%f", &parsedLat); err != nil {
			continue
		}
		if _, err := fmt.Sscanf(lngStr, "%f", &parsedLng); err != nil {
			continue
		}
		if parsedLat >= -90 && parsedLat <= 90 && parsedLng >= -180 && parsedLng <= 180 {
			return parsedLat, parsedLng, true
		}
	}
	return 0, 0, false
}

// extractAccountIDsFromHistory scans assistant messages in history for
// "account://{uuid}" markdown link patterns and returns unique IDs.
func extractAccountIDsFromHistory(history []ai.ChatMessage) []string {
	seen := make(map[string]struct{})
	var ids []string
	accountLinkPrefix := "account://"
	for i := len(history) - 1; i >= 0; i-- {
		msg := history[i]
		if msg.Role != "assistant" {
			continue
		}
		content := msg.Content
		idx := 0
		for {
			pos := strings.Index(content[idx:], accountLinkPrefix)
			if pos < 0 {
				break
			}
			start := idx + pos + len(accountLinkPrefix)
			end := start
			for end < len(content) && content[end] != ')' && content[end] != ' ' && content[end] != '\n' {
				end++
			}
			id := content[start:end]
			if len(id) == 36 { // UUID length
				if _, exists := seen[id]; !exists {
					seen[id] = struct{}{}
					ids = append(ids, id)
				}
			}
			idx = start
		}
		if len(ids) > 0 {
			break // Use IDs from the most recent assistant message that has them
		}
	}
	return ids
}

// buildWaypointsFromAccountIDs fetches accounts by ID and builds route waypoints.
// Invalid or not-found IDs are silently skipped.
func (s *Service) buildWaypointsFromAccountIDs(accountIDs []string, userCtx *domainauth.UserContext) []route_optimization_domain.Waypoint {
	var waypoints []route_optimization_domain.Waypoint
	for _, id := range accountIDs {
		acc, err := s.accountRepo.FindByID(id)
		if err != nil || acc == nil {
			continue
		}
		accountOwner := ""
		if acc.AssignedTo != nil {
			accountOwner = *acc.AssignedTo
		}
		if !s.canAccessOwner(userCtx, "account", accountOwner) {
			continue
		}
		// Require either a lat/lng or at least an address to geocode later
		if acc.Latitude == nil || acc.Longitude == nil {
			continue
		}
		name := acc.Name
		accID := acc.ID
		address := acc.Address
		if acc.City != "" {
			if address != "" {
				address = address + ", " + acc.City
			} else {
				address = acc.City
			}
		}
		waypoints = append(waypoints, route_optimization_domain.Waypoint{
			Lat:         *acc.Latitude,
			Lng:         *acc.Longitude,
			Address:     address,
			AccountID:   &accID,
			AccountName: &name,
		})
		if len(waypoints) >= 20 { // Hard cap to avoid overly long routes from AI
			break
		}
	}
	return waypoints
}

// buildCRUDContext gathers relevant entity data from the database when the user
// expresses an intent to create or update a CRM entity. It scans the
// conversation history for entity IDs/names mentioned by the LLM and fetches
// their details so the LLM can emit a proper TOOL_CALL with real IDs.
func (s *Service) buildCRUDContext(messageLower string, history []ai.ChatMessage, userID string, userCtx *domainauth.UserContext) string {
	var parts []string

	// 1. Extract entity IDs from conversation history (all entity types).
	entityIDs := s.extractEntityIDsFromHistory(history)

	// 2. Load contacts referenced in history or by name in the current message.
	contactCtx := s.loadContactsForCRUD(messageLower, entityIDs, userID, userCtx)
	if contactCtx != "" {
		parts = append(parts, contactCtx)
	}

	// 3. Load accounts referenced in history.
	accountCtx := s.loadAccountsForCRUD(messageLower, entityIDs, userID, userCtx)
	if accountCtx != "" {
		parts = append(parts, accountCtx)
	}

	// 4. Load deals referenced in history.
	dealCtx := s.loadDealsForCRUD(messageLower, entityIDs, userID, userCtx)
	if dealCtx != "" {
		parts = append(parts, dealCtx)
	}

	// 5. Load tasks referenced by history/name for task updates.
	if isTaskCRUDIntent(messageLower) {
		taskCtx := s.loadTasksForCRUD(messageLower, entityIDs, userID, userCtx)
		if taskCtx != "" {
			parts = append(parts, taskCtx)
		}
	}

	// 6. Load schedules referenced by history/name for schedule updates.
	if isScheduleCRUDIntent(messageLower) {
		scheduleCtx := s.loadSchedulesForCRUD(messageLower, entityIDs, userID, userCtx)
		if scheduleCtx != "" {
			parts = append(parts, scheduleCtx)
		}
	}

	// 7. Load leads/statuses when the user mentions lead mutation or lead-linked activity.
	if strings.Contains(messageLower, "lead") && isLeadCRUDIntent(messageLower) {
		leadCtx := s.loadLeadStatusesForCRUD(messageLower, entityIDs, userID, userCtx)
		if leadCtx != "" {
			parts = append(parts, leadCtx)
		}
	}

	// 8. Load pipeline stages if the user mentions moving a deal stage.
	if strings.Contains(messageLower, "deal") || strings.Contains(messageLower, "stage") || strings.Contains(messageLower, "pindah") || strings.Contains(messageLower, "pipeline") {
		stageCtx := s.loadPipelineStagesForCRUD()
		if stageCtx != "" {
			parts = append(parts, stageCtx)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	header := "CRUD CONTEXT — Use the following real entity data to construct the TOOL_CALL. Do NOT ask the user for IDs that are already available below.\n\n"
	return header + strings.Join(parts, "\n\n")
}

func isTaskCRUDIntent(messageLower string) bool {
	return (strings.Contains(messageLower, "task") || strings.Contains(messageLower, "tugas")) &&
		(strings.Contains(messageLower, "status") ||
			strings.Contains(messageLower, "ubah") ||
			strings.Contains(messageLower, "update") ||
			strings.Contains(messageLower, "complete") ||
			strings.Contains(messageLower, "completed") ||
			strings.Contains(messageLower, "selesai"))
}

func isScheduleCRUDIntent(messageLower string) bool {
	return (strings.Contains(messageLower, "schedule") ||
		strings.Contains(messageLower, "jadwal") ||
		strings.Contains(messageLower, "meeting") ||
		strings.Contains(messageLower, "rapat") ||
		strings.Contains(messageLower, "reschedule")) &&
		(strings.Contains(messageLower, "ubah") ||
			strings.Contains(messageLower, "update") ||
			strings.Contains(messageLower, "ganti") ||
			strings.Contains(messageLower, "menjadi") ||
			strings.Contains(messageLower, "tanggal") ||
			strings.Contains(messageLower, "jam") ||
			strings.Contains(messageLower, "reschedule"))
}

func isTaskStatusUpdateIntent(messageLower string) bool {
	return isTaskCRUDIntent(messageLower) &&
		(strings.Contains(messageLower, " ke ") ||
			strings.Contains(messageLower, " menjadi ") ||
			strings.Contains(messageLower, "complete") ||
			strings.Contains(messageLower, "completed") ||
			strings.Contains(messageLower, "selesai"))
}

func taskStatusFilterFromMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "pending"):
		return "pending"
	case strings.Contains(messageLower, "in_progress"), strings.Contains(messageLower, "in-progress"):
		return "in_progress"
	case strings.Contains(messageLower, "completed"), strings.Contains(messageLower, "done"), strings.Contains(messageLower, "selesai"):
		return "completed"
	case strings.Contains(messageLower, "cancelled"), strings.Contains(messageLower, "canceled"), strings.Contains(messageLower, "batal"):
		return "cancelled"
	default:
		return ""
	}
}

func taskSearchTermFromMessage(messageLower string) string {
	if !(strings.Contains(messageLower, "task") || strings.Contains(messageLower, "tugas")) {
		return ""
	}

	replacer := strings.NewReplacer(
		"ubah status task ", "",
		"update status task ", "",
		"ubah task ", "",
		"update task ", "",
		"status task ", "",
		"task ", "",
		"tugas ", "",
	)
	candidate := strings.TrimSpace(replacer.Replace(messageLower))

	stopPhrases := []string{
		" ke completed", " ke complete", " ke done", " ke selesai", " ke pending", " ke cancelled", " ke canceled",
		" menjadi completed", " menjadi complete", " menjadi done", " menjadi selesai", " menjadi pending", " menjadi cancelled", " menjadi canceled",
		" jadi completed", " jadi complete", " jadi done", " jadi selesai", " jadi pending", " jadi cancelled", " jadi canceled",
	}
	for _, stop := range stopPhrases {
		if idx := strings.Index(candidate, stop); idx >= 0 {
			candidate = strings.TrimSpace(candidate[:idx])
			break
		}
	}

	candidate = strings.Trim(candidate, ".,;:!? ")
	if candidate == "" || candidate == messageLower {
		return ""
	}
	return candidate
}

func isLeadCRUDIntent(messageLower string) bool {
	for _, keyword := range []string{
		"status", "ubah", "update",
		"activity", "aktivitas", "catat aktivitas", "log activity",
		"visit", "kunjungan", "log visit",
		"bant", "qualification", "kualifikasi", "budget", "authority", "need", "timeline",
		"tambah", "tambahkan", "add",
	} {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}
	return false
}

func isMyProductSalesIntent(messageLower string) bool {
	return isProductSalesAnalyticsIntent(messageLower) && hasOwnSalesQualifier(messageLower)
}

func isProductSalesAnalyticsIntent(messageLower string) bool {
	hasProduct := strings.Contains(messageLower, "product") ||
		strings.Contains(messageLower, "produk") ||
		strings.Contains(messageLower, "obat")
	hasSales := strings.Contains(messageLower, "terjual") ||
		strings.Contains(messageLower, "dijual") ||
		strings.Contains(messageLower, "penjualan") ||
		strings.Contains(messageLower, "penjuualan") ||
		strings.Contains(messageLower, "pejualan") ||
		strings.Contains(messageLower, "pejuaalan") ||
		strings.Contains(messageLower, "jual") ||
		strings.Contains(messageLower, "sold") ||
		strings.Contains(messageLower, "sales") ||
		strings.Contains(messageLower, "revenue") ||
		strings.Contains(messageLower, "pendapatan") ||
		strings.Contains(messageLower, "omzet") ||
		strings.Contains(messageLower, "kontribusi") ||
		strings.Contains(messageLower, "paling banyak") ||
		strings.Contains(messageLower, "paling laku") ||
		strings.Contains(messageLower, "paling sering") ||
		strings.Contains(messageLower, "jarang")
	if !hasSales {
		return false
	}

	hasSalesDataQuery := strings.Contains(messageLower, "data penjualan") ||
		strings.Contains(messageLower, "data penjuualan") ||
		strings.Contains(messageLower, "dara penjuualan") ||
		strings.Contains(messageLower, "laporan penjualan")

	if hasTrendOrChartTerm(messageLower) && !hasProduct {
		return false
	}

	return hasProduct || hasSalesDataQuery
}

func hasOwnSalesQualifier(messageLower string) bool {
	return strings.Contains(messageLower, "saya") ||
		strings.Contains(messageLower, "olehku") ||
		strings.Contains(messageLower, "milik saya") ||
		strings.Contains(messageLower, "my ") ||
		strings.Contains(messageLower, " mine")
}

func (s *Service) productSalesScopeKind(userCtx *domainauth.UserContext) string {
	if userCtx == nil {
		return "own"
	}
	scopedUserIDs := s.scopedUserIDs(userCtx, "product_analysis")
	if scopedUserIDs == nil {
		return "global"
	}
	if len(scopedUserIDs) == 1 && scopedUserIDs[0] == userCtx.UserID {
		return "own"
	}
	return "team"
}

func buildNoProductSalesMessage(periodLabel string, scopeKind string) string {
	periodLabel = strings.TrimSpace(periodLabel)
	if periodLabel == "" {
		periodLabel = "periode yang diminta"
	}
	switch scopeKind {
	case "global":
		return fmt.Sprintf("Dari hasil seluruh data penjualan yang dapat Anda akses, belum ada produk terjual pada %s.", periodLabel)
	case "team":
		return fmt.Sprintf("Dari hasil data penjualan tim Anda, belum ada produk terjual pada %s.", periodLabel)
	default:
		return fmt.Sprintf("Dari hasil data penjualan yang Anda miliki, Anda belum menjual suatu produk pada %s.", periodLabel)
	}
}

func buildNoUserProductSalesMessage(periodLabel string) string {
	return buildNoProductSalesMessage(periodLabel, "own")
}

// extractEntityIDsFromHistory scans assistant messages in conversation history
// for all entity type links (account://UUID, contact://UUID, deal://UUID, etc.)
// and returns them grouped by type.
func (s *Service) extractEntityIDsFromHistory(history []ai.ChatMessage) map[string][]string {
	result := map[string][]string{}
	prefixes := []string{"account://", "contact://", "deal://", "lead://", "task://", "visit://", "schedule://"}

	for i := len(history) - 1; i >= 0; i-- {
		msg := history[i]
		if msg.Role != "assistant" {
			continue
		}
		content := msg.Content
		for _, prefix := range prefixes {
			entityType := strings.TrimSuffix(prefix, "://")
			idx := 0
			for {
				pos := strings.Index(content[idx:], prefix)
				if pos < 0 {
					break
				}
				start := idx + pos + len(prefix)
				end := start
				for end < len(content) && content[end] != ')' && content[end] != ' ' && content[end] != '\n' && content[end] != '|' && content[end] != ']' {
					end++
				}
				id := content[start:end]
				if len(id) == 36 { // UUID length
					result[entityType] = appendUnique(result[entityType], id)
				}
				idx = start
			}
		}

		for _, match := range aiActionPattern.FindAllString(content, -1) {
			actionJSON := strings.TrimSpace(strings.TrimPrefix(match, "<!-- ACTION:"))
			actionJSON = strings.TrimSpace(strings.TrimSuffix(actionJSON, "-->"))
			var action map[string]string
			if err := json.Unmarshal([]byte(actionJSON), &action); err != nil {
				continue
			}
			entityType := strings.ToLower(strings.TrimSpace(action["entity"]))
			entityID := strings.TrimSpace(action["entityId"])
			if entityType != "" && len(entityID) == 36 {
				result[entityType] = appendUnique(result[entityType], entityID)
			}
		}
	}
	return result
}

// loadContactsForCRUD loads contacts by IDs from history and/or by name
// mentioned in the current message.
func (s *Service) loadContactsForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("contact", userID, userCtx)
	if !allowed {
		return ""
	}

	var contacts []interface{}

	// Load by ID from history
	for _, cid := range entityIDs["contact"] {
		c, err := s.contactRepo.FindByID(cid)
		if err == nil && c != nil {
			if acc, accErr := s.accountRepo.FindByID(c.AccountID); accErr == nil && acc != nil {
				accountOwner := ""
				if acc.AssignedTo != nil {
					accountOwner = *acc.AssignedTo
				}
				if s.canAccessOwner(userCtx, "account", accountOwner) {
					contacts = append(contacts, c)
				}
			}
		}
	}

	// Also try to find contacts by name mentioned in the message.
	// Extract names from conversation history assistant messages to match.
	nameHints := extractNamesFromHistory(messageLower)
	if len(nameHints) > 0 {
		allContacts, _, err := s.contactRepo.List(&contact.ListContactsRequest{
			Page:          1,
			PerPage:       50,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "contact"),
		})
		if err == nil {
			for _, c := range allContacts {
				fullName := strings.ToLower(c.Name)
				for _, hint := range nameHints {
					if strings.Contains(fullName, hint) || strings.Contains(hint, fullName) {
						contacts = append(contacts, c)
						break
					}
				}
			}
		}
	}

	if len(contacts) == 0 {
		return ""
	}
	contactsJSON, _ := json.Marshal(contacts)
	return fmt.Sprintf("AVAILABLE CONTACTS (use these IDs for contact_id in TOOL_CALL):\n%s", string(contactsJSON))
}

// loadAccountsForCRUD loads accounts by IDs from history.
func (s *Service) loadAccountsForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("account", userID, userCtx)
	if !allowed {
		return ""
	}

	var accounts []interface{}
	seen := map[string]bool{}
	for _, aid := range entityIDs["account"] {
		a, err := s.accountRepo.FindByID(aid)
		if err == nil && a != nil {
			accountOwner := ""
			if a.AssignedTo != nil {
				accountOwner = *a.AssignedTo
			}
			if s.canAccessOwner(userCtx, "account", accountOwner) && !seen[a.ID] {
				accounts = append(accounts, a)
				seen[a.ID] = true
			}
		}
	}

	for _, term := range extractNamesFromHistory(messageLower) {
		results, err := s.searchAccountsForTool([]string{term}, userCtx)
		if err != nil {
			continue
		}
		for i := range results {
			a := results[i]
			if seen[a.ID] {
				continue
			}
			accountOwner := ""
			if a.AssignedTo != nil {
				accountOwner = *a.AssignedTo
			}
			if !s.canAccessOwner(userCtx, "account", accountOwner) {
				continue
			}
			accountCopy := a
			accounts = append(accounts, &accountCopy)
			seen[a.ID] = true
			if len(accounts) >= 10 {
				break
			}
		}
		if len(accounts) >= 10 {
			break
		}
	}

	if len(accounts) == 0 {
		return ""
	}
	accountsJSON, _ := json.Marshal(accounts)
	return fmt.Sprintf("AVAILABLE ACCOUNTS (use these IDs for account_id in TOOL_CALL):\n%s", string(accountsJSON))
}

// loadDealsForCRUD loads deals by IDs from history.
func (s *Service) loadDealsForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
	if !allowed {
		return ""
	}

	var deals []interface{}
	seen := map[string]bool{}
	for _, did := range entityIDs["deal"] {
		d, err := s.dealRepo.FindByID(did)
		if err == nil && d != nil {
			dealOwner := ""
			if d.AssignedTo != nil {
				dealOwner = *d.AssignedTo
			}
			if s.canAccessOwner(userCtx, "deal", dealOwner) && !seen[d.ID] {
				deals = append(deals, d)
				seen[d.ID] = true
			}
		}
	}

	for _, term := range extractNamesFromHistory(messageLower) {
		results, _, err := s.dealRepo.List(&pipeline.ListDealsRequest{
			Page:          1,
			PerPage:       5,
			Search:        term,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
		})
		if err != nil {
			continue
		}
		for i := range results {
			d := results[i]
			if seen[d.ID] {
				continue
			}
			dealOwner := ""
			if d.AssignedTo != nil {
				dealOwner = *d.AssignedTo
			}
			if !s.canAccessOwner(userCtx, "deal", dealOwner) {
				continue
			}
			dealCopy := d
			deals = append(deals, &dealCopy)
			seen[d.ID] = true
			if len(deals) >= 10 {
				break
			}
		}
		if len(deals) >= 10 {
			break
		}
	}

	if len(deals) == 0 {
		return ""
	}
	dealsJSON, _ := json.Marshal(deals)
	return fmt.Sprintf("AVAILABLE DEALS (use these IDs for deal_id or id in TOOL_CALL):\n%s", string(dealsJSON))
}

func (s *Service) loadTasksForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("task", userID, userCtx)
	if !allowed || s.taskRepo == nil {
		return ""
	}

	var tasks []interface{}
	seen := map[string]bool{}
	for _, taskID := range entityIDs["task"] {
		taskEntity, err := s.taskRepo.FindByID(taskID)
		if err != nil || taskEntity == nil {
			continue
		}
		taskOwner := ""
		if taskEntity.AssignedTo != nil {
			taskOwner = *taskEntity.AssignedTo
		}
		if s.canAccessOwner(userCtx, "task", taskOwner) && !seen[taskEntity.ID] {
			tasks = append(tasks, taskEntity)
			seen[taskEntity.ID] = true
		}
	}

	searchTerms := []string{}
	if search := taskSearchTermFromMessage(messageLower); search != "" {
		searchTerms = append(searchTerms, search)
	}
	searchTerms = append(searchTerms, extractNamesFromHistory(messageLower)...)

	for _, term := range uniqueNonEmpty(searchTerms) {
		results, _, err := s.taskRepo.List(&task.ListTasksRequest{
			Page:          1,
			PerPage:       10,
			Search:        term,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
		})
		if err != nil || len(results) == 0 {
			results, _, err = s.taskRepo.List(&task.ListTasksRequest{
				Page:          1,
				PerPage:       100,
				ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
			})
			if err != nil {
				continue
			}
			results = selectBestTaskMatches(results, []string{term})
		}
		for i := range results {
			taskEntity := results[i]
			if seen[taskEntity.ID] {
				continue
			}
			taskOwner := ""
			if taskEntity.AssignedTo != nil {
				taskOwner = *taskEntity.AssignedTo
			}
			if !s.canAccessOwner(userCtx, "task", taskOwner) {
				continue
			}
			taskCopy := taskEntity
			tasks = append(tasks, &taskCopy)
			seen[taskEntity.ID] = true
			if len(tasks) >= 10 {
				break
			}
		}
		if len(tasks) >= 10 {
			break
		}
	}

	if len(tasks) == 0 {
		return ""
	}
	tasksJSON, _ := json.Marshal(tasks)
	return fmt.Sprintf("AVAILABLE TASKS (use these IDs for task_id or id in TOOL_CALL):\n%s", string(tasksJSON))
}

func (s *Service) loadSchedulesForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("schedule", userID, userCtx)
	if !allowed || s.scheduleService == nil {
		return ""
	}

	var schedules []scheduledomain.ScheduleResponse
	seen := map[string]bool{}
	for _, scheduleID := range entityIDs["schedule"] {
		scheduleEntity, err := s.scheduleService.GetScheduleByID(scheduleID)
		if err != nil || scheduleEntity == nil {
			continue
		}
		if s.canAccessOwner(userCtx, "schedule", scheduleEntity.UserID) && !seen[scheduleEntity.ID] {
			schedules = append(schedules, *scheduleEntity)
			seen[scheduleEntity.ID] = true
		}
	}

	searchTerms := []string{}
	if strings.Contains(messageLower, "meeting") {
		searchTerms = append(searchTerms, "meeting")
	}
	if strings.Contains(messageLower, "rapat") {
		searchTerms = append(searchTerms, "rapat")
	}
	searchTerms = append(searchTerms, extractNamesFromHistory(messageLower)...)

	for _, term := range uniqueNonEmpty(searchTerms) {
		results, _, err := s.scheduleService.ListSchedules(&scheduledomain.ListSchedulesRequest{
			Page:          1,
			PerPage:       10,
			Search:        term,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"),
		})
		if err != nil || len(results) == 0 {
			results, _, err = s.scheduleService.ListSchedules(&scheduledomain.ListSchedulesRequest{
				Page:          1,
				PerPage:       50,
				ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"),
			})
			if err != nil {
				continue
			}
			results = selectBestScheduleMatches(results, []string{term})
		}
		for _, scheduleEntity := range results {
			if seen[scheduleEntity.ID] || !s.canAccessOwner(userCtx, "schedule", scheduleEntity.UserID) {
				continue
			}
			schedules = append(schedules, scheduleEntity)
			seen[scheduleEntity.ID] = true
			if len(schedules) >= 10 {
				break
			}
		}
		if len(schedules) >= 10 {
			break
		}
	}

	if len(schedules) == 0 {
		return ""
	}
	schedulesJSON, _ := json.Marshal(schedules)
	return fmt.Sprintf("AVAILABLE SCHEDULES (use these IDs for schedule_id or id in TOOL_CALL):\n%s", string(schedulesJSON))
}

// loadLeadStatusesForCRUD loads leads by IDs and available lead statuses.
func (s *Service) loadLeadStatusesForCRUD(messageLower string, entityIDs map[string][]string, userID string, userCtx *domainauth.UserContext) string {
	allowed, _ := s.checkDataPrivacy("lead", userID, userCtx)
	if !allowed {
		return ""
	}

	var parts []string

	// Load leads from history
	var leads []interface{}
	seen := map[string]bool{}
	for _, lid := range entityIDs["lead"] {
		l, err := s.leadRepo.FindByID(lid)
		if err == nil && l != nil {
			leadOwner := ""
			if l.AssignedTo != nil {
				leadOwner = *l.AssignedTo
			}
			if s.canAccessOwner(userCtx, "lead", leadOwner) && !seen[l.ID] {
				leads = append(leads, l)
				seen[l.ID] = true
			}
		}
	}
	for _, term := range extractNamesFromHistory(messageLower) {
		results, _, err := s.leadRepo.List(&lead.ListLeadsRequest{
			Page:          1,
			PerPage:       5,
			Search:        term,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "lead"),
		})
		if err != nil {
			continue
		}
		for i := range results {
			l := results[i]
			leadOwner := ""
			if l.AssignedTo != nil {
				leadOwner = *l.AssignedTo
			}
			if seen[l.ID] || !s.canAccessOwner(userCtx, "lead", leadOwner) {
				continue
			}
			leadCopy := l
			leads = append(leads, &leadCopy)
			seen[l.ID] = true
			if len(leads) >= 10 {
				break
			}
		}
		if len(leads) >= 10 {
			break
		}
	}
	if len(leads) > 0 {
		leadsJSON, _ := json.Marshal(leads)
		parts = append(parts, fmt.Sprintf("AVAILABLE LEADS:\n%s", string(leadsJSON)))
	}

	if s.leadStatusRepo != nil {
		if statuses, err := s.leadStatusRepo.ListAll(); err == nil && len(statuses) > 0 {
			statusesJSON, _ := json.Marshal(statuses)
			parts = append(parts, fmt.Sprintf("AVAILABLE LEAD STATUSES (use these IDs for lead_status_id, or use code values like new/contacted/interested/qualified/proposal_sent/converted/lost):\n%s", string(statusesJSON)))
		}
	}

	return strings.Join(parts, "\n")
}

// loadPipelineStagesForCRUD loads all pipeline stages so the LLM can reference
// correct stage IDs when creating deals or moving deal stages.
func (s *Service) loadPipelineStagesForCRUD() string {
	if s.pipelineService == nil {
		return ""
	}
	stages, err := s.pipelineService.ListStages(&pipeline.ListPipelineStagesRequest{})
	if err != nil || len(stages) == 0 {
		return ""
	}
	stagesJSON, _ := json.Marshal(stages)
	return fmt.Sprintf("AVAILABLE PIPELINE STAGES (use these IDs for stage_id in TOOL_CALL):\n%s", string(stagesJSON))
}

// extractNamesFromHistory extracts probable person/entity names from the
// user's current message. Splits on common keywords and returns lowercase
// name fragments that can be fuzzy-matched against DB records.
func extractNamesFromHistory(messageLower string) []string {
	var names []string

	// Common patterns: "untuk Dr Maria", "for Dr. Sari", "kontak Ahmad"
	patterns := []string{
		"untuk ", "for ", "pada ", "kontak ", "contact ",
		"leads ", "lead ", "deal ", "account ", "akun ", "brick ",
		"prospect ", "propect ", "prospek ", "customer ", "pelanggan ",
		"rs ", "rs. ", "rumah sakit ",
		"sales ", "jadwal ", "schedule ", "task ",
		"dokter ", "dr ", "dr. ",
		"follow up ", "follow-up ", "followup ",
		"hubungi ", "call ",
	}

	for _, p := range patterns {
		idx := strings.Index(messageLower, p)
		if idx < 0 {
			continue
		}
		rest := strings.TrimSpace(messageLower[idx+len(p):])
		// Take words until we hit a stop word or punctuation
		words := strings.Fields(rest)
		var nameWords []string
		for _, w := range words {
			cleaned := strings.TrimRight(w, ".,;:!?")
			if isStopWord(cleaned) {
				break
			}
			nameWords = append(nameWords, cleaned)
			if len(nameWords) >= 4 {
				break
			}
		}
		if len(nameWords) > 0 {
			names = append(names, strings.Join(nameWords, " "))
		}
	}

	return names
}

// isStopWord returns true for common Indonesian/English stop words that
// should terminate name extraction.
func isStopWord(w string) bool {
	stops := map[string]bool{
		"dan": true, "dan,": true, "atau": true, "di": true, "ke": true,
		"dari": true, "yang": true, "untuk": true, "dengan": true,
		"and": true, "or": true, "the": true, "to": true, "for": true,
		"with": true, "in": true, "on": true, "at": true,
		"buat": true, "tambah": true, "create": true, "update": true,
		"ubah": true, "jadi": true, "menjadi": true, "status": true,
		"memiliki": true, "punya": true, "stage": true, "target": true,
		"nilai": true, "value": true,
		"segera": true, "minggu": true, "bulan": true, "hari": true,
		"ini": true, "itu": true, "nya": true,
	}
	return stops[w]
}

// appendUnique adds val to slice only if not already present.
func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}

type aiDateRange struct {
	Start     string
	End       string
	Period    string
	Label     string
	HasFilter bool
}

var isoDatePattern = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)
var localizedDatePattern = regexp.MustCompile(`\b(\d{1,2})[/-](\d{1,2})[/-](\d{4})\b`)
var yearRangePattern = regexp.MustCompile(`\b(20\d{2})\s*(?:-|sampai|hingga|to|s/d|sd)\s*(20\d{2})\b`)
var startYearToNowPattern = regexp.MustCompile(`(?:dari\s+)?(?:awal\s+)?(?:tahun\s+)?(20\d{2})\s*(?:sampai|hingga|ke|to|s/d|sd)\s*(?:saat\s+ini|sekarang|hari\s+ini|now|today)\b`)
var monthToNowPattern = regexp.MustCompile(`(?i)(?:dari\s+)?(?:bulan\s+)?(january|januari|february|februari|march|maret|april|may|mei|june|juni|july|juli|august|agustus|september|october|oktober|november|december|desember)\s+(20\d{2})\s*(?:-|sampai|hingga|to|s/d|sd)\s*(?:bulan\s+)?(?:ini|saat\s+ini|sekarang|today|now)(?:\s*(20\d{2}))?\b`)
var previousYearsToNowPattern = regexp.MustCompile(`(?:dari\s+)?(?:(\d{1,2}|satu|dua|tiga|empat|lima|enam|tujuh|delapan|sembilan|sepuluh)\s+)?tahun\s+(?:kemarin|lalu)\s*(?:sampai|hingga|ke|to|s/d|sd)\s*(?:saat\s+ini|sekarang|hari\s+ini|now|today)\b`)
var explicitYearPattern = regexp.MustCompile(`\btahun\s+(20\d{2})\b|\b(20\d{2})\b`)
var relativeYearsPattern = regexp.MustCompile(`\b(\d{1,2}|satu|dua|tiga|empat|lima|enam|tujuh|delapan|sembilan|sepuluh)\s+tahun\s+(?:terakhir|kebelakang|ke belakang|belakangan|kemarin|lalu|last)\b`)

func parseAIDateRange(messageLower string, now time.Time) aiDateRange {
	parsedDates := isoDatePattern.FindAllString(messageLower, -1)
	if len(parsedDates) >= 2 {
		return aiDateRange{
			Start:     parsedDates[0],
			End:       parsedDates[1],
			Label:     parsedDates[0] + " s/d " + parsedDates[1],
			HasFilter: true,
		}
	}
	if localizedDates := localizedDatePattern.FindAllStringSubmatch(messageLower, -1); len(localizedDates) >= 2 {
		start, startOK := normalizeLocalizedDate(localizedDates[0])
		end, endOK := normalizeLocalizedDate(localizedDates[1])
		if startOK && endOK {
			return aiDateRange{
				Start:     start,
				End:       end,
				Label:     start + " s/d " + end,
				HasFilter: true,
			}
		}
	}
	if matches := yearRangePattern.FindStringSubmatch(messageLower); len(matches) == 3 {
		startYear, _ := strconv.Atoi(matches[1])
		endYear, _ := strconv.Atoi(matches[2])
		if startYear > endYear {
			startYear, endYear = endYear, startYear
		}
		start := time.Date(startYear, time.January, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(endYear, time.December, 31, 0, 0, 0, 0, now.Location())
		if end.After(now) {
			end = now
		}
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     fmt.Sprintf("%d s/d %d", startYear, endYear),
			HasFilter: true,
		}
	}
	if matches := startYearToNowPattern.FindStringSubmatch(messageLower); len(matches) == 2 {
		startYear, _ := strconv.Atoi(matches[1])
		start := time.Date(startYear, time.January, 1, 0, 0, 0, 0, now.Location())
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       now.Format("2006-01-02"),
			Label:     fmt.Sprintf("awal %d sampai saat ini", startYear),
			HasFilter: true,
		}
	}
	if matches := monthToNowPattern.FindStringSubmatch(messageLower); len(matches) >= 3 {
		monthName := strings.ToLower(matches[1])
		year, _ := strconv.Atoi(matches[2])
		monthMap := map[string]time.Month{
			"january": time.January, "januari": time.January,
			"february": time.February, "februari": time.February,
			"march": time.March, "maret": time.March,
			"april": time.April,
			"may":   time.May, "mei": time.May,
			"june": time.June, "juni": time.June,
			"july": time.July, "juli": time.July,
			"august": time.August, "agustus": time.August,
			"september": time.September,
			"october":   time.October, "oktober": time.October,
			"november": time.November,
			"december": time.December, "desember": time.December,
		}
		if month, ok := monthMap[monthName]; ok {
			end := now
			if len(matches) == 4 && strings.TrimSpace(matches[3]) != "" {
				if parsedEndYear, err := strconv.Atoi(matches[3]); err == nil {
					candidateEnd := time.Date(parsedEndYear, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
					if candidateEnd.Before(end) {
						end = candidateEnd
					}
				}
			}
			start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
			return aiDateRange{
				Start:     start.Format("2006-01-02"),
				End:       end.Format("2006-01-02"),
				Label:     fmt.Sprintf("%s %d sampai saat ini", strings.Title(monthName), year),
				HasFilter: true,
			}
		}
	}
	if matches := previousYearsToNowPattern.FindStringSubmatch(messageLower); len(matches) == 2 {
		years := 1
		if matches[1] != "" {
			years = parseIndonesianCount(matches[1])
		}
		if years <= 0 {
			years = 1
		}
		startYear := now.Year() - years
		start := time.Date(startYear, time.January, 1, 0, 0, 0, 0, now.Location())
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       now.Format("2006-01-02"),
			Label:     fmt.Sprintf("awal %d sampai saat ini", startYear),
			HasFilter: true,
		}
	}
	if matches := relativeYearsPattern.FindStringSubmatch(messageLower); len(matches) == 2 {
		years := parseIndonesianCount(matches[1])
		if years > 1 {
			start := now.AddDate(-years, 0, 0)
			return aiDateRange{
				Start:     start.Format("2006-01-02"),
				End:       now.Format("2006-01-02"),
				Label:     fmt.Sprintf("%d tahun terakhir", years),
				HasFilter: true,
			}
		}
	}

	if strings.Contains(messageLower, "hari ini") || strings.Contains(messageLower, "today") {
		date := now.Format("2006-01-02")
		return aiDateRange{Start: date, End: date, Period: "today", Label: "hari ini", HasFilter: true}
	}
	if strings.Contains(messageLower, "minggu ini") || strings.Contains(messageLower, "this week") {
		return aiDateRange{Period: "week", Label: "minggu ini", HasFilter: true}
	}
	if strings.Contains(messageLower, "bulan ini") || strings.Contains(messageLower, "this month") {
		return aiDateRange{Period: "month", Label: "bulan ini", HasFilter: true}
	}
	if strings.Contains(messageLower, "tahun ini") || strings.Contains(messageLower, "this year") {
		return aiDateRange{Period: "year", Label: "tahun ini", HasFilter: true}
	}
	if (strings.Contains(messageLower, "satu tahun") || strings.Contains(messageLower, "1 tahun") || strings.Contains(messageLower, "12 bulan")) &&
		(strings.Contains(messageLower, "kebelakang") || strings.Contains(messageLower, "ke belakang") || strings.Contains(messageLower, "terakhir") || strings.Contains(messageLower, "last")) {
		start := now.AddDate(-1, 0, 0)
		end := now
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     "12 bulan terakhir",
			HasFilter: true,
		}
	}
	if strings.Contains(messageLower, "tahun lalu") || strings.Contains(messageLower, "tahun kemarin") || strings.Contains(messageLower, "last year") {
		lastYear := now.AddDate(-1, 0, 0).Year()
		start := time.Date(lastYear, time.January, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(lastYear, time.December, 31, 0, 0, 0, 0, now.Location())
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     fmt.Sprintf("tahun lalu (%d)", lastYear),
			HasFilter: true,
		}
	}

	if strings.Contains(messageLower, "minggu lalu") || strings.Contains(messageLower, "last week") {
		start := now.AddDate(0, 0, -7)
		end := now.AddDate(0, 0, -1)
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     "minggu lalu",
			HasFilter: true,
		}
	}
	if strings.Contains(messageLower, "bulan lalu") || strings.Contains(messageLower, "last month") {
		lastMonth := now.AddDate(0, -1, 0)
		start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     "bulan lalu",
			HasFilter: true,
		}
	}

	monthMap := map[string]time.Month{
		"january": time.January, "januari": time.January,
		"february": time.February, "februari": time.February,
		"march": time.March, "maret": time.March,
		"april": time.April,
		"may":   time.May, "mei": time.May,
		"june": time.June, "juni": time.June,
		"july": time.July, "juli": time.July,
		"august": time.August, "agustus": time.August,
		"september": time.September,
		"october":   time.October, "oktober": time.October,
		"november": time.November,
		"december": time.December, "desember": time.December,
	}
	for name, month := range monthMap {
		if !strings.Contains(messageLower, name) {
			continue
		}
		year := now.Year()
		yearPattern := regexp.MustCompile(name + `\s+(\d{4})`)
		if matches := yearPattern.FindStringSubmatch(messageLower); len(matches) == 2 {
			if parsedYear, err := strconv.Atoi(matches[1]); err == nil {
				year = parsedYear
			}
		}
		start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return aiDateRange{
			Start:     start.Format("2006-01-02"),
			End:       end.Format("2006-01-02"),
			Label:     fmt.Sprintf("%s %d", strings.Title(name), year),
			HasFilter: true,
		}
	}

	if !localizedDatePattern.MatchString(messageLower) && len(parsedDates) == 0 {
		if matches := explicitYearPattern.FindStringSubmatch(messageLower); len(matches) == 3 {
			yearText := matches[1]
			if yearText == "" {
				yearText = matches[2]
			}
			if year, err := strconv.Atoi(yearText); err == nil {
				start := time.Date(year, time.January, 1, 0, 0, 0, 0, now.Location())
				end := time.Date(year, time.December, 31, 0, 0, 0, 0, now.Location())
				if end.After(now) {
					end = now
				}
				return aiDateRange{
					Start:     start.Format("2006-01-02"),
					End:       end.Format("2006-01-02"),
					Label:     fmt.Sprintf("tahun %d", year),
					HasFilter: true,
				}
			}
		}
	}

	return aiDateRange{}
}

func normalizeLocalizedDate(match []string) (string, bool) {
	if len(match) != 4 {
		return "", false
	}
	day, dayErr := strconv.Atoi(match[1])
	month, monthErr := strconv.Atoi(match[2])
	year, yearErr := strconv.Atoi(match[3])
	if dayErr != nil || monthErr != nil || yearErr != nil {
		return "", false
	}
	parsed := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if parsed.Year() != year || int(parsed.Month()) != month || parsed.Day() != day {
		return "", false
	}
	return parsed.Format("2006-01-02"), true
}

func parseIndonesianCount(value string) int {
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}
	switch value {
	case "satu":
		return 1
	case "dua":
		return 2
	case "tiga":
		return 3
	case "empat":
		return 4
	case "lima":
		return 5
	case "enam":
		return 6
	case "tujuh":
		return 7
	case "delapan":
		return 8
	case "sembilan":
		return 9
	case "sepuluh":
		return 10
	default:
		return 0
	}
}

func normalizeDateRangeForRequest(dr aiDateRange) (start string, end string, period string) {
	if dr.Start != "" || dr.End != "" {
		return dr.Start, dr.End, ""
	}
	return "", "", dr.Period
}

func aiDateRangeToTimes(dr aiDateRange, now time.Time) (time.Time, time.Time) {
	if dr.Start != "" || dr.End != "" {
		var startDate, endDate time.Time
		if dr.Start != "" {
			if parsed, err := time.Parse(dateFormat, dr.Start); err == nil {
				startDate = parsed
			}
		}
		if dr.End != "" {
			if parsed, err := time.Parse(dateFormat, dr.End); err == nil {
				endDate = parsed.Add(24*time.Hour - time.Nanosecond)
			}
		}
		return startDate, endDate
	}

	switch dr.Period {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start, start.Add(24*time.Hour - time.Nanosecond)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(weekday - 1))
		return start, start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return start, start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	default:
		return time.Time{}, time.Time{}
	}
}

func summarizePerformanceByBrick(items []sales_overview.SalesPerformanceListResponse, usersByID map[string]string, brickNames map[string]string) []map[string]interface{} {
	type summary struct {
		BrickID            string
		BrickName          string
		SalesCount         int
		TotalRevenue       int64
		DealsClosed        int
		VisitsCompleted    int
		TasksCompleted     int
		TargetAmount       int64
		TotalProspects     int
		WonProspects       int
		LostProspects      int
		ConversionRateSum  float64
		AchievementRateSum float64
	}

	byBrick := map[string]*summary{}
	for _, item := range items {
		brickID := usersByID[item.UserID]
		if brickID == "" {
			brickID = "unassigned"
		}
		brickName := brickNames[brickID]
		if brickName == "" {
			brickName = "Unassigned"
		}
		entry := byBrick[brickID]
		if entry == nil {
			entry = &summary{BrickID: brickID, BrickName: brickName}
			byBrick[brickID] = entry
		}
		entry.SalesCount++
		entry.TotalRevenue += item.TotalRevenue
		entry.DealsClosed += item.DealsClosed
		entry.VisitsCompleted += item.VisitsCompleted
		entry.TasksCompleted += item.TasksCompleted
		entry.TargetAmount += item.TargetAmount
		entry.TotalProspects += item.TotalProspects
		entry.WonProspects += item.WonProspects
		entry.LostProspects += item.LostProspects
		entry.ConversionRateSum += item.ConversionRate
		entry.AchievementRateSum += item.TargetAchievementPercentage
	}

	results := make([]map[string]interface{}, 0, len(byBrick))
	for _, item := range byBrick {
		avgConversion := 0.0
		avgAchievement := 0.0
		if item.SalesCount > 0 {
			avgConversion = item.ConversionRateSum / float64(item.SalesCount)
			avgAchievement = item.AchievementRateSum / float64(item.SalesCount)
		}
		results = append(results, map[string]interface{}{
			"brick_id":                        item.BrickID,
			"brick_name":                      item.BrickName,
			"sales_count":                     item.SalesCount,
			"total_revenue":                   item.TotalRevenue,
			"deals_closed":                    item.DealsClosed,
			"visits_completed":                item.VisitsCompleted,
			"tasks_completed":                 item.TasksCompleted,
			"target_amount":                   item.TargetAmount,
			"total_prospects":                 item.TotalProspects,
			"won_prospects":                   item.WonProspects,
			"lost_prospects":                  item.LostProspects,
			"average_conversion_rate":         avgConversion,
			"average_target_achievement_rate": avgAchievement,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		left, _ := results[i]["total_revenue"].(int64)
		right, _ := results[j]["total_revenue"].(int64)
		return left > right
	})

	return results
}

func (s *Service) buildWonLostDealsContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data deals/pipeline tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}

	dateRange := parseAIDateRange(messageLower, time.Now())
	start, end, _ := normalizeDateRangeForRequest(dateRange)
	label := "periode aktif"
	if dateRange.Label != "" {
		label = dateRange.Label
	}

	baseReq := pipeline.ListDealsRequest{
		Page:          1,
		PerPage:       50,
		DateFrom:      start,
		DateTo:        end,
		ScopedUserIDs: s.scopedUserIDs(userCtx, "deal"),
	}

	wonReq := baseReq
	wonReq.Status = "won"
	wonDeals, wonTotal, wonErr := s.dealRepo.List(&wonReq)
	if wonErr != nil {
		return "", "⚠️ Data deal won tidak dapat diakses saat ini."
	}

	lostReq := baseReq
	lostReq.Status = "lost"
	lostDeals, lostTotal, lostErr := s.dealRepo.List(&lostReq)
	if lostErr != nil {
		return "", "⚠️ Data deal lost tidak dapat diakses saat ini."
	}

	payload := map[string]interface{}{
		"period_label": label,
		"start_date":   start,
		"end_date":     end,
		"won_summary": map[string]interface{}{
			"total": wonTotal,
			"items": s.formatDealsForAI(wonDeals),
		},
		"lost_summary": map[string]interface{}{
			"total": lostTotal,
			"items": s.formatDealsForAI(lostDeals),
		},
	}

	raw, _ := json.Marshal(payload)
	return fmt.Sprintf("REAL WON/LOST DEALS DATA:\n%s\n\nPresent won and lost deals separately for the requested period. Use Markdown tables, include totals, highlight trends, and use [Title](deal://id) for deal links.", string(raw)), ""
}

func (s *Service) buildReportContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	if s.reportService == nil {
		return "", "⚠️ Layanan reports belum diinisialisasi untuk AI Assistant."
	}

	allowed, _ := s.checkDataPrivacy("report", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke reports tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}

	dateRange := parseAIDateRange(messageLower, time.Now())
	start, end, _ := normalizeDateRangeForRequest(dateRange)
	req := &reportdomain.ReportRequest{
		StartDate: start,
		EndDate:   end,
		Limit:     25,
	}

	label := "periode aktif"
	if dateRange.Label != "" {
		label = dateRange.Label
	}

	payload := map[string]interface{}{
		"period_label": label,
		"start_date":   start,
		"end_date":     end,
	}

	switch {
	case strings.Contains(messageLower, "visit") || strings.Contains(messageLower, "kunjungan"):
		req.ScopedUserIDs = s.scopedUserIDs(userCtx, "visit_report")
		visitReport, err := s.reportService.GetVisitReportReport(req)
		if err != nil {
			return "", "⚠️ Data laporan kunjungan tidak tersedia atau tidak dapat diakses saat ini."
		}
		payload["report_type"] = "visit_report"
		payload["visit_report"] = visitReport
	case strings.Contains(messageLower, "sales performance") ||
		strings.Contains(messageLower, "performa sales") ||
		strings.Contains(messageLower, "performa penjualan") ||
		strings.Contains(messageLower, "sales") ||
		strings.Contains(messageLower, "penjualan"):
		req.ScopedUserIDs = s.scopedUserIDs(userCtx, "sales_performance")
		performanceReport, err := s.reportService.GetSalesPerformanceReport(req)
		if err != nil {
			return "", "⚠️ Data laporan sales performance tidak tersedia atau tidak dapat diakses saat ini."
		}
		payload["report_type"] = "sales_performance"
		payload["sales_performance"] = performanceReport
	case strings.Contains(messageLower, "pipeline") ||
		strings.Contains(messageLower, "funnel") ||
		strings.Contains(messageLower, "deal"):
		req.ScopedUserIDs = s.scopedUserIDs(userCtx, "deal")
		pipelineReq := *req
		pipelineReq.EntityType = "deal"
		pipelineReport, err := s.reportService.GetPipelineReport(&pipelineReq)
		if err != nil {
			return "", "⚠️ Data laporan pipeline/deals tidak tersedia atau tidak dapat diakses saat ini."
		}
		payload["report_type"] = "pipeline"
		payload["pipeline"] = pipelineReport
	default:
		req.ScopedUserIDs = s.scopedUserIDs(userCtx, "visit_report")
		visitReport, visitErr := s.reportService.GetVisitReportReport(req)
		pipelineReq := *req
		pipelineReq.EntityType = "deal"
		pipelineReq.ScopedUserIDs = s.scopedUserIDs(userCtx, "deal")
		pipelineReport, pipelineErr := s.reportService.GetPipelineReport(&pipelineReq)
		req.ScopedUserIDs = s.scopedUserIDs(userCtx, "sales_performance")
		performanceReport, performanceErr := s.reportService.GetSalesPerformanceReport(req)
		if visitErr != nil && pipelineErr != nil && performanceErr != nil {
			return "", "⚠️ Data reports tidak tersedia atau tidak dapat diakses saat ini."
		}
		payload["report_type"] = "summary"
		if visitErr == nil {
			payload["visit_report"] = visitReport
		}
		if pipelineErr == nil {
			payload["pipeline"] = pipelineReport
		}
		if performanceErr == nil {
			payload["sales_performance"] = performanceReport
		}
	}

	raw, _ := json.Marshal(payload)
	return fmt.Sprintf("REAL REPORT DATA:\n%s\n\nGenerate a concise report in Markdown using ONLY this data. Include period, key totals, notable breakdowns, and recommended next actions. Do not invent rows or values.", string(raw)), ""
}

func (s *Service) buildUserManagementContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	allowed, _ := s.checkDataPrivacy("user", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data users tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}
	if s.userRepo == nil {
		return "", "⚠️ Layanan user management belum tersedia untuk AI Assistant."
	}

	req := &userdomain.ListUsersRequest{
		Page:          1,
		PerPage:       20,
		ScopedUserIDs: s.scopedUserIDs(userCtx, "user"),
	}
	if strings.Contains(messageLower, "inactive") || strings.Contains(messageLower, "nonaktif") {
		req.Status = "inactive"
	} else if strings.Contains(messageLower, "active") || strings.Contains(messageLower, "aktif") {
		req.Status = "active"
	}

	roleFilter := userRoleFilterTerm(messageLower)
	if roleFilter != "" {
		if s.roleRepo == nil {
			return "", "⚠️ Layanan role management belum tersedia untuk memfilter data users berdasarkan role."
		}
		roleID, ok := s.resolveRoleIDForUserFilter(roleFilter)
		if !ok {
			return "", fmt.Sprintf("Tidak ada role yang cocok dengan **%s**.", roleFilter)
		}
		req.RoleID = roleID
	}
	if search := userManagementSearchTerm(messageLower); search != "" && roleFilter == "" {
		req.Search = search
	}

	users, total, err := s.userRepo.List(req)
	if err != nil {
		return "", "⚠️ Tidak dapat mengakses data users dari database. Data mungkin tidak tersedia."
	}
	if len(users) == 0 {
		return "", "Tidak ada data users sesuai akses dan filter yang diminta."
	}

	usersJSON, _ := json.Marshal(s.formatUsersForAI(users))
	return fmt.Sprintf("REAL USERS DATA (showing %d of %d users, scoped to logged-in user's RBAC):\n%s\n\nPresent in a Markdown table. Show name, email, role, group, status, and created_at. Never show password, token, or internal credential fields. Never invent users. Include a navigate action card to /master-data/users when useful.", len(users), total, string(usersJSON)), ""
}

func (s *Service) buildRoleManagementContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	allowed, _ := s.checkDataPrivacy("role", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data roles tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}
	if s.roleRepo == nil {
		return "", "⚠️ Layanan role management belum tersedia untuk AI Assistant."
	}

	roles, err := s.roleRepo.List()
	if err != nil {
		return "", "⚠️ Tidak dapat mengakses data roles dari database. Data mungkin tidak tersedia."
	}
	if len(roles) == 0 {
		return "", "Tidak ada data roles sesuai akses dan filter yang diminta."
	}

	filtered := make([]map[string]interface{}, 0, len(roles))
	search := roleManagementSearchTerm(messageLower)
	for _, roleEntity := range roles {
		if search != "" &&
			!strings.Contains(strings.ToLower(roleEntity.Name), search) &&
			!strings.Contains(strings.ToLower(roleEntity.Code), search) {
			continue
		}
		filtered = append(filtered, map[string]interface{}{
			"id":           roleEntity.ID,
			"name":         roleEntity.Name,
			"code":         roleEntity.Code,
			"description":  roleEntity.Description,
			"status":       roleEntity.Status,
			"is_protected": roleEntity.IsProtected,
			"user_count":   roleEntity.UserCount,
			"created_at":   roleEntity.CreatedAt,
		})
	}
	if len(filtered) == 0 {
		return "", "Tidak ada data roles sesuai akses dan filter yang diminta."
	}

	rolesJSON, _ := json.Marshal(filtered)
	return fmt.Sprintf("REAL ROLES DATA (showing %d of %d roles):\n%s\n\nPresent in a Markdown table. Show name, code, status, protected, user_count, and description. Do not show permission internals unless the user explicitly asks for role permissions. Never invent roles. Include a navigate action card to /master-data/roles when useful.", len(filtered), len(roles), string(rolesJSON)), ""
}

func userManagementSearchTerm(messageLower string) string {
	cleaned := strings.NewReplacer(
		"berikan", "",
		"berika", "",
		"beri", "",
		"jadi", "",
		"tampilkan", "",
		"lihat", "",
		"tunjukkan", "",
		"dara", "",
		"data", "",
		"users", "",
		"user", "",
		"pengguna", "",
		"daftar", "",
		"list", "",
		"yang", "",
		"ada", "",
	).Replace(messageLower)
	cleaned = normalizeManagementSearchTerm(cleaned)
	if len(cleaned) < 3 {
		return ""
	}
	if strings.Contains(cleaned, "semua") || strings.Contains(cleaned, "all") {
		return ""
	}
	return cleaned
}

func userRoleFilterTerm(messageLower string) string {
	rolePattern := regexp.MustCompile(`(?i)\b(?:role|roles|peran)\s+([a-zA-Z_ -]+)`)
	match := rolePattern.FindStringSubmatch(messageLower)
	if len(match) < 2 {
		return ""
	}

	cleaned := strings.NewReplacer(
		"yang", "",
		"ada", "",
		"aktif", "",
		"active", "",
		"inactive", "",
		"nonaktif", "",
		"status", "",
		"dengan", "",
	).Replace(strings.ToLower(match[1]))
	cleaned = normalizeManagementSearchTerm(cleaned)
	words := strings.Fields(cleaned)
	if len(words) > 2 {
		words = words[:2]
	}
	return normalizeRoleLookupTerm(strings.Join(words, " "))
}

func normalizeRoleLookupTerm(value string) string {
	value = normalizeManagementSearchTerm(strings.ToLower(value))
	value = strings.Trim(value, " .,:;-")
	switch value {
	case "sales manager", "sales-manager", "manager sales", "salesmanager":
		return "sales_manager"
	case "sales representative", "sales rep", "sales-representative", "salesrep":
		return "sales"
	}
	return strings.ReplaceAll(value, " ", "_")
}

func (s *Service) resolveRoleIDForUserFilter(roleTerm string) (string, bool) {
	if s.roleRepo == nil || roleTerm == "" {
		return "", false
	}

	if roleEntity, err := s.roleRepo.FindByCode(roleTerm); err == nil && roleEntity != nil && roleEntity.ID != "" {
		return roleEntity.ID, true
	}

	roles, err := s.roleRepo.List()
	if err != nil {
		return "", false
	}
	for _, roleEntity := range roles {
		if strings.EqualFold(roleEntity.Code, roleTerm) ||
			strings.EqualFold(roleEntity.Name, roleTerm) ||
			normalizeRoleLookupTerm(roleEntity.Name) == roleTerm {
			return roleEntity.ID, true
		}
	}
	return "", false
}

func roleManagementSearchTerm(messageLower string) string {
	cleaned := strings.NewReplacer(
		"berikan", "",
		"berika", "",
		"beri", "",
		"jadi", "",
		"tampilkan", "",
		"lihat", "",
		"tunjukkan", "",
		"dara", "",
		"data", "",
		"roles", "",
		"role", "",
		"peran", "",
		"daftar", "",
		"list", "",
		"yang", "",
		"ada", "",
	).Replace(messageLower)
	cleaned = normalizeManagementSearchTerm(cleaned)
	if len(cleaned) < 3 || strings.Contains(cleaned, "semua") || strings.Contains(cleaned, "all") {
		return ""
	}
	return cleaned
}

func (s *Service) buildGroupManagementContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	allowed, _ := s.checkDataPrivacy("group", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data groups tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}
	if s.groupRepo == nil {
		return "", "⚠️ Layanan group management belum tersedia untuk AI Assistant."
	}

	req := &groupdomain.ListGroupsRequest{
		Page:    1,
		PerPage: 20,
	}
	if strings.Contains(messageLower, "inactive") || strings.Contains(messageLower, "nonaktif") {
		req.Status = "inactive"
	} else if strings.Contains(messageLower, "active") || strings.Contains(messageLower, "aktif") {
		req.Status = "active"
	}
	if search := groupManagementSearchTerm(messageLower); search != "" {
		req.Search = search
	}

	groups, total, err := s.groupRepo.List(req)
	if err != nil {
		return "", "⚠️ Tidak dapat mengakses data groups dari database. Data mungkin tidak tersedia."
	}
	if len(groups) == 0 {
		return "", "Tidak ada data groups sesuai akses dan filter yang diminta."
	}

	formatted := make([]map[string]interface{}, 0, len(groups))
	for _, groupEntity := range groups {
		formatted = append(formatted, map[string]interface{}{
			"id":          groupEntity.ID,
			"name":        groupEntity.Name,
			"code":        groupEntity.Code,
			"description": groupEntity.Description,
			"status":      groupEntity.Status,
			"created_at":  groupEntity.CreatedAt,
		})
	}

	groupsJSON, _ := json.Marshal(formatted)
	return fmt.Sprintf("REAL GROUPS DATA (showing %d of %d groups):\n%s\n\nPresent in a Markdown table. Show name, code, status, description, and created_at. Never invent groups. Include a navigate action card to /master-data/groups when useful.", len(groups), total, string(groupsJSON)), ""
}

func groupManagementSearchTerm(messageLower string) string {
	cleaned := strings.NewReplacer(
		"berikan", "",
		"berika", "",
		"beri", "",
		"jadi", "",
		"tampilkan", "",
		"lihat", "",
		"tunjukkan", "",
		"dara", "",
		"data", "",
		"groups", "",
		"group", "",
		"grup", "",
		"daftar", "",
		"list", "",
		"yang", "",
		"ada", "",
	).Replace(messageLower)
	cleaned = normalizeManagementSearchTerm(cleaned)
	if len(cleaned) < 3 || strings.Contains(cleaned, "semua") || strings.Contains(cleaned, "all") {
		return ""
	}
	return cleaned
}

func normalizeManagementSearchTerm(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func (s *Service) buildBrickManagementContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	allowed, _ := s.checkDataPrivacy("brick_management", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data bricks/territories tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}
	if s.brickRepo == nil {
		return "", "⚠️ Layanan brick management belum tersedia untuk AI Assistant."
	}

	req := &brickdomain.ListBricksRequest{
		Page:    1,
		PerPage: 20,
	}
	if strings.Contains(messageLower, "inactive") || strings.Contains(messageLower, "nonaktif") {
		req.Status = "inactive"
	} else if strings.Contains(messageLower, "active") || strings.Contains(messageLower, "aktif") {
		req.Status = "active"
	}
	if search := brickManagementSearchTerm(messageLower); search != "" {
		req.Search = search
	}

	var bricks []brickdomain.Brick
	var total int64
	var err error

	switch {
	case userCtx == nil:
		return "", "⚠️ Konteks user tidak tersedia untuk membaca data bricks."
	case isAIGlobalRole(userCtx) || userCtx.IsGlobalScope("bricks"):
		bricks, total, err = s.brickRepo.List(req)
	case userCtx.IsTeamScope("bricks"):
		req.ManagerID = &userCtx.UserID
		bricks, total, err = s.brickRepo.List(req)
	default:
		if s.userRepo == nil {
			return "", "⚠️ Layanan user management belum tersedia untuk mengecek akses brick user."
		}
		currentUser, findErr := s.userRepo.FindByID(userCtx.UserID)
		if findErr != nil || currentUser == nil || currentUser.BrickID == nil || *currentUser.BrickID == "" {
			return "", "Tidak ada data bricks sesuai akses user Anda."
		}
		brickEntity, brickErr := s.brickRepo.FindByID(*currentUser.BrickID)
		if brickErr != nil || brickEntity == nil {
			return "", "Tidak ada data bricks sesuai akses user Anda."
		}
		bricks = []brickdomain.Brick{*brickEntity}
		total = 1
	}

	if err != nil {
		return "", "⚠️ Tidak dapat mengakses data bricks dari database. Data mungkin tidak tersedia."
	}
	if len(bricks) == 0 {
		return "", "Tidak ada data bricks sesuai akses dan filter yang diminta."
	}

	revenueByBrickID, revenuePeriodLabel := s.buildBrickRevenueMap(bricks, messageLower, userCtx)
	bricksJSON, _ := json.Marshal(s.formatBricksForAI(bricks, revenueByBrickID, revenuePeriodLabel))
	return fmt.Sprintf("REAL BRICKS/TERRITORIES DATA (showing %d of %d bricks, scoped to logged-in user's RBAC):\n%s\n\nPresent in a Markdown table. Show name, code, province, regency, manager_name, status, total_revenue_formatted, deals_closed, visits_completed, target_amount_formatted, and revenue_period_label. Use ONLY the revenue fields from this JSON for income/penghasilan questions. Never invent territory data, managers, or revenue. Include a navigate action card to /master-data/bricks when useful.", len(bricks), total, string(bricksJSON)), ""
}

type brickRevenueSummary struct {
	TotalRevenue          int64
	TotalRevenueFormatted string
	DealsClosed           int
	VisitsCompleted       int
	TasksCompleted        int
	TargetAmount          int64
	TargetAmountFormatted string
}

func (s *Service) buildBrickRevenueMap(bricks []brickdomain.Brick, messageLower string, userCtx *domainauth.UserContext) (map[string]brickRevenueSummary, string) {
	result := make(map[string]brickRevenueSummary, len(bricks))
	for _, brickEntity := range bricks {
		result[brickEntity.ID] = brickRevenueSummary{
			TotalRevenueFormatted: formatCurrencyRupiah(0),
			TargetAmountFormatted: formatCurrencyRupiah(0),
		}
	}
	if s.salesOverviewService == nil || len(bricks) == 0 {
		return result, "semua periode"
	}

	dateRange := parseAIDateRange(messageLower, time.Now())
	start, end, period := normalizeDateRangeForRequest(dateRange)
	periodLabel := "semua periode"
	if dateRange.Label != "" {
		periodLabel = dateRange.Label
	}

	for _, brickEntity := range bricks {
		req := &sales_overview.ListSalesPerformanceRequest{
			Page:          1,
			PerPage:       100,
			StartDate:     start,
			EndDate:       end,
			Period:        period,
			BrickID:       brickEntity.ID,
			SortBy:        "revenue",
			Order:         "desc",
			ScopedUserIDs: s.scopedUserIDs(userCtx, "sales_performance"),
		}
		items, _, err := s.salesOverviewService.ListSalesPerformance(req)
		if err != nil {
			continue
		}

		summary := result[brickEntity.ID]
		for _, item := range items {
			summary.TotalRevenue += item.TotalRevenue
			summary.DealsClosed += item.DealsClosed
			summary.VisitsCompleted += item.VisitsCompleted
			summary.TasksCompleted += item.TasksCompleted
			summary.TargetAmount += item.TargetAmount
		}
		summary.TotalRevenueFormatted = formatCurrencyRupiah(summary.TotalRevenue)
		summary.TargetAmountFormatted = formatCurrencyRupiah(summary.TargetAmount)
		result[brickEntity.ID] = summary
	}

	return result, periodLabel
}

func brickManagementSearchTerm(messageLower string) string {
	messageLower = strings.ToLower(messageLower)
	cleaned := strings.NewReplacer(
		"berikan", "",
		"berika", "",
		"beri", "",
		"jadi", "",
		"tampilkan", "",
		"lihat", "",
		"tunjukkan", "",
		"dara", "",
		"data", "",
		"bricks", "",
		"brick", "",
		"territory", "",
		"territories", "",
		"wilayah", "",
		"area mapping", "",
		"area", "",
		"daftar", "",
		"list", "",
		"yang", "",
		"ada", "",
		"dan", "",
		"siapa", "",
		"manager", "",
		"manajer", "",
		"serta", "",
		"penghasilan", "",
		"pendapatan", "",
		"revenue", "",
		"income", "",
		"tersebut", "",
		"berapa", "",
		"/", " ",
	).Replace(messageLower)
	cleaned = normalizeManagementSearchTerm(cleaned)
	if len(cleaned) < 3 {
		return ""
	}
	if strings.Contains(cleaned, "semua") || strings.Contains(cleaned, "all") {
		return ""
	}
	return cleaned
}

func (s *Service) buildTargetContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	if s.monthlyTargetService == nil {
		return "", "⚠️ Layanan target belum diinisialisasi untuk AI Assistant."
	}

	allowed, _ := s.checkDataPrivacy("target", userID, userCtx)
	if !allowed {
		return "", "⚠️ Akses ke data target tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
	}

	now := time.Now()
	targetMonths := targetMonthsFromMessage(messageLower, now)
	if len(targetMonths) == 0 {
		dateRange := parseAIDateRange(messageLower, now)
		start, _, _ := normalizeDateRangeForRequest(dateRange)
		targetTime := now
		if parsed, err := time.Parse(dateFormat, start); err == nil {
			targetTime = parsed
		}
		targetMonths = []time.Time{targetTime}
	}

	scope := targetScopeFromMessage(messageLower)
	ownerFilter := targetOwnerFilterFromMessage(messageLower)
	scopedUserIDs := s.scopedUserIDs(userCtx, "target")
	rows := make([]map[string]interface{}, 0)
	var grandTotalRupiah int64
	var totalRecords int

	for _, targetTime := range targetMonths {
		req := &monthlytargetdomain.ListMonthlyTargetsRequest{
			Page:          1,
			PerPage:       100,
			Year:          intPtr(targetTime.Year()),
			Month:         intPtr(int(targetTime.Month())),
			Scope:         scope,
			ScopedUserIDs: scopedUserIDs,
		}

		targets, pagination, err := s.monthlyTargetService.List(req)
		if err != nil {
			return "", "⚠️ Data target tidak tersedia atau tidak dapat diakses saat ini."
		}

		totalRecords += pagination.Total
		for _, target := range targets {
			ownerName := targetOwnerName(target, scope)
			if ownerFilter != "" && !strings.EqualFold(ownerName, ownerFilter) {
				continue
			}

			amountRupiah := target.TargetAmount / 100
			grandTotalRupiah += amountRupiah
			rows = append(rows, map[string]interface{}{
				"id":               target.ID,
				"owner":            ownerName,
				"scope":            scope,
				"user_id":          target.UserID,
				"group_id":         target.GroupID,
				"brick_id":         target.BrickID,
				"year":             target.Year,
				"month":            target.Month,
				"month_name":       targetTime.Format("January"),
				"amount_idr":       amountRupiah,
				"amount_formatted": formatCurrencyRupiah(target.TargetAmount),
			})
		}
	}

	if len(rows) == 0 {
		return "", "⚠️ Tidak ada data target untuk periode yang diminta."
	}

	payload := map[string]interface{}{
		"scope":                  scope,
		"owner_filter":           ownerFilter,
		"total_records":          len(rows),
		"total_records_in_scope": totalRecords,
		"total_amount_idr":       grandTotalRupiah,
		"total_amount_formatted": "Rp " + formatNumberRupiah(float64(grandTotalRupiah)),
		"targets":                rows,
	}
	raw, _ := json.Marshal(payload)
	return fmt.Sprintf("REAL MONTHLY TARGET DATA:\n%s\n\nPresent target data in Markdown using ONLY this JSON. Currency values are already normalized to Rupiah in `amount_idr`, `amount_formatted`, `total_amount_idr`, and `total_amount_formatted`. Never use or infer minor-unit/cents values. The `scope` field is the single source of truth for the rows; do not combine user, group, and brick targets. State that total_amount_formatted is the total for this scope and selected month(s) only. Do not invent rows, labels, amounts, unavailable months, or all-scope totals.", string(raw)), ""
}

func isDealValueTargetIntent(messageLower string) bool {
	if strings.Contains(messageLower, "target deal") ||
		strings.Contains(messageLower, "deal target") ||
		strings.Contains(messageLower, "nilai target") ||
		strings.Contains(messageLower, "target nilai") ||
		strings.Contains(messageLower, "target opportunity") ||
		strings.Contains(messageLower, "opportunity target") {
		return true
	}

	hasDealTerm := strings.Contains(messageLower, "deal") ||
		strings.Contains(messageLower, "opportunity") ||
		strings.Contains(messageLower, "pipeline") ||
		strings.Contains(messageLower, "prospect") ||
		strings.Contains(messageLower, "prospek") ||
		strings.Contains(messageLower, "propect")
	hasWriteTerm := strings.Contains(messageLower, "buat") ||
		strings.Contains(messageLower, "tambahkan") ||
		strings.Contains(messageLower, "tambah") ||
		strings.Contains(messageLower, "create") ||
		strings.Contains(messageLower, "add")
	hasMoneyTerm := strings.Contains(messageLower, "juta") ||
		strings.Contains(messageLower, "jt") ||
		strings.Contains(messageLower, "rp") ||
		strings.Contains(messageLower, "idr")

	return strings.Contains(messageLower, "target") && hasDealTerm && (hasWriteTerm || hasMoneyTerm)
}

func targetScopeFromMessage(messageLower string) string {
	normalized := strings.ToLower(messageLower)
	switch {
	case strings.Contains(normalized, "brick"),
		strings.Contains(normalized, "territory"),
		strings.Contains(normalized, "wilayah"),
		strings.Contains(normalized, "area"):
		return "brick"
	case strings.Contains(normalized, "group"),
		strings.Contains(normalized, "grup"),
		strings.Contains(normalized, "team"),
		strings.Contains(normalized, "tim"):
		return "group"
	default:
		return "user"
	}
}

func targetMonthsFromMessage(messageLower string, now time.Time) []time.Time {
	monthNames := []struct {
		names []string
		month time.Month
	}{
		{names: []string{"january", "januari"}, month: time.January},
		{names: []string{"february", "februari"}, month: time.February},
		{names: []string{"march", "maret"}, month: time.March},
		{names: []string{"april"}, month: time.April},
		{names: []string{"may", "mei"}, month: time.May},
		{names: []string{"june", "juni"}, month: time.June},
		{names: []string{"july", "juli"}, month: time.July},
		{names: []string{"august", "agustus"}, month: time.August},
		{names: []string{"september"}, month: time.September},
		{names: []string{"october", "oktober"}, month: time.October},
		{names: []string{"november"}, month: time.November},
		{names: []string{"december", "desember"}, month: time.December},
	}

	year := now.Year()
	if match := regexp.MustCompile(`\b(20\d{2}|21\d{2})\b`).FindStringSubmatch(messageLower); len(match) == 2 {
		if parsedYear, err := strconv.Atoi(match[1]); err == nil {
			year = parsedYear
		}
	}

	seen := map[time.Month]bool{}
	result := make([]time.Time, 0)
	for _, item := range monthNames {
		for _, name := range item.names {
			if !strings.Contains(messageLower, name) || seen[item.month] {
				continue
			}
			seen[item.month] = true
			result = append(result, time.Date(year, item.month, 1, 0, 0, 0, 0, now.Location()))
		}
	}

	return result
}

func targetOwnerFilterFromMessage(messageLower string) string {
	match := regexp.MustCompile(`\buntuk\s+(.+?)(?:\s+menjadi\b|\s+jadi\b|\s+ke\b|\s+sebesar\b|\s+dengan\s+nilai\b|\s+pada\s+bulan\b|\s+di\s+bulan\b|\s+bulan\b|\s+month\b|\s+tahun\b|\s+year\b|$)`).FindStringSubmatch(messageLower)
	if len(match) != 2 {
		return ""
	}

	candidate := strings.TrimSpace(match[1])
	candidate = strings.Trim(candidate, ".,;:!? ")
	candidate = strings.TrimPrefix(candidate, "data target ")
	candidate = strings.TrimPrefix(candidate, "target ")
	candidate = strings.TrimPrefix(candidate, "monthly sales ")
	candidate = strings.TrimSuffix(candidate, " pada")
	candidate = strings.TrimSuffix(candidate, " di")
	candidate = strings.TrimSpace(candidate)

	if candidate == "" ||
		candidate == "sales" ||
		candidate == "sales rep" ||
		candidate == "sales reps" ||
		strings.Contains(candidate, "semua ") ||
		strings.Contains(candidate, "seluruh ") ||
		strings.Contains(candidate, "all ") ||
		strings.Contains(candidate, "para ") {
		return ""
	}

	return candidate
}

func targetOwnerName(target monthlytargetdomain.MonthlyTargetResponse, scope string) string {
	switch scope {
	case "brick":
		if target.Brick != nil && target.Brick.Name != "" {
			return target.Brick.Name
		}
		if target.BrickID != nil {
			return *target.BrickID
		}
	case "group":
		if target.Group != nil && target.Group.Name != "" {
			return target.Group.Name
		}
		if target.GroupID != nil {
			return *target.GroupID
		}
	default:
		if target.User != nil && target.User.Name != "" {
			return target.User.Name
		}
		if target.UserID != nil {
			return *target.UserID
		}
	}
	return "Unknown"
}

func intPtr(value int) *int {
	return &value
}

func (s *Service) buildPerformanceContext(messageLower string, userID string, userCtx *domainauth.UserContext) (string, string) {
	if s.salesOverviewService == nil {
		return "", "⚠️ Layanan sales performance belum diinisialisasi."
	}

	allowed, _ := s.checkDataPrivacy("sales_performance", userID, userCtx)
	if !allowed {
		return "", "⚠️ Modul Sales Performance tidak diaktifkan di pengaturan AI atau Anda tidak memiliki akses."
	}

	dateRange := parseAIDateRange(messageLower, time.Now())
	start, end, period := normalizeDateRangeForRequest(dateRange)
	req := &sales_overview.ListSalesPerformanceRequest{
		Page:          1,
		PerPage:       100,
		StartDate:     start,
		EndDate:       end,
		Period:        period,
		SortBy:        "revenue",
		Order:         "desc",
		ScopedUserIDs: s.scopedUserIDs(userCtx, "sales_performance"),
	}

	items, total, err := s.salesOverviewService.ListSalesPerformance(req)
	if err != nil {
		return "", "⚠️ Data sales performance tidak tersedia atau tidak dapat diakses saat ini."
	}

	if len(items) == 0 {
		label := "periode yang diminta"
		if dateRange.Label != "" {
			label = dateRange.Label
		}
		return "", fmt.Sprintf("Tidak ada data sales performance untuk %s sesuai scope dan filter user login.", label)
	}

	if len(items) > 25 {
		items = items[:25]
	}

	label := "periode aktif"
	if dateRange.Label != "" {
		label = dateRange.Label
	}

	if strings.Contains(messageLower, "brick") || strings.Contains(messageLower, "wilayah") {
		usersByID := map[string]string{}
		brickIDs := make([]string, 0)
		seenBrick := map[string]bool{}
		for _, item := range items {
			if s.userRepo == nil {
				continue
			}
			userEntity, findErr := s.userRepo.FindByID(item.UserID)
			if findErr != nil || userEntity == nil || userEntity.BrickID == nil {
				continue
			}
			brickID := *userEntity.BrickID
			usersByID[item.UserID] = brickID
			if !seenBrick[brickID] {
				brickIDs = append(brickIDs, brickID)
				seenBrick[brickID] = true
			}
		}

		brickNames := map[string]string{}
		if s.brickRepo != nil && len(brickIDs) > 0 {
			if bricks, brickErr := s.brickRepo.FindByIDs(brickIDs); brickErr == nil {
				for i := range bricks {
					brickNames[bricks[i].ID] = bricks[i].Name
				}
			}
		}

		payload := map[string]interface{}{
			"period_label": label,
			"total_sales":  total,
			"items":        summarizePerformanceByBrick(items, usersByID, brickNames),
		}
		raw, _ := json.Marshal(payload)
		return fmt.Sprintf("REAL SALES PERFORMANCE BY BRICK:\n%s\n\nPresent performance grouped by brick. Show brick name, sales count, revenue, deals closed, visits, tasks, target, and achievement.", string(raw)), ""
	}

	var trend interface{}
	var yearlyTrend interface{}
	if wantsLineChart(messageLower) {
		trendMode := "monthly"
		if strings.Contains(messageLower, "bulan ini") || strings.Contains(messageLower, "minggu") || strings.Contains(messageLower, "harian") || strings.Contains(messageLower, "daily") {
			trendMode = "rolling_30d"
		}
		if overview, overviewErr := s.salesOverviewService.GetMonthlySalesOverview(start, end, trendMode, req.ScopedUserIDs); overviewErr == nil && overview != nil && len(overview.MonthlyData) > 0 {
			trend = overview
			yearlyTrend = summarizeYearlySalesTrend(overview.MonthlyData)
		}
	}

	trendSummary := map[string]interface{}{}
	if overview, ok := trend.(*sales_overview.MonthlySalesOverviewResponse); ok && len(overview.MonthlyData) > 0 {
		trendSummary["monthly_points"] = len(overview.MonthlyData)
		trendSummary["first_period_label"] = overview.MonthlyData[0].PeriodLabel
		trendSummary["last_period_label"] = overview.MonthlyData[len(overview.MonthlyData)-1].PeriodLabel
		trendSummary["first_period_start"] = overview.MonthlyData[0].PeriodStart.Format("2006-01-02")
		trendSummary["last_period_end"] = overview.MonthlyData[len(overview.MonthlyData)-1].PeriodEnd.Format("2006-01-02")
	}

	auditMode := hasSalesIssueAnalysisTerm(messageLower)
	raw, _ := json.Marshal(map[string]interface{}{
		"period_label":  label,
		"audit_mode":    auditMode,
		"total_sales":   total,
		"items":         items,
		"trend":         trend,
		"yearly_trend":  yearlyTrend,
		"trend_summary": trendSummary,
	})
	return fmt.Sprintf(`REAL SALES PERFORMANCE BY SALES REP:
%s

Present performance per sales rep. Use Markdown table, sort by revenue, and mention top and bottom performers.
If audit_mode is true, identify likely sales issues from the provided metrics only: target gaps, low achievement rate, zero/low revenue, low visit/task activity, weak conversion/deals closed, and unusual differences between sales reps. Do not invent causes; label them as hypotheses from CRM metrics and recommend concrete follow-up checks.
If trend is present and the user asks for monthly/yearly development, prioritize trend.monthly_data and yearly_trend before the sales-rep table. Use every row in trend.monthly_data for the line CHART marker, with period_label as label and total_revenue as value. INCLUDE months with 0 revenue in the chart so the user knows you checked the entire period. If there are many months with 0 revenue, explicitly state to the user that no sales were recorded in those earlier months, confirming that the full date range was scanned. Do not reduce the chart to only the latest month when monthly_data has multiple rows. Do not say chart data is unavailable when trend.monthly_data is present. If trend_summary.monthly_points is greater than 1, state that the chart covers multiple periods and do not claim only one month exists.`, string(raw)), ""
}

func wantsLineChart(messageLower string) bool {
	return (strings.Contains(messageLower, "grafik") || strings.Contains(messageLower, "chart")) &&
		(strings.Contains(messageLower, "line") || strings.Contains(messageLower, "tren") || strings.Contains(messageLower, "trend") || strings.Contains(messageLower, "nya"))
}

type aiYearlySalesTrend struct {
	Year         int   `json:"year"`
	TotalRevenue int64 `json:"total_revenue"`
	TotalDeals   int   `json:"total_deals"`
	TotalVisits  int   `json:"total_visits"`
	TotalTasks   int   `json:"total_tasks"`
	TargetAmount int64 `json:"target_amount"`
}

func summarizeYearlySalesTrend(monthlyData []sales_overview.MonthlySalesData) []aiYearlySalesTrend {
	byYear := map[int]*aiYearlySalesTrend{}
	years := make([]int, 0)
	for _, row := range monthlyData {
		summary := byYear[row.Year]
		if summary == nil {
			summary = &aiYearlySalesTrend{Year: row.Year}
			byYear[row.Year] = summary
			years = append(years, row.Year)
		}
		summary.TotalRevenue += row.TotalRevenue
		summary.TotalDeals += row.TotalDeals
		summary.TotalVisits += row.TotalVisits
		summary.TotalTasks += row.TotalTasks
		summary.TargetAmount += row.TargetAmount
	}
	sort.Ints(years)
	result := make([]aiYearlySalesTrend, 0, len(years))
	for _, year := range years {
		result = append(result, *byYear[year])
	}
	return result
}
