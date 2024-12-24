package api

import (
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/step"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s *HTTPServer) GetStepById(w http.ResponseWriter, r *http.Request, stepId StepId) {
	step, err := s.stepSvc.GetByID(r.Context(), step.GetStepByIDQuery{ID: stepId})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelStepToDTO(step))
}

func (s *HTTPServer) ListStepsByWorkflowExecutionId(w http.ResponseWriter, r *http.Request,
	workflowExecutionId WorkflowExecutionId, params ListStepsByWorkflowExecutionIdParams) {
	sorts, err := xsort.NewList(params.Sort)
	if err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(err, "invalid sort"))
	}

	steps, err := s.stepSvc.List(r.Context(), step.ListStepQuery{
		WorkflowExecutionID: workflowExecutionId,
		Sorts:               sorts,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]StepResponse, len(steps))
	for i, item := range steps {
		items[i] = modelStepToDTO(item)
	}

	s.respondJSON(w, http.StatusOK, items)
}

func modelStepToDTO(s model.Step) StepResponse {
	return StepResponse{
		Id:                  s.ID,
		WorkflowExecutionId: s.WorkflowExecutionID,
		Env:                 s.Env,
		Status:              string(s.Status),
		Node:                modelWorkflowNodeToDTO(s.Node),
		StartedAt:           s.StartedAt,
		CompletedAt:         s.CompletedAt,
	}
}
