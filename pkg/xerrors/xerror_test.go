package xerrors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

func TestErrorMethod(t *testing.T) {
	err := xerrors.ThrowError(nil, "msg", xerrors.WithCode("code"))
	expected := "Code=code Message=msg"
	assert.Equal(t, expected, err.Error())

	err = xerrors.ThrowError(err, "subMsg", xerrors.WithCode("subCode"))
	subExptected := "Code=subCode Message=subMsg Parent=(Code=code Message=msg)"
	assert.Equal(t, subExptected, err.Error())
}
