package ai

import (
	"fmt"
	"strings"
	"time"
)

// domainKeywords maps each domain to keywords used for backend-side intent detection.
// This serves as a fallback when the frontend does not provide a domain hint.
var domainKeywords = map[string][]string{
	"route_optimization": {
		"rute", "route", "optimize", "optimasi", "navigasi", "navigation",
		"jarak", "distance", "shortest path", "jalur", "efisien",
		"location", "lokasi", "maps", "peta", "arah", "direction",
		"kunjungan optimal", "visit route",
	},
	"sales": {
		"lead", "leads", "pipeline", "deal", "deals", "opportunity",
		"visit", "kunjungan", "task", "tugas", "schedule", "jadwal",
		"follow up", "follow-up", "konversi", "conversion", "prospect",
		"prospek", "sales", "penjualan", "closed won", "closed lost",
		"stage", "tahap", "aktivitas", "activity", "laporan kunjungan",
	},
	"inventory": {
		"product", "produk", "inventory", "inventaris", "stok", "stock",
		"obat", "medicine", "pharmaceutical", "farmasi", "harga", "price",
		"katalog", "catalog", "barang", "item",
	},
	"customers": {
		"account", "akun", "customer", "pelanggan", "contact", "kontak",
		"rumah sakit", "hospital", "klinik", "clinic", "apotek", "pharmacy",
		"dokter", "doctor", "rs", "rsud", "faskes", "facility",
	},
	"analytics": {
		"analytics", "analitik", "performance", "performa", "report",
		"laporan", "statistik", "statistic", "dashboard", "chart",
		"grafik", "graph", "trend", "tren", "kpi", "metric",
		"revenue", "pendapatan", "target", "pencapaian", "achievement",
		"forecast", "prediksi", "prediction", "growth", "pertumbuhan",
		"conversion rate", "rata-rata", "average", "comparison",
	},
	"management": {
		"user", "pengguna", "role", "peran", "group", "grup",
		"brick", "wilayah", "territory",
		"permission", "izin", "hak akses",
		"admin", "management", "manajemen", "assign", "struktur",
	},
}

// getDomainPrompt returns the domain-specific prompt fragment for the given domain.
func getDomainPrompt(domain string) string {
	switch domain {
	case "route_optimization":
		return routeOptimizationDomainPrompt
	case "sales":
		return salesDomainPrompt
	case "inventory":
		return inventoryDomainPrompt
	case "customers":
		return customersDomainPrompt
	case "analytics":
		return analyticsDomainPrompt
	case "management":
		return managementDomainPrompt
	default:
		return generalDomainPrompt
	}
}

// generalDomainPrompt is used when no specific domain is detected.
const generalDomainPrompt = `
ACTIVE MODULE: GENERAL
You are a general CRM assistant. Help the user with any CRM-related question.
If the query clearly relates to a specific domain (Sales, Inventory, Customers, Analytics, Management, Route Optimization), respond with domain-appropriate knowledge.
For ambiguous queries, ask the user to clarify which module they need help with.

Available modules you can assist with:
- Sales: Leads, Pipeline/Deals, Visit Reports, Tasks, Schedules, Activities
- Inventory: Products and product catalog
- Customers: Accounts (healthcare facilities) and Contacts
- Analytics: Sales performance, product analytics, reports, forecasting
- Management: Users, Roles, Groups, Bricks/Territories, Targets
- Route Optimization: Visit route planning and travel efficiency`

// detectDomainFromMessage performs keyword-based domain detection on the backend.
// Returns the best-matching domain or "general" if no clear match.
func detectDomainFromMessage(message string) string {
	messageLower := strings.ToLower(message)

	type scored struct {
		domain string
		score  int
	}

	var results []scored

	for domain, keywords := range domainKeywords {
		score := 0
		for _, kw := range keywords {
			if strings.Contains(messageLower, kw) {
				// Longer keywords get higher weight to prioritise specific matches
				weight := 1
				if len(kw) > 6 {
					weight = 2
				}
				score += weight
			}
		}
		if score > 0 {
			results = append(results, scored{domain: domain, score: score})
		}
	}

	if len(results) == 0 {
		return "general"
	}

	// Find highest scoring domain
	best := results[0]
	for _, r := range results[1:] {
		if r.score > best.score {
			best = r
		}
	}

	return best.domain
}

