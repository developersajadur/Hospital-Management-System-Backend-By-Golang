package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"hospital_management_system/config"
	"hospital_management_system/internal/delivery/http/routes"
	"hospital_management_system/internal/infra/cloudinary"
	"hospital_management_system/internal/infra/db/postgres_db"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/infra/rabbitmq"
	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/pkg/helpers"
		"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors" 
)

func RunServer() {
	// Connect and migrate DB
	postgres_db.ConnectDB()
	postgres_db.Migration(postgres_db.DB)

		// Initialize Cloudinary
	_, err := cloudinary.NewCloudinary()
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}


	// Initialize email repository
	emailRepo := repository.EmailNewRepository(postgres_db.DB)

	// Initialize RabbitMQ publisher
	emailPublisher, err := rabbitmq.NewPublisher(config.ENV.RabbitMqUrl, "email_queue")
	if err != nil {
		log.Fatalf("Failed to create email publisher: %v", err)
	}
	defer emailPublisher.Close()

	// Start RabbitMQ email consumer
	emailPort, _ := strconv.Atoi(config.ENV.EmailPort)
	go rabbitmq.StartConsumer(
		config.ENV.RabbitMqUrl,
		"email_queue",
		config.ENV.EmailHost,
		emailPort,
		config.ENV.Email,
		config.ENV.EmailAppPassword,
		emailRepo,
	)

	// Setup Chi router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middlewares.LoggingMiddleware)
		r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	rateLimiterConfig := middlewares.RateLimiterConfig{
		Limit:  5,
		Period: 1 * time.Second,
	}
	r.Use(func(next http.Handler) http.Handler {
		return middlewares.RateLimiter(next, rateLimiterConfig)
	})

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helpers.Success(w, http.StatusOK, "Welcome To Hospital Management Server", nil)
	})

	// Mount API v1 routes
	const apiV1Prefix = "/api/v1"
	r.Route(apiV1Prefix, func(api chi.Router) {
		routes.SetupRoutes(api, postgres_db.DB)
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.ENV.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("Server running at port %s\n", config.ENV.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
