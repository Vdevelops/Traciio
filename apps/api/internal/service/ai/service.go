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
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	route_optimization_domain "github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	dashboardservice "github.com/gilabs/crm-healthcare/api/internal/service/dashboard"
	leadservice "github.com/gilabs/crm-healthcare/api/internal/service/lead"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	routeoptimizationservice "github.com/gilabs/crm-healthcare/api/internal/service/route_optimization"
	salesoverviewservice "github.com/gilabs/crm-healthcare/api/internal/service/sales_overview"
	scheduleservice "github.com/gilabs/crm-healthcare/api/internal/service/schedule"
	taskservice "github.com/gilabs/crm-healthcare/api/internal/service/task"
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
	brickRepo                interfaces.BrickRepository
	settingsRepo             interfaces.AISettingsRepository
	permService              *permissionservice.Service
	dashboardService         *dashboardservice.Service         // For analytics data
	routeOptimizationService *routeoptimizationservice.Service // For creating real routes
	salesOverviewService     *salesoverviewservice.Service
	// CRUD tool services
	leadService     *leadservice.Service
	taskService     *taskservice.Service
	pipelineService *pipelineservice.Service
	scheduleService *scheduleservice.Service
	apiKey          string
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
	brickRepo interfaces.BrickRepository,
	settingsRepo interfaces.AISettingsRepository,
	permService *permissionservice.Service,
	dashboardService *dashboardservice.Service,
	routeOptimizationService *routeoptimizationservice.Service,
	salesOverviewService *salesoverviewservice.Service,
	leadSvc *leadservice.Service,
	taskSvc *taskservice.Service,
	pipelineSvc *pipelineservice.Service,
	scheduleSvc *scheduleservice.Service,
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
		brickRepo:                brickRepo,
		settingsRepo:             settingsRepo,
		permService:              permService,
		dashboardService:         dashboardService,
		routeOptimizationService: routeOptimizationService,
		salesOverviewService:     salesOverviewService,
		leadService:              leadSvc,
		taskService:              taskSvc,
		pipelineService:          pipelineSvc,
		scheduleService:          scheduleSvc,
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
	var dataPrivacy ai_settings.DataPrivacySettings
	if settings.DataPrivacy != nil {
		if err := json.Unmarshal(settings.DataPrivacy, &dataPrivacy); err != nil {
			return true, nil // Default to allow if parsing fails
		}
	} else {
		// Default: allow all data types and enable all modules
		dataPrivacy = ai_settings.DataPrivacySettings{
			// Sales domain
			AllowLeads:        true,
			AllowDeals:        true,
			AllowVisitReports: true,
			AllowActivities:   true,
			AllowTasks:        true,
			AllowSchedule:     true,
			AllowPipelines:    true,
			// Customer domain
			AllowAccounts: true,
			AllowContacts: true,
			// Inventory domain
			AllowProducts: true,
			// Analytics domain
			AllowSalesPerformance: true,
			AllowProductAnalysis:  true,
			AllowReports:          true,
			// Management domain
			AllowUsers:           true,
			AllowRoles:           true,
			AllowGroups:          true,
			AllowBrickManagement: true,
			AllowTarget:          true,
			// Route Optimization domain
			AllowRouteOptimization: true,
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
	if strings.Contains(messageLower, "data privacy") || strings.Contains(messageLower, "privacy") ||
		strings.Contains(messageLower, "data privasi") || strings.Contains(messageLower, "privasi") ||
		strings.Contains(messageLower, "akses data") || strings.Contains(messageLower, "data yang bisa") {
		// Get data privacy settings
		var dataPrivacy ai_settings.DataPrivacySettings

		if settings.DataPrivacy != nil {
			if err := json.Unmarshal(settings.DataPrivacy, &dataPrivacy); err == nil {
				privacyInfo := buildPrivacyInfoMessage(dataPrivacy, false)

				// Return direct response about data privacy settings
				return &ai.ChatResponse{
					Message: privacyInfo + "\n\nIni adalah pengaturan data privacy dan modul analytics yang aktif di sistem. Anda dapat mengubah pengaturan ini melalui halaman AI Settings.",
					Tokens:  0, // No tokens consumed for this internal response
				}, nil
			}
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
			Message: privacyInfo + "\n\nAnda dapat mengatur data privacy dan enable/disable modul analytics melalui halaman AI Settings.",
			Tokens:  0,
		}, nil
	}

	// Load context data if provided
	var contextData string
	var dataAccessInfo string

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

		// ── CRUD INTENT DETECTION (HIGHEST PRIORITY) ──────────────────────────
		// When the user explicitly asks to create/update an entity, we must
		// enrich the context with relevant entity data so the LLM can emit a
		// proper TOOL_CALL with real IDs instead of asking the user for them.
		isCRUDIntent := strings.Contains(messageLower, "buat") ||
			strings.Contains(messageLower, "buatkan") ||
			strings.Contains(messageLower, "tambah") ||
			strings.Contains(messageLower, "create") ||
			strings.Contains(messageLower, "add") ||
			strings.Contains(messageLower, "update") ||
			strings.Contains(messageLower, "ubah") ||
			strings.Contains(messageLower, "pindah") ||
			strings.Contains(messageLower, "jadwalkan") ||
			strings.Contains(messageLower, "ya") && (strings.Contains(messageLower, "buat") || strings.Contains(messageLower, "follow")) ||
			strings.Contains(messageLower, "follow up") ||
			strings.Contains(messageLower, "follow-up") ||
			strings.Contains(messageLower, "followup")

		if isCRUDIntent && contextData == "" {
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
			(strings.Contains(messageLower, "performa") && strings.Contains(messageLower, "brick")) ||
			(strings.Contains(messageLower, "performance") && strings.Contains(messageLower, "brick")) ||
			(strings.Contains(messageLower, "revenue") && strings.Contains(messageLower, "target")))

		if isSalesPerformanceQuery {
			if performanceContext, performanceInfo := s.buildPerformanceContext(messageLower, userID, userCtx); performanceContext != "" {
				contextData = performanceContext
				contextType = "sales_performance"
			} else if performanceInfo != "" {
				dataAccessInfo = performanceInfo
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
				} else {
					if err != nil {
						fmt.Printf("Error fetching deals: %v\n", err)
					}
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
				fmt.Printf("Fetching leads - Error: %v, Count: %d, Total: %d, Status: %s\n", err, len(leads), total, req.Status)
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
					contextData = instruction
					contextType = "lead"
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
					if dataAccessInfo == "" {
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
					if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data visit reports dari database. Data mungkin tidak tersedia."
					}
				}
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
				// Build request with optional status filter
				req := &task.ListTasksRequest{
					Page:          1,
					PerPage:       20,
					ScopedUserIDs: s.scopedUserIDs(userCtx, "task"),
				}

				// Extract status filter from message if mentioned
				if strings.Contains(messageLower, "pending") {
					req.Status = "pending"
				} else if strings.Contains(messageLower, "in_progress") || strings.Contains(messageLower, "in-progress") {
					req.Status = "in_progress"
				} else if strings.Contains(messageLower, "completed") || strings.Contains(messageLower, "done") {
					req.Status = "completed"
				} else if strings.Contains(messageLower, "cancelled") {
					req.Status = "cancelled"
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
					if dataAccessInfo == "" {
						dataAccessInfo = "⚠️ Tidak dapat mengakses data tasks dari database. Data mungkin tidak tersedia."
					}
				}
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
		// If timezone is invalid, use UTC
		loc = time.UTC
		fmt.Printf("Warning: Invalid timezone '%s', using UTC instead\n", timezone)
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

	// Available models from the UI dropdown (case-insensitive matching)
	// Models available: Llama-3.1-8B, Qwen-3-32B, GPT-OSS-120B, ZAI GLM 4.6, Llama-3.3-70B, Qwen3-235B (Instruct)
	availableModels := map[string]string{
		// Llama models
		"llama-3.1-8b":  "llama-3.1-8b",
		"llama-3.1-70b": "llama-3.1-70b",
		"llama-3.3-70b": "llama-3.3-70b",
		"llama-3-8b":    "llama-3.1-8b",  // Normalize to available model
		"llama-3-70b":   "llama-3.3-70b", // Normalize to available model
		"llama3-8b":     "llama-3.1-8b",
		"llama3.1-8b":   "llama-3.1-8b",
		"llama3.3-70b":  "llama-3.3-70b",

		// Qwen models
		"qwen-3-32b": "qwen-3-32b",
		"qwen3-235b": "qwen3-235b",

		// GPT-OSS model
		"gpt-oss-120b": "gpt-oss-120b",
		"gpt-oss":      "gpt-oss-120b",

		// ZAI GLM model
		"zai-glm-4.6": "zai-glm-4.6",
		"zai-glm":     "zai-glm-4.6",
		"zai glm 4.6": "zai-glm-4.6",
		"zai_glm_4.6": "zai-glm-4.6",
	}

	// Check if model is in the available models map
	if normalizedModel, exists := availableModels[selectedModel]; exists {
		selectedModel = normalizedModel
	} else {
		// Model not found in available models
		// Check if it's a GPT model (not GPT-OSS)
		if strings.HasPrefix(selectedModel, "gpt-") && selectedModel != "gpt-oss-120b" {
			return nil, fmt.Errorf("model '%s' tidak didukung. Model yang tersedia: llama-3.1-8b, llama-3.3-70b, qwen-3-32b, qwen3-235b, gpt-oss-120b, zai-glm-4.6. Silakan pilih model yang valid.", originalModel)
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
			return nil, fmt.Errorf("model '%s' tidak ditemukan atau tidak tersedia. Model yang tersedia: llama-3.1-8b, llama-3.1-70b. Silakan pilih model yang valid.", selectedModel)
		}

		// Check if error is about GPT models
		if strings.Contains(errorStr, "gpt-") {
			return nil, fmt.Errorf("model '%s' tidak didukung. Cerebras API hanya mendukung model Cerebras (contoh: llama-3.1-8b, llama-3.1-70b). Silakan pilih model Cerebras yang valid.", selectedModel)
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

	// 5. Load lead statuses if the user mentions updating lead status.
	if strings.Contains(messageLower, "lead") && (strings.Contains(messageLower, "status") || strings.Contains(messageLower, "ubah") || strings.Contains(messageLower, "update")) {
		leadCtx := s.loadLeadStatusesForCRUD(messageLower, entityIDs, userID, userCtx)
		if leadCtx != "" {
			parts = append(parts, leadCtx)
		}
	}

	// 6. Load pipeline stages if the user mentions moving a deal stage.
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

// extractEntityIDsFromHistory scans assistant messages in conversation history
// for all entity type links (account://UUID, contact://UUID, deal://UUID, etc.)
// and returns them grouped by type.
func (s *Service) extractEntityIDsFromHistory(history []ai.ChatMessage) map[string][]string {
	result := map[string][]string{}
	prefixes := []string{"account://", "contact://", "deal://", "lead://", "task://", "visit://"}

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
		results, _, err := s.accountRepo.List(&account.ListAccountsRequest{
			Page:          1,
			PerPage:       5,
			Search:        term,
			ScopedUserIDs: s.scopedUserIDs(userCtx, "account"),
		})
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
		"untuk ", "for ", "kontak ", "contact ",
		"lead ", "deal ", "account ", "akun ", "brick ",
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
		"may": time.May, "mei": time.May,
		"june": time.June, "juni": time.June,
		"july": time.July, "juli": time.July,
		"august": time.August, "agustus": time.August,
		"september": time.September,
		"october": time.October, "oktober": time.October,
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

	return aiDateRange{}
}

func normalizeDateRangeForRequest(dr aiDateRange) (start string, end string, period string) {
	if dr.Start != "" || dr.End != "" {
		return dr.Start, dr.End, ""
	}
	return "", "", dr.Period
}

func summarizePerformanceByBrick(items []sales_overview.SalesPerformanceListResponse, usersByID map[string]string, brickNames map[string]string) []map[string]interface{} {
	type summary struct {
		BrickID              string
		BrickName            string
		SalesCount           int
		TotalRevenue         int64
		DealsClosed          int
		VisitsCompleted      int
		TasksCompleted       int
		TargetAmount         int64
		TotalProspects       int
		WonProspects         int
		LostProspects        int
		ConversionRateSum    float64
		AchievementRateSum   float64
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
		return "", "⚠️ Tidak ada data performance untuk periode yang diminta."
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

	raw, _ := json.Marshal(map[string]interface{}{
		"period_label": label,
		"total_sales":  total,
		"items":        items,
	})
	return fmt.Sprintf("REAL SALES PERFORMANCE BY SALES REP:\n%s\n\nPresent performance per sales rep. Use Markdown table, sort by revenue, and mention top and bottom performers.", string(raw)), ""
}
