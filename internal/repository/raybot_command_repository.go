package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

var _ model.RaybotCommandRepository = (*raybotCommandRepository)(nil)

type raybotCommandRepository struct {
	store *db.Store
}

func NewRaybotCommandRepository(s *db.Store) *raybotCommandRepository {
	return &raybotCommandRepository{
		store: s,
	}
}

func (r raybotCommandRepository) Get(ctx context.Context, id uuid.UUID) (model.RaybotCommand, error) {
	row, err := r.store.GetRaybotCommand(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RaybotCommand{}, xerrors.ThrowNotFound(err, "command not found")
		}
		return model.RaybotCommand{}, err
	}

	return rowRaybotCommandToModel(*row)
}

func (r raybotCommandRepository) List(ctx context.Context, raybotId uuid.UUID, p paging.Params, sorts []xsort.Sort) (*paging.List[model.RaybotCommand], error) {
	params := db.ListRaybotCommandsParams{
		RaybotID: raybotId,
		Limit:    p.Limit(),
		Offset:   p.Offset(),
		Sorts:    sorts,
	}

	row, err := r.store.ListRaybotCommands(ctx, params)
	if err != nil {
		return nil, err
	}

	var items []model.RaybotCommand
	for _, item := range row.Items {
		m, err := rowRaybotCommandToModel(item)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}

	return paging.NewList(items, row.TotalItem), nil
}

func (r raybotCommandRepository) Create(ctx context.Context, cmd model.RaybotCommand) error {
	param := db.CreateRaybotCommandParams{
		ID:        cmd.ID,
		RaybotID:  cmd.RaybotID,
		Type:      string(cmd.Type),
		Status:    string(cmd.Status),
		CreatedAt: pgtype.Timestamptz{Time: cmd.CreatedAt, Valid: true},
	}
	if cmd.Input != nil {
		// Compact input to save space
		var buffer bytes.Buffer
		if err := json.Compact(&buffer, cmd.Input); err != nil {
			return fmt.Errorf("failed to compact input: %w", err)
		}
		param.Input = buffer.Bytes()
	}

	if err := r.store.CreateRaybotCommand(ctx, param); err != nil {
		return err
	}

	return nil
}

func (r raybotCommandRepository) Update(
	ctx context.Context,
	cmdID uuid.UUID,
	raybotStatus model.RaybotStatus,
	fn func(raybotCmd *model.RaybotCommand) error,
) error {
	return r.store.WithTx(ctx, func(s db.Store) error {
		row, err := r.store.GetRaybotCommandForUpdate(ctx, cmdID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return xerrors.ThrowNotFound(err, "command not found")
			}
			return err
		}

		cmd, err := rowRaybotCommandToModel(*row)
		if err != nil {
			return err
		}

		err = fn(&cmd)
		if err != nil {
			return err
		}

		// Update raybot command
		params := db.UpdateRaybotCommandParams{
			ID:        cmd.ID,
			RaybotID:  cmd.RaybotID,
			Type:      string(cmd.Type),
			Status:    string(cmd.Status),
			Input:     cmd.Input,
			Output:    cmd.Output,
			CreatedAt: pgtype.Timestamptz{Time: cmd.CreatedAt, Valid: true},
		}
		if cmd.CompletedAt != nil {
			params.CompletedAt = pgtype.Timestamptz{Time: *cmd.CompletedAt, Valid: true}
		}

		err = r.store.UpdateRaybotCommand(ctx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return xerrors.ThrowNotFound(err, "command not found")
			}
			return err
		}

		// Update raybot status
		err = r.store.UpdateRaybotStatus(ctx, db.UpdateRaybotStatusParams{
			ID:        cmd.RaybotID,
			Status:    string(raybotStatus),
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return xerrors.ThrowNotFound(err, "raybot not found")
			}
			return err
		}

		return nil
	})
}

func (r raybotCommandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRaybotCommand(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "command not found")
	}

	return err
}

func rowRaybotCommandToModel(row db.RaybotCommand) (model.RaybotCommand, error) {
	m := model.RaybotCommand{
		ID:        row.ID,
		RaybotID:  row.RaybotID,
		Status:    model.RaybotCommandStatus(row.Status),
		Type:      model.RaybotCommandType(row.Type),
		Input:     row.Input,
		Output:    row.Output,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.CompletedAt.Valid {
		m.CompletedAt = &row.CompletedAt.Time
	}

	return m, nil
}
