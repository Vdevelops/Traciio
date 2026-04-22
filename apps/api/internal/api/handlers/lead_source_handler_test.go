package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	lead_source_service "github.com/gilabs/crm-healthcare/api/internal/service/lead_source"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestLeadSourceHandler() (*LeadSourceHandler, *MockLeadSourceRepository) {
	mockRepo := &MockLeadSourceRepository{}
	// For integration tests of handler, we mock the service dependencies.
	// However, the handler takes the service interface.
	// We are instantiating the REAL service with MOCK repo.
	
	svc := lead_source_service.NewService(mockRepo, nil)
	return NewLeadSourceHandler(svc), mockRepo
}

func TestLeadSourceHandler_Create_Success(t *testing.T) {
	handler, mockRepo := setupTestLeadSourceHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/lead-sources", handler.Create)
	
	reqBody := lead_source.CreateLeadSourceRequest{
		Name: "Facebook Ads",
		Code: "FB_ADS",
		Order: 1,
	}
	
	mockRepo.FindByCodeFunc = func(code string) (*lead_source.LeadSource, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockRepo.CreateFunc = func(ls *lead_source.LeadSource) error {
		ls.ID = "ls-1"
		return nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/lead-sources", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestLeadSourceHandler_List_Success(t *testing.T) {
	handler, mockRepo := setupTestLeadSourceHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/lead-sources", handler.List)
	
	mockRepo.ListFunc = func(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error) {
		return []*lead_source.LeadSource{
			{ID: "ls1", Name: "Source 1", CreatedAt: time.Now()},
			{ID: "ls2", Name: "Source 2", CreatedAt: time.Now()},
		}, 2, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/lead-sources?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestLeadSourceHandler_ListAll_Success(t *testing.T) {
	handler, mockRepo := setupTestLeadSourceHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/lead-sources/all", handler.ListAll)
	
	mockRepo.ListAllFunc = func() ([]*lead_source.LeadSource, error) {
		return []*lead_source.LeadSource{
			{ID: "ls1", Name: "Source 1"},
			{ID: "ls2", Name: "Source 2"},
		}, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/lead-sources/all", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}
