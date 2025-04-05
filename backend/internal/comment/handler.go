package comment

import (
	"backend/internal"
	"backend/internal/problem"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type Getter interface {
	GetAll(ctx context.Context) ([]Comment, error)
	GetById(ctx context.Context, id pgtype.UUID) (Comment, error)
	GetByPost(ctx context.Context, postId pgtype.UUID) ([]Comment, error)
}

type Store interface {
	Create(ctx context.Context, arg CreateParams) (Comment, error)
	Update(ctx context.Context, arg UpdateParams) (Comment, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type Handler struct {
	logger    *zap.Logger
	tracer    trace.Tracer
	validator *validator.Validate
	getter    Getter
	store     Store
}

func NewHandler(logger *zap.Logger, getter Getter, store Store) *Handler {
	return &Handler{
		logger: logger,
		tracer: otel.Tracer("comment/handler"),
		getter: getter,
		store:  store,
	}
}

type GetByIdRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type GetByPostRequest struct {
	PostId string `json:"post_id" validate:"required,uuid"`
}

type Response struct {
	ID        string `json:"id"`
	PostId    string `json:"post_id"`
	AuthorId  string `json:"author_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) GetAllCommentHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetAllCommentEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	commentList, err := h.getter.GetAll(r.Context())

	// Handle error if fetching comment list fails
	if err != nil {
		logger.Error("Error fetching comment list", zap.Error(err))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert commentList to Response
	var response []Response
	for _, comment := range commentList {
		response = append(response, Response{
			ID:        comment.ID.String(),
			PostId:    comment.PostID.String(),
			AuthorId:  comment.AuthorID.String(),
			Title:     comment.Title.String,
			Content:   comment.Content.String,
			CreatedAt: comment.CreatedAt.Time.String(),
		})
	}

	// Write response
	internal.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetCommentByIdEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Extract ID from request
	var req GetByIdRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &req)

	if err != nil {
		logger.Error("Error decoding request body", zap.Error(err), zap.Any("body", r.Body))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Verify and transform ID
	var uuid pgtype.UUID
	err = uuid.Scan(req.ID)
	if err != nil {
		logger.Error("Error parsing UUID", zap.Error(err), zap.String("id", req.ID))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Fetch comment by ID
	comment, err := h.getter.GetById(r.Context(), uuid)
	if err != nil {
		logger.Error("Error fetching comment", zap.Error(err), zap.String("id", req.ID))
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

func (h *Handler) GetCommentByPostHandler(w http.ResponseWriter, r *http.Request) {
	traceCtx, span := h.tracer.Start(r.Context(), "GetCommentByPostEndpoint")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, h.logger)

	// Extract PostID from request
	var req GetByPostRequest
	err := internal.ParseAndValidateRequestBody(traceCtx, h.validator, r, &req)
	if err != nil {
		logger.Error("Error decoding request body", zap.Error(err), zap.Any("body", r.Body))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Verify and transform PostID
	var uuid pgtype.UUID
	err = uuid.Scan(req.PostId)
	if err != nil {
		logger.Error("Error parsing UUID", zap.Error(err), zap.String("post_id", req.PostId))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Fetch comments by PostID
	var response []Response
	comments, err := h.getter.GetByPost(r.Context(), uuid)
	if err != nil {
		logger.Error("Error fetching comments by post id", zap.Error(err), zap.String("post_id", req.PostId))
		problem.WriteError(traceCtx, w, err, logger)
		return
	}

	// Convert comments to Response
	for _, comment := range comments {
		response = append(response, Response{
			ID:        comment.ID.String(),
			PostId:    comment.PostID.String(),
			AuthorId:  comment.AuthorID.String(),
			Title:     comment.Title.String,
			Content:   comment.Content.String,
			CreatedAt: comment.CreatedAt.Time.String(),
		})
	}

	// Write response
	internal.WriteJSONResponse(w, http.StatusOK, response)
}
