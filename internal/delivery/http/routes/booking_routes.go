package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	bookingCreateRoute = "/create"
	getBookingListRoute   = "/get-all"
	getBookingByIDRoute   = "/get/{id}"
	changeBookingStatusRoute = "/{id}/status"
	deleteBookingDeleteRoute = "/delete/{id}"
)

func RegisterBookingRoutes(r chi.Router, handler *handlers.BookingHandler, userUC usecase.UserUsecase) {
	const bookingRoutePrefix = "/bookings"

	r.Route(bookingRoutePrefix, func(r chi.Router) {

		// Patient routes → Create a booking
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RolePatient}))
			r.Post(bookingCreateRoute, handler.Create)
		})

		// Admin + Doctor routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin, models.RoleDoctor}))

			r.Get(getBookingListRoute, handler.GetAll)
			r.Get(getBookingByIDRoute, handler.GetByID)
			r.Put(changeBookingStatusRoute, handler.UpdateStatus)
			r.Delete(deleteBookingDeleteRoute, handler.Delete)
		})
	})
}
