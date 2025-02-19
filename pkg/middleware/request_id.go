package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDCtxKey int

// CtxRequestIDKey is the key that holds the unique request ID in a request context.
const CtxRequestIDKey requestIDCtxKey = iota

// RequestIDHeader is the name of the HTTP Header which contains the request id.
var RequestIDHeader = "X-Request-ID"

// RequestID is a middleware that injects a request ID into the context of each request.
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), CtxRequestIDKey, reqID)
		w.Header().Set(RequestIDHeader, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(CtxRequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
