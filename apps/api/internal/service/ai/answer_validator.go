package ai

import (
	"fmt"
	"regexp"
	"strings"
)

func validateGroundedAIAnswer(message string, contextData string, dataAccessInfo string) string {
	finalMessage := strings.TrimSpace(message)
	trimmedContext := strings.TrimSpace(contextData)
	trimmedAccessInfo := strings.TrimSpace(dataAccessInfo)

	if trimmedContext == "" && trimmedAccessInfo != "" {
		return trimmedAccessInfo
	}
	if trimmedContext == "" {
		return finalMessage
	}

	lower := strings.ToLower(finalMessage)
	if strings.Contains(lower, "contoh data") ||
		strings.Contains(lower, "sample data") ||
		strings.Contains(lower, "misalnya terdapat") ||
		strings.Contains(lower, "asumsi") {
		return "Saya tidak bisa menggunakan data contoh atau asumsi untuk pertanyaan ini. Berdasarkan data scoped yang tersedia, silakan gunakan angka dan entity yang ada pada hasil database saja."
	}

	if strings.Contains(trimmedContext, "EXTERNAL INTELLIGENCE") && hasNumberOnlyExternalCitation(finalMessage) {
		if sources := externalSourceLinksFromContext(trimmedContext); sources != "" && !strings.Contains(finalMessage, "### Sumber Eksternal") {
			finalMessage = strings.TrimSpace(finalMessage) + "\n\n### Sumber Eksternal\n" + sources
		}
	}

	return finalMessage
}

func hasNumberOnlyExternalCitation(message string) bool {
	lower := strings.ToLower(message)
	patterns := []string{"sumber", "source"}
	if !strings.Contains(lower, "http://") && !strings.Contains(lower, "https://") {
		for _, pattern := range patterns {
			if strings.Contains(lower, pattern) {
				return true
			}
		}
	}
	return regexp.MustCompile(`(?i)\b(sumber|source)\s*\d+\b`).MatchString(message)
}

func externalSourceLinksFromContext(contextData string) string {
	titlePattern := regexp.MustCompile(`(?m)^\s*\d+\.\s+Title:\s*(.+)$`)
	urlPattern := regexp.MustCompile(`(?m)^\s*URL:\s*(https?://\S+)\s*$`)
	titles := titlePattern.FindAllStringSubmatch(contextData, -1)
	urls := urlPattern.FindAllStringSubmatch(contextData, -1)
	limit := len(titles)
	if len(urls) < limit {
		limit = len(urls)
	}
	if limit > 5 {
		limit = 5
	}
	if limit == 0 {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < limit; i++ {
		title := strings.TrimSpace(titles[i][1])
		url := strings.TrimSpace(urls[i][1])
		if title == "" || url == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", title, url))
	}
	return strings.TrimSpace(sb.String())
}
