package usecase

import (
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Usecase defines OTP business logic
type OtpUsecase interface {
	GenerateAndSaveOTP(email string, purpose string) (*models.OTP, error)
	ValidateOTP(userID uuid.UUID, code string, purpose string) error
}

type otpUsecase struct {
	repo    repository.OtpRepository
	emailUc EmailUsecase
	userUc  UserUsecase
}

func OtpNewUsecase(repo repository.OtpRepository, emailUc EmailUsecase, userUc UserUsecase) OtpUsecase {
	return &otpUsecase{repo: repo, emailUc: emailUc, userUc: userUc}
}

// GenerateAndSaveOTP creates, saves, and returns a new OTP
func (u *otpUsecase) GenerateAndSaveOTP(email string, purpose string) (*models.OTP, error) {
	otpCode := utils.GenerateOTP()
	expiration := time.Now().Add(5 * time.Minute)

	otp := &models.OTP{
		Email:     email,
		Code:      otpCode,
		Purpose:   purpose,
		ExpiresAt: expiration,
	}

	if err := u.repo.SaveOTP(otp); err != nil {
		return nil, helpers.NewAppError(500, "Failed to save OTP")
	}

	user, err := u.userUc.FindByEmail(email)
	if err != nil {
		return nil, helpers.NewAppError(http.StatusNotFound, "User not found")
	}
	// Render email template
	body, err := utils.RenderEmailTemplate("templates/otp_email.html", map[string]string{
		"Name": user.Name,
		"Code": otpCode,
	})
	if err != nil {
		log.Println("Failed to render email template:", err)
	}

	// Asynchronous tasks: Create email record
	go func() {
		_, err := u.emailUc.CreateEmail(
			user.ID,
			user.Email,
			"OTP Verification",
			body,
			models.EmailTypeOTP,
		)
		if err != nil {
			log.Println("Failed to create email record:", err)
			return
		}
	}()

	return otp, nil
}

// ValidateOTP checks validity and marks as used
func (u *otpUsecase) ValidateOTP(userID uuid.UUID, code string, purpose string) error {
	otp, err := u.repo.GetOTPByCode(userID, code, purpose)
	if err != nil {
		return helpers.NewAppError(400, "Invalid OTP")
	}

	if time.Now().After(otp.ExpiresAt) {
		return helpers.NewAppError(400, "OTP expired")
	}

	if err := u.repo.MarkOTPUsed(otp.ID); err != nil {
		return helpers.NewAppError(500, "Failed to mark OTP as used")
	}

	return nil
}
