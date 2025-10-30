package handlers

import (
	"hospital_management_system/internal/dto"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils"
	"hospital_management_system/internal/usecase"
	"net/http"
)

type OtpHandler struct {
	otpUc usecase.OtpUsecase
}

func OtpNewHandler(uc usecase.OtpUsecase) *OtpHandler {
	return &OtpHandler{otpUc: uc}
}


func (h *OtpHandler) GenerateAndSaveOTP (w http.ResponseWriter, r *http.Request) {
	var req dto.GenerateAndSaveOTPRequest
	utils.BodyDecoder(w, r, &req)

	_, err := h.otpUc.GenerateAndSaveOTP(req.Email, req.Purpose)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusCreated, "Otp Send successfully", nil)
}

func (h *OtpHandler) ValidateOTP (w http.ResponseWriter, r *http.Request) {
	var req dto.OtpValidateRequest
		utils.BodyDecoder(w, r, &req)
	 err := h.otpUc.ValidateOTP(req.Email, req.Code)
	if err != nil {
		helpers.Error(w, err)
		return
	}
	helpers.Success(w, http.StatusCreated, "Otp Validate Successfully", nil)

}