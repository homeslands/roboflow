package raybotcommand

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type GetRaybotCommandByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetRaybotCommandByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListRaybotCommandQuery struct {
	RaybotID uuid.UUID `validate:"required,uuid"`

	PagingParams paging.Params
	Sorts        []xsort.Sort
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(status|created_at|completed_at)$")
)

func (q ListRaybotCommandQuery) Validate() error {
	if err := validator.Validate(q); err != nil {
		return err
	}

	for _, sort := range q.Sorts {
		if !allowedSortFieldsRegexp.MatchString(sort.Col) {
			return xerrors.ThrowInvalidArgument(nil, fmt.Sprintf("invalid sort field: %s", sort.Col))
		}
	}

	return nil
}
