package user

import (
	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
)

type DoctorInfo struct {
	Specialization string  `json:"specialization"`
	Experience     int     `json:"experience"`
	Fee            float64 `json:"fee"`
	ProfileImage   string  `json:"profile_image"`
	Status         string  `json:"status"`
}

// RegisterRequest represents the data needed to register a user
type RegisterRequest struct {
	Name     string                  `json:"name"`
	Email    string                  `json:"email"`
	Phone    string                  `json:"phone"`
	Password string                  `json:"password"`
	Role     string                  `json:"role"` // doctor, patient, admin
	Doctor   *doctor.DoctorCreateRequest   `json:"doctor,omitempty"`
	Patient  *patient.PatientCreateRequest `json:"patient,omitempty"` // new field
}


// LoginRequest represents input for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

