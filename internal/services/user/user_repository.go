package user

import (
	"errors"

	"gorm.io/gorm"
)

// Repository interface
type Repository interface {
	Register(user *User) (*User, error)
	RegisterTx(tx *gorm.DB, user *User) (*User, error)
	FindByEmail(email string) (*User, error)        // user with doctor preloaded
	FindByEmailTx(tx *gorm.DB, email string) (*User, error)
	FindByID(id string) (*User, error)             // user with doctor preloaded
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create a new user
func (r *repository) Register(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	// preload related data
	switch user.Role {
	case RoleDoctor:
		r.db.Preload("Doctor").First(user, "id = ?", user.ID)
	case RolePatient:
		r.db.Preload("Patient").First(user, "id = ?", user.ID)
	}

	return user, nil
}

func (r *repository) RegisterTx(tx *gorm.DB, user *User) (*User, error) {
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}

	// preload related data
	switch user.Role {
	case RoleDoctor:
		tx.Preload("Doctor").First(user, "id = ?", user.ID)
	case RolePatient:
		tx.Preload("Patient").First(user, "id = ?", user.ID)
	}

	return user, nil
}

// Find user by email and preload doctor if role is doctor
func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Preload("Doctor").Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Find user by email inside a transaction
func (r *repository) FindByEmailTx(tx *gorm.DB, email string) (*User, error) {
	var user User
	err := tx.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Find user by ID and preload doctor if role is doctor
func (r *repository) FindByID(id string) (*User, error) {
	var user User
	err := r.db.Preload("Doctor").Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
