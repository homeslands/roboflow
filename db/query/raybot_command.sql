-- name: GetRaybotCommand :one
SELECT * FROM raybot_commands
WHERE id = $1;

-- name: GetRaybotCommandForUpdate :one
SELECT * FROM raybot_commands
WHERE id = $1 FOR UPDATE;

-- name: GetRaybotCommandStatus :one
SELECT status FROM raybot_commands
WHERE id = $1;

-- name: CreateRaybotCommand :exec
INSERT INTO raybot_commands (
    id,
    raybot_id,
    type,
    status,
    input,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateRaybotCommand :exec
UPDATE raybot_commands
SET
	id = $1,
	raybot_id = $2,
	type = $3,
	status = $4,
	input = $5,
	output = $6,
	created_at = $7,
	completed_at = $8
WHERE id = $1;

-- name: DeleteRaybotCommand :exec
DELETE FROM raybot_commands
WHERE id = $1;
