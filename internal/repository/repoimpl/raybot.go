package repoimpl

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ repository.RaybotRepository = (*raybotRepository)(nil)

	raybotNameConstraint = "name"

	ErrRaybotNotFound    = xerror.NotFound(nil, "raybot.notFound", "raybot not found")
	ErrNameAlreadyExists = xerror.Conflict(nil, "raybot.nameAlreadyExists", "name already exists")
)

type raybotRepository struct {
	queries sqlcpg.Queries
}

func newRaybotRepository(queries sqlcpg.Queries) *raybotRepository {
	return &raybotRepository{queries: queries}
}

func (r *raybotRepository) GetRaybot(ctx context.Context, db sqldb.SQLDB, id string) (raybot.Raybot, error) {
	row, err := r.queries.RaybotGetByID(ctx, db, id)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return raybot.Raybot{}, ErrRaybotNotFound
		}
		if sqldb.IsUniqueViolationError(err, raybotNameConstraint) {
			return raybot.Raybot{}, ErrNameAlreadyExists
		}
		return raybot.Raybot{}, fmt.Errorf("queries get raybot by id: %w", err)
	}

	return raybotRowToModel(row), nil
}

func (r *raybotRepository) ListRaybots(
	ctx context.Context,
	db sqldb.SQLDB,
	pagingParams paging.Params,
	sorts []sort.Sort,
	isOnline *bool,
	controlMode *raybot.ControlMode,
) (paging.List[raybot.Raybot], error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("raybots").
		Limit(uint64(pagingParams.Limit())).
		Offset(uint64(pagingParams.Offset()))

	if isOnline != nil {
		query = query.Where("is_online = ?", *isOnline)
	}
	if controlMode != nil {
		query = query.Where("control_mode = ?", string(*controlMode))
	}
	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("queries list raybots: %w", err)
	}
	defer rows.Close()

	items := make([]raybot.Raybot, 0, pagingParams.Limit())
	for rows.Next() {
		var row sqlcpg.Raybot
		if err := rows.Scan(&row.ID, &row.Name, &row.ControlMode, &row.IsOnline, &row.IpAddress, &row.LastConnectedAt, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return paging.List[raybot.Raybot]{}, fmt.Errorf("scan raybot: %w", err)
		}

		items = append(items, raybotRowToModel(row))
	}

	if err := rows.Err(); err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("rows error: %w", err)
	}

	countQuery := psql.Select("COUNT(*)").From("raybots")
	if isOnline != nil {
		countQuery = countQuery.Where("is_online = ?", *isOnline)
	}
	if controlMode != nil {
		countQuery = countQuery.Where("control_mode = ?", string(*controlMode))
	}

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[raybot.Raybot]{}, fmt.Errorf("queries count raybots: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r *raybotRepository) CreateRaybot(ctx context.Context, db sqldb.SQLDB, raybot raybot.Raybot) error {
	if err := r.queries.RaybotInsert(ctx, db, sqlcpg.RaybotInsertParams{
		ID:              raybot.ID,
		Name:            raybot.Name,
		ControlMode:     string(raybot.ControlMode),
		IsOnline:        raybot.IsOnline,
		LastConnectedAt: raybot.LastConnectedAt,
		IpAddress:       raybot.IPAddress,
		CreatedAt:       raybot.CreatedAt,
		UpdatedAt:       raybot.UpdatedAt,
	}); err != nil {
		if sqldb.IsUniqueViolationError(err, raybotNameConstraint) {
			return ErrNameAlreadyExists
		}
		return fmt.Errorf("queries create raybot: %w", err)
	}

	return nil
}

func (r *raybotRepository) UpdateRaybot(ctx context.Context, db sqldb.SQLDB, params repository.UpdateRaybotParams) (raybot.Raybot, error) {
	row, err := r.queries.RaybotUpdate(ctx, db, sqlcpg.RaybotUpdateParams{
		ID:                 params.ID,
		Name:               params.Name,
		SetName:            params.SetName,
		ControlMode:        string(params.ControlMode),
		SetControlMode:     params.SetControlMode,
		IsOnline:           params.IsOnline,
		SetIsOnline:        params.SetIsOnline,
		IpAddress:          params.IPAddress,
		SetIpAddress:       params.SetIPAddress,
		LastConnectedAt:    params.LastConnectedAt,
		SetLastConnectedAt: params.SetLastConnectedAt,
	})
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return raybot.Raybot{}, ErrRaybotNotFound
		}
		if sqldb.IsUniqueViolationError(err, raybotNameConstraint) {
			return raybot.Raybot{}, ErrNameAlreadyExists
		}
		return raybot.Raybot{}, fmt.Errorf("queries update raybot: %w", err)
	}

	return raybotRowToModel(row), nil
}

func (r *raybotRepository) DeleteRaybot(ctx context.Context, db sqldb.SQLDB, id string) error {
	if err := r.queries.RaybotDelete(ctx, db, id); err != nil {
		if sqldb.IsNoRowsError(err) {
			return ErrRaybotNotFound
		}
		return fmt.Errorf("queries delete raybot: %w", err)
	}

	return nil
}

func raybotRowToModel(row sqlcpg.Raybot) raybot.Raybot {
	return raybot.Raybot{
		ID:              row.ID,
		Name:            row.Name,
		ControlMode:     raybot.ControlMode(row.ControlMode),
		IsOnline:        row.IsOnline,
		IPAddress:       row.IpAddress,
		LastConnectedAt: row.LastConnectedAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}
