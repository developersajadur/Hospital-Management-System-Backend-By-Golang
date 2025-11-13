package repository

import (
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) (*models.Service, error)
	GetByID(id string) (*models.Service, error)
	GetByName(name string) (*models.Service, error)
	GetAll() ([]models.Service, error)
	Update(service *models.Service) (*models.Service, error)
	Delete(id string) error
}

type serviceRepo struct {
	db *gorm.DB
}

func ServiceNewRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepo{db: db}
}

func (r *serviceRepo) Create(service *models.Service) (*models.Service, error) {
	if err := r.db.Create(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (r *serviceRepo) GetByID(id string) (*models.Service, error) {
	var service models.Service
	if err := r.db.Where("id = ? AND is_deleted = FALSE", id).First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepo) GetByName(name string) (*models.Service, error) {
	var service models.Service
	if err := r.db.Where("name = ? AND is_deleted = FALSE", name).First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepo) GetAll() ([]models.Service, error) {
	var services []models.Service
	if err := r.db.Where("is_deleted = FALSE").
		Order("created_at DESC").
		Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *serviceRepo) Update(service *models.Service) (*models.Service, error) {
	if err := r.db.Model(&models.Service{}).
		Where("id = ? AND is_deleted = FALSE", service.ID).
		Updates(service).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(service, "id = ?", service.ID).Error; err != nil {
		return nil, err
	}

	return service, nil
}

func (r *serviceRepo) Delete(id string) error {
	return r.db.Model(&models.Service{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}
