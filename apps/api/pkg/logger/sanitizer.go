package logger

import (
	"fmt"
	"strings"
)

// SanitizeForLog removes sensitive information from log messages
func SanitizeForLog(message string) string {
	// Don't process empty messages
	if message == "" {
		return message
	}

	sanitized := message

	// List of sensitive keywords to redact
	sensitiveKeywords := []string{
		"password",
		"token",
		"secret",
		"api_key",
		"apikey",
		"authorization",
		"bearer",
		"refresh_token",
		"access_token",
		"csrf_token",
		"session",
		"cookie",
		"private_key",
		"encryption_key",
	}

	// Check if message contains sensitive data
	lowerMessage := strings.ToLower(message)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerMessage, keyword) {
			// Replace the value after the keyword
			// Pattern: "keyword": "value" or keyword=value or keyword: value
			patterns := []struct {
				start string
				end   string
			}{
				{fmt.Sprintf(`"%s"`, keyword), `"`},
				{fmt.Sprintf(`"%s":`, keyword), `,`},
				{fmt.Sprintf(`%s=`, keyword), ` `},
				{fmt.Sprintf(`%s:`, keyword), ` `},
			}

			for _, pattern := range patterns {
				if strings.Contains(lowerMessage, strings.ToLower(pattern.start)) {
					sanitized = redactSensitiveValue(sanitized, keyword)
				}
			}
		}
	}

	return sanitized
}

// redactSensitiveValue redacts sensitive values from log messages
func redactSensitiveValue(message string, keyword string) string {
	sanitized := message
	
	// Only process if keyword is present
	if !strings.Contains(strings.ToLower(message), strings.ToLower(keyword)) {
		return sanitized
	}

	// For JSON-like patterns
	if strings.Contains(message, fmt.Sprintf(`"%s"`, keyword)) || strings.Contains(message, fmt.Sprintf(`%s:`, keyword)) {
		// Find the position and redact
		idx := strings.Index(strings.ToLower(sanitized), strings.ToLower(keyword))
		if idx == -1 {
			return sanitized
		}

		// Look for the value part
		start := idx + len(keyword)
		for start < len(sanitized) && (sanitized[start] == ':' || sanitized[start] == '=' || sanitized[start] == ' ' || sanitized[start] == '"') {
			start++
		}
		
		end := start
		inQuotes := false
		if start > 0 && sanitized[start-1] == '"' {
			inQuotes = true
		}
		
		for end < len(sanitized) {
			if inQuotes && sanitized[end] == '"' {
				break
			}
			if !inQuotes && (sanitized[end] == ' ' || sanitized[end] == ',' || sanitized[end] == '\n' || sanitized[end] == '}') {
				break
			}
			end++
		}
		
		if end > start {
			sanitized = sanitized[:start] + "***REDACTED***" + sanitized[end:]
		}
	}


	return sanitized
}
