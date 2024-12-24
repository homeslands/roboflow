package raybot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx context.Context, cmd CreateRaybotCommand) (model.Raybot, error)
	Delete(ctx context.Context, cmd DeleteRaybotCommand) error
	UpdateState(ctx context.Context, cmd UpdateStateCommand) error

	GetByID(ctx context.Context, q GetRaybotByIDQuery) (model.Raybot, error)
	List(ctx context.Context, q ListRaybotQuery) (*paging.List[model.Raybot], error)
	GetState(ctx context.Context, q GetStatusQuery) (model.RaybotStatus, error)
}

type service struct {
	raybotRepo model.RaybotRepository
}

func NewService(raybotRepo model.RaybotRepository) *service {
	return &service{
		raybotRepo: raybotRepo,
	}
}

func (s service) Create(ctx context.Context, cmd CreateRaybotCommand) (model.Raybot, error) {
	if err := cmd.Validate(); err != nil {
		return model.Raybot{}, err
	}

	now := time.Now()
	modelRaybot := model.Raybot{
		ID:        uuid.New(),
		Name:      cmd.Name,
		Token:     uuid.NewString(),
		Status:    model.RaybotStatusOffline,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.raybotRepo.Create(ctx, modelRaybot); err != nil {
		return model.Raybot{}, err
	}

	return modelRaybot, nil
}

func (s service) Delete(ctx context.Context, cmd DeleteRaybotCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.raybotRepo.Delete(ctx, cmd.ID)
}

func (s service) UpdateState(ctx context.Context, cmd UpdateStateCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.raybotRepo.UpdateState(ctx, cmd.ID, cmd.State)
}

func (s service) GetByID(ctx context.Context, q GetRaybotByIDQuery) (model.Raybot, error) {
	if err := q.Validate(); err != nil {
		return model.Raybot{}, err
	}

	return s.raybotRepo.Get(ctx, q.ID)
}

func (s service) List(ctx context.Context, q ListRaybotQuery) (*paging.List[model.Raybot], error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	return s.raybotRepo.List(ctx, q.PagingParams, q.Sorts, q.Status)
}

func (s service) GetState(ctx context.Context, q GetStatusQuery) (model.RaybotStatus, error) {
	if err := q.Validate(); err != nil {
		return "", err
	}

	return s.raybotRepo.GetState(ctx, q.ID)
}
