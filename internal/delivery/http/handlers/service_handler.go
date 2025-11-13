package handlers

import (
	"encoding/json"
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type ServiceHandler struct {
	serviceUC usecase.ServiceUsecase
}

func ServiceNewHandler(serviceUC usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{serviceUC: serviceUC}
}

// POST /services
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid JSON body"))
		return
	}

	service, err := h.serviceUC.Create(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "Service created successfully", service)
}

// GET /services/{id}
func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	service, err := h.serviceUC.GetByID(id)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Service retrieved successfully", service)
}

// GET /services
func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	services, err := h.serviceUC.GetAll()
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusOK, "All services fetched successfully", services)
}

// PUT /services/{id}
func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")
	var req dto.UpdateServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid JSON body"))
		return
	}

	service, err := h.serviceUC.Update(id, &req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Service updated successfully", service)
}

// DELETE /services/{id}
func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := utils.Param(r, "id")

	if err := h.serviceUC.Delete(id); err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Service deleted successfully", nil)
}
