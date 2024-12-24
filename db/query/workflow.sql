-- name: GetWorkflow :one
SELECT * FROM workflows
WHERE id = $1;

-- name: CreateWorkflow :exec
INSERT INTO workflows (
    id,
    name,
    description,
    definition,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateWorkflow :one
UPDATE workflows
SET
    name = $1,
    description = $2,
    definition = $3,
    updated_at = $4
WHERE id = $5
RETURNING *;

-- name: DeleteWorkflow :exec
DELETE FROM workflows
WHERE id = $1;
