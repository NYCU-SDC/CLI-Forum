package comment

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	commentQueries *Queries
}

func NewService(logger *zap.Logger, db *pgxpool.Pool) *Service {
	return &Service{
		logger:         logger,
		commentQueries: New(db),
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Comment, error) {
	comments, err := s.commentQueries.FindAll(ctx)

	if err != nil {
		s.logger.Error("Error fetching all comments", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

func (s *Service) GetById(ctx context.Context, id pgtype.UUID) (Comment, error) {
	comment, err := s.commentQueries.FindByID(ctx, id)

	if err != nil {
		s.logger.Error("Error fetching comment by ID", zap.Error(err), zap.String("id", id.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Create(ctx context.Context, arg CreateParams) (Comment, error) {
	comment, err := s.commentQueries.Create(ctx, arg)

	if err != nil {
		s.logger.Error("Error creating comment", zap.Error(err), zap.String("post_id", arg.PostID.String()), zap.String("author_id", arg.AuthorID.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Update(ctx context.Context, arg UpdateParams) (Comment, error) {
	comment, err := s.commentQueries.Update(ctx, arg)
	if err != nil {
		s.logger.Error("Error updating comment", zap.Error(err), zap.String("id", arg.ID.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Delete(ctx context.Context, id pgtype.UUID) error {
	err := s.commentQueries.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Error deleting comment", zap.Error(err), zap.String("id", id.String()))
		return err
	}
	return nil
}
