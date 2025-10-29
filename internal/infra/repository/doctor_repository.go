package repository

import (
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	Create(doctor *models.Doctor) (*models.Doctor, error)
	CreateTx(tx *gorm.DB, doctor *models.Doctor) (*models.Doctor, error) // transaction support
}

type doctorRepo struct {
	db *gorm.DB
}

func DoctorNewRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(doctor *models.Doctor) (*models.Doctor, error) {
	if err := r.db.Create(doctor).Error; err != nil {
		return nil, err
	}
	return doctor, nil
}

func (r *doctorRepo) CreateTx(tx *gorm.DB, doctor *models.Doctor) (*models.Doctor, error) {
	if err := tx.Create(doctor).Error; err != nil {
		return nil, err
	}
	return doctor, nil
}
