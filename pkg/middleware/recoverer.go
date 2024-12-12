package middleware

// ref: https://github.com/go-chi/chi/blob/master/middleware/recoverer.go

import (
	"encoding/json"
	"log/slog"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Recoverer prints a stack trace of the last function call.
func Recoverer(next http.Handler, isProd bool, log *slog.Logger) http.Handler {
	errorResponse := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "internal_server_error",
		Message: "Internal Server Error",
	}

	response, err := json.Marshal(errorResponse)
	if err != nil {
		panic(err)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				if isProd {
					log.Error("panic", slog.Any("recover", rvr))
				} else {
					chimiddleware.PrintPrettyStack(rvr)
				}

				if r.Header.Get("Connection") != "Upgrade" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(response)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
