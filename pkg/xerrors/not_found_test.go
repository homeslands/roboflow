package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestNotFoundError(t *testing.T) {
	var notFoundError interface{} = new(xerrors.NotFoundError)
	_, ok := notFoundError.(xerrors.NotFound)
	assert.True(t, ok)
}

func TestIsNotFound(t *testing.T) {
	err := xerrors.ThrowNotFound(nil, "msg")
	ok := xerrors.IsNotFound(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsNotFound(err)
	assert.False(t, ok)
}
