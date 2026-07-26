package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	productanalyticsdomain "github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	scheduledomain "github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

func (s *Service) retrieveAIContext(plan aiQueryPlan, message string, userID string, userCtx *domainauth.UserContext, now time.Time) aiRetrievedContext {
	messageLower := strings.ToLower(message)
	retrieved := aiRetrievedContext{
		Intent:      plan.Intent,
		Domain:      plan.Domain,
		ContextType: plan.ContextType,
		PeriodLabel: plan.PeriodLabel,
		ScopeLabel:  s.aiScopeLabel(plan, userCtx),
	}

	switch plan.Intent {
	case "my_product_sales":
		return s.retrieveProductSalesContext(plan, userID, userCtx, now)
	case "market_trend_proxy":
		return s.retrieveMarketTrendProxyContext(plan, userID, userCtx, now)
	case "sales_performance_summary":
		if context, info := s.buildPerformanceContext(messageLower, userID, userCtx); context != "" {
			retrieved.SourceDataText = context
		} else {
			retrieved.AccessInfo = info
		}
	case "target_achievement":
		if context, info := s.buildTargetContext(messageLower, userID, userCtx); context != "" {
			retrieved.SourceDataText = context
		} else {
			retrieved.AccessInfo = info
		}
	case "won_lost_deals":
		if context, info := s.buildWonLostDealsContext(messageLower, userID, userCtx); context != "" {
			retrieved.SourceDataText = context
		} else {
			retrieved.AccessInfo = info
		}
	case "lead_lookup":
		return s.retrieveLeadsContext(plan, messageLower, userID, userCtx)
	case "pipeline_status":
		return s.retrieveDealsContext(plan, messageLower, userID, userCtx)
	case "account_profile":
		return s.retrieveAccountsContext(plan, userID, userCtx)
	case "contact_lookup":
		return s.retrieveContactsContext(plan, userID, userCtx)
	case "visit_report_lookup":
		return s.retrieveVisitReportsContext(plan, messageLower, userID, userCtx)
	case "schedule_lookup":
		return s.retrieveSchedulesContext(plan, messageLower, userID, userCtx)
	case "task_lookup":
		return s.retrieveTasksContext(plan, messageLower, userID, userCtx)
	case "product_catalog":
		return s.retrieveProductsContext(plan, messageLower, userID, userCtx)
	}

	return retrieved
}

func (s *Service) retrieveMarketTrendProxyContext(plan aiQueryPlan, userID string, userCtx *domainauth.UserContext, now time.Time) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("product_analysis", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data product analytics tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	if s.productAnalyticsService == nil {
		retrieved.AccessInfo = "⚠️ Layanan product analytics belum tersedia untuk membuat proxy tren pasar."
		return retrieved
	}

	startDate, endDate := aiDateRangeToTimes(plan.DateRange, now)
	if !plan.DateRange.HasFilter {
		endDate = now
		startDate = now.AddDate(0, -11, 0)
		retrieved.PeriodLabel = "12 bulan terakhir"
	}
	scopedUserIDs := s.scopedUserIDs(userCtx, "product_analysis")
	monthlySales, err := s.productAnalyticsService.GetMonthlySales(startDate, endDate, scopedUserIDs)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data tren penjualan internal untuk proxy analisa pasar."
		return retrieved
	}
	if monthlySales == nil || len(monthlySales.MonthlySales) == 0 {
		retrieved.AccessInfo = "Tidak ada data tren penjualan internal sesuai akses dan periode yang diminta. Saya belum bisa membuat grafik tren pasar berbasis data tanpa minimal data penjualan, lead, atau pipeline."
		return retrieved
	}

	raw, _ := json.Marshal(monthlySales)
	retrieved.SourceDataText = fmt.Sprintf(`INTERNAL MARKET TREND PROXY:
- External market research data is not available in this CRM context.
- Use this internal sales trend as a market-demand proxy, not as external market-size data.
- Period: %s
- Scope: %s
- Build a line CHART marker using month labels and total_revenue values.
- Explain clearly that the chart reflects internal sales performance / demand signal.
- Do not claim it represents the whole pharmaceutical market.

Monthly product sales trend:
%s`, retrieved.PeriodLabel, retrieved.ScopeLabel, string(raw))
	return retrieved
}

