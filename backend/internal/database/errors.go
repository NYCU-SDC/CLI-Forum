package database

import (
	errorPkg "backend/internal/error"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	PGErrUniqueViolation     = "23505"
	PGErrForeignKeyViolation = "23503"
	PGErrDeadlockDetected    = "40P01"
)

var (
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDeadlockDetected    = errors.New("deadlock detected")
	ErrQueryTimeout        = errors.New("query timed out")
)

type InternalServerError struct {
	Source error
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Source.Error())
}

func WrapDBError(err error, logger *zap.Logger, message string) error {
	if err == nil {
		return nil
	}

	logger.Warn("Wrapping database error", zap.Error(err))

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %v", errorPkg.ErrNotFound, err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrQueryTimeout
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGErrUniqueViolation:
			logger.Warn("Unique constraint violation", zap.String("table", pgErr.TableName), zap.String("constraint", pgErr.ConstraintName))
			return fmt.Errorf("%w: %v", ErrUniqueViolation, err)
		case PGErrForeignKeyViolation:
			logger.Warn("Foreign key violation", zap.String("table", pgErr.TableName), zap.String("constraint", pgErr.ConstraintName))
			return fmt.Errorf("%w: %v", ErrForeignKeyViolation, err)
		case PGErrDeadlockDetected:
			logger.Warn("Deadlock detected", zap.String("table", pgErr.TableName), zap.String("constraint", pgErr.ConstraintName))
			return fmt.Errorf("%w: %v", ErrDeadlockDetected, err)
		}
	}

	logger.Error(message, zap.Error(err))

	return InternalServerError{Source: err}
}

func WrapDBErrorWithKeyValue(err error, table, key, value string, logger *zap.Logger) error {
	if err == nil {
		return nil
	}

	logger.Warn("Wrapping database error with key value", zap.Error(err), zap.String("table", table), zap.String("key", key), zap.String("value", value))

	if errors.Is(err, pgx.ErrNoRows) {
		return errorPkg.NewNotFoundError(table, key, value, "")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrQueryTimeout
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGErrUniqueViolation:
			return fmt.Errorf("%w: %v", ErrUniqueViolation, err)
		case PGErrForeignKeyViolation:
			return fmt.Errorf("%w: %v", ErrForeignKeyViolation, err)
		case PGErrDeadlockDetected:
			return fmt.Errorf("%w: %v", ErrDeadlockDetected, err)
		}
	}

	return InternalServerError{Source: err}
}
