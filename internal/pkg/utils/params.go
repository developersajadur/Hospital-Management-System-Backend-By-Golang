package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Param retrieves URL parameters using chi router.
func Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
