package api

import (
	"encoding/json"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/workflow"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s *HTTPServer) GetWorkflowById(w http.ResponseWriter, r *http.Request, workflowId WorkflowId) {
	modelWorkflow, err := s.workflowSvc.GetByID(r.Context(), workflow.GetWorkflowByIDQuery{
		ID: workflowId,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelWorkflowToDTO(modelWorkflow))
}

func (s *HTTPServer) ListWorkflows(w http.ResponseWriter, r *http.Request, params ListWorkflowsParams) {
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

	workflowList, err := s.workflowSvc.List(r.Context(), workflow.ListWorkflowQuery{
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]PagingWorkflowItemResponse, len(workflowList.Items))
	for i, item := range workflowList.Items {
		items[i] = PagingWorkflowItemResponse{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	s.respondJSON(w, http.StatusOK, PagingWorkflowResponse{
		Items:      items,
		TotalItems: workflowList.TotalItem,
	})
}

func (s *HTTPServer) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	modelWorkflow, err := s.workflowSvc.Create(r.Context(), workflow.CreateWorkflowCommand{
		Name:        req.Name,
		Description: req.Description,
		Definition:  dtoWorkflowDefinitionToModel(req.Definition),
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusCreated, modelWorkflowToDTO(modelWorkflow))
}

func (s *HTTPServer) UpdateWorkflowById(w http.ResponseWriter, r *http.Request, workflowId WorkflowId) {
	var req UpdateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	modelWorkflow, err := s.workflowSvc.Update(r.Context(), workflow.UpdateWorkflowCommand{
		ID:          workflowId,
		Name:        req.Name,
		Description: req.Description,
		Definition:  dtoWorkflowDefinitionToModel(req.Definition),
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelWorkflowToDTO(modelWorkflow))
}
func (s *HTTPServer) DeleteWorkflowById(w http.ResponseWriter, r *http.Request, workflowId WorkflowId) {
	if err := s.workflowSvc.Delete(r.Context(), workflow.DeleteWorkflowCommand{ID: workflowId}); err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusNoContent, nil)
}

func (s *HTTPServer) RunWorkflowById(w http.ResponseWriter, r *http.Request, workflowId WorkflowId) {
	var req RunWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	workflowExecID, err := s.workflowSvc.Run(r.Context(), workflow.RunWorkflowCommand{
		ID:  workflowId,
		Env: req.Env,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusCreated, RunWorkflowResponse{
		WorkflowExecutionId: workflowExecID,
	})
}

func modelWorkflowToDTO(w model.Workflow) WorkflowResponse {
	r := WorkflowResponse{
		Id:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}

	if w.Definition == nil {
		return r
	}

	edges := make([]WorkflowEdge, 0, len(w.Definition.Edges))
	for _, edge := range w.Definition.Edges {
		edges = append(edges, modelWorkflowEdgeToDTO(edge))
	}
	nodes := make([]WorkflowNode, 0, len(w.Definition.Nodes))
	for _, node := range w.Definition.Nodes {
		nodes = append(nodes, modelWorkflowNodeToDTO(node))
	}

	r.Definition = &WorkflowDefinition{
		Edges: edges,
		Nodes: nodes,
		Viewport: ViewPort{
			X:    w.Definition.ViewPort.X,
			Y:    w.Definition.ViewPort.Y,
			Zoom: w.Definition.ViewPort.Zoom,
		},
		Position: w.Definition.Position,
		Zoom:     w.Definition.Zoom,
	}

	return r
}

func modelWorkflowEdgeToDTO(e model.WorkflowEdge) WorkflowEdge {
	return WorkflowEdge{
		Id:           e.ID,
		Type:         e.Type,
		Source:       e.Source,
		Target:       e.Target,
		SourceHandle: e.SourceHandle,
		TargetHandle: e.TargetHandle,
		Label:        e.Label,
		Animated:     e.Animated,
		SourceX:      e.SourceX,
		SourceY:      e.SourceY,
		TargetX:      e.TargetX,
		TargetY:      e.TargetY,
	}
}

func modelWorkflowNodeToDTO(n model.WorkflowNode) WorkflowNode {
	return WorkflowNode{
		Id:          n.ID,
		Type:        string(n.Type),
		Initialized: n.Initialized,
		Position: Position{
			n.Position.X,
			n.Position.Y,
		},
		Definition: NodeDefinition{
			Type:   TaskType(n.Definition.Type),
			Fields: modelNodeFieldsToDTO(n.Definition.Fields),
		},
	}
}

func dtoWorkflowDefinitionToModel(d WorkflowDefinition) model.WorkflowDefinition {
	edges := make([]model.WorkflowEdge, 0, len(d.Edges))
	for _, edge := range d.Edges {
		edges = append(edges, dtoWorkflowEdgeToModel(edge))
	}
	nodes := make([]model.WorkflowNode, 0, len(d.Nodes))
	for _, node := range d.Nodes {
		nodes = append(nodes, dtoWorkflowNodeToModel(node))
	}

	m := model.WorkflowDefinition{
		Edges:    edges,
		Nodes:    nodes,
		Zoom:     d.Zoom,
		Position: d.Position,
		ViewPort: model.ViewPort{
			X:    d.Viewport.X,
			Y:    d.Viewport.Y,
			Zoom: d.Viewport.Zoom,
		},
	}

	return m
}

func dtoWorkflowEdgeToModel(e WorkflowEdge) model.WorkflowEdge {
	return model.WorkflowEdge{
		ID:           e.Id,
		Type:         e.Type,
		Source:       e.Source,
		Target:       e.Target,
		SourceHandle: e.SourceHandle,
		TargetHandle: e.TargetHandle,
		Label:        e.Label,
		Animated:     e.Animated,
		SourceX:      e.SourceX,
		SourceY:      e.SourceY,
		TargetX:      e.TargetX,
		TargetY:      e.TargetY,
	}
}

func dtoWorkflowNodeToModel(n WorkflowNode) model.WorkflowNode {
	return model.WorkflowNode{
		ID:          n.Id,
		Type:        model.NodeType(n.Type),
		Initialized: n.Initialized,
		Position: struct {
			X float32 `json:"x" validate:"required"`
			Y float32 `json:"y" validate:"required"`
		}{
			X: n.Position.X,
			Y: n.Position.Y,
		},
		Definition: model.NodeDefinition{
			Type:   model.TaskType(n.Definition.Type),
			Fields: dtoNodeFieldsToModel(n.Definition.Fields),
		},
	}
}

func modelNodeFieldsToDTO(fields map[string]model.NodeField) map[string]NodeField {
	m := make(map[string]NodeField, len(fields))
	for k, f := range fields {
		m[k] = NodeField{
			UseEnv: f.UseEnv,
			Value:  f.Value,
		}
	}

	return m
}

func dtoNodeFieldsToModel(fields map[string]NodeField) map[string]model.NodeField {
	m := make(map[string]model.NodeField, len(fields))
	for k, f := range fields {
		m[k] = model.NodeField{
			UseEnv: f.UseEnv,
			Value:  f.Value,
		}
	}

	return m
}
