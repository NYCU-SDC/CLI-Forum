package comment

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
)

//go:generate mockery --name=Store
type Store interface {
	GetAll(ctx context.Context) ([]Comment, error)
	GetById(ctx context.Context, id uuid.UUID) (Comment, error)
	GetByPost(ctx context.Context, postId uuid.UUID) ([]Comment, error)
	Create(ctx context.Context, arg CreateRequest) (Comment, error)
}

type Handler struct {
	logger    *zap.Logger
	tracer    trace.Tracer
	validator *validator.Validate
	store     Store
}

func NewHandler(validator *validator.Validate, logger *zap.Logger, store Store) *Handler {
	return &Handler{
		logger:    logger,
		tracer:    otel.Tracer("comment/handler"),
		validator: validator,
		store:     store,
	}
}

type CreateRequest struct {
	PostID   uuid.UUID `json:"post_id"`
	AuthorID uuid.UUID `json:"author_id"`
	Title    string    `json:"title" validate:"required"`
	Content  string    `json:"content" validate:"required"`
}

type Response struct {
	ID        string `json:"id"`
	PostId    string `json:"post_id"`
	AuthorId  string `json:"author_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetAllCommentEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	commentList, err := h.store.GetAll(r.Context())

	// Handle error if fetching comment list fails
	if err != nil {
		logger.Error("Error fetching comment list", zap.Error(err))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert commentList to Response
	response := make([]Response, len(commentList))
	for i, comment := range commentList {
		response[i] = Response{
			ID:        comment.ID.String(),
			PostId:    comment.PostID.String(),
			AuthorId:  comment.AuthorID.String(),
			Title:     comment.Title.String,
			Content:   comment.Content.String,
			CreatedAt: comment.CreatedAt.Time.String(),
		}
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetCommentByIdEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	commentID := r.PathValue("id")

	// Verify and transform ID to UUID
	var id uuid.UUID
	err := id.Scan(commentID)
	if err != nil {
		logger.Error("Error parsing UUID", zap.Error(err), zap.String("id", commentID))
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrInvalidUUID, err), logger)
		return
	}

	comment, err := h.store.GetById(traceCtx, id)
	if err != nil {
		logger.Error("Error fetching comment", zap.Error(err), zap.String("id", commentID))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert comment to Response
	response := Response{
		ID:        comment.ID.String(),
		PostId:    comment.PostID.String(),
		AuthorId:  comment.AuthorID.String(),
		Title:     comment.Title.String,
		Content:   comment.Content.String,
		CreatedAt: comment.CreatedAt.Time.String(),
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetByPostHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetCommentByPostEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	postID := r.PathValue("post_id")

	// Verify and transform ID to UUID
	var id uuid.UUID
	err := id.Scan(postID)
	if err != nil {
		logger.Error("Error parsing UUID", zap.Error(err), zap.String("post_id", postID))
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrInvalidUUID, err), logger)
		return
	}

	comments, err := h.store.GetByPost(traceCtx, id)
	if err != nil {
		logger.Error("Error fetching comments by post id", zap.Error(err), zap.String("post_id", id.String()))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert comments to Response
	response := make([]Response, len(comments))
	for i, comment := range comments {
		response[i] = Response{
			ID:        comment.ID.String(),
			PostId:    comment.PostID.String(),
			AuthorId:  comment.AuthorID.String(),
			Title:     comment.Title.String,
			Content:   comment.Content.String,
			CreatedAt: comment.CreatedAt.Time.String(),
		}
	}

	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "CreateCommentEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Parse and validate requestBody body
	var req CreateRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &req)
	if err != nil {
		logger.Error("Error decoding requestBody body", zap.Error(err), zap.Any("body", r.Body))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	postID := r.PathValue("post_id")

	// Verify and transform PostID to UUID
	var id uuid.UUID
	err = id.Scan(postID)
	if err != nil {
		logger.Error("Error parsing UUID", zap.Error(err), zap.String("post_id", postID))
		problem.WriteError(traceCtx, w, fmt.Errorf("%w: %v", errorPkg.ErrInvalidUUID, err), logger)
		return
	}
	req.PostID = id

	// Convert AuthorId to UUID
	var authorId uuid.UUID
	u, err := jwt.GetUserFromContext(r.Context())
	if err != nil {
		logger.DPanic("Can't find user in context, this should never happen")
		problem.WriteError(traceCtx, w, err, logger)
	}
	err = authorId.Scan(u.ID)
	if err != nil {
		logger.Error("Error getting author id from context", zap.Error(err), zap.String("author_id", u.ID))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}
	req.AuthorID = authorId

	comment, err := h.store.Create(traceCtx, req)
	if err != nil {
		logger.Error("Error creating comment", zap.Error(err))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert comment to Response
	response := Response{
		ID:        comment.ID.String(),
		PostId:    comment.PostID.String(),
		AuthorId:  comment.AuthorID.String(),
		Title:     comment.Title.String,
		Content:   comment.Content.String,
		CreatedAt: comment.CreatedAt.Time.String(),
	}

	// Write response
	internal.WriteJSONResponse(w, http.StatusOK, response)
}
