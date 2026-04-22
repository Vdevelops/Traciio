package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	brickservice "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestBrickHandler() (*BrickHandler, *MockBrickRepository, *MockUserRepository) {
	mockBrickRepo := &MockBrickRepository{}
	mockUserRepo := &MockUserRepository{}
	
	// Assuming NewService uses *MockBrickRepository and *MockUserRepository
	// Need to check if NewService in brick/service.go accepts interfaces or concrete types.
	// brick/service.go: NewService(brickRepo interfaces.BrickRepository, userRepo interfaces.UserRepository, db *gorm.DB)
	// We can pass nil for DB if not used, or mock if needed.
	
	svc := brickservice.NewService(mockBrickRepo, mockUserRepo, nil, nil)
	return NewBrickHandler(svc), mockBrickRepo, mockUserRepo
}

func TestBrickHandler_Create_Success(t *testing.T) {
	handler, mockRepo, _ := setupTestBrickHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/bricks", handler.Create)
	
	reqBody := brick.CreateBrickRequest{
		Name:     "West Java Zone 1",
		Code:     "WJ-01",
		Province: "West Java",
		Regency:  "Bandung",
		Status:   "active",
	}
	
	mockRepo.FindByCodeFunc = func(code string) (*brick.Brick, error) {
		return nil, gorm.ErrRecordNotFound
	}

	mockRepo.FindByRegencyAndProvinceFunc = func(regency, province string) (*brick.Brick, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockRepo.CreateFunc = func(b *brick.Brick) error {
		b.ID = "brick-1"
		return nil
	}
	
	mockRepo.FindByIDFunc = func(id string) (*brick.Brick, error) {
		return &brick.Brick{
			ID:        id,
			Name:      reqBody.Name,
			Code:      reqBody.Code,
			Province:  reqBody.Province,
			Regency:   reqBody.Regency,
			Status:    reqBody.Status,
			CreatedAt: time.Now(),
		}, nil
	}
	
	mockRepo.CountSalesByBrickIDFunc = func(brickID string) (int64, error) {
		return 0, nil
	}
	
	jsonVal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/bricks", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestBrickHandler_List_Success(t *testing.T) {
	handler, mockRepo, _ := setupTestBrickHandler()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/bricks", handler.List)
	
	mockRepo.ListFunc = func(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) {
		return []brick.Brick{
			{ID: "b1", Name: "Brick 1"},
			{ID: "b2", Name: "Brick 2"},
		}, 2, nil
	}
	
	mockRepo.CountSalesByBrickIDFunc = func(brickID string) (int64, error) {
		return 0, nil
	}
	
	req, _ := http.NewRequest("GET", "/api/v1/bricks?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}
