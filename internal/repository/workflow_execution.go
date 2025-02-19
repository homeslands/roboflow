package repository

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type UpdateWorkflowExecutionParams struct {
	ID             string
	Status         workflowexecution.Status
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

type WorkflowExecutionRepository interface {
	// GetWorkflowExecution gets a WorkflowExecution by its ID.
	GetWorkflowExecution(ctx context.Context, db sqldb.SQLDB, id string) (workflowexecution.WorkflowExecution, error)

	// ListWorkflowExecutionsByWorkflowID lists all WorkflowExecutions by Workflow ID.
	ListWorkflowExecutionsByWorkflowID(
		ctx context.Context,
		db sqldb.SQLDB,
		pagingParams paging.Params,
		sorts []sort.Sort,
		workflowID string,
	) (paging.List[workflowexecution.WorkflowExecution], error)

	// CreateWorkflowExecution creates a new WorkflowExecution.
	CreateWorkflowExecution(ctx context.Context, db sqldb.SQLDB, workflowExecution workflowexecution.WorkflowExecution) error

	// UpdateWorkflowExecution updates a WorkflowExecution.
	UpdateWorkflowExecution(ctx context.Context, db sqldb.SQLDB, params UpdateWorkflowExecutionParams) (workflowexecution.WorkflowExecution, error)
}
