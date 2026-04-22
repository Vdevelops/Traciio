package lead

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Lead represents a sales lead
type Lead struct {
	ID             string                  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FirstName      string                  `gorm:"type:varchar(100);not null;index:idx_leads_fts,type:gin,expression:to_tsvector('english'\\, first_name || ' ' || COALESCE(last_name\\, '') || ' ' || COALESCE(company_name\\, '') || ' ' || email)" json:"first_name"`
	LastName       string                  `gorm:"type:varchar(100)" json:"last_name"`
	CompanyName    string                  `gorm:"type:varchar(255)" json:"company_name"`
	Email          string                  `gorm:"type:varchar(255);not null;index" json:"email"`
	Phone          string                  `gorm:"type:varchar(20);index" json:"phone"`
	JobTitle       string                  `gorm:"type:varchar(100)" json:"job_title"`
	Industry       string                  `gorm:"type:varchar(100)" json:"industry"`
	LeadSource     string                  `gorm:"type:varchar(100);not null" json:"lead_source"`              // website, referral, cold_call, event, etc.
	LeadStatus     string                  `gorm:"type:varchar(50);not null;default:'new'" json:"lead_status"` // Legacy: new, contacted, qualified, converted, lost (deprecated, use LeadStatusID)
	LeadStatusID   *string                 `gorm:"type:uuid;index" json:"lead_status_id"`                      // Foreign key to lead_statuses table
	LeadStatusRef  *lead_status.LeadStatus `gorm:"foreignKey:LeadStatusID" json:"lead_status_ref,omitempty"`   // Relasi ke LeadStatus
	LeadScore      int                     `gorm:"type:integer;default:0" json:"lead_score"`                   // 0-100
	Probability    int                     `gorm:"type:integer;default:10" json:"probability"`                 // Win probability 0-100
	EstimatedValue int64                   `gorm:"type:bigint;default:0" json:"estimated_value"`               // Estimated deal value
	// BANT Qualification fields
	BudgetConfirmed     bool           `gorm:"type:boolean;default:false" json:"budget_confirmed"`
	BudgetAmount        *int64         `gorm:"type:bigint" json:"budget_amount,omitempty"`
	AuthorityConfirmed  bool           `gorm:"type:boolean;default:false" json:"authority_confirmed"`
	AuthorityPerson     string         `gorm:"type:varchar(255)" json:"authority_person,omitempty"`
	NeedConfirmed       bool           `gorm:"type:boolean;default:false" json:"need_confirmed"`
	NeedDescription     string         `gorm:"type:text" json:"need_description,omitempty"`
	TimelineConfirmed   bool           `gorm:"type:boolean;default:false" json:"timeline_confirmed"`
	ExpectedCloseDate   *time.Time     `gorm:"type:date" json:"expected_close_date,omitempty"`
	AssignedTo          *string        `gorm:"type:uuid;index" json:"assigned_to"` // Sales rep ID
	AssignedUser        *UserRef       `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
	AccountID           *string        `gorm:"type:uuid;index" json:"account_id"` // Created account after conversion
	Account             *AccountRef    `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ContactID           *string        `gorm:"type:uuid;index" json:"contact_id"` // Created contact after conversion
	Contact             *ContactRef    `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
	OpportunityID       *string        `gorm:"type:uuid;index" json:"opportunity_id"` // Converted opportunity/deal
	Opportunity         *DealRef       `gorm:"foreignKey:OpportunityID" json:"opportunity,omitempty"`
	ConvertedAt         *time.Time     `gorm:"type:timestamp" json:"converted_at"`
	ConvertedBy         *string        `gorm:"type:uuid" json:"converted_by"`
	ConvertedPipelineID *string        `gorm:"type:uuid;index" json:"converted_pipeline_id,omitempty"`                // The deal created during conversion
	ConversionMetadata  datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"conversion_metadata,omitempty"` // Stores conversion history
	Notes               string         `gorm:"type:text" json:"notes"`
	Address             string         `gorm:"type:text" json:"address"`
	City                string         `gorm:"type:varchar(100)" json:"city"`
	Province            string         `gorm:"type:varchar(100)" json:"province"`
	PostalCode          string         `gorm:"type:varchar(20)" json:"postal_code"`
	Country             string         `gorm:"type:varchar(100);default:'Indonesia';index" json:"country"`
	Website             string         `gorm:"type:varchar(255);index" json:"website"`
	CreatedBy           string         `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Lead
func (Lead) TableName() string {
	return "leads"
}

// getStringValue returns string value from pointer, or empty string if nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BeforeCreate hook to generate UUID
func (l *Lead) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	if l.Country == "" {
		l.Country = "Indonesia"
	}
	if l.LeadStatus == "" {
		l.LeadStatus = "new"
	}
	return nil
}

// UserRef represents user reference in lead
type UserRef struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// TableName specifies the table name for UserRef
func (UserRef) TableName() string {
	return "users"
}

// AccountRef represents account reference in lead
type AccountRef struct {
	ID   string `gorm:"type:uuid;primary_key" json:"id"`
	Name string `json:"name"`
}

// TableName specifies the table name for AccountRef
func (AccountRef) TableName() string {
	return "accounts"
}

// ContactRef represents contact reference in lead
type ContactRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// TableName specifies the table name for ContactRef
func (ContactRef) TableName() string {
	return "contacts"
}

// DealRef represents deal reference in lead
type DealRef struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Title string `json:"title"`
	Value int64  `json:"value"`
}

// TableName specifies the table name for DealRef
func (DealRef) TableName() string {
	return "deals"
}

// LeadResponse represents lead response DTO
type LeadResponse struct {
	ID            string                 `json:"id"`
	FirstName     string                 `json:"first_name"`
	LastName      string                 `json:"last_name"`
	CompanyName   string                 `json:"company_name"`
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	JobTitle      string                 `json:"job_title"`
	Industry      string                 `json:"industry"`
	LeadSource    string                 `json:"lead_source"`
	LeadStatus    string                 `json:"lead_status"` // Legacy field
	LeadStatusID  string                 `json:"lead_status_id,omitempty"`
	LeadStatusRef *LeadStatusRefResponse `json:"lead_status_ref,omitempty"`
	LeadScore     int                    `json:"lead_score"`
	AssignedTo    string                 `json:"assigned_to"`
	AssignedUser  *UserRefResponse       `json:"assigned_user,omitempty"`
	AccountID     string                 `json:"account_id"`
	Account       *AccountRefResponse    `json:"account,omitempty"`
	ContactID     string                 `json:"contact_id"`
	Contact       *ContactRefResponse    `json:"contact,omitempty"`
	OpportunityID string                 `json:"opportunity_id"`
	Opportunity   *DealRefResponse       `json:"opportunity,omitempty"`
	ConvertedAt   *time.Time             `json:"converted_at"`
	ConvertedBy   string                 `json:"converted_by"`
	Notes         string                 `json:"notes"`
	Address       string                 `json:"address"`
	City          string                 `json:"city"`
	Province      string                 `json:"province"`
	PostalCode    string                 `json:"postal_code"`
	Country       string                 `json:"country"`
	Website       string                 `json:"website"`
	CreatedBy     string                 `json:"created_by"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// UserRefResponse represents user in lead response
type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// AccountRefResponse represents account in lead response
type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ContactRefResponse represents contact in lead response
type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// DealRefResponse represents deal in lead response
type DealRefResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Value          int64  `json:"value"`
	ValueFormatted string `json:"value_formatted,omitempty"`
}

