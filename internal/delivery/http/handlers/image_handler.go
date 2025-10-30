// internal/handler/image_handler.go
package handlers

import (
	"fmt"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"hospital_management_system/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageUc usecase.ImageUsecase
}

func ImageNewHandler(imageUc usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{imageUc: imageUc}
}

// UploadImage handles image upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Failed to parse form data"))
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "No image file provided"))
		return
	}
	defer file.Close()

		jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
	if err != nil || jwtClaims == nil {
		helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	// Get form values
	imageType := r.FormValue("image_type")

	// Parse user ID
	userID, err := uuid.Parse(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	// Set default image type
	if imageType == "" {
		imageType = "general"
	}

	// Create request DTO
	req := &dto.ImageUploadRequest{
		UserID:    userID,
		ImageType: imageType,
	}

	// Upload image
	response, err := h.imageUc.UploadImage(file, fileHeader, req)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			helpers.Error(w, appErr)
			return
		}
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to upload image"))
		return
	}

	helpers.Success(w, http.StatusCreated, "Image uploaded successfully", response)
}

// UploadMultipleImages handles multiple image uploads
func (h *ImageHandler) UploadMultipleImages(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 50MB for multiple images)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Failed to parse form data"))
		return
	}

	// Get JWT claims
	jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
	if err != nil || jwtClaims == nil {
		helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	// Get image type
	imageType := r.FormValue("image_type")
	if imageType == "" {
		imageType = "general"
	}

	// Get all files from form
	files := r.MultipartForm.File["images"] // Note: "images" plural
	if len(files) == 0 {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "No images provided"))
		return
	}

	// Limit number of images (e.g., max 10)
	const maxImages = 10
	if len(files) > maxImages {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, fmt.Sprintf("Maximum %d images allowed", maxImages)))
		return
	}

	// Upload all images
	uploadedImages := make([]*models.Image, 0, len(files))
	var uploadErrors []string

	for i, fileHeader := range files {
		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: Failed to open", i+1))
			continue
		}

		// Create request DTO
		req := &dto.ImageUploadRequest{
			UserID:    userID,
			ImageType: imageType,
		}

		// Upload image
		response, err := h.imageUc.UploadImage(file, fileHeader, req)
		file.Close()

		if err != nil {
			if appErr, ok := err.(*helpers.AppError); ok {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: %s", i+1, appErr.Message))
			} else {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: Upload failed", i+1))
			}
			continue
		}

		uploadedImages = append(uploadedImages, response)
	}

	// Prepare response
	responseData := map[string]interface{}{
		"uploaded_count": len(uploadedImages),
		"total_count":    len(files),
		"images":         uploadedImages,
	}

	if len(uploadErrors) > 0 {
		responseData["errors"] = uploadErrors
	}

	if len(uploadedImages) == 0 {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "All uploads failed"))
		return
	}

	message := fmt.Sprintf("%d of %d images uploaded successfully", len(uploadedImages), len(files))
	helpers.Success(w, http.StatusCreated, message, responseData)
}

// GetImage retrieves image by ID
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid image ID"))
		return
	}

	image, err := h.imageUc.GetImageByID(id)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			helpers.Error(w, appErr)
			return
		}
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to retrieve image"))
		return
	}

	helpers.Success(w, http.StatusOK, "Image retrieved successfully", image)
}

// GetUserImages retrieves all images for a user
func (h *ImageHandler) GetUserImages(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	images, err := h.imageUc.GetUserImages(userID, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			helpers.Error(w, appErr)
			return
		}
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to retrieve images"))
		return
	}

	helpers.Success(w, http.StatusOK, "Images retrieved successfully", images)
}

// DeleteImage deletes an image
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid image ID"))
		return
	}

	err = h.imageUc.DeleteImage(id)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			helpers.Error(w, appErr)
			return
		}
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to delete image"))
		return
	}

	helpers.Success(w, http.StatusOK, "Image deleted successfully", nil)
}