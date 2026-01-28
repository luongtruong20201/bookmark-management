package jwt

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTGenerator_GenerateeToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		keyPath        string
		inputContent   jwt.MapClaims
		expectedOutput string
		expectedErr    error
	}{
		{
			name:    "valid key path",
			keyPath: filepath.FromSlash("./private_test.pem"),
			inputContent: jwt.MapClaims{
				"name":    "truonglq",
				"address": "HY",
			},
			expectedOutput: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiSFkiLCJuYW1lIjoidHJ1b25nbHEifQ.cQ45zW2G5WXuIKHTrGaRmZX4o75Vp7uQIaGqBKS7GTK12dp8HZdd6XoCoBYd27a3ebZSarndubCt6QD5TZgUMLN5iDelqzSITyyoUgjLLfFydoWvfeEs68BzZac0kXtHd0BCggge5PtY2SWP-XR-kWStmlfRzTau3InqnZNJU2tgB9Z28X4lve2rgiDPLQhZ6EzR_stj9bWLJQjIzVj0kabcL2w2N-1NhIuDiJmD9trxradvmOXqP4Ktk4iQSupts0Iitmtr8fuUIqB23uFG5dieS10hMZL3MurcED2WZRv7eJhF3PbgTzD18Aand-UcComP29E0O61ntj-c7Bgj3g",
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			generator, err := NewJWTGenerator(tc.keyPath)
			assert.Equal(t, err, tc.expectedErr)
			res, err := generator.GenerateToken(tc.inputContent)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
