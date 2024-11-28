package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/raybot/model"
	"github.com/tuanvumaihuynh/roboflow/internal/raybot/port"
)

var _ port.RaybotService = (*RaybotService)(nil)

type RaybotRepository interface {
	GetRaybot(ctx context.Context, id uuid.UUID) (*model.Raybot, error)
	ListRaybots(ctx context.Context) ([]*model.Raybot, error)
	CreateRaybot(ctx context.Context, raybot *model.Raybot) error
	DeleteRaybot(ctx context.Context, id uuid.UUID) error

	GetRaybotStatus(ctx context.Context, id uuid.UUID) (model.RaybotStatus, error)
	UpdateRaybotStatus(ctx context.Context, id uuid.UUID, status model.RaybotStatus) error
}

type RaybotService struct {
	raybotRepo RaybotRepository
}

func NewRaybotService(raybotRepo RaybotRepository) *RaybotService {
	if raybotRepo == nil {
		panic("nil raybotRepo")
	}
	return &RaybotService{
		raybotRepo: raybotRepo,
	}
}

func (s RaybotService) GetRaybot(ctx context.Context, id uuid.UUID) (*model.Raybot, error) {
	return s.raybotRepo.GetRaybot(ctx, id)
}

func (s RaybotService) ListRaybots(ctx context.Context) ([]*model.Raybot, error) {
	return s.raybotRepo.ListRaybots(ctx)
}

func (s RaybotService) CreateRaybot(ctx context.Context, raybot *model.Raybot) (*model.Raybot, error) {
	err := s.raybotRepo.CreateRaybot(ctx, raybot)
	if err != nil {
		return nil, err
	}
	return raybot, nil
}

func (s RaybotService) DeleteRaybot(ctx context.Context, id uuid.UUID) error {
	return s.raybotRepo.DeleteRaybot(ctx, id)
}

func (s RaybotService) UpdateRaybotStatus(ctx context.Context, id uuid.UUID, status model.RaybotStatus) error {
	return s.raybotRepo.UpdateRaybotStatus(ctx, id, status)
}

func (s RaybotService) GetRaybotStatus(ctx context.Context, id uuid.UUID) (model.RaybotStatus, error) {
	return s.raybotRepo.GetRaybotStatus(ctx, id)
}
