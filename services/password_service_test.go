package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService_GeneratePassword(t *testing.T) {
	testCases := []struct {
		name        string
		expectedLen int
		expectedErr error
	}{
		{name: "normal case",
			expectedLen: 10,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testSvc := NewPassword()
			pass, err := testSvc.GeneratePassword()
			assert.Equal(t, tc.expectedLen, len(pass))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
