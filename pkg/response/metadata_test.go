package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginationMetadata(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		page         int
		pageSize     int
		total        int64
		expectedPage int
		expectedSize int
		expectedTotal int64
	}{
		{
			name:              "success - exact division",
			page:              1,
			pageSize:          10,
			total:             20,
			expectedPage:      1,
			expectedSize:      10,
			expectedTotal:     20,
		},
		{
			name:              "success - requires rounding up",
			page:              2,
			pageSize:          10,
			total:             25,
			expectedPage:      2,
			expectedSize:      10,
			expectedTotal:     25,
		},
		{
			name:              "success - single page",
			page:              1,
			pageSize:          10,
			total:             5,
			expectedPage:      1,
			expectedSize:      10,
			expectedTotal:     5,
		},
		{
			name:              "success - empty result",
			page:              1,
			pageSize:          10,
			total:             0,
			expectedPage:      1,
			expectedSize:      10,
			expectedTotal:     0,
		},
		{
			name:              "success - large total",
			page:              5,
			pageSize:          50,
			total:             1234,
			expectedPage:      5,
			expectedSize:      50,
			expectedTotal:     1234,
		},
		{
			name:              "success - page size larger than total",
			page:              1,
			pageSize:          100,
			total:             50,
			expectedPage:      1,
			expectedSize:      100,
			expectedTotal:     50,
		},
		{
			name:              "success - single item per page",
			page:              3,
			pageSize:          1,
			total:             5,
			expectedPage:      3,
			expectedSize:      1,
			expectedTotal:     5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := NewPaginationMetadata(tc.page, tc.pageSize, tc.total)

			assert.Equal(t, tc.expectedPage, result.Page)
			assert.Equal(t, tc.expectedSize, result.PageSize)
			assert.Equal(t, tc.expectedTotal, result.Total)
		})
	}
}

func TestPaginationMetadata_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PaginationMetadata has correct JSON tags", func(t *testing.T) {
		t.Parallel()

		metadata := PaginationMetadata{
			Page:     1,
			PageSize: 10,
			Total:    100,
		}

		assert.Equal(t, 1, metadata.Page)
		assert.Equal(t, 10, metadata.PageSize)
		assert.Equal(t, int64(100), metadata.Total)
	})

	t.Run("PaginationMetadata with zero values", func(t *testing.T) {
		t.Parallel()

		metadata := PaginationMetadata{}

		assert.Equal(t, 0, metadata.Page)
		assert.Equal(t, 0, metadata.PageSize)
		assert.Equal(t, int64(0), metadata.Total)
	})
}

