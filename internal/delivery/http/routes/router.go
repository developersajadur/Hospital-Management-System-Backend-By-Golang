package routes

import (
	"log"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"hospital_management_system/config"
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/usecase"
)

func SetupRoutes(r chi.Router, db *gorm.DB) {
	// Initialize Doctor dependencies
	doctorRepo := repository.DoctorNewRepository(db)
	doctorUsecase := usecase.DoctorNewUsecase(doctorRepo)
	// doctorHandler := doctor.NewHandler(doctorUsecase)

	// Initialize Patient dependencies
	patientRepo := repository.PatientNewRepository(db)
	patientUsecase := usecase.PatientNewUsecase(patientRepo)
	// doctorHandler := doctor.NewHandler(doctorUsecase)

	// Initialize User dependencies
	userRepo := repository.UserNewRepository(db)
	userUsecase := usecase.UserNewUsecase(userRepo, doctorUsecase, patientUsecase)

	// Initialize Email dependencies
	emailRepo := repository.EmailNewRepository(db)
	emailUsecase := usecase.EmailNewUsecase(emailRepo)

	// Initialize OTP dependencies
	otpRepo := repository.OtpNewRepository(db)
	otpUsecase := usecase.OtpNewUsecase(otpRepo, emailUsecase, userUsecase)

	// Initialize RabbitMQ publisher dependencies
	publisher, err := rabbitmq.NewPublisher(config.ENV.RabbitMqUrl, "email_queue")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}


	userHandler := handlers.UserNewHandler(userUsecase, otpUsecase, emailUsecase, publisher)
	otpHandler := handlers.OtpNewHandler(otpUsecase)

	// Register routes
	RegisterUserRoutes(r, userHandler, userUsecase)
	RegisterOtpRoutes(r, otpHandler, otpUsecase)
	// doctor.RegisterRoutes(r, doctorHandler, doctorUsecase)

}
