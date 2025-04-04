package comment

import (
	"backend/internal"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var ErrEntryNotFound = errors.New("entry not found")

type Service struct {
	logger         *zap.Logger
	tracer         trace.Tracer
	commentQuerier Querier
}

func NewService(logger *zap.Logger, queries Querier) *Service {
	return &Service{
		logger:         logger,
		tracer:         otel.Tracer("comment/service"),
		commentQuerier: queries,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetAll")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comments, err := s.commentQuerier.FindAll(ctx)

	if err != nil {
		logger.Error("Error fetching all comments", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

func (s *Service) GetById(ctx context.Context, id pgtype.UUID) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetById")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.commentQuerier.FindByID(ctx, id)

	if err != nil {
		// Required entry not found
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", id.String()))
			return Comment{}, ErrEntryNotFound
		}
		// Other errors
		logger.Error("Error fetching comment by ID", zap.Error(err), zap.String("id", id.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Create(ctx context.Context, arg CreateParams) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "Create")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.commentQuerier.Create(ctx, arg)

	if err != nil {
		logger.Error("Error creating comment", zap.Error(err), zap.String("post_id", arg.PostID.String()), zap.String("author_id", arg.AuthorID.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Update(ctx context.Context, arg UpdateParams) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "Update")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.commentQuerier.Update(ctx, arg)
	if err != nil {
		// Required entry not found
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", arg.ID.String()))
			return Comment{}, ErrEntryNotFound
		}
		// Other errors
		logger.Error("Error updating comment", zap.Error(err), zap.String("id", arg.ID.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Delete(ctx context.Context, id pgtype.UUID) error {
	traceCtx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	err := s.commentQuerier.Delete(ctx, id)

	if err != nil {
		// Required entry not found
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", id.String()))
			return ErrEntryNotFound
		}
		// Other errors
		logger.Error("Error deleting comment", zap.Error(err), zap.String("id", id.String()))
		return err
	}
	return nil
}
