package event

import (
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

const (
	TopicWorkflowExecutionRun = "workflow.execution.run"
)

type WorkflowExecutionRun struct {
	WorkflowExecutionID uuid.UUID
	Steps               []Step
}

type Step struct {
	ID     uuid.UUID
	Type   model.TaskType
	Fields map[string]string
}
