package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
)

// RefreshTokenRepository defines the interface for refresh token repository
type RefreshTokenRepository interface {
	// Create creates a new refresh token
	Create(token *refresh_token.RefreshToken) error
	
	// FindByTokenID finds a refresh token by token ID (jti)
	FindByTokenID(tokenID string) (*refresh_token.RefreshToken, error)
	
	// FindByUserID finds all refresh tokens for a user
	FindByUserID(userID string) ([]*refresh_token.RefreshToken, error)
	
	// Revoke marks a refresh token as revoked
	Revoke(tokenID string) error
	
	// RevokeByUserID revokes all refresh tokens for a user
	RevokeByUserID(userID string) error
	
	// DeleteExpired deletes expired refresh tokens (cleanup job)
	DeleteExpired() error
}

