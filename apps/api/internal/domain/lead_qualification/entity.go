package lead_qualification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LeadQualificationChecklist represents the BANT qualification tracking for a lead
type LeadQualificationChecklist struct {
	ID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LeadID string `gorm:"type:uuid;not null;uniqueIndex:uq_lead_qualification_lead_id" json:"lead_id"`

	// Budget tracking
	BudgetTargetAmount   int64  `gorm:"type:bigint;default:0" json:"budget_target_amount"`
	BudgetTargetCurrency string `gorm:"type:varchar(3);not null;default:'IDR'" json:"budget_target_currency"`
	BudgetConfirmed      bool   `gorm:"type:boolean;not null;default:false" json:"budget_confirmed"`
	BudgetNotes          string `gorm:"type:text;default:''" json:"budget_notes"`

	// Authority tracking
	AuthorityTargetPerson string `gorm:"type:varchar(255);default:''" json:"authority_target_person"`
	AuthorityTargetRole   string `gorm:"type:varchar(100);default:''" json:"authority_target_role"`
	AuthorityConfirmed    bool   `gorm:"type:boolean;not null;default:false" json:"authority_confirmed"`
	AuthorityNotes        string `gorm:"type:text;default:''" json:"authority_notes"`

	// Need tracking
	NeedTargetProducts datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"need_target_products"`          // [{product_id, product_name}]
	NeedPriorityLevel  string         `gorm:"type:varchar(20);not null;default:'medium'" json:"need_priority_level"` // low, medium, high, critical
	NeedConfirmed      bool           `gorm:"type:boolean;not null;default:false" json:"need_confirmed"`
	NeedNotes          string         `gorm:"type:text;default:''" json:"need_notes"`

	// Timeline tracking
	TimelineTargetDate  *time.Time `gorm:"type:date" json:"timeline_target_date"`
	TimelineFlexibility string     `gorm:"type:varchar(20);not null;default:'flexible'" json:"timeline_flexibility"` // fixed, flexible, urgent
	TimelineConfirmed   bool       `gorm:"type:boolean;not null;default:false" json:"timeline_confirmed"`
	TimelineNotes       string     `gorm:"type:text;default:''" json:"timeline_notes"`

	// Overall qualification
	QualificationScore  int    `gorm:"type:integer;not null;default:0" json:"qualification_score"`              // 0-100
	QualificationStatus string `gorm:"type:varchar(20);not null;default:'pending'" json:"qualification_status"` // pending, qualified, unqualified

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for LeadQualificationChecklist
func (LeadQualificationChecklist) TableName() string {
	return "lead_qualification_checklist"
}

// BeforeCreate hook to generate UUID
func (q *LeadQualificationChecklist) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}

// CalculateScore computes the qualification score (0-100) based on confirmed BANT items
func (q *LeadQualificationChecklist) CalculateScore() int {
	score := 0
	if q.BudgetConfirmed {
		score += 25
	}
	if q.AuthorityConfirmed {
		score += 25
	}
	if q.NeedConfirmed {
		score += 25
	}
	if q.TimelineConfirmed {
		score += 25
	}
	return score
}

// CalculateStatus determines qualification status based on score
func (q *LeadQualificationChecklist) CalculateStatus() string {
	score := q.CalculateScore()
	if score >= 75 {
		return "qualified"
	}
	if score >= 25 {
		return "pending"
	}
	return "unqualified"
}

// BeforeSave recalculates score and status before persisting
func (q *LeadQualificationChecklist) BeforeSave(tx *gorm.DB) error {
	q.QualificationScore = q.CalculateScore()
	q.QualificationStatus = q.CalculateStatus()
	return nil
}
