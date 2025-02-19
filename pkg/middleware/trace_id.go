package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type traceIDCtxKey int

// CtxTraceIDKey is the key that holds the unique trace ID in a request context.
const CtxTraceIDKey traceIDCtxKey = iota

// TraceIDHeader is the name of the HTTP Header which contains the trace id.
var TraceIDHeader = "X-Trace-ID"

// TraceID is a middleware that injects a trace ID into the context of each request.
func TraceID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			// Panic?
			traceID = uuid.Must(uuid.NewV7()).String()
		}

		ctx := context.WithValue(r.Context(), CtxTraceIDKey, traceID)
		w.Header().Set(TraceIDHeader, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// GetTraceID returns a trace ID from the given context if one is present.
// Returns the empty string if a trace ID cannot be found.
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(CtxTraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
