package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

type ListStepsParams struct {
	WorkflowExecutionID uuid.UUID
	Sorts               []xsort.Sort
}

func (s *Store) ListSteps(ctx context.Context, arg ListStepsParams) ([]Step, error) {
	query := s.StmtBuilder.Select("*").
		From("steps").
		Where("workflow_execution_id = ?", arg.WorkflowExecutionID.String())

	for _, sort := range arg.Sorts {
		query = sort.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []Step{}, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return []Step{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := []Step{}
	for rows.Next() {
		var i Step
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowExecutionID,
			&i.Env,
			&i.Node,
			&i.Status,
			&i.StartedAt,
			&i.CompletedAt,
		); err != nil {
			return []Step{}, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}

	return items, nil
}
