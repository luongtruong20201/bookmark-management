// Package request provides helpers for binding and validating HTTP request inputs.
// This file contains pagination-related structures and utilities.
package request

const (
	// DefaultPage is the default page number when not provided or invalid.
	DefaultPage = 1
	// DefaultPageSize is the default page size when not provided or invalid.
	DefaultPageSize = 10
	// MaxPageSize is the maximum allowed page size.
	MaxPageSize = 100
)

// PaginationBase represents base pagination query parameters.
// It provides page-based pagination with page and pageSize parameters.
type PaginationBase struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"pageSize" binding:"omitempty,gte=1,lte=100"`
}

// ValidateAndNormalize validates and normalizes pagination parameters.
// It ensures page and pageSize are within valid ranges and applies defaults if needed.
// Returns the normalized page and pageSize values.
func (p *PaginationBase) ValidateAndNormalize() (page, pageSize int) {
	page = p.Page
	if page < 1 {
		page = DefaultPage
	}

	pageSize = p.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return page, pageSize
}

// ToOffsetLimit converts page and pageSize to database offset and limit values.
// It first validates and normalizes the pagination parameters, then calculates offset and limit.
// Returns offset and limit values suitable for database queries.
func (p *PaginationBase) ToOffsetLimit() (offset, limit int) {
	page, pageSize := p.ValidateAndNormalize()
	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, limit
}

// PaginationQuery represents pagination query parameters for list endpoints.
// It is an alias for PaginationBase to maintain backward compatibility.
type PaginationQuery struct {
	PaginationBase
}

// PaginationWithSort represents pagination query parameters with sorting support.
// It embeds PaginationBase and adds sort field and sort order.
type PaginationWithSort struct {
	PaginationBase
	SortBy    string `form:"sortBy" binding:"omitempty"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc ASC DESC"`
}

// GetSortOrder returns the normalized sort order as a lowercase string ("asc" or "desc").
// It handles case-insensitive input and defaults to "asc" if the sortOrder is empty or invalid.
//
// Returns:
//   - string: The normalized sort order, either "asc" or "desc"
//
// Example:
//
//	p := PaginationWithSort{SortOrder: "DESC"}
//	order := p.GetSortOrder() // returns "desc"
func (p *PaginationWithSort) GetSortOrder() string {
	order := p.SortOrder
	if order == "" {
		return "asc"
	}

	if order == "ASC" || order == "asc" {
		return "asc"
	}

	if order == "DESC" || order == "desc" {
		return "desc"
	}

	return "asc"
}
