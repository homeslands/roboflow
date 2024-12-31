package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type WorkflowExecutionStepStatus string

const (
	WorkflowExecutionStepStatusPending   WorkflowExecutionStepStatus = "PENDING"
	WorkflowExecutionStepStatusRunning   WorkflowExecutionStepStatus = "RUNNING"
	WorkflowExecutionStepStatusCompleted WorkflowExecutionStepStatus = "COMPLETED"
	WorkflowExecutionStepStatusFailed    WorkflowExecutionStepStatus = "FAILED"
)

func (s WorkflowExecutionStepStatus) Validate() error {
	switch s {
	case WorkflowExecutionStepStatusPending:
	case WorkflowExecutionStepStatusRunning:
	case WorkflowExecutionStepStatusCompleted:
	case WorkflowExecutionStepStatusFailed:
	default:
		return xerrors.ThrowInvalidArgument(nil, "invalid step status")
	}
	return nil
}

type Step struct {
	ID                  uuid.UUID
	WorkflowExecutionID uuid.UUID
	Env                 map[string]string
	Node                WorkflowNode
	Status              WorkflowExecutionStepStatus
	StartedAt           *time.Time
	CompletedAt         *time.Time
}

type StepRepository interface {
	Get(ctx context.Context, id uuid.UUID) (Step, error)
	List(ctx context.Context, workflowExecutionID uuid.UUID, sorts []xsort.Sort) ([]Step, error)
	Update(ctx context.Context, step Step) error
}
