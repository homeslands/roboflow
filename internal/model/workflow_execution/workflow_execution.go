package workflowexecution

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
)

type Status string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *Status) UnmarshalText(text []byte) error {
	workflowExecutionStatus := Status(text)
	if _, ok := StatusMap[workflowExecutionStatus]; !ok {
		return fmt.Errorf("invalid WorkflowExecutionStatus: %s", text)
	}
	*s = workflowExecutionStatus
	return nil
}

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
	StatusCancelled Status = "CANCELLED"
)

var StatusMap = map[Status]struct{}{
	StatusPending:   {},
	StatusRunning:   {},
	StatusCompleted: {},
	StatusFailed:    {},
	StatusCancelled: {},
}

type WorkflowExecution struct {
	ID          string
	WorkflowID  string
	Status      Status
	Data        workflow.Data
	Inputs      map[string]any
	Outputs     map[string]any
	Error       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

func NewWorkflowExecution(workflowID string, data workflow.Data, inputs map[string]any) WorkflowExecution {
	now := time.Now()
	return WorkflowExecution{
		ID:          uuid.NewString(),
		WorkflowID:  workflowID,
		Status:      StatusPending,
		Data:        data,
		Inputs:      inputs,
		Outputs:     make(map[string]any),
		Error:       nil,
		StartedAt:   nil,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
