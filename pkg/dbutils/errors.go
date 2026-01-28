// Package dbutils provides helpers for normalizing and classifying database errors.
// It translates low-level driver or GORM errors into consistent application error types.
package dbutils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var errorFilters = []func(err error) (bool, error){
	filterDuplicationType,
	filterRecordNotFound,
}

// CatchDBErr inspects a low-level database error and attempts to map it to a
// well-defined application error (for example duplication or not-found).
// If no mapping rule matches, the original error is returned unchanged.
func CatchDBErr(err error) error {
	if err == nil {
		return nil
	}

	for _, filter := range errorFilters {
		check, newErr := filter(err)
		if check {
			return newErr
		}
	}

	return err
}

var (
	ErrDuplicationType = errors.New("duplicate type")
	ErrNotFoundType    = errors.New("not found type")
)

// filterDuplicationType detects unique-constraint violations and maps them
// to ErrDuplicationType.
func filterDuplicationType(err error) (bool, error) {
	return strings.Contains(err.Error(), "unique constraint"), ErrDuplicationType
}

// filterRecordNotFound detects GORM's ErrRecordNotFound and maps it
// to ErrNotFoundType.
func filterRecordNotFound(err error) (bool, error) {
	return errors.Is(err, gorm.ErrRecordNotFound), ErrNotFoundType
}
