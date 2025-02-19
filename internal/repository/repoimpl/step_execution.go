package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
)

var (
	_ repository.StepExecutionRepository = (*stepExecutionRepository)(nil)
)

type stepExecutionRepository struct {
	queries sqlcpg.Queries
}

func newStepExecutionRepository(queries sqlcpg.Queries) *stepExecutionRepository {
	return &stepExecutionRepository{queries: queries}
}

func (r stepExecutionRepository) GetStepExecution(ctx context.Context, db sqldb.SQLDB, id string) (stepexecution.StepExecution, error) {
	row, err := r.queries.StepExecutionGet(ctx, db, id)
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("queries get step execution: %w", err)
	}

	return stepExecutionRowToModel(row)
}

func (r stepExecutionRepository) ListStepsByWorkflowExecutionID(ctx context.Context, db sqldb.SQLDB, workflowExecutionID string) ([]stepexecution.StepExecution, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("step_executions").
		Where(sq.Eq{"workflow_execution_id": workflowExecutionID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("queries list steps by workflow execution id: %w", err)
	}
	defer rows.Close()

	items := make([]stepexecution.StepExecution, 0)
	for rows.Next() {
		var i sqlcpg.StepExecution
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowExecutionID,
			&i.Status,
			&i.Inputs,
			&i.Outputs,
			&i.Error,
			&i.StartedAt,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan step execution: %w", err)
		}

		item, err := stepExecutionRowToModel(i)
		if err != nil {
			return nil, fmt.Errorf("convert step execution row to model: %w", err)
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return items, nil
}

func (r stepExecutionRepository) BatchCreateStepExecutions(ctx context.Context, db sqldb.SQLDB, steps []stepexecution.StepExecution) error {
	arg := make([]sqlcpg.StepExecutionBatchInsertParams, 0, len(steps))
	for _, step := range steps {
		inputs, err := json.Marshal(step.Inputs)
		if err != nil {
			return fmt.Errorf("marshal inputs: %w", err)
		}

		outputs, err := json.Marshal(step.Outputs)
		if err != nil {
			return fmt.Errorf("marshal outputs: %w", err)
		}

		node, err := json.Marshal(step.Node)
		if err != nil {
			return fmt.Errorf("marshal node: %w", err)
		}

		arg = append(arg, sqlcpg.StepExecutionBatchInsertParams{
			ID:                  step.ID,
			WorkflowExecutionID: step.WorkflowExecutionID,
			Status:              string(step.Status),
			Node:                node,
			Inputs:              inputs,
			Outputs:             outputs,
			Error:               step.Error,
			CreatedAt:           step.CreatedAt,
			UpdatedAt:           step.UpdatedAt,
			StartedAt:           step.StartedAt,
			CompletedAt:         step.CompletedAt,
		})
	}

	count, err := r.queries.StepExecutionBatchInsert(ctx, db, arg)
	if err != nil {
		return fmt.Errorf("queries batch create steps: %w", err)
	}
	if count != int64(len(steps)) {
		return fmt.Errorf("batch create steps: %w", err)
	}

	return nil
}

func (r stepExecutionRepository) UpdateStepExecution(ctx context.Context, db sqldb.SQLDB, params repository.UpdateStepExecutionParams) (stepexecution.StepExecution, error) {
	inputs, err := json.Marshal(params.Inputs)
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("marshal inputs: %w", err)
	}

	outputs, err := json.Marshal(params.Outputs)
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("marshal outputs: %w", err)
	}

	row, err := r.queries.StepExecutionUpdate(ctx, db, sqlcpg.StepExecutionUpdateParams{
		ID:             params.ID,
		Status:         string(params.Status),
		SetStatus:      params.SetStatus,
		Inputs:         inputs,
		SetInputs:      params.SetInputs,
		Outputs:        outputs,
		SetOutputs:     params.SetOutputs,
		Error:          params.Error,
		SetError:       params.SetError,
		StartedAt:      params.StartedAt,
		SetStartedAt:   params.SetStartedAt,
		CompletedAt:    params.CompletedAt,
		SetCompletedAt: params.SetCompletedAt,
	})
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("queries update step: %w", err)
	}

	return stepExecutionRowToModel(row)
}

func stepExecutionRowToModel(row sqlcpg.StepExecution) (stepexecution.StepExecution, error) {
	inputs := map[string]any{}
	if err := json.Unmarshal(row.Inputs, &inputs); err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("unmarshal inputs: %w", err)
	}

	outputs := map[string]any{}
	if err := json.Unmarshal(row.Outputs, &outputs); err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("unmarshal outputs: %w", err)
	}

	node := node.Node{}
	if err := json.Unmarshal(row.Node, &node); err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("unmarshal node: %w", err)
	}

	return stepexecution.StepExecution{
		ID:                  row.ID,
		WorkflowExecutionID: row.WorkflowExecutionID,
		Status:              stepexecution.Status(row.Status),
		Node:                node,
		Inputs:              inputs,
		Outputs:             outputs,
		Error:               row.Error,
		StartedAt:           row.StartedAt,
		CompletedAt:         row.CompletedAt,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}, nil

}
