package jwt

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		keyPath        string
		input          string
		expectedOutput jwt.MapClaims
		expectedErrStr string
	}{
		{
			name:           "invalid key path",
			keyPath:        "in-valid.pem",
			input:          "",
			expectedOutput: nil,
			expectedErrStr: "no such file or directory",
		},
		{
			name:    "success",
			keyPath: filepath.FromSlash("./public.pem"),
			input:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiSFkiLCJuYW1lIjoidHJ1b25nbHEifQ.cQ45zW2G5WXuIKHTrGaRmZX4o75Vp7uQIaGqBKS7GTK12dp8HZdd6XoCoBYd27a3ebZSarndubCt6QD5TZgUMLN5iDelqzSITyyoUgjLLfFydoWvfeEs68BzZac0kXtHd0BCggge5PtY2SWP-XR-kWStmlfRzTau3InqnZNJU2tgB9Z28X4lve2rgiDPLQhZ6EzR_stj9bWLJQjIzVj0kabcL2w2N-1NhIuDiJmD9trxradvmOXqP4Ktk4iQSupts0Iitmtr8fuUIqB23uFG5dieS10hMZL3MurcED2WZRv7eJhF3PbgTzD18Aand-UcComP29E0O61ntj-c7Bgj3g",
			expectedOutput: jwt.MapClaims{
				"name":    "truonglq",
				"address": "HY",
			},
			expectedErrStr: "",
		},
		{
			name:           "invalid token",
			keyPath:        filepath.FromSlash("./public.pem"),
			input:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiSFkiLCJuYW1lIjoidHJ1b25nbHEifQ.cQ45zW2G5WXuIKHTrGaRmZX4o75Vp7uQIaGqBKS7GTK12dp8HZdd6XoCoBYd27a3ebZSarndubCt6QD5TZgUMLN5iDelqzSITyyoUgjLLfFydoWvfeEs68BzZac0kXtHd0BCggge5PtY2SWP-XR-kWStmlfRzTau3InqnZNJU2tgB9Z28X4lve2rgiDPLQhZ6EzR_stj9bWLJQjIzVj0kabcL2w2N-1NhIuDiJmD9trxradvmOXqP4Ktk4iQSupts0Iitmtr8fuUIqB23uFG5dieS10hMZL3MurcED2WZRv7eJhF3PbgTzD18Aand-UcComP29E0O61ntj-c7gghi",
			expectedOutput: nil,
			expectedErrStr: "invalid token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := NewJWTValidator(tc.keyPath)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}

			claims, err := validator.ValidateToken(tc.input)
			assert.Equal(t, claims, tc.expectedOutput)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedErrStr)
			}
		})
	}
}
