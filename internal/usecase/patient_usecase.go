package usecase

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientUsecase interface {
	Create(req *dto.PatientCreateRequest) (*models.Patient, error)
	CreateTx(txTx interface{}, req *dto.PatientCreateRequest) (*models.Patient, error) 
}

type patientUsecase struct {
	repo repository.PatientRepository
}

func PatientNewUsecase(repo repository.PatientRepository) PatientUsecase {
	return &patientUsecase{repo: repo}
}

// Create patient
func (u *patientUsecase) Create(req *dto.PatientCreateRequest) (*models.Patient, error) {
	patient := &models.Patient{
		UserID:         uuidFromString(req.UserID),
		Age:            req.Age,
		Gender:         models.Gender(req.Gender),
		Address:        req.Address,
		MedicalHistory: req.MedicalHistory,
	}
	return u.repo.Create(patient)
}

// Create patient in transaction
func (u *patientUsecase) CreateTx(tx interface{}, req *dto.PatientCreateRequest) (*models.Patient, error) {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, helpers.NewAppError(500, "Invalid transaction object")
	}

	patient := &models.Patient{
		UserID:         uuidFromString(req.UserID),
		Age:            req.Age,
		Gender:         models.Gender(req.Gender),
		Address:        req.Address,
		MedicalHistory: req.MedicalHistory,
	}

	return u.repo.CreateTx(gormTx, patient)
}

// helper to convert string to UUID
func uuidFromString(id string) uuid.UUID {
	uid, _ := uuid.Parse(id)
	return uid
}
