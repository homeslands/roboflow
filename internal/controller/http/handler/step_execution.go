package handler

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
)

type stepExecutionHandler struct {
	stepExecutionSvc service.StepExecutionService
}

func newStepExecutionHandler(stepExecutionSvc service.StepExecutionService) *stepExecutionHandler {
	return &stepExecutionHandler{stepExecutionSvc: stepExecutionSvc}
}

func (h stepExecutionHandler) StepExecutionGet(ctx context.Context, request gen.StepExecutionGetRequestObject) (gen.StepExecutionGetResponseObject, error) {
	stepExecution, err := h.stepExecutionSvc.GetStepExecution(ctx, service.GetStepExecutionParams{
		ID: request.StepExecutionId,
	})
	if err != nil {
		return nil, fmt.Errorf("get step execution: %w", err)
	}

	res, err := converter.ToStepExecutionResponse(stepExecution)
	if err != nil {
		return nil, fmt.Errorf("convert to step execution response: %w", err)
	}

	return gen.StepExecutionGet200JSONResponse(res), nil
}

//nolint:revive
func (h stepExecutionHandler) StepExecutionListByWorkflowExecutionId(ctx context.Context, request gen.StepExecutionListByWorkflowExecutionIdRequestObject) (gen.StepExecutionListByWorkflowExecutionIdResponseObject, error) {
	stepExecutions, err := h.stepExecutionSvc.ListStepsByWorkflowExecutionID(ctx, service.ListStepsByWorkflowExecutionIDParams{
		WorkflowExecutionID: request.WorkflowExecutionId,
	})
	if err != nil {
		return nil, fmt.Errorf("list steps by workflow execution id: %w", err)
	}

	items := make([]gen.StepExecutionResponse, len(stepExecutions))
	for i, stepExecution := range stepExecutions {
		res, err := converter.ToStepExecutionResponse(stepExecution)
		if err != nil {
			return nil, fmt.Errorf("convert to step execution response: %w", err)
		}

		items[i] = res
	}

	return gen.StepExecutionListByWorkflowExecutionId200JSONResponse(items), nil
}
