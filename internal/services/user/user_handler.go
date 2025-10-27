package user

import (
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/pkg/utils/jwt"
	"net/http"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{usecase: uc}
}

// Register creates a new user (doctor, patient, admin)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	utils.BodyDecoder(w, r, &req)

	// Only admins should be allowed to create doctor or admin users
	if req.Role == RoleDoctor || req.Role == RoleAdmin {
		jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
		if err != nil || jwtClaims == nil || jwtClaims.Role != RoleAdmin {
			helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized: Only admin can create doctor/admin"))
			return
		}
	}

	user, err := h.usecase.Register(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusCreated, "User registered successfully", user)
}

// Login authenticates a user and returns a JWT
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

// GetProfile returns the authenticated user's profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
	if err != nil || jwtClaims == nil {
		helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	user, err := h.usecase.GetUserByIdForAuth(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "User profile fetched successfully", user)
}
