package internal

import (
	errorPkg "backend/internal/error"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func ParseAndValidateRequestBody(ctx context.Context, v *validator.Validate, r *http.Request, s interface{}) error {
	_, span := otel.Tracer("internal/handler").Start(ctx, "ParseAndValidateRequestBody")
	defer span.End()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Error closing request body:", err)
		}
	}()

	err = json.Unmarshal(bodyBytes, s)
	if err != nil {
		span.RecordError(err)
		return err
	}

	err = v.Struct(s)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseUUID(value string) (uuid.UUID, error) {
	var res uuid.UUID
	err := res.Scan(value)
	if err != nil {
		return res, errorPkg.ErrInvalidUUID
	}
	return res, nil
}
