package schedule

import (
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

func setupTest(t *testing.T) (*Service, *MockScheduleRepository, *MockTaskRepository, *MockUserRepository) {
	scheduleRepo := &MockScheduleRepository{}
	taskRepo := &MockTaskRepository{}
	userRepo := &MockUserRepository{}

	service := NewService(
		scheduleRepo,
		taskRepo,
		userRepo,
		nil, // googleCalendarToken service nil to avoid external calls
	)

	return service, scheduleRepo, taskRepo, userRepo
}

func TestService_CreateSchedule_Success(t *testing.T) {
	service, scheduleRepo, taskRepo, userRepo := setupTest(t)

	taskID := "task-1"
	userID := "user-1"
	
	// Mock Task check
	taskRepo.FindByIDFunc = func(id string) (*task.Task, error) {
		return &task.Task{ID: id, AssignedTo: &userID}, nil
	}

	// Mock User check
	userRepo.FindByIDFunc = func(id string) (*user.User, error) {
		return &user.User{ID: id}, nil
	}

	// Mock Create
	scheduleRepo.CreateFunc = func(sch *schedule.Schedule) error {
		if sch.Title != "Meeting" {
			t.Errorf("expected title Meeting, got %s", sch.Title)
		}
		if *sch.TaskID != taskID {
			t.Errorf("expected task ID %s, got %v", taskID, sch.TaskID)
		}
		return nil
	}
	
	// Mock FindByID reload
	scheduleRepo.FindByIDFunc = func(id string) (*schedule.Schedule, error) {
		return &schedule.Schedule{
			ID:          id,
			Title:       "Meeting",
			TaskID:      &taskID,
			ScheduledAt: time.Now().Add(24 * time.Hour),
		}, nil
	}

	req := &schedule.CreateScheduleRequest{
		TaskID:      taskID,
		Title:       "Meeting",
		ScheduledAt: time.Now().Add(24 * time.Hour),
		SyncToGoogleCalendar: false, // Important: keep false to avoid nil panic
	}

	resp, err := service.CreateSchedule(req, "creator-1")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
}

func TestService_ListSchedules_Success(t *testing.T) {
	service, scheduleRepo, _, _ := setupTest(t)

	scheduleRepo.ListFunc = func(req *schedule.ListSchedulesRequest) ([]schedule.Schedule, int64, error) {
		return []schedule.Schedule{
			{ID: "sch-1", Title: "Sch 1"},
			{ID: "sch-2", Title: "Sch 2"},
		}, 2, nil
	}

	req := &schedule.ListSchedulesRequest{
		Page:    1,
		PerPage: 10,
	}

	resp, pagination, err := service.ListSchedules(req)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if pagination.Total != 2 {
		t.Errorf("expected total 2, got %d", pagination.Total)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 schedules, got %d", len(resp))
	}
}

func TestService_GetScheduleByID_Success(t *testing.T) {
	service, scheduleRepo, _, _ := setupTest(t)

	sid := "sch-1"
	scheduleRepo.FindByIDFunc = func(id string) (*schedule.Schedule, error) {
		return &schedule.Schedule{ID: id, Title: "Test"}, nil
	}

	resp, err := service.GetScheduleByID(sid)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.ID != sid {
		t.Errorf("expected ID %s, got %s", sid, resp.ID)
	}
}

func TestService_DeleteSchedule_Success(t *testing.T) {
	service, scheduleRepo, _, _ := setupTest(t)

	sid := "sch-1"
	// Mock FindByID to return non-synced schedule
	scheduleRepo.FindByIDFunc = func(id string) (*schedule.Schedule, error) {
		return &schedule.Schedule{
			ID:                       id,
			GoogleCalendarSyncStatus: "not_synced",
		}, nil
	}

	scheduleRepo.DeleteFunc = func(id string) error {
		if id != sid {
			t.Errorf("expected delete ID %s, got %s", sid, id)
		}
		return nil
	}

	err := service.DeleteSchedule(sid)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}