func (s *Service) retrieveProductSalesContext(plan aiQueryPlan, userID string, userCtx *domainauth.UserContext, now time.Time) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("product_analysis", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data product analytics tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	if s.productAnalyticsService == nil {
		retrieved.AccessInfo = "⚠️ Layanan product analytics belum tersedia untuk AI Assistant."
		return retrieved
	}

	startDate, endDate := aiDateRangeToTimes(plan.DateRange, now)
	scopeKind := s.productSalesScopeKind(userCtx)
	scopedUserIDs := s.scopedUserIDs(userCtx, "product_analysis")
	var products []*productanalyticsdomain.ProductListItem
	var total int64
	var err error
	if scopeKind == "own" || plan.PreferAuthenticatedUser {
		products, total, err = s.productAnalyticsService.GetUserProductSales(userID, startDate, endDate, "total_sold", "desc", 1, 20)
		scopeKind = "own"
	} else {
		products, total, err = s.productAnalyticsService.GetProductsList(startDate, endDate, "", "total_sold", "desc", 1, 20, scopedUserIDs)
	}
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data product analytics dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(products) == 0 {
		noDataMessage := buildNoProductSalesMessage(plan.PeriodLabel, scopeKind)
		if plan.FallbackPlanIfNoData {
			retrieved.ScopeLabel = scopeKind
			retrieved.SourceDataText = buildNoProductSalesFallbackPlanContext(noDataMessage, plan)
			return retrieved
		}
		retrieved.AccessInfo = noDataMessage
		return retrieved
	}
	productsJSON, _ := json.Marshal(products)
	retrieved.ScopeLabel = scopeKind
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{
		Title: "Product Sales",
		JSON:  string(productsJSON),
		Shown: len(products),
		Total: total,
		Notes: []string{
			"Sorted by total_sold desc.",
			"Show product name, SKU, category, total_sold, total_revenue, total_profit, sales_count, and last_sold_at when present.",
			"Never say sales data is unavailable when this block is present.",
		},
	})
	return retrieved
}

func buildNoProductSalesFallbackPlanContext(noDataMessage string, plan aiQueryPlan) string {
	years := plan.FallbackPlanYears
	if years < 1 {
		years = 1
	}
	return fmt.Sprintf(`PRODUCT SALES DATA CHECK:
- Result: no product sales data found.
- Message to user: %s
- Period checked: %s

FALLBACK REQUEST:
- The user explicitly asked for a sales plan if data is not available.
- Build a practical %d-year sales plan.
- Do not invent historical sales numbers.
- The plan should include yearly focus, target strategy, pipeline actions, product/catalog cleanup, lead/account prioritization, visit cadence, KPI tracking, and review rhythm.
- Make it clear that the plan is strategic guidance because no sales records were found for the requested period.`, noDataMessage, plan.PeriodLabel, years)
}

