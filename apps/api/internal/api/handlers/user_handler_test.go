package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	userservice "github.com/gilabs/crm-healthcare/api/internal/service/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestUserHandler() (*UserHandler, *MockUserRepository, *MockRoleRepository) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	mockGroupRepo := &MockGroupRepository{}
	mockBrickRepo := &MockBrickRepository{}
	mockTargetRepo := &MockMonthlyTargetRepository{}

	svc := userservice.NewService(mockUserRepo, mockRoleRepo, mockGroupRepo, mockBrickRepo, mockTargetRepo, nil)
	return NewUserHandler(svc, nil, nil), mockUserRepo, mockRoleRepo
}

func TestUserHandler_Create_Success(t *testing.T) {
	handler, mockUserRepo, mockRoleRepo := setupTestUserHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/users", handler.Create)
	
	testRoleID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := user.CreateUserRequest{
		Email:    "new@example.com",
		Password: "password123",
		Name:     "New User",
		RoleID:   testRoleID,
	}
	
	// Mock Expectations
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		if id == testRoleID {
			return &role.Role{ID: testRoleID, Code: "user"}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	mockUserRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockUserRepo.CreateFunc = func(u *user.User) error {
		u.ID = "generated-id"
		return nil
	}
	
	mockUserRepo.FindByIDFunc = func(id string) (*user.User, error) {
		return &user.User{
			ID:     id,
			Email:  reqBody.Email,
			Name:   reqBody.Name,
			RoleID: reqBody.RoleID,
			Role:   &role.Role{ID: testRoleID, Code: "user"},
			Status: "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	
	// Execute
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	// Verify
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if w.Code != http.StatusCreated {
		b, _ := json.Marshal(response)
		t.Errorf("expected status 201, got %d. Resp: %s", w.Code, string(b))
	}
	
	if success, ok := response["success"].(bool); !ok || !success {
		b, _ := json.Marshal(response)
		t.Errorf("expected success true, got %v. Resp: %s", response["success"], string(b))
	}
	
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		b, _ := json.Marshal(response)
		t.Errorf("expected data to be map, got %T. Resp: %s", response["data"], string(b))
		return
	}
	
	if data["email"] != "new@example.com" {
		t.Errorf("expected email %s, got %v", "new@example.com", data["email"])
	}
}

func TestUserHandler_List_Success(t *testing.T) {
	handler, mockUserRepo, _ := setupTestUserHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/users", handler.List)
	
	mockUserRepo.ListFunc = func(req *user.ListUsersRequest) ([]user.User, int64, error) {
		users := []user.User{
			{ID: "1", Email: "a@example.com", Name: "User A", Status: "active", RoleID: "r1", Role: &role.Role{Code: "user"}},
			{ID: "2", Email: "b@example.com", Name: "User B", Status: "active", RoleID: "r1", Role: &role.Role{Code: "user"}},
		}
		return users, 2, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/users?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("expected success true, got %v", response["success"])
	}
	
	data := response["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 users, got %d", len(data))
	}
}
