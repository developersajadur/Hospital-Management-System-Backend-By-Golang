package postgres_db

import (
	"log"
	"os"

	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
	userModel "hospital_management_system/internal/services/user/model"

	// "hospital_management_system/internal/services/patient"
	// other services can be imported here

	"gorm.io/gorm"
)

func Migration(DB *gorm.DB) {
	// Important: migrate users first, because doctors reference users
	err := DB.AutoMigrate(
		&userModel.User{},   // User table first
		&doctor.Doctor{}, // Doctor table second
		&patient.Patient{},
		// &doctor.DoctorAvailability{},
		// &doctor.DoctorSlot{},
		// &room.Room{},
		// &service_entity.Service{},
		// &booking.Booking{},
		// &payment.Payment{},
		// &otp.OTP{},
		// &email.Email{},
	)
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
		os.Exit(1)
	}

	log.Println("Database migrated successfully")
}
