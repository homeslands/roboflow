package api

import (
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/service/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s *HTTPServer) GetWorkflowExecutionById(w http.ResponseWriter, r *http.Request, workflowExecutionId WorkflowExecutionId) {
	modelWorkflowExecution, err := s.workflowExecutionSvc.GetByID(r.Context(), workflowexecution.GetWorkflowExecutionByIDQuery{
		ID: workflowExecutionId,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelWorkflowExecutionToDTO(modelWorkflowExecution))
}

func (s *HTTPServer) GetWorkflowExecutionStatusById(w http.ResponseWriter, r *http.Request, workflowExecutionId WorkflowExecutionId) {
	modelWorkflowExecution, err := s.workflowExecutionSvc.GetByID(r.Context(), workflowexecution.GetWorkflowExecutionByIDQuery{
		ID: workflowExecutionId,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: string(modelWorkflowExecution.Status),
	})
}

func (s *HTTPServer) ListWorkflowExecutionsByWorkflowId(w http.ResponseWriter, r *http.Request, workflowId WorkflowId, params ListWorkflowExecutionsByWorkflowIdParams) {
	pagingParams := paging.NewParams(
		params.PageSize,
		params.Page,
		paging.WithMaxPageSize(100),
	)

	sorts, err := xsort.NewList(params.Sort)
	if err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(err, "invalid sort"))
		return
	}

	workflowExecutionList, err := s.workflowExecutionSvc.List(r.Context(), workflowexecution.ListWorkflowExecutionQuery{
		WorkflowID:   workflowId,
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]WorkflowExecutionResponse, len(workflowExecutionList.Items))
	for i, item := range workflowExecutionList.Items {
		items[i] = modelWorkflowExecutionToDTO(item)
	}

	s.respondJSON(w, http.StatusOK, PagingWorkflowExecutionResponse{
		Items:      items,
		TotalItems: workflowExecutionList.TotalItem,
	})
}

func modelWorkflowExecutionToDTO(m model.WorkflowExecution) WorkflowExecutionResponse {
	res := WorkflowExecutionResponse{
		Id:          m.ID,
		WorkflowId:  m.WorkflowID,
		Env:         m.Env,
		Status:      string(m.Status),
		CreatedAt:   m.CreatedAt,
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
	}

	edges := make([]WorkflowEdge, 0, len(m.Definition.Edges))
	for _, edge := range m.Definition.Edges {
		edges = append(edges, modelWorkflowEdgeToDTO(edge))
	}
	nodes := make([]WorkflowNode, 0, len(m.Definition.Nodes))
	for _, node := range m.Definition.Nodes {
		nodes = append(nodes, modelWorkflowNodeToDTO(node))
	}

	res.Definition = WorkflowDefinition{
		Edges: edges,
		Nodes: nodes,
		Viewport: ViewPort{
			X:    m.Definition.ViewPort.X,
			Y:    m.Definition.ViewPort.Y,
			Zoom: m.Definition.ViewPort.Zoom,
		},
		Position: m.Definition.Position,
		Zoom:     m.Definition.Zoom,
	}

	return res
}
