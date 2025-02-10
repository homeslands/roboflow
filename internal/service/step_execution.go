package service

import (
	"context"
	"time"

	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
)

type GetStepExecutionParams struct {
	ID string `validate:"required,uuid"`
}

type ListStepsByWorkflowExecutionIDParams struct {
	WorkflowExecutionID string `validate:"required,uuid"`
}

type UpdateStepExecutionParams struct {
	ID             string               `validate:"required,uuid"`
	Status         stepexecution.Status `validate:"required_if=SetStatus true,omitempty,enum"`
	SetStatus      bool
	Inputs         map[string]any `validate:"required_if=SetInputs true"`
	SetInputs      bool
	Outputs        map[string]any `validate:"required_if=SetOutputs true"`
	SetOutputs     bool
	Error          *string `validate:"required_if=SetError true,omitempty,min=1,max=100"`
	SetError       bool
	StartedAt      *time.Time
	SetStartedAt   bool
	CompletedAt    *time.Time
	SetCompletedAt bool
}

type StepExecutionService interface {
	// GetStepExecution gets a StepExecution by ID.
	GetStepExecution(ctx context.Context, params GetStepExecutionParams) (stepexecution.StepExecution, error)

	// ListStepsByWorkflowExecutionID lists all Steps by WorkflowExecution ID.
	ListStepsByWorkflowExecutionID(ctx context.Context, params ListStepsByWorkflowExecutionIDParams) ([]stepexecution.StepExecution, error)

	// UpdateStepExecution updates a StepExecution.
	UpdateStepExecution(ctx context.Context, params UpdateStepExecutionParams) (stepexecution.StepExecution, error)
}
