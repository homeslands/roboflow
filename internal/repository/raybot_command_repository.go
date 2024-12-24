package repository

import (
	"context"
	"encoding/json"
	"errors"

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
	inputBytes, err := json.Marshal(cmd.Input)
	if err != nil {
		return xerrors.ThrowInvalidArgument(err, "failed to marshal input")
	}
	param := db.CreateRaybotCommandParams{
		ID:        cmd.ID,
		RaybotID:  cmd.RaybotID,
		Type:      string(cmd.Type),
		Status:    string(cmd.Status),
		Input:     inputBytes,
		CreatedAt: pgtype.Timestamptz{Time: cmd.CreatedAt, Valid: true},
	}
	if err := r.store.CreateRaybotCommand(ctx, param); err != nil {
		return err
	}

	return nil
}

func (r raybotCommandRepository) Update(ctx context.Context, cmd model.RaybotCommand) error {
	param := db.UpdateRaybotCommandParams{
		ID:     cmd.ID,
		Status: string(cmd.Status),
		Output: cmd.Output.(map[string]interface{}),
	}
	if cmd.CompletedAt != nil {
		param.CompletedAt = pgtype.Timestamptz{Time: *cmd.CompletedAt, Valid: true}
	}
	if err := r.store.UpdateRaybotCommand(ctx, param); err != nil {
		return err
	}

	return nil
}

func (r raybotCommandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRaybotCommand(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "command not found")
	}

	return err
}

func rowRaybotCommandToModel(row db.RaybotCommand) (model.RaybotCommand, error) {
	input, err := unmarshalCommandInput(row.Input, model.RaybotCommandType(row.Type))
	if err != nil {
		return model.RaybotCommand{}, err
	}

	m := model.RaybotCommand{
		ID:        row.ID,
		RaybotID:  row.RaybotID,
		Status:    model.RaybotCommandStatus(row.Status),
		Type:      model.RaybotCommandType(row.Type),
		Input:     input,
		Output:    row.Output,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.CompletedAt.Valid {
		m.CompletedAt = &row.CompletedAt.Time
	}

	return m, nil
}

func unmarshalCommandInput(input []byte, t model.RaybotCommandType) (any, error) {
	switch t {
	case model.RaybotCommandTypeMoveToLocation:
		var v model.MoveToLocationInput
		if err := json.Unmarshal(input, &v); err != nil {
			return nil, err
		}
		return v, nil
	case model.RaybotCommandTypeCheckQrCode:
		var v model.CheckQRCodeInput
		if err := json.Unmarshal(input, &v); err != nil {
			return nil, err
		}
		return v, nil
	case model.RaybotCommandTypeLiftBox, model.RaybotCommandTypeDropBox:
		var v model.LiftDropBoxInput
		if err := json.Unmarshal(input, &v); err != nil {
			return nil, err
		}
		return v, nil
	case model.RaybotCommandTypeStop,
		model.RaybotCommandTypeMoveForward,
		model.RaybotCommandTypeMoveBackward,
		model.RaybotCommandTypeOpenBox,
		model.RaybotCommandTypeCloseBox,
		model.RaybotCommandTypeWaitGetItem:
		return nil, nil
	default:
		return nil, xerrors.ThrowInvalidArgument(nil, "invalid command type")
	}
}
