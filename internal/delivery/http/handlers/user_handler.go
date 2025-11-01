package handlers

import (
	"log"
	"net/http"

	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/pkg/utils/jwt"
	"hospital_management_system/internal/usecase"
)

// Handler handles user-related HTTP requests
type UserHandler struct {
	userUc   usecase.UserUsecase
	otpUC     usecase.OtpUsecase
	emailUC   usecase.EmailUsecase
	publisher *rabbitmq.Publisher
}

// NewHandler creates a new User Handler
func UserNewHandler(
	userUc   usecase.UserUsecase,
	otpUC     usecase.OtpUsecase,
	emailUC   usecase.EmailUsecase,
	publisher *rabbitmq.Publisher,
) *UserHandler {
	return &UserHandler{
		userUc:   userUc,
		 otpUC:     otpUC,
		emailUC:   emailUC,
		publisher: publisher,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	utils.BodyDecoder(w, r, &req)


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
	otpRecord, err := h.otpUC.GenerateAndSaveOTP(user.Email, models.OTPPurposeRegister)
	if err != nil {
		log.Println("Failed to generate/save OTP:", err)
		helpers.Error(w, helpers.NewAppError(http.StatusInternalServerError, "Failed to send OTP"))
		return
	}

	// Render email template
	body, err := utils.RenderEmailTemplate("templates/otp_email.html", map[string]string{
		"Name": user.Name,
		"Code": otpRecord.Code,
	})
	if err != nil {
		log.Println("Failed to render email template:", err)
	}

	// Asynchronous tasks: Create email record + Queue job
	go func() {
		job := helpers.EmailJob{
			To:      user.Email,
			Subject: "OTP Verification",
			Body:    body,
		}
		if err := h.publisher.Publish(job); err != nil {
			log.Println("Failed to publish email job:", err)
		}
	}()

	// Respond to client immediately
	helpers.Success(w, http.StatusCreated, "User registered successfully. OTP sent to email.", user)
}

// Login authenticates a user and returns a JWT
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	utils.BodyDecoder(w, r, &req)

	token, err := h.userUc.Login(&req)
	if err != nil {
		helpers.Error(w, helpers.NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	helpers.Success(w, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
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
