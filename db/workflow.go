package db

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type ListWorkflowsParams struct {
	Limit  int32
	Offset int32
	Sorts  []xsort.Sort
}

func (s *Store) ListWorkflows(ctx context.Context, arg ListWorkflowsParams) (*paging.List[Workflow], error) {
	query := s.StmtBuilder.Select("*").
		From("workflows").
		Limit(uint64(arg.Limit)).
		Offset(uint64(arg.Offset))

	for _, sort := range arg.Sorts {
		query = sort.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return &paging.List[Workflow]{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return &paging.List[Workflow]{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := []Workflow{}
	for rows.Next() {
		var i Workflow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Definition,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return &paging.List[Workflow]{}, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}

	totalSql, totalArgs, err := s.StmtBuilder.Select("COUNT(*)").
		From("workflows").
		ToSql()
	if err != nil {
		return &paging.List[Workflow]{}, fmt.Errorf("failed to build total query: %w", err)
	}

	var totalItem int64
	if err := s.db.QueryRow(ctx, totalSql, totalArgs...).Scan(&totalItem); err != nil {
		return &paging.List[Workflow]{}, fmt.Errorf("failed to get total item: %w", err)
	}

	return paging.NewList(items, totalItem), nil
}
