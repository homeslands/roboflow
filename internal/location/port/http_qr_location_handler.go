package port

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/internal/location/dto"
	"github.com/tuanvumaihuynh/roboflow/internal/location/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/response"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type QrLocationService interface {
	GetQRLocation(ctx context.Context, id uuid.UUID) (*model.QrLocation, error)
	ListQRLocations(ctx context.Context) ([]*model.QrLocation, error)
	CreateQRLocation(ctx context.Context, qrLocation *model.QrLocation) (*model.QrLocation, error)
	UpdateQRLocation(ctx context.Context, qrLocation *model.QrLocation) (*model.QrLocation, error)
	DeleteQRLocation(ctx context.Context, id uuid.UUID) error
}

type QrLocationHandler struct {
	qrLocationService QrLocationService
	logger            *zap.Logger
}

func NewQrLocationHandler(qrLocationService QrLocationService, logger *zap.Logger) *QrLocationHandler {
	if qrLocationService == nil {
		panic("nil qrLocationService")
	}
	if logger == nil {
		panic("nil logger")
	}
	return &QrLocationHandler{
		qrLocationService: qrLocationService,
		logger:            logger,
	}
}

func (h QrLocationHandler) HandleGetQRLocation(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	qrLocation, err := h.qrLocationService.GetQRLocation(ctx, id)
	if err != nil {
		h.logger.Error("error getting qr location", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.MapQRLocationToResponse(qrLocation))
}

func (h QrLocationHandler) HandleListQRLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qrLocations, err := h.qrLocationService.ListQRLocations(ctx)
	if err != nil {
		h.logger.Error("error listing qr locations", zap.Error(err))
		response.Error(w, err)
		return
	}

	res := make([]dto.QRLocationResponse, 0)
	for _, qrLocation := range qrLocations {
		res = append(res, dto.MapQRLocationToResponse(qrLocation))
	}

	response.JSON(w, http.StatusOK, res)
}

func (h QrLocationHandler) HandleCreateQRLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQRLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	qrLocation, err := model.NewQrLocation(req.Name, req.QrCode)
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}
	qrLocation, err = h.qrLocationService.CreateQRLocation(ctx, qrLocation)
	if err != nil {
		h.logger.Error("error creating qr location", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.MapQRLocationToResponse(qrLocation))
}

func (h QrLocationHandler) HandleUpdateQRLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateQRLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	qrLocation, err := model.NewQrLocation(req.Name, req.QrCode)
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}
	qrLocation.ID = id
	qrLocation, err = h.qrLocationService.UpdateQRLocation(ctx, qrLocation)
	if err != nil {
		h.logger.Error("error updating qr location", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.MapQRLocationToResponse(qrLocation))
}

func (h QrLocationHandler) HandleDeleteQRLocation(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	if err := h.qrLocationService.DeleteQRLocation(ctx, id); err != nil {
		h.logger.Error("error deleting qr location", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
