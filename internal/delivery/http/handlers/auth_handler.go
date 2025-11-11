package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

// Handler handles user-related HTTP requests
type AuthHandler struct {
	authUc   usecase.AuthUsecase
}

// NewHandler creates a new User Handler
func AuthNewHandler(
	authUc   usecase.AuthUsecase,
) *AuthHandler {
	return &AuthHandler{
		authUc:   authUc,
	}
}

// Login authenticates a user and returns a JWT
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	utils.BodyDecoder(w, r, &req)

	token, err := h.authUc.Login(&req)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	helpers.Success(w, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}