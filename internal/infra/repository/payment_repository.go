package repository

import (
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByTranID(tranID string) (*models.Payment, error)
	Update(payment *models.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func PaymentNewRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByTranID(tranID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("tran_id = ?", tranID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}
