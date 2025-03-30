package comment

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"net/http"
)

type Getter interface {
	GetAll(ctx context.Context) ([]Comment, error)
	GetById(ctx context.Context, id pgtype.UUID) (Comment, error)
}

type Store interface {
	Create(ctx context.Context, arg CreateParams) (Comment, error)
	Update(ctx context.Context, arg UpdateParams) (Comment, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type Handler struct {
	logger *zap.Logger
	getter Getter
	store  Store
}

func NewHandler(logger *zap.Logger, getter Getter, store Store) *Handler {
	return &Handler{
		logger: logger,
		getter: getter,
		store:  store,
	}
}

type GetByIdRequest struct {
	ID string `json:"id"`
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
	commentList, err := h.getter.GetAll(r.Context())

	// Handle error if fetching comment list fails
	if err != nil {
		h.logger.Error("Error fetching comment list", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Error encoding response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from request
	var req GetByIdRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Error decoding request body", zap.Error(err), zap.Any("body", r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify and transform ID
	var uuid pgtype.UUID
	err = uuid.Scan(req.ID)
	if err != nil {
		h.logger.Error("Error parsing UUID", zap.Error(err), zap.String("id", req.ID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch comment by ID
	comment, err := h.getter.GetById(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, ErrEntryNotFound) {
			h.logger.Error("Comment not found", zap.Error(err), zap.String("id", req.ID))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("Error fetching comment", zap.Error(err), zap.String("id", req.ID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Error encoding response", zap.Error(err), zap.Any("response", response))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
