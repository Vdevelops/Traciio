package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper function to create service with mocks
func setupTestService() (*Service, *MockAuthRepository, *MockRefreshTokenRepository, *MockRoleRepository) {
	mockAuthRepo := &MockAuthRepository{}
	mockTokenRepo := &MockRefreshTokenRepository{}
	mockRoleRepo := &MockRoleRepository{} 
	
	// JWT Manager (Real)
	jwtManager := jwt.NewJWTManager("test-secret", time.Minute*15, time.Hour*24)
	
	// We pass nil for permission service to test pure auth logic first.
	// We can instantiate it in specific tests if needed.
	svc := NewService(mockAuthRepo, mockTokenRepo, jwtManager, nil)
	
	return svc, mockAuthRepo, mockTokenRepo, mockRoleRepo
}

func TestService_Login_Success(t *testing.T) {
	svc, mockAuthRepo, mockTokenRepo, _ := setupTestService()
	
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	mockUser := &user.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Status:   "active",
		Role: &role.Role{
			Code: "sales_rep",
		},
	}
	
	// Mock expectations
	mockAuthRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		if email == "test@example.com" {
			return mockUser, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	mockTokenRepo.CreateFunc = func(token *refresh_token.RefreshToken) error {
		if token.UserID != "user-1" {
			return errors.New("wrong user id")
		}
		return nil
	}
	
	// Execute
	req := &auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	resp, err := svc.Login(req)
	
	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Token == "" {
		t.Error("expected access token to be generated")
	}
	if resp.RefreshToken == "" {
		t.Error("expected refresh token to be generated")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email %s, got %s", "test@example.com", resp.User.Email)
	}
}

func TestService_Login_InvalidPassword(t *testing.T) {
	svc, mockAuthRepo, _, _ := setupTestService()
	
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	mockUser := &user.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Status:   "active",
	}
	
	mockAuthRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return mockUser, nil
	}
	
	req := &auth.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	
	resp, err := svc.Login(req)
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
	if resp != nil {
		t.Error("expected nil response on failure")
	}
}

func TestService_Login_UserNotFound(t *testing.T) {
	svc, mockAuthRepo, _, _ := setupTestService()
	
	mockAuthRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	req := &auth.LoginRequest{
		Email:    "unknown@example.com",
		Password: "password123",
	}
	
	resp, err := svc.Login(req)
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials (mapped from RecordNotFound), got %v", err)
	}
	if resp != nil {
		t.Error("expected nil response")
	}
}

func TestService_Login_InactiveUser(t *testing.T) {
	svc, mockAuthRepo, _, _ := setupTestService()
	
	mockUser := &user.User{
		ID:     "user-1",
		Email:  "inactive@example.com",
		Status: "inactive",
	}
	
	mockAuthRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return mockUser, nil
	}
	
	req := &auth.LoginRequest{
		Email:    "inactive@example.com",
		Password: "password123",
	}
	
	_, err := svc.Login(req)
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrUserInactive) {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestService_RefreshToken_Success(t *testing.T) {
	svc, mockAuthRepo, mockTokenRepo, _ := setupTestService()
	jwtManager := jwt.NewJWTManager("test-secret", time.Minute*15, time.Hour*24)
	
	userID := "user-1"
	// Generate a valid refresh token
	tokenStr, _ := jwtManager.GenerateRefreshToken(userID)
	tokenID, _ := jwtManager.ExtractRefreshTokenID(tokenStr)
	
	// Mock Token in DB
	mockTokenEntity := &refresh_token.RefreshToken{
		TokenID:   tokenID,
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	
	// Mock User in DB
	mockUser := &user.User{
		ID:     userID,
		Status: "active",
		Role: &role.Role{Code: "user"},
	}
	
	mockTokenRepo.FindByTokenIDFunc = func(id string) (*refresh_token.RefreshToken, error) {
		if id == tokenID {
			return mockTokenEntity, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	mockTokenRepo.RevokeFunc = func(id string) error {
		if id == tokenID {
			mockTokenEntity.Revoked = true
			return nil
		}
		return errors.New("token not found")
	}
	
	mockTokenRepo.CreateFunc = func(token *refresh_token.RefreshToken) error {
		return nil // Success creating new token
	}
	
	mockAuthRepo.FindByIDFunc = func(id string) (*user.User, error) {
		if id == userID {
			return mockUser, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	// Execute
	resp, err := svc.RefreshToken(tokenStr)
	
	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Token == "" {
		t.Error("expected new access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected new refresh token")
	}
	if resp.RefreshToken == tokenStr {
		t.Error("expected DIFFERENT refresh token (rotation)")
	}
}

