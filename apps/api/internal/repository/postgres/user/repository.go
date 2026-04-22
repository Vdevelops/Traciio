package user

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) interfaces.UserRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Role").Preload("Role.Permissions").Preload("Group").Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Role").Preload("Role.Permissions").Preload("Group").Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) List(req *user.ListUsersRequest) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	query := r.db.Model(&user.User{}).Preload("Role").Preload("Group")

	// Apply filters
	if req.Search != "" {
		// Optimized: Use Full Text Search instead of LIKE %...%
		// Matches name or email using GIN index
		// Also checks status and role name (kept as fallback/OR for now, though less efficient)
		query = query.Where(
			"to_tsvector('english', users.name || ' ' || users.email) @@ plainto_tsquery('english', ?) OR users.status = ? OR EXISTS (SELECT 1 FROM roles WHERE roles.id = users.role_id AND roles.name = ?)",
			req.Search, req.Search, req.Search,
		)
	}

	if req.Status != "" {
		query = query.Where("users.status = ?", req.Status)
	}

	if req.RoleID != "" {
		query = query.Where("users.role_id = ?", req.RoleID)
	}

	if req.GroupID != "" {
		query = query.Where("users.group_id = ?", req.GroupID)
	}

	if req.BrickID != "" {
		query = query.Where("users.brick_id = ?", req.BrickID)
	}

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("users.id IN ?", req.ScopedUserIDs)
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

	// Fetch data
	err := query.Order("users.created_at DESC").Offset(offset).Limit(perPage).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *repository) Update(u *user.User) error {
	// Clear the Role and Group associations to prevent GORM from overwriting IDs
	u.Role = nil
	u.Group = nil
	return r.db.Save(u).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&user.User{}).Error
}

func (r *repository) CountUsersByRoleID(roleID string) (int64, error) {
	var count int64
	err := r.db.Model(&user.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

func (r *repository) GetUsersByGroupID(groupID string) ([]user.User, error) {
	var users []user.User
	err := r.db.Model(&user.User{}).Preload("Role").Preload("Group").Where("group_id = ?", groupID).Find(&users).Error
	return users, err
}

func (r *repository) GetUsersByBrickID(brickID string) ([]user.User, error) {
	var users []user.User
	err := r.db.Model(&user.User{}).Preload("Role").Where("brick_id = ?", brickID).Find(&users).Error
	return users, err
}

func (r *repository) GetUsersByRoleID(roleID string) ([]string, error) {
	var userIDs []string
	err := r.db.Model(&user.User{}).Where("role_id = ?", roleID).Pluck("id", &userIDs).Error
	return userIDs, err
}

