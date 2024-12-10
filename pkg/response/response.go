package response

import (
	"encoding/json"
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

// JSON writes the response as JSON.
func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}
	json.NewEncoder(w).Encode(data)
}

// Error writes the error response as JSON.
func Error(w http.ResponseWriter, err error) {
	res := ErrorToResponse(err)
	JSON(w, res.StatusCode, res)
}
