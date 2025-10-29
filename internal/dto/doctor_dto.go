package dto

import "hospital_management_system/internal/models"

type DoctorCreateRequest struct {
	UserID         string       `json:"user_id" binding:"required,uuid"`
	Specialization string       `json:"specialization" binding:"required"`
	Experience     int          `json:"experience" binding:"required"`
	Fee            float64      `json:"fee" binding:"required"`
	ProfileImage   string       `json:"profile_image"`
	Status         models.DoctorStatus `json:"status" binding:"required,oneof=active inactive on_leave"`
}

// DoctorUpdateRequest represents the payload to update an existing doctor

type DoctorUpdateRequest struct {
	Specialization string       `json:"specialization,omitempty"`
	Experience     int          `json:"experience,omitempty"`
	Fee            float64      `json:"fee,omitempty"`
	ProfileImage   string       `json:"profile_image,omitempty"`
	Status         models.DoctorStatus `json:"status,omitempty"`
}
