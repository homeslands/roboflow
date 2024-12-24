package paging

const (
	defaultPage     = int32(1)
	defaultPageSize = int32(10)
)

type Params struct {
	PageSize int32
	Page     int32
}

// NewParams creates a new Params instance with page size and page.
// If page size is nil or less than 0, it will use defaultPageSize(1).
// If page is nil, it will use defaultPage(10).
func NewParams(pageSize, page *int32, options ...ParamsOption) Params {
	p := Params{
		PageSize: defaultPageSize,
		Page:     defaultPage,
	}
	if pageSize != nil && *pageSize > 0 {
		p.PageSize = *pageSize
	}
	if page != nil && *page > 0 {
		p.Page = *page
	}

	for _, option := range options {
		option(&p)
	}

	return p
}

func (p Params) Offset() int32 {
	return (p.Page - 1) * p.PageSize
}

func (p Params) Limit() int32 {
	return p.PageSize
}

type ParamsOption func(*Params)

func WithDefaultPageSize(size int32) ParamsOption {
	return func(p *Params) {
		p.PageSize = size
	}
}

func WithDefaultPage(page int32) ParamsOption {
	return func(p *Params) {
		p.Page = page
	}
}

func WithMaxPageSize(size int32) ParamsOption {
	return func(p *Params) {
		if p.PageSize > size {
			p.PageSize = size
		}
	}
}
