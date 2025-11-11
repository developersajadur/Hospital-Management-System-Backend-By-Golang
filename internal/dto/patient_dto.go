package dto


type PatientCreateRequest struct {
	UserID         string `json:"user_id"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	MedicalHistory string `json:"medical_history,omitempty"`
}