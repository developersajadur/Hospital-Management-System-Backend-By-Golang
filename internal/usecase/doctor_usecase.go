package usecase

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoctorUsecase interface {
	Create(req *dto.DoctorCreateRequest) (*models.Doctor, error)
	CreateTx(tx *gorm.DB, req *dto.DoctorCreateRequest) (*models.Doctor, error)
}

type doctorUsecase struct {
	repo repository.DoctorRepository
}

func DoctorNewUsecase(repo repository.DoctorRepository) DoctorUsecase {
	return &doctorUsecase{repo: repo}
}

func (u *doctorUsecase) Create(req *dto.DoctorCreateRequest) (*models.Doctor, error) {
	return u.CreateTx(nil, req) // call CreateTx without transaction
}

// CreateTx supports optional transaction
func (u *doctorUsecase) CreateTx(tx *gorm.DB, req *dto.DoctorCreateRequest) (*models.Doctor, error) {
	// Convert UserID string to UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, helpers.NewAppError(400, "Invalid user_id")
	}

	if req.Experience < 0 {
		return nil, helpers.NewAppError(400, "Experience cannot be negative")
	}
	if req.Fee < 0 {
		return nil, helpers.NewAppError(400, "Fee cannot be negative")
	}

	doctor := &models.Doctor{
		UserID:         userID,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Fee:            req.Fee,
		ProfileImageId:   req.ProfileImageId,
		Status:         req.Status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	var created *models.Doctor
	if tx != nil {
		created, err = u.repo.CreateTx(tx, doctor)
	} else {
		created, err = u.repo.Create(doctor)
	}

	if err != nil {
		return nil, helpers.NewAppError(500, err.Error())
	}
	return created, nil
}
