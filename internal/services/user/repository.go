
package repository

import (
	"errors"
	"hospital_management_system/internal/services/user/model"

	"gorm.io/gorm"
)

type Repository interface {
	DB() *gorm.DB
	Register(user *model.User) (*model.User, error)
	RegisterTx(tx *gorm.DB, user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByEmailTx(tx *gorm.DB, email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DB() *gorm.DB {
	return r.db
}

func (r *repository) Register(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) RegisterTx(tx *gorm.DB, user *model.User) (*model.User, error) {
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *repository) FindByEmailTx(tx *gorm.DB, email string) (*model.User, error) {
	var user model.User
	err := tx.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *repository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
