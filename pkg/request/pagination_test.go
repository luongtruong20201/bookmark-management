package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationBase_ValidateAndNormalize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		page         int
		pageSize     int
		expectedPage int
		expectedSize int
	}{
		{
			name:         "success - valid page and pageSize",
			page:         2,
			pageSize:     20,
			expectedPage: 2,
			expectedSize: 20,
		},
		{
			name:         "success - default page when page is 0",
			page:         0,
			pageSize:     10,
			expectedPage: DefaultPage,
			expectedSize: 10,
		},
		{
			name:         "success - default page when page is negative",
			page:         -1,
			pageSize:     10,
			expectedPage: DefaultPage,
			expectedSize: 10,
		},
		{
			name:         "success - default pageSize when pageSize is 0",
			page:         1,
			pageSize:     0,
			expectedPage: 1,
			expectedSize: DefaultPageSize,
		},
		{
			name:         "success - default pageSize when pageSize is negative",
			page:         1,
			pageSize:     -5,
			expectedPage: 1,
			expectedSize: DefaultPageSize,
		},
		{
			name:         "success - pageSize exceeds max, should cap at MaxPageSize",
			page:         1,
			pageSize:     200,
			expectedPage: 1,
			expectedSize: MaxPageSize,
		},
		{
			name:         "success - both page and pageSize need defaults",
			page:         0,
			pageSize:     0,
			expectedPage: DefaultPage,
			expectedSize: DefaultPageSize,
		},
		{
			name:         "success - pageSize at max boundary",
			page:         1,
			pageSize:     MaxPageSize,
			expectedPage: 1,
			expectedSize: MaxPageSize,
		},
		{
			name:         "success - pageSize just below max",
			page:         1,
			pageSize:     MaxPageSize - 1,
			expectedPage: 1,
			expectedSize: MaxPageSize - 1,
		},
		{
			name:         "success - pageSize just above max",
			page:         1,
			pageSize:     MaxPageSize + 1,
			expectedPage: 1,
			expectedSize: MaxPageSize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := &PaginationBase{
				Page:     tc.page,
				PageSize: tc.pageSize,
			}

			page, pageSize := p.ValidateAndNormalize()

			assert.Equal(t, tc.expectedPage, page)
			assert.Equal(t, tc.expectedSize, pageSize)
		})
	}
}

func TestPaginationBase_ToOffsetLimit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "success - first page",
			page:           1,
			pageSize:       10,
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "success - second page",
			page:           2,
			pageSize:       10,
			expectedOffset: 10,
			expectedLimit:  10,
		},
		{
			name:           "success - third page",
			page:           3,
			pageSize:       20,
			expectedOffset: 40,
			expectedLimit:  20,
		},
		{
			name:           "success - page with normalization (page 0 becomes 1)",
			page:           0,
			pageSize:       10,
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "success - pageSize normalization (0 becomes default)",
			page:           2,
			pageSize:       0,
			expectedOffset: DefaultPageSize,
			expectedLimit:  DefaultPageSize,
		},
		{
			name:           "success - pageSize capped at max",
			page:           1,
			pageSize:       200,
			expectedOffset: 0,
			expectedLimit:  MaxPageSize,
		},
		{
			name:           "success - large page number",
			page:           100,
			pageSize:       50,
			expectedOffset: 4950, // (100-1) * 50
			expectedLimit:  50,
		},
		{
			name:           "success - single item per page",
			page:           5,
			pageSize:       1,
			expectedOffset: 4,
			expectedLimit:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := &PaginationBase{
				Page:     tc.page,
				PageSize: tc.pageSize,
			}

			offset, limit := p.ToOffsetLimit()

			assert.Equal(t, tc.expectedOffset, offset)
			assert.Equal(t, tc.expectedLimit, limit)
		})
	}
}

func TestPaginationQuery_EmbeddedMethods(t *testing.T) {
	t.Parallel()

	t.Run("PaginationQuery can use PaginationBase methods", func(t *testing.T) {
		t.Parallel()

		p := &PaginationQuery{
			PaginationBase: PaginationBase{
				Page:     2,
				PageSize: 20,
			},
		}

		page, pageSize := p.ValidateAndNormalize()
		assert.Equal(t, 2, page)
		assert.Equal(t, 20, pageSize)

		offset, limit := p.ToOffsetLimit()
		assert.Equal(t, 20, offset) // (2-1) * 20
		assert.Equal(t, 20, limit)
	})
}

func TestPaginationWithSort_GetSortOrder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		sortOrder      string
		expectedResult string
	}{
		{
			name:           "success - lowercase asc",
			sortOrder:      "asc",
			expectedResult: "asc",
		},
		{
			name:           "success - uppercase ASC",
			sortOrder:      "ASC",
			expectedResult: "asc",
		},
		{
			name:           "success - lowercase desc",
			sortOrder:      "desc",
			expectedResult: "desc",
		},
		{
			name:           "success - uppercase DESC",
			sortOrder:      "DESC",
			expectedResult: "desc",
		},
		{
			name:           "success - empty string defaults to asc",
			sortOrder:      "",
			expectedResult: "asc",
		},
		{
			name:           "success - invalid value defaults to asc",
			sortOrder:      "invalid",
			expectedResult: "asc",
		},
		{
			name:           "success - mixed case asc",
			sortOrder:      "AsC",
			expectedResult: "asc",
		},
		{
			name:           "success - mixed case desc",
			sortOrder:      "DeSc",
			expectedResult: "asc", // Only exact matches work
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := &PaginationWithSort{
				PaginationBase: PaginationBase{
					Page:     1,
					PageSize: 10,
				},
				SortOrder: tc.sortOrder,
			}

			result := p.GetSortOrder()
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestPaginationWithSort_EmbeddedMethods(t *testing.T) {
	t.Parallel()

	t.Run("PaginationWithSort can use PaginationBase methods", func(t *testing.T) {
		t.Parallel()

		p := &PaginationWithSort{
			PaginationBase: PaginationBase{
				Page:     3,
				PageSize: 15,
			},
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		page, pageSize := p.ValidateAndNormalize()
		assert.Equal(t, 3, page)
		assert.Equal(t, 15, pageSize)

		offset, limit := p.ToOffsetLimit()
		assert.Equal(t, 30, offset) // (3-1) * 15
		assert.Equal(t, 15, limit)

		sortOrder := p.GetSortOrder()
		assert.Equal(t, "desc", sortOrder)
	})
}

func TestPaginationBase_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("very large page number", func(t *testing.T) {
		t.Parallel()

		p := &PaginationBase{
			Page:     10000,
			PageSize: 10,
		}

		offset, limit := p.ToOffsetLimit()
		assert.Equal(t, 99990, offset) // (10000-1) * 10
		assert.Equal(t, 10, limit)
	})

	t.Run("pageSize at boundary values", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name         string
			pageSize     int
			expectedSize int
		}{
			{
				name:         "pageSize = 1",
				pageSize:     1,
				expectedSize: 1,
			},
			{
				name:         "pageSize = MaxPageSize",
				pageSize:     MaxPageSize,
				expectedSize: MaxPageSize,
			},
			{
				name:         "pageSize > MaxPageSize",
				pageSize:     MaxPageSize + 10,
				expectedSize: MaxPageSize,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := &PaginationBase{
					Page:     1,
					PageSize: tc.pageSize,
				}

				_, pageSize := p.ValidateAndNormalize()
				assert.Equal(t, tc.expectedSize, pageSize)
			})
		}
	})
}
