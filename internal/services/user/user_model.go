package model

import (
	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RolePatient = "patient"
	RoleDoctor  = "doctor"
	RoleAdmin   = "admin"
)

type User struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string          `gorm:"not null" json:"name"`
	Email      string          `gorm:"uniqueIndex;not null" json:"email"`
	Phone      string          `gorm:"uniqueIndex;not null" json:"phone"`
	Password   string          `gorm:"not null" json:"-"`
	Role       string          `gorm:"not null" json:"role"`
	IsVerified bool            `gorm:"default:false" json:"is_verified"`
	IsBlocked  bool            `gorm:"default:false" json:"is_blocked"`
	IsDeleted  bool            `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`

	Doctor  *doctor.Doctor  `gorm:"foreignKey:UserID" json:"doctor,omitempty"`
	Patient *patient.Patient `gorm:"foreignKey:UserID" json:"patient,omitempty"`
}


// BeforeCreate hook: auto-generate UUID and timestamps
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	if u.Role == "" {
		u.Role = RolePatient
	}
	return nil
}

// BeforeUpdate hook: update timestamp
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
