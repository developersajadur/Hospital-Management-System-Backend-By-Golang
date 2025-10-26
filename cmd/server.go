package server

import (
	"fmt"
	"log"
	"time"

	"hospital_management_system/config"
	"hospital_management_system/internal/infra/db/postgres_db"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/router"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Connect to PostgreSQL
	postgres_db.ConnectDB()

	// Migrate database
	postgres_db.Migration(postgres_db.DB)

	// Setup router
	r := router.SetupRoutes(postgres_db.DB)

	// Apply global middleware (Rate Limiter example)
	r.Use(middlewares.NewRateLimiterMiddleware(middlewares.RateLimiterConfig{
		Limit:  5,
		Period: 1 * time.Second,
	}))

	fmt.Printf("Server running at port %s\n", config.ENV.Port)

	// Run server
	if err := r.Run(":" + config.ENV.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
