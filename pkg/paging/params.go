package paging

const (
	defaultPage     = 1
	defaultPageSize = 10
)

// Params represents the parameters for paging.
type Params struct {
	PageSize uint `validate:"required,min=1"`
	Page     uint `validate:"required,min=1"`
}

// NewParams creates a new Params instance with page size and page.
// If page size is nil or less than 0, it will use defaultPageSize(1).
// If page is nil, it will use defaultPage(10).
func NewParams(pageSize, page *uint, options ...ParamsOption) Params {
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

func (p Params) Offset() uint {
	return (p.Page - 1) * p.PageSize
}

func (p Params) Limit() uint {
	return p.PageSize
}

type ParamsOption func(*Params)

func WithDefaultPageSize(size uint) ParamsOption {
	return func(p *Params) {
		p.PageSize = size
	}
}

func WithDefaultPage(page uint) ParamsOption {
	return func(p *Params) {
		p.Page = page
	}
}

func WithMaxPageSize(size uint) ParamsOption {
	return func(p *Params) {
		if p.PageSize > size {
			p.PageSize = size
		}
	}
}
