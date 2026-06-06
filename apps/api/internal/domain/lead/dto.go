package lead

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

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
	LeadStatus    string                 `json:"lead_status"`
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

type UserRefResponse struct{ ID, Name, Email, AvatarURL string }
type AccountRefResponse struct{ ID, Name string }
type ContactRefResponse struct{ ID, Name, Email, Phone string }
type DealRefResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Value          int64  `json:"value"`
	ValueFormatted string `json:"value_formatted,omitempty"`
}

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
		resp.AssignedUser = &UserRefResponse{ID: l.AssignedUser.ID, Name: l.AssignedUser.Name, Email: l.AssignedUser.Email, AvatarURL: l.AssignedUser.AvatarURL}
	}
	if l.Account != nil {
		resp.Account = &AccountRefResponse{ID: l.Account.ID, Name: l.Account.Name}
	}
	if l.Contact != nil {
		resp.Contact = &ContactRefResponse{ID: l.Contact.ID, Name: l.Contact.Name, Email: l.Contact.Email, Phone: l.Contact.Phone}
	}
	if l.LeadStatusRef != nil {
		resp.LeadStatusRef = &LeadStatusRefResponse{ID: l.LeadStatusRef.ID, Name: l.LeadStatusRef.Name, Code: l.LeadStatusRef.Code, Description: l.LeadStatusRef.Description, Score: l.LeadStatusRef.Score, Color: l.LeadStatusRef.Color, Order: l.LeadStatusRef.Order, IsActive: l.LeadStatusRef.IsActive, IsDefault: l.LeadStatusRef.IsDefault, IsConverted: l.LeadStatusRef.IsConverted}
	}
	if l.Opportunity != nil {
		resp.Opportunity = &DealRefResponse{ID: l.Opportunity.ID, Title: l.Opportunity.Title, Value: l.Opportunity.Value, ValueFormatted: formatCurrency(l.Opportunity.Value)}
	}
	return resp
}

func formatCurrency(amount int64) string { return currency.FormatCurrency(amount) }

type CreateLeadRequest struct {
	FirstName, LastName, CompanyName, Email, Phone, JobTitle, Industry, LeadSource, LeadStatus, LeadStatusID, AssignedTo, Notes, Address, City, Province, PostalCode, Country, Website string
	LeadScore                                                                                                                                                                          int `json:"lead_score" binding:"omitempty,min=0,max=100"`
}

type UpdateLeadRequest struct {
	FirstName, LastName, CompanyName, Email, Phone, JobTitle, Industry, LeadSource, LeadStatus, LeadStatusID, AssignedTo, Notes, Address, City, Province, PostalCode, Country, Website string
	LeadScore                                                                                                                                                                          *int `json:"lead_score" binding:"omitempty,min=0,max=100"`
}

type ConvertLeadRequest struct {
	OpportunityTitle       string     `json:"opportunity_title" binding:"required,min=1,max=255"`
	OpportunityDescription string     `json:"opportunity_description" binding:"omitempty"`
	StageID                string     `json:"stage_id" binding:"required,uuid"`
	Value                  *int64     `json:"value" binding:"omitempty,min=0"`
	Probability            *int       `json:"probability" binding:"omitempty,min=0,max=100"`
	ExpectedCloseDate      *time.Time `json:"expected_close_date" binding:"omitempty"`
	AccountID              string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID              string     `json:"contact_id" binding:"omitempty,uuid"`
}

type ConvertLeadResponse struct {
	Lead        *LeadResponse `json:"lead"`
	Opportunity interface{}   `json:"opportunity"`
	Account     interface{}   `json:"account,omitempty"`
	Contact     interface{}   `json:"contact,omitempty"`
}

type CreateAccountFromLeadRequest struct {
	CategoryID    string `json:"category_id" binding:"omitempty,uuid"`
	CreateContact bool   `json:"create_contact" binding:"omitempty"`
}

type CreateAccountFromLeadResponse struct {
	Lead    *LeadResponse `json:"lead"`
	Account interface{}   `json:"account"`
	Contact interface{}   `json:"contact,omitempty"`
}

type ListLeadsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Status        string   `form:"status" binding:"omitempty"`
	Source        string   `form:"source" binding:"omitempty"`
	AssignedTo    string   `form:"assigned_to" binding:"omitempty,uuid"`
	Search        string   `form:"search" binding:"omitempty"`
	Sort          string   `form:"sort" binding:"omitempty"`
	Order         string   `form:"order" binding:"omitempty,oneof=asc desc"`
	ScopedUserIDs []string `form:"-" json:"-"`
}

type LeadAnalyticsRequest struct {
	StartDate, EndDate, AssignedTo string
}

type LeadAnalyticsResponse struct {
	TotalLeads              int64         `json:"total_leads"`
	LeadsByStatus           []StatusCount `json:"leads_by_status"`
	LeadsBySource           []SourceCount `json:"leads_by_source"`
	ConversionRate          float64       `json:"conversion_rate"`
	AverageTimeToConversion *float64      `json:"average_time_to_conversion,omitempty"`
	QualifiedLeadsCount     int64         `json:"qualified_leads_count"`
	ConvertedLeadsCount     int64         `json:"converted_leads_count"`
}

type StatusCount struct {
	Status     string
	Count      int64
	Percentage float64
}
type SourceCount struct {
	Source         string
	Count          int64
	Percentage     float64
	ConversionRate float64
}
type LeadFormDataResponse struct {
	LeadSources  []LeadSourceOption `json:"lead_sources"`
	LeadStatuses []LeadStatusOption `json:"lead_statuses"`
	Users        []UserOption       `json:"users"`
	Industries   []string           `json:"industries"`
	Provinces    []string           `json:"provinces"`
	Defaults     LeadFormDefaults   `json:"defaults"`
}
type LeadSourceOption struct{ ID, Value, Label string }
type LeadStatusOption struct{ ID, Value, Label string }
type UserOption struct{ ID, Name, Email string }
type LeadFormDefaults struct {
	Country, LeadStatus string
	LeadScore           int
}
type LeadMobileFormDataResponse struct {
	LeadSources  []LeadSourceOption `json:"lead_sources"`
	LeadStatuses []LeadStatusOption `json:"lead_statuses"`
	Industries   []string           `json:"industries"`
	Provinces    []string           `json:"provinces"`
}
