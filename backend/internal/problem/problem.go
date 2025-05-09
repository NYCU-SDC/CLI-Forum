package problem

import (
	"backend/internal/database"
	errorPkg "backend/internal/error"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

// Problem represents a problem detail as defined in RFC 7807
type Problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`

	// Type indicates the URI that identifies the problem type.
	// In production, this would point to the project's documentation.
	// For demonstration purposes, we use an MDN URI here.
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func WriteError(ctx context.Context, w http.ResponseWriter, err error, logger *zap.Logger) {
	_, span := otel.Tracer("problem/problem").Start(ctx, "WriteError")
	defer span.End()

	if err == nil {
		return
	}

	var problem Problem
	var notFoundError errorPkg.NotFoundError
	var validationErrors validator.ValidationErrors
	var internalDbError database.InternalServerError
	switch {
	case errors.As(err, &notFoundError):
		problem = NewNotFoundProblem(err.Error())
	case errors.As(err, &validationErrors):
		problem = NewValidateProblem(validationErrors.Error())
	case errors.Is(err, errorPkg.ErrUserAlreadyExists):
		problem = NewValidateProblem("User already exists")
	case errors.Is(err, errorPkg.ErrCredentialInvalid):
		problem = NewUnauthorizedProblem("Invalid username or password")
	case errors.Is(err, errorPkg.ErrForbidden):
		problem = NewForbiddenProblem("Make sure you have the right permissions")
	case errors.Is(err, errorPkg.ErrUnauthorized):
		problem = NewUnauthorizedProblem("You must be logged in to access this resource")
	case errors.Is(err, errorPkg.ErrInvalidUUID):
		problem = NewValidateProblem(err.Error())
	case errors.As(err, &internalDbError):
		problem = NewInternalServerProblem("Internal server error")
	case errors.Is(err, errorPkg.ErrInvalidUUID):
		problem = NewValidateProblem("Invalid UUID format")
	default:
		problem = NewInternalServerProblem("Internal server error")
	}

	logger.Warn("Handling "+problem.Title, zap.String("problem", problem.Title), zap.Error(err), zap.Int("status", problem.Status), zap.String("type", problem.Type), zap.String("detail", problem.Detail))

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	jsonBytes, err := json.Marshal(problem)
	if err != nil {
		logger.Error("Failed to marshal problem response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		logger.Error("Failed to write problem response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewInternalServerProblem(detail string) Problem {
	return Problem{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Type:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500",
		Detail: detail,
	}
}

func NewNotFoundProblem(detail string) Problem {
	return Problem{
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Type:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404",
		Detail: detail,
	}
}

func NewValidateProblem(detail string) Problem {
	return Problem{
		Title:  "Validation Problem",
		Status: http.StatusBadRequest,
		Type:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400",
		Detail: detail,
	}
}

func NewUnauthorizedProblem(detail string) Problem {
	return Problem{
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
		Type:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401",
		Detail: detail,
	}
}

func NewForbiddenProblem(detail string) Problem {
	return Problem{
		Title:  "Forbidden",
		Status: http.StatusForbidden,
		Type:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403",
		Detail: detail,
	}
}
