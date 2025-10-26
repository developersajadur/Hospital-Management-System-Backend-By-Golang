package user

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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	utils.BodyDecoder(w, r, &req)

	user, err := h.usecase.Register(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, 201, "User registered successfully", user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	utils.BodyDecoder(w, r, &req)

	token, err := h.usecase.Login(&req)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	helpers.Success(w, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}
