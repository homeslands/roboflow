package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/location/model"
	"github.com/tuanvumaihuynh/roboflow/internal/location/port"
)

var _ port.QrLocationService = (*QrLocationService)(nil)

type QrLocationRepository interface {
	GetQRLocation(ctx context.Context, id uuid.UUID) (*model.QrLocation, error)
	ListQRLocations(ctx context.Context) ([]*model.QrLocation, error)
	CreateQRLocation(ctx context.Context, qrLocation *model.QrLocation) error
	UpdateQRLocation(ctx context.Context, qrLocation *model.QrLocation) error
	DeleteQRLocation(ctx context.Context, id uuid.UUID) error
}

type QrLocationService struct {
	qrLocationRepo QrLocationRepository
}

func NewQrLocationService(qrLocationRepo QrLocationRepository) *QrLocationService {
	if qrLocationRepo == nil {
		panic("nil qrLocationRepo")
	}
	return &QrLocationService{
		qrLocationRepo: qrLocationRepo,
	}
}

func (s QrLocationService) GetQRLocation(ctx context.Context, id uuid.UUID) (*model.QrLocation, error) {
	return s.qrLocationRepo.GetQRLocation(ctx, id)
}

func (s QrLocationService) ListQRLocations(ctx context.Context) ([]*model.QrLocation, error) {
	return s.qrLocationRepo.ListQRLocations(ctx)
}

func (s QrLocationService) CreateQRLocation(ctx context.Context, qrLocation *model.QrLocation) (*model.QrLocation, error) {
	err := s.qrLocationRepo.CreateQRLocation(ctx, qrLocation)
	if err != nil {
		return nil, err
	}
	return qrLocation, nil
}

func (s QrLocationService) UpdateQRLocation(ctx context.Context, qrLocation *model.QrLocation) (*model.QrLocation, error) {
	err := s.qrLocationRepo.UpdateQRLocation(ctx, qrLocation)
	if err != nil {
		return nil, err
	}
	return qrLocation, nil
}

func (s QrLocationService) DeleteQRLocation(ctx context.Context, id uuid.UUID) error {
	return s.qrLocationRepo.DeleteQRLocation(ctx, id)
}
