package ai

import (
	"strings"
	"time"
)

type aiQueryPlan struct {
	Intent                  string
	Domain                  string
	DataTypes               []string
	PeriodLabel             string
	DateRange               aiDateRange
	ContextType             string
	RequireClarification    bool
	ClarificationReason     string
	HandledByLegacyFlow     bool
	PreferAuthenticatedUser bool
	FallbackPlanIfNoData    bool
	FallbackPlanYears       int
	SortBy                  string
}

func (p aiQueryPlan) hasDataTypes() bool {
	return len(p.DataTypes) > 0
}

func (s *Service) planAIQuery(message string, domain string, now time.Time) aiQueryPlan {
	messageLower := strings.ToLower(strings.TrimSpace(message))
	resolvedDomain := domain
	if resolvedDomain == "" || resolvedDomain == "auto" {
		resolvedDomain = detectDomainFromMessage(messageLower)
	}

	plan := aiQueryPlan{
		Intent:      "general_crm_help",
		Domain:      resolvedDomain,
		ContextType: "",
		DateRange:   parseAIDateRange(messageLower, now),
	}
	plan.PeriodLabel = "semua periode"
	if plan.DateRange.HasFilter && plan.DateRange.Label != "" {
		plan.PeriodLabel = plan.DateRange.Label
	}

	if isProductSalesAnalyticsIntent(messageLower) {
		plan.Intent = "my_product_sales"
		plan.Domain = "analytics"
		plan.DataTypes = []string{"product_analysis"}
		plan.ContextType = "product_analysis"
		plan.PreferAuthenticatedUser = isMyScopeQuestion(messageLower)
		plan.FallbackPlanIfNoData = wantsSalesPlanIfNoData(messageLower)
		plan.FallbackPlanYears = salesPlanYearsFromMessage(messageLower)
		plan.SortBy = productSalesSortByFromMessage(messageLower)
		return plan
	}

	if isExternalIntelligencePlannerIntent(messageLower) {
		plan.Intent = "external_market_intelligence"
		plan.Domain = "analytics"
		if isProductRegulatoryExternalIntent(messageLower) {
			plan.DataTypes = []string{"product", "external_intelligence"}
		} else {
			plan.DataTypes = []string{"product_analysis", "external_intelligence"}
		}
		plan.ContextType = "external_market_intelligence"
		return plan
	}

	if isMarketTrendPlannerIntent(messageLower) {
		plan.Intent = "market_trend_proxy"
		plan.Domain = "analytics"
		plan.DataTypes = []string{"product_analysis"}
		plan.ContextType = "market_trend_proxy"
		return plan
	}

	if isSalesPerformancePlannerIntent(messageLower) {
		plan.Intent = "sales_performance_summary"
		plan.Domain = "analytics"
		plan.DataTypes = []string{"sales_performance"}
		plan.ContextType = "sales_performance"
		return plan
	}

	if isTargetPlannerIntent(messageLower) && !isDealValueTargetIntent(messageLower) {
		plan.Intent = "target_achievement"
		plan.Domain = "analytics"
		plan.DataTypes = []string{"target"}
		plan.ContextType = "target"
		return plan
	}

	if isWonLostPlannerIntent(messageLower) {
		plan.Intent = "won_lost_deals"
		plan.Domain = "sales"
		plan.DataTypes = []string{"deal"}
		plan.ContextType = "deal"
		return plan
	}

	if isLeadPlannerIntent(messageLower) {
		plan.Intent = "lead_lookup"
		plan.Domain = "sales"
		plan.DataTypes = []string{"lead"}
		plan.ContextType = "lead"
		return plan
	}

	if isDealPlannerIntent(messageLower) {
		plan.Intent = "pipeline_status"
		plan.Domain = "sales"
		plan.DataTypes = []string{"deal"}
		plan.ContextType = "deal"
		return plan
	}

	if isSchedulePlannerIntent(messageLower) {
		plan.Intent = "schedule_lookup"
		plan.Domain = "sales"
		plan.DataTypes = []string{"schedule"}
		plan.ContextType = "schedule"
		return plan
	}

	if isTaskPlannerIntent(messageLower) {
		plan.Intent = "task_lookup"
		plan.Domain = "sales"
		plan.DataTypes = []string{"task"}
		plan.ContextType = "task"
		return plan
	}

	if isVisitPlannerIntent(messageLower) {
		plan.Intent = "visit_report_lookup"
		plan.Domain = "sales"
		plan.DataTypes = []string{"visit_report"}
		plan.ContextType = "visit_report"
		return plan
	}

	if isAccountPlannerIntent(messageLower) {
		plan.Intent = "account_profile"
		plan.Domain = "customers"
		plan.DataTypes = []string{"account"}
		plan.ContextType = "account"
		return plan
	}

	if isContactPlannerIntent(messageLower) {
		plan.Intent = "contact_lookup"
		plan.Domain = "customers"
		plan.DataTypes = []string{"contact"}
		plan.ContextType = "contact"
		return plan
	}

	if isProductCatalogPlannerIntent(messageLower) {
		plan.Intent = "product_catalog"
		plan.Domain = "inventory"
		plan.DataTypes = []string{"product"}
		plan.ContextType = "product"
		return plan
	}

	if isForecastPlannerIntent(messageLower) || isRoutePlannerIntent(messageLower) || isManagementPlannerIntent(messageLower) || isCRUDPlannerIntent(messageLower) {
		plan.HandledByLegacyFlow = true
	}

	return plan
}

