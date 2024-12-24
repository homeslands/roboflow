package qrlocation

import (
	"regexp"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type GetQrLocationByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetQrLocationByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListQrLocationQuery struct {
	PagingParams paging.Params
	Sorts        []xsort.Sort
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(name|qr_code|created_at|updated_at)$")
)

func (q ListQrLocationQuery) Validate() error {
	for _, sort := range q.Sorts {
		if !allowedSortFieldsRegexp.MatchString(sort.Col) {
			return xerrors.ThrowInvalidArgument(nil, "invalid sort field")
		}
	}
	return nil
}
