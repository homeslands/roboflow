package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/internal/location/model"
	"github.com/tuanvumaihuynh/roboflow/internal/location/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type QrLocationRepository struct {
	store *db.Store
}

var _ service.QrLocationRepository = (*QrLocationRepository)(nil)

func NewQrRepository(s *db.Store) *QrLocationRepository {
	return &QrLocationRepository{
		store: s,
	}
}

func (r QrLocationRepository) GetQRLocation(ctx context.Context, id uuid.UUID) (*model.QrLocation, error) {
	row, err := r.store.GetQRLocation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.ThrowNotFound(err, "qr location not found")
		}
		return nil, err
	}

	return rowToQrLocation(row), nil
}

func (r QrLocationRepository) ListQRLocations(ctx context.Context) ([]*model.QrLocation, error) {
	rows, err := r.store.ListQRLocations(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.QrLocation
	for _, row := range rows {
		result = append(result, rowToQrLocation(row))
	}

	return result, nil
}

func (r QrLocationRepository) CreateQRLocation(ctx context.Context, qrLocation *model.QrLocation) error {
	err := r.store.CreateQRLocation(ctx,
		db.CreateQRLocationParams{
			ID:        qrLocation.ID,
			Name:      qrLocation.Name,
			QrCode:    qrLocation.QRCode,
			CreatedAt: pgtype.Timestamptz{Time: qrLocation.CreatedAt, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: qrLocation.UpdatedAt, Valid: true},
		},
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return xerrors.ThrowAlreadyExists(err, "qr code already exists")
	}

	return err
}

func (r QrLocationRepository) UpdateQRLocation(ctx context.Context, qrLocation *model.QrLocation) error {
	err := r.store.UpdateQRLocation(ctx,
		db.UpdateQRLocationParams{
			ID:        qrLocation.ID,
			Name:      qrLocation.Name,
			QrCode:    qrLocation.QRCode,
			UpdatedAt: pgtype.Timestamptz{Time: qrLocation.UpdatedAt, Valid: true},
		},
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return xerrors.ThrowAlreadyExists(err, "qr code already exists")
	}

	return err
}

func (r QrLocationRepository) DeleteQRLocation(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteQRLocation(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "qr location not found")
	}

	return err
}

func rowToQrLocation(row *db.QrLocation) *model.QrLocation {
	return &model.QrLocation{
		ID:        row.ID,
		Name:      row.Name,
		QRCode:    row.QrCode,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
