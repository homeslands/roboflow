package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type ListRaybotCommandsParams struct {
	RaybotID uuid.UUID
	Limit    int32
	Offset   int32
	Sorts    []xsort.Sort
}

func (s *Store) ListRaybotCommands(ctx context.Context, arg ListRaybotCommandsParams) (*paging.List[RaybotCommand], error) {
	query := s.StmtBuilder.Select("*").
		From("raybot_commands").
		Where(sq.Eq{"raybot_id": arg.RaybotID}).
		Limit(uint64(arg.Limit)).
		Offset(uint64(arg.Offset))

	for _, sort := range arg.Sorts {
		query = sort.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return &paging.List[RaybotCommand]{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return &paging.List[RaybotCommand]{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := []RaybotCommand{}
	for rows.Next() {
		var i RaybotCommand
		if err := rows.Scan(
			&i.ID,
			&i.RaybotID,
			&i.Type,
			&i.Status,
			&i.Input,
			&i.Output,
			&i.CreatedAt,
			&i.CompletedAt,
		); err != nil {
			return &paging.List[RaybotCommand]{}, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}

	totalSql, totalArgs, err := s.StmtBuilder.Select("COUNT(*)").
		From("raybot_commands").
		Where(sq.Eq{"raybot_id": arg.RaybotID}).
		ToSql()
	if err != nil {
		return &paging.List[RaybotCommand]{}, fmt.Errorf("failed to build total query: %w", err)
	}

	var total int64
	if err := s.db.QueryRow(ctx, totalSql, totalArgs...).Scan(&total); err != nil {
		return &paging.List[RaybotCommand]{}, fmt.Errorf("failed to get total: %w", err)
	}

	return paging.NewList(items, total), nil
}
