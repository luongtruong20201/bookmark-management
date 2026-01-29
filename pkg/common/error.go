// Package common provides shared helper utilities used across the application.
// It contains small, generic helpers that don't belong to a specific domain.
package common

// HandleError panics if the provided error is non-nil.
// It is a convenience helper for quickly failing on unexpected errors.
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
