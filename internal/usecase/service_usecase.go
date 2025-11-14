package usecase

import (
	"errors"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"net/http"

	"gorm.io/gorm"
)

type ServiceUsecase interface {
	Create(req *dto.CreateServiceRequest) (*models.Service, error)
	GetByID(id string) (*models.Service, error)
	GetAll() ([]models.Service, error)
	Update(id string, req *dto.UpdateServiceRequest) (*models.Service, error)
	Delete(id string) error
}

type serviceUsecase struct {
	repo repository.ServiceRepository
}

func ServiceNewUsecase(repo repository.ServiceRepository) ServiceUsecase {
	return &serviceUsecase{repo: repo}
}

func (u *serviceUsecase) Create(req *dto.CreateServiceRequest) (*models.Service, error) {

	existing, _ := u.repo.GetByName(req.Name)
	if existing != nil {
		return nil, helpers.NewAppError(http.StatusConflict, "Service with this name already exists")
	}

	service := &models.Service{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Duration:    req.Duration,
	}

	return u.repo.Create(service)
}

func (u *serviceUsecase) GetByID(id string) (*models.Service, error) {
	service, err := u.repo.GetServiceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewAppError(http.StatusNotFound, "Service not found")
		}
		return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
	}

	return service, nil
}

func (u *serviceUsecase) GetAll() ([]models.Service, error) {
	return u.repo.GetAll()
}

func (u *serviceUsecase) Update(id string, req *dto.UpdateServiceRequest) (*models.Service, error) {
	service := &models.Service{ID: models.UUIDFromString(id)}

	if req.Name != nil {
		service.Name = *req.Name
	}
	if req.Price != nil {
		service.Price = *req.Price
	}
	if req.Description != nil {
		service.Description = *req.Description
	}
	if req.Duration != nil {
		service.Duration = *req.Duration
	}

	return u.repo.Update(service)
}

func (u *serviceUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
