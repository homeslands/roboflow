package step

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type GetStepByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetStepByIDQuery) Validate() error {
	return validator.Validate(q)
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(status|started_at|completed_at)$")
)

type ListStepQuery struct {
	WorkflowExecutionID uuid.UUID `validate:"required,uuid"`
	Sorts               []xsort.Sort
}

func (q ListStepQuery) Validate() error {
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
