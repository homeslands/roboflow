package pagination

import (
	"net/http"
	"strconv"
)

type PaginationOption func(options *paginationOptions)

type paginationOptions struct {
	defaultPage     int
	defaultPageSize int
	minPage         int
	maxPage         int
	minPageSize     int
	maxPageSize     int
}

func WithDefaultPage(page int) PaginationOption {
	return func(opts *paginationOptions) {
		opts.defaultPage = page
	}
}

func WithDefaultPageSize(pageSize int) PaginationOption {
	return func(opts *paginationOptions) {
		opts.defaultPageSize = pageSize
	}
}

func WithPageRange(min, max int) PaginationOption {
	return func(opts *paginationOptions) {
		opts.minPage = min
		opts.maxPage = max
	}
}

func WithPageSizeRange(min, max int) PaginationOption {
	return func(opts *paginationOptions) {
		opts.minPageSize = min
		opts.maxPageSize = max
	}
}

func GetPage(r *http.Request, opts ...PaginationOption) int {
	options := &paginationOptions{
		defaultPage: 1,
		minPage:     1,
		maxPage:     65535,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(options)
	}

	// Get page from query
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return options.defaultPage
	}

	//	Check page range
	pageInt := int(page)
	if options.minPage > 0 && pageInt < options.minPage {
		return options.minPage
	}
	if options.maxPage > 0 && pageInt > options.maxPage {
		return options.maxPage
	}

	return page
}

func GetPageSize(r *http.Request, opts ...PaginationOption) int {
	options := &paginationOptions{
		defaultPageSize: 10,
		minPageSize:     1,
		maxPageSize:     100,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(options)
	}

	// Get page size from query
	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return options.defaultPageSize
	}

	// Check page size range
	pageSizeInt := int(pageSize)
	if options.minPageSize > 0 && pageSizeInt < options.minPageSize {
		return options.minPageSize
	}
	if options.maxPageSize > 0 && pageSizeInt > options.maxPageSize {
		return options.maxPageSize
	}

	return pageSize
}
