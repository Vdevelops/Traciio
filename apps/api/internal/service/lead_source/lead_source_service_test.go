package lead_source

import (
	"errors"
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	"gorm.io/gorm"
)

func TestService_Create_Success(t *testing.T) {
	mockRepo := &MockLeadSourceRepository{}
	// mock DB is not easily mocked as service struct expects *gorm.DB.
	// However, if we look at service implementation, it only uses repo methods, 
	// EXCEPT for Delete which uses s.db for counting usages.
	// For Create, Update, FindByID, List, it uses repo.
	
	// Issue: Service struct has `db *gorm.DB` field.
	// func NewService(repo lead_source_repo.Repository, db *gorm.DB) Service
	
	// We can pass nil for db if the method tested doesn't use it.
	// Create method uses s.repo.FindByCode and s.repo.Create. It does NOT use s.db.
	
	service := NewService(mockRepo, nil)

	req := &lead_source.CreateLeadSourceRequest{
		Name: "Referral",
		Code: "REF",
	}

	mockRepo.FindByCodeFunc = func(code string) (*lead_source.LeadSource, error) {
		return nil, gorm.ErrRecordNotFound
	}

	mockRepo.CreateFunc = func(ls *lead_source.LeadSource) error {
		ls.ID = "ls-1"
		return nil
	}

	resp, err := service.Create(req, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "ls-1" {
		t.Errorf("expected ID ls-1, got %s", resp.ID)
	}
}

func TestService_Create_AlreadyExists(t *testing.T) {
	mockRepo := &MockLeadSourceRepository{}
	service := NewService(mockRepo, nil)

	req := &lead_source.CreateLeadSourceRequest{
		Name: "Referral",
		Code: "REF",
	}

	mockRepo.FindByCodeFunc = func(code string) (*lead_source.LeadSource, error) {
		return &lead_source.LeadSource{ID: "existing-1", Code: "REF"}, nil
	}

	_, err := service.Create(req, "user-1")
	if !errors.Is(err, ErrLeadSourceCodeExists) {
		t.Errorf("expected ErrLeadSourceCodeExists, got %v", err)
	}
}

func TestService_Update_Success(t *testing.T) {
	mockRepo := &MockLeadSourceRepository{}
	service := NewService(mockRepo, nil)

	mockRepo.FindByIDFunc = func(id string) (*lead_source.LeadSource, error) {
		return &lead_source.LeadSource{
			ID:        "ls-1",
			Name:      "Old Name",
			Code:      "OLD",
			CreatedAt: time.Now(),
		}, nil
	}
	
	mockRepo.FindByCodeFunc = func(code string) (*lead_source.LeadSource, error) {
		return nil, gorm.ErrRecordNotFound
	}

	mockRepo.UpdateFunc = func(ls *lead_source.LeadSource) error {
		return nil
	}

	newName := "New Name"
	_, err := service.Update("ls-1", &lead_source.UpdateLeadSourceRequest{Name: newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Note: TestService_Delete would require a mock DB for the Count query.
// Since we cannot easily mock *gorm.DB without a driver or go-sqlmock which requires modifying NewService to take an interface or *gorm.DB from sqlmock,
// we will skip Testing Delete logic that depends on direct DB access in this unit test file, 
// or we would need to refactor service to use a repository method for checking usage.
// Refactoring is out of scope for "Writing Tests", unless necessary.
// For now, we test what we can.
