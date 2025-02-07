-- name: RaybotCommandGetByID :one
SELECT * FROM raybot_commands
WHERE id = @id;

-- name: RaybotCommandStatusGetByID :one
SELECT status FROM raybot_commands
WHERE id = @id;

-- name: RaybotCommandInsert :exec
INSERT INTO raybot_commands (
    id,
    raybot_id,
    type,
    status,
    inputs,
	outputs,
	error,
	completed_at,
	created_at,
	updated_at
)
VALUES (
    @id,
    @raybot_id,
    @type,
    @status,
    @inputs,
    @outputs,
    @error,
    @completed_at,
    @created_at,
    @updated_at
);

-- name: RaybotCommandUpdate :one
UPDATE raybot_commands
SET
	status = CASE WHEN @set_status::boolean THEN @status ELSE status END,
	outputs = CASE WHEN @set_outputs::boolean THEN @outputs ELSE outputs END,
	error = CASE WHEN @set_error::boolean THEN @error ELSE error END,
	completed_at = CASE WHEN @set_completed_at::boolean THEN @completed_at ELSE completed_at END,
	updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: RaybotCommandDeleteByRaybotID :exec
DELETE FROM raybot_commands
WHERE raybot_id = @raybot_id;

-- name: RaybotCommandGetOneByRaybotIDAndStatus :one
SELECT * FROM raybot_commands
WHERE raybot_id = @raybot_id AND status = @status
ORDER BY created_at DESC
LIMIT 1;

-- name: RaybotCommandMarkFailed :exec
UPDATE raybot_commands
SET
    status = 'FAILED',
    error = @error,
    updated_at = NOW()
WHERE raybot_id = @raybot_id
    AND (status = 'PENDING' OR status = 'IN_PROGRESS');
