package auth

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) interfaces.AuthRepository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Role").Preload("Role.Permissions").Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByID(id string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Role").Preload("Role.Permissions").Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *repository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

