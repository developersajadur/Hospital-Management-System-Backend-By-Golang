package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OTPPurposeRegister      = "register"
	OTPPurposePasswordReset = "password_reset"
	OTPPurposeBookingVerify = "booking_verify"
)

type OTP struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"type:varchar(10);not null" json:"code"`
	Purpose   string    `gorm:"type:varchar(50);not null" json:"purpose"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook to auto-set UUID and timestamps
func (o *OTP) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to refresh updated timestamp
func (o *OTP) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}
