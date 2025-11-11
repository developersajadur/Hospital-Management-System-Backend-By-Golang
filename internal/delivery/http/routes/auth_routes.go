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

		// // Protected routes
		// r.Group(func(r chi.Router) {
		// 	// Any authenticated user can access profile
		// 	r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin, models.RoleDoctor, models.RolePatient}))
		// 	r.Get(profileRoute, handler.GetProfile)
		// })

		// // Admin-only routes for creating admin/doctor accounts
		// r.Group(func(r chi.Router) {
		// 	r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin}))
		// 	r.Post(registerAdminRoute, handler.Register)
		// 	r.Post(registerDoctorRoute, handler.Register)
		// })
	})
}
