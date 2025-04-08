package comment

import (
	"backend/internal"
	"backend/internal/database"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	tracer trace.Tracer
	query  *Queries
}

func NewService(logger *zap.Logger, db DBTX) *Service {
	return &Service{
		logger: logger,
		tracer: otel.Tracer("comment/service"),
		query:  New(db),
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetAll")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comments, err := s.query.FindAll(ctx)

	if err != nil {
		err = database.WrapDBError(err, logger, "get all comments")
		span.RecordError(err)
		return nil, err
	}
	return comments, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetById")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.query.FindByID(ctx, id)

	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "comments", "id", id.String(), logger, "get comment by id")
		span.RecordError(err)
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) GetByPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByPost")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comments, err := s.query.FindByPostID(traceCtx, postId)

	if err != nil {
		err = database.WrapDBError(err, logger, "get comments by post ID")
		span.RecordError(err)
		return nil, err
	}

	return comments, nil
}

func (s *Service) Create(ctx context.Context, arg CreateRequest) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "Create")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.query.Create(ctx, CreateParams{
		PostID:   arg.PostID,
		AuthorID: arg.AuthorID,
		Title:    pgtype.Text{String: arg.Title, Valid: true},
		Content:  pgtype.Text{String: arg.Content, Valid: true},
	})

	if err != nil {
		err = database.WrapDBError(err, logger, "create comment")
		span.RecordError(err)
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	traceCtx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	err := s.query.Delete(ctx, id)

	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "comments", "id", id.String(), logger, "delete comment")
		span.RecordError(err)
		return err
	}
	return nil
}
