package patient

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(patient *Patient) (*Patient, error)
	CreateTx(tx *gorm.DB, patient *Patient) (*Patient, error)
	FindByUserID(userID string) (*Patient, error)
	FindByUserIDTx(tx *gorm.DB, userID string) (*Patient, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create a patient
func (r *repository) Create(patient *Patient) (*Patient, error) {
	if err := r.db.Create(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

// Create patient inside a transaction
func (r *repository) CreateTx(tx *gorm.DB, patient *Patient) (*Patient, error) {
	if err := tx.Create(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

// Find patient by user ID
func (r *repository) FindByUserID(userID string) (*Patient, error) {
	var patient Patient
	err := r.db.Where("user_id = ?", userID).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &patient, err
}

// Find patient by user ID inside a transaction
func (r *repository) FindByUserIDTx(tx *gorm.DB, userID string) (*Patient, error) {
	var patient Patient
	err := tx.Where("user_id = ?", userID).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &patient, err
}
