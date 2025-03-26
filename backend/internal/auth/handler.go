package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

type UserService interface {
	Register(ctx context.Context, u RegisterRequest) (string, error)
	Login(ctx context.Context, u LoginRequest) (string, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(service *Service) *Handler {
	// build interfaces from the service
	return &Handler{
		userService: service,
	}
}

type RegisterRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type Response struct {
	Token string
}

func (h Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into the RegisterRequest struct
	var u RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Call the Register function
	token, err := h.userService.Register(context.Background(), u)
	if err != nil {
		if err.Error() == "user_already_exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u LoginRequest

	// Decode the request body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Call the Login function
	token, err := h.userService.Login(context.Background(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