func (s *Service) retrieveLeadsContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("lead", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data leads tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	req := &lead.ListLeadsRequest{Page: 1, PerPage: 20, ScopedUserIDs: s.scopedUserIDs(userCtx, "lead")}
	start, end, _ := normalizeDateRangeForRequest(plan.DateRange)
	req.StartDate = start
	req.EndDate = end
	req.Status = leadStatusFilterFromPlannerMessage(messageLower)
	leads, total, err := s.leadRepo.List(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data leads dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(leads) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	if len(leads) > 15 {
		leads = leads[:15]
	}
	leadsJSON, _ := json.Marshal(s.formatLeadsForAI(leads))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Leads", JSON: string(leadsJSON), Shown: len(leads), Total: total, Notes: []string{"Use [Name](lead://id) for clickable lead links."}})
	return retrieved
}

func (s *Service) retrieveDealsContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("deal", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data deals/pipeline tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	req := &pipeline.ListDealsRequest{Page: 1, PerPage: 20, ScopedUserIDs: s.scopedUserIDs(userCtx, "deal")}
	start, end, _ := normalizeDateRangeForRequest(plan.DateRange)
	req.DateFrom = start
	req.DateTo = end
	req.Status = dealStatusFilterFromPlannerMessage(messageLower)
	deals, total, err := s.dealRepo.List(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data pipeline/deals dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(deals) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	if len(deals) > 15 {
		deals = deals[:15]
	}
	dealsJSON, _ := json.Marshal(s.formatDealsForAI(deals))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Pipeline Deals", JSON: string(dealsJSON), Shown: len(deals), Total: total, Notes: []string{"Use [Title](deal://id) for clickable deal links."}})
	return retrieved
}

func (s *Service) retrieveAccountsContext(plan aiQueryPlan, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("account", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data accounts tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	accounts, total, err := s.accountRepo.List(&account.ListAccountsRequest{Page: 1, PerPage: 10, ScopedUserIDs: s.scopedUserIDs(userCtx, "account")})
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data accounts dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(accounts) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	accountsJSON, _ := json.Marshal(s.formatAccountsForAI(accounts))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Accounts", JSON: string(accountsJSON), Shown: len(accounts), Total: total, Notes: []string{"Use [Name](account://id) for clickable account links."}})
	return retrieved
}

func (s *Service) retrieveContactsContext(plan aiQueryPlan, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("contact", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data contacts tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	contacts, total, err := s.contactRepo.List(&contact.ListContactsRequest{Page: 1, PerPage: 10, ScopedUserIDs: s.scopedUserIDs(userCtx, "contact")})
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data contacts dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(contacts) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	contactsJSON, _ := json.Marshal(contacts)
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Contacts", JSON: string(contactsJSON), Shown: len(contacts), Total: int64(total), Notes: []string{"Use [Name](contact://id) for clickable contact links."}})
	return retrieved
}

func (s *Service) retrieveVisitReportsContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("visit_report", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data visit reports tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	req := &visit_report.ListVisitReportsRequest{Page: 1, PerPage: 10, ScopedUserIDs: s.scopedUserIDs(userCtx, "visit_report"), Status: visitStatusFilterFromPlannerMessage(messageLower)}
	start, end, _ := normalizeDateRangeForRequest(plan.DateRange)
	req.StartDate = start
	req.EndDate = end
	visitReports, total, err := s.visitReportRepo.List(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data visit reports dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(visitReports) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	visitJSON, _ := json.Marshal(s.formatVisitReportsForAI(visitReports))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Visit Reports", JSON: string(visitJSON), Shown: len(visitReports), Total: int64(total), Notes: []string{"Use [Name](visit://id) for clickable visit report links."}})
	return retrieved
}

func (s *Service) retrieveSchedulesContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("schedule", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data schedules tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	if s.scheduleService == nil {
		retrieved.AccessInfo = "⚠️ Layanan schedule belum tersedia untuk AI Assistant."
		return retrieved
	}
	req := &scheduledomain.ListSchedulesRequest{Page: 1, PerPage: 20, ScopedUserIDs: s.scopedUserIDs(userCtx, "schedule"), Status: scheduleStatusFilterFromPlannerMessage(messageLower)}
	start, end, _ := normalizeDateRangeForRequest(plan.DateRange)
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
	schedules, pagination, err := s.scheduleService.ListSchedules(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data schedules dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(schedules) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	total := len(schedules)
	if pagination != nil {
		total = pagination.Total
	}
	schedulesJSON, _ := json.Marshal(schedules)
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Schedules", JSON: string(schedulesJSON), Shown: len(schedules), Total: int64(total), Notes: []string{"Use [Title](schedule://id) for clickable schedule links."}})
	return retrieved
}

func (s *Service) retrieveTasksContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("task", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data tasks tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	req := &task.ListTasksRequest{Page: 1, PerPage: 20, ScopedUserIDs: s.scopedUserIDs(userCtx, "task")}
	if search := taskSearchTermFromMessage(messageLower); search != "" {
		req.Search = search
	}
	if !isTaskStatusUpdateIntent(messageLower) {
		req.Status = taskStatusFilterFromMessage(messageLower)
	}
	tasks, total, err := s.taskRepo.List(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data tasks dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(tasks) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	tasksJSON, _ := json.Marshal(s.formatTasksForAI(tasks))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Tasks", JSON: string(tasksJSON), Shown: len(tasks), Total: total, Notes: []string{"Use [Title](task://id) for clickable task links."}})
	return retrieved
}

func (s *Service) retrieveProductsContext(plan aiQueryPlan, messageLower string, userID string, userCtx *domainauth.UserContext) aiRetrievedContext {
	retrieved := aiRetrievedContext{Intent: plan.Intent, Domain: plan.Domain, ContextType: plan.ContextType, PeriodLabel: plan.PeriodLabel, ScopeLabel: s.aiScopeLabel(plan, userCtx)}
	allowed, _ := s.checkDataPrivacy("product", userID, userCtx)
	if !allowed {
		retrieved.AccessInfo = "⚠️ Akses ke data products tidak diizinkan berdasarkan pengaturan privasi data atau permission yang Anda miliki."
		return retrieved
	}
	if s.productRepo == nil {
		retrieved.AccessInfo = "⚠️ Layanan product belum tersedia untuk AI Assistant."
		return retrieved
	}
	req := &product.ListProductsRequest{Page: 1, PerPage: 20}
	if strings.Contains(messageLower, "active") || strings.Contains(messageLower, "aktif") {
		req.Status = "active"
	} else if strings.Contains(messageLower, "inactive") || strings.Contains(messageLower, "nonaktif") {
		req.Status = "inactive"
	}
	products, total, err := s.productRepo.List(req)
	if err != nil {
		retrieved.AccessInfo = "⚠️ Tidak dapat mengakses data products dari database. Data mungkin tidak tersedia."
		return retrieved
	}
	if len(products) == 0 {
		retrieved.AccessInfo = noDataAIMessage(plan, "")
		return retrieved
	}
	productsJSON, _ := json.Marshal(s.formatProductsForAI(products))
	retrieved.DataBlocks = append(retrieved.DataBlocks, aiContextDataBlock{Title: "Products", JSON: string(productsJSON), Shown: len(products), Total: total, Notes: []string{"Do not invent stock, margin, or sales numbers unless present in context."}})
	return retrieved
}

func (s *Service) aiScopeLabel(plan aiQueryPlan, userCtx *domainauth.UserContext) string {
	if userCtx == nil || len(plan.DataTypes) == 0 {
		return "unknown"
	}
	rule, ok := aiDataAccessRules[plan.DataTypes[0]]
	if !ok {
		return "own"
	}
	return fmt.Sprintf("%s", userCtx.GetScope(rule.Resource))
}

func leadStatusFilterFromPlannerMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "new"):
		return "new"
	case strings.Contains(messageLower, "contacted"):
		return "contacted"
	case strings.Contains(messageLower, "proposal sent") || strings.Contains(messageLower, "proposal_sent"):
		return "proposal_sent"
	case strings.Contains(messageLower, "interested"):
		return "interested"
	case strings.Contains(messageLower, "qualified"):
		return "qualified"
	case strings.Contains(messageLower, "converted"):
		return "converted"
	case strings.Contains(messageLower, "lost"):
		return "lost"
	default:
		return ""
	}
}

func dealStatusFilterFromPlannerMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "closed won") || strings.Contains(messageLower, "won"):
		return "won"
	case strings.Contains(messageLower, "closed lost") || strings.Contains(messageLower, "lost"):
		return "lost"
	default:
		return ""
	}
}

func visitStatusFilterFromPlannerMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "approved"):
		return "approved"
	case strings.Contains(messageLower, "submitted"):
		return "submitted"
	case strings.Contains(messageLower, "draft"):
		return "draft"
	case strings.Contains(messageLower, "rejected"):
		return "rejected"
	case strings.Contains(messageLower, "completed") || strings.Contains(messageLower, "selesai"):
		return "completed"
	default:
		return ""
	}
}

func scheduleStatusFilterFromPlannerMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "pending"):
		return "pending"
	case strings.Contains(messageLower, "submitted"):
		return "submitted"
	case strings.Contains(messageLower, "confirmed"):
		return "confirmed"
	case strings.Contains(messageLower, "completed") || strings.Contains(messageLower, "selesai"):
		return "completed"
	case strings.Contains(messageLower, "cancelled") || strings.Contains(messageLower, "batal"):
		return "cancelled"
	case strings.Contains(messageLower, "rejected"):
		return "rejected"
	default:
		return ""
	}
}
