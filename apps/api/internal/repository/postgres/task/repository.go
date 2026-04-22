package task

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

// Field selections for optimized queries (reduce memory usage)
const (
	userFields    = "id, name, email, avatar_url"
	accountFields = "id, name"
	contactFields = "id, name, email, phone"
	dealFields    = "id, title, value, status"
	leadFields    = "id, first_name, last_name, email"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new task repository
func NewRepository(db *gorm.DB) interfaces.TaskRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*task.Task, error) {
	var t task.Task
	err := r.db.
		Preload("AssignedUser").
		Preload("AssignedFromUser").
		Preload("Account").
		Preload("Contact").
		Preload("Deal").
		Preload("Lead").
		Where("id = ?", id).
		First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) List(req *task.ListTasksRequest) ([]task.Task, int64, error) {
	var tasks []task.Task
	var total int64

	query := r.db.Model(&task.Task{})

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("assigned_to IN ?", req.ScopedUserIDs)
	}

	// Apply filters
	if req.Search != "" {
		// Use FTS search
		searchQuery := strings.Join(strings.Fields(req.Search), " & ") + ":*"
		query = query.Where("to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ to_tsquery('english', ?)", searchQuery)
	}

	if req.Status != "" {
		// Support comma-separated status values (e.g., "pending,in_progress")
		statuses := strings.Split(req.Status, ",")
		if len(statuses) > 1 {
			// Multiple statuses - use IN clause
			query = query.Where("status IN ?", statuses)
		} else {
			// Single status - use equality
			query = query.Where("status = ?", req.Status)
		}
	}

	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if req.AssignedTo != "" {
		query = query.Where("assigned_to = ?", req.AssignedTo)
	}

	if req.AccountID != "" {
		query = query.Where("account_id = ?", req.AccountID)
	}

	if req.ContactID != "" {
		query = query.Where("contact_id = ?", req.ContactID)
	}

	if req.DealID != "" {
		query = query.Where("deal_id = ?", req.DealID)
	}

	if req.LeadID != "" {
		query = query.Where("lead_id = ?", req.LeadID)
	}

	if req.TaskSource != "" {
		query = query.Where("task_source = ?", req.TaskSource)
	}

	if req.IsSchedule != nil {
		query = query.Where("is_schedule_task = ?", *req.IsSchedule)
	}

	if req.DueDateFrom != nil {
		query = query.Where("due_date >= ?", *req.DueDateFrom)
	}

	if req.DueDateTo != nil {
		query = query.Where("due_date <= ?", *req.DueDateTo)
	}

	if req.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *req.CreatedFrom)
	}

	if req.CreatedTo != nil {
		// Add 1 day to include the end date fully if it's just a date
		// But if it's passed as a timestamp from valid logic, just use it
		// Here we assume it's passed correctly, or we can ensuring it covers the day
		query = query.Where("created_at <= ?", *req.CreatedTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	// Fetch data with preload
	err := query.
		Preload("AssignedUser").
		Preload("AssignedFromUser").
		Preload("Account").
		Preload("Contact").
		Preload("Deal").
		Preload("Lead").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *repository) Create(t *task.Task) error {
	return r.db.Create(t).Error
}

func (r *repository) Update(t *task.Task) error {
	// Clear relations to avoid updating them
	t.AssignedUser = nil
	t.AssignedFromUser = nil
	t.Account = nil
	t.Contact = nil
	t.Deal = nil
	t.Lead = nil

	return r.db.Model(t).Omit("AssignedUser", "AssignedFromUser", "Account", "Contact", "Deal", "Lead").Updates(t).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&task.Task{}).Error
}
