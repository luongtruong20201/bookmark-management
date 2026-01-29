package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name: "error",
			err:  errors.New("ok"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.err == nil {
				assert.NotPanics(t, func() {
					HandleError(tc.err)
				})
			} else {
				assert.PanicsWithError(t, "ok", func() {
					HandleError(tc.err)
				})
			}
		})
	}
}
