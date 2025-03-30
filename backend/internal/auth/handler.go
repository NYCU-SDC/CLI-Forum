package auth

import (
	"backend/internal/user"
	"context"
	"go.uber.org/zap"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserStore interface {
	Create(ctx context.Context, name, password string) (user.User, error)
	GetByName(ctx context.Context, name string) (user.User, error)
}

type Handler struct {
	Logger    *zap.Logger
	UserStore UserStore
}

func NewHandler(logger *zap.Logger, userStore UserStore) *Handler {
	return &Handler{
		Logger:    logger,
		UserStore: userStore,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}
