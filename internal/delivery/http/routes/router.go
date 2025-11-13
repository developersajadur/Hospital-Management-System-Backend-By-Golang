package routes

import (
	"log"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"hospital_management_system/config"
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/usecase"
)

func SetupRoutes(r chi.Router, db *gorm.DB, cloudinaryUploader *helpers.CloudinaryUploader) {
	// Initialize RabbitMQ publisher dependencies
	publisher, err := rabbitmq.NewPublisher(config.ENV.RabbitMqUrl, "email_queue")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}

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


	// Initialize Auth dependencies

	authUsecase := usecase.AuthNewUsecase(userRepo)
	authHandler := handlers.AuthNewHandler(authUsecase)

	// Initialize Email dependencies
	emailRepo := repository.EmailNewRepository(db)
	emailUsecase := usecase.EmailNewUsecase(emailRepo)

	// Initialize OTP dependencies
	otpRepo := repository.OtpNewRepository(db)
	otpUsecase := usecase.OtpNewUsecase(otpRepo, emailUsecase, userUsecase, publisher)

	userHandler := handlers.UserNewHandler(userUsecase, otpUsecase, emailUsecase, publisher, cloudinaryUploader)
	otpHandler := handlers.OtpNewHandler(otpUsecase)

	// Initialize Image dependencies
	imageRepo := repository.ImageNewRepository(db)
	imageUsecase := usecase.ImageNewUsecase(imageRepo, cloudinaryUploader)
	imageHandler := handlers.ImageNewHandler(imageUsecase)

	// Initialize Room dependencies
	roomRepo := repository.RoomNewRepository(db)
	roomUsecase := usecase.RoomNewUsecase(roomRepo)
	roomHandler := handlers.RoomNewHandler(roomUsecase, cloudinaryUploader)


	// Initialize Service dependencies
	serviceRepo := repository.ServiceNewRepository(db)
	serviceUsecase := usecase.ServiceNewUsecase(serviceRepo)
	serviceHandler := handlers.ServiceNewHandler(serviceUsecase)

	// Register routes
	RegisterUserRoutes(r, userHandler, userUsecase)
	RegisterOtpRoutes(r, otpHandler, otpUsecase)
	RegisterImageRoutes(r, imageHandler, userUsecase)
	RegisterRoomRoutes(r, roomHandler, userUsecase)
	RegisterAuthRoutes(r, authHandler, userUsecase)
	RegisterServiceRoutes(r, serviceHandler, userUsecase)
	// doctor.RegisterRoutes(r, doctorHandler, doctorUsecase)

}
