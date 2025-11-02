package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomType string

const (
	RoomTypeGeneral RoomType = "general"
	RoomTypeICU     RoomType = "icu"
	RoomTypeVIP     RoomType = "vip"
)

type Room struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RoomNumber   string     `gorm:"type:varchar(50);not null;unique" json:"room_number"`
	Type         RoomType   `gorm:"type:varchar(20);not null" json:"type"`
	PricePerDay  float64    `gorm:"type:decimal(10,2);not null" json:"price_per_day"`
	Availability bool       `gorm:"default:true" json:"availability"`
	Features     string     `gorm:"type:text" json:"features,omitempty"`
	IsDeleted    bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

func (r *Room) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}


// UUIDFromString safely converts a string to uuid.UUID.
// Returns uuid.Nil if invalid or empty.
func UUIDFromString(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}