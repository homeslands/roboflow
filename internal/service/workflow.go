package service

import (
	"context"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type GetWorkflowParams struct {
	ID string `validate:"required,uuid"`
}

type ListWorkflowsParams struct {
	PagingParams paging.Params `validate:"required"`
	Sorts        []sort.Sort   `validate:"sort=name is_draft created_at updated_at"`
}

type CreateWorkflowParams struct {
	Name        string        `validate:"required,alphanumspace,min=1,max=100"`
	Description *string       `validate:"omitempty,min=1,max=100"`
	Data        workflow.Data `validate:"required"`
}

type UpdateWorkflowParams struct {
	ID             string `validate:"required,uuid"`
	Name           string `validate:"required_if=SetName true,omitempty,alphanumspace,min=1,max=100"`
	SetName        bool
	Description    *string `validate:"required_if=SetDescription true,omitempty,min=1,max=100"`
	SetDescription bool
	IsDraft        bool `validate:"required_if=SetIsDraft true"`
	SetIsDraft     bool
	Data           workflow.Data `validate:"required_if=SetData true,omitempty"`
	SetData        bool
}

type DeleteWorkflowParams struct {
	ID string `validate:"required,uuid"`
}

type RunWorkflowParams struct {
	ID               string         `validate:"required,uuid"`
	RuntimeVariables map[string]any `validate:"required,dive"`
}

type WorkflowService interface {
	// GetWorkflow gets a workflow by its ID.
	GetWorkflow(ctx context.Context, params GetWorkflowParams) (workflow.Workflow, error)

	// ListWorkflows lists all workflows.
	ListWorkflows(ctx context.Context, params ListWorkflowsParams) (paging.List[workflow.Workflow], error)

	// CreateWorkflow creates a new workflow.
	CreateWorkflow(ctx context.Context, params CreateWorkflowParams) (workflow.Workflow, error)

	// UpdateWorkflow updates a workflow.
	UpdateWorkflow(ctx context.Context, params UpdateWorkflowParams) (workflow.Workflow, error)

	// DeleteWorkflow deletes a workflow.
	DeleteWorkflow(ctx context.Context, params DeleteWorkflowParams) error

	// RunWorkflow runs a workflow.
	RunWorkflow(ctx context.Context, params RunWorkflowParams) (workflowExecutionID string, err error)
}
