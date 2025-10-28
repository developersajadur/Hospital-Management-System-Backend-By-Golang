package router

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
	userAPI "hospital_management_system/internal/services/user/handler"
	userRepository "hospital_management_system/internal/services/user/repository"
	userUsecase "hospital_management_system/internal/services/user/usecase"
)

func SetupRoutes(r chi.Router, db *gorm.DB) {
    // Doctor domain
    doctorRepo := doctor.NewRepository(db)
    doctorUsecase := doctor.NewUsecase(doctorRepo)
    // doctorHandler := doctor.NewHandler(doctorUsecase)

       // Doctor domain
    patientRepo := patient.NewRepository(db)
    patientUsecase := patient.NewUsecase(patientRepo)
    // doctorHandler := doctor.NewHandler(doctorUsecase)

    // User domain, inject doctor usecase
    userRepo := userRepository.NewRepository(db)
    userUsecase := userUsecase.NewUsecase(userRepo, doctorUsecase, patientUsecase)
    userHandler := userAPI.NewHandler(userUsecase)

    // Register routes
    userAPI.RegisterRoutes(r, userHandler, userUsecase)
    // doctor.RegisterRoutes(r, doctorHandler, doctorUsecase)
}
