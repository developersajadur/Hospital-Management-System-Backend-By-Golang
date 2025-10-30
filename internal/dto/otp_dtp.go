package dto

type GenerateAndSaveOTPRequest struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
}

type OtpValidateRequest struct {
	Email   string `json:"email"`
	Code string `json:"code"`
}
