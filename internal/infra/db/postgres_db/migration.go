package postgres_db

import (
	"hospital_management_system/internal/models"
	"log"
	"os"

	"gorm.io/gorm"
)

func Migration(DB *gorm.DB) {
	// Important: migrate users first, because doctors reference users
	err := DB.AutoMigrate(
		&models.User{},   // User table first
		&models.Doctor{}, // Doctor table second
		&models.Patient{},
		// &doctor.DoctorAvailability{},
		// &doctor.DoctorSlot{},
		&models.Room{},
		&models.Service{},
		// &booking.Booking{},
		// &payment.Payment{},
		&models.OTP{},
		&models.Email{},
		&models.Image{},
	)
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
		os.Exit(1)
	}

	log.Println("Database migrated successfully")
}