func isExternalIntelligencePlannerIntent(messageLower string) bool {
	hasExternalTerm := strings.Contains(messageLower, "internet") ||
		strings.Contains(messageLower, "luar database") ||
		strings.Contains(messageLower, "luar db") ||
		strings.Contains(messageLower, "eksternal") ||
		strings.Contains(messageLower, "external") ||
		strings.Contains(messageLower, "berita") ||
		strings.Contains(messageLower, "news") ||
		strings.Contains(messageLower, "web")
	hasMarketNeed := strings.Contains(messageLower, "pasar") ||
		strings.Contains(messageLower, "market") ||
		strings.Contains(messageLower, "kompetitor") ||
		strings.Contains(messageLower, "competitor") ||
		strings.Contains(messageLower, "regulasi") ||
		strings.Contains(messageLower, "regulation") ||
		strings.Contains(messageLower, "tren") ||
		strings.Contains(messageLower, "trend") ||
		strings.Contains(messageLower, "analisis") ||
		strings.Contains(messageLower, "analisa") ||
		strings.Contains(messageLower, "analysis") ||
		strings.Contains(messageLower, "insight") ||
		strings.Contains(messageLower, "rekomendasi")
	return hasExternalTerm && hasMarketNeed
}

func isProductRegulatoryExternalIntent(messageLower string) bool {
	hasProductTerm := strings.Contains(messageLower, "produk") ||
		strings.Contains(messageLower, "product") ||
		strings.Contains(messageLower, "obat") ||
		strings.Contains(messageLower, "sku")
	hasRegulatoryTerm := strings.Contains(messageLower, "regulasi") ||
		strings.Contains(messageLower, "regulation") ||
		strings.Contains(messageLower, "safety") ||
		strings.Contains(messageLower, "alert") ||
		strings.Contains(messageLower, "medwatch") ||
		strings.Contains(messageLower, "label") ||
		strings.Contains(messageLower, "bpom") ||
		strings.Contains(messageLower, "fda") ||
		strings.Contains(messageLower, "ema") ||
		strings.Contains(messageLower, "diperbarui") ||
		strings.Contains(messageLower, "update")
	return hasProductTerm && hasRegulatoryTerm
}

