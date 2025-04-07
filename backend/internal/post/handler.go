package post

import (
	"backend/internal"
	"backend/internal/jwt"
	"backend/internal/problem"
	"context"
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

	posts, err := h.postStore.GetAll(traceCtx)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}
	response := make([]Response, len(posts))
	for index, post := range posts {
		response[index] = GenerateResponse(post)
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetPostHandler")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	pathID := r.PathValue("id")
	postID, err := internal.ParseUUID(pathID)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	post, err := h.postStore.GetByID(traceCtx, postID)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	response := GenerateResponse(post)
	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "CreateEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	var request CreateRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	user, err := jwt.GetUserFromContext(traceCtx)
	if err != nil {
		logger.DPanic("Can't find user in context, this should never happen")
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	authorID, err := internal.ParseUUID(user.ID)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}
	request.AuthorID = authorID

	post, err := h.postStore.Create(traceCtx, request)
	if err != nil {
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	logger.Info("Created post", zap.String("id", post.ID.String()))

	response := GenerateResponse(post)
	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func GenerateResponse(post Post) Response {
	return Response{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Title:    post.Title.String,
		Content:  post.Content.String,
		CreateAt: post.CreateAt.Time.Format(time.RFC3339),
	}
}
