package routes

import (
	"hospital_management_system/internal/delivery/http/handlers"
	"hospital_management_system/internal/infra/middlewares"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const (
	uploadImage          = "/upload"
	uploadMultipleImages = "/upload-multiple"
	getImage             = "/{id}"
	getUserImages        = "/user/{user_id}"
	deleteImage          = "/delete/{id}"
)

func RegisterImageRoutes(r chi.Router, handler *handlers.ImageHandler, userUC usecase.UserUsecase) {
	const imageRoutePrefix = "/images"

	r.Route(imageRoutePrefix, func(r chi.Router) {

		// Protected routes
		r.Group(func(r chi.Router) {
			// Any authenticated user can access profile
			r.Use(middlewares.Auth(userUC, []string{models.RoleAdmin, models.RoleDoctor, models.RolePatient}))
			r.Post(uploadImage, handler.UploadImage)
			r.Post(uploadMultipleImages, handler.UploadMultipleImages)
			r.Get(getImage, handler.GetImage)
			r.Get(getUserImages, handler.GetUserImages)
			r.Delete(deleteImage, handler.DeleteImage)
		})
	})

}
