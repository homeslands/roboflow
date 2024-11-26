package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestUnavailableError(t *testing.T) {
	var err interface{} = new(xerrors.UnavailableError)
	_, ok := err.(xerrors.Unavailable)
	assert.True(t, ok)
}

func TestIsUnavailable(t *testing.T) {
	err := xerrors.ThrowUnavailable(nil, "msg")
	ok := xerrors.IsUnavailable(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsUnavailable(err)
	assert.False(t, ok)
}
