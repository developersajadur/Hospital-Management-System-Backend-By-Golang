package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Patient struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE" json:"user_id"`
	ProfileImageId *uuid.UUID `gorm:"type:uuid" json:"profile_image_id,omitempty"`
	Age            int        `gorm:"not null" json:"age"`
	Gender         Gender     `gorm:"type:varchar(20);not null" json:"gender"`
	Address        string     `gorm:"type:text;not null" json:"address"`
	MedicalHistory string     `gorm:"type:text" json:"medical_history,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Image *Image `gorm:"foreignKey:ProfileImageId" json:"image,omitempty"`
	User  User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook
func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// BeforeUpdate hook
func (p *Patient) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
