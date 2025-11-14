package repository

import (
	"hospital_management_system/internal/models"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(b *models.Booking) (*models.Booking, error)
	GetByID(id string) (*models.Booking, error)
	GetAll() ([]models.Booking, error)
	Update(b *models.Booking) (*models.Booking, error)
	Delete(id string) error

	CountServiceBookingsForDay(serviceID string, day string) (int64, error)
}

type bookingRepo struct {
	db *gorm.DB
}

func BookingNewRepository(db *gorm.DB) BookingRepository {
	return &bookingRepo{db: db}
}

func (r *bookingRepo) Create(b *models.Booking) (*models.Booking, error) {
	return b, r.db.Create(b).Error
}

func (r *bookingRepo) GetByID(id string) (*models.Booking, error) {
	var b models.Booking
	if err := r.db.Where("id = ? AND is_deleted = FALSE", id).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *bookingRepo) GetAll() ([]models.Booking, error) {
	var list []models.Booking
	err := r.db.Where("is_deleted = FALSE").Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *bookingRepo) Update(b *models.Booking) (*models.Booking, error) {
	if err := r.db.Model(&models.Booking{}).
		Where("id = ?", b.ID).
		Updates(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}

func (r *bookingRepo) Delete(id string) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *bookingRepo) CountServiceBookingsForDay(serviceID string, day string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("service_id = ? AND DATE(scheduled_at) = ? AND is_deleted = FALSE", serviceID, day).
		Count(&count).Error
	return count, err
}
