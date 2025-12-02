package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingType string

const (
	BookingTypeRoom    BookingType = "room"
	BookingTypeService BookingType = "service"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCompleted BookingStatus = "completed"
	BookingCanceled  BookingStatus = "canceled"
)

type Booking struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	BookingType BookingType   `gorm:"type:varchar(20);not null" json:"booking_type"`
	Status      BookingStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	PatientID uuid.UUID `gorm:"type:uuid;not null" json:"patient_id"`
	Patient   Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	// Patient   Patient   `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`

	RoomID       *uuid.UUID `gorm:"type:uuid" json:"room_id,omitempty"`
	Room         *Room      `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CheckInDate  *time.Time `json:"check_in_date,omitempty"`
	CheckOutDate *time.Time `json:"check_out_date,omitempty"`
	TotalPrice   *float64   `gorm:"type:decimal(10,2)" json:"total_price,omitempty"`

	ServiceID   *uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	Service     *Service   `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`

	SerialNumber *int `json:"serial_number,omitempty"`

	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *Booking) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
