package converter

import (
	"encoding/json"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
)

func ToWorkflowResponse(m workflow.Workflow) (gen.WorkflowResponse, error) {
	data, err := json.Marshal(m.Data)
	if err != nil {
		return gen.WorkflowResponse{}, fmt.Errorf("json marshal workflow data: %w", err)
	}

	return gen.WorkflowResponse{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		IsDraft:     m.IsDraft,
		Data:        data,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

func ToWorkflowItemListResponse(m workflow.Workflow) gen.WorkflowItemListResponse {
	return gen.WorkflowItemListResponse{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		IsDraft:     m.IsDraft,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToWorkflowRunResponse(workflowExecutionID string) gen.RunWorkflowResponse {
	return gen.RunWorkflowResponse{
		Id: workflowExecutionID,
	}
}
