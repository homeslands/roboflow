package raybot

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type GetRaybotByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetRaybotByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListRaybotQuery struct {
	PagingParams paging.Params
	Sorts        []xsort.Sort
	Status       *model.RaybotStatus
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(name|status|last_connected_at|created_at|updated_at)$")
)

func (q ListRaybotQuery) Validate() error {
	for _, sort := range q.Sorts {
		if !allowedSortFieldsRegexp.MatchString(sort.Col) {
			return xerrors.ThrowInvalidArgument(nil, fmt.Sprintf("invalid sort field: %s", sort.Col))
		}
	}
	if q.Status != nil {
		if err := q.Status.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type GetStatusQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetStatusQuery) Validate() error {
	return validator.Validate(q)
}
