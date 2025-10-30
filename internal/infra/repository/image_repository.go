// internal/infra/repository/image_repository.go
package repository

import (
	"hospital_management_system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(image *models.Image) error
	FindByID(id uuid.UUID) (*models.Image, error)
	FindByUserID(userID uuid.UUID, page, pageSize int) ([]models.Image, int64, error)
	FindByPublicID(publicID string) (*models.Image, error)
	Update(image *models.Image) error
	Delete(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error
}

type imageRepo struct {
	db *gorm.DB
}

func ImageNewRepository(db *gorm.DB) ImageRepository {
	return &imageRepo{db: db}
}

// Create inserts new image record
func (r *imageRepo) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

// FindByID retrieves image by ID
func (r *imageRepo) FindByID(id uuid.UUID) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("id = ? AND is_deleted = false", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// FindByUserID retrieves all images for a user with pagination
func (r *imageRepo) FindByUserID(userID uuid.UUID, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.Model(&models.Image{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Where("user_id = ? AND is_deleted = false", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&images).Error

	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// FindByPublicID retrieves image by Cloudinary public ID
func (r *imageRepo) FindByPublicID(publicID string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("public_id = ? AND is_deleted = false", publicID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// Update updates image record
func (r *imageRepo) Update(image *models.Image) error {
	return r.db.Save(image).Error
}

// Delete permanently deletes image record
func (r *imageRepo) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.Image{}, "id = ?", id).Error
}

// SoftDelete marks image as deleted
func (r *imageRepo) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Image{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.DeletedAt{},
		}).Error
}