package currency

import (
	"fmt"
	"strings"
)

// FormatCurrency formats integer (sen) to formatted currency string (Rupiah)
// amount is in smallest currency unit (sen), e.g., 1000000 = Rp 10,000
func FormatCurrency(amount int64) string {
	// Convert to Rupiah (divide by 100 if stored in sen)
	rupiah := float64(amount) / 100.0
	// Format with thousand separator
	formatted := FormatNumber(rupiah)
	return "Rp " + formatted
}

// FormatNumber formats number with thousand separator (Indonesian format with dots)
func FormatNumber(n float64) string {
	// Convert to int64 to remove decimal places
	amount := int64(n)

	// Handle zero case
	if amount == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if amount < 0 {
		negative = true
		amount = -amount
	}

	// Convert to string
	str := fmt.Sprintf("%d", amount)
	length := len(str)

	// Add thousand separators (dot for Indonesian format)
	// We'll build the result by inserting dots every 3 digits from right
	var parts []string
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{str[start:i]}, parts...)
	}

	result := strings.Join(parts, ".")
	if negative {
		result = "-" + result
	}

	return result
}
