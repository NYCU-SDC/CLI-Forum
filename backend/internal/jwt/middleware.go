package jwt

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Verifier interface {
	Verify(tokenString string) error
}

type MiddlewareService struct {
	logger   *zap.Logger
	verifier Verifier
}

func NewMiddleware(service Service) MiddlewareService {
	return MiddlewareService{
		logger:   service.logger,
		verifier: service,
	}
}

func (m MiddlewareService) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWT Middleware")

		// Get the token from the Authorization header
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(token, "Bearer ")

		// Verify the token
		err := m.verifier.Verify(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Call the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}
