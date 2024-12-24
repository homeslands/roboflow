package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type Workflow struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Definition  *WorkflowDefinition
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkflowRepository interface {
	Get(ctx context.Context, id uuid.UUID) (Workflow, error)
	List(ctx context.Context, p paging.Params, sorts []xsort.Sort) (*paging.List[Workflow], error)
	Create(ctx context.Context, workflow Workflow) error
	Update(ctx context.Context, workflow Workflow) (Workflow, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
