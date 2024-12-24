package paging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

func TestNewParams_Defaults(t *testing.T) {
	params := paging.NewParams(nil, nil)

	assert.Equal(t, int32(10), params.PageSize)
	assert.Equal(t, int32(1), params.Page)
	assert.Equal(t, int32(0), params.Offset())
	assert.Equal(t, int32(10), params.Limit())
}

func TestNewParams_CustomValues(t *testing.T) {
	pageSize := int32(20)
	page := int32(3)

	params := paging.NewParams(&pageSize, &page)

	assert.Equal(t, pageSize, params.PageSize)
	assert.Equal(t, page, params.Page)
	assert.Equal(t, int32(40), params.Offset())
	assert.Equal(t, int32(20), params.Limit())
}

func TestNewParams_WithOptions(t *testing.T) {
	pageSize := int32(50)
	page := int32(2)

	params := paging.NewParams(
		&pageSize, &page,
		paging.WithDefaultPage(5),
		paging.WithMaxPageSize(30),
	)

	assert.Equal(t, int32(30), params.PageSize) // Limited by WithMaxPageSize
	assert.Equal(t, int32(5), params.Page)      // Overridden by WithDefaultPage
	assert.Equal(t, int32(120), params.Offset())
}

func TestNewParams_InvalidPageSize(t *testing.T) {
	invalidPageSize := int32(-5)
	params := paging.NewParams(&invalidPageSize, nil)

	assert.Equal(t, int32(10), params.PageSize) // Falls back to defaultPageSize
	assert.Equal(t, int32(1), params.Page)
}

func TestNewParams_WithZeroPage(t *testing.T) {
	pageSize := int32(15)
	page := int32(0)

	params := paging.NewParams(&pageSize, &page)

	assert.Equal(t, pageSize, params.PageSize)
	assert.Equal(t, int32(0), params.Offset())
	assert.Equal(t, int32(15), params.Limit())
}
