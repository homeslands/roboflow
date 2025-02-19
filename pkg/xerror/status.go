package xerror

import "net/http"

type Status string

const (
	// HTTP Status code 500.
	StatusUnknown Status = "Unknown"
	// HTTP Status code 401.
	StatusUnauthorized Status = "Unauthorized"
	// HTTP status code 403.
	StatusForbidden Status = "Forbidden"
	// HTTP status code 404.
	StatusNotFound Status = "Not found"
	// HTTP status code 422.
	StatusUnprocessableEntity Status = "Unprocessable entity"
	// HTTP status code 409.
	StatusConflict Status = "Conflict"
	// HTTP status code 429.
	StatusTooManyRequests Status = "Too many requests"
	// HTTP status code 400.
	StatusBadRequest Status = "Bad request"
	// HTTP status code 400.
	StatusValidationFailed Status = "Validation failed"
	// HTTP status code 500.
	StatusInternalServerError Status = "Internal server error"
	// HTTP status code 504.
	StatusTimeout Status = "Timeout"
	// HTTP status code 501.
	StatusNotImplemented Status = "Not implemented"
	// HTTP status code 502.
	StatusBadGateway Status = "Bad gateway"
)

func (s Status) HTTPStatus() int {
	switch s {
	case StatusUnauthorized:
		return http.StatusUnauthorized
	case StatusForbidden:
		return http.StatusForbidden
	case StatusNotFound:
		return http.StatusNotFound
	case StatusUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case StatusConflict:
		return http.StatusConflict
	case StatusTooManyRequests:
		return http.StatusTooManyRequests
	case StatusBadRequest:
		return http.StatusBadRequest
	case StatusValidationFailed:
		return http.StatusBadRequest
	case StatusUnknown, StatusInternalServerError:
		return http.StatusInternalServerError
	case StatusTimeout:
		return http.StatusGatewayTimeout
	case StatusNotImplemented:
		return http.StatusNotImplemented
	case StatusBadGateway:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func (s Status) String() string {
	return string(s)
}
