package xsort

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewList(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   []Sort
		expectErr  bool
		errMessage string
	}{
		{
			name:  "Single column, ascending",
			input: "created_at",
			expected: []Sort{
				{col: "created_at", order: OrderASC},
			},
			expectErr: false,
		},
		{
			name:  "Single column, descending",
			input: "-updated_at",
			expected: []Sort{
				{col: "updated_at", order: OrderDESC},
			},
			expectErr: false,
		},
		{
			name:  "Multiple columns",
			input: "name,-created_at,updated_at",
			expected: []Sort{
				{col: "name", order: OrderASC},
				{col: "created_at", order: OrderDESC},
				{col: "updated_at", order: OrderASC},
			},
			expectErr: false,
		},
		{
			name:      "Empty input",
			input:     "",
			expected:  nil,
			expectErr: false,
		},
		{
			name:       "Invalid input with trailing space",
			input:      "name ,created_at",
			expectErr:  true,
			errMessage: "invalid sort column",
		},
		{
			name:       "Invalid input with leading space",
			input:      " -created_at",
			expectErr:  true,
			errMessage: "invalid sort column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewList(&tt.input)

			if tt.expectErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errMessage, "error message does not match")
			} else {
				require.NoError(t, err, "unexpected error occurred")
				assert.Equal(t, tt.expected, result, "result mismatch for input '%s'", tt.input)
			}
		})
	}
}
