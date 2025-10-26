// router/router.go
package router

import (
    "github.com/go-chi/chi/v5"
    "hospital_management_system/internal/services/user"
    "gorm.io/gorm"
)

func SetupRoutes(r chi.Router, db *gorm.DB) {
    // Initialize domains and register their routes
    userRepo := user.NewRepository(db)
    userUsecase := user.NewUsecase(userRepo)
    userHandler := user.NewHandler(userUsecase)
    user.RegisterRoutes(r, userHandler)
}
