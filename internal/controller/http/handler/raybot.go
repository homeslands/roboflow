package handler

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/ptr"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type raybotHandler struct {
	raybotSvc service.RaybotService
}

func newRaybotHandler(raybotSvc service.RaybotService) *raybotHandler {
	return &raybotHandler{raybotSvc: raybotSvc}
}

func (h raybotHandler) RaybotGet(ctx context.Context, request gen.RaybotGetRequestObject) (gen.RaybotGetResponseObject, error) {
	m, err := h.raybotSvc.GetRaybot(ctx, service.GetRaybotParams{
		ID: request.RaybotId,
	})
	if err != nil {
		return nil, err
	}

	return gen.RaybotGet200JSONResponse(converter.ToRaybotResponse(m)), nil
}

func (h raybotHandler) RaybotList(ctx context.Context, request gen.RaybotListRequestObject) (gen.RaybotListResponseObject, error) {
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

	var controlMode *raybot.ControlMode
	if request.Params.ControlMode != nil {
		controlMode = ptr.New(raybot.ControlMode(*request.Params.ControlMode))
	}

	mp, err := h.raybotSvc.ListRaybots(ctx, service.ListRaybotsParams{
		PagingParams: pagingParams,
		Sorts:        sorts,
		IsOnline:     request.Params.IsOnline,
		ControlMode:  controlMode,
	})
	if err != nil {
		return nil, fmt.Errorf("raybot service list raybots: %w", err)
	}

	items := make([]gen.RaybotResponse, len(mp.Items))
	for i, item := range mp.Items {
		items[i] = converter.ToRaybotResponse(item)
	}

	return gen.RaybotList200JSONResponse{
		Items:      items,
		TotalItems: mp.TotalItems,
	}, nil
}

func (h raybotHandler) RaybotCreate(ctx context.Context, request gen.RaybotCreateRequestObject) (gen.RaybotCreateResponseObject, error) {
	m, err := h.raybotSvc.CreateRaybot(ctx, service.CreateRaybotParams{
		Name: request.Body.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("raybot service create raybot: %w", err)
	}

	return gen.RaybotCreate201JSONResponse(converter.ToRaybotResponse(m)), nil
}

func (h raybotHandler) RaybotDelete(ctx context.Context, request gen.RaybotDeleteRequestObject) (gen.RaybotDeleteResponseObject, error) {
	err := h.raybotSvc.DeleteRaybot(ctx, service.DeleteRaybotParams{
		ID: request.RaybotId,
	})
	if err != nil {
		return nil, fmt.Errorf("raybot service delete raybot: %w", err)
	}

	return gen.RaybotDelete204Response{}, nil
}
