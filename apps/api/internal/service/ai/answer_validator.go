package ai

import "strings"

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

	return finalMessage
}
