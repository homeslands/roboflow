package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestInternalError(t *testing.T) {
	var err interface{} = new(xerrors.InternalError)
	_, ok := err.(xerrors.Internal)
	assert.True(t, ok)
}

func TestIsInternal(t *testing.T) {
	err := xerrors.ThrowInternal(nil, "msg")
	ok := xerrors.IsInternal(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsInternal(err)
	assert.False(t, ok)
}
