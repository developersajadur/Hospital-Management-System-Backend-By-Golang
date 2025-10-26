package middlewares

import "net/http"



type Middleware func(http.Handler) http.Handler

func Middlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