// LeadStatusRefResponse represents lead status in lead response
type LeadStatusRefResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Score       int    `json:"score"`
	Color       string `json:"color"`
	Order       int    `json:"order"`
	IsActive    bool   `json:"is_active"`
	IsDefault   bool   `json:"is_default"`
	IsConverted bool   `json:"is_converted"`
}

// ToLeadResponse converts Lead to LeadResponse
func (l *Lead) ToLeadResponse() *LeadResponse {
	resp := &LeadResponse{
		ID:            l.ID,
		FirstName:     l.FirstName,
		LastName:      l.LastName,
		CompanyName:   l.CompanyName,
		Email:         l.Email,
		Phone:         l.Phone,
		JobTitle:      l.JobTitle,
		Industry:      l.Industry,
		LeadSource:    l.LeadSource,
		LeadStatus:    l.LeadStatus,
		LeadStatusID:  getStringValue(l.LeadStatusID),
		LeadScore:     l.LeadScore,
		AssignedTo:    getStringValue(l.AssignedTo),
		AccountID:     getStringValue(l.AccountID),
		ContactID:     getStringValue(l.ContactID),
		OpportunityID: getStringValue(l.OpportunityID),
		ConvertedAt:   l.ConvertedAt,
		ConvertedBy:   getStringValue(l.ConvertedBy),
		Notes:         l.Notes,
		Address:       l.Address,
		City:          l.City,
		Province:      l.Province,
		PostalCode:    l.PostalCode,
		Country:       l.Country,
		Website:       l.Website,
		CreatedBy:     l.CreatedBy,
		CreatedAt:     l.CreatedAt,
		UpdatedAt:     l.UpdatedAt,
	}

	if l.AssignedUser != nil {
		resp.AssignedUser = &UserRefResponse{
			ID:        l.AssignedUser.ID,
			Name:      l.AssignedUser.Name,
			Email:     l.AssignedUser.Email,
			AvatarURL: l.AssignedUser.AvatarURL,
		}
	}

	if l.Account != nil {
		resp.Account = &AccountRefResponse{
			ID:   l.Account.ID,
			Name: l.Account.Name,
		}
	}

	if l.Contact != nil {
		resp.Contact = &ContactRefResponse{
			ID:    l.Contact.ID,
			Name:  l.Contact.Name,
			Email: l.Contact.Email,
			Phone: l.Contact.Phone,
		}
	}

	if l.LeadStatusRef != nil {
		resp.LeadStatusRef = &LeadStatusRefResponse{
			ID:          l.LeadStatusRef.ID,
			Name:        l.LeadStatusRef.Name,
			Code:        l.LeadStatusRef.Code,
			Description: l.LeadStatusRef.Description,
			Score:       l.LeadStatusRef.Score,
			Color:       l.LeadStatusRef.Color,
			Order:       l.LeadStatusRef.Order,
			IsActive:    l.LeadStatusRef.IsActive,
			IsDefault:   l.LeadStatusRef.IsDefault,
			IsConverted: l.LeadStatusRef.IsConverted,
		}
	}

	if l.Opportunity != nil {
		resp.Opportunity = &DealRefResponse{
			ID:             l.Opportunity.ID,
			Title:          l.Opportunity.Title,
			Value:          l.Opportunity.Value,
			ValueFormatted: formatCurrency(l.Opportunity.Value),
		}
	}

	return resp
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	return currency.FormatCurrency(amount)
}

