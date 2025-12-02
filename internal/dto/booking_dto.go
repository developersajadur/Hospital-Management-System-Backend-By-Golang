package dto

import "time"

type CreateBookingRequest struct {
	BookingType  string     `json:"booking_type" validate:"required,oneof=room service"`
	PatientID    string     `json:"patient_id" validate:"required"`

	RoomID       *string    `json:"room_id,omitempty"`
	CheckInDate  *time.Time `json:"check_in_date,omitempty"`
	CheckOutDate *time.Time `json:"check_out_date,omitempty"`

	ServiceID    *string    `json:"service_id,omitempty"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	TotalPrice   *float64   `json:"total_price" validate:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed completed canceled"`
}
