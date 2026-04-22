package events

import "time"

// TaskCreatedEvent is emitted when a new task is created
type TaskCreatedEvent struct {
	TaskID       string     `json:"task_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	AssignedTo   string     `json:"assigned_to"`
	AssignedFrom string     `json:"assigned_from"`
	AccountID    string     `json:"account_id,omitempty"`
	ContactID    string     `json:"contact_id,omitempty"`
	DealID       string     `json:"deal_id,omitempty"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TaskAssignedEvent is emitted when a task is assigned to a user
type TaskAssignedEvent struct {
	TaskID       string    `json:"task_id"`
	Title        string    `json:"title"`
	OldAssignee  string    `json:"old_assignee,omitempty"`
	NewAssignee  string    `json:"new_assignee"`
	AssignedBy   string    `json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
}

// TaskCompletedEvent is emitted when a task is marked as completed
type TaskCompletedEvent struct {
	TaskID      string    `json:"task_id"`
	Title       string    `json:"title"`
	AssignedTo  string    `json:"assigned_to"`
	CompletedBy string    `json:"completed_by"`
	CompletedAt time.Time `json:"completed_at"`
}
