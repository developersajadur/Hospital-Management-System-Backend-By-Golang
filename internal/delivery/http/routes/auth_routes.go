package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	loginRoute    = "/login"
)

func RegisterAuthRoutes(r chi.Router, handler *handlers.AuthHandler, userUC usecase.UserUsecase) {
	const userRoutePrefix = "/auth"

	r.Route(userRoutePrefix, func(r chi.Router) {
		// Public routes
		r.Post(loginRoute, handler.Login)
	})
}
