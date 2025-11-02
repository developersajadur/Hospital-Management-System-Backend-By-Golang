package usecase

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
)

type RoomUsecase interface {
	Create(req *dto.CreateRoomRequest) (*models.Room, error)
	GetByRoomNumber(roomNumber string) (*models.Room, error)
	GetRoomsWithFilters(roomType string, available *bool) ([]models.Room, error)
	Update(id string, req *dto.UpdateRoomRequest) (*models.Room, error)
	Delete(id string) error
}

type roomUsecase struct {
	repo repository.RoomRepository
}

func RoomNewUsecase(repo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{repo: repo}
}

func (u *roomUsecase) Create(req *dto.CreateRoomRequest) (*models.Room, error) {
	room := &models.Room{
		RoomNumber:   req.RoomNumber,
		Type:         models.RoomType(req.Type),
		PricePerDay:  req.PricePerDay,
		Availability: req.Availability,
		Features:     req.Features,
	}

	return u.repo.Create(room)
}

func (u *roomUsecase) GetByRoomNumber(roomNumber string) (*models.Room, error) {
	return u.repo.GetByRoomNumber(roomNumber)
}

func (u *roomUsecase) GetRoomsWithFilters(roomType string, available *bool) ([]models.Room, error) {
	return u.repo.GetRoomsWithFilters(roomType, available)
}

func (u *roomUsecase) Update(id string, req *dto.UpdateRoomRequest) (*models.Room, error) {
	room := &models.Room{ID: models.UUIDFromString(id)}

	if req.RoomNumber != nil {
		room.RoomNumber = *req.RoomNumber
	}
	if req.Type != nil {
		room.Type = models.RoomType(*req.Type)
	}
	if req.PricePerDay != nil {
		room.PricePerDay = *req.PricePerDay
	}
	if req.Availability != nil {
		room.Availability = *req.Availability
	}
	if req.Features != nil {
		room.Features = *req.Features
	}

	return u.repo.Update(room)
}

func (u *roomUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
