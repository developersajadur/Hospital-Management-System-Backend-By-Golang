package handlers

import (
	"encoding/json"
	"net/http"

	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"hospital_management_system/internal/usecase"
)

// Handler handles user-related HTTP requests
type UserHandler struct {
	userUc    usecase.UserUsecase
	otpUC     usecase.OtpUsecase
	emailUC   usecase.EmailUsecase
	publisher *rabbitmq.Publisher
	uploader  *helpers.CloudinaryUploader
}

// NewHandler creates a new User Handler
func UserNewHandler(
	userUc usecase.UserUsecase,
	otpUC usecase.OtpUsecase,
	emailUC usecase.EmailUsecase,
	publisher *rabbitmq.Publisher,
	uploader *helpers.CloudinaryUploader,
) *UserHandler {
	return &UserHandler{
		userUc:    userUc,
		otpUC:     otpUC,
		emailUC:   emailUC,
		publisher: publisher,
		uploader:  uploader,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid form data"))
		return
	}

	// Unmarshal user data from "data" field
	var req dto.RegisterRequest
	jsonData := r.FormValue("data")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid JSON data"))
		return
	}

	// Handle image upload
	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		uploadOpts := &helpers.UploadOptions{Folder: "user_profiles"}
		uploadedImage, err := h.uploader.UploadImage(file, fileHeader, uploadOpts)
		if err != nil {
			helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to upload image"))
			return
		}
		// fmt.Printf(uploadedImage.URL)
		// Update profile image URL in request
		if req.Doctor != nil {
			req.Doctor.ProfileImageURL = uploadedImage.URL
		} else if req.Patient != nil {
			req.Patient.ProfileImageURL = uploadedImage.URL
		}
	} else if err != http.ErrMissingFile {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, "Invalid image file"))
		return
	}

	// Register the user
	user, err := h.userUc.Register(&req)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	if user == nil {
		helpers.Error(w, err)
		return
	}

	// Generate and save OTP
	_, err = h.otpUC.GenerateAndSaveOTP(user.Email, models.OTPPurposeRegister)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to send OTP"))
		return
	}

	// Respond to client immediately
	helpers.Success(w, http.StatusCreated, "User registered successfully. OTP sent to email.", user)
}

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	jwtClaims, err := jwt.GetUserDataFromReqJWT(r)
	if err != nil || jwtClaims == nil {
		helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	user, err := h.userUc.FindByID(jwtClaims.UserID)
	if err != nil {
		helpers.Error(w, err)
		return
	}

	helpers.Success(w, http.StatusOK, "User profile fetched successfully", user)
}