// BuildModularSystemPrompt composes the final system prompt from modular parts.
// Flow: Core Prompt + Domain Prompt + Time Context + Model Info + Context Data
//
// Parameters:
//   - domain: frontend domain hint (may be empty, triggers backend detection)
//   - message: current user message (used for backend domain detection fallback)
//   - contextID: optional entity ID for specific context
//   - contextType: entity type (lead, deal, account, etc.)
//   - contextData: serialised entity data
//   - dataAccessInfo: privacy/permission warning messages
//   - model: selected AI model name
//   - provider: AI provider name
//   - currentTime: current server time in configured timezone
//   - timezone: timezone string
func BuildModularSystemPrompt(
	domain string,
	message string,
	contextID string,
	contextType string,
	contextData string,
	dataAccessInfo string,
	accessContext string,
	model string,
	provider string,
	currentTime time.Time,
	timezone string,
) string {
	// Resolve domain: prefer frontend hint, fallback to backend detection
	resolvedDomain := domain
	if resolvedDomain == "" || resolvedDomain == "auto" {
		resolvedDomain = detectDomainFromMessage(message)
	}

	// 1. Core prompt (always included)
	var sb strings.Builder
	sb.WriteString(coreSystemPrompt)

	// 2. Domain-specific prompt
	sb.WriteString("\n\n")
	sb.WriteString(getDomainPrompt(resolvedDomain))

	// 3. Time context
	sb.WriteString(buildTimeContext(currentTime, timezone))

	// 4. Model/provider info
	sb.WriteString(buildModelInfo(model, provider))

	// 5. Permission and scope context
	if accessContext != "" {
		sb.WriteString(accessContext)
	}

	// 6. Context data and access info
	if contextID != "" && contextType != "" && contextData != "" {
		// Specific entity context
		contextLabel := getContextLabel(contextType)
		sb.WriteString(fmt.Sprintf("\n\n=== CURRENT CONTEXT: %s ===\nID: %s\nData:\n%s", contextLabel, contextID, contextData))
		sb.WriteString(buildDataInstructions())
	} else if contextData != "" {
		// Data without specific context ID (e.g., list queries)
		sb.WriteString(fmt.Sprintf("\n\n=== DATA PROVIDED ===\n%s", contextData))
		sb.WriteString(buildDataInstructions())
	} else {
		// No data available
		sb.WriteString("\n\nIMPORTANT: You do NOT have access to real data from the system. If the user asks for data (leads, accounts, contacts, deals, visit reports), inform them that data is not available for this query. NEVER create example or sample data.")
	}

	// 7. Data access warnings
	if dataAccessInfo != "" {
		sb.WriteString("\n\n")
		sb.WriteString(dataAccessInfo)
	}

	return sb.String()
}

// buildTimeContext generates the time-aware context block.
func buildTimeContext(currentTime time.Time, timezone string) string {
	dateStr := currentTime.Format("2006-01-02")
	timeStr := currentTime.Format("15:04:05")
	weekdayStr := currentTime.Weekday().String()
	monthStr := currentTime.Month().String()
	yearStr := fmt.Sprintf("%d", currentTime.Year())

	ctx := fmt.Sprintf("\n\nCURRENT DATE AND TIME:\n- Date: %s (%s)\n- Time: %s\n- Timezone: %s\n- Year: %s\n- Month: %s\n",
		dateStr, weekdayStr, timeStr, timezone, yearStr, monthStr)

	// Indonesian holidays context
	now := currentTime
	year := now.Year()

	christmas := time.Date(year, 12, 25, 0, 0, 0, 0, currentTime.Location())
	daysUntilChristmas := int(christmas.Sub(now).Hours() / 24)
	newYear := time.Date(year+1, 1, 1, 0, 0, 0, 0, currentTime.Location())
	daysUntilNewYear := int(newYear.Sub(now).Hours() / 24)

	if daysUntilChristmas >= 0 && daysUntilChristmas <= 60 {
		ctx += fmt.Sprintf("- Christmas is in %d days (December 25, %d)\n", daysUntilChristmas, year)
	}
	if daysUntilNewYear >= 0 && daysUntilNewYear <= 60 {
		ctx += fmt.Sprintf("- New Year is in %d days (January 1, %d)\n", daysUntilNewYear, year+1)
	}

	ctx += "\nTIME-AWARE RULES:\n"
	ctx += "- Use current date/time for contextually appropriate responses\n"
	ctx += "- If user asks current date, respond with exact date from above\n"
	ctx += "- Consider seasonal factors in pharmaceutical sales based on current month\n"

	return ctx
}

// buildModelInfo generates the AI model configuration block.
func buildModelInfo(model string, provider string) string {
	return fmt.Sprintf("\n\nCURRENT AI CONFIGURATION:\n- Provider: %s\n- Model: %s\n\nIf asked about your model or provider, respond: 'Saya menggunakan model %s dari provider %s.'",
		provider, model, model, provider)
}

// getContextLabel maps context type to a human-readable label.
func getContextLabel(contextType string) string {
	switch contextType {
	case "lead":
		return "LEAD"
	case "visit_report":
		return "VISIT REPORT"
	case "deal":
		return "DEAL/OPPORTUNITY"
	case "contact":
		return "CONTACT"
	case "account":
		return "ACCOUNT (HEALTHCARE FACILITY)"
	case "sales_performance":
		return "SALES PERFORMANCE ANALYTICS"
	case "brick_management":
		return "BRICK/TERRITORY MANAGEMENT"
	case "product_analysis":
		return "PRODUCT ANALYSIS"
	case "product":
		return "PRODUCT"
	case "user":
		return "USER MANAGEMENT"
	case "groups":
		return "GROUPS/SEGMENTATION"
	case "target":
		return "TARGET/QUOTA MANAGEMENT"
	case "schedule":
		return "SCHEDULE/VISIT PLANNING"
	default:
		return strings.ToUpper(contextType)
	}
}

// buildDataInstructions returns compact data-handling instructions appended after context data.
func buildDataInstructions() string {
	return `

DATA HANDLING RULES:
1. Use ONLY the data provided above - NEVER fabricate data
2. Present in Markdown table format with clickable links: [Name](type://ID)
3. NEVER show IDs as separate columns
4. If data is empty or doesn't match the query, honestly inform the user
5. For conversion rates: show calculation steps with actual numbers
6. After presenting data, ALWAYS provide insights, follow-up questions, and recommendations`
}
