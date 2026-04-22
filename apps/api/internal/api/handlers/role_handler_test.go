package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	roleservice "github.com/gilabs/crm-healthcare/api/internal/service/role"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestRoleHandler() (*RoleHandler, *MockRoleRepository) {
	mockRoleRepo := &MockRoleRepository{}
	mockUserRepo := &MockUserRepository{}
	
	svc := roleservice.NewService(mockRoleRepo, mockUserRepo)
	return NewRoleHandler(svc), mockRoleRepo
}

func TestRoleHandler_Create_Success(t *testing.T) {
	handler, mockRoleRepo := setupTestRoleHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/roles", handler.Create)
	
	reqBody := role.CreateRoleRequest{
		Name: "New Role",
		Code: "new_role",
	}
	
	mockRoleRepo.FindByCodeFunc = func(code string) (*role.Role, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockRoleRepo.CreateFunc = func(r *role.Role) error {
		r.ID = "role-id-1"
		return nil
	}
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id, Code: reqBody.Code, Name: reqBody.Name}, nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("expected success true, got %v", response["success"])
	}
	
	data := response["data"].(map[string]interface{})
	if data["code"] != "new_role" {
		t.Errorf("expected code new_role, got %v", data["code"])
	}
}

func TestRoleHandler_List_Success(t *testing.T) {
	handler, mockRoleRepo := setupTestRoleHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/roles", handler.List)
	
	mockRoleRepo.ListFunc = func() ([]role.Role, error) {
		return []role.Role{
			{ID: "1", Code: "admin", Name: "Admin"},
			{ID: "2", Code: "user", Name: "User"},
		}, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/roles", nil)
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
		t.Errorf("expected 2 roles, got %d", len(data))
	}
}
