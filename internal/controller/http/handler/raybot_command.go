package handler

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type raybotCommandHandler struct {
	raybotCommandSvc service.RaybotCommandService
}

func newRaybotCommandHandler(raybotCommandSvc service.RaybotCommandService) *raybotCommandHandler {
	return &raybotCommandHandler{raybotCommandSvc: raybotCommandSvc}
}

func (h raybotCommandHandler) RaybotCommandGet(ctx context.Context, request gen.RaybotCommandGetRequestObject) (gen.RaybotCommandGetResponseObject, error) {
	m, err := h.raybotCommandSvc.GetRaybotCommand(ctx, service.GetRaybotCommandParams{
		ID: request.RaybotCommandId,
	})
	if err != nil {
		return nil, err
	}

	return gen.RaybotCommandGet200JSONResponse(converter.ToRaybotCommandResponse(m)), nil
}

func (h raybotCommandHandler) RaybotCommandList(ctx context.Context, request gen.RaybotCommandListRequestObject) (gen.RaybotCommandListResponseObject, error) {
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

	mp, err := h.raybotCommandSvc.ListRaybotCommandsByRaybotID(ctx, service.ListRaybotCommandsByRaybotIDParams{
		RaybotID:     request.RaybotId,
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		return nil, fmt.Errorf("raybot command service list raybot commands by raybot id: %w", err)
	}

	items := make([]gen.RaybotCommandResponse, len(mp.Items))
	for i, item := range mp.Items {
		items[i] = converter.ToRaybotCommandResponse(item)
	}

	return gen.RaybotCommandList200JSONResponse{
		Items:      items,
		TotalItems: mp.TotalItems,
	}, nil
}

func (h raybotCommandHandler) RaybotCommandCreate(ctx context.Context, request gen.RaybotCommandCreateRequestObject) (gen.RaybotCommandCreateResponseObject, error) {
	m, err := h.raybotCommandSvc.CreateRaybotCommand(ctx, service.CreateRaybotCommandParams{
		RaybotID: request.RaybotId,
		Type:     raybotcommand.Type(request.Body.Type),
		Inputs:   raybotcommand.NewInputs(request.Body.Inputs),
	})
	if err != nil {
		return nil, fmt.Errorf("raybot command service create raybot command: %w", err)
	}

	return gen.RaybotCommandCreate201JSONResponse(converter.ToRaybotCommandResponse(m)), nil
}
