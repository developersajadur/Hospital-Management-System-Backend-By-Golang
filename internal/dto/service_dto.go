package dto

type CreateServiceRequest struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description,omitempty"`
	Duration    int     `json:"duration" validate:"required,gt=0"` // minutes
}

type UpdateServiceRequest struct {
	Name        *string  `json:"name,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Description *string  `json:"description,omitempty"`
	Duration    *int     `json:"duration,omitempty"`
}
