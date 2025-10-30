// internal/usecase/image_usecase.go
package usecase

import (
	"context"
	"fmt"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type ImageUsecase interface {
	UploadImage(file multipart.File, fileHeader *multipart.FileHeader, req *dto.ImageUploadRequest) (*models.Image, error)
UploadMultipleImages(files []multipart.File, fileHeaders []*multipart.FileHeader, req *dto.ImageUploadRequest) ([]*models.Image, []error)
	GetImageByID(id uuid.UUID) (*models.Image, error)
	GetUserImages(userID uuid.UUID, page, pageSize int) (*dto.ListResponse, error)
	DeleteImage(id uuid.UUID) error
}

type imageUsecase struct {
	repo repository.ImageRepository
	cld  *cloudinary.Cloudinary
}

func ImageNewUsecase(repo repository.ImageRepository, cld *cloudinary.Cloudinary) ImageUsecase {
	return &imageUsecase{
		repo: repo,
		cld:  cld,
	}
}

// UploadImage uploads image to Cloudinary and saves record to database
func (u *imageUsecase) UploadImage(file multipart.File, fileHeader *multipart.FileHeader, req *dto.ImageUploadRequest) (*models.Image, error) {
	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, helpers.NewAppError(400, "Invalid file type. Only JPEG, PNG, and WebP allowed")
	}

	// Validate file size (max 10MB)
	const imageMaxSize = int64(10 << 20) // 10MB
	if fileHeader.Size > imageMaxSize {
		return nil, helpers.NewAppError(400, "Image size too large. Max 10MB")
	}

	// Set default folder if not provided
		folder := fmt.Sprintf("uploads/%s", req.ImageType)
	

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Upload to Cloudinary
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       strings.TrimSuffix(fileName, ext),
		Folder:         folder,
		ResourceType:   "image",
		Transformation: "q_auto,f_auto",
	})

	if err != nil {
		return nil, helpers.NewAppError(500, fmt.Sprintf("Failed to upload image: %v", err))
	}

	// Create image record
	image := &models.Image{
		UserID:    req.UserID,
		URL:       uploadResult.SecureURL,
		PublicID:  uploadResult.PublicID,
		FileName:  fileHeader.Filename,
		FileSize:  fileHeader.Size,
		FileType:  contentType,
		Width:     uploadResult.Width,
		Height:    uploadResult.Height,
		ImageType: req.ImageType,
	}

	// Save to database
	if err := u.repo.Create(image); err != nil {
		// If database save fails, delete from Cloudinary
		u.deleteFromCloudinary(uploadResult.PublicID)
		return nil, helpers.NewAppError(500, "Failed to save image record")
	}

	return image, nil
}

// UploadMultipleImages uploads multiple images concurrently
func (u *imageUsecase) UploadMultipleImages(files []multipart.File, fileHeaders []*multipart.FileHeader, req *dto.ImageUploadRequest) ([]*models.Image, []error) {
	var (
		images []*models.Image
		errors []error
		mu     sync.Mutex
		wg     sync.WaitGroup
	)

	// Use semaphore to limit concurrent uploads
	semaphore := make(chan struct{}, 3) // Max 3 concurrent uploads

	for i := range files {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			image, err := u.UploadImage(files[idx], fileHeaders[idx], req)
			
			mu.Lock()
			if err != nil {
				errors = append(errors, err)
			} else {
				images = append(images, image)
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
func (u *imageUsecase) DeleteImage(id uuid.UUID) error {
	// Get image record
	image, err := u.repo.FindByID(id)
	if err != nil {
		return helpers.NewAppError(404, "Image not found")
	}

	// Delete from Cloudinary
	if err := u.deleteFromCloudinary(image.PublicID); err != nil {
		return helpers.NewAppError(500, "Failed to delete image from Cloudinary")
	}

	// Delete from database
	if err := u.repo.Delete(id); err != nil {
		return helpers.NewAppError(500, "Failed to delete image record")
	}

	return nil
}

// deleteFromCloudinary helper function to delete image from Cloudinary
func (u *imageUsecase) deleteFromCloudinary(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})

	return err
}