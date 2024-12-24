-- name: GetWorkflowExecution :one
SELECT * FROM workflow_executions WHERE id = $1;

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

-- name: UpdateWorkflowExecution :one
UPDATE workflow_executions
SET	status = $2,
	started_at = $3,
	completed_at = $4
WHERE id = $1
RETURNING *;


-- name: DeleteWorkflowExecution :exec
DELETE FROM workflow_executions WHERE id = $1;
