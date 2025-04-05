package post

import (
	"backend/internal"
	"backend/internal/database"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//go:generate mockery --name Querier
type Querier interface {
	FindAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id uuid.UUID) (Post, error)
	Create(ctx context.Context, arg CreateParams) (Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, arg UpdateParams) (Post, error)
}

type Service struct {
	logger *zap.Logger
	tracer trace.Tracer
	query  Querier
}

func NewService(logger *zap.Logger, db *pgxpool.Pool) Service {
	return Service{
		logger: logger,
		tracer: otel.Tracer("post/service"),
		query:  New(db),
	}
}

func (s Service) GetAll(ctx context.Context) ([]Post, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	posts, err := s.query.FindAll(traceCtx)
	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)
		logger.Error("Failed to get all posts", zap.Error(err))
		return nil, err
	}

	return posts, nil
}

func (s Service) GetByID(ctx context.Context, id uuid.UUID) (Post, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	post, err := s.query.FindByID(traceCtx, id)
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "id", id.String(), logger)
		span.RecordError(err)
		logger.Error("Failed to get post by ID", zap.Error(err))
		return Post{}, err
	}
	return post, err
}

func (s Service) Create(ctx context.Context, r CreateRequest) (Post, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	createdPost, err := s.query.Create(traceCtx, CreateParams{
		AuthorID: r.AuthorID,
		Title:    pgtype.Text{String: r.Title, Valid: true},
		Content:  pgtype.Text{String: r.Content, Valid: true},
	})
	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)
		logger.Error("Failed to create post", zap.Error(err))
		return Post{}, err
	}
	return createdPost, nil
}
