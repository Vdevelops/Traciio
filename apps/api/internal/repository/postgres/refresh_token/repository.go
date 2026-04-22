package refresh_token

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new refresh token repository
func NewRepository(db *gorm.DB) interfaces.RefreshTokenRepository {
	return &repository{db: db}
}

func (r *repository) Create(token *refresh_token.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *repository) FindByTokenID(tokenID string) (*refresh_token.RefreshToken, error) {
	var token refresh_token.RefreshToken
	err := r.db.Where("token_id = ?", tokenID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *repository) FindByUserID(userID string) ([]*refresh_token.RefreshToken, error) {
	var tokens []*refresh_token.RefreshToken
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *repository) Revoke(tokenID string) error {
	now := time.Now()
	return r.db.Model(&refresh_token.RefreshToken{}).
		Where("token_id = ?", tokenID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		}).Error
}

func (r *repository) RevokeByUserID(userID string) error {
	now := time.Now()
	return r.db.Model(&refresh_token.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		}).Error
}

func (r *repository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&refresh_token.RefreshToken{}).Error
}