func productSalesSortByFromMessage(messageLower string) string {
	switch {
	case strings.Contains(messageLower, "revenue") ||
		strings.Contains(messageLower, "pendapatan") ||
		strings.Contains(messageLower, "omzet") ||
		strings.Contains(messageLower, "kontribusi"):
		return "revenue"
	case strings.Contains(messageLower, "profit") ||
		strings.Contains(messageLower, "keuntungan") ||
		strings.Contains(messageLower, "margin"):
		return "profit"
	case strings.Contains(messageLower, "nama"):
		return "name"
	default:
		return "total_sold"
	}
}

func isMarketTrendPlannerIntent(messageLower string) bool {
	hasMarketTerm := strings.Contains(messageLower, "pasar") ||
		strings.Contains(messageLower, "market")
	hasTrendTerm := strings.Contains(messageLower, "tren") ||
		strings.Contains(messageLower, "trend") ||
		strings.Contains(messageLower, "analisa") ||
		strings.Contains(messageLower, "analisis") ||
		strings.Contains(messageLower, "analysis") ||
		strings.Contains(messageLower, "grafik") ||
		strings.Contains(messageLower, "chart")
	return hasMarketTerm && hasTrendTerm
}

func wantsSalesPlanIfNoData(messageLower string) bool {
	hasNoDataCondition := strings.Contains(messageLower, "jika tidak ada") ||
		strings.Contains(messageLower, "kalau tidak ada") ||
		strings.Contains(messageLower, "bila tidak ada") ||
		strings.Contains(messageLower, "if none") ||
		strings.Contains(messageLower, "if no data")
	hasPlanRequest := strings.Contains(messageLower, "plan") ||
		strings.Contains(messageLower, "rencana") ||
		strings.Contains(messageLower, "strategi")
	return hasNoDataCondition && hasPlanRequest
}

func salesPlanYearsFromMessage(messageLower string) int {
	switch {
	case strings.Contains(messageLower, "5 tahun") || strings.Contains(messageLower, "lima tahun"):
		return 5
	case strings.Contains(messageLower, "3 tahun") || strings.Contains(messageLower, "tiga tahun"):
		return 3
	default:
		return 1
	}
}

func isMyScopeQuestion(messageLower string) bool {
	return strings.Contains(messageLower, "saya") ||
		strings.Contains(messageLower, "milik saya") ||
		strings.Contains(messageLower, "punya saya") ||
		strings.Contains(messageLower, "my ") ||
		strings.Contains(messageLower, "me ")
}

func isSalesPerformancePlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "sales performance") ||
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
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "kontribusi")) ||
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "berkontribusi")) ||
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "contribution")) ||
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "contributor")) ||
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "grafik")) ||
		(strings.Contains(messageLower, "sales") && strings.Contains(messageLower, "chart")) ||
		(strings.Contains(messageLower, "penjualan") && hasTrendOrChartTerm(messageLower)) ||
		(strings.Contains(messageLower, "penjualan") && hasSalesIssueAnalysisTerm(messageLower)) ||
		(strings.Contains(messageLower, "sales") && hasSalesIssueAnalysisTerm(messageLower)) ||
		(strings.Contains(messageLower, "revenue") && strings.Contains(messageLower, "target"))
}

func hasTrendOrChartTerm(messageLower string) bool {
	return strings.Contains(messageLower, "tren") ||
		strings.Contains(messageLower, "trend") ||
		strings.Contains(messageLower, "grafik") ||
		strings.Contains(messageLower, "chart") ||
		strings.Contains(messageLower, "line")
}

func hasSalesIssueAnalysisTerm(messageLower string) bool {
	return strings.Contains(messageLower, "kesalahan") ||
		strings.Contains(messageLower, "salah") ||
		strings.Contains(messageLower, "penyimpangan") ||
		strings.Contains(messageLower, "anomali") ||
		strings.Contains(messageLower, "masalah") ||
		strings.Contains(messageLower, "error") ||
		strings.Contains(messageLower, "audit") ||
		strings.Contains(messageLower, "evaluasi")
}

func isTargetPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "target") ||
		strings.Contains(messageLower, "quota") ||
		strings.Contains(messageLower, "kuota")
}

