package api

import (
	"encoding/json"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s *HTTPServer) GetRaybotById(w http.ResponseWriter, r *http.Request, raybotId RaybotId) {
	modelRaybot, err := s.raybotSvc.GetByID(r.Context(), raybot.GetRaybotByIDQuery{ID: raybotId})
	if err != nil {
		s.respondError(w, r, err)
		return
	}
	s.respondJSON(w, http.StatusOK, modelRaybotToDTO(modelRaybot))
}

func (s *HTTPServer) ListRaybots(w http.ResponseWriter, r *http.Request, params ListRaybotsParams) {
	pagingParams := paging.NewParams(
		params.PageSize,
		params.Page,
		paging.WithMaxPageSize(100),
	)

	sorts, err := xsort.NewList(params.Sort)
	if err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(err, "invalid sort"))
	}
	q := raybot.ListRaybotQuery{
		PagingParams: pagingParams,
		Sorts:        sorts,
	}
	if params.Status != nil {
		status := model.RaybotStatus(*params.Status)
		q.Status = &status
	}

	raybotList, err := s.raybotSvc.List(r.Context(), q)
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]RaybotResponse, 0, len(raybotList.Items))
	for _, r := range raybotList.Items {
		items = append(items, modelRaybotToDTO(r))
	}

	s.respondJSON(w, http.StatusOK, PagingRaybotResponse{
		Items:      items,
		TotalItems: raybotList.TotalItem,
	})
}

func (s *HTTPServer) CreateRaybot(w http.ResponseWriter, r *http.Request) {
	var req CreateRaybotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	modelRaybot, err := s.raybotSvc.Create(r.Context(), raybot.CreateRaybotCommand{
		Name: req.Name,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusCreated, modelRaybotToDTO(modelRaybot))
}

func (s *HTTPServer) DeleteRaybotById(w http.ResponseWriter, r *http.Request, raybotId RaybotId) {
	if err := s.raybotSvc.Delete(r.Context(), raybot.DeleteRaybotCommand{ID: raybotId}); err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusNoContent, nil)
}

func modelRaybotToDTO(r model.Raybot) RaybotResponse {
	return RaybotResponse{
		Id:              r.ID,
		Name:            r.Name,
		Token:           r.Token,
		Status:          string(r.Status),
		IpAddress:       r.IpAddress,
		LastConnectedAt: r.LastConnectedAt,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
