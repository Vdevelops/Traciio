package ai

import (
	"fmt"
	"strings"
)

type aiRetrievedContext struct {
	Intent         string
	Domain         string
	ContextType    string
	PeriodLabel    string
	ScopeLabel     string
	DataBlocks     []aiContextDataBlock
	AccessInfo     string
	SourceDataText string
}

type aiContextDataBlock struct {
	Title string
	JSON  string
	Total int64
	Shown int
	Notes []string
}

func (ctx aiRetrievedContext) hasData() bool {
	return len(ctx.DataBlocks) > 0 || strings.TrimSpace(ctx.SourceDataText) != ""
}

func composeAIContext(plan aiQueryPlan, retrieved aiRetrievedContext) string {
	if !retrieved.hasData() {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("=== AI QUERY PLAN ===\n")
	sb.WriteString(fmt.Sprintf("- Intent: %s\n", plan.Intent))
	sb.WriteString(fmt.Sprintf("- Domain: %s\n", plan.Domain))
	if plan.PeriodLabel != "" {
		sb.WriteString(fmt.Sprintf("- Period: %s\n", plan.PeriodLabel))
	}
	if retrieved.ScopeLabel != "" {
		sb.WriteString(fmt.Sprintf("- Scope: %s\n", retrieved.ScopeLabel))
	}
	if len(plan.DataTypes) > 0 {
		sb.WriteString(fmt.Sprintf("- Data Types: %s\n", strings.Join(plan.DataTypes, ", ")))
	}

	if strings.TrimSpace(retrieved.SourceDataText) != "" {
		sb.WriteString("\n=== VERIFIED DATA ===\n")
		sb.WriteString(retrieved.SourceDataText)
		sb.WriteString("\n")
	} else {
		for _, block := range retrieved.DataBlocks {
			sb.WriteString("\n=== VERIFIED DATA: ")
			sb.WriteString(block.Title)
			sb.WriteString(" ===\n")
			if block.Total > 0 {
				sb.WriteString(fmt.Sprintf("- Showing: %d of %d\n", block.Shown, block.Total))
			} else if block.Shown > 0 {
				sb.WriteString(fmt.Sprintf("- Showing: %d\n", block.Shown))
			}
			for _, note := range block.Notes {
				if note != "" {
					sb.WriteString("- ")
					sb.WriteString(note)
					sb.WriteString("\n")
				}
			}
			sb.WriteString("Data:\n")
			sb.WriteString(block.JSON)
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n=== ANSWER CONTRACT ===\n")
	sb.WriteString("- Gunakan hanya data di bagian VERIFIED DATA sebagai sumber kebenaran.\n")
	sb.WriteString("- Jangan membuat angka, nama entity, status, revenue, quantity, atau ID yang tidak ada di data.\n")
	sb.WriteString("- Jika data kosong atau akses ditolak, jelaskan secara singkat bahwa data tidak tersedia sesuai akses/filter user.\n")
	sb.WriteString("- Sebutkan periode dan scope yang digunakan jika relevan.\n")
	sb.WriteString("- Jika ada bagian EXTERNAL INTELLIGENCE, bedakan dengan jelas data CRM internal vs sumber eksternal dan sertakan URL sumber saat memakai informasi eksternal.\n")
	sb.WriteString("- Untuk sumber eksternal, jangan tulis hanya `sumber 1` atau `sumber 4`; gunakan Markdown link langsung seperti [Judul sumber](https://domain/path) pada kalimat yang relevan.\n")
	sb.WriteString("- Tampilkan data tabular hanya dengan kolom bisnis yang berguna; jangan tampilkan ID sebagai kolom terpisah.\n")
	sb.WriteString("- Link entity boleh memakai format internal yang sudah ada seperti [Nama](lead://id), [Title](deal://id), [Nama](account://id), [Title](task://id), atau [Title](schedule://id).\n")

	return sb.String()
}

func noDataAIMessage(plan aiQueryPlan, info string) string {
	if strings.TrimSpace(info) != "" {
		return info
	}
	period := "filter yang diminta"
	if plan.PeriodLabel != "" {
		period = plan.PeriodLabel
	}
	return fmt.Sprintf("Tidak ada data untuk intent `%s` pada %s sesuai akses data Anda.", plan.Intent, period)
}
