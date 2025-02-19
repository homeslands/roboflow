package converter

import (
	"encoding/json"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
)

func ToWorkflowExecutionResponse(m workflowexecution.WorkflowExecution) (gen.WorkflowExecutionResponse, error) {
	data, err := json.Marshal(m.Data)
	if err != nil {
		return gen.WorkflowExecutionResponse{}, fmt.Errorf("marshal workflow execution data: %w", err)
	}

	return gen.WorkflowExecutionResponse{
		Id:          m.ID,
		WorkflowId:  m.WorkflowID,
		Status:      string(m.Status),
		Data:        data,
		Inputs:      m.Inputs,
		Outputs:     m.Outputs,
		Error:       m.Error,
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}
