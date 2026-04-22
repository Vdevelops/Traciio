package google_calendar_token

import (
	"context"
	"errors"
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/encryption"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

var (
	ErrTokenNotFound      = errors.New("Google Calendar token not found")
	ErrTokenExpired       = errors.New("Google Calendar token expired and refresh failed")
	ErrEncryptionFailed   = errors.New("failed to encrypt token")
	ErrDecryptionFailed   = errors.New("failed to decrypt token")
	ErrInvalidCredentials = errors.New("invalid Google Calendar credentials")
)

type Service struct {
	tokenRepo interfaces.GoogleCalendarTokenRepository
	config    *config.GoogleCalendarConfig
	encKey    string
}

func NewService(
	tokenRepo interfaces.GoogleCalendarTokenRepository,
	config *config.GoogleCalendarConfig,
	encKey string,
) *Service {
	return &Service{
		tokenRepo: tokenRepo,
		config:    config,
		encKey:    encKey,
	}
}

// GetOAuth2Config returns OAuth2 configuration for Google Calendar
func (s *Service) GetOAuth2Config() *oauth2.Config {
	scopes := []string{"https://www.googleapis.com/auth/calendar.events"}
	if s.config.Scopes != "" {
		// Parse comma-separated scopes
		scopes = []string{s.config.Scopes}
	}

	log.Printf("[Google Calendar Service] GetOAuth2Config - ClientID: %s, RedirectURL: %s, Scopes: %v",
		s.config.ClientID, s.config.RedirectURL, scopes)

	if s.config.ClientID == "" {
		log.Printf("[Google Calendar Service] GetOAuth2Config - WARNING: ClientID is empty")
	}
	if s.config.ClientSecret == "" {
		log.Printf("[Google Calendar Service] GetOAuth2Config - WARNING: ClientSecret is empty")
	}
	if s.config.RedirectURL == "" {
		log.Printf("[Google Calendar Service] GetOAuth2Config - WARNING: RedirectURL is empty")
	}

	return &oauth2.Config{
		ClientID:     s.config.ClientID,
		ClientSecret: s.config.ClientSecret,
		RedirectURL:  s.config.RedirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

// StoreToken stores or updates a user's Google Calendar token (encrypted)
func (s *Service) StoreToken(userID string, token *oauth2.Token) error {
	log.Printf("[Google Calendar Service] StoreToken - userID: %s, AccessToken length: %d, RefreshToken length: %d",
		userID, len(token.AccessToken), len(token.RefreshToken))

	// Encrypt tokens
	log.Printf("[Google Calendar Service] StoreToken - encrypting tokens...")
	encryptedAccessToken, err := encryption.Encrypt(token.AccessToken, s.encKey)
	if err != nil {
		log.Printf("[Google Calendar Service] StoreToken - ERROR: Failed to encrypt access token: %v", err)
		return ErrEncryptionFailed
	}

	encryptedRefreshToken, err := encryption.Encrypt(token.RefreshToken, s.encKey)
	if err != nil {
		log.Printf("[Google Calendar Service] StoreToken - ERROR: Failed to encrypt refresh token: %v", err)
		return ErrEncryptionFailed
	}

	log.Printf("[Google Calendar Service] StoreToken - tokens encrypted successfully")

	// Check if token exists
	log.Printf("[Google Calendar Service] StoreToken - checking for existing token...")
	existingToken, err := s.tokenRepo.FindByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[Google Calendar Service] StoreToken - ERROR: Failed to find existing token: %v", err)
		return err
	}

	if existingToken != nil {
		log.Printf("[Google Calendar Service] StoreToken - updating existing token...")
		// Update existing token
		existingToken.AccessToken = encryptedAccessToken
		existingToken.RefreshToken = encryptedRefreshToken
		existingToken.TokenType = token.TokenType
		existingToken.ExpiresAt = token.Expiry
		if token.Extra("scope") != nil {
			if scope, ok := token.Extra("scope").(string); ok {
				existingToken.Scope = scope
			}
		}
		if err := s.tokenRepo.Update(existingToken); err != nil {
			log.Printf("[Google Calendar Service] StoreToken - ERROR: Failed to update token: %v", err)
			return err
		}
		log.Printf("[Google Calendar Service] StoreToken - token updated successfully")
		return nil
	}

	log.Printf("[Google Calendar Service] StoreToken - creating new token...")
	// Create new token
	newToken := &google_calendar_token.GoogleCalendarToken{
		UserID:       userID,
		AccessToken:  encryptedAccessToken,
		RefreshToken: encryptedRefreshToken,
		TokenType:    token.TokenType,
		ExpiresAt:    token.Expiry,
	}
	if token.Extra("scope") != nil {
		if scope, ok := token.Extra("scope").(string); ok {
			newToken.Scope = scope
		}
	}
	if err := s.tokenRepo.Create(newToken); err != nil {
		log.Printf("[Google Calendar Service] StoreToken - ERROR: Failed to create token: %v", err)
		return err
	}

	log.Printf("[Google Calendar Service] StoreToken - token created successfully")
	return nil
}

// GetToken retrieves and decrypts a user's Google Calendar token
func (s *Service) GetToken(userID string) (*oauth2.Token, error) {
	tokenEntity, err := s.tokenRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// Decrypt tokens
	accessToken, err := encryption.Decrypt(tokenEntity.AccessToken, s.encKey)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	refreshToken, err := encryption.Decrypt(tokenEntity.RefreshToken, s.encKey)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// Build oauth2.Token
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenEntity.TokenType,
		Expiry:       tokenEntity.ExpiresAt,
	}

	// Add scope if available
	if tokenEntity.Scope != "" {
		token = token.WithExtra(map[string]interface{}{
			"scope": tokenEntity.Scope,
		})
	}

	return token, nil
}

// RefreshToken refreshes a user's Google Calendar token
func (s *Service) RefreshToken(ctx context.Context, userID string) (*oauth2.Token, error) {
	tokenEntity, err := s.tokenRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// Decrypt refresh token
	refreshToken, err := encryption.Decrypt(tokenEntity.RefreshToken, s.encKey)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// Create OAuth2 config
	oauth2Config := s.GetOAuth2Config()

	// Create token source for refresh
	tokenSource := oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	// Get new token
	newToken, err := tokenSource.Token()
	if err != nil {
		// If refresh fails, mark token as expired
		return nil, ErrTokenExpired
	}

	// Store updated token
	if err := s.StoreToken(userID, newToken); err != nil {
		return nil, err
	}

	return newToken, nil
}

// GetTokenSource returns an oauth2.TokenSource that automatically refreshes tokens
func (s *Service) GetTokenSource(ctx context.Context, userID string) (oauth2.TokenSource, error) {
	oauth2Config := s.GetOAuth2Config()

	// Get current token
	token, err := s.GetToken(userID)
	if err != nil {
		return nil, err
	}

	// Check if token needs refresh
	tokenEntity, err := s.tokenRepo.FindByUserID(userID)
	if err == nil && tokenEntity.NeedsRefresh() {
		refreshedToken, err := s.RefreshToken(ctx, userID)
		if err != nil {
			return nil, err
		}
		token = refreshedToken
	}

	// Create reusable token source
	return oauth2Config.TokenSource(ctx, token), nil
}

// DeleteToken deletes a user's Google Calendar token
func (s *Service) DeleteToken(userID string) error {
	return s.tokenRepo.DeleteByUserID(userID)
}

// GetCalendarService creates a Google Calendar service with automatic token refresh
func (s *Service) GetCalendarService(ctx context.Context, userID string) (*calendar.Service, error) {
	tokenSource, err := s.GetTokenSource(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create HTTP client with token source
	httpClient := oauth2.NewClient(ctx, tokenSource)

	// Create Calendar service
	service, err := calendar.New(httpClient)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// HandleOAuth2Callback handles OAuth2 callback and stores token
func (s *Service) HandleOAuth2Callback(ctx context.Context, code string) (*oauth2.Token, error) {
	log.Printf("[Google Calendar Service] HandleOAuth2Callback - code length: %d", len(code))

	oauth2Config := s.GetOAuth2Config()

	// Exchange code for token
	log.Printf("[Google Calendar Service] HandleOAuth2Callback - exchanging code for token...")
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("[Google Calendar Service] HandleOAuth2Callback - ERROR: Exchange failed: %v", err)
		return nil, err
	}

	log.Printf("[Google Calendar Service] HandleOAuth2Callback - SUCCESS: Token received, AccessToken length: %d, RefreshToken length: %d",
		len(token.AccessToken), len(token.RefreshToken))

	return token, nil
}

// GetAuthURL returns the OAuth2 authorization URL
func (s *Service) GetAuthURL(state string) string {
	log.Printf("[Google Calendar Service] GetAuthURL - state: %s", state)
	oauth2Config := s.GetOAuth2Config()
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Printf("[Google Calendar Service] GetAuthURL - generated URL: %s", authURL)
	return authURL
}

// GetOAuth2ConfigForPlatform returns OAuth2 configuration for specific platform
func (s *Service) GetOAuth2ConfigForPlatform(platform string) *oauth2.Config {
	config := s.GetOAuth2Config()

	if platform == "mobile" {
		// Use HTTPS redirect URL for mobile (same as web)
		// Backend will handle the callback, exchange code, then forward to mobile via deep link
		config.RedirectURL = s.config.RedirectURL
		log.Printf("[Google Calendar Service] GetOAuth2ConfigForPlatform - Mobile redirect (HTTPS): %s", config.RedirectURL)
	}
	// For web, use the configured redirect URL

	return config
}

// GetAuthURLForPlatform returns the OAuth2 authorization URL for specific platform
func (s *Service) GetAuthURLForPlatform(state, platform string) string {
	log.Printf("[Google Calendar Service] GetAuthURLForPlatform - state: %s, platform: %s", state, platform)
	oauth2Config := s.GetOAuth2ConfigForPlatform(platform)
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Printf("[Google Calendar Service] GetAuthURLForPlatform - generated URL: %s", authURL)
	return authURL
}

// ExchangeCode exchanges authorization code for token (used by mobile)
func (s *Service) ExchangeCode(ctx context.Context, code, platform string) (*oauth2.Token, error) {
	log.Printf("[Google Calendar Service] ExchangeCode - code length: %d, platform: %s", len(code), platform)

	oauth2Config := s.GetOAuth2ConfigForPlatform(platform)

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("[Google Calendar Service] ExchangeCode - ERROR: Exchange failed: %v", err)
		return nil, err
	}

	log.Printf("[Google Calendar Service] ExchangeCode - SUCCESS: Token received, AccessToken length: %d, RefreshToken length: %d",
		len(token.AccessToken), len(token.RefreshToken))

	return token, nil
}
