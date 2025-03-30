package jwt

import (
	"backend/internal"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type Verifier interface {
	Parse(ctx context.Context, tokenString string) (User, error)
}

type MiddlewareService struct {
	logger *zap.Logger
	tracer trace.Tracer

	verifier Verifier
}

func NewMiddleware(verifier Verifier, logger *zap.Logger) MiddlewareService {
	name := "middleware/jwt"
	tracer := otel.Tracer(name)

	return MiddlewareService{
		tracer:   tracer,
		logger:   logger,
		verifier: verifier,
	}
}

func (m MiddlewareService) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceCtx, span := m.tracer.Start(r.Context(), "JWTMiddleware")
		defer span.End()
		logger := internal.LoggerWithContext(traceCtx, m.logger)

		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Warn("Authorization header required")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.verifier.Parse(traceCtx, token)
		if err != nil {
			logger.Warn("Authorization header invalid", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logger.Debug("Authorization header valid")
		r = r.WithContext(context.WithValue(traceCtx, "user", user))
		next.ServeHTTP(w, r)
	}
}
