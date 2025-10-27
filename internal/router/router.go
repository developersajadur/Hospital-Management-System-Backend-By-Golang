package router

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"hospital_management_system/internal/services/doctor"
	"hospital_management_system/internal/services/patient"
	"hospital_management_system/internal/services/user"
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
    userRepo := user.NewRepository(db)
    userUsecase := user.NewUsecase(userRepo, doctorUsecase, patientUsecase)
    userHandler := user.NewHandler(userUsecase)

    // Register routes
    user.RegisterRoutes(r, userHandler, userUsecase)
    // doctor.RegisterRoutes(r, doctorHandler, doctorUsecase)
}
