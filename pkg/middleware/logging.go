package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// ref: https://github.com/moul/chizap/blob/main/chizap.go

// Opts contains the middleware configuration.
type Opts struct {
	// WithReferer enables logging the "Referer" HTTP header value.
	WithReferer bool

	// WithUserAgent enables logging the "User-Agent" HTTP header value.
	WithUserAgent bool
}

// Logging returns a logger middleware for chi, that implements the http.Handler interface.
func Logging(log *slog.Logger) func(next http.Handler) http.Handler {
	if log == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				log.Info("HTTP request",
					slog.Int("status", ww.Status()),
					slog.String("path", r.URL.Path),
					slog.Duration("latency", time.Since(t1)),
					slog.String("req_id", GetReqID(r.Context())))
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
