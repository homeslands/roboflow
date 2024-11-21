package xerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestPermissionDeniedError(t *testing.T) {
	var err interface{} = new(xerrors.PermissionDeniedError)
	_, ok := err.(xerrors.PermissionDenied)
	assert.True(t, ok)
}

func TestIsPermissionDenied(t *testing.T) {
	err := xerrors.ThrowPermissionDenied(nil, "msg")
	ok := xerrors.IsPermissionDenied(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = xerrors.IsPermissionDenied(err)
	assert.False(t, ok)
}
