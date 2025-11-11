package usecase

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)


type AuthUsecase interface {
	Login(req *dto.LoginRequest) (string, error)
}

type authUsecase struct {
	repo       repository.UserRepository
}

func AuthNewUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		repo:     repo,
	}
}




func (u *authUsecase) Login(req *dto.LoginRequest) (string, error) {
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
		if !user.IsVerified {
		return "", helpers.NewAppError(403, "User is not verify")
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(user.ID.String(), user.Email, user.Role, 24*time.Hour)
	if err != nil {
		return "", helpers.NewAppError(500, "Failed to generate token")
	}

	return token, nil
}

