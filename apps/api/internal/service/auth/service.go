package auth

import (
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserInactive         = errors.New("user is inactive")
	ErrRefreshTokenInvalid  = errors.New("refresh token is invalid")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")
	ErrRoleNotAllowed       = errors.New("role not allowed for mobile login")
)

type Service struct {
	repo              interfaces.AuthRepository
	refreshTokenRepo  interfaces.RefreshTokenRepository
	jwtManager        *jwt.JWTManager
	permService       *permissionservice.Service
}

func NewService(repo interfaces.AuthRepository, refreshTokenRepo interfaces.RefreshTokenRepository, jwtManager *jwt.JWTManager, permService *permissionservice.Service) *Service {
	return &Service{
		repo:             repo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		permService:      permService,
	}
}

// Login authenticates a user and returns tokens
func (s *Service) Login(req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Validate that user has a role (role must exist)
	if user.Role == nil {
		return nil, ErrUserNotFound // User role is missing, treat as invalid user
	}

	// Get role code
	roleCode := user.Role.Code

	// Get permissions
	var permissions []string
	if s.permService != nil {
		permissions, _ = s.permService.GetPermissionsByRole(roleCode)
		// We ignore error here to allow login even if permission fetch fails (fail open/degraded?) 
		// Actually better to fail closed or just return empty list. Empty list is safer.
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, roleCode)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Extract token ID (jti) from refresh token
	tokenID, err := s.jwtManager.ExtractRefreshTokenID(refreshToken)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	refreshTokenEntity := &refresh_token.RefreshToken{
		UserID:    user.ID,
		TokenID:   tokenID,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenTTL()),
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(refreshTokenEntity); err != nil {
		return nil, err
	}

	// Calculate expires in (seconds)
	expiresIn := int(s.jwtManager.AccessTokenTTL().Seconds())

	// Convert to auth response format
	userResp := user.ToUserResponse()
	authUserResp := &auth.UserResponse{
		ID:          userResp.ID,
		Email:       userResp.Email,
		Name:        userResp.Name,
		AvatarURL:   userResp.AvatarURL,
		Role:        roleCode,
		Permissions: permissions,
		Status:      userResp.Status,
		CreatedAt:   userResp.CreatedAt,
		UpdatedAt:   userResp.UpdatedAt,
	}

	return &auth.LoginResponse{
		User:         authUserResp,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken refreshes an access token using refresh token with token rotation
func (s *Service) RefreshToken(refreshToken string) (*auth.LoginResponse, error) {
	// Validate refresh token and extract user ID and token ID
	userID, tokenID, err := s.jwtManager.ValidateRefreshTokenWithID(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Check if token exists in database
	tokenEntity, err := s.refreshTokenRepo.FindByTokenID(tokenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, err
	}

	// Check if token is revoked
	if tokenEntity.Revoked {
		return nil, ErrRefreshTokenRevoked
	}

	// Check if token is expired
	if tokenEntity.IsExpired() {
		return nil, ErrRefreshTokenExpired
	}

	// Verify user ID matches
	if tokenEntity.UserID != userID {
		return nil, ErrRefreshTokenInvalid
	}

	// Find user
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	// Validate that user has a role (role must exist)
	if user.Role == nil {
		return nil, ErrUserNotFound // User role is missing, treat as invalid user
	}

	// Get role code
	roleCode := user.Role.Code

	// Get permissions
	var permissions []string
	if s.permService != nil {
		permissions, _ = s.permService.GetPermissionsByRole(roleCode)
	}

	// Revoke old refresh token (token rotation)
	if err := s.refreshTokenRepo.Revoke(tokenID); err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, roleCode)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Extract token ID (jti) from new refresh token
	newTokenID, err := s.jwtManager.ExtractRefreshTokenID(newRefreshToken)
	if err != nil {
		return nil, err
	}

	// Store new refresh token in database
	newRefreshTokenEntity := &refresh_token.RefreshToken{
		UserID:    user.ID,
		TokenID:   newTokenID,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenTTL()),
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(newRefreshTokenEntity); err != nil {
		return nil, err
	}

	expiresIn := int(s.jwtManager.AccessTokenTTL().Seconds())

	// Convert to auth response format
	userResp := user.ToUserResponse()
	authUserResp := &auth.UserResponse{
		ID:          userResp.ID,
		Email:       userResp.Email,
		Name:        userResp.Name,
		AvatarURL:   userResp.AvatarURL,
		Role:        roleCode,
		Permissions: permissions,
		Status:      userResp.Status,
		CreatedAt:   userResp.CreatedAt,
		UpdatedAt:   userResp.UpdatedAt,
	}

	return &auth.LoginResponse{
		User:         authUserResp,
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Logout revokes a refresh token
func (s *Service) Logout(refreshToken string) error {
	// Extract token ID from refresh token
	tokenID, err := s.jwtManager.ExtractRefreshTokenID(refreshToken)
	if err != nil {
		// If token is invalid, we can't revoke it, but we don't return error
		// This allows logout to succeed even if token is already invalid
		return nil
	}

	// Revoke the token
	return s.refreshTokenRepo.Revoke(tokenID)
}

// MobileLogin authenticates a user for mobile app and returns tokens (only sales role allowed)
func (s *Service) MobileLogin(req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Get role code
	roleCode := "user"
	mobileAccess := false
	if user.Role != nil {
		roleCode = user.Role.Code
		mobileAccess = user.Role.MobileAccess
	}

	// Check if user's role has mobile access enabled
	if !mobileAccess {
		return nil, ErrRoleNotAllowed
	}

	// Get permissions
	var permissions []string
	if s.permService != nil {
		permissions, _ = s.permService.GetPermissionsByRole(roleCode)
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, roleCode)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Extract token ID (jti) from refresh token
	tokenID, err := s.jwtManager.ExtractRefreshTokenID(refreshToken)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	refreshTokenEntity := &refresh_token.RefreshToken{
		UserID:    user.ID,
		TokenID:   tokenID,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenTTL()),
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(refreshTokenEntity); err != nil {
		return nil, err
	}

	// Calculate expires in (seconds)
	expiresIn := int(s.jwtManager.AccessTokenTTL().Seconds())

	// Convert to auth response format
	userResp := user.ToUserResponse()
	authUserResp := &auth.UserResponse{
		ID:          userResp.ID,
		Email:       userResp.Email,
		Name:        userResp.Name,
		AvatarURL:   userResp.AvatarURL,
		Role:        roleCode,
		Permissions: permissions,
		Status:      userResp.Status,
		CreatedAt:   userResp.CreatedAt,
		UpdatedAt:   userResp.UpdatedAt,
	}

	return &auth.LoginResponse{
		User:         authUserResp,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

