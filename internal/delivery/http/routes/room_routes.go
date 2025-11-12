package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func RegisterRoomRoutes(r chi.Router, handler *handlers.RoomHandler, userUC usecase.UserUsecase) {
	const prefix = "/rooms"

	r.Route(prefix, func(r chi.Router) {
		// Admin-only protected routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin}))
			r.Post("/create", handler.Create)
			r.Patch("/update/{id}", handler.Update)
			r.Delete("/{id}", handler.Delete)
		})

		// Public routes
		r.Get("/", handler.GetRooms)
		r.Get("/{room_number}", handler.GetByRoomNumber)
	})
}
