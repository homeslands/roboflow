package handler

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
)

type workflowExecutionHandler struct {
	workflowExecutionSvc service.WorkflowExecutionService
}

func newWorkflowExecutionHandler(workflowExecutionSvc service.WorkflowExecutionService) *workflowExecutionHandler {
	return &workflowExecutionHandler{workflowExecutionSvc: workflowExecutionSvc}
}

func (h workflowExecutionHandler) WorkflowExecutionGet(ctx context.Context, request gen.WorkflowExecutionGetRequestObject) (gen.WorkflowExecutionGetResponseObject, error) {
	workflowExecution, err := h.workflowExecutionSvc.GetWorkflowExecution(ctx, service.GetWorkflowExecutionParams{
		ID: request.WorkflowExecutionId,
	})
	if err != nil {
		return nil, fmt.Errorf("get workflow execution: %w", err)
	}

	res, err := converter.ToWorkflowExecutionResponse(workflowExecution)
	if err != nil {
		return nil, fmt.Errorf("convert to workflow execution response: %w", err)
	}

	return gen.WorkflowExecutionGet200JSONResponse(res), nil
}

func (h workflowExecutionHandler) WorkflowExecutionList(ctx context.Context, request gen.WorkflowExecutionListRequestObject) (gen.WorkflowExecutionListResponseObject, error) {
	workflowExecutions, err := h.workflowExecutionSvc.ListWorkflowExecutionsByWorkflowID(ctx, service.ListWorkflowExecutionsByWorkflowIDParams{
		WorkflowID: request.WorkflowId,
	})
	if err != nil {
		return nil, fmt.Errorf("list workflow executions by workflow id: %w", err)
	}

	items := make([]gen.WorkflowExecutionResponse, len(workflowExecutions.Items))
	for i, item := range workflowExecutions.Items {
		res, err := converter.ToWorkflowExecutionResponse(item)
		if err != nil {
			return nil, fmt.Errorf("convert to workflow execution response: %w", err)
		}

		items[i] = res
	}

	return gen.WorkflowExecutionList200JSONResponse{
		Items:      items,
		TotalItems: workflowExecutions.TotalItems,
	}, nil
}
