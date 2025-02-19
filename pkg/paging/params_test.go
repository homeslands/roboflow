package paging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

func TestNewParams_Defaults(t *testing.T) {
	params := paging.NewParams(nil, nil)

	assert.Equal(t, uint(10), params.PageSize)
	assert.Equal(t, uint(1), params.Page)
	assert.Equal(t, uint(0), params.Offset())
	assert.Equal(t, uint(10), params.Limit())
}

func TestNewParams_CustomValues(t *testing.T) {
	pageSize := uint(20)
	page := uint(3)

	params := paging.NewParams(&pageSize, &page)

	assert.Equal(t, pageSize, params.PageSize)
	assert.Equal(t, page, params.Page)
	assert.Equal(t, uint(40), params.Offset())
	assert.Equal(t, uint(20), params.Limit())
}

func TestNewParams_WithOptions(t *testing.T) {
	pageSize := uint(50)
	page := uint(2)

	params := paging.NewParams(
		&pageSize, &page,
		paging.WithDefaultPage(5),
		paging.WithMaxPageSize(30),
	)

	assert.Equal(t, uint(30), params.PageSize) // Limited by WithMaxPageSize
	assert.Equal(t, uint(5), params.Page)      // Overridden by WithDefaultPage
	assert.Equal(t, uint(120), params.Offset())
}

func TestNewParams_WithZeroPage(t *testing.T) {
	pageSize := uint(15)
	page := uint(0)

	params := paging.NewParams(&pageSize, &page)

	assert.Equal(t, pageSize, params.PageSize)
	assert.Equal(t, uint(0), params.Offset())
	assert.Equal(t, uint(15), params.Limit())
}
