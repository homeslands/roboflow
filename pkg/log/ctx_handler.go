package log

import (
	"context"
	"log/slog"

	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
)

type ContextHandler struct {
	attrs []slog.Attr
	h     slog.Handler
}

func NewContextHandler(h slog.Handler) *ContextHandler {
	return &ContextHandler{h: h}
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := middleware.GetReqID(ctx); traceID != "" {
		r.AddAttrs(slog.String("trace_id", traceID))
	}

	if requestID := middleware.GetReqID(ctx); requestID != "" {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	return h.h.Handle(ctx, r)
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// WithAttrs returns a new Handler with the given attributes added to the existing attributes.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//nolint:gocritic
	newAttrs := append(h.attrs, attrs...)
	return &ContextHandler{
		attrs: newAttrs,
		h:     h.h,
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		attrs: h.attrs,
		h:     h.h.WithGroup(name),
	}
}
