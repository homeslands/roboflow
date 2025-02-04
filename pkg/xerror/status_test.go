package xerror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		status         Status
		expectedStatus int
	}{
		{StatusUnauthorized, http.StatusUnauthorized},
		{StatusForbidden, http.StatusForbidden},
		{StatusNotFound, http.StatusNotFound},
		{StatusUnprocessableEntity, http.StatusUnprocessableEntity},
		{StatusConflict, http.StatusConflict},
		{StatusTooManyRequests, http.StatusTooManyRequests},
		{StatusBadRequest, http.StatusBadRequest},
		{StatusValidationFailed, http.StatusBadRequest},
		{StatusUnknown, http.StatusInternalServerError},
		{StatusInternalServerError, http.StatusInternalServerError},
		{StatusTimeout, http.StatusGatewayTimeout},
		{StatusNotImplemented, http.StatusNotImplemented},
		{StatusBadGateway, http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			assert.Equal(t, tt.expectedStatus, tt.status.HTTPStatus())
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		status         Status
		expectedString string
	}{
		{StatusUnauthorized, "Unauthorized"},
		{StatusForbidden, "Forbidden"},
		{StatusNotFound, "Not found"},
		{StatusUnprocessableEntity, "Unprocessable entity"},
		{StatusConflict, "Conflict"},
		{StatusTooManyRequests, "Too many requests"},
		{StatusBadRequest, "Bad request"},
		{StatusValidationFailed, "Validation failed"},
		{StatusUnknown, "Unknown"},
		{StatusInternalServerError, "Internal server error"},
		{StatusTimeout, "Timeout"},
		{StatusNotImplemented, "Not implemented"},
		{StatusBadGateway, "Bad gateway"},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			assert.Equal(t, tt.expectedString, tt.status.String())
		})
	}
}
