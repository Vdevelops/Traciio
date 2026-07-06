package ai

// AnalyzeVisitReportRequest represents request to analyze visit report.
type AnalyzeVisitReportRequest struct {
	VisitReportID string `json:"visit_report_id" binding:"required,uuid"`
}

// AnalyzeDealRequest represents request to analyze deal.
type AnalyzeDealRequest struct {
	DealID string `json:"deal_id" binding:"required,uuid"`
}

// AnalyzeContactRequest represents request to analyze contact.
type AnalyzeContactRequest struct {
	ContactID string `json:"contact_id" binding:"required,uuid"`
}

// AnalyzeAccountRequest represents request to analyze account.
type AnalyzeAccountRequest struct {
	AccountID string `json:"account_id" binding:"required,uuid"`
}

// AnalyzePipelineRequest represents request to analyze pipeline.
type AnalyzePipelineRequest struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// ChatRequest represents chat request.
type ChatRequest struct {
	Message             string        `json:"message" binding:"required,min=1"`
	Context             string        `json:"context,omitempty"`
	ContextType         string        `json:"context_type,omitempty"`
	ConversationHistory []ChatMessage `json:"conversation_history,omitempty"`
	Model               string        `json:"model,omitempty"`
	Domain              string        `json:"domain,omitempty"`
}

// ChatResponse represents chat response.
type ChatResponse struct {
	Message string `json:"message"`
	Tokens  int    `json:"tokens,omitempty"`
}

// InsightResponse represents generic insight response.
type InsightResponse struct {
	Type   InsightType `json:"type"`
	Data   interface{} `json:"data"`
	Tokens int         `json:"tokens,omitempty"`
}
