package doctor

import (
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"net/http"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req DoctorCreateRequest
	utils.BodyDecoder(w, r, &req)

	doctor, err := h.usecase.Create(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "Doctor created successfully", doctor)
}
