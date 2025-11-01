package validators

import (
	"fmt"
	"hospital_management_system/internal/pkg/helpers"
	"mime/multipart"
)

func ValidateImage(fileHeader *multipart.FileHeader) error {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return helpers.NewAppError(400, "Invalid file type. Only JPEG, PNG, and WebP allowed")
	}

	const maxSize = int64(10 << 20) // 10MB
	if fileHeader.Size > maxSize {
		return helpers.NewAppError(400, fmt.Sprintf("Image size too large. Max %d MB", maxSize>>20))
	}

	return nil
}
