package auth

import (
	"backend/internal/user"
	"go.uber.org/zap"
	"net/http"
)

func Middleware(next http.HandlerFunc, logger *zap.Logger, requiredRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("user") == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u := r.Context().Value("user").(user.User)
		for _, role := range requiredRoles {
			// Always true when user logged in
			if u.HasRole(role) {
				break
			}
		}

		next(w, r)
	}
}
