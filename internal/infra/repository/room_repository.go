package repository

import (
	"hospital_management_system/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *models.Room) (*models.Room, error)
	GetByRoomNumber(roomNumber string) (*models.Room, error)
	GetRoomsWithFilters(roomType string, available *bool) ([]models.Room, error)
	Update(room *models.Room) (*models.Room, error)
	Delete(id string) error
}

type roomRepo struct {
	db *gorm.DB
}

func RoomNewRepository(db *gorm.DB) RoomRepository {
	return &roomRepo{db: db}
}

// Create a new room
func (r *roomRepo) Create(room *models.Room) (*models.Room, error) {
	if err := r.db.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

// Get room by room number (ignores deleted)
func (r *roomRepo) GetByRoomNumber(roomNumber string) (*models.Room, error) {
	var room models.Room
	if err := r.db.Where("room_number = ? AND is_deleted = FALSE", roomNumber).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// Get rooms by optional filters
func (r *roomRepo) GetRoomsWithFilters(roomType string, available *bool) ([]models.Room, error) {
	var rooms []models.Room
	query := r.db.Model(&models.Room{}).Where("is_deleted = FALSE")

	if roomType != "" {
		query = query.Where("type = ?", roomType)
	}

	if available != nil {
		query = query.Where("availability = ?", *available)
	}

	if err := query.Order("created_at DESC").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// Update only provided fields
func (r *roomRepo) Update(room *models.Room) (*models.Room, error) {
	if err := r.db.Model(&models.Room{}).Where("id = ? AND is_deleted = FALSE", room.ID).Updates(room).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(room, "id = ?", room.ID).Error; err != nil {
		return nil, err
	}
	return room, nil
}

// Soft delete
func (r *roomRepo) Delete(id string) error {
	return r.db.Model(&models.Room{}).Where("id = ?", id).Update("is_deleted", true).Error
}
