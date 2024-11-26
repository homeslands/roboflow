package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestAlreadyExistsError(t *testing.T) {
	var alreadyExistsError interface{} = new(xerrors.AlreadyExistsError)
	_, ok := alreadyExistsError.(xerrors.AlreadyExists)
	assert.True(t, ok)
}

func TestIsErrorAlreadyExists(t *testing.T) {
	err := xerrors.ThrowAlreadyExists(nil, "msg")
	ok := xerrors.IsErrorAlreadyExists(err)
	assert.True(t, ok)

	err = errors.New("Already Exists!")
	ok = xerrors.IsErrorAlreadyExists(err)
	assert.False(t, ok)
}
