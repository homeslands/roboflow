package repository

import (
	"context"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type UpdateWorkflowParams struct {
	ID             string
	Name           string
	SetName        bool
	Description    *string
	SetDescription bool
	IsDraft        bool
	SetIsDraft     bool
	IsValid        bool
	SetIsValid     bool
	Data           workflow.Data
	SetData        bool
}

type WorkflowRepository interface {
	// GetWorkflow gets a Workflow by its ID.
	GetWorkflow(ctx context.Context, db sqldb.SQLDB, id string) (workflow.Workflow, error)

	// ListWorkflows lists all Workflows.
	ListWorkflows(ctx context.Context, db sqldb.SQLDB, pagingParams paging.Params, sorts []sort.Sort) (paging.List[workflow.Workflow], error)

	// CreateWorkflow creates a new Workflow.
	CreateWorkflow(ctx context.Context, db sqldb.SQLDB, workflow workflow.Workflow) error

	// UpdateWorkflow updates a Workflow.
	UpdateWorkflow(ctx context.Context, db sqldb.SQLDB, params UpdateWorkflowParams) (workflow.Workflow, error)

	// DeleteWorkflow deletes a Workflow.
	DeleteWorkflow(ctx context.Context, db sqldb.SQLDB, id string) error
}
