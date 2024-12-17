package xsort

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const (
	OrderASC  = "ASC"
	OrderDESC = "DESC"
)

type Sort struct {
	col   string
	order string
}

// Attach attaches sort to builder.
func (s Sort) Attach(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(s.col + " " + s.order)
}

// NewList parses a comma-separated string into a slice of Sort objects.
// The input string must follow the pattern: "col1,-col2,...", where:
//   - Each item represents a column to sort by.
//   - A column prefixed with a minus sign ("-") indicates descending order.
//   - Columns without a prefix indicate ascending order.
func NewList(s *string) ([]Sort, error) {
	if s == nil || len(*s) == 0 {
		return nil, nil
	}

	orderBys := strings.Split(*s, ",")
	sorts := make([]Sort, len(orderBys))

	for i, r := range orderBys {
		orderBy := string(r)
		if strings.HasPrefix(orderBy, " ") || strings.HasSuffix(orderBy, " ") {
			return nil, fmt.Errorf("invalid sort column: %s", orderBy)
		}

		order := OrderASC
		if strings.HasPrefix(orderBy, "-") {
			order = OrderDESC
			orderBy = orderBy[1:]
		}

		sorts[i] = Sort{
			col:   orderBy,
			order: order,
		}
	}

	return sorts, nil
}
