package paging

type List[T any] struct {
	Items     []T
	TotalItem int64
}

// NewList creates a new List instance with total item.
func NewList[T any](items []T, totalItem int64) *List[T] {
	return &List[T]{
		Items:     items,
		TotalItem: totalItem,
	}
}
