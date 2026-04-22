package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	groupservice "github.com/gilabs/crm-healthcare/api/internal/service/group"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestGroupHandler() (*GroupHandler, *MockGroupRepository, *MockUserRepository) {
	mockGroupRepo := &MockGroupRepository{}
	mockUserRepo := &MockUserRepository{}
	
	svc := groupservice.NewService(mockGroupRepo, mockUserRepo)
	return NewGroupHandler(svc), mockGroupRepo, mockUserRepo
}

func TestGroupHandler_Create_Success(t *testing.T) {
	handler, mockGroupRepo, _ := setupTestGroupHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/groups", handler.Create)
	
	reqBody := group.CreateGroupRequest{
		Name: "Field Sales",
		Code: "GRP-FIELD",
	}
	
	mockGroupRepo.FindByCodeFunc = func(code string) (*group.Group, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockGroupRepo.CreateFunc = func(g *group.Group) error {
		g.ID = "grp-1"
		return nil
	}
	
	mockGroupRepo.FindByIDFunc = func(id string) (*group.Group, error) {
		return &group.Group{
			ID:        id,
			Name:      reqBody.Name,
			Code:      reqBody.Code,
			CreatedAt: time.Now(),
		}, nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGroupHandler_List_Success(t *testing.T) {
	handler, mockGroupRepo, _ := setupTestGroupHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/groups", handler.List)
	
	mockGroupRepo.ListFunc = func(req *group.ListGroupsRequest) ([]group.Group, int64, error) {
		return []group.Group{
			{ID: "g1", Name: "Group 1"},
			{ID: "g2", Name: "Group 2"},
		}, 2, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/groups?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}
