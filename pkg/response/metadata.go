// Package response provides standardized API response structures and utilities.
// This file contains pagination metadata structures for list endpoints.
package response

type PaginationMetadata struct {
	// Page is the current page number (1-indexed).
	Page int `json:"page"`
	// PageSize is the number of items per page, exposed to clients as "limit".
	PageSize int `json:"limit"`
	// Total is the total number of records available.
	Total int64 `json:"total"`
}

// NewPaginationMetadata creates a new PaginationMetadata instance from the provided parameters.
// It calculates the total pages by dividing the total count by page size and rounding up.
//
// Parameters:
//   - page: The current page number (should be normalized/validated before calling)
//   - pageSize: The number of items per page (should be normalized/validated before calling)
//   - total: The total number of records available
//
// Returns:
//   - PaginationMetadata: A new metadata instance with calculated totalPages
//
// Example:
//
//	metadata := NewPaginationMetadata(2, 10, 25)
//	// Result: Page=2, PageSize=10, Total=25, TotalPages=3
func NewPaginationMetadata(page, pageSize int, total int64) PaginationMetadata {
	return PaginationMetadata{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}
