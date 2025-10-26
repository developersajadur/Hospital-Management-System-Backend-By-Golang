package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Register(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Register(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *repository) FindByID(id string) (*User, error) {
	var user User
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
