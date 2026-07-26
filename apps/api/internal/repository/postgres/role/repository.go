package role

import (
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new role repository
func NewRepository(db *gorm.DB) interfaces.RoleRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*roledomain.Role, error) {
	var ro roledomain.Role
	err := r.db.Preload("Permissions").
		Select("*, (SELECT COUNT(*) FROM users WHERE users.role_id = roles.id) as user_count").
		Where("id = ?", id).First(&ro).Error
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *repository) FindByCode(code string) (*roledomain.Role, error) {
	var ro roledomain.Role
	err := r.db.Preload("Permissions").
		Select("*, (SELECT COUNT(*) FROM users WHERE users.role_id = roles.id) as user_count").
		Where("code = ?", code).First(&ro).Error
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *repository) List() ([]roledomain.Role, error) {
	var roles []roledomain.Role
	err := r.db.Preload("Permissions").
		Select("*, (SELECT COUNT(*) FROM users WHERE users.role_id = roles.id) as user_count").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) Create(ro *roledomain.Role) error {
	return r.db.Create(ro).Error
}

func (r *repository) Update(ro *roledomain.Role) error {
	return r.db.Save(ro).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&roledomain.Role{}).Error
}

func (r *repository) AssignPermissions(roleID string, permissionIDs []string) error {
	// First, clear existing permissions
	if err := r.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// Then assign new permissions
	if len(permissionIDs) > 0 {
		for _, permID := range permissionIDs {
			if err := r.db.Exec(
				"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
				roleID, permID,
			).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *repository) GetPermissions(roleID string) ([]string, error) {
	var permissionIDs []string
	err := r.db.Table("role_permissions").
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	if err != nil {
		return nil, err
	}
	return permissionIDs, nil
}

// GetScopesByRoleID returns all data-visibility scope entries for a given role
func (r *repository) GetScopesByRoleID(roleID string) ([]roledomain.RoleScope, error) {
	var scopes []roledomain.RoleScope
	if err := r.db.Where("role_id = ?", roleID).Find(&scopes).Error; err != nil {
		return nil, err
	}
	return scopes, nil
}

// UpsertScopes creates or updates scope entries for a role.
// Each item is matched by (role_id, resource) and the scope value is set accordingly.
func (r *repository) UpsertScopes(roleID string, scopes []roledomain.RoleScopeItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range scopes {
			var existing roledomain.RoleScope
			err := tx.Where("role_id = ? AND resource = ?", roleID, item.Resource).First(&existing).Error
			if err == nil {
				// Update existing scope
				if err := tx.Model(&existing).Update("scope", item.Scope).Error; err != nil {
					return err
				}
			} else {
				// Create new scope entry
				newScope := roledomain.RoleScope{
					RoleID:   roleID,
					Resource: item.Resource,
					Scope:    item.Scope,
				}
				if err := tx.Create(&newScope).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
