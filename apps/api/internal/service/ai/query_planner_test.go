package ai

import (
	"strings"
	"testing"
	"time"
)

func TestPlanAIQueryPrioritizesProductSalesOverCatalog(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("produk saya yang paling banyak terjual bulan ini", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "my_product_sales" {
		t.Fatalf("expected my_product_sales intent, got %s", plan.Intent)
	}
	if plan.ContextType != "product_analysis" {
		t.Fatalf("expected product_analysis context type, got %s", plan.ContextType)
	}
	if !plan.PreferAuthenticatedUser {
		t.Fatal("expected authenticated user scope preference")
	}
}

func TestPlanAIQueryDetectsTeamProductRevenueContribution(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("Produk apa yang paling besar kontribusi revenue tim bulan ini?", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "my_product_sales" {
		t.Fatalf("expected my_product_sales intent, got %s", plan.Intent)
	}
	if plan.ContextType != "product_analysis" {
		t.Fatalf("expected product_analysis context type, got %s", plan.ContextType)
	}
	if plan.SortBy != "revenue" {
		t.Fatalf("expected revenue sort, got %s", plan.SortBy)
	}
	if plan.PreferAuthenticatedUser {
		t.Fatal("expected team/global scope, not own-user scope")
	}
	if plan.PeriodLabel != "bulan ini" {
		t.Fatalf("expected current month period, got %s", plan.PeriodLabel)
	}
}

func TestPlanAIQueryDetectsSalesContributionChart(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("tampilkan sales yang berkontribusi dan buatkan grafik line nya", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "sales_performance_summary" {
		t.Fatalf("expected sales_performance_summary intent, got %s", plan.Intent)
	}
	if plan.ContextType != "sales_performance" {
		t.Fatalf("expected sales_performance context type, got %s", plan.ContextType)
	}
	if !wantsLineChart("tampilkan sales yang berkontribusi dan buatkan grafik line nya") {
		t.Fatal("expected line chart request to be detected")
	}
}

func TestPlanAIQueryDetectsSalesDataFallbackPlan(t *testing.T) {
	svc := &Service{}
	now := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)

	plan := svc.planAIQuery("apakah ada data penjualan tahun kemarin jika ada datanya berikan data nya jika tidak ada datanya maka berikan saya plan untuk penjualan 5 tahun kedepan", "", now)

	if plan.Intent != "my_product_sales" {
		t.Fatalf("expected my_product_sales intent, got %s", plan.Intent)
	}
	if !plan.FallbackPlanIfNoData {
		t.Fatal("expected fallback plan request to be detected")
	}
	if plan.FallbackPlanYears != 5 {
		t.Fatalf("expected 5 fallback years, got %d", plan.FallbackPlanYears)
	}
	if plan.PeriodLabel != "tahun lalu (2025)" {
		t.Fatalf("expected previous-year period label, got %s", plan.PeriodLabel)
	}
	if plan.DateRange.Start != "2025-01-01" || plan.DateRange.End != "2025-12-31" {
		t.Fatalf("expected 2025 date range, got %s to %s", plan.DateRange.Start, plan.DateRange.End)
	}
}

func TestPlanAIQueryDetectsMarketTrendProxy(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("buatkan saya grafik tren analisa pasar", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "market_trend_proxy" {
		t.Fatalf("expected market_trend_proxy intent, got %s", plan.Intent)
	}
	if plan.ContextType != "market_trend_proxy" {
		t.Fatalf("expected market_trend_proxy context type, got %s", plan.ContextType)
	}
	if len(plan.DataTypes) != 1 || plan.DataTypes[0] != "product_analysis" {
		t.Fatalf("expected product_analysis data type, got %#v", plan.DataTypes)
	}
}

func TestPlanAIQueryDetectsExternalIntelligence(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("apakah AI dapat akses data internet untuk insight tren pasar dan rekomendasi", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "external_market_intelligence" {
		t.Fatalf("expected external_market_intelligence intent, got %s", plan.Intent)
	}
	if plan.ContextType != "external_market_intelligence" {
		t.Fatalf("expected external_market_intelligence context type, got %s", plan.ContextType)
	}
	if len(plan.DataTypes) != 2 || plan.DataTypes[1] != "external_intelligence" {
		t.Fatalf("expected external intelligence data type, got %#v", plan.DataTypes)
	}
}

func TestPlanAIQueryDetectsProductRegulatoryExternalIntelligence(t *testing.T) {
	svc := &Service{}

	plan := svc.planAIQuery("Produk mana yang perlu diperbarui karena ada update regulasi/safety alert eksternal?", "", time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC))

	if plan.Intent != "external_market_intelligence" {
		t.Fatalf("expected external_market_intelligence intent, got %s", plan.Intent)
	}
	if len(plan.DataTypes) != 2 || plan.DataTypes[0] != "product" || plan.DataTypes[1] != "external_intelligence" {
		t.Fatalf("expected product + external intelligence data types, got %#v", plan.DataTypes)
	}
}

func TestNoProductSalesFallbackPlanContext(t *testing.T) {
	plan := aiQueryPlan{
		PeriodLabel:          "tahun lalu (2025)",
		FallbackPlanIfNoData: true,
		FallbackPlanYears:    5,
	}

	context := buildNoProductSalesFallbackPlanContext("Tidak ada data penjualan.", plan)

	for _, expected := range []string{
		"Result: no product sales data found",
		"Build a practical 5-year sales plan",
		"Do not invent historical sales numbers",
		"tahun lalu (2025)",
	} {
		if !strings.Contains(context, expected) {
			t.Fatalf("expected context to contain %q, got:\n%s", expected, context)
		}
	}
}

func TestComposeAIContextIncludesGroundingContract(t *testing.T) {
	plan := aiQueryPlan{
		Intent:      "lead_lookup",
		Domain:      "sales",
		DataTypes:   []string{"lead"},
		PeriodLabel: "bulan ini",
	}
	retrieved := aiRetrievedContext{
		Intent:      "lead_lookup",
		Domain:      "sales",
		ContextType: "lead",
		PeriodLabel: "bulan ini",
		ScopeLabel:  "own",
		DataBlocks: []aiContextDataBlock{
			{
				Title: "Leads",
				JSON:  `[{"id":"lead-1","name":"RS Jakarta"}]`,
				Shown: 1,
				Total: 1,
			},
		},
	}

	context := composeAIContext(plan, retrieved)

	for _, expected := range []string{
		"=== AI QUERY PLAN ===",
		"Intent: lead_lookup",
		"=== VERIFIED DATA: Leads ===",
		"Gunakan hanya data di bagian VERIFIED DATA",
		"Jangan membuat angka",
	} {
		if !strings.Contains(context, expected) {
			t.Fatalf("expected context to contain %q, got:\n%s", expected, context)
		}
	}
}

func TestValidateGroundedAIAnswerRejectsSampleData(t *testing.T) {
	message := "Contoh data: total revenue Rp 100.000.000 berdasarkan asumsi."
	context := `=== VERIFIED DATA === [{"name":"Product A","total_revenue":5000}]`

	validated := validateGroundedAIAnswer(message, context, "")

	if strings.Contains(strings.ToLower(validated), "rp 100.000.000") {
		t.Fatalf("expected hallucinated sample amount to be removed, got %s", validated)
	}
	if !strings.Contains(validated, "tidak bisa menggunakan data contoh") {
		t.Fatalf("expected grounded rejection message, got %s", validated)
	}
}

func TestValidateGroundedAIAnswerAddsExternalSourceLinks(t *testing.T) {
	message := "FDA mengeluarkan recall Cetirizine karena masalah keamanan (sumber 4)."
	context := `EXTERNAL INTELLIGENCE:
- Sources:
  1. Title: FDA recall notice for Cetirizine
     URL: https://www.fda.gov/safety/recalls-market-withdrawals-safety-alerts/example
     Host: fda.gov
     Snippet: Recall notice.`

	validated := validateGroundedAIAnswer(message, context, "")

	if !strings.Contains(validated, "[FDA recall notice for Cetirizine](https://www.fda.gov/safety/recalls-market-withdrawals-safety-alerts/example)") {
		t.Fatalf("expected validated answer to include external source URL, got %s", validated)
	}
}
