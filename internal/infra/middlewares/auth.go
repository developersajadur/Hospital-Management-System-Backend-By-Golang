package middlewares

import (
	"context"
	"hospital_management_system/internal/pkg/helpers"
	"hospital_management_system/internal/pkg/utils/jwt"
	"net/http"
	"time"
)

type contextKey string

const UserContextKey = contextKey("user")

type User struct {
	ID         string
    Email string
	Role       string
	IsBlocked  bool
	IsVerified bool
	IsDeleted  bool
}

type UserFetcher interface {
	GetUserByIdForAuth(id string) (*User, error)
}

func Auth(userFetcher UserFetcher, roles []string) func(http.Handler) http.Handler {

	allowedRoles := make(map[string]bool)
	for _, r := range roles {
		allowedRoles[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			jwtUser, err := jwt.GetUserDataFromReqJWT(r)
			if err != nil {
				helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, ("Invalid token")))
				return
			}

			if time.Now().Unix() > jwtUser.Exp {
				helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, ("Token has expired")))
				return
			}

			user, err := userFetcher.GetUserByIdForAuth(jwtUser.UserID)
			if err != nil || user == nil {
				helpers.Error(w, helpers.NewAppError(http.StatusForbidden, ("User not found")))
				return

			}

			if user.IsBlocked {
				helpers.Error(w, helpers.NewAppError(http.StatusForbidden, ("User is blocked")))
				return
			}

			// if !user.IsVerified {

			// 	helpers.Error(w, helpers.NewAppError(http.StatusForbidden, ("User is not verified")))
			// 	return
			// }

			if user.IsDeleted {

				helpers.Error(w, helpers.NewAppError(http.StatusNotFound, ("User not found")))
				return
			}

			// Role check
			if !allowedRoles[user.Role] {

				helpers.Error(w, helpers.NewAppError(http.StatusForbidden, ("Unauthorized: insufficient role")))
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
