package converter

import (
	"encoding/json"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
)

func ToStepExecutionResponse(m stepexecution.StepExecution) (gen.StepExecutionResponse, error) {
	node, err := json.Marshal(m.Node)
	if err != nil {
		return gen.StepExecutionResponse{}, fmt.Errorf("marshal node: %w", err)
	}

	return gen.StepExecutionResponse{
		Id:                  m.ID,
		WorkflowExecutionId: m.WorkflowExecutionID,
		Status:              string(m.Status),
		Node:                node,
		Inputs:              m.Inputs,
		Outputs:             m.Outputs,
		Error:               m.Error,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}, nil
}
