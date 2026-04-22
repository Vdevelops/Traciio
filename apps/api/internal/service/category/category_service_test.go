package category

import (
	"errors"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"gorm.io/gorm"
)

func TestService_Create_Success(t *testing.T) {
	mockRepo := &MockCategoryRepository{}
	service := NewService(mockRepo)

	req := &category.CreateCategoryRequest{
		Name: "Medicine",
		Code: "MED",
	}

	mockRepo.FindByCodeFunc = func(code string) (*category.Category, error) {
		return nil, gorm.ErrRecordNotFound
	}

	mockRepo.CreateFunc = func(c *category.Category) error {
		c.ID = "cat-1"
		return nil
	}

	mockRepo.FindByIDFunc = func(id string) (*category.Category, error) {
		return &category.Category{
			ID:        "cat-1",
			Name:      req.Name,
			Code:      req.Code,
			CreatedAt: time.Now(),
		}, nil
	}

	resp, err := service.Create(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "cat-1" {
		t.Errorf("expected ID cat-1, got %s", resp.ID)
	}
}

func TestService_Create_AlreadyExists(t *testing.T) {
	mockRepo := &MockCategoryRepository{}
	service := NewService(mockRepo)

	req := &category.CreateCategoryRequest{
		Name: "Medicine",
		Code: "MED",
	}

	mockRepo.FindByCodeFunc = func(code string) (*category.Category, error) {
		return &category.Category{ID: "existing"}, nil
	}

	_, err := service.Create(req)
	if !errors.Is(err, ErrCategoryAlreadyExists) {
		t.Errorf("expected ErrCategoryAlreadyExists, got %v", err)
	}
}
