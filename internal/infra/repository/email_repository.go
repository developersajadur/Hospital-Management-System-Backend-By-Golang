package repository

import (
	"hospital_management_system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailRepository interface {
	CreateEmail(email *models.Email) error
	UpdateEmailStatus(id uuid.UUID, status models.EmailStatus, errMsg *string) error
}

type emailRepo struct {
	db *gorm.DB
}

func EmailNewRepository(db *gorm.DB) EmailRepository {
	return &emailRepo{db: db}
}

func (r *emailRepo) CreateEmail(email *models.Email) error {
	return r.db.Create(email).Error
}

func (r *emailRepo) UpdateEmailStatus(id uuid.UUID, status models.EmailStatus, errMsg *string) error {
	return r.db.Model(&models.Email{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  status,
			"error":   errMsg,
			"sent_at": gorm.Expr("NOW()"),
		}).Error
}
