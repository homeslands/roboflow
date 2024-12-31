package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type WorkflowExecutionStatus string

const (
	WorkflowExecutionStatusPending   WorkflowExecutionStatus = "PENDING"
	WorkflowExecutionStatusRunning   WorkflowExecutionStatus = "RUNNING"
	WorkflowExecutionStatusCompleted WorkflowExecutionStatus = "COMPLETED"
	WorkflowExecutionStatusFailed    WorkflowExecutionStatus = "FAILED"
	WorkflowExecutionStatusCanceled  WorkflowExecutionStatus = "CANCELED"
)

func (s WorkflowExecutionStatus) Validate() error {
	switch s {
	case WorkflowExecutionStatusPending:
	case WorkflowExecutionStatusRunning:
	case WorkflowExecutionStatusCompleted:
	case WorkflowExecutionStatusFailed:
	case WorkflowExecutionStatusCanceled:
	default:
		return xerrors.ThrowInvalidArgument(nil, "invalid workflow execution status")
	}
	return nil
}

type WorkflowExecution struct {
	ID          uuid.UUID
	WorkflowID  uuid.UUID
	Status      WorkflowExecutionStatus
	Env         map[string]string
	Definition  WorkflowDefinition
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Steps       *[]Step
}

type WorkflowExecutionRepository interface {
	Get(ctx context.Context, id uuid.UUID) (WorkflowExecution, error)
	GetStatus(ctx context.Context, id uuid.UUID) (WorkflowExecutionStatus, error)
	List(ctx context.Context, workflowID uuid.UUID, p paging.Params, sorts []xsort.Sort) (*paging.List[WorkflowExecution], error)
	Create(ctx context.Context, workflowExecution WorkflowExecution) error
	Update(
		ctx context.Context,
		id uuid.UUID,
		fn func(wfe *WorkflowExecution) error,
	) error
	Delete(ctx context.Context, id uuid.UUID) error
}
