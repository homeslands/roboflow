-- name: GetStepByID :one
SELECT * FROM steps
WHERE id = $1;

-- name: ListSteps :many
SELECT * FROM steps
WHERE workflow_execution_id = $1;

-- name: BulkInsertSteps :copyfrom
INSERT INTO steps (id, workflow_execution_id, env, node, status)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateStep :exec
UPDATE steps
SET
    status = $2,
    started_at = $3,
    completed_at = $4
WHERE id = $1;
