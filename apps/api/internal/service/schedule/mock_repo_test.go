package schedule

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MockScheduleRepository
type MockScheduleRepository struct {
	FindByIDFunc func(id string) (*schedule.Schedule, error)
	ListFunc     func(req *schedule.ListSchedulesRequest) ([]schedule.Schedule, int64, error)
	CreateFunc   func(sch *schedule.Schedule) error
	UpdateFunc   func(sch *schedule.Schedule) error
	DeleteFunc   func(id string) error
}

func (m *MockScheduleRepository) FindByID(id string) (*schedule.Schedule, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockScheduleRepository) List(req *schedule.ListSchedulesRequest) ([]schedule.Schedule, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockScheduleRepository) Create(sch *schedule.Schedule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(sch)
	}
	return nil
}
func (m *MockScheduleRepository) Update(sch *schedule.Schedule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(sch)
	}
	return nil
}
func (m *MockScheduleRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// Added missing method
func (m *MockScheduleRepository) FindByDateRange(start, end time.Time, userID *string) ([]schedule.Schedule, error) {
	return nil, nil
}
// Added missing method
func (m *MockScheduleRepository) FindByUserID(userID string, start, end *time.Time) ([]schedule.Schedule, error) {
	return nil, nil
}
// Added missing method
func (m *MockScheduleRepository) FindConflicts(userID string, startTime, endTime time.Time, excludeID string) ([]schedule.Schedule, error) {
	return nil, nil
}

// MockTaskRepository
type MockTaskRepository struct {
	FindByIDFunc func(id string) (*task.Task, error)
}
func (m *MockTaskRepository) FindByID(id string) (*task.Task, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockTaskRepository) List(req *task.ListTasksRequest) ([]task.Task, int64, error) { return nil, 0, nil }
func (m *MockTaskRepository) Create(task *task.Task) error { return nil }
func (m *MockTaskRepository) Update(task *task.Task) error { return nil }
func (m *MockTaskRepository) Delete(id string) error { return nil }
func (m *MockTaskRepository) GetStatsByStatus(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockTaskRepository) GetStatsByPriority(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockTaskRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }


// MockUserRepository
type MockUserRepository struct {
	FindByIDFunc func(id string) (*user.User, error)
}
func (m *MockUserRepository) FindByID(id string) (*user.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) { return nil, nil }
func (m *MockUserRepository) Create(user *user.User) error { return nil }
func (m *MockUserRepository) Update(user *user.User) error { return nil }
func (m *MockUserRepository) Delete(id string) error { return nil }
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) { return nil, 0, nil }
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) { return 0, nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) { return nil, nil }
