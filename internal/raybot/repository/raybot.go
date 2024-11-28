package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/internal/raybot/model"
	"github.com/tuanvumaihuynh/roboflow/internal/raybot/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type RaybotRepository struct {
	store *db.Store
}

func NewRaybotRepository(s *db.Store) *RaybotRepository {
	return &RaybotRepository{
		store: s,
	}
}

var _ service.RaybotRepository = (*RaybotRepository)(nil)

func (r RaybotRepository) GetRaybot(ctx context.Context, id uuid.UUID) (*model.Raybot, error) {
	row, err := r.store.GetRaybot(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.ThrowNotFound(err, "raybot not found")
		}
		return nil, err
	}

	return rowToModel(row), nil
}

func (r RaybotRepository) ListRaybots(ctx context.Context) ([]*model.Raybot, error) {
	rows, err := r.store.ListRaybots(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Raybot
	for _, row := range rows {
		result = append(result, rowToModel(row))
	}

	return result, nil
}

func (r RaybotRepository) CreateRaybot(ctx context.Context, raybot *model.Raybot) error {
	params := db.CreateRaybotParams{
		ID:        raybot.ID,
		Name:      raybot.Name,
		Status:    string(raybot.Status),
		Token:     raybot.Token,
		CreatedAt: pgtype.Timestamptz{Time: raybot.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: raybot.UpdatedAt, Valid: true},
	}
	err := r.store.CreateRaybot(ctx, params)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return xerrors.ThrowAlreadyExists(err, "raybot name already exists")
	}

	return err
}

func (r RaybotRepository) UpdateRaybot(ctx context.Context, raybot *model.Raybot) error {
	params := db.UpdateRaybotParams{
		ID:        raybot.ID,
		Name:      raybot.Name,
		Token:     raybot.Token,
		Status:    string(raybot.Status),
		UpdatedAt: pgtype.Timestamptz{Time: raybot.UpdatedAt, Valid: true},
	}

	err := r.store.UpdateRaybot(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "raybot not found")
	}

	return err
}

func (r RaybotRepository) DeleteRaybot(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRaybot(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "raybot not found")
	}

	return err
}

func (r RaybotRepository) GetRaybotStatus(ctx context.Context, id uuid.UUID) (model.RaybotStatus, error) {
	row, err := r.store.GetRaybotStatus(ctx, id)
	if err != nil {
		return "", err
	}

	status := model.RaybotStatus(row)
	return status, nil
}

func (r RaybotRepository) UpdateRaybotStatus(ctx context.Context, id uuid.UUID, status model.RaybotStatus) error {
	params := db.UpdateRaybotStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: pgtype.Timestamptz{Time: time.Now()},
	}

	err := r.store.UpdateRaybotStatus(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func rowToModel(row *db.Raybot) *model.Raybot {
	return &model.Raybot{
		ID:        row.ID,
		Name:      row.Name,
		Token:     row.Token,
		Status:    model.RaybotStatus(row.Status),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
