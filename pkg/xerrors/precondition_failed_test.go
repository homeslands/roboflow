package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestPreconditionFailedError(t *testing.T) {
	var err interface{} = new(xerrors.PreconditionFailedError)
	_, ok := err.(xerrors.PreconditionFailed)
	assert.True(t, ok)
}

func TestIsPreconditionFailed(t *testing.T) {
	err := xerrors.ThrowPreconditionFailed(nil, "msg")
	ok := xerrors.IsPreconditionFailed(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsPreconditionFailed(err)
	assert.False(t, ok)
}