func isWonLostPlannerIntent(messageLower string) bool {
	return (strings.Contains(messageLower, "deal") || strings.Contains(messageLower, "deals") || strings.Contains(messageLower, "pipeline")) &&
		(strings.Contains(messageLower, "won") || strings.Contains(messageLower, "lost") || strings.Contains(messageLower, "win"))
}

func isLeadPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "lead") ||
		strings.Contains(messageLower, "prospek") ||
		strings.Contains(messageLower, "calon pelanggan")
}

func isDealPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "pipeline") ||
		strings.Contains(messageLower, "sales funnel") ||
		strings.Contains(messageLower, "funnel") ||
		strings.Contains(messageLower, "deal") ||
		strings.Contains(messageLower, "opportunity") ||
		strings.Contains(messageLower, "kesempatan")
}

func isSchedulePlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "schedule") ||
		strings.Contains(messageLower, "schedules") ||
		strings.Contains(messageLower, "jadwal")
}

func isTaskPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "task") ||
		strings.Contains(messageLower, "tugas")
}

func isVisitPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "visit") ||
		strings.Contains(messageLower, "kunjungan") ||
		strings.Contains(messageLower, "laporan kunjungan")
}

func isAccountPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "account") ||
		strings.Contains(messageLower, "akun") ||
		strings.Contains(messageLower, "rumah sakit") ||
		strings.Contains(messageLower, "klinik") ||
		strings.Contains(messageLower, "apotek") ||
		strings.Contains(messageLower, "facility")
}

func isContactPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "contact") ||
		strings.Contains(messageLower, "kontak") ||
		strings.Contains(messageLower, "dokter") ||
		strings.Contains(messageLower, "apoteker")
}

func isProductCatalogPlannerIntent(messageLower string) bool {
	if isProductSalesAnalyticsIntent(messageLower) {
		return false
	}
	return strings.Contains(messageLower, "product") ||
		strings.Contains(messageLower, "produk") ||
		strings.Contains(messageLower, "inventory") ||
		strings.Contains(messageLower, "inventaris") ||
		strings.Contains(messageLower, "obat") ||
		strings.Contains(messageLower, "sku")
}

func isForecastPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "forecast") ||
		strings.Contains(messageLower, "grafik forecast") ||
		strings.Contains(messageLower, "prediksi") ||
		strings.Contains(messageLower, "ramalan")
}

func isRoutePlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "rute") ||
		strings.Contains(messageLower, "route") ||
		strings.Contains(messageLower, "optimasi") ||
		strings.Contains(messageLower, "optimize") ||
		strings.Contains(messageLower, "navigasi") ||
		strings.Contains(messageLower, "navigation") ||
		strings.Contains(messageLower, "kunjungan optimal") ||
		strings.Contains(messageLower, "visit route") ||
		strings.Contains(messageLower, "jalur") ||
		strings.Contains(messageLower, "perjalanan")
}

func isManagementPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "user") ||
		strings.Contains(messageLower, "pengguna") ||
		strings.Contains(messageLower, "role") ||
		strings.Contains(messageLower, "peran") ||
		strings.Contains(messageLower, "group") ||
		strings.Contains(messageLower, "grup") ||
		strings.Contains(messageLower, "brick") ||
		strings.Contains(messageLower, "territory") ||
		strings.Contains(messageLower, "wilayah") ||
		strings.Contains(messageLower, "area mapping")
}

func isCRUDPlannerIntent(messageLower string) bool {
	return strings.Contains(messageLower, "buat") ||
		strings.Contains(messageLower, "buatkan") ||
		strings.Contains(messageLower, "tambah") ||
		strings.Contains(messageLower, "create") ||
		strings.Contains(messageLower, "add") ||
		strings.Contains(messageLower, "catat") ||
		strings.Contains(messageLower, "update") ||
		strings.Contains(messageLower, "ubah") ||
		strings.Contains(messageLower, "jadwalkan") ||
		strings.Contains(messageLower, "follow up") ||
		strings.Contains(messageLower, "follow-up") ||
		strings.Contains(messageLower, "followup")
}
