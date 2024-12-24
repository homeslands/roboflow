package workflowexecution

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type UpdateWorkflowExecutionCommand struct {
	ID          uuid.UUID                     `validate:"required,uuid"`
	Status      model.WorkflowExecutionStatus `validate:"required"`
	StartedAt   *time.Time                    `validate:"omitempty"`
	CompletedAt *time.Time                    `validate:"omitempty"`
}

func (c UpdateWorkflowExecutionCommand) Validate() error {
	if err := validator.Validate(c); err != nil {
		return err
	}

	return c.Status.Validate()
}

type DeleteWorkflowExecutionCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteWorkflowExecutionCommand) Validate() error {
	return validator.Validate(c)
}

type GetWorkflowExecutionByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetWorkflowExecutionByIDQuery) Validate() error {
	return validator.Validate(q)
}

type GetWorkflowExecutionStatusByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetWorkflowExecutionStatusByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListWorkflowExecutionQuery struct {
	WorkflowID   uuid.UUID `validate:"required,uuid"`
	PagingParams paging.Params
	Sorts        []xsort.Sort
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(status|created_at|started_at|completed_at)$")
)

func (q ListWorkflowExecutionQuery) Validate() error {
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
