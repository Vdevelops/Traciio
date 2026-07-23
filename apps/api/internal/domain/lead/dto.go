package lead

import (
	"encoding/json"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

type LeadResponse struct {
	ID                 string                 `json:"id"`
	FirstName          string                 `json:"first_name"`
	LastName           string                 `json:"last_name"`
	CompanyName        string                 `json:"company_name"`
	Email              string                 `json:"email"`
	Phone              string                 `json:"phone"`
	JobTitle           string                 `json:"job_title"`
	Industry           string                 `json:"industry"`
	LeadSource         string                 `json:"lead_source"`
	LeadStatus         string                 `json:"lead_status"`
	LeadStatusID       string                 `json:"lead_status_id,omitempty"`
	LeadStatusRef      *LeadStatusRefResponse `json:"lead_status_ref,omitempty"`
	LeadScore          int                    `json:"lead_score"`
	Probability        int                    `json:"probability"`
	EstimatedValue     int64                  `json:"estimated_value"`
	StatusReason       string                 `json:"status_reason,omitempty"`
	BudgetConfirmed    bool                   `json:"budget_confirmed"`
	BudgetAmount       *int64                 `json:"budget_amount,omitempty"`
	AuthorityConfirmed bool                   `json:"authority_confirmed"`
	AuthorityPerson    string                 `json:"authority_person,omitempty"`
	NeedConfirmed      bool                   `json:"need_confirmed"`
	NeedDescription    string                 `json:"need_description,omitempty"`
	TimelineConfirmed  bool                   `json:"timeline_confirmed"`
	ExpectedCloseDate  *time.Time             `json:"expected_close_date,omitempty"`
	AssignedTo         string                 `json:"assigned_to"`
	AssignedUser       *UserRefResponse       `json:"assigned_user,omitempty"`
	AccountID          string                 `json:"account_id"`
	Account            *AccountRefResponse    `json:"account,omitempty"`
	ContactID          string                 `json:"contact_id"`
	Contact            *ContactRefResponse    `json:"contact,omitempty"`
	OpportunityID      string                 `json:"opportunity_id"`
	Opportunity        *DealRefResponse       `json:"opportunity,omitempty"`
	ConvertedAt        *time.Time             `json:"converted_at"`
	ConvertedBy        string                 `json:"converted_by"`
	Notes              string                 `json:"notes"`
	Address            string                 `json:"address"`
	City               string                 `json:"city"`
	Province           string                 `json:"province"`
	PostalCode         string                 `json:"postal_code"`
	Country            string                 `json:"country"`
	Latitude           *float64               `json:"latitude,omitempty"`
	Longitude          *float64               `json:"longitude,omitempty"`
	Website            string                 `json:"website"`
	CreatedBy          string                 `json:"created_by"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type UserRefResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}
type AccountRefResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ContactRefResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
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
		ID:                 l.ID,
		FirstName:          l.FirstName,
		LastName:           l.LastName,
		CompanyName:        l.CompanyName,
		Email:              l.Email,
		Phone:              l.Phone,
		JobTitle:           l.JobTitle,
		Industry:           l.Industry,
		LeadSource:         l.LeadSource,
		LeadStatus:         l.LeadStatus,
		LeadStatusID:       getStringValue(l.LeadStatusID),
		LeadScore:          l.LeadScore,
		Probability:        l.Probability,
		EstimatedValue:     l.EstimatedValue,
		StatusReason:       extractLatestStatusReason(l.ConversionMetadata),
		BudgetConfirmed:    l.BudgetConfirmed,
		BudgetAmount:       l.BudgetAmount,
		AuthorityConfirmed: l.AuthorityConfirmed,
		AuthorityPerson:    l.AuthorityPerson,
		NeedConfirmed:      l.NeedConfirmed,
		NeedDescription:    l.NeedDescription,
		TimelineConfirmed:  l.TimelineConfirmed,
		ExpectedCloseDate:  l.ExpectedCloseDate,
		AssignedTo:         getStringValue(l.AssignedTo),
		AccountID:          getStringValue(l.AccountID),
		ContactID:          getStringValue(l.ContactID),
		OpportunityID:      getStringValue(l.OpportunityID),
		ConvertedAt:        l.ConvertedAt,
		ConvertedBy:        getStringValue(l.ConvertedBy),
		Notes:              l.Notes,
		Address:            l.Address,
		City:               l.City,
		Province:           l.Province,
		PostalCode:         l.PostalCode,
		Country:            l.Country,
		Latitude:           l.Latitude,
		Longitude:          l.Longitude,
		Website:            l.Website,
		CreatedBy:          l.CreatedBy,
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
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

func extractLatestStatusReason(metadata []byte) string {
	if len(metadata) == 0 {
		return ""
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return ""
	}

	reason, _ := parsed["latest_status_reason"].(string)
	return reason
}

type CreateLeadRequest struct {
	FirstName          string   `json:"first_name" binding:"omitempty,max=100"`
	LastName           string   `json:"last_name" binding:"omitempty,max=100"`
	CompanyName        string   `json:"company_name" binding:"omitempty,max=255"`
	Email              string   `json:"email" binding:"required,email,max=255"`
	Phone              string   `json:"phone" binding:"omitempty,max=20"`
	JobTitle           string   `json:"job_title" binding:"omitempty,max=100"`
	Industry           string   `json:"industry" binding:"omitempty,max=100"`
	LeadSource         string   `json:"lead_source" binding:"required,max=100"`
	LeadStatus         string   `json:"lead_status" binding:"omitempty,max=50"`
	LeadStatusID       string   `json:"lead_status_id" binding:"omitempty,uuid"`
	LeadScore          int      `json:"lead_score" binding:"omitempty,min=0,max=100"`
	Probability        int      `json:"probability" binding:"omitempty,min=0,max=100"`
	EstimatedValue     int64    `json:"estimated_value" binding:"omitempty,min=0"`
	BudgetConfirmed    bool     `json:"budget_confirmed" binding:"omitempty"`
	BudgetAmount       *int64   `json:"budget_amount" binding:"omitempty,min=0"`
	AuthorityConfirmed bool     `json:"authority_confirmed" binding:"omitempty"`
	AuthorityPerson    string   `json:"authority_person" binding:"omitempty,max=255"`
	NeedConfirmed      bool     `json:"need_confirmed" binding:"omitempty"`
	NeedDescription    string   `json:"need_description" binding:"omitempty"`
	TimelineConfirmed  bool     `json:"timeline_confirmed" binding:"omitempty"`
	AssignedTo         string   `json:"assigned_to" binding:"omitempty,uuid"`
	Notes              string   `json:"notes" binding:"omitempty"`
	Address            string   `json:"address" binding:"omitempty"`
	City               string   `json:"city" binding:"omitempty,max=100"`
	Province           string   `json:"province" binding:"omitempty,max=100"`
	PostalCode         string   `json:"postal_code" binding:"omitempty,max=20"`
	Country            string   `json:"country" binding:"omitempty,max=100"`
	Latitude           *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude          *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	Website            string   `json:"website" binding:"omitempty,max=255"`
}

type UpdateLeadRequest struct {
	FirstName          string   `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName           string   `json:"last_name" binding:"omitempty,max=100"`
	CompanyName        string   `json:"company_name" binding:"omitempty,max=255"`
	Email              string   `json:"email" binding:"omitempty,email,max=255"`
	Phone              string   `json:"phone" binding:"omitempty,max=20"`
	JobTitle           string   `json:"job_title" binding:"omitempty,max=100"`
	Industry           string   `json:"industry" binding:"omitempty,max=100"`
	LeadSource         string   `json:"lead_source" binding:"omitempty,max=100"`
	LeadStatus         string   `json:"lead_status" binding:"omitempty,max=50"`
	LeadStatusID       string   `json:"lead_status_id" binding:"omitempty,uuid"`
	LeadScore          *int     `json:"lead_score" binding:"omitempty,min=0,max=100"`
	Probability        *int     `json:"probability" binding:"omitempty,min=0,max=100"`
	EstimatedValue     *int64   `json:"estimated_value" binding:"omitempty,min=0"`
	StatusReason       string   `json:"status_reason" binding:"omitempty,max=500"`
	BudgetConfirmed    *bool    `json:"budget_confirmed" binding:"omitempty"`
	BudgetAmount       *int64   `json:"budget_amount" binding:"omitempty,min=0"`
	AuthorityConfirmed *bool    `json:"authority_confirmed" binding:"omitempty"`
	AuthorityPerson    string   `json:"authority_person" binding:"omitempty,max=255"`
	NeedConfirmed      *bool    `json:"need_confirmed" binding:"omitempty"`
	NeedDescription    string   `json:"need_description" binding:"omitempty"`
	TimelineConfirmed  *bool    `json:"timeline_confirmed" binding:"omitempty"`
	AssignedTo         string   `json:"assigned_to" binding:"omitempty,uuid"`
	Notes              string   `json:"notes" binding:"omitempty"`
	Address            string   `json:"address" binding:"omitempty"`
	City               string   `json:"city" binding:"omitempty,max=100"`
	Province           string   `json:"province" binding:"omitempty,max=100"`
	PostalCode         string   `json:"postal_code" binding:"omitempty,max=20"`
	Country            string   `json:"country" binding:"omitempty,max=100"`
	Latitude           *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude          *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	Website            string   `json:"website" binding:"omitempty,max=255"`
}

type ConvertLeadRequest struct {
	OpportunityTitle       string                          `json:"opportunity_title" binding:"required,min=1,max=255"`
	OpportunityDescription string                          `json:"opportunity_description" binding:"omitempty"`
	StageID                string                          `json:"stage_id" binding:"required,uuid"`
	Value                  *int64                          `json:"value" binding:"omitempty,min=0"`
	StatusReason           string                          `json:"status_reason" binding:"omitempty,max=500"`
	ProductItems           []ConvertLeadProductItemRequest `json:"product_items" binding:"omitempty,dive"`
}

type ConvertLeadProductItemRequest struct {
	ProductID      string `json:"product_id" binding:"required,uuid"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	UnitPrice      *int64 `json:"unit_price" binding:"omitempty,min=0"`
	DiscountAmount *int64 `json:"discount_amount" binding:"omitempty,min=0"`
	Notes          string `json:"notes" binding:"omitempty,max=500"`
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
	StartDate     string   `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate       string   `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
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
type UserOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
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
