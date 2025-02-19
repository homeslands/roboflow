package log

import (
	"context"
	"log/slog"
	"os"
)

type loggerCtxKey int

// CtxLoggerKey is the context key for the logger.
const CtxLoggerKey loggerCtxKey = iota

// WithLogger returns a new context that includes the provided logger.
// Useful for propagating logger configuration across different parts of the application.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerKey, logger)
}

// FromContext returns the logger from context.
// If no logger is found, it returns the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxLoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// NewLogger creates a new logger based on the provided configuration.
func NewLogger(cfg Config) *slog.Logger {
	var handler slog.Handler
	if cfg.Format == FormatJSON {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: cfg.AddSource,
		})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: cfg.AddSource,
		})
	}

	ctxHandler := NewContextHandler(handler)
	return slog.New(ctxHandler)
}

// CloneLogger creates a new logger with the same configuration as the provided logger.
func CloneLogger(logger *slog.Logger) *slog.Logger {
	return slog.New(logger.Handler())
}
