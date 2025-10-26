package router

import (
	"hospital_management_system/internal/services/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Versioned API group
	api := r.Group("/api/v1")

	// Initialize domains and register their routes
	userRepo := user.NewRepository(db)
	userUsecase := user.NewUsecase(userRepo)
	userHandler := user.NewHandler(userUsecase)
	user.RegisterRoutes(api, userHandler)


	return r
}
