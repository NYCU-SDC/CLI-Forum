package post

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	postQueries *Queries
	logger      *zap.SugaredLogger
}

func NewService(db *pgxpool.Pool, logger *zap.SugaredLogger) Service {
	return Service{
		postQueries: New(db),
		logger:      logger,
	}
}

func (s Service) GetAll(ctx context.Context) ([]Post, error) {
	posts, err := s.postQueries.FindAll(ctx)
	if err != nil {
		s.logger.Errorw("Error finding all posts", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func (s Service) GetPost(ctx context.Context, id pgtype.UUID) (Post, error) {
	post, err := s.postQueries.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Error finding post by ID", zap.Error(err))
		return Post{}, err
	}
	return post, err
}

func (s Service) CreatePost(ctx context.Context, r CreateRequest) (Post, error) {
	createdPost, err := s.postQueries.Create(ctx, CreateParams{
		Title:   pgtype.Text{String: r.Title},
		Content: pgtype.Text{String: r.Content},
	})
	if err != nil {
		s.logger.Errorw("Error creating post", zap.Error(err))
		return Post{}, err
	}
	return createdPost, nil
}
