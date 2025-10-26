package middlewares

import (
    "context"
    "net/http"
    "time"

    "hospital_management_system/internal/pkg/helpers"
    "hospital_management_system/internal/pkg/utils"
    "hospital_management_system/internal/services/user"
)

type contextKey string

const UserContextKey = contextKey("user")

func Auth(userUC user.Usecase, roles []string) func(http.Handler) http.Handler {
    allowedRoles := make(map[string]bool)
    for _, role := range roles {
        allowedRoles[role] = true
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Missing token"))
                return
            }

            claims, err := utils.VerifyJWT(token)
            if err != nil || claims == nil || time.Now().After(claims.ExpiresAt.Time) {
                helpers.Error(w, helpers.NewAppError(http.StatusUnauthorized, "Invalid or expired token"))
                return
            }

            u, err := userUC.GetUserById(claims.UserID)
            if err != nil || u == nil {
                helpers.Error(w, helpers.NewAppError(http.StatusNotFound, "User not found"))
                return
            }

            if u.IsDeleted {
                helpers.Error(w, helpers.NewAppError(http.StatusNotFound, "User deleted"))
                return
            }
            if u.IsBlocked {
                helpers.Error(w, helpers.NewAppError(http.StatusForbidden, "User is blocked"))
                return
            }
            if !u.IsVerified {
                helpers.Error(w, helpers.NewAppError(http.StatusForbidden, "User is not verified"))
                return
            }

            if !allowedRoles[u.Role] {
                helpers.Error(w, helpers.NewAppError(http.StatusForbidden, "Access denied"))
                return
            }

            ctx := context.WithValue(r.Context(), UserContextKey, u)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
