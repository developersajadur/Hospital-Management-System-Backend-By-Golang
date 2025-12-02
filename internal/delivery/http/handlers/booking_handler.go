package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type BookingHandler struct {
	bookingUC usecase.BookingUsecase
}

func BookingNewHandler(bookingUC usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUC: bookingUC}
}

// POST /bookings
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingRequest
	utils.BodyDecoder(w, r, &req)

	booking, err := h.bookingUC.Create(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "Booking created successfully", booking)
}

// GET /bookings/{id}
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	booking, err := h.bookingUC.GetByID(id)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Booking retrieved", booking)
}

// GET /bookings
func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.bookingUC.GetAll()
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Bookings fetched", list)
}

// PUT /bookings/{id}/status
func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	var req dto.UpdateBookingStatusRequest
	utils.BodyDecoder(w, r, &req)

	booking, err := h.bookingUC.UpdateStatus(id, &req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Status updated", booking)
}

// DELETE /bookings/{id}
func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	if err := h.bookingUC.Delete(id); err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Booking deleted", nil)
}
