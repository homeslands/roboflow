package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestUnauthenticatedError(t *testing.T) {
	var err interface{} = new(xerrors.UnauthenticatedError)
	_, ok := err.(xerrors.Unauthenticated)
	assert.True(t, ok)
}

func TestIsUnauthenticated(t *testing.T) {
	err := xerrors.ThrowUnauthenticated(nil, "msg")
	ok := xerrors.IsUnauthenticated(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsUnauthenticated(err)
	assert.False(t, ok)
}
