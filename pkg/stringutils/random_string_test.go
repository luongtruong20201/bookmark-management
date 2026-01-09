package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGenerator_GenerateCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedLen   int
		expectedError error
	}{
		{
			name:          "success",
			expectedLen:   10,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keyGen := NewKeyGen()
			code, err := keyGen.GenerateCode(10)
			assert.Equal(t, err, tc.expectedError)
			assert.Equal(t, len(code), 10)
		})
	}
}
