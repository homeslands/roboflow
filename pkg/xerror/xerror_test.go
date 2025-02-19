package xerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewXError(t *testing.T) {
	t.Run("With parent error", func(t *testing.T) {
		parentErr := errors.New("parent error")
		xErr := NewXError(parentErr, StatusBadRequest, "test.invalidInput", "Invalid input provided")

		assert.Equal(t, StatusBadRequest, xErr.Status(), "Status should match")
		assert.Equal(t, "test.invalidInput", xErr.MsgID(), "MsgID should match")
		assert.Equal(t, "Invalid input provided", xErr.Msg(), "Message should match")
		assert.Equal(t, parentErr, xErr.Parent(), "Parent error should match")

		expectedErrorString := "ID=test.invalidInput, Msg=Invalid input provided, Parent=(parent error)"
		assert.Equal(t, expectedErrorString, xErr.Error(), "Error string should match")
	})

	t.Run("Without parent error", func(t *testing.T) {
		xErr := NewXError(nil, StatusBadRequest, "test.invalidInput", "Invalid input provided")

		assert.Equal(t, StatusBadRequest, xErr.Status(), "Status should match")
		assert.Equal(t, "test.invalidInput", xErr.MsgID(), "MsgID should match")
		assert.Equal(t, "Invalid input provided", xErr.Msg(), "Message should match")
		assert.Nil(t, xErr.Parent(), "Parent error should be nil")

		expectedErrorString := "ID=test.invalidInput, Msg=Invalid input provided"
		assert.Equal(t, expectedErrorString, xErr.Error(), "Error string should match")
	})
}

func TestWithParent(t *testing.T) {
	xErr := NewXError(nil, StatusNotFound, "resource.notFound", "Resource not found")
	newParent := errors.New("new parent error")
	xErr.WithParent(newParent)

	assert.Equal(t, newParent, xErr.Parent(), "Parent error should be updated")
}

func TestUnwrap(t *testing.T) {
	parentErr := errors.New("parent error")
	xErr := NewXError(parentErr, StatusInternalServerError, "server.error", "An internal server error occurred")

	assert.Equal(t, parentErr, xErr.Unwrap(), "Unwrapped error should match the parent error")
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name           string
		constructor    func(parent error, msgID, msg string) XError
		expectedStatus Status
	}{
		{"Unauthorized", Unauthorized, StatusUnauthorized},
		{"Forbidden", Forbidden, StatusForbidden},
		{"NotFound", NotFound, StatusNotFound},
		{"UnprocessableEntity", UnprocessableEntity, StatusUnprocessableEntity},
		{"Conflict", Conflict, StatusConflict},
		{"TooManyRequests", TooManyRequests, StatusTooManyRequests},
		{"BadRequest", BadRequest, StatusBadRequest},
		{"InternalServerError", InternalServerError, StatusInternalServerError},
		{"Timeout", Timeout, StatusTimeout},
		{"NotImplemented", NotImplemented, StatusNotImplemented},
		{"BadGateway", BadGateway, StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parentErr := errors.New("parent error")
			xErr := tt.constructor(parentErr, "test.msgID", "test message")

			assert.Equal(t, tt.expectedStatus, xErr.Status(), "Status should match")
			assert.Equal(t, "test.msgID", xErr.MsgID(), "MsgID should match")
			assert.Equal(t, "test message", xErr.Msg(), "Message should match")
			assert.Equal(t, parentErr, xErr.Parent(), "Parent error should match")
		})
	}
}

func TestValidationFailed(t *testing.T) {
	parentErr := errors.New("parent error")
	xErr := ValidationFailed(parentErr, "Validation failed for field 'name'")

	assert.Equal(t, StatusValidationFailed, xErr.Status(), "Status should be ValidationFailed")
	assert.Equal(t, "validationFailed", xErr.MsgID(), "MsgID should be 'validationFailed'")
	assert.Equal(t, "Validation failed for field 'name'", xErr.Msg(), "Message should match")
	assert.Equal(t, parentErr, xErr.Parent(), "Parent error should match")
}
