// internal/services/email/usecase.go
package usecase

import (
	"fmt"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"

	"github.com/google/uuid"
)

type EmailUsecase interface {
	CreateEmail(userID uuid.UUID, to, subject, body string, typ models.EmailType) (models.Email, error)
}

type emailUsecase struct {
	repo repository.EmailRepository
}

func EmailNewUsecase(repo repository.EmailRepository) EmailUsecase {
	return &emailUsecase{repo: repo}
}

func (u *emailUsecase) CreateEmail(userID uuid.UUID, to, subject, body string, typ models.EmailType) (models.Email, error) {

	email := &models.Email{
		UserID:  userID,
		Email:   to,
		Subject: subject,
		Body:    body,
		Type:    typ,
		Status:  models.EmailStatusPending,
	}

	if err := u.repo.CreateEmail(email); err != nil {
		return models.Email{}, fmt.Errorf("failed to create email: %w", err)
	}

	return *email, nil
}
