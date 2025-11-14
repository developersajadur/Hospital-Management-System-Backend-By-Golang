package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)


const (
	getAllServicesRoute = "/get-all"
	getServiceByIDRoute = "/get/{id}"
	CreateServiceRoute  = "/create"
	UpdateServiceRoute  = "/update/{id}"
	DeleteServiceRoute  = "/delete/{id}"
)

func RegisterServiceRoutes(r chi.Router, handler *handlers.ServiceHandler, userUC usecase.UserUsecase) {
	const serviceRoutePrefix = "/services"

	r.Route(serviceRoutePrefix, func(r chi.Router) {

		r.Get(getAllServicesRoute, handler.GetAll)
		r.Get(getServiceByIDRoute, handler.GetByID)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin}))
			r.Post(CreateServiceRoute, handler.Create)
			r.Patch(UpdateServiceRoute, handler.Update)
			r.Delete(DeleteServiceRoute, handler.Delete)
		})
	})
}
