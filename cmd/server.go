package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hospital_management_system/config"
	"hospital_management_system/internal/infra/db/postgres_db"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/router"

	"github.com/go-chi/chi/v5"
)

func RunServer() {
    // Connect and migrate DB
    postgres_db.ConnectDB()
    postgres_db.Migration(postgres_db.DB)

    // Setup Chi router
    r := chi.NewRouter()

    // Global middleware
    r.Use(middlewares.LoggingMiddleware)

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

    // Mount API v1 routes directly
    const apiV1Prefix = "/api/v1"
    r.Route(apiV1Prefix, func(api chi.Router) {
        router.SetupRoutes(api, postgres_db.DB)
    })

    // Create server
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
