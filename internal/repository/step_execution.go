package repository

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
)

type UpdateStepExecutionParams struct {
	ID             string
	Status         stepexecution.Status
	SetStatus      bool
	Inputs         map[string]any
	SetInputs      bool
	Outputs        map[string]any
	SetOutputs     bool
	Error          *string
	SetError       bool
	StartedAt      *time.Time
	SetStartedAt   bool
	CompletedAt    *time.Time
	SetCompletedAt bool
}
type StepExecutionRepository interface {
	// GetStepExecution gets a StepExecution by ID.
	GetStepExecution(ctx context.Context, db sqldb.SQLDB, id string) (stepexecution.StepExecution, error)

	// ListStepsByWorkflowExecutionID lists all Steps by WorkflowExecution ID.
	ListStepsByWorkflowExecutionID(ctx context.Context, db sqldb.SQLDB, workflowExecutionID string) ([]stepexecution.StepExecution, error)

	// BatchCreateStepExecutions creates multiple Steps.
	BatchCreateStepExecutions(ctx context.Context, db sqldb.SQLDB, steps []stepexecution.StepExecution) error

	// UpdateStepExecution updates a Step.
	UpdateStepExecution(ctx context.Context, db sqldb.SQLDB, params UpdateStepExecutionParams) (stepexecution.StepExecution, error)
}
