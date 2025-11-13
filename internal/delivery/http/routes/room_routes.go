package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)


const (
	createRoomRoute    = "/create"
	updateRoomRoute    = "/update/{id}"
	deleteRoomRoute    = "/delete/{id}"
	getRoomByNumberRoute = "/get/{room_number}"
	getAllRoomsRoute   = "/get-all"
)

func RegisterRoomRoutes(r chi.Router, handler *handlers.RoomHandler, userUC usecase.UserUsecase) {
	const prefix = "/rooms"

	r.Route(prefix, func(r chi.Router) {
		// Admin-only protected routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin}))
			r.Post(createRoomRoute, handler.Create)
			r.Patch(updateRoomRoute, handler.Update)
			r.Delete(deleteRoomRoute, handler.Delete)
		})

		// Public routes
		r.Get(getAllRoomsRoute, handler.GetRooms)
		r.Get(getRoomByNumberRoute, handler.GetByRoomNumber)
	})
}
