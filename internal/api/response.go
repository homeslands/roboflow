package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	StatusCode int               `json:"-"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    []ValidationError `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// responseJSON writes the response as JSON to the response writer.
func (s *HTTPServer) respondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}

	json.NewEncoder(w).Encode(data)
}

// respondError writes the error response to the response writer.
func (s *HTTPServer) respondError(w http.ResponseWriter, r *http.Request, err error) {
	res, isInternalErr := ErrorToResponse(err)
	if isInternalErr {
		s.log.ErrorContext(r.Context(), err.Error(), slog.Any("error", err))
	} else {
		s.log.DebugContext(r.Context(), err.Error(), slog.Any("error", err))
	}

	s.respondJSON(w, res.StatusCode, res)
}
