package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	contactroleservice "github.com/gilabs/crm-healthcare/api/internal/service/contact_role"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestContactRoleHandler() (*ContactRoleHandler, *MockContactRoleRepository) {
	mockRepo := &MockContactRoleRepository{}
	
	svc := contactroleservice.NewService(mockRepo)
	return NewContactRoleHandler(svc), mockRepo
}

func TestContactRoleHandler_Create_Success(t *testing.T) {
	handler, mockRepo := setupTestContactRoleHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/contact-roles", handler.Create)
	
	reqBody := contact_role.CreateContactRoleRequest{
		Name: "Manager",
		Code: "MGR",
	}
	
	mockRepo.FindByCodeFunc = func(code string) (*contact_role.ContactRole, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockRepo.CreateFunc = func(cr *contact_role.ContactRole) error {
		cr.ID = "cr-1"
		return nil
	}
	
	mockRepo.FindByIDFunc = func(id string) (*contact_role.ContactRole, error) {
		return &contact_role.ContactRole{
			ID:        id,
			Name:      reqBody.Name,
			Code:      reqBody.Code,
			CreatedAt: time.Now(),
		}, nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/contact-roles", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestContactRoleHandler_List_Success(t *testing.T) {
	handler, mockRepo := setupTestContactRoleHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/contact-roles", handler.List)
	
	mockRepo.ListFunc = func() ([]contact_role.ContactRole, error) {
		return []contact_role.ContactRole{
			{ID: "cr1", Name: "Role 1"},
			{ID: "cr2", Name: "Role 2"},
		}, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/contact-roles", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}
