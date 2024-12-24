package workflowexecution

import (
	"time"

	"github.com/google/uuid"
	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
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
