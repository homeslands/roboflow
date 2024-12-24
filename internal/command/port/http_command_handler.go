package port

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/internal/command/dto"
	"github.com/tuanvumaihuynh/roboflow/internal/command/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/response"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type Service interface {
	GetCommand(ctx context.Context, id uuid.UUID) (*model.Command, error)
	ListCommands(ctx context.Context, raybotId uuid.UUID) ([]*model.Command, error)
	CreateCommand(ctx context.Context, cmd *model.Command) (*model.Command, error)
	UpdateCommand(ctx context.Context, cmd *model.Command) (*model.Command, error)
	DeleteCommand(ctx context.Context, id uuid.UUID) error
}

type CommandHandler struct {
	commandService Service
	logger         *zap.Logger
}

func NewCommandHandler(commandService Service, logger *zap.Logger) *CommandHandler {
	if commandService == nil {
		panic("nil commandService")
	}
	if logger == nil {
		panic("nil logger")
	}
	return &CommandHandler{
		commandService: commandService,
		logger:         logger,
	}
}

func (h CommandHandler) HandleGetCommand(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	cmd, err := h.commandService.GetCommand(ctx, id)
	if err != nil {
		h.logger.Error("error getting command", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.MapCommandToResponse(cmd))
}

func (h CommandHandler) HandleListCommands(w http.ResponseWriter, r *http.Request) {
	raybotId, err := uuid.Parse(chi.URLParam(r, "raybotId"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	cmds, err := h.commandService.ListCommands(ctx, raybotId)
	if err != nil {
		h.logger.Error("error listing commands", zap.Error(err))
		response.Error(w, err)
		return
	}

	res := make([]dto.CommandResponse, 0)
	for _, cmd := range cmds {
		res = append(res, dto.MapCommandToResponse(cmd))
	}

	response.JSON(w, http.StatusOK, res)
}

func (h CommandHandler) HandleCreateCommand(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	cmd := model.NewCommand(req.RaybotID, req.Type)
	cmd, err := h.commandService.CreateCommand(ctx, cmd)
	if err != nil {
		h.logger.Error("error creating command", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.MapCommandToResponse(cmd))
}

func (h CommandHandler) HandleDeleteCommand(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, xerrors.ThrowInvalidArgument(nil, err.Error()))
		return
	}

	ctx := r.Context()
	if err := h.commandService.DeleteCommand(ctx, id); err != nil {
		h.logger.Error("error deleting command", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
