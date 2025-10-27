package doctor

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(doctor *Doctor) (*Doctor, error)
	CreateTx(tx *gorm.DB, doctor *Doctor) (*Doctor, error) // transaction support
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(doctor *Doctor) (*Doctor, error) {
	if err := r.db.Create(doctor).Error; err != nil {
		return nil, err
	}
	return doctor, nil
}

func (r *repository) CreateTx(tx *gorm.DB, doctor *Doctor) (*Doctor, error) {
	if err := tx.Create(doctor).Error; err != nil {
		return nil, err
	}
	return doctor, nil
}
