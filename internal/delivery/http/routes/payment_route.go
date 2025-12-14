package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	initPaymentRoute   = "/init"
	successPaymentRoute = "/success"
	failPaymentRoute    = "/fail"
	cancelPaymentRoute  = "/cancel"
	getAllPaymentsRoute = "/get-all"
)

func RegisterPaymentRoutes(r chi.Router, handler *handlers.PaymentHandler, userUC usecase.UserUsecase) {
	const prefix = "/payments"

	r.Route(prefix, func(r chi.Router) {

		// Protected routes (User must be logged in)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{
				models.RolePatient,
				models.RoleAdmin,
				models.RoleDoctor,
			}))
			r.Post(initPaymentRoute, handler.Init) // user initiates payment
		})
 		// admin gets all payments
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(userUC, []string{
				models.RoleAdmin,
			}))
			r.Get(getAllPaymentsRoute, handler.GetAll)
		})

		// SSLCommerz callback routes (public)
		r.Post(successPaymentRoute, handler.Success)
		r.Post(failPaymentRoute, handler.Fail)
		r.Post(cancelPaymentRoute, handler.Fail)
	})
}
