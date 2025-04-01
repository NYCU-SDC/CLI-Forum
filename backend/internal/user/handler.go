package user

import (
	"backend/internal"
	"backend/internal/jwt"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type Store interface {
	Create(ctx context.Context, name, password string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	UpdateName(ctx context.Context, id uuid.UUID, name string) (User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	Validator *validator.Validate
	Logger    *zap.Logger
	Store     Store
}

func NewHandler(validator *validator.Validate, logger *zap.Logger, store Store) *Handler {
	return &Handler{
		Validator: validator,
		Logger:    logger,
		Store:     store,
	}
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	logger := internal.LoggerWithContext(r.Context(), h.Logger)
	logger.Debug("CreateHandler called")

	user := r.Context().Value(internal.UserContextKey).(jwt.User)
	fmt.Printf("Is you %s !", user.Username)
}
