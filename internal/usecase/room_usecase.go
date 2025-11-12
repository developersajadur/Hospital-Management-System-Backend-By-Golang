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

// Create Room
func (u *roomUsecase) Create(req *dto.CreateRoomRequest) (*models.Room, error) {
	existing, _ := u.repo.GetByRoomNumber(req.RoomNumber)
	if existing != nil {
		return nil, helpers.NewAppError(409, "Room with this number already exists")
	}

	room := &models.Room{
		RoomNumber:   req.RoomNumber,
		Type:         models.RoomType(req.Type),
		PricePerDay:  req.PricePerDay,
		Availability: req.Availability,
		Features:     req.Features,
		Image:        req.Image,
	}

	return u.repo.Create(room)
}

// Get a specific room
func (u *roomUsecase) GetByRoomNumber(roomNumber string) (*models.Room, error) {
	room, err := u.repo.GetByRoomNumber(roomNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewAppError(http.StatusNotFound, "Room not found")
		}
		return nil, helpers.NewAppError(http.StatusInternalServerError, "Database error")
	}
	return room, nil
}


// Filtered list of rooms
func (u *roomUsecase) GetRoomsWithFilters(roomType string, available *bool) ([]models.Room, error) {
	return u.repo.GetRoomsWithFilters(roomType, available)
}

// Update room details
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
	if req.Image != nil {
		room.Image = req.Image
	}

	return u.repo.Update(room)
}

// Soft delete
func (u *roomUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
