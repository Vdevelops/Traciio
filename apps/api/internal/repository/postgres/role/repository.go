package role

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
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

func (r *repository) GetMobilePermissions(roleID string, role *roledomain.Role) (*roledomain.GetMobilePermissionsResponse, error) {
	// Check if user is admin - admin has ALL permissions
	isAdmin := role.Code == "admin"

	// Get permissions for the role
	var permissions []permission.Permission
	if isAdmin {
		// Admin: Get ALL permissions
		if err := r.db.Find(&permissions).Error; err != nil {
			return nil, err
		}
	} else {
		// Non-admin: Get permissions for the role
		permissions = role.Permissions
	}

	// Create a map of permissions by menu code for quick lookup
	permissionMap := make(map[string]map[string]bool) // menu_code -> action -> has_access

	// Mobile menu mapping: menu name/URL pattern -> mobile menu code
	mobileMenuMapping := map[string]string{
		"dashboard":     "dashboard",
		"task":          "task",
		"tasks":         "task",
		"accounts":      "accounts",
		"contacts":      "contacts",
		"visit_reports": "visit_reports",
		"visit-reports": "visit_reports",
		"schedules":     "schedules",
	}

	// Initialize map for all mobile menus
	mobileMenuCodes := []string{"dashboard", "task", "accounts", "contacts", "visit_reports", "schedules"}
	for _, menuCode := range mobileMenuCodes {
		permissionMap[menuCode] = make(map[string]bool)
		permissionMap[menuCode]["VIEW"] = false
		permissionMap[menuCode]["CREATE"] = false
		permissionMap[menuCode]["EDIT"] = false
		permissionMap[menuCode]["DELETE"] = false
	}

	// Build permission map based on role's permissions
	for _, p := range permissions {
		if p.MenuID != nil {
			// Get menu URL to identify which mobile menu this belongs to
			var menuURL string
			err := r.db.Table("menus").Where("id = ?", *p.MenuID).Pluck("url", &menuURL).Error
			if err == nil && menuURL != "" {
				// Extract menu name from URL
				menuName := ""
				parts := strings.Split(strings.Trim(menuURL, "/"), "/")
				if len(parts) > 0 {
					menuName = parts[len(parts)-1]
				}

				// Also check if URL contains the menu name
				menuURLLower := strings.ToLower(menuURL)

				// Check if this menu maps to a mobile menu
				mobileMenuCode := ""
				if code, exists := mobileMenuMapping[menuName]; exists {
					mobileMenuCode = code
				} else if strings.Contains(menuURLLower, "dashboard") {
					mobileMenuCode = "dashboard"
				} else if strings.Contains(menuURLLower, "task") {
					mobileMenuCode = "task"
				} else if strings.Contains(menuURLLower, "account") {
					mobileMenuCode = "accounts"
				} else if strings.Contains(menuURLLower, "contact") {
					mobileMenuCode = "contacts"
				} else if strings.Contains(menuURLLower, "visit") {
					mobileMenuCode = "visit_reports"
				} else if strings.Contains(menuURLLower, "schedule") {
					mobileMenuCode = "schedules"
				}

				if mobileMenuCode != "" && permissionMap[mobileMenuCode] != nil {
					// Map action to permission
					action := p.Action
					permissionMap[mobileMenuCode][action] = true

					// Also map accounts permissions to contacts (contacts are typically part of accounts)
					if mobileMenuCode == "accounts" && permissionMap["contacts"] != nil {
						permissionMap["contacts"][action] = true
					}
				}
			}
		}
	}

	// Build response
	result := &roledomain.GetMobilePermissionsResponse{
		Menus: make([]roledomain.MobileMenuPermission, 0),
	}

	for _, menuCode := range mobileMenuCodes {
		actions := make([]string, 0)
		if permissionMap[menuCode]["VIEW"] {
			actions = append(actions, "VIEW")
		}
		if permissionMap[menuCode]["CREATE"] {
			actions = append(actions, "CREATE")
		}
		if permissionMap[menuCode]["EDIT"] {
			actions = append(actions, "EDIT")
		}
		if permissionMap[menuCode]["DELETE"] {
			actions = append(actions, "DELETE")
		}

		result.Menus = append(result.Menus, roledomain.MobileMenuPermission{
			Menu:    menuCode,
			Actions: actions,
		})
	}

	return result, nil
}

