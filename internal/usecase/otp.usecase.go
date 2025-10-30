package usecase

import (
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Usecase defines OTP business logic
type OtpUsecase interface {
	GenerateAndSaveOTP(email string, purpose string) (*models.OTP, error)
	ValidateOTP(email string, code string) error
}

type otpUsecase struct {
	repo      repository.OtpRepository
	emailUc   EmailUsecase
	userUc    UserUsecase
	publisher *rabbitmq.Publisher
}

func OtpNewUsecase(repo repository.OtpRepository, emailUc EmailUsecase, userUc UserUsecase, publisher *rabbitmq.Publisher) OtpUsecase {
	return &otpUsecase{repo: repo, emailUc: emailUc, userUc: userUc, publisher: publisher}
}

// GenerateAndSaveOTP creates, saves, and returns a new OTP
func (u *otpUsecase) GenerateAndSaveOTP(email string, purpose string) (*models.OTP, error) {
	user, err := u.userUc.FindByEmail(email)
	if err != nil {
		return nil, helpers.NewAppError(http.StatusNotFound, "User not found")
	}
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

	// Render email template
	body, err := utils.RenderEmailTemplate("templates/otp_email.html", map[string]string{
		"Name": user.Name,
		"Code": otpCode,
	})
	if err != nil {
		log.Println("Failed to render email template:", err)
	}

	// Asynchronous tasks: Create email record and publish to RabbitMQ
	go func() {
		emailRecord, err := u.emailUc.CreateEmail(
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

		// Publish email job to RabbitMQ
		job := helpers.EmailJob{
			EmailID: emailRecord.ID,
			To:      emailRecord.Email,
			Subject: emailRecord.Subject,
			Body:    emailRecord.Body,
		}
		if err := u.publisher.Publish(job); err != nil {
			log.Println("Failed to publish email job:", err)
		}
	}()

	return otp, nil
}

// ValidateOTP checks validity, marks as used, and verifies user
func (u *otpUsecase) ValidateOTP(email string, code string) error {
	otp, err := u.repo.GetOTPByCodeAndEmail(email, code)
	if err != nil {
		return helpers.NewAppError(400, "Invalid OTP")
	}
	
	if time.Now().After(otp.ExpiresAt) {
		return helpers.NewAppError(400, "OTP expired")
	}

	// Use transaction to ensure both operations succeed or fail together
	err = u.repo.Transaction(func(tx *gorm.DB) error {
		// Mark OTP as used
		if err := u.repo.MarkOTPUsed(tx, otp.ID); err != nil {
			return err
		}

		// Mark user as verified
		if err := u.repo.MarkUserVerified(tx, email); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return helpers.NewAppError(500, "Failed to verify OTP")
	}

	return nil
}