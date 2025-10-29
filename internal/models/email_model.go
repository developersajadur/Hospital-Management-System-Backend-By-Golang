package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailType string
type EmailStatus string

const (
	EmailTypeOTP                 EmailType = "otp"
	EmailTypeBookingConfirmation EmailType = "booking_confirmation"
	EmailTypePasswordReset       EmailType = "password_reset"
	EmailTypeProfileUpdate       EmailType = "profile_update"
	EmailTypePaymentReceipt      EmailType = "payment_receipt"
	EmailTypeOther               EmailType = "other"

	EmailStatusPending EmailStatus = "pending"
	EmailStatusSent    EmailStatus = "sent"
	EmailStatusFailed  EmailStatus = "failed"
)

type Email struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Email     string         `gorm:"type:varchar(255);not null" json:"email"`
	Subject   string         `gorm:"type:varchar(255);not null" json:"subject"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	Type      EmailType      `gorm:"type:varchar(50);not null" json:"type"`
	Status    EmailStatus    `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	Error     *string        `gorm:"type:text" json:"error"`
	SentAt    *time.Time     `json:"sent_at"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate hook: auto-generate UUID and timestamps
func (u *Email) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate hook: update timestamp
func (u *Email) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
