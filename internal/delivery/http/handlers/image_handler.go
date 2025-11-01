// internal/handler/image_handler.go
package handlers

import (
	"context"
	"fmt"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"hospital_management_system/internal/usecase"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageUc usecase.ImageUsecase
}

func ImageNewHandler(imageUc usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{imageUc: imageUc}
}

// UploadImage handles a single image upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Content-Type:", r.Header.Get("Content-Type"))

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// fmt.Println(err)
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Failed to parse form data"))
		return
	}

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

	userID, err := uuid.Parse(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	imageType := r.FormValue("image_type")
	if imageType == "" {
		imageType = "general"
	}

	req := &dto.ImageUploadRequest{
		UserID:    userID,
		ImageType: imageType,
	}

	ctx := context.Background()
	response, err := h.imageUc.UploadImage(ctx, file, fileHeader, req)
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
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		// fmt.Println(err)
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Failed to parse form data"))
		return
	}

	jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
	if err != nil || jwtClaims == nil {
		helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	userID, err := uuid.Parse(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	imageType := r.FormValue("image_type")
	if imageType == "" {
		imageType = "general"
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "No images provided"))
		return
	}

	const maxImages = 10
	if len(files) > maxImages {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, fmt.Sprintf("Maximum %d images allowed", maxImages)))
		return
	}

	req := &dto.ImageUploadRequest{
		UserID:    userID,
		ImageType: imageType,
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var uploadedImages []*models.Image
	var uploadErrors []string

	semaphore := make(chan struct{}, 3)

	for i, fileHeader := range files {
		wg.Add(1)
		go func(idx int, fh *multipart.FileHeader) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			file, err := fh.Open()
			if err != nil {
				mu.Lock()
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: failed to open", idx+1))
				mu.Unlock()
				return
			}
			defer file.Close()

			img, err := h.imageUc.UploadImage(ctx, file, fh, req)
			mu.Lock()
			if err != nil {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: %v", idx+1, err))
			} else {
				uploadedImages = append(uploadedImages, img)
			}
			mu.Unlock()
		}(i, fileHeader)
	}

	wg.Wait()

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

	err = h.imageUc.DeleteImage(r.Context(), id)
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