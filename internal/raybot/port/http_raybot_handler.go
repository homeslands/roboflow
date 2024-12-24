package port

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/internal/raybot/dto"
	"github.com/tuanvumaihuynh/roboflow/internal/raybot/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/response"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type RaybotService interface {
	GetRaybot(ctx context.Context, id uuid.UUID) (*model.Raybot, error)
	ListRaybots(ctx context.Context) ([]*model.Raybot, error)
	CreateRaybot(ctx context.Context, r *model.Raybot) (*model.Raybot, error)
	DeleteRaybot(ctx context.Context, id uuid.UUID) error

	GetRaybotStatus(ctx context.Context, id uuid.UUID) (model.RaybotStatus, error)
	UpdateRaybotStatus(ctx context.Context, id uuid.UUID, status model.RaybotStatus) error
}

type RaybotHandler struct {
	raybotService RaybotService
	logger        *zap.Logger
}

func NewRaybotHandler(raybotSvc RaybotService, logger *zap.Logger) *RaybotHandler {
	if raybotSvc == nil {
		panic("nil raybotSvc")
	}
	if logger == nil {
		panic("nil logger")
	}
	return &RaybotHandler{
		raybotService: raybotSvc,
		logger:        logger,
	}
}

func (h RaybotHandler) HandleGetRaybot(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	raybot, err := h.raybotService.GetRaybot(ctx, id)
	if err != nil {
		h.logger.Error("error getting raybot", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.MapRaybotToResponse(raybot))
}

func (h RaybotHandler) HandleListRaybots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	raybots, err := h.raybotService.ListRaybots(ctx)
	if err != nil {
		h.logger.Error("error listing raybots", zap.Error(err))
		response.Error(w, err)
		return
	}

	res := make([]*dto.RaybotResponse, 0)
	for _, raybot := range raybots {
		res = append(res, dto.MapRaybotToResponse(raybot))
	}

	response.JSON(w, http.StatusOK, res)
}

func (h RaybotHandler) HandleCreateRaybot(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRaybotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	raybot, err := model.NewRaybot(req.Name)
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}
	raybot, err = h.raybotService.CreateRaybot(ctx, raybot)
	if err != nil {
		h.logger.Error("error creating raybot", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.MapRaybotToResponse(raybot))
}

func (h RaybotHandler) HandleDeleteRaybot(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	if err := h.raybotService.DeleteRaybot(ctx, id); err != nil {
		h.logger.Error("error deleting raybot", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
