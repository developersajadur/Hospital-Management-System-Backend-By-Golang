package user

import (
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (string, error)
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Register(req *RegisterRequest) (*User, error) {
	existing, _ := u.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, helpers.NewAppError(409, "User already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, helpers.NewAppError(500, "Failed to hash password")
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),
		Role:     req.Role,
	}

	return u.repo.Register(user)
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
	token, err := utils.GenerateJWT(user.ID.String(), user.Email, user.Role, 24*time.Hour)
	if err != nil {
		return "", helpers.NewAppError(500, "Failed to generate token")
	}

	return token, nil
}
