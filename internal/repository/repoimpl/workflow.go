package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	workflow "github.com/tuanvumaihuynh/roboflow/internal/model/workflow"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ repository.WorkflowRepository = (*workflowRepository)(nil)

	ErrWorkflowNotFound = xerror.NotFound(nil, "workflow.notFound", "workflow not found")
)

type workflowRepository struct {
	queries sqlcpg.Queries
}

func newWorkflowRepository(queries sqlcpg.Queries) *workflowRepository {
	return &workflowRepository{queries: queries}
}

func (r workflowRepository) GetWorkflow(ctx context.Context, db sqldb.SQLDB, id string) (workflow.Workflow, error) {
	row, err := r.queries.WorkflowGetByID(ctx, db, id)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return workflow.Workflow{}, ErrWorkflowNotFound
		}
		return workflow.Workflow{}, fmt.Errorf("queries get workflow by id: %w", err)
	}

	m, err := workflowRowToModel(row)
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("workflow row to model: %w", err)
	}

	return m, nil
}

func (r workflowRepository) ListWorkflows(ctx context.Context, db sqldb.SQLDB, pagingParams paging.Params, sorts []sort.Sort) (paging.List[workflow.Workflow], error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("workflows").
		Limit(uint64(pagingParams.Limit())).
		Offset(uint64(pagingParams.Offset()))

	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("queries list workflows: %w", err)
	}
	defer rows.Close()

	items := make([]workflow.Workflow, 0, pagingParams.Limit())
	for rows.Next() {
		var i sqlcpg.Workflow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IsDraft,
			&i.IsValid,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return paging.List[workflow.Workflow]{}, fmt.Errorf("scan workflow: %w", err)
		}

		item, err := workflowRowToModel(i)
		if err != nil {
			return paging.List[workflow.Workflow]{}, fmt.Errorf("workflow row to model: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("rows error: %w", err)
	}

	countQuery := psql.Select("COUNT(*)").From("workflows")

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[workflow.Workflow]{}, fmt.Errorf("queries count workflows: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r workflowRepository) CreateWorkflow(ctx context.Context, db sqldb.SQLDB, workflow workflow.Workflow) error {
	data, err := json.Marshal(workflow.Data)
	if err != nil {
		return fmt.Errorf("marshal workflow data: %w", err)
	}

	err = r.queries.WorkflowInsert(ctx, db, sqlcpg.WorkflowInsertParams{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		IsDraft:     workflow.IsDraft,
		IsValid:     workflow.IsValid,
		Data:        data,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("queries create workflow: %w", err)
	}

	return nil
}

func (r workflowRepository) UpdateWorkflow(ctx context.Context, db sqldb.SQLDB, params repository.UpdateWorkflowParams) (workflow.Workflow, error) {
	data, err := json.Marshal(params.Data)
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("marshal workflow data: %w", err)
	}

	row, err := r.queries.WorkflowUpdate(ctx, db, sqlcpg.WorkflowUpdateParams{
		ID:             params.ID,
		Name:           params.Name,
		SetName:        params.SetName,
		Description:    params.Description,
		SetDescription: params.SetDescription,
		IsDraft:        params.IsDraft,
		SetIsDraft:     params.SetIsDraft,
		IsValid:        params.IsValid,
		SetIsValid:     params.SetIsValid,
		Data:           data,
		SetData:        params.SetData,
	})
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return workflow.Workflow{}, ErrWorkflowNotFound
		}
		return workflow.Workflow{}, fmt.Errorf("queries update workflow: %w", err)
	}

	ret, err := workflowRowToModel(row)
	if err != nil {
		return workflow.Workflow{}, fmt.Errorf("workflow row to model: %w", err)
	}

	return ret, nil
}

func (r *workflowRepository) DeleteWorkflow(ctx context.Context, db sqldb.SQLDB, id string) error {
	err := r.queries.WorkflowDelete(ctx, db, id)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return ErrWorkflowNotFound
		}
		return fmt.Errorf("queries delete workflow: %w", err)
	}

	return nil
}

func workflowRowToModel(row sqlcpg.Workflow) (workflow.Workflow, error) {
	var data workflow.Data
	if err := json.Unmarshal(row.Data, &data); err != nil {
		return workflow.Workflow{}, fmt.Errorf("unmarshal workflow data: %w", err)
	}

	return workflow.Workflow{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		IsDraft:     row.IsDraft,
		IsValid:     row.IsValid,
		Data:        data,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}
