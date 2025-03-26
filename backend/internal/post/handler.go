package post

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockery --name Finder
type Finder interface {
	GetAll(ctx context.Context) ([]Post, error)
	GetPost(ctx context.Context, id pgtype.UUID) (Post, error)
	CreatePost(ctx context.Context, request CreateRequest) (Post, error)
}

type Handler struct {
	Finder Finder
	logger *zap.SugaredLogger
}

func NewHandler(f Finder, logger *zap.SugaredLogger) Handler {
	return Handler{
		Finder: f,
		logger: logger,
	}
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get all posts from the service
	posts, err := h.Finder.GetAll(r.Context())
	if err != nil {
		h.logger.Errorw("Error getting all posts", "error", err)
		return
	}

	// Create the response
	var response []Response
	for _, post := range posts {
		response = append(response, Response{
			ID:       post.ID.String(),
			AuthorID: post.AuthorID.String(),
			Title:    post.Title.String,
			Content:  post.Content.String,
			CreateAt: post.CreateAt.Time.String(),
		})
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Errorw("error encoding response", "error", err)
		return
	}
}

func (h Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	// Get the post id from the query string
	postID := r.URL.Query().Get("id")
	if postID == "" {
		h.logger.Error("missing post id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Scan the post id into a pgtype.UUID
	var id pgtype.UUID
	err := id.Scan(postID)
	if err != nil {
		h.logger.Errorw("error scanning post id", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the post from the service
	post, err := h.Finder.GetPost(r.Context(), id)
	if err != nil {
		h.logger.Errorw("Error getting post", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the response
	response := Response{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Title:    post.Title.String,
		Content:  post.Content.String,
		CreateAt: post.CreateAt.Time.String(),
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Errorw("Error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	var createRequest CreateRequest
	err := decoder.Decode(&createRequest)
	if err != nil {
		h.logger.Errorw("Error decoding body", "error", err)
		return
	}
	defer r.Body.Close()

	// Create the post
	post, err := h.Finder.CreatePost(r.Context(), createRequest)
	if err != nil {
		h.logger.Errorw("Error creating post", "error", err)
		return
	}

	h.logger.Infow("Created post", "id", post.ID.String())

	// Create the response
	response := Response{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Title:    post.Title.String,
		Content:  post.Content.String,
		CreateAt: post.CreateAt.Time.String(),
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Errorw("Error encoding response", "error", err)
		return
	}
}
