// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: workflow_execution.sql

package sqlcpg

import (
	"context"
	"encoding/json"
	"time"
)

const workflowExecutionGetByID = `-- name: WorkflowExecutionGetByID :one
SELECT id, workflow_id, status, data, inputs, outputs, error, created_at, updated_at, started_at, completed_at FROM workflow_executions
WHERE id = $1
`

func (q *Queries) WorkflowExecutionGetByID(ctx context.Context, db DBTX, id string) (WorkflowExecution, error) {
	row := db.QueryRow(ctx, workflowExecutionGetByID, id)
	var i WorkflowExecution
	err := row.Scan(
		&i.ID,
		&i.WorkflowID,
		&i.Status,
		&i.Data,
		&i.Inputs,
		&i.Outputs,
		&i.Error,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.StartedAt,
		&i.CompletedAt,
	)
	return i, err
}

const workflowExecutionInsert = `-- name: WorkflowExecutionInsert :exec
INSERT INTO workflow_executions (
	id,
	workflow_id,
	status,
	data,
	inputs,
	outputs,
	error,
	created_at,
	started_at,
	completed_at
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10
)
`

type WorkflowExecutionInsertParams struct {
	ID          string          `json:"id"`
	WorkflowID  string          `json:"workflow_id"`
	Status      string          `json:"status"`
	Data        json.RawMessage `json:"data"`
	Inputs      json.RawMessage `json:"inputs"`
	Outputs     json.RawMessage `json:"outputs"`
	Error       *string         `json:"error"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at"`
}

func (q *Queries) WorkflowExecutionInsert(ctx context.Context, db DBTX, arg WorkflowExecutionInsertParams) error {
	_, err := db.Exec(ctx, workflowExecutionInsert,
		arg.ID,
		arg.WorkflowID,
		arg.Status,
		arg.Data,
		arg.Inputs,
		arg.Outputs,
		arg.Error,
		arg.CreatedAt,
		arg.StartedAt,
		arg.CompletedAt,
	)
	return err
}

const workflowExecutionUpdate = `-- name: WorkflowExecutionUpdate :one
UPDATE workflow_executions
SET
	status = CASE WHEN $1::boolean THEN $2 ELSE status END,
	inputs = CASE WHEN $3::boolean THEN $4 ELSE inputs END,
	outputs = CASE WHEN $5::boolean THEN $6 ELSE outputs END,
	error = CASE WHEN $7::boolean THEN $8 ELSE error END,
	started_at = CASE WHEN $9::boolean THEN $10 ELSE started_at END,
	completed_at = CASE WHEN $11::boolean THEN $12 ELSE completed_at END,
	updated_at = NOW()
WHERE id = $13
RETURNING id, workflow_id, status, data, inputs, outputs, error, created_at, updated_at, started_at, completed_at
`

type WorkflowExecutionUpdateParams struct {
	SetStatus      bool            `json:"set_status"`
	Status         string          `json:"status"`
	SetInputs      bool            `json:"set_inputs"`
	Inputs         json.RawMessage `json:"inputs"`
	SetOutputs     bool            `json:"set_outputs"`
	Outputs        json.RawMessage `json:"outputs"`
	SetError       bool            `json:"set_error"`
	Error          *string         `json:"error"`
	SetStartedAt   bool            `json:"set_started_at"`
	StartedAt      *time.Time      `json:"started_at"`
	SetCompletedAt bool            `json:"set_completed_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	ID             string          `json:"id"`
}

func (q *Queries) WorkflowExecutionUpdate(ctx context.Context, db DBTX, arg WorkflowExecutionUpdateParams) (WorkflowExecution, error) {
	row := db.QueryRow(ctx, workflowExecutionUpdate,
		arg.SetStatus,
		arg.Status,
		arg.SetInputs,
		arg.Inputs,
		arg.SetOutputs,
		arg.Outputs,
		arg.SetError,
		arg.Error,
		arg.SetStartedAt,
		arg.StartedAt,
		arg.SetCompletedAt,
		arg.CompletedAt,
		arg.ID,
	)
	var i WorkflowExecution
	err := row.Scan(
		&i.ID,
		&i.WorkflowID,
		&i.Status,
		&i.Data,
		&i.Inputs,
		&i.Outputs,
		&i.Error,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.StartedAt,
		&i.CompletedAt,
	)
	return i, err
}
