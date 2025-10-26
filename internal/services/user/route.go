package user

import "github.com/go-chi/chi/v5"


const (
	registerRoute = "/register"
	loginRoute    = "/auth/login"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	const userRoutePrifix = "/users"
	r.Route(userRoutePrifix, func(r chi.Router) {
		r.Post(registerRoute, handler.Register)
		r.Post(loginRoute, handler.Login)
	})
}
