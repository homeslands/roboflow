package api

import (
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func (s *HTTPServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (s *HTTPServer) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusNotFound, ErrorResponse{
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
		Message:    "Route not found",
	})
}

func (s *HTTPServer) HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Code:       "method_not_allowed",
		Message:    "Method not allowed",
	})
}

func (s *HTTPServer) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	s.respondError(w, r, xerrors.ThrowInvalidArgument(err, err.Error()))
}
