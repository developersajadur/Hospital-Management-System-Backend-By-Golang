package handlers

import (
	"encoding/json"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type RoomHandler struct {
	roomUC   usecase.RoomUsecase
	uploader *helpers.CloudinaryUploader
}

func RoomNewHandler(roomUC usecase.RoomUsecase, uploader *helpers.CloudinaryUploader) *RoomHandler {
	return &RoomHandler{roomUC: roomUC, uploader: uploader}
}

// POST /rooms/create
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid form data"))
		return
	}

	var req dto.CreateRoomRequest
	jsonData := r.FormValue("data")
	if jsonData == "" {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Missing 'data' field in form"))
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid JSON data"))
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		uploadOpts := &helpers.UploadOptions{Folder: "room_images"}
		uploadedImage, err := h.uploader.UploadImage(file, fileHeader, uploadOpts)
		if err != nil {
			helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to upload image"))
			return
		}
		req.Image = &uploadedImage.URL
	}

	room, err := h.roomUC.Create(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "Room created successfully", room)
}

// GET /rooms/{room_number}
func (h *RoomHandler) GetByRoomNumber(w http.ResponseWriter, r *http.Request) {
	roomNumber := utils.Param(r, "room_number")
	room, err := h.roomUC.GetByRoomNumber(roomNumber)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusOK, "Room retrieved successfully", room)
}

// GET /rooms?type=icu&available=true
func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	roomType := r.URL.Query().Get("type")
	availableParam := r.URL.Query().Get("available")
	var available *bool
	if availableParam == "true" {
		val := true
		available = &val
	} else if availableParam == "false" {
		val := false
		available = &val
	}

	rooms, err := h.roomUC.GetRoomsWithFilters(roomType, available)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusOK, "Rooms fetched successfully", rooms)
}

// PUT /rooms/{id}
func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid form data"))
		return
	}

	id := utils.Param(r, "id")
	var req dto.UpdateRoomRequest
	jsonData := r.FormValue("data")
	if jsonData == "" {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Missing 'data' field in form"))
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid JSON data"))
		return
	}

	// Handle image upload
	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		uploadOpts := &helpers.UploadOptions{Folder: "room_images"}
		uploadedImage, err := h.uploader.UploadImage(file, fileHeader, uploadOpts)
		if err != nil {
			helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to upload image"))
			return
		}
		// fmt.Printf(uploadedImage.URL)
		// Update profile image URL in request
		req.Image = &uploadedImage.URL
	} else if err != http.ErrMissingFile {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid image file"))
		return
	}

	room, err := h.roomUC.Update(id, &req)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusOK, "Room updated successfully", room)
}

// DELETE /rooms/{id}
func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	if err := h.roomUC.Delete(id); err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusOK, "Room deleted successfully", nil)
}
