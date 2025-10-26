package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string         `gorm:"not null" json:"name"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Phone      string         `gorm:"uniqueIndex;not null" json:"phone"`
	Password   string         `gorm:"not null" json:"-"`
	Role       string         `gorm:"not null" json:"role"` // patient, doctor, admin
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	IsBlocked  bool           `gorm:"default:false" json:"is_blocked"`
	IsDeleted  bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// BeforeCreate hook: auto-generate UUID and set defaults
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = RolePatient // default role if not set
	}

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	// Ensure booleans default to false
	u.IsVerified = false
	u.IsBlocked = false
	u.IsDeleted = false

	return
}

// BeforeUpdate hook: update timestamp
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

// Roles constants
const (
	RolePatient = "patient"
	RoleDoctor  = "doctor"
	RoleAdmin   = "admin"
)
