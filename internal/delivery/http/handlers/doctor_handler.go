package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type DoctorHandler struct {
	doctorUc usecase.DoctorUsecase
}

func DoctorNewHandler(uc usecase.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{doctorUc: uc}
}

func (h *DoctorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.DoctorCreateRequest
	utils.BodyDecoder(w, r, &req)

	doctor, err := h.doctorUc.Create(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "Doctor created successfully", doctor)
}
