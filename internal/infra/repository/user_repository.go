package repository

import (
	"errors"
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

// Repository interface
type UserRepository interface {
		GetDB() *gorm.DB
	Register(user *models.User) (*models.User, error)
	RegisterTx(tx *gorm.DB, user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)        // user with doctor preloaded
	FindByEmailTx(tx *gorm.DB, email string) (*models.User, error)
	FindByID(id string) (*models.User, error)             // user with doctor preloaded
}

type userRepo struct {
	db *gorm.DB
}

func UserNewRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}
func (r *userRepo) GetDB() *gorm.DB {
	return r.db
}
// Create a new user
func (r *userRepo) Register(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) RegisterTx(tx *gorm.DB, user *models.User) (*models.User, error) {
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Find user by email
func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Find user by email inside a transaction
func (r *userRepo) FindByEmailTx(tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Find user by ID
func (r *userRepo) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
