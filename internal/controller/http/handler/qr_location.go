package handler

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/converter"
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
)

type qrLocationHandler struct {
	qrLocationSvc service.QRLocationService
}

func newQRLocationHandler(qrLocationSvc service.QRLocationService) *qrLocationHandler {
	return &qrLocationHandler{qrLocationSvc: qrLocationSvc}
}

func (h qrLocationHandler) QrLocationGet(ctx context.Context, request gen.QrLocationGetRequestObject) (gen.QrLocationGetResponseObject, error) {
	m, err := h.qrLocationSvc.GetQRLocation(ctx, service.GetQRLocationParams{
		ID: request.QrLocationId,
	})
	if err != nil {
		return nil, fmt.Errorf("qr location service get qr location: %w", err)
	}

	return gen.QrLocationGet200JSONResponse(converter.ToQRLocationResponse(m)), nil

}

func (h qrLocationHandler) QrLocationList(ctx context.Context, request gen.QrLocationListRequestObject) (gen.QrLocationListResponseObject, error) {
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

	mp, err := h.qrLocationSvc.ListQRLocations(ctx, service.ListQRLocationsParams{
		PagingParams: pagingParams,
		Sorts:        sorts,
	})
	if err != nil {
		return nil, fmt.Errorf("qr location service list qr locations: %w", err)
	}

	items := make([]gen.QRLocationResponse, len(mp.Items))
	for i, qrLocation := range mp.Items {
		items[i] = converter.ToQRLocationResponse(qrLocation)
	}

	return gen.QrLocationList200JSONResponse(
		gen.QRLocationsListResponse{
			Items:      items,
			TotalItems: mp.TotalItems,
		},
	), nil
}

func (h qrLocationHandler) QrLocationCreate(ctx context.Context, request gen.QrLocationCreateRequestObject) (gen.QrLocationCreateResponseObject, error) {
	m, err := h.qrLocationSvc.CreateQRLocation(ctx, service.CreateQRLocationParams{
		Name:     request.Body.Name,
		QRCode:   request.Body.QrCode,
		Metadata: request.Body.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("qr location service create qr location: %w", err)
	}

	return gen.QrLocationCreate201JSONResponse(converter.ToQRLocationResponse(m)), nil
}

func (h qrLocationHandler) QrLocationUpdate(ctx context.Context, request gen.QrLocationUpdateRequestObject) (gen.QrLocationUpdateResponseObject, error) {
	m, err := h.qrLocationSvc.UpdateQRLocation(ctx, service.UpdateQRLocationParams{
		ID:          request.QrLocationId,
		Name:        request.Body.Name,
		SetName:     true,
		QRCode:      request.Body.QrCode,
		SetQRCode:   true,
		Metadata:    request.Body.Metadata,
		SetMetadata: true,
	})
	if err != nil {
		return nil, fmt.Errorf("qr location service update qr location: %w", err)
	}

	return gen.QrLocationUpdate200JSONResponse(converter.ToQRLocationResponse(m)), nil
}

func (h qrLocationHandler) QrLocationDelete(ctx context.Context, request gen.QrLocationDeleteRequestObject) (gen.QrLocationDeleteResponseObject, error) {
	err := h.qrLocationSvc.DeleteQRLocation(ctx, service.DeleteQRLocationParams{
		ID: request.QrLocationId,
	})
	if err != nil {
		return nil, fmt.Errorf("qr location service delete qr location: %w", err)
	}

	return gen.QrLocationDelete204Response{}, nil
}
