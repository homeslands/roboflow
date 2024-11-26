package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestInvalidArgumentError(t *testing.T) {
	var invalidArgumentError interface{} = new(xerrors.InvalidArgumentError)
	_, ok := invalidArgumentError.(xerrors.InvalidArgument)
	assert.True(t, ok)
}

func TestIsErrorInvalidArgument(t *testing.T) {
	err := xerrors.ThrowInvalidArgument(nil, "msg")
	ok := xerrors.IsErrorInvalidArgument(err)
	assert.True(t, ok)

	err = errors.New("I am invalid!")
	ok = xerrors.IsErrorInvalidArgument(err)
	assert.False(t, ok)
}
