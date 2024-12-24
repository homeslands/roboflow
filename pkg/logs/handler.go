package logs

import (
	"context"
	"log/slog"

	"github.com/tuanvumaihuynh/roboflow/pkg/middleware"
)

type Handler struct {
	slog.Handler
}

func (h Handler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		r.Add("req_id", slog.StringValue(requestID))
	}

	return h.Handler.Handle(ctx, r)
}

func (c Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c
}
