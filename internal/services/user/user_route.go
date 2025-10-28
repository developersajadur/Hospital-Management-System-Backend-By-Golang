
package handler

import (
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/services/user/model"
	"hospital_management_system/internal/services/user/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	registerRoute = "/register"
	loginRoute    = "/auth/login"
	profileRoute  = "/profile"
)

func RegisterRoutes(r chi.Router, handler *Handler, userUC usecase.Usecase) {
	const userRoutePrefix = "/users"

	r.Route(userRoutePrefix, func(r chi.Router) {
		// Public routes
		r.Post(registerRoute, handler.Register) // patient registration and general
		r.Post(loginRoute, handler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			// Any authenticated user can access profile
			r.Use(middlewares.Auth(userUC, []string{model.RoleAdmin, model.RoleDoctor, model.RolePatient}))
			r.Get(profileRoute, handler.GetProfile)
		})

		// Admin-only routes for creating admin/doctor accounts
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{model.RoleAdmin}))
			r.Post("/register/admin", handler.Register)
			r.Post("/register/doctor", handler.Register)
		})
	})
}
