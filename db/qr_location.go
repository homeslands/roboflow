package db

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type ListQrLocationsParams struct {
	Limit  int32
	Offset int32
	Sorts  []xsort.Sort
}

func (s *Store) ListQrLocations(ctx context.Context, arg ListQrLocationsParams) (*paging.List[QrLocation], error) {
	query := s.StmtBuilder.Select("*").
		From("qr_locations").
		Limit(uint64(arg.Limit)).
		Offset(uint64(arg.Offset))

	for _, sort := range arg.Sorts {
		query = sort.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return &paging.List[QrLocation]{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return &paging.List[QrLocation]{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := []QrLocation{}
	for rows.Next() {
		var i QrLocation
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.QrCode,
			&i.Metadata,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return &paging.List[QrLocation]{}, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}

	totalSql, totalArgs, err := s.StmtBuilder.Select("COUNT(*)").
		From("qr_locations").
		ToSql()
	if err != nil {
		return &paging.List[QrLocation]{}, fmt.Errorf("failed to build total query: %w", err)
	}

	var total int64
	if err := s.db.QueryRow(ctx, totalSql, totalArgs...).Scan(&total); err != nil {
		return &paging.List[QrLocation]{}, fmt.Errorf("failed to get total: %w", err)
	}

	return paging.NewList(items, total), nil
}
