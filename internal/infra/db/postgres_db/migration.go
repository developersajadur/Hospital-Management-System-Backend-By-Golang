package postgres_db

import (
	"log"
	"os"

	"hospital_management_system/internal/services/user"

	"gorm.io/gorm"
)

func Migration(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&user.User{},
		// &patient.Patient{},
		// &doctor.Doctor{},
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
