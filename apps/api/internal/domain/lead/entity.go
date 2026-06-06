package lead

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
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
