package usecase

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	Register(req *dto.RegisterRequest) (*models.User, error)
	FindByID(id string) (*models.User, error)   
	FindByEmail(email string) (*models.User, error)
}

type userUsecase struct {
	repo       repository.UserRepository
	doctorUC   DoctorUsecase // inject doctor usecase
	patientUC PatientUsecase
}

func UserNewUsecase(repo repository.UserRepository, doctorUC DoctorUsecase, patientUC PatientUsecase) UserUsecase {
	return &userUsecase{
		repo:     repo,
		doctorUC: doctorUC,
		patientUC: patientUC,
	}
}

func (u *userUsecase) Register(req *dto.RegisterRequest) (*models.User, error) {
	var createdUser *models.User

	txErr := u.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. Check existing user
		existing, _ := u.repo.FindByEmailTx(tx, req.Email)
		if existing != nil {
			return helpers.NewAppError(409, "User already exists")
		}

		// 2. Hash password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return helpers.NewAppError(500, "Failed to hash password")
		}

		// 3. Create user
		user := &models.User{
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: string(hashed),
			Role:     req.Role,
		}

		createdUser, err = u.repo.RegisterTx(tx, user)
		if err != nil {
			return err
		}

		// role-specific creation (doctor/patient)
		switch req.Role {
		case models.RoleDoctor:
			if req.Doctor != nil {
				doctorReq := &dto.DoctorCreateRequest{
					UserID:         createdUser.ID.String(),
					Specialization: req.Doctor.Specialization,
					Experience:     req.Doctor.Experience,
					Fee:            req.Doctor.Fee,
					ProfileImageURL:   req.Doctor.ProfileImageURL,
					Status:         req.Doctor.Status,
				}
				createdDoctor, err := u.doctorUC.CreateTx(tx, doctorReq)
				if err != nil {
					return helpers.NewAppError(500, "Failed to create doctor profile: "+err.Error())
				}
				createdUser.Doctor = createdDoctor
			}
		case models.RolePatient:
			if req.Patient != nil {
				patientReq := &dto.PatientCreateRequest{
					UserID:        createdUser.ID.String(),
					Age:           req.Patient.Age,
					Gender:        req.Patient.Gender,
					Address:       req.Patient.Address,
					ProfileImageURL:   req.Patient.ProfileImageURL,
					MedicalHistory: req.Patient.MedicalHistory,
				}
				createdPatient, err := u.patientUC.CreateTx(tx, patientReq)
				if err != nil {
					return helpers.NewAppError(500, "Failed to create patient profile: "+err.Error())
				}
				createdUser.Patient = createdPatient
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return createdUser, nil
}


func (u *userUsecase) FindByID(id string) (*models.User, error) {
		user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) FindByEmail(email string) (*models.User, error) {
		user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}