package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	sendOTPForVerifyAccount = "/send-verify-otp"
	validateOTP = "/validate-otp"
)

func RegisterOtpRoutes(r chi.Router, handler *handlers.OtpHandler, otpUc usecase.OtpUsecase) {
	const otpRoutePrefix = "/otps"

	r.Route(otpRoutePrefix, func(r chi.Router) {
		r.Post(sendOTPForVerifyAccount, handler.GenerateAndSaveOTP)
		r.Post(validateOTP, handler.ValidateOTP)

	})
}
