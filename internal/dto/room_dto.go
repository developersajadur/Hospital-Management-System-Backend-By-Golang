package dto

type CreateRoomRequest struct {
	RoomNumber   string  `json:"room_number" validate:"required"`
	Type         string  `json:"type" validate:"required,oneof=general icu vip"`
	PricePerDay  float64 `json:"price_per_day" validate:"required,gt=0"`
	Availability bool    `json:"availability"`
	Features     string  `json:"features,omitempty"`
	Image        *string  `json:"image,omitempty"`
}

type UpdateRoomRequest struct {
	RoomNumber   *string  `json:"room_number,omitempty"`
	Type         *string  `json:"type,omitempty" validate:"omitempty,oneof=general icu vip"`
	PricePerDay  *float64 `json:"price_per_day,omitempty" validate:"omitempty,gt=0"`
	Availability *bool    `json:"availability,omitempty"`
	Features     *string  `json:"features,omitempty"`
	Image        *string   `json:"image,omitempty"`
}
