package workflow

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

type CreateWorkflowCommand struct {
	Name        string                   `validate:"required,alphanumspace,min=1,max=100"`
	Description *string                  `validate:"omitempty,max=1000"`
	Definition  model.WorkflowDefinition `validate:"required"`
}

func (c CreateWorkflowCommand) Validate() error {
	if err := validator.Validate(c); err != nil {
		return err
	}

	for _, node := range c.Definition.Nodes {
		if err := validateNode(node); err != nil {
			return err
		}
	}

	return nil
}

type UpdateWorkflowCommand struct {
	ID          uuid.UUID                `validate:"required,uuid"`
	Name        string                   `validate:"required,alphanumspace,min=1,max=100"`
	Description *string                  `validate:"omitempty,max=1000"`
	Definition  model.WorkflowDefinition `validate:"required"`
}

func (c UpdateWorkflowCommand) Validate() error {
	if err := validator.Validate(c); err != nil {
		return err
	}

	for _, node := range c.Definition.Nodes {
		if err := validateNode(node); err != nil {
			return err
		}
	}

	return nil
}

type DeleteWorkflowCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteWorkflowCommand) Validate() error {
	return validator.Validate(c)
}

type RunWorkflowCommand struct {
	ID  uuid.UUID         `validate:"required,uuid"`
	Env map[string]string `validate:"omitempty"`
}

func (c RunWorkflowCommand) Validate() error {
	return validator.Validate(c)
}

type GetWorkflowByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetWorkflowByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListWorkflowQuery struct {
	PagingParams paging.Params
	Sorts        []xsort.Sort
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(name|created_at|updated_at)$")
)

func (q ListWorkflowQuery) Validate() error {
	for _, sort := range q.Sorts {
		if !allowedSortFieldsRegexp.MatchString(sort.Col) {
			return xerrors.ThrowInvalidArgument(nil, fmt.Sprintf("invalid sort field: %s", sort.Col))
		}
	}

	return nil
}

func validateNode(node model.WorkflowNode) error {
	if err := node.Type.Validate(); err != nil {
		return err
	}

	if err := node.Definition.Type.Validate(); err != nil {
		return err
	}

	return nil
}
