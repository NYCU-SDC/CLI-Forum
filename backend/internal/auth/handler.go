package auth

import (
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type LoginResponse struct {
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
	token, err := Register(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(RegisterResponse{Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u LoginRequest

	// Decode the request body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Call the Login function
	token, err := Login(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(RegisterResponse{Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
