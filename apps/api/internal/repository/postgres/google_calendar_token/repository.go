package google_calendar_token

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Google Calendar token repository
func NewRepository(db *gorm.DB) interfaces.GoogleCalendarTokenRepository {
	return &repository{db: db}
}

func (r *repository) FindByUserID(userID string) (*google_calendar_token.GoogleCalendarToken, error) {
	var token google_calendar_token.GoogleCalendarToken
	err := r.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *repository) Create(token *google_calendar_token.GoogleCalendarToken) error {
	// Use upsert to handle soft-deleted records with the same user_id.
	// If a (soft-)deleted row exists for this user, restore and update it instead of failing.
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"access_token", "refresh_token", "token_type", "expires_at", "scope", "updated_at", "deleted_at",
		}),
	}).Create(token).Error
}

func (r *repository) Update(token *google_calendar_token.GoogleCalendarToken) error {
	return r.db.Save(token).Error
}

func (r *repository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&google_calendar_token.GoogleCalendarToken{}).Error
}

