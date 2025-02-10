package repoimpl

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ repository.RaybotCommandRepository = (*raybotCommandRepository)(nil)

	ErrRaybotCommandNotFound = xerror.BadRequest(nil, "raybotCommand.notFound", "raybot command not found")
)

type raybotCommandRepository struct {
	queries sqlcpg.Queries
}

func newRaybotCommandRepository(queries sqlcpg.Queries) *raybotCommandRepository {
	return &raybotCommandRepository{queries: queries}
}

func (r raybotCommandRepository) GetRaybotCommand(ctx context.Context, db sqldb.SQLDB, id string) (raybotcommand.RaybotCommand, error) {
	row, err := r.queries.RaybotCommandGetByID(ctx, db, id)
	if err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("queries get raybot command by id: %w", err)
	}

	return raybotCommandRowToModel(row), nil
}

func (r raybotCommandRepository) ListRaybotCommandsByRaybotID(ctx context.Context, db sqldb.SQLDB, raybotID string, pagingParams paging.Params, sorts []sort.Sort) (paging.List[raybotcommand.RaybotCommand], error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("raybot_commands").
		Limit(uint64(pagingParams.Limit())).
		Offset(uint64(pagingParams.Offset())).
		Where(sq.Eq{"raybot_id": raybotID})

	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("queries list raybot commands by raybot id: %w", err)
	}
	defer rows.Close()

	items := make([]raybotcommand.RaybotCommand, 0, pagingParams.Limit())
	for rows.Next() {
		var i sqlcpg.RaybotCommand
		if err := rows.Scan(
			&i.ID,
			&i.RaybotID,
			&i.Type,
			&i.Status,
			&i.Inputs,
			&i.Outputs,
			&i.Error,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("scan raybot command: %w", err)
		}

		items = append(items, raybotCommandRowToModel(i))
	}
	if err := rows.Err(); err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("rows error: %w", err)
	}

	countQuery := psql.Select("COUNT(*)").From("raybot_commands").Where(sq.Eq{"raybot_id": raybotID})

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("queries count raybot commands: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r raybotCommandRepository) CreateRaybotCommand(ctx context.Context, db sqldb.SQLDB, raybotCommand raybotcommand.RaybotCommand) error {
	err := r.queries.RaybotCommandInsert(ctx, db, sqlcpg.RaybotCommandInsertParams{
		ID:          raybotCommand.ID,
		RaybotID:    raybotCommand.RaybotID,
		Type:        string(raybotCommand.Type),
		Status:      string(raybotCommand.Status),
		Inputs:      raybotCommand.Inputs.Raw(),
		Outputs:     raybotCommand.Outputs.Raw(),
		Error:       raybotCommand.Error,
		CompletedAt: raybotCommand.CompletedAt,
		CreatedAt:   raybotCommand.CreatedAt,
		UpdatedAt:   raybotCommand.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("queries create raybot command: %w", err)
	}

	return nil
}

func (r raybotCommandRepository) UpdateRaybotCommand(ctx context.Context, db sqldb.SQLDB, params repository.UpdateRaybotCommandParams) (raybotcommand.RaybotCommand, error) {
	row, err := r.queries.RaybotCommandUpdate(ctx, db, sqlcpg.RaybotCommandUpdateParams{
		ID:             params.ID,
		Status:         string(params.Status),
		SetStatus:      params.SetStatus,
		Outputs:        params.Outputs.Raw(),
		SetOutputs:     params.SetOutputs,
		Error:          params.Error,
		SetError:       params.SetError,
		CompletedAt:    params.CompletedAt,
		SetCompletedAt: params.SetCompletedAt,
	})
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return raybotcommand.RaybotCommand{}, ErrRaybotCommandNotFound
		}
		return raybotcommand.RaybotCommand{}, fmt.Errorf("queries update raybot command: %w", err)
	}

	return raybotCommandRowToModel(row), nil
}

func (r raybotCommandRepository) DeleteRaybotCommandsByRaybotID(ctx context.Context, db sqldb.SQLDB, raybotID string) error {
	err := r.queries.RaybotCommandDeleteByRaybotID(ctx, db, raybotID)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return ErrRaybotCommandNotFound
		}
		return fmt.Errorf("queries delete raybot commands by raybot id: %w", err)
	}

	return nil
}

func (r raybotCommandRepository) MarkRaybotCommandFailed(ctx context.Context, db sqldb.SQLDB, params repository.MarkRaybotCommandFailedParams) error {
	err := r.queries.RaybotCommandMarkFailed(ctx, db, sqlcpg.RaybotCommandMarkFailedParams{
		RaybotID: params.RaybotID,
		Error:    &params.Error,
	})
	if err != nil {
		return fmt.Errorf("queries mark raybot command failed: %w", err)
	}

	return nil
}

func raybotCommandRowToModel(row sqlcpg.RaybotCommand) raybotcommand.RaybotCommand {
	return raybotcommand.RaybotCommand{
		ID:          row.ID,
		RaybotID:    row.RaybotID,
		Type:        raybotcommand.Type(row.Type),
		Status:      raybotcommand.Status(row.Status),
		Inputs:      raybotcommand.NewInputs(row.Inputs),
		Outputs:     raybotcommand.NewOutputs(row.Outputs),
		Error:       row.Error,
		CompletedAt: row.CompletedAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
