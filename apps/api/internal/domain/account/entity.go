package account

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account represents an account entity (Hospital, Clinic, Pharmacy)
type Account struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null;index;index:idx_accounts_fts,type:gin,expression:to_tsvector('english'\\, name || ' ' || COALESCE(address\\, '') || ' ' || COALESCE(city\\, ''))" json:"name"`
	CategoryID string    `gorm:"type:uuid;not null;index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Address    string    `gorm:"type:text" json:"address"`
	City       string    `gorm:"type:varchar(100);index" json:"city"`
	Province   string    `gorm:"type:varchar(100)" json:"province"`
	Phone      string    `gorm:"type:varchar(20)" json:"phone"`
	Email      string    `gorm:"type:varchar(255)" json:"email"`
	Latitude   *float64  `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude  *float64  `gorm:"type:decimal(11,8)" json:"longitude"`
	PostalCode string    `gorm:"type:varchar(20);default:''" json:"postal_code"`
	Country    string    `gorm:"type:varchar(100);not null;default:'Indonesia'" json:"country"`
	Website    string    `gorm:"type:varchar(255);default:''" json:"website"`
	Industry   string    `gorm:"type:varchar(100);default:''" json:"industry"`
	Status     string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	AssignedTo *string   `gorm:"type:uuid;index" json:"assigned_to"`
	BrickID    *string   `gorm:"type:uuid;index" json:"brick_id"`
	// ContactCount is populated via a SQL subquery in List() — not stored in DB
	ContactCount int            `gorm:"->;column:contact_count" json:"contact_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Category represents account category (imported from category package)
type Category struct {
	ID          string `gorm:"type:uuid;primary_key" json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	BadgeColor  string `json:"badge_color"`
	Status      string `json:"status"`
}

// TableName specifies the table name for Account
func (Account) TableName() string {
	return "accounts"
}

// BeforeCreate hook to generate UUID
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// ToAccountResponse converts Account to AccountResponse
func (a *Account) ToAccountResponse() *AccountResponse {
	resp := &AccountResponse{
		ID:           a.ID,
		Name:         a.Name,
		CategoryID:   a.CategoryID,
		Address:      a.Address,
		City:         a.City,
		Province:     a.Province,
		Phone:        a.Phone,
		Email:        a.Email,
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
		PostalCode:   a.PostalCode,
		Country:      a.Country,
		Website:      a.Website,
		Industry:     a.Industry,
		Status:       a.Status,
		AssignedTo:   a.AssignedTo,
		BrickID:      a.BrickID,
		ContactCount: a.ContactCount,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.Category != nil {
		resp.Category = &CategoryResponse{
			ID:          a.Category.ID,
			Name:        a.Category.Name,
			Code:        a.Category.Code,
			Description: a.Category.Description,
			BadgeColor:  a.Category.BadgeColor,
			Status:      a.Category.Status,
		}
	}
	return resp
}
