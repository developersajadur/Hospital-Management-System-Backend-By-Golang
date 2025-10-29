package repository

import (
	"errors"
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *models.Patient) (*models.Patient, error)
	CreateTx(tx *gorm.DB, patient *models.Patient) (*models.Patient, error)
	FindByUserID(userID string) (*models.Patient, error)
	FindByUserIDTx(tx *gorm.DB, userID string) (*models.Patient, error)
}

type patientRepo struct {
	db *gorm.DB
}

func PatientNewRepository(db *gorm.DB) PatientRepository {
	return &patientRepo{db: db}
}

// Create a patient
func (r *patientRepo) Create(patient *models.Patient) (*models.Patient, error) {
	if err := r.db.Create(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

// Create patient inside a transaction
func (r *patientRepo) CreateTx(tx *gorm.DB, patient *models.Patient) (*models.Patient, error) {
	if err := tx.Create(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

// Find patient by user ID
func (r *patientRepo) FindByUserID(userID string) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Where("user_id = ?", userID).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &patient, err
}

// Find patient by user ID inside a transaction
func (r *patientRepo) FindByUserIDTx(tx *gorm.DB, userID string) (*models.Patient, error) {
	var patient models.Patient
	err := tx.Where("user_id = ?", userID).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &patient, err
}
