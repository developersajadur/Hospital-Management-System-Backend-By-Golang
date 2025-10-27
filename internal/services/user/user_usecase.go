package user

import (
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (string, error)
	GetUserByIdForAuth(id string) (*middlewares.User, error)
}

type usecase struct {
	repo       Repository
	doctorUC   doctor.Usecase // inject doctor usecase
	patientUC patient.Usecase
}

func NewUsecase(repo Repository, doctorUC doctor.Usecase, patientUC patient.Usecase) Usecase {
	return &usecase{
		repo:     repo,
		doctorUC: doctorUC,
		patientUC: patientUC,
	}
}

func (u *usecase) Register(req *RegisterRequest) (*User, error) {
	var createdUser *User

	txErr := u.repo.(*repository).db.Transaction(func(tx *gorm.DB) error {
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
		user := &User{
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

		// 4. Create role-specific profile only for doctor or patient
		switch req.Role {
		case RoleDoctor:
			if req.Doctor != nil {
				doctorReq := &doctor.DoctorCreateRequest{
					UserID:         createdUser.ID.String(),
					Specialization: req.Doctor.Specialization,
					Experience:     req.Doctor.Experience,
					Fee:            req.Doctor.Fee,
					ProfileImage:   req.Doctor.ProfileImage,
					Status:         req.Doctor.Status,
				}

				_, err := u.doctorUC.CreateTx(tx, doctorReq)
				if err != nil {
					return helpers.NewAppError(500, "Failed to create doctor profile: "+err.Error())
				}
			}

		case RolePatient:
			if req.Patient != nil {
				patientReq := &patient.PatientCreateRequest{
					UserID:        createdUser.ID.String(),
					Age:           req.Patient.Age,
					Gender:        req.Patient.Gender,
					Address:       req.Patient.Address,
					MedicalHistory: req.Patient.MedicalHistory,
				}

				_, err := u.patientUC.CreateTx(tx, patientReq)
				if err != nil {
					return helpers.NewAppError(500, "Failed to create patient profile: "+err.Error())
				}
			}

		case RoleAdmin:
			// Admin does not need any additional table, only user
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// 5. Preload role-specific data before returning
	query := u.repo.(*repository).db
	switch createdUser.Role {
	case RoleDoctor:
		query = query.Preload("Doctor")
	case RolePatient:
		query = query.Preload("Patient")
	case RoleAdmin:
		// Admin has no additional table, just return user
	}

	if err := query.First(&createdUser, "id = ?", createdUser.ID).Error; err != nil {
		return nil, err
	}

	return createdUser, nil
}


func (u *usecase) Login(req *LoginRequest) (string, error) {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return "", helpers.NewAppError(406, "You have given a wrong email or password!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", helpers.NewAppError(406, "You have given a wrong email or password!")
	}

	if user.IsBlocked {
		return "", helpers.NewAppError(403, "User is blocked")
	}
	if user.IsDeleted {
		return "", helpers.NewAppError(403, "User is deleted")
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(user.ID.String(), user.Email, user.Role, 24*time.Hour)
	if err != nil {
		return "", helpers.NewAppError(500, "Failed to generate token")
	}

	return token, nil
}

func (u *usecase) GetUserByIdForAuth(id string) (*middlewares.User, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &middlewares.User{
		ID:         user.ID.String(),
		Email: user.Email,
		Role:       user.Role,
		IsBlocked:  user.IsBlocked,
		IsVerified: user.IsVerified,
		IsDeleted:  user.IsDeleted,
	}, nil
}
