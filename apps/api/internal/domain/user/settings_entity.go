package user

// GetSettingsSummaryRequest represents request parameters for settings summary
type GetSettingsSummaryRequest struct {
	StartDate string `form:"start_date"` // Optional: Filter start date (YYYY-MM-DD)
	EndDate   string `form:"end_date"`   // Optional: Filter end date (YYYY-MM-DD)
}

// SettingsStats represents extended statistics for settings page
type SettingsStats struct {
	Visits                     int     `json:"visits"`        // Visit reports count
	Deals                      int     `json:"deals"`         // Total deals count
	Tasks                      int     `json:"tasks"`         // Tasks count
	TotalRevenue               int64   `json:"total_revenue"` // Total revenue from all deals
	DealsWon                   int     `json:"deals_won"`     // Count of won deals
	DealsLost                  int     `json:"deals_lost"`    // Count of lost deals
	DealsOpen                  int     `json:"deals_open"`    // Count of open/active deals
	TotalRevenueFormatted      string  `json:"total_revenue_formatted"`      // Formatted revenue string
	ConversionRate             float64 `json:"conversion_rate"`              // Conversion rate percentage
	AverageDealValue           int64   `json:"average_deal_value"`           // Average deal value
	AverageDealValueFormatted  string  `json:"average_deal_value_formatted"` // Formatted average deal value
}

// SettingsSummaryResponse represents the complete settings summary response
type SettingsSummaryResponse struct {
	User         *UserResponse        `json:"user"`
	Stats        *SettingsStats       `json:"stats"`
	Activities   []ProfileActivity    `json:"activities"`
	Transactions []ProfileTransaction `json:"transactions"`
}
