package internal

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

func TraceMiddleware(next http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	name := "middleware/trace"
	tracer := otel.Tracer(name)
	propagator := otel.GetTextMapPropagator()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		upstream := trace.SpanFromContext(ctx).SpanContext()

		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()

		span.SetAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("query", r.URL.RawQuery),
		)
		span.AddEvent("HTTPRequestStarted")

		if upstream.HasTraceID() {
			logger.Debug("Upstream trace available", zap.String("trace_id", upstream.TraceID().String()))
		} else {
			logger.Debug("No upstream trace available, creating a new one", zap.String("trace_id", span.SpanContext().TraceID().String()))
		}

		next(w, r.WithContext(ctx))
	}
}

func RecoverMiddleware(next http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	name := "middleware/trace"
	tracer := otel.Tracer(name)

	return func(w http.ResponseWriter, r *http.Request) {
		traceCtx, span := tracer.Start(r.Context(), "RecoverMiddleware")
		defer func() {
			if err := recover(); err != nil {
				span.AddEvent("PanicRecovered", trace.WithAttributes(attribute.String("panic", fmt.Sprintf("%v", err))))
				logger.Error("Recovered from panic", zap.Any("error", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

			span.End()
		}()

		next(w, r.WithContext(traceCtx))
	}
}