func (r *repository) UpdateMobilePermissions(roleID string, req *roledomain.UpdateMobilePermissionsRequest) error {
	// Validate: Dashboard hanya bisa memiliki VIEW action
	for i := range req.Menus {
		if req.Menus[i].Menu == "dashboard" {
			// Filter hanya VIEW action untuk dashboard
			viewOnlyActions := make([]string, 0)
			for _, action := range req.Menus[i].Actions {
				if action == "VIEW" {
					viewOnlyActions = append(viewOnlyActions, action)
				}
			}
			req.Menus[i].Actions = viewOnlyActions
		}
	}

	// Get role permissions
	rolePermissions, err := r.GetPermissions(roleID)
	if err != nil {
		return err
	}

	// Get all permissions with menu info
	var allPermissions []permission.Permission
	if err := r.db.Preload("Menu").Find(&allPermissions).Error; err != nil {
		return err
	}

	// Create map of mobile menu -> actions from request
	requestedPermissions := make(map[string]map[string]bool)
	for _, menu := range req.Menus {
		requestedPermissions[menu.Menu] = make(map[string]bool)
		for _, action := range menu.Actions {
			requestedPermissions[menu.Menu][action] = true
		}
	}

	// Mobile menu mapping
	mobileMenuMapping := map[string]string{
		"dashboard":     "dashboard",
		"task":          "task",
		"tasks":         "task",
		"accounts":      "accounts",
		"contacts":      "contacts",
		"visit_reports": "visit_reports",
		"visit-reports": "visit_reports",
		"schedules":     "schedules",
	}

	// Find permissions to add/remove based on mobile menu permissions
	permissionsToKeep := make(map[string]bool)
	for _, permID := range rolePermissions {
		permissionsToKeep[permID] = true
	}

	// For each mobile menu in request, find matching permissions
	for menuCode, actions := range requestedPermissions {
		// Find permissions that match this mobile menu
		for _, perm := range allPermissions {
			if perm.MenuID == nil {
				continue
			}

			var menuURL string
			err := r.db.Table("menus").Where("id = ?", *perm.MenuID).Pluck("url", &menuURL).Error
			if err != nil || menuURL == "" {
				continue
			}

			// Check if this permission belongs to the mobile menu
			menuName := ""
			parts := strings.Split(strings.Trim(menuURL, "/"), "/")
			if len(parts) > 0 {
				menuName = parts[len(parts)-1]
			}
			menuURLLower := strings.ToLower(menuURL)

			matchedMenuCode := ""
			if code, exists := mobileMenuMapping[menuName]; exists {
				matchedMenuCode = code
			} else if strings.Contains(menuURLLower, "dashboard") {
				matchedMenuCode = "dashboard"
			} else if strings.Contains(menuURLLower, "task") {
				matchedMenuCode = "task"
			} else if strings.Contains(menuURLLower, "account") {
				matchedMenuCode = "accounts"
			} else if strings.Contains(menuURLLower, "contact") {
				matchedMenuCode = "contacts"
			} else if strings.Contains(menuURLLower, "visit") {
				matchedMenuCode = "visit_reports"
			} else if strings.Contains(menuURLLower, "schedule") {
				matchedMenuCode = "schedules"
			}

			// If this permission matches the mobile menu and action
			if matchedMenuCode == menuCode {
				// Check if action matches
				if actions[perm.Action] {
					// Keep this permission
					permissionsToKeep[perm.ID] = true
				} else {
					// Remove this permission if it's a mobile menu permission
					delete(permissionsToKeep, perm.ID)
				}
			}
		}
	}

	// Convert map to slice
	finalPermissionIDs := make([]string, 0)
	for permID, keep := range permissionsToKeep {
		if keep {
			finalPermissionIDs = append(finalPermissionIDs, permID)
		}
	}

	// Update role permissions
	return r.AssignPermissions(roleID, finalPermissionIDs)
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
