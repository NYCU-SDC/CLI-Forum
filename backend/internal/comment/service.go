package comment

import (
	"backend/internal"
	"backend/internal/database"
	errorPkg "backend/internal/error"
	"context"
	"errors"
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
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		logger.Error("Error fetching all comments", zap.Error(err))
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
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		// Required entry not found
		if errors.Is(err, errorPkg.ErrNotFound) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", id.String()))
			return Comment{}, err
		}
		// Other errors
		logger.Error("Error fetching comment by ID", zap.Error(err), zap.String("id", id.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) GetByPost(ctx context.Context, postId pgtype.UUID) ([]Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByPost")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comments, err := s.query.FindByPostID(traceCtx, postId)

	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		logger.Error("Error fetching comments in post", zap.Error(err), zap.String("postId", postId.String()))
		return nil, err
	}

	return comments, nil
}

func (s *Service) Create(ctx context.Context, arg CreateParams) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "Create")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.query.Create(ctx, arg)

	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		logger.Error("Error creating comment", zap.Error(err), zap.String("post_id", arg.PostID.String()), zap.String("author_id", arg.AuthorID.String()))
		return Comment{}, err
	}
	return comment, nil
}

func (s *Service) Update(ctx context.Context, arg UpdateParams) (Comment, error) {
	traceCtx, span := s.tracer.Start(ctx, "Update")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	comment, err := s.query.Update(ctx, arg)
	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		// Required entry not found
		if errors.Is(err, errorPkg.ErrNotFound) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", arg.ID.String()))
			return Comment{}, err
		}
		// Other errors
		logger.Error("Error updating comment", zap.Error(err), zap.String("id", arg.ID.String()))
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
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		// Required entry not found
		if errors.Is(err, errorPkg.ErrNotFound) {
			logger.Error("Comment not found", zap.Error(err), zap.String("id", id.String()))
			return err
		}
		// Other errors
		logger.Error("Error deleting comment", zap.Error(err), zap.String("id", id.String()))
		return err
	}
	return nil
}
