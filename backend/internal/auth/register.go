package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into the RegisterRequest struct
	var u RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Call the Register function
	err = Register(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Registration successful, generate a token
	expirationTime := time.Now().Add(120 * time.Hour)
	claims := &Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	json.NewEncoder(w).Encode(RegisterResponse{Token: tokenString})
}

func Register(u RegisterRequest) error {
	// TODO: Implement registration with the database
	return nil
}
