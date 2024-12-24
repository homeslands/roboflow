package api

import (
	"encoding/json"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s HTTPServer) GetRaybotCommandById(w http.ResponseWriter, r *http.Request, raybotCommandId RaybotCommandId) {
	modelRaybotCommand, err := s.raybotCommandSvc.GetByID(r.Context(), raybotcommand.GetRaybotCommandByIDQuery{ID: raybotCommandId})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelRaybotCommandToDTO(modelRaybotCommand))
}
func (s HTTPServer) ListRaybotCommands(w http.ResponseWriter, r *http.Request, raybotId RaybotId, params ListRaybotCommandsParams) {
	pagingParams := paging.NewParams(
		params.PageSize,
		params.Page,
		paging.WithMaxPageSize(100),
	)

	sorts, err := xsort.NewList(params.Sort)
	if err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(err, "invalid sort"))
	}
	raybotCommandList, err := s.raybotCommandSvc.List(r.Context(), raybotcommand.ListRaybotCommandQuery{
		RaybotID:     raybotId,
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]RaybotCommandResponse, 0, len(raybotCommandList.Items))
	for _, r := range raybotCommandList.Items {
		items = append(items, modelRaybotCommandToDTO(r))
	}

	s.respondJSON(w, http.StatusOK, PagingRaybotCommandResponse{
		Items:      items,
		TotalItems: raybotCommandList.TotalItem,
	})
}

func (s HTTPServer) CreateRaybotCommand(w http.ResponseWriter, r *http.Request, raybotId RaybotId) {
	var req CreateRaybotCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, err)
		return
	}

	modelRaybotCommand, err := s.raybotCommandSvc.Create(r.Context(), raybotcommand.CreateRaybotCommandCommand{
		RaybotID: raybotId,
		Type:     model.RaybotCommandType(req.Type),
		Input:    req.Input,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusCreated, modelRaybotCommandToDTO(modelRaybotCommand))
}

func modelRaybotCommandToDTO(r model.RaybotCommand) RaybotCommandResponse {
	return RaybotCommandResponse{
		Id:          r.ID,
		RaybotId:    r.RaybotID,
		Type:        string(r.Type),
		Status:      string(r.Status),
		Input:       r.Input,
		Output:      r.Output,
		CreatedAt:   r.CreatedAt,
		CompletedAt: r.CompletedAt,
	}
}
