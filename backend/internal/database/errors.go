package database

import (
	"backend/internal"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PGErrUniqueViolation     = "23505"
	PGErrForeignKeyViolation = "23503"
	PGErrDeadlockDetected    = "40P01"
)

var (
	ErrDuplicateKey        = errors.New("duplicate key value")
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

func WrapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %v", internal.ErrNotFound, err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrQueryTimeout
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGErrUniqueViolation:
			return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		case PGErrForeignKeyViolation:
			return fmt.Errorf("%w: %v", ErrForeignKeyViolation, err)
		case PGErrDeadlockDetected:
			return fmt.Errorf("%w: %v", ErrDeadlockDetected, err)
		}
	}

	return InternalServerError{Source: err}
}

func WrapDBErrorWithKeyValue(err error, table, key, value string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return internal.NewNotFoundError(table, key, value, "")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrQueryTimeout
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGErrUniqueViolation:
			return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		case PGErrForeignKeyViolation:
			return fmt.Errorf("%w: %v", ErrForeignKeyViolation, err)
		case PGErrDeadlockDetected:
			return fmt.Errorf("%w: %v", ErrDeadlockDetected, err)
		}
	}

	return InternalServerError{Source: err}
}
