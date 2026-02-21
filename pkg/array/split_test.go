package array

import "testing"

func TestSplitIntoChunks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		input          []int
		size           int
		expectedChunks [][]int
	}{
		{
			name:           "empty slice",
			input:          []int{},
			size:           3,
			expectedChunks: [][]int{{}},
		},
		{
			name:           "size greater than length",
			input:          []int{1, 2, 3},
			size:           5,
			expectedChunks: [][]int{{1, 2, 3}},
		},
		{
			name:           "exact multiple of size",
			input:          []int{1, 2, 3, 4, 5, 6},
			size:           2,
			expectedChunks: [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:           "with remainder",
			input:          []int{1, 2, 3, 4, 5},
			size:           2,
			expectedChunks: [][]int{{1, 2}, {3, 4}, {5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			chunks := SplitIntoChunks(tc.input, tc.size)

			if len(chunks) != len(tc.expectedChunks) {
				t.Fatalf("expected %d chunks, got %d", len(tc.expectedChunks), len(chunks))
			}

			for i := range tc.expectedChunks {
				if len(chunks[i]) != len(tc.expectedChunks[i]) {
					t.Fatalf("expected chunk %d length %d, got %d", i, len(tc.expectedChunks[i]), len(chunks[i]))
				}
				for j := range tc.expectedChunks[i] {
					if chunks[i][j] != tc.expectedChunks[i][j] {
						t.Fatalf("expected chunks[%d][%d]=%d, got %d", i, j, tc.expectedChunks[i][j], chunks[i][j])
					}
				}
			}
		})
	}
}
