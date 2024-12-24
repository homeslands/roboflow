package repository

import (
	"context"
	"errors"
	"strings"
	"time"

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

var _ model.RaybotRepository = (*raybotRepository)(nil)

type raybotRepository struct {
	store *db.Store
}

func NewRaybotRepository(s *db.Store) *raybotRepository {
	return &raybotRepository{
		store: s,
	}
}

func (r raybotRepository) Get(ctx context.Context, id uuid.UUID) (model.Raybot, error) {
	row, err := r.store.GetRaybot(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Raybot{}, xerrors.ThrowNotFound(err, "raybot not found")
		}
		return model.Raybot{}, err
	}

	return rowRaybotToModel(*row), nil
}

func (r raybotRepository) List(ctx context.Context, p paging.Params, sorts []xsort.Sort, status *model.RaybotStatus) (*paging.List[model.Raybot], error) {
	params := db.ListRaybotsParams{
		Limit:  p.Limit(),
		Offset: p.Offset(),
		Sorts:  sorts,
	}
	if status != nil {
		params.Status = (*string)(status)
	}

	row, err := r.store.ListRaybots(ctx, params)
	if err != nil {
		return nil, err
	}

	var items []model.Raybot
	for _, item := range row.Items {
		items = append(items, rowRaybotToModel(item))
	}

	return paging.NewList(items, row.TotalItem), nil
}

func (r raybotRepository) Create(ctx context.Context, raybot model.Raybot) error {
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
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "name") {
		return xerrors.ThrowAlreadyExists(err, "raybot name already exists")
	}

	return err
}

func (r raybotRepository) Update(ctx context.Context, raybot model.Raybot) (model.Raybot, error) {
	row, err := r.store.UpdateRaybot(ctx, db.UpdateRaybotParams{
		ID:        raybot.ID,
		Name:      raybot.Name,
		Token:     raybot.Token,
		Status:    string(raybot.Status),
		UpdatedAt: pgtype.Timestamptz{Time: raybot.UpdatedAt, Valid: true},
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "name") {
		return model.Raybot{}, xerrors.ThrowAlreadyExists(err, "raybot name already exists")
	} else if errors.Is(err, pgx.ErrNoRows) {
		return model.Raybot{}, xerrors.ThrowNotFound(err, "raybot not found")
	}

	return rowRaybotToModel(*row), err
}

func (r raybotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRaybot(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "raybot not found")
	}

	return err
}

func (r raybotRepository) GetState(ctx context.Context, id uuid.UUID) (model.RaybotStatus, error) {
	row, err := r.store.GetRaybotStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", xerrors.ThrowNotFound(err, "raybot not found")
		}

		return "", err
	}

	status := model.RaybotStatus(row)
	return status, nil
}

func (r raybotRepository) UpdateState(ctx context.Context, id uuid.UUID, status model.RaybotStatus) error {
	params := db.UpdateRaybotStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: pgtype.Timestamptz{Time: time.Now()},
	}

	err := r.store.UpdateRaybotStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return xerrors.ThrowNotFound(err, "raybot not found")
		}

		return err
	}

	return nil
}

func rowRaybotToModel(row db.Raybot) model.Raybot {
	m := model.Raybot{
		ID:        row.ID,
		Name:      row.Name,
		Token:     row.Token,
		Status:    model.RaybotStatus(row.Status),
		IpAddress: row.IpAddress,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.LastConnectedAt.Valid {
		m.LastConnectedAt = &row.LastConnectedAt.Time
	}

	return m
}
