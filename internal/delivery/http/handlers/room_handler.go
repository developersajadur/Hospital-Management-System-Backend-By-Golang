package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type RoomHandler struct {
	roomUC usecase.RoomUsecase
}

func RoomNewHandler(roomUC usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{roomUC: roomUC}
}

// POST /rooms/create
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomRequest
	utils.BodyDecoder(w, r, &req)

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
	id := utils.Param(r, "id")
	var req dto.UpdateRoomRequest
	utils.BodyDecoder(w, r, &req)

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
