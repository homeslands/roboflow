package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type ListRaybotsParams struct {
	Limit  int32
	Offset int32
	Sorts  []xsort.Sort
	Status *string
}

func (s *Store) ListRaybots(ctx context.Context, arg ListRaybotsParams) (*paging.List[Raybot], error) {
	query := s.StmtBuilder.Select("*").
		From("raybots").
		Limit(uint64(arg.Limit)).
		Offset(uint64(arg.Offset))

	for _, sort := range arg.Sorts {
		query = sort.Attach(query)
	}
	if arg.Status != nil && *arg.Status != "" {
		query = query.Where(sq.Eq{"status": arg.Status})
	}

	raw, args, err := query.ToSql()
	if err != nil {
		return &paging.List[Raybot]{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, raw, args...)
	if err != nil {
		return &paging.List[Raybot]{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := []Raybot{}
	for rows.Next() {
		var i Raybot
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Token,
			&i.Status,
			&i.IpAddress,
			&i.LastConnectedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return &paging.List[Raybot]{}, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}

	query = s.StmtBuilder.Select("COUNT(*)").
		From("raybots")

	if arg.Status != nil && *arg.Status != "" {
		query = query.Where(sq.Eq{"status": arg.Status})
	}
	raw, countArgs, err := query.ToSql()
	if err != nil {
		return &paging.List[Raybot]{}, fmt.Errorf("failed to build total query: %w", err)
	}

	var total int64
	if err := s.db.QueryRow(ctx, raw, countArgs...).Scan(&total); err != nil {
		return &paging.List[Raybot]{}, fmt.Errorf("failed to scan total: %w", err)
	}

	return paging.NewList(items, total), nil
}
