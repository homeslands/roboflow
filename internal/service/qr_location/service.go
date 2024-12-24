package qrlocation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx context.Context, cmd CreateQrLocationCommand) (model.QRLocation, error)
	Update(ctx context.Context, cmd UpdateQRLocationCommand) (model.QRLocation, error)
	Delete(ctx context.Context, cmd DeleteQrLocationCommand) error

	GetByID(ctx context.Context, q GetQrLocationByIDQuery) (model.QRLocation, error)
	List(ctx context.Context, q ListQrLocationQuery) (*paging.List[model.QRLocation], error)
}

type service struct {
	qrLocationRepo model.QRLocationRepository
}

func NewService(qrLocationRepo model.QRLocationRepository) *service {
	return &service{
		qrLocationRepo: qrLocationRepo,
	}
}

func (s service) Create(ctx context.Context, cmd CreateQrLocationCommand) (model.QRLocation, error) {
	if err := cmd.Validate(); err != nil {
		return model.QRLocation{}, err
	}

	now := time.Now()
	modelQrLocation := model.QRLocation{
		ID:        uuid.New(),
		Name:      cmd.Name,
		QRCode:    cmd.QRCode,
		Metadata:  cmd.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.qrLocationRepo.Create(ctx, modelQrLocation); err != nil {
		return model.QRLocation{}, err
	}

	return modelQrLocation, nil
}

func (s service) Update(ctx context.Context, cmd UpdateQRLocationCommand) (model.QRLocation, error) {
	if err := cmd.Validate(); err != nil {
		return model.QRLocation{}, err
	}

	loc, err := s.qrLocationRepo.Update(ctx, model.QRLocation{
		ID:        cmd.ID,
		Name:      cmd.Name,
		QRCode:    cmd.QRCode,
		Metadata:  cmd.Metadata,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return model.QRLocation{}, err
	}

	return loc, nil
}

func (s service) Delete(ctx context.Context, cmd DeleteQrLocationCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.qrLocationRepo.Delete(ctx, cmd.ID)
}

func (s service) GetByID(ctx context.Context, q GetQrLocationByIDQuery) (model.QRLocation, error) {
	if err := q.Validate(); err != nil {
		return model.QRLocation{}, err
	}

	return s.qrLocationRepo.Get(ctx, q.ID)
}

func (s service) List(ctx context.Context, q ListQrLocationQuery) (*paging.List[model.QRLocation], error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	return s.qrLocationRepo.List(ctx, q.PagingParams, q.Sorts)
}
