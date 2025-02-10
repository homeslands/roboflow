package serviceimpl

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.RaybotService = (*raybotService)(nil)

type raybotService struct {
	raybotRepo    repository.RaybotRepository
	sqlDBProvider sqldb.Provider
	validator     validator.Validator
}

func newRaybotService(
	raybotRepo repository.RaybotRepository,
	sqlDBProvider sqldb.Provider,
	validator validator.Validator,
) *raybotService {
	return &raybotService{
		raybotRepo:    raybotRepo,
		sqlDBProvider: sqlDBProvider,
		validator:     validator,
	}
}

func (s raybotService) GetRaybot(ctx context.Context, params service.GetRaybotParams) (raybot.Raybot, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybot.Raybot{}, fmt.Errorf("validate params: %w", err)
	}

	rb, err := s.raybotRepo.GetRaybot(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return raybot.Raybot{}, fmt.Errorf("repo get raybot: %w", err)
	}

	return rb, nil
}

func (s raybotService) ListRaybots(ctx context.Context, params service.ListRaybotsParams) (paging.List[raybot.Raybot], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("validate params: %w", err)
	}

	rb, err := s.raybotRepo.ListRaybots(
		ctx,
		s.sqlDBProvider.DB(),
		params.PagingParams,
		params.Sorts,
		params.IsOnline,
		params.ControlMode,
	)
	if err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("repo list raybots: %w", err)
	}

	return rb, nil
}

func (s raybotService) CreateRaybot(ctx context.Context, params service.CreateRaybotParams) (raybot.Raybot, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybot.Raybot{}, fmt.Errorf("validate params: %w", err)
	}

	rb := raybot.NewRaybot(params.Name)

	err := s.raybotRepo.CreateRaybot(ctx, s.sqlDBProvider.DB(), rb)
	if err != nil {
		return raybot.Raybot{}, fmt.Errorf("repo create raybot: %w", err)
	}

	return rb, nil
}

func (s raybotService) UpdateRaybot(ctx context.Context, params service.UpdateRaybotParams) (raybot.Raybot, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybot.Raybot{}, fmt.Errorf("validate params: %w", err)
	}

	rb, err := s.raybotRepo.UpdateRaybot(ctx, s.sqlDBProvider.DB(), repository.UpdateRaybotParams{
		ID:                 params.ID,
		Name:               params.Name,
		SetName:            params.SetName,
		ControlMode:        params.ControlMode,
		SetControlMode:     params.SetControlMode,
		IsOnline:           params.IsOnline,
		SetIsOnline:        params.SetIsOnline,
		IPAddress:          params.IPAddress,
		SetIPAddress:       params.SetIPAddress,
		LastConnectedAt:    params.LastConnectedAt,
		SetLastConnectedAt: params.SetLastConnectedAt,
	})
	if err != nil {
		return raybot.Raybot{}, fmt.Errorf("repo update raybot: %w", err)
	}

	return rb, nil
}

func (s raybotService) DeleteRaybot(ctx context.Context, params service.DeleteRaybotParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.raybotRepo.DeleteRaybot(ctx, s.sqlDBProvider.DB(), params.ID); err != nil {
		return fmt.Errorf("delete raybot: %w", err)
	}

	return nil
}
