// internal/dto/image_dto.go
package dto

import "github.com/google/uuid"

// ImageUploadRequest represents image upload request
type ImageUploadRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ImageType string    `json:"image_type"` // profile, document, general
}