package post

import (
	"backend/internal"
	errorPkg "backend/internal/error"
	"backend/internal/jwt"
	"backend/internal/problem"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateRequest struct {
	AuthorID uuid.UUID `json:"author_id"`
	Title    string    `json:"title"   validate:"required"`
	Content  string    `json:"content" validate:"required"`
}

type Response struct {
	ID       string `json:"id"`
	AuthorID string `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
}

//go:generate mockery --name Store
type Store interface {
	GetAll(ctx context.Context) ([]Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (Post, error)
	Create(ctx context.Context, request CreateRequest) (Post, error)
}

type Handler struct {
	validator *validator.Validate
	logger    *zap.Logger
	tracer    trace.Tracer

	postStore Store
}

func NewHandler(v *validator.Validate, logger *zap.Logger, s Store) Handler {
	return Handler{
		tracer:    otel.Tracer("post/handler"),
		validator: v,
		logger:    logger,
		postStore: s,
	}
}

func (h Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetAllEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Get all posts from the service
	posts, err := h.postStore.GetAll(r.Context())
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Write the response
	var response []Response
	for _, post := range posts {
		response = append(response, Response{
			ID:       post.ID.String(),
			AuthorID: post.AuthorID.String(),
			Title:    post.Title.String,
			Content:  post.Content.String,
			CreateAt: post.CreateAt.Time.Format(time.RFC3339),
		})
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetPostHandler")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Get the post id from the path
	postID := r.PathValue("id")

	// Scan the post id into uuid.UUID
	var id uuid.UUID
	err := id.Scan(postID)
	if err != nil {
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrInvalidUUID, err), logger)
		return
	}

	// Get the post from the service
	post, err := h.postStore.GetByID(r.Context(), id)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Write the response
	response := Response{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Title:    post.Title.String,
		Content:  post.Content.String,
		CreateAt: post.CreateAt.Time.Format(time.RFC3339),
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "CreateEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Parse and validate the request body
	var request CreateRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Get the user from the context
	user := r.Context().Value(internal.UserContextKey).(jwt.User)
	var authorID uuid.UUID
	err = authorID.Scan(user.ID)
	if err != nil {
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrInvalidUUID, err), logger)
		return
	}
	request.AuthorID = authorID

	// Create the post
	post, err := h.postStore.Create(r.Context(), request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	logger.Info("Created post", zap.String("id", post.ID.String()))

	// Write the response
	response := Response{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Title:    post.Title.String,
		Content:  post.Content.String,
		CreateAt: post.CreateAt.Time.Format(time.RFC3339),
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}
