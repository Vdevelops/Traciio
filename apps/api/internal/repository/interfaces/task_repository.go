package interfaces

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
)

// TaskRepository defines the interface for task repository
type TaskRepository interface {
	// FindByID finds a task by ID
	FindByID(id string) (*task.Task, error)
	
	// List returns a list of tasks with pagination
	List(req *task.ListTasksRequest) ([]task.Task, int64, error)
	
	// Create creates a new task
	Create(task *task.Task) error
	
	// Update updates a task
	Update(task *task.Task) error
	
	// Delete soft deletes a task
	Delete(id string) error
}

// ReminderRepository defines the interface for reminder repository
type ReminderRepository interface {
	// FindByID finds a reminder by ID
	FindByID(id string) (*reminder.Reminder, error)
	
	// FindByTaskID finds reminders by task ID
	FindByTaskID(taskID string) ([]reminder.Reminder, error)
	
	// List returns a list of reminders with pagination
	List(req *reminder.ListRemindersRequest) ([]reminder.Reminder, int64, error)
	
	// Create creates a new reminder
	Create(reminder *reminder.Reminder) error
	
	// Update updates a reminder
	Update(reminder *reminder.Reminder) error
	
	// Delete soft deletes a reminder
	Delete(id string) error
	
	// FindPendingReminders finds reminders that need to be sent
	FindPendingReminders(beforeTime time.Time) ([]reminder.Reminder, error)
	
	// MarkAsSent marks a reminder as sent
	MarkAsSent(id string, sentAt time.Time) error
}

