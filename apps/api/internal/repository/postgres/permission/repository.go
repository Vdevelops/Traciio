package permission

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new permission repository
func NewRepository(db *gorm.DB) interfaces.PermissionRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*permission.Permission, error) {
	var p permission.Permission
	err := r.db.Preload("Menu").Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindByCode(code string) (*permission.Permission, error) {
	var p permission.Permission
	err := r.db.Preload("Menu").Where("code = ?", code).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List() ([]permission.Permission, error) {
	var permissions []permission.Permission
	err := r.db.Preload("Menu").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *repository) GetByMenuID(menuID string) ([]permission.Permission, error) {
	var permissions []permission.Permission
	err := r.db.Where("menu_id = ?", menuID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *repository) GetByRoleID(roleID string) ([]permission.Permission, error) {
	var permissions []permission.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}


// MenuRepository implementation
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository
func NewMenuRepository(db *gorm.DB) interfaces.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindByID(id string) (*permission.Menu, error) {
	var m permission.Menu
	err := r.db.Preload("Parent").Preload("Children").Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepository) FindByURL(url string) (*permission.Menu, error) {
	var m permission.Menu
	err := r.db.Preload("Parent").Preload("Children").Where("url = ?", url).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepository) List() ([]permission.Menu, error) {
	var menus []permission.Menu
	err := r.db.Preload("Parent").Preload("Children").Order("\"order\" ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) GetRootMenus() ([]permission.Menu, error) {
	var menus []permission.Menu
	err := r.db.Where("parent_id IS NULL").Preload("Children").Order("\"order\" ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) Create(m *permission.Menu) error {
	return r.db.Create(m).Error
}

func (r *menuRepository) Update(m *permission.Menu) error {
	return r.db.Save(m).Error
}

func (r *menuRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&permission.Menu{}).Error
}

