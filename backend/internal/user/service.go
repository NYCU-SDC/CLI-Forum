package user

import (
	"backend/internal"
	"backend/internal/database"
	errorPkg "backend/internal/error"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// HasRole is a placeholder for the actual role checking logic, because the role storage is not in project scope.
func (u User) HasRole(role string) bool {
	return true
}

type Service struct {
	logger *zap.Logger
	tracer trace.Tracer
	query  *Queries
}

func NewService(logger *zap.Logger, db DBTX) *Service {
	return &Service{
		logger: logger,
		tracer: otel.Tracer("user/service"),
		query:  New(db),
	}
}

func (s *Service) Create(ctx context.Context, name, password string) (User, error) {
	traceCtx, span := s.tracer.Start(ctx, "Create")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	user, err := s.query.Create(ctx, CreateParams{
		Name:     name,
		Password: password,
	})
	if err != nil {
		err = database.WrapDBError(err, logger)
		span.RecordError(err)

		if errors.Is(err, database.ErrUniqueViolation) {
			logger.Debug("User already exists", zap.String("name", name))
			return User{}, errorPkg.ErrUserAlreadyExists
		}

		logger.Error("Failed to create user", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	user, err := s.query.GetByID(ctx, id)
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "id", id.String(), logger)
		span.RecordError(err)
		logger.Error("Failed to get user by ID", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (User, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetByName")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	user, err := s.query.GetByName(ctx, name)
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "name", name, logger)
		span.RecordError(err)
		logger.Error("Failed to get user by name", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (s *Service) UpdateName(ctx context.Context, id uuid.UUID, name string) (User, error) {
	traceCtx, span := s.tracer.Start(ctx, "UpdateName")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	user, err := s.query.UpdateName(ctx, UpdateNameParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "id", id.String(), logger)
		span.RecordError(err)
		logger.Error("Failed to update user name", zap.Error(err))
		return User{}, err
	}

	logger.Debug("Updated user", zap.String("id", id.String()), zap.String("name", name))

	return user, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	traceCtx, span := s.tracer.Start(ctx, "UpdatePassword")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	count, err := s.query.UpdatePassword(ctx, UpdatePasswordParams{
		ID:       id,
		Password: password,
	})
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "id", id.String(), logger)
		span.RecordError(err)
		logger.Error("Failed to update user password", zap.Error(err))
		return err
	}

	logger.Debug("Updated user password", zap.String("id", id.String()), zap.Int64("affected_rows", count))

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	traceCtx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()
	logger := internal.LoggerWithContext(traceCtx, s.logger)

	count, err := s.query.Delete(ctx, id)
	if err != nil {
		err = database.WrapDBErrorWithKeyValue(err, "users", "id", id.String(), logger)
		span.RecordError(err)
		logger.Error("Failed to delete user", zap.Error(err))
		return err
	}

	logger.Debug("Deleted user", zap.String("id", id.String()), zap.Int64("affected_rows", count))

	return nil
}