// CreateLeadRequest represents create lead request DTO
type CreateLeadRequest struct {
	FirstName    string `json:"first_name" binding:"required,min=1,max=100"`
	LastName     string `json:"last_name" binding:"omitempty,max=100"`
	CompanyName  string `json:"company_name" binding:"omitempty,max=255"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"omitempty,max=20"`
	JobTitle     string `json:"job_title" binding:"omitempty,max=100"`
	Industry     string `json:"industry" binding:"omitempty,max=100"`
	LeadSource   string `json:"lead_source" binding:"required"`
	LeadStatus   string `json:"lead_status" binding:"omitempty"` // Deprecated for create/update, prefer LeadStatusID
	LeadStatusID string `json:"lead_status_id" binding:"omitempty,uuid"`
	LeadScore    int    `json:"lead_score" binding:"omitempty,min=0,max=100"`
	AssignedTo   string `json:"assigned_to" binding:"omitempty,uuid"`
	Notes        string `json:"notes" binding:"omitempty"`
	Address      string `json:"address" binding:"omitempty"`
	City         string `json:"city" binding:"omitempty,max=100"`
	Province     string `json:"province" binding:"omitempty,max=100"`
	PostalCode   string `json:"postal_code" binding:"omitempty,max=20"`
	Country      string `json:"country" binding:"omitempty,max=100"`
	Website      string `json:"website" binding:"omitempty,url"`
}

// UpdateLeadRequest represents update lead request DTO
type UpdateLeadRequest struct {
	FirstName    string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName     string `json:"last_name" binding:"omitempty,max=100"`
	CompanyName  string `json:"company_name" binding:"omitempty,max=255"`
	Email        string `json:"email" binding:"omitempty,email"`
	Phone        string `json:"phone" binding:"omitempty,max=20"`
	JobTitle     string `json:"job_title" binding:"omitempty,max=100"`
	Industry     string `json:"industry" binding:"omitempty,max=100"`
	LeadSource   string `json:"lead_source" binding:"omitempty"`
	LeadStatus   string `json:"lead_status" binding:"omitempty"` // Deprecated for update, prefer LeadStatusID
	LeadStatusID string `json:"lead_status_id" binding:"omitempty,uuid"`
	LeadScore    *int   `json:"lead_score" binding:"omitempty,min=0,max=100"`
	AssignedTo   string `json:"assigned_to" binding:"omitempty,uuid"`
	Notes        string `json:"notes" binding:"omitempty"`
	Address      string `json:"address" binding:"omitempty"`
	City         string `json:"city" binding:"omitempty,max=100"`
	Province     string `json:"province" binding:"omitempty,max=100"`
	PostalCode   string `json:"postal_code" binding:"omitempty,max=20"`
	Country      string `json:"country" binding:"omitempty,max=100"`
	Website      string `json:"website" binding:"omitempty,url"`
}

