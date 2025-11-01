package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoctorStatus string

const (
	DoctorActive   DoctorStatus = "active"
	DoctorInactive DoctorStatus = "inactive"
	DoctorOnLeave  DoctorStatus = "on_leave"
)

type Doctor struct {
    ID             uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
    UserID         uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
    Specialization string       `gorm:"type:varchar(255);not null" json:"specialization"`
    Experience     int          `gorm:"not null" json:"experience"`
    Fee            float64      `gorm:"type:decimal(10,2);not null" json:"fee"`
    ProfileImageId   string       `gorm:"type:uuid" json:"profile_image_id"`
    Status         DoctorStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
    CreatedAt      time.Time    `json:"created_at"`
    UpdatedAt      time.Time    `json:"updated_at"`

		
	// Relations
	Image Image `gorm:"foreignKey:ProfileImageId" json:"image,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}


// BeforeCreate hook: auto-generate UUID
func (d *Doctor) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	return nil
}

// BeforeUpdate hook: update timestamp
func (d *Doctor) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
