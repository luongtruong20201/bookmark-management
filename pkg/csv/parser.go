// Package csv provides utilities for parsing CSV content, including reading
// from multipart form file uploads.
package csv

import (
	"mime/multipart"

	"github.com/gocarina/gocsv"
)

// ParseFromMultipartFile opens the uploaded file from src, unmarshals its CSV
// content into data using gocsv, then closes the file. data must be a pointer to
// a struct or slice of structs with csv struct tags defining column mapping.
// Returns any error from opening the file or from CSV unmarshaling.
func ParseFromMultipartFile(src *multipart.FileHeader, data any) error {
	file, err := src.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.Unmarshal(file, data); err != nil {
		return err
	}

	return nil
}
