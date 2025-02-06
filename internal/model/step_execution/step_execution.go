package stepexecution

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
)

type Status string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *Status) UnmarshalText(text []byte) error {
	stt := Status(text)
	if _, ok := StatusMap[stt]; !ok {
		return fmt.Errorf("invalid StepStatus: %s", text)
	}
	*s = stt
	return nil
}

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

var StatusMap = map[Status]struct{}{
	StatusPending:   {},
	StatusRunning:   {},
	StatusCompleted: {},
	StatusFailed:    {},
}

type StepExecution struct {
	ID                  string
	WorkflowExecutionID string
	Status              Status
	Node                node.Node
	Inputs              map[string]any
	Outputs             map[string]any
	Error               *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	StartedAt           *time.Time
	CompletedAt         *time.Time
}

func NewStepExecution(workflowExecutionID string, node node.Node, inputs map[string]any) StepExecution {
	now := time.Now()
	return StepExecution{
		ID:                  uuid.NewString(),
		WorkflowExecutionID: workflowExecutionID,
		Status:              StatusPending,
		Node:                node,
		Inputs:              inputs,
		Outputs:             make(map[string]any),
		Error:               nil,
		StartedAt:           nil,
		CompletedAt:         nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
