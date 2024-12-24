package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/internal/api"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

// FieldErrorMock is a mock implementation of validator.FieldError for testing purposes.
type mockFieldError struct {
	validator.FieldError
	tag   string
	field string
}

func (e mockFieldError) Tag() string { return e.tag }

func (e mockFieldError) Field() string { return e.field }

func TestErrorToResponse(t *testing.T) {
	type expectedResponse struct {
		res        api.ErrorResponse
		isInternal bool
	}
	tests := []struct {
		name     string
		err      error
		expected expectedResponse
	}{
		{
			name: "nil error",
			err:  nil,
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Code:       "internal_error",
					Message:    "Internal server error",
				},
				isInternal: true,
			},
		},
		{
			name: "AlreadyExistsError",
			err:  xerrors.ThrowAlreadyExists(nil, "Already exists"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusConflict,
					Code:       "already_exists",
					Message:    "Already exists",
				},
				isInternal: false,
			},
		},
		{
			name: "InternalError",
			err:  &xerrors.InternalError{},
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Code:       "internal_server_error",
					Message:    "Internal server error",
				},
				isInternal: true,
			},
		},
		{
			name: "InvalidArgumentError",
			err:  xerrors.ThrowInvalidArgument(nil, "Invalid argument"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Code:       "invalid_argument",
					Message:    "Invalid argument",
				},
				isInternal: false,
			},
		},
		{
			name: "PreconditionFailedError",
			err:  xerrors.ThrowPreconditionFailed(nil, "Precondition failed"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Code:       "precondition_failed",
					Message:    "Precondition failed",
				},
				isInternal: false,
			},
		},
		{
			name: "NotFoundError",
			err:  xerrors.ThrowNotFound(nil, "Not found"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Code:       "not_found",
					Message:    "Not found",
				},
				isInternal: false,
			},
		},
		{
			name: "PermissionDeniedError",
			err:  xerrors.ThrowPermissionDenied(nil, "Permission denied"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusForbidden,
					Code:       "permission_denied",
					Message:    "Permission denied",
				},
				isInternal: false,
			},
		},
		{
			name: "UnauthenticatedError",
			err:  xerrors.ThrowUnauthenticated(nil, "Unauthenticated"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusUnauthorized,
					Code:       "unauthenticated",
					Message:    "Unauthenticated",
				},
				isInternal: false,
			},
		},
		{
			name: "UnavailableError",
			err:  xerrors.ThrowUnavailable(nil, "Service unavailable"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusServiceUnavailable,
					Code:       "unavailable",
					Message:    "Service unavailable",
				},
				isInternal: true,
			},
		},
		{
			name: "Unknown error type",
			err:  errors.New("unknown error"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Code:       "internal_error",
					Message:    "Internal server error",
				},
				isInternal: true,
			},
		},
		{
			name: "XError wrapped",
			err:  xerrors.ThrowError(errors.New("wrapped error"), "Outer error"),
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Code:       "internal_error",
					Message:    "Internal server error",
				},
				isInternal: true,
			},
		},
		{
			name: "ValidationError",
			err: validator.ValidationErrors{
				&mockFieldError{
					tag:   "required",
					field: "Field1",
				},
				&mockFieldError{
					tag:   "email",
					field: "Field2",
				},
			},
			expected: expectedResponse{
				res: api.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Code:       "validation_error",
					Message:    "Validation error",
					Details: []api.ValidationError{
						{Field: "Field1", Message: "'Field1' required"},
						{Field: "Field2", Message: "'Field2' email"},
					},
				},
				isInternal: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isInternal := api.ErrorToResponse(tt.err)
			assert.Equal(t, tt.expected, expectedResponse{result, isInternal})
		})
	}
}
