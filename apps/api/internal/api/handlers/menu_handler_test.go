package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gin-gonic/gin"
)

func setupTestMenuHandler() (*MenuHandler, *MockMenuRepository) {
	mockRepo := &MockMenuRepository{}
	return NewMenuHandler(mockRepo), mockRepo
}

func TestMenuHandler_List_Success(t *testing.T) {
	handler, mockRepo := setupTestMenuHandler()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/menus", handler.List)

	mockRepo.ListFunc = func() ([]permission.Menu, error) {
		return []permission.Menu{
			{ID: "m1", Name: "Dashboard", URL: "/dashboard", Icon: "home"},
			{ID: "m2", Name: "Settings", URL: "/settings", Icon: "settings"},
		}, nil
	}

	req, _ := http.NewRequest("GET", "/api/v1/menus", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("expected success true, got %v", response["success"])
	}

	data := response["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 menus, got %d", len(data))
	}
}

func TestMenuHandler_GetByID_Success(t *testing.T) {
	handler, mockRepo := setupTestMenuHandler()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/menus/:id", handler.GetByID)

	mockRepo.FindByIDFunc = func(id string) (*permission.Menu, error) {
		if id == "m1" {
			return &permission.Menu{ID: "m1", Name: "Dashboard", URL: "/dashboard"}, nil
		}
		return nil, nil // Should act as not found if error handling checks nil
	}

	req, _ := http.NewRequest("GET", "/api/v1/menus/m1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
