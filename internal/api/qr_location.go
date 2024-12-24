package api

import (
	"encoding/json"
	"net/http"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/service/qr_location"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func (s *HTTPServer) GetQRLocationById(w http.ResponseWriter, r *http.Request, qrLocationId QRLocationId) {
	modelQRLocation, err := s.qrLocationSvc.GetByID(r.Context(), qrlocation.GetQrLocationByIDQuery{ID: qrLocationId})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelQRLocationToDTO(modelQRLocation))
}

func (s *HTTPServer) ListQRLocations(w http.ResponseWriter, r *http.Request, params ListQRLocationsParams) {
	pagingParams := paging.NewParams(
		params.PageSize,
		params.Page,
		paging.WithMaxPageSize(10000),
	)

	sorts, err := xsort.NewList(params.Sort)
	if err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(err, "invalid sort"))
	}

	qrLocationList, err := s.qrLocationSvc.List(r.Context(), qrlocation.ListQrLocationQuery{
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	items := make([]QRLocationResponse, 0, len(qrLocationList.Items))
	for _, r := range qrLocationList.Items {
		items = append(items, modelQRLocationToDTO(r))
	}

	s.respondJSON(w, http.StatusOK, PagingQRLocationResponse{
		Items:      items,
		TotalItems: qrLocationList.TotalItem,
	})
}

func (s *HTTPServer) CreateQRLocation(w http.ResponseWriter, r *http.Request) {
	var req CreateQRLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	modelQRLocation, err := s.qrLocationSvc.Create(r.Context(), qrlocation.CreateQrLocationCommand{
		Name:     req.Name,
		QRCode:   req.QrCode,
		Metadata: req.Metadata,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusCreated, modelQRLocationToDTO(modelQRLocation))
}

func (s *HTTPServer) UpdateQRLocationById(w http.ResponseWriter, r *http.Request, qrLocationId QRLocationId) {
	var req UpdateQRLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, r, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	modelQRLocation, err := s.qrLocationSvc.Update(r.Context(), qrlocation.UpdateQRLocationCommand{
		ID:       qrLocationId,
		Name:     req.Name,
		QRCode:   req.QrCode,
		Metadata: req.Metadata,
	})
	if err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, modelQRLocationToDTO(modelQRLocation))
}

func (s *HTTPServer) DeleteQRLocationById(w http.ResponseWriter, r *http.Request, qrLocationId QRLocationId) {
	if err := s.qrLocationSvc.Delete(r.Context(), qrlocation.DeleteQrLocationCommand{ID: qrLocationId}); err != nil {
		s.respondError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, nil)
}

func modelQRLocationToDTO(l model.QRLocation) QRLocationResponse {
	return QRLocationResponse{
		Id:        l.ID,
		Name:      l.Name,
		QrCode:    l.QRCode,
		Metadata:  l.Metadata,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}
