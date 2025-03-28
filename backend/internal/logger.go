package internal

import (
	"context"
	"go.uber.org/zap"
)

// ZapProductionConfig returns a zap.Config same as zap.NewProduction() but without sampling
func ZapProductionConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func LoggerWithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if ctx == nil {
		return logger
	}

	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return logger.With(zap.String("trace_id", traceID))
	}

	return logger
}
