package dto



type GenerateAndSaveOTPRequest struct {
	Email     string                 `json:"email"`
	Purpose    string                  `json:"purpose"`
}
