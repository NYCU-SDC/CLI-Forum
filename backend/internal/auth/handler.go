package auth

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type UserService interface {
	Register(ctx context.Context, u RegisterRequest) (string, error)
	Login(ctx context.Context, u LoginRequest) (string, error)
}

type Handler struct {
	logger      *zap.Logger
	userService UserService
}

func NewHandler(service *Service) *Handler {
	// build interfaces from the service
	return &Handler{
		logger:      service.logger,
		userService: service,
	}
}

type RegisterRequest struct {
	Username string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error when closing body", zap.Error(err))
		}
	}(r.Body)

	// check if the username or password is empty
	validate := validator.New()
	err = validate.Struct(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Call the Register function
	token, err := h.userService.Register(context.Background(), u)
	if err != nil {
		if err.Error() == "user_already_exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		h.logger.Error("error when registering user", zap.String("username", u.Username), zap.String("password", u.Password), zap.Error(err))
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error closing request body", zap.Error(err))
		}
	}(r.Body)

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
