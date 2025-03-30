package internal

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

func TraceMiddleware(next func(w http.ResponseWriter, r *http.Request), logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
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
