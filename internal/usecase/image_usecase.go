package usecase

import (
	"context"
	"fmt"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/validators"
	"mime/multipart"
	"sync"

	"github.com/google/uuid"
)

type ImageUsecase interface {
	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, req *dto.ImageUploadRequest) (*models.Image, error)
	UploadMultipleImages(ctx context.Context, files []multipart.File, fileHeaders []*multipart.FileHeader, req *dto.ImageUploadRequest) ([]*models.Image, []error)
	GetImageByID(id uuid.UUID) (*models.Image, error)
	GetUserImages(userID uuid.UUID, page, pageSize int) (*dto.ListResponse, error)
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type imageUsecase struct {
	repo               repository.ImageRepository
	cloudinaryUploader *helpers.CloudinaryUploader
}

func ImageNewUsecase(repo repository.ImageRepository, cloudinaryUploader *helpers.CloudinaryUploader) ImageUsecase {
	return &imageUsecase{
		repo:               repo,
		cloudinaryUploader: cloudinaryUploader,
	}
}

func (u *imageUsecase) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, req *dto.ImageUploadRequest) (*models.Image, error) {
	if err := validators.ValidateImage(fileHeader); err != nil {
		return nil, err
	}

	folder := "uploads/" + req.ImageType
	uploaded, err := u.cloudinaryUploader.UploadImage(file, fileHeader, &helpers.UploadOptions{Folder: folder})
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	image := &models.Image{
		UserID:    req.UserID,
		URL:       uploaded.URL,
		PublicID:  uploaded.PublicID,
		FileName:  uploaded.FileName,
		FileSize:  uploaded.FileSize,
		FileType:  uploaded.FileType,
		Width:     uploaded.Width,
		Height:    uploaded.Height,
		ImageType: req.ImageType,
	}

	if err := u.repo.Create(image); err != nil {
		_ = u.cloudinaryUploader.Delete(uploaded.PublicID)
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	return image, nil
}

func (u *imageUsecase) UploadMultipleImages(ctx context.Context, files []multipart.File, fileHeaders []*multipart.FileHeader, req *dto.ImageUploadRequest) ([]*models.Image, []error) {
	var (
		images []*models.Image
		errors []error
		mu     sync.Mutex
		wg     sync.WaitGroup
	)
	semaphore := make(chan struct{}, 3)

	for i := range files {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			img, err := u.UploadImage(ctx, files[idx], fileHeaders[idx], req)
			mu.Lock()
			if err != nil {
				errors = append(errors, err)
			} else {
				images = append(images, img)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return images, errors
}


// GetImageByID retrieves image by ID
func (u *imageUsecase) GetImageByID(id uuid.UUID) (*models.Image, error) {
	image, err := u.repo.FindByID(id)
	if err != nil {
		return nil, helpers.NewAppError(404, "Image not found")
	}

	return image, nil
}

// GetUserImages retrieves all images for a user with pagination
func (u *imageUsecase) GetUserImages(userID uuid.UUID, page, pageSize int) (*dto.ListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	images, total, err := u.repo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, helpers.NewAppError(500, "Failed to retrieve images")
	}

	// Convert []models.Image to []interface{}
	data := make([]interface{}, len(images))
	for i, img := range images {
		data[i] = img
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := &dto.ListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// DeleteImage deletes image from both Cloudinary and database
func (u *imageUsecase) DeleteImage(ctx context.Context, id uuid.UUID) error {
	// Get image record
	image, err := u.repo.FindByID(id)
	if err != nil {
		return helpers.NewAppError(404, "Image not found")
	}

	// Delete from Cloudinary
	if err := u.cloudinaryUploader.Delete(image.PublicID); err != nil {
		return helpers.NewAppError(500, "Failed to delete image from Cloudinary")
	}

	// Delete from database
	if err := u.repo.Delete(id); err != nil {
		return helpers.NewAppError(500, "Failed to delete image record")
	}

	return nil
}