package service

import (
	"context"

	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type GetWorkflowExecutionParams struct {
	ID string `validate:"required,uuid"`
}

type ListWorkflowExecutionsByWorkflowIDParams struct {
	WorkflowID   string        `validate:"required,uuid"`
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=status started_at completed_at created_at updated_at"`
}

type ProcessRunWorkflowExecutionParams struct {
	WorkflowExecutionID string `validate:"required,uuid"`
}

type WorkflowExecutionService interface {
	// GetWorkflowExecution gets a WorkflowExecution by its ID.
	GetWorkflowExecution(ctx context.Context, params GetWorkflowExecutionParams) (workflowexecution.WorkflowExecution, error)

	// ListWorkflowExecutionsByWorkflowID lists all WorkflowExecutions by Workflow ID.
	ListWorkflowExecutionsByWorkflowID(ctx context.Context, params ListWorkflowExecutionsByWorkflowIDParams) (paging.List[workflowexecution.WorkflowExecution], error)

	// ProcessRunWorkflowExecution processes a run WorkflowExecution.
	ProcessRunWorkflowExecution(ctx context.Context, params ProcessRunWorkflowExecutionParams) error
}
