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
