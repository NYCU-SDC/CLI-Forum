package auth

import (
	"backend/internal"
	"backend/internal/jwt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

func Middleware(next http.HandlerFunc, logger *zap.Logger, requiredRoles ...string) http.HandlerFunc {
	name := "middleware/auth"
	tracer := otel.Tracer(name)

	return func(w http.ResponseWriter, r *http.Request) {
		traceCtx, span := tracer.Start(r.Context(), "AuthMiddleware")
		defer span.End()
		logger = internal.LoggerWithContext(traceCtx, logger)

		if r.Context().Value("user") == nil {
			logger.Debug("User not found in context")
			span.AddEvent("UserNotFoundInContext")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u := r.Context().Value("user").(jwt.User)
		for _, role := range requiredRoles {
			if u.HasRole(role) {
				logger.Debug("User passes role check", zap.String("role", role))
				span.AddEvent("UserPassesRoleCheck", trace.WithAttributes(attribute.String("role", role)))
				break
			}
		}

		next(w, r)
	}
}
