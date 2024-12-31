-- name: GetWorkflowExecution :one
SELECT * FROM workflow_executions WHERE id = $1;

-- name: GetWorkflowExecutionForUpdate :one
SELECT * FROM workflow_executions WHERE id = $1 FOR UPDATE;

-- name: GetWorkflowExecutionStatus :one
SELECT status FROM workflow_executions WHERE id = $1;

-- name: CreateWorkflowExecution :exec
INSERT INTO workflow_executions (
	id,
    workflow_id,
    status,
    env,
    definition,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateWorkflowExecution :exec
UPDATE workflow_executions
SET
	id = $1,
	workflow_id = $2,
	status = $3,
	env = $4,
	definition = $5,
	created_at = $6,
	started_at = $7,
	completed_at = $8
WHERE id = $1;


-- name: DeleteWorkflowExecution :exec
DELETE FROM workflow_executions WHERE id = $1;
