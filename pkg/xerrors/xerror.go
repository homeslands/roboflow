package xerrors

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Error = (*XError)(nil)

type XError struct {
	Parent  error
	Message string
	Code    string
}

type Option func(*XError)

func WithCode(code string) Option {
	return func(err *XError) {
		err.Code = code
	}
}

// ThrowError creates a new XError.
// Default error code is "unknown".
func ThrowError(parent error, message string, options ...Option) error {
	x := CreateXError(parent, "unknown", message)
	for _, opt := range options {
		opt(x)
	}
	return x
}

func CreateXError(parent error, code, message string) *XError {
	return &XError{
		Parent:  parent,
		Code:    code,
		Message: message,
	}
}

func (err *XError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("Code=%s Message=%s Parent=(%v)", err.Code, err.Message, err.Parent)
	}
	return fmt.Sprintf("Code=%s Message=%s", err.Code, err.Message)
}

func (err *XError) Unwrap() error {
	return err.GetParent()
}

func (err *XError) GetParent() error {
	return err.Parent
}

func (err *XError) GetMessage() string {
	return err.Message
}

func (err *XError) SetMessage(msg string) {
	err.Message = msg
}

func (err *XError) GetCode() string {
	return err.Code
}

func (err *XError) Is(target error) bool {
	t, ok := target.(*XError)
	if !ok {
		return false
	}
	if t.Code != "" && t.Code != err.Code {
		return false
	}
	if t.Message != "" && t.Message != err.Message {
		return false
	}
	if t.Parent != nil && !errors.Is(err.Parent, t.Parent) {
		return false
	}

	return true
}

func (err *XError) As(target interface{}) bool {
	_, ok := target.(**XError)
	if !ok {
		return false
	}
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.ValueOf(err))
	return true
}
