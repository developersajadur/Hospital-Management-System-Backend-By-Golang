// internal/models/image_model.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	URL          string         `gorm:"type:varchar(500);not null" json:"url"`
	PublicID     string         `gorm:"type:varchar(255);not null" json:"public_id"`
	FileName     string         `gorm:"type:varchar(255)" json:"file_name"`
	FileSize     int64          `gorm:"type:bigint" json:"file_size"`
	FileType     string         `gorm:"type:varchar(50)" json:"file_type"`
	Width        int            `gorm:"type:int" json:"width"`
	Height       int            `gorm:"type:int" json:"height"`
	ImageType    string         `gorm:"type:varchar(50);default:'general'" json:"image_type"` // profile, document, general
	IsDeleted    bool           `gorm:"default:false;index" json:"is_deleted"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook
func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// TableName specifies table name
func (Image) TableName() string {
	return "images"
}