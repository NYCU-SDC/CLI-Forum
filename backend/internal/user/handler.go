package user

import (
	"backend/internal"
	"context"
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
	Logger *zap.Logger
	Store  Store
}

func NewHandler(logger *zap.Logger, store Store) *Handler {
	return &Handler{
		Logger: logger,
		Store:  store,
	}
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	logger := internal.LoggerWithContext(r.Context(), h.Logger)
	logger.Debug("CreateHandler called")

	panic("test panic")

	w.Write([]byte("CreateHandler"))
}
