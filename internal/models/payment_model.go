package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentInitiated PaymentStatus = "initiated"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCanceled  PaymentStatus = "canceled"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID     uuid.UUID     `gorm:"type:uuid;not null" json:"booking_id"`
	Booking       *Booking      `gorm:"foreignKey:BookingID" json:"-"`
	TranID        string        `gorm:"type:varchar(191);uniqueIndex;not null" json:"tran_id"`
	Amount        float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        PaymentStatus `gorm:"type:varchar(20);not null" json:"status"`
	Method        string        `gorm:"type:varchar(100)" json:"method,omitempty"`
	BankTranID    string        `gorm:"type:varchar(191)" json:"bank_tran_id,omitempty"`
	ValidationID  string        `gorm:"type:varchar(191)" json:"validation_id,omitempty"`
	TransactionAt *time.Time    `json:"transaction_at,omitempty"`
	IsDeleted     bool          `gorm:"default:false" json:"is_deleted"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
