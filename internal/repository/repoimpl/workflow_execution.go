package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ repository.WorkflowExecutionRepository = (*workflowExecutionRepository)(nil)

	ErrWorkflowExecutionNotFound = xerror.NotFound(nil, "workflowExecution.notFound", "workflow execution not found")
)

type workflowExecutionRepository struct {
	queries sqlcpg.Queries
}

func newWorkflowExecutionRepository(queries sqlcpg.Queries) *workflowExecutionRepository {
	return &workflowExecutionRepository{queries: queries}
}

func (r workflowExecutionRepository) GetWorkflowExecution(ctx context.Context, db sqldb.SQLDB, id string) (workflowexecution.WorkflowExecution, error) {
	row, err := r.queries.WorkflowExecutionGetByID(ctx, db, id)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return workflowexecution.WorkflowExecution{}, ErrWorkflowExecutionNotFound
		}
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("queries get workflow execution by id: %w", err)
	}

	ret, err := workflowExecutionRowToModel(row)
	if err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("convert workflow execution row to model: %w", err)
	}

	return ret, nil
}

func (r workflowExecutionRepository) ListWorkflowExecutionsByWorkflowID(
	ctx context.Context,
	db sqldb.SQLDB,
	pagingParams paging.Params,
	sorts []sort.Sort,
	workflowID string,
) (paging.List[workflowexecution.WorkflowExecution], error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("workflow_executions").
		Limit(uint64(pagingParams.Limit())).
		Offset(uint64(pagingParams.Offset())).
		Where(sq.Eq{"workflow_id": workflowID})

	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("queries list workflow executions: %w", err)
	}
	defer rows.Close()

	items := make([]workflowexecution.WorkflowExecution, 0, pagingParams.Limit())
	for rows.Next() {
		var i sqlcpg.WorkflowExecution
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowID,
			&i.Status,
			&i.Data,
			&i.Inputs,
			&i.Outputs,
			&i.Error,
			&i.CreatedAt,
			&i.StartedAt,
			&i.CompletedAt,
		); err != nil {
			return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("scan workflow execution: %w", err)
		}

		item, err := workflowExecutionRowToModel(i)
		if err != nil {
			return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("convert workflow execution row to model: %w", err)
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("rows error: %w", err)
	}

	countQuery := psql.Select("COUNT(*)").
		From("workflow_executions").
		Where(sq.Eq{"workflow_id": workflowID})

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("queries count workflow executions: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r workflowExecutionRepository) CreateWorkflowExecution(ctx context.Context, db sqldb.SQLDB, workflowExecution workflowexecution.WorkflowExecution) error {
	data, err := json.Marshal(workflowExecution.Data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	inputs, err := json.Marshal(workflowExecution.Inputs)
	if err != nil {
		return fmt.Errorf("marshal inputs: %w", err)
	}

	outputs, err := json.Marshal(workflowExecution.Outputs)
	if err != nil {
		return fmt.Errorf("marshal outputs: %w", err)
	}

	err = r.queries.WorkflowExecutionInsert(ctx, db, sqlcpg.WorkflowExecutionInsertParams{
		ID:          workflowExecution.ID,
		WorkflowID:  workflowExecution.WorkflowID,
		Status:      string(workflowExecution.Status),
		Data:        data,
		Inputs:      inputs,
		Outputs:     outputs,
		Error:       workflowExecution.Error,
		CreatedAt:   workflowExecution.CreatedAt,
		StartedAt:   workflowExecution.StartedAt,
		CompletedAt: workflowExecution.CompletedAt,
	})
	if err != nil {
		return fmt.Errorf("queries create workflow execution: %w", err)
	}

	return nil
}

func (r workflowExecutionRepository) UpdateWorkflowExecution(ctx context.Context, db sqldb.SQLDB, params repository.UpdateWorkflowExecutionParams) (workflowexecution.WorkflowExecution, error) {
	inputs, err := json.Marshal(params.Inputs)
	if err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("marshal inputs: %w", err)
	}

	outputs, err := json.Marshal(params.Outputs)
	if err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("marshal outputs: %w", err)
	}

	row, err := r.queries.WorkflowExecutionUpdate(ctx, db, sqlcpg.WorkflowExecutionUpdateParams{
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
		if sqldb.IsNoRowsError(err) {
			return workflowexecution.WorkflowExecution{}, ErrWorkflowExecutionNotFound
		}
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("queries update workflow execution: %w", err)
	}

	ret, err := workflowExecutionRowToModel(row)
	if err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("convert workflow execution row to model: %w", err)
	}

	return ret, nil
}

func workflowExecutionRowToModel(row sqlcpg.WorkflowExecution) (workflowexecution.WorkflowExecution, error) {
	data := workflow.Data{}
	if err := json.Unmarshal(row.Data, &data); err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("unmarshal data: %w", err)
	}

	inputs := map[string]any{}
	if err := json.Unmarshal(row.Inputs, &inputs); err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("unmarshal inputs: %w", err)
	}

	outputs := map[string]any{}
	if err := json.Unmarshal(row.Outputs, &outputs); err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("unmarshal outputs: %w", err)
	}

	return workflowexecution.WorkflowExecution{
		ID:          row.ID,
		WorkflowID:  row.WorkflowID,
		Status:      workflowexecution.Status(row.Status),
		Data:        data,
		Inputs:      inputs,
		Outputs:     outputs,
		Error:       row.Error,
		StartedAt:   row.StartedAt,
		CompletedAt: row.CompletedAt,
	}, nil
}
