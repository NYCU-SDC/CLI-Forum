package auth

import (
	"fmt"
	"net/http"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWT Middleware")
		// TODO: Implement JWT verification here
		next.ServeHTTP(w, r)
	})
}
