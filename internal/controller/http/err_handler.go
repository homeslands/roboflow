package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/apierr"
)

func (s HTTPService) handleRequestError(w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	res := apierr.New(err)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Warn("error encoding error request", slog.Any("error", err))
	}
}

func (s HTTPService) handleResponseError(w http.ResponseWriter, r *http.Request, err error) {
	res := apierr.New(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)

	logLevel := slog.LevelInfo
	if res.StatusCode >= 500 {
		logLevel = slog.LevelError
	} else if res.StatusCode >= 400 {
		logLevel = slog.LevelWarn
	}
	slog.Log(r.Context(), logLevel, "http error", slog.Any("error", err))

	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("error encoding error response", slog.Any("error", err))
	}
}
