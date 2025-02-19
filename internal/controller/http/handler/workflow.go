package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type workflowHandler struct {
	workflowSvc service.WorkflowService
}

func newWorkflowHandler(workflowSvc service.WorkflowService) *workflowHandler {
	return &workflowHandler{workflowSvc: workflowSvc}
}

func (h workflowHandler) WorkflowGet(ctx context.Context, request gen.WorkflowGetRequestObject) (gen.WorkflowGetResponseObject, error) {
	m, err := h.workflowSvc.GetWorkflow(ctx, service.GetWorkflowParams{
		ID: request.WorkflowId,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service get workflow: %w", err)
	}

	res, err := converter.ToWorkflowResponse(m)
	if err != nil {
		return nil, fmt.Errorf("converter to workflow response: %w", err)
	}

	return gen.WorkflowGet200JSONResponse(res), nil
}

func (h workflowHandler) WorkflowList(ctx context.Context, request gen.WorkflowListRequestObject) (gen.WorkflowListResponseObject, error) {
	pagingParams := paging.NewParams(
		request.Params.PageSize,
		request.Params.Page,
		paging.WithMaxPageSize(1000),
	)

	var sorts []sort.Sort
	var err error
	if request.Params.Sort != nil {
		sorts, err = sort.NewListFromString(*request.Params.Sort)
		if err != nil {
			return nil, fmt.Errorf("sort new list from string: %w", err)
		}
	}

	mp, err := h.workflowSvc.ListWorkflows(ctx, service.ListWorkflowsParams{
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service list workflows: %w", err)
	}

	items := make([]gen.WorkflowItemListResponse, len(mp.Items))
	for i, item := range mp.Items {
		items[i] = converter.ToWorkflowItemListResponse(item)
	}

	return gen.WorkflowList200JSONResponse{
		Items:      items,
		TotalItems: mp.TotalItems,
	}, nil
}

func (h workflowHandler) WorkflowCreate(ctx context.Context, request gen.WorkflowCreateRequestObject) (gen.WorkflowCreateResponseObject, error) {
	data := workflow.Data{}
	err := json.Unmarshal(request.Body.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal workflow data: %w", err)
	}

	m, err := h.workflowSvc.CreateWorkflow(ctx, service.CreateWorkflowParams{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Data:        data,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service create workflow: %w", err)
	}

	res, err := converter.ToWorkflowResponse(m)
	if err != nil {
		return nil, fmt.Errorf("converter to workflow response: %w", err)
	}

	return gen.WorkflowCreate201JSONResponse(res), nil
}

func (h workflowHandler) WorkflowUpdate(ctx context.Context, request gen.WorkflowUpdateRequestObject) (gen.WorkflowUpdateResponseObject, error) {
	data := workflow.Data{}
	err := json.Unmarshal(request.Body.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal workflow data: %w", err)
	}

	m, err := h.workflowSvc.UpdateWorkflow(ctx, service.UpdateWorkflowParams{
		ID:             request.WorkflowId,
		Name:           request.Body.Name,
		SetName:        true,
		Description:    request.Body.Description,
		SetDescription: true,
		Data:           data,
		SetData:        true,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service update workflow: %w", err)
	}

	res, err := converter.ToWorkflowResponse(m)
	if err != nil {
		return nil, fmt.Errorf("converter to workflow response: %w", err)
	}

	return gen.WorkflowUpdate200JSONResponse(res), nil
}

func (h workflowHandler) WorkflowDelete(ctx context.Context, request gen.WorkflowDeleteRequestObject) (gen.WorkflowDeleteResponseObject, error) {
	err := h.workflowSvc.DeleteWorkflow(ctx, service.DeleteWorkflowParams{
		ID: request.WorkflowId,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service delete workflow: %w", err)
	}

	return gen.WorkflowDelete204Response{}, nil
}

func (h workflowHandler) WorkflowRun(ctx context.Context, request gen.WorkflowRunRequestObject) (gen.WorkflowRunResponseObject, error) {
	workflowExecutionID, err := h.workflowSvc.RunWorkflow(ctx, service.RunWorkflowParams{
		ID: request.WorkflowId,
	})
	if err != nil {
		return nil, fmt.Errorf("workflow service run workflow: %w", err)
	}

	return gen.WorkflowRun201JSONResponse(converter.ToWorkflowRunResponse(workflowExecutionID)), nil
}
