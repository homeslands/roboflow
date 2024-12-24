package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

var _ model.QRLocationRepository = (*qrLocationRepository)(nil)

type qrLocationRepository struct {
	store *db.Store
}

func NewQRLocationRepository(s *db.Store) *qrLocationRepository {
	return &qrLocationRepository{
		store: s,
	}
}

func (r qrLocationRepository) Get(ctx context.Context, id uuid.UUID) (model.QRLocation, error) {
	row, err := r.store.GetQRLocation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.QRLocation{}, xerrors.ThrowNotFound(err, "qr location not found")
		}
		return model.QRLocation{}, err
	}

	return rowQrLocationToModel(*row), nil
}

func (r qrLocationRepository) List(ctx context.Context, p paging.Params, sorts []xsort.Sort) (*paging.List[model.QRLocation], error) {
	params := db.ListQrLocationsParams{
		Limit:  p.Limit(),
		Offset: p.Offset(),
		Sorts:  sorts,
	}
	row, err := r.store.ListQrLocations(ctx, params)
	if err != nil {
		return nil, err
	}

	var items []model.QRLocation
	for _, item := range row.Items {
		items = append(items, rowQrLocationToModel(item))
	}

	return paging.NewList(items, row.TotalItem), nil
}

func (r qrLocationRepository) ExistByQRCode(ctx context.Context, qrCode string) (bool, error) {
	exists, err := r.store.ExistsQRLocationByQRCode(ctx, qrCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r qrLocationRepository) Create(ctx context.Context, loc model.QRLocation) error {
	err := r.store.CreateQRLocation(ctx,
		db.CreateQRLocationParams{
			ID:        loc.ID,
			Name:      loc.Name,
			QrCode:    loc.QRCode,
			Metadata:  loc.Metadata,
			CreatedAt: pgtype.Timestamptz{Time: loc.CreatedAt, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: loc.UpdatedAt, Valid: true},
		},
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "qr_code") {
		return xerrors.ThrowAlreadyExists(err, "qr code already exists")
	}

	return err
}

func (r qrLocationRepository) Update(ctx context.Context, loc model.QRLocation) (model.QRLocation, error) {
	row, err := r.store.UpdateQRLocation(ctx, db.UpdateQRLocationParams{
		ID:        loc.ID,
		Name:      loc.Name,
		QrCode:    loc.QRCode,
		Metadata:  loc.Metadata,
		UpdatedAt: pgtype.Timestamptz{Time: loc.UpdatedAt, Valid: true},
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "qr_code") {
		return model.QRLocation{}, xerrors.ThrowAlreadyExists(err, "qr code already exists")
	} else if errors.Is(err, pgx.ErrNoRows) {
		return model.QRLocation{}, xerrors.ThrowNotFound(err, "qr location not found")
	}

	return rowQrLocationToModel(*row), nil
}

func (r qrLocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteQRLocation(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "qr location not found")
	}

	return err
}

func rowQrLocationToModel(row db.QrLocation) model.QRLocation {
	return model.QRLocation{
		ID:        row.ID,
		Name:      row.Name,
		QRCode:    row.QrCode,
		Metadata:  row.Metadata,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
