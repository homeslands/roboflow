package pagination

import "encoding/json"

type PagingResponse[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
}

func (p PagingResponse[T]) MarshalJSON() ([]byte, error) {
	if p.Items == nil {
		p.Items = []T{}
	}
	return json.Marshal(p)
}
