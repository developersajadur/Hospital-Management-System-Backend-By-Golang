package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	registerRoute = "/register"
	registerPatientRoute = registerRoute + "/patient"
	registerAdminRoute = registerRoute + "/admin"
	registerDoctorRoute = registerRoute + "/doctor"
	profileRoute  = "/profile"
)

func RegisterUserRoutes(r chi.Router, handler *handlers.UserHandler, userUC usecase.UserUsecase) {
	const userRoutePrefix = "/users"

	r.Route(userRoutePrefix, func(r chi.Router) {
		// Public routes
		r.Post(registerPatientRoute, handler.Register) // patient registration and general

		// Protected routes
		r.Group(func(r chi.Router) {
			// Any authenticated user can access profile
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin, models.RoleDoctor, models.RolePatient}))
			r.Get(profileRoute, handler.GetProfile)
		})

		// Admin-only routes for creating admin/doctor accounts
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin}))
			r.Post(registerAdminRoute, handler.Register)
			r.Post(registerDoctorRoute, handler.Register)
		})
	})
}