// ConvertLeadRequest represents convert lead to opportunity request DTO
// Will automatically create Account + Contact + Opportunity
type ConvertLeadRequest struct {
	OpportunityTitle       string     `json:"opportunity_title" binding:"required,min=1,max=255"`
	OpportunityDescription string     `json:"opportunity_description" binding:"omitempty"`
	StageID                string     `json:"stage_id" binding:"required,uuid"`
	Value                  *int64     `json:"value" binding:"omitempty,min=0"`
	Probability            *int       `json:"probability" binding:"omitempty,min=0,max=100"`
	ExpectedCloseDate      *time.Time `json:"expected_close_date" binding:"omitempty"`
	// CreateAccount and CreateContact removed - now automatic
	// AccountID and ContactID can be provided to use existing records (optional)
	AccountID string `json:"account_id" binding:"omitempty,uuid"`
	ContactID string `json:"contact_id" binding:"omitempty,uuid"`
}

// ConvertLeadResponse represents convert lead response DTO
type ConvertLeadResponse struct {
	Lead        *LeadResponse `json:"lead"`
	Opportunity interface{}   `json:"opportunity"`       // DealResponse from pipeline package
	Account     interface{}   `json:"account,omitempty"` // AccountResponse from account package
	Contact     interface{}   `json:"contact,omitempty"` // ContactResponse from contact package
}

// CreateAccountFromLeadRequest represents create account from lead request DTO
type CreateAccountFromLeadRequest struct {
	CategoryID    string `json:"category_id" binding:"omitempty,uuid"` // If not provided, will use first available category
	CreateContact bool   `json:"create_contact" binding:"omitempty"`   // Also create contact from lead
}

// CreateAccountFromLeadResponse represents create account from lead response DTO
type CreateAccountFromLeadResponse struct {
	Lead    *LeadResponse `json:"lead"`
	Account interface{}   `json:"account"`           // AccountResponse from account package
	Contact interface{}   `json:"contact,omitempty"` // ContactResponse from contact package (if created)
}

// ListLeadsRequest represents list leads query parameters
type ListLeadsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Status        string   `form:"status" binding:"omitempty"`
	Source        string   `form:"source" binding:"omitempty"`
	AssignedTo    string   `form:"assigned_to" binding:"omitempty,uuid"`
	Search        string   `form:"search" binding:"omitempty"`
	Sort          string   `form:"sort" binding:"omitempty"`
	Order         string   `form:"order" binding:"omitempty,oneof=asc desc"`
	ScopedUserIDs []string `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}

// LeadAnalyticsRequest represents lead analytics query parameters
type LeadAnalyticsRequest struct {
	StartDate  string `form:"start_date" binding:"omitempty"`
	EndDate    string `form:"end_date" binding:"omitempty"`
	AssignedTo string `form:"assigned_to" binding:"omitempty,uuid"`
}

// LeadAnalyticsResponse represents lead analytics response DTO
type LeadAnalyticsResponse struct {
	TotalLeads              int64         `json:"total_leads"`
	LeadsByStatus           []StatusCount `json:"leads_by_status"`
	LeadsBySource           []SourceCount `json:"leads_by_source"`
	ConversionRate          float64       `json:"conversion_rate"`
	AverageTimeToConversion *float64      `json:"average_time_to_conversion,omitempty"`
	QualifiedLeadsCount     int64         `json:"qualified_leads_count"`
	ConvertedLeadsCount     int64         `json:"converted_leads_count"`
}

// StatusCount represents count by status
type StatusCount struct {
	Status     string  `json:"status"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// SourceCount represents count by source
type SourceCount struct {
	Source         string  `json:"source"`
	Count          int64   `json:"count"`
	Percentage     float64 `json:"percentage"`
	ConversionRate float64 `json:"conversion_rate"`
}

// LeadFormDataResponse represents form data for creating a lead
type LeadFormDataResponse struct {
	LeadSources  []LeadSourceOption `json:"lead_sources"`
	LeadStatuses []LeadStatusOption `json:"lead_statuses"`
	Users        []UserOption       `json:"users"`
	Industries   []string           `json:"industries"`
	Provinces    []string           `json:"provinces"`
	Defaults     LeadFormDefaults   `json:"defaults"`
}

// LeadSourceOption represents a lead source option
type LeadSourceOption struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

// LeadStatusOption represents a lead status option
type LeadStatusOption struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

// UserOption represents a user option for assigned_to
type UserOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LeadFormDefaults represents default values for lead form
type LeadFormDefaults struct {
	Country    string `json:"country"`
	LeadStatus string `json:"lead_status"`
	LeadScore  int    `json:"lead_score"`
}

// LeadMobileFormDataResponse represents optimized form data for mobile
type LeadMobileFormDataResponse struct {
	LeadSources  []LeadSourceOption `json:"lead_sources"`
	LeadStatuses []LeadStatusOption `json:"lead_statuses"`
	Industries   []string           `json:"industries"`
	Provinces    []string           `json:"provinces"`
}
