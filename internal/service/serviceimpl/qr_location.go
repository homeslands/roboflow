package serviceimpl

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/model/qr_location"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.QRLocationService = (*qrLocationService)(nil)

type qrLocationService struct {
	qrLocationRepo repository.QRLocationRepository
	sqlDBProvider  sqldb.Provider
	validator      validator.Validator
}

func newQRLocationService(
	qrLocationRepo repository.QRLocationRepository,
	sqlDBProvider sqldb.Provider,
	validator validator.Validator,
) *qrLocationService {
	return &qrLocationService{
		qrLocationRepo: qrLocationRepo,
		sqlDBProvider:  sqlDBProvider,
		validator:      validator,
	}
}

func (s qrLocationService) GetQRLocation(ctx context.Context, params service.GetQRLocationParams) (qrlocation.QRLocation, error) {
	if err := s.validator.Validate(params); err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("validate params: %w", err)
	}

	qrLocation, err := s.qrLocationRepo.GetQRLocation(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("repo get qr location by id: %w", err)
	}

	return qrLocation, nil
}

func (s qrLocationService) ListQRLocations(ctx context.Context, params service.ListQRLocationsParams) (paging.List[qrlocation.QRLocation], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("validate params: %w", err)
	}

	qrLocations, err := s.qrLocationRepo.ListQRLocations(ctx, s.sqlDBProvider.DB(), params.PagingParams, params.Sorts)
	if err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("repo list qr locations: %w", err)
	}

	return qrLocations, nil
}

func (s qrLocationService) CreateQRLocation(ctx context.Context, params service.CreateQRLocationParams) (qrlocation.QRLocation, error) {
	if err := s.validator.Validate(params); err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("validate params: %w", err)
	}

	qrLocation := qrlocation.NewQRLocation(params.Name, params.QRCode, params.Metadata)

	err := s.qrLocationRepo.CreateQRLocation(ctx, s.sqlDBProvider.DB(), qrLocation)
	if err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("repo create qr location: %w", err)
	}

	return qrLocation, nil
}

func (s qrLocationService) UpdateQRLocation(ctx context.Context, params service.UpdateQRLocationParams) (qrlocation.QRLocation, error) {
	if err := s.validator.Validate(params); err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("validate params: %w", err)
	}

	qrLocation, err := s.qrLocationRepo.UpdateQRLocation(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateQRLocationParams{
			ID:          params.ID,
			Name:        params.Name,
			SetName:     params.SetName,
			QRCode:      params.QRCode,
			SetQRCode:   params.SetQRCode,
			Metadata:    params.Metadata,
			SetMetadata: params.SetMetadata,
		},
	)

	if err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("repo update qr location: %w", err)
	}

	return qrLocation, nil
}

func (s qrLocationService) DeleteQRLocation(ctx context.Context, params service.DeleteQRLocationParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	err := s.qrLocationRepo.DeleteQRLocation(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return fmt.Errorf("repo delete qr location: %w", err)
	}

	return nil
}
