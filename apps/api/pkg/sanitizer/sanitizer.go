// Package sanitizer provides utilities for input sanitization and XSS protection
package sanitizer

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Script tag pattern for removing script tags
	scriptTagPattern = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	
	// Event handler pattern for removing inline event handlers
	eventHandlerPattern = regexp.MustCompile(`(?i)on\w+\s*=`)
	
	// JavaScript protocol pattern
	jsProtocolPattern = regexp.MustCompile(`(?i)javascript:`)
	
	// Data URI pattern for potentially dangerous data URIs
	dangerousDataURIPattern = regexp.MustCompile(`(?i)data:text/html`)
	
	// SQL injection patterns (basic detection)
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|eval|expression)\s+(.*)(from|into|table|database|procedure)`),
		regexp.MustCompile(`(?i)(--|;|\/\*|\*\/|xp_|sp_)`),
		regexp.MustCompile(`(?i)'(\s)*(or|and)(\s)*'`),
	}
)

// SanitizeHTML removes potentially dangerous HTML content
func SanitizeHTML(input string) string {
	if input == "" {
		return ""
	}

	// Remove script tags
	sanitized := scriptTagPattern.ReplaceAllString(input, "")
	
	// Remove inline event handlers
	sanitized = eventHandlerPattern.ReplaceAllString(sanitized, "")
	
	// Remove javascript: protocols
	sanitized = jsProtocolPattern.ReplaceAllString(sanitized, "")
	
	// Remove dangerous data URIs
	sanitized = dangerousDataURIPattern.ReplaceAllString(sanitized, "")
	
	// HTML escape remaining content
	sanitized = html.EscapeString(sanitized)
	
	return sanitized
}

// SanitizeInput performs general input sanitization
// Use this for text inputs that should not contain HTML
func SanitizeInput(input string) string {
	if input == "" {
		return ""
	}
	
	// HTML escape to prevent XSS
	return html.EscapeString(input)
}

// SanitizeSQL performs basic SQL injection prevention
// Note: This is NOT a replacement for parameterized queries
// Always use ORM or prepared statements as primary defense
func SanitizeSQL(input string) string {
	if input == "" {
		return ""
	}
	
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// DetectSQLInjection checks if input contains potential SQL injection patterns
// Returns true if suspicious pattern detected
func DetectSQLInjection(input string) bool {
	if input == "" {
		return false
	}
	
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	
	return false
}

// SanitizeFilename removes potentially dangerous characters from filenames
func SanitizeFilename(filename string) string {
	if filename == "" {
		return ""
	}
	
	// Replace path separators
	sanitized := strings.ReplaceAll(filename, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	
	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	
	// Remove control characters
	var result strings.Builder
	for _, r := range sanitized {
		if !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// SanitizeEmail performs basic email sanitization
func SanitizeEmail(email string) string {
	if email == "" {
		return ""
	}
	
	// Trim whitespace
	email = strings.TrimSpace(email)
	
	// Convert to lowercase
	email = strings.ToLower(email)
	
	// Remove any HTML tags
	email = html.EscapeString(email)
	
	return email
}

// SanitizeURL removes potentially dangerous URL schemes
func SanitizeURL(url string) string {
	if url == "" {
		return ""
	}
	
	// Trim whitespace
	url = strings.TrimSpace(url)
	
	// Remove javascript: protocol
	url = jsProtocolPattern.ReplaceAllString(url, "")
	
	// Remove dangerous data URIs
	url = dangerousDataURIPattern.ReplaceAllString(url, "")
	
	return url
}

// TruncateString safely truncates a string to maxLength
// This prevents buffer overflow attacks
func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	
	// Use rune count to handle multi-byte characters correctly
	runes := []rune(input)
	if len(runes) <= maxLength {
		return input
	}
	
	return string(runes[:maxLength])
}

// RemoveControlCharacters removes control characters from input
func RemoveControlCharacters(input string) string {
	if input == "" {
		return ""
	}
	
	var result strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) || r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// SanitizeForLog removes sensitive information from log data
func SanitizeForLog(input string) string {
	if input == "" {
		return ""
	}
	
	// List of sensitive field patterns
	sensitivePatterns := []struct {
		pattern *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`(?i)"password"\s*:\s*"[^"]*"`), `"password":"***REDACTED***"`},
		{regexp.MustCompile(`(?i)"token"\s*:\s*"[^"]*"`), `"token":"***REDACTED***"`},
		{regexp.MustCompile(`(?i)"secret"\s*:\s*"[^"]*"`), `"secret":"***REDACTED***"`},
		{regexp.MustCompile(`(?i)"api_key"\s*:\s*"[^"]*"`), `"api_key":"***REDACTED***"`},
		{regexp.MustCompile(`(?i)"apikey"\s*:\s*"[^"]*"`), `"apikey":"***REDACTED***"`},
		{regexp.MustCompile(`(?i)"authorization"\s*:\s*"Bearer\s+[^"]*"`), `"authorization":"Bearer ***REDACTED***"`},
		{regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9\-._~+/]+=*`), `Bearer ***REDACTED***`},
	}
	
	sanitized := input
	for _, sp := range sensitivePatterns {
		sanitized = sp.pattern.ReplaceAllString(sanitized, sp.replace)
	}
	
	return sanitized
}

// ValidateNoScriptInjection validates that input doesn't contain script injection
// Returns true if input is safe, false if suspicious
func ValidateNoScriptInjection(input string) bool {
	if input == "" {
		return true
	}
	
	lowered := strings.ToLower(input)
	
	// Check for script tags
	if strings.Contains(lowered, "<script") {
		return false
	}
	
	// Check for javascript: protocol
	if strings.Contains(lowered, "javascript:") {
		return false
	}
	
	// Check for event handlers
	if eventHandlerPattern.MatchString(input) {
		return false
	}
	
	// Check for dangerous data URIs
	if strings.Contains(lowered, "data:text/html") {
		return false
	}
	
	return true
}
